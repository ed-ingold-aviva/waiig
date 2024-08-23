package lexer

import "monkey/token"

// unitTokens are tokens that are exactly one character long
// and whose lexical representation is the same as the string on the token type
var unitTokens = make(map[byte]token.Token)

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
		if len(tokenType) != 1 {
			panic("Unit tokens must have a single character")
		}
		unitTokens[tokenType[0]] = token.Token{
			Type:    tokenType,
			Literal: string(tokenType),
		}
	}
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

var eqToken = token.Token{
	Type:    token.EQ,
	Literal: "==",
}

var neqToken = token.Token{
	Type:    token.NOT_EQ,
	Literal: "!=",
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
