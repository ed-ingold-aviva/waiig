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

// getStringForCond consumes the lexerState for character's that pass the
// required condition. If a character fails the cond, it is not consumed.
//
// The resultant bytes are converted to a string.
//
// The init parameter allows you to initialise the byte array that is built up
// so those bytes will end up at the start of the returned string.
func (l *lexerState) getStringForCond(cond func(byte) bool, init ...byte) string {
	for {
		if !cond(l.peek()) {
			return string(init)
		}
		init = append(init, l.next())
	}
}

func consumeWhitespace(state *lexerState) {
	_ = state.getStringForCond(isWhitespace)
	return
}

func getDoubleToken(ch byte, state *lexerState) (tok token.Token, tokenFound bool) {
	double := string([]byte{ch, state.peek()})

	if doubleToken, ok := doubleTokens[double]; ok {
		tok = doubleToken
		tokenFound = true
		// need to consume second character as well
		_ = state.next()
	}
	return
}

func getStringToken(ch byte, state *lexerState) token.Token {
	literal := state.getStringForCond(isLetter, ch)

	var tok token.Token
	if keywordToken, ok := keywordTokens[literal]; ok {
		tok = keywordToken
	} else {
		tok = identToken(literal)
	}
	return tok
}

func getDigitToken(ch byte, state *lexerState) token.Token {
	return digitToken(state.getStringForCond(isDigit, ch))
}

func nextToken(state *lexerState) token.Token {
	consumeWhitespace(state)
	ch := state.next()

	var tok token.Token
	if ch == EofByte {
		tok = eofToken
	} else if doubleToken, ok := getDoubleToken(ch, state); ok {
		// Need to process the double tokens before unit ones
		tok = doubleToken
	} else if unitToken, ok := unitTokens[ch]; ok {
		tok = unitToken
	} else if isLetter(ch) {
		tok = getStringToken(ch, state)
	} else if isDigit(ch) {
		tok = getDigitToken(ch, state)
	} else {
		tok = illegalToken(ch)
	}
	return tok
}

func Lex(input string) iter.Seq[token.Token] {
	state := newLexerState(input)

	return func(yield func(token.Token) bool) {
		for {
			tok := nextToken(state)

			if !yield(tok) || tok.Type == token.EOF || tok.Type == token.ILLEGAL {
				break
			}
		}
	}

}
