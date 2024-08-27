package lexer

import "monkey/token"

// doubleTokens are tokens that have exactly two characters
// and start with a non, digit, non character symbol
var doubleTokens = map[string]token.Token{
	"!=": {
		Type:    token.NOT_EQ,
		Literal: "!=",
	},
	"==": {
		Type:    token.EQ,
		Literal: "==",
	},
}

// unitTokens are tokens that are exactly one character long
// and whose lexical representation is the same as the string on the token type
var unitTokens = map[byte]token.Token{
	'=': {
		Type:    token.ASSIGN,
		Literal: "=",
	},
	'+': {
		Type:    token.PLUS,
		Literal: "+",
	},
	'-': {
		Type:    token.MINUS,
		Literal: "-",
	},
	'!': {
		Type:    token.BANG,
		Literal: "!",
	},
	'*': {
		Type:    token.ASTERISK,
		Literal: "*",
	},
	'/': {
		Type:    token.SLASH,
		Literal: "/",
	},
	',': {
		Type:    token.COMMA,
		Literal: ",",
	},
	';': {
		Type:    token.SEMICOLON,
		Literal: ";",
	},
	'(': {
		Type:    token.LPAREN,
		Literal: "(",
	},
	')': {
		Type:    token.RPAREN,
		Literal: ")",
	},
	'{': {
		Type:    token.LBRACE,
		Literal: "{",
	},
	'}': {
		Type:    token.RBRACE,
		Literal: "}",
	},
	'<': {
		Type:    token.LT,
		Literal: "<",
	},
	'>': {
		Type:    token.GT,
		Literal: ">",
	},
}

// keywordTokens are tokens that are words that are more
// than one character
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

var eofToken = token.Token{
	Type: token.EOF,
}

func identToken(literal string) token.Token {
	return token.Token{
		Type:    token.IDENT,
		Literal: literal,
	}
}

func digitToken(literal string) token.Token {
	return token.Token{
		Type:    token.INT,
		Literal: literal,
	}
}

func illegalToken(literalByte byte) token.Token {
	return token.Token{
		Type:    token.ILLEGAL,
		Literal: string(literalByte),
	}
}
