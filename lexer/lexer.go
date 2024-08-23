package lexer

import (
	"iter"
	"monkey/token"
)

const (
	EofByte = 0
)

func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n' || ch == '\v' || ch == '\f'
}

type lexerState struct {
	toProcess string
}

func newLexerState(toProcess string) *lexerState {
	return &lexerState{
		toProcess: toProcess,
	}
}

func (l *lexerState) finished() bool {
	return l.toProcess == ""
}

func (l *lexerState) next() byte {
	if l.finished() {
		return EofByte
	}
	result := l.toProcess[0]
	l.toProcess = l.toProcess[1:]
	return result
}

func (l *lexerState) peek() byte {
	if l.finished() {
		return EofByte
	}
	return l.toProcess[0]
}

func (l *lexerState) getNextBytes(cond func(byte) bool) (resultBytes []byte) {
	for {
		if !cond(l.peek()) {
			return
		}
		resultBytes = append(resultBytes, l.next())
	}
}

var unitTokens = make(map[string]token.Token)

func init() {
	for _, tokenType := range []token.TokenType{
		token.ASSIGN,
		token.PLUS,
		token.MINUS,
		token.BANG,
		token.ASTERISK,
		token.SLASH,
		token.COMMA,
		token.SEMICOLON,
		token.LPAREN,
		token.RPAREN,
		token.LBRACE,
		token.RBRACE,
		token.LT,
		token.GT,
	} {
		strToken := string(tokenType)
		unitTokens[strToken] = token.Token{
			Type:    tokenType,
			Literal: strToken,
		}
	}
}

var keywordTokens = map[string]token.Token{
	"let": {
		Type:    token.LET,
		Literal: "let",
	},
	"fn": {
		Type:    token.FUNCTION,
		Literal: "fn",
	},
	"true": {
		Type:    token.TRUE,
		Literal: "true",
	},
	"false": {
		Type:    token.FALSE,
		Literal: "false",
	},
	"if": {
		Type:    token.IF,
		Literal: "if",
	},
	"else": {
		Type:    token.ELSE,
		Literal: "else",
	},
	"return": {
		Type:    token.RETURN,
		Literal: "return",
	},
}

func consumeWhitespace(state *lexerState) {
	_ = state.getNextBytes(isWhitespace)
	return
}

func getDouble(ch byte, state *lexerState) (tok token.Token, ok bool) {
	nextChar := state.peek()
	if ch == '=' && nextChar == '=' {
		tok = token.Token{Type: token.EQ, Literal: "=="}
		ok = true
	} else if ch == '!' && nextChar == '=' {
		tok = token.Token{Type: token.NOT_EQ, Literal: "!="}
		ok = true
	}
	if ok {
		_ = state.next()
	}
	return
}

func getStringToken(ch byte, state *lexerState) token.Token {
	nextStringBytes := state.getNextBytes(isLetter)

	literal := string(append([]byte{ch}, nextStringBytes...))

	var tok token.Token
	if keywordToken, ok := keywordTokens[literal]; ok {
		tok = keywordToken
	} else {
		tok = token.Token{Type: token.IDENT, Literal: literal}
	}
	return tok
}

func getDigitToken(ch byte, state *lexerState) token.Token {
	nextDigitBytes := state.getNextBytes(isDigit)
	literal := string(append([]byte{ch}, nextDigitBytes...))
	return token.Token{Type: token.INT, Literal: literal}
}

func nextToken(state *lexerState) token.Token {
	consumeWhitespace(state)
	ch := state.next()

	var tok token.Token
	if ch == EofByte {
		tok = token.Token{Type: token.EOF}
	} else if doubleToken, ok := getDouble(ch, state); ok {
		tok = doubleToken
	} else if unitToken, ok := unitTokens[string(ch)]; ok {
		tok = unitToken
	} else if isLetter(ch) {
		tok = getStringToken(ch, state)
	} else if isDigit(ch) {
		tok = getDigitToken(ch, state)
	} else {
		tok = token.Token{Type: token.ILLEGAL, Literal: string(ch)}
	}
	return tok
}

func Lex(input string) iter.Seq[token.Token] {
	state := newLexerState(input)

	return func(yield func(token.Token) bool) {
		var tok token.Token
		for {
			tok = nextToken(state)

			if !yield(tok) || tok.Type == token.EOF || tok.Type == token.ILLEGAL {
				break
			}
		}
	}

}
