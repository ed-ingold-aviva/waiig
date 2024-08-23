package lexer

import (
	"iter"
	"monkey/token"
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

type lexerState string

func (l lexerState) finished() bool {
	return l == ""
}

func (l lexerState) next() (byte, lexerState) {
	if l.finished() {
		return 0, l
	}
	return l[0], l[1:]
}

func (l lexerState) peek() byte {
	if l.finished() {
		return 0
	}
	return l[0]
}

func (l lexerState) getNextBytes(cond func(byte) bool) (resultBytes []byte, resultState lexerState) {
	resultState = l
	for {
		if !cond(resultState.peek()) {
			return
		}
		var resultByte byte
		resultByte, resultState = resultState.next()
		resultBytes = append(resultBytes, resultByte)
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

func consumeWhitespace(state lexerState) lexerState {
	_, newState := state.getNextBytes(isWhitespace)
	return newState
}

func getDouble(ch byte, state lexerState) (tok token.Token, newState lexerState, ok bool) {
	nextChar := state.peek()
	if ch == '=' && nextChar == '=' {
		tok = token.Token{Type: token.EQ, Literal: "=="}
		ok = true
	} else if ch == '!' && nextChar == '=' {
		tok = token.Token{Type: token.NOT_EQ, Literal: "!="}
		ok = true
	}
	if ok {
		_, newState = state.next()
	}

	return
}

func getUnitToken(ch byte) (token token.Token, ok bool) {
	token, ok = unitTokens[string(ch)]
	return
}

func getStringToken(ch byte, state lexerState) (token.Token, lexerState) {
	nextStringBytes, nextState := state.getNextBytes(isLetter)

	literal := string(append([]byte{ch}, nextStringBytes...))

	var tok token.Token
	if keywordToken, ok := keywordTokens[literal]; ok {
		tok = keywordToken
	} else {
		tok.Type = token.IDENT
		tok.Literal = literal
	}
	return tok, nextState
}

func getDigitToken(ch byte, state lexerState) (token.Token, lexerState) {
	nextDigitBytes, nextState := state.getNextBytes(isDigit)
	return token.Token{Type: token.INT, Literal: string(append([]byte{ch}, nextDigitBytes...))}, nextState
}

func nextToken(state lexerState) (token.Token, lexerState) {
	var ch byte
	var tok token.Token

	ch, state = consumeWhitespace(state).next()

	if ch == 0 {
		tok = token.Token{
			Type:    token.EOF,
			Literal: "",
		}
	} else if doubleToken, newState, ok := getDouble(ch, state); ok {
		tok = doubleToken
		state = newState
	} else if unitToken, ok := getUnitToken(ch); ok {
		tok = unitToken
	} else if isLetter(ch) {
		tok, state = getStringToken(ch, state)
	} else if isDigit(ch) {
		tok, state = getDigitToken(ch, state)
	} else {
		tok.Type = token.ILLEGAL
		tok.Literal = string(ch)
	}
	return tok, state
}

func Lex(input string) iter.Seq[token.Token] {
	state := lexerState(input)

	return func(yield func(token.Token) bool) {
		var tok token.Token
		for {
			tok, state = nextToken(state)

			if !yield(tok) || tok.Type == token.EOF || tok.Type == token.ILLEGAL {
				break
			}
		}
	}

}
