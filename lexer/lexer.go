package lexer

import (
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

type lexerState struct {
	input    string
	position int
}

func (l lexerState) finished() bool {
	return l.position >= len(l.input)
}

func (l lexerState) next() (byte, lexerState) {
	if l.finished() {
		return 0, l
	}
	next := l.input[l.position]
	nextState := lexerState{
		input:    l.input,
		position: l.position + 1,
	}
	return next, nextState
}

func (l lexerState) peek() byte {
	if l.finished() {
		return 0
	}
	return l.input[l.position]
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

func (l lexerState) consumeWhitespace() lexerState {
	for {
		if !isWhitespace(l.peek()) {
			return l
		}
		_, l = l.next()
	}
}

var eqToken = token.Token{
	Type:    token.EQ,
	Literal: "==",
}

var neqToken = token.Token{
	Type:    token.NOT_EQ,
	Literal: "!=",
}

var eofToken = token.Token{
	Type:    token.EOF,
	Literal: "",
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

func getDouble(ch byte, state lexerState) (token token.Token, newState lexerState, ok bool) {
	nextChar := state.peek()
	if ch == '=' && nextChar == '=' {
		token = eqToken
		ok = true
	} else if ch == '!' && nextChar == '=' {
		token = neqToken
		ok = true
	}
	if ok {
		_, newState = state.next()
	}

	return
}

func getStringToken(ch byte, state lexerState) (token.Token, lexerState) {
	nextStringBytes, nextState := state.getNextBytes(isLetter)

	literal := string(append([]byte{ch}, nextStringBytes...))

	if keywordToken, ok := keywordTokens[literal]; ok {
		return keywordToken, nextState
	}
	return token.Token{Type: token.IDENT, Literal: literal}, nextState
}

func getDigitToken(ch byte, state lexerState) (token.Token, lexerState) {
	nextDigitBytes, nextState := state.getNextBytes(isDigit)
	return token.Token{Type: token.INT, Literal: string(append([]byte{ch}, nextDigitBytes...))}, nextState
}

type Lexer struct {
	state lexerState
}

func (l *Lexer) NextToken() (tok token.Token) {
	var ch byte

	l.state = l.state.consumeWhitespace()
	ch, l.state = l.state.next()

	if ch == 0 {
		tok = eofToken
	} else if doubleToken, newState, ok := getDouble(ch, l.state); ok {
		tok = doubleToken
		l.state = newState
	} else if unitToken, ok := unitTokens[string(ch)]; ok {
		tok = unitToken
	} else if isLetter(ch) {
		tok, l.state = getStringToken(ch, l.state)
	} else if isDigit(ch) {
		tok, l.state = getDigitToken(ch, l.state)
	} else {
		tok = token.Token{
			Type:    token.ILLEGAL,
			Literal: string(ch),
		}
	}
	return
}

func New(input string) *Lexer {
	return &Lexer{state: lexerState{
		input: input,
	}}
}
