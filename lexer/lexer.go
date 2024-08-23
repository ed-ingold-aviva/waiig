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

func buildString(initialByte byte, rest []byte) string {
	return string(append([]byte{initialByte}, rest...))
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

func consumeWhitespace(state *lexerState) {
	_ = state.getNextBytes(isWhitespace)
	return
}

func getDouble(ch byte, state *lexerState) (tok token.Token, ok bool) {
	nextChar := state.peek()
	if ch == '=' && nextChar == '=' {
		tok = eqToken
		ok = true
	} else if ch == '!' && nextChar == '=' {
		tok = neqToken
		ok = true
	}
	if ok {
		_ = state.next()
	}
	return
}

func getStringToken(ch byte, state *lexerState) token.Token {
	nextStringBytes := state.getNextBytes(isLetter)

	literal := buildString(ch, nextStringBytes)
	var tok token.Token
	if keywordToken, ok := keywordTokens[literal]; ok {
		tok = keywordToken
	} else {
		tok = identToken(literal)
	}
	return tok
}

func getDigitToken(ch byte, state *lexerState) token.Token {
	nextDigitBytes := state.getNextBytes(isDigit)
	return digitToken(buildString(ch, nextDigitBytes))
}

func nextToken(state *lexerState) token.Token {
	consumeWhitespace(state)
	ch := state.next()

	var tok token.Token
	if ch == EofByte {
		tok = eofToken
	} else if doubleToken, ok := getDouble(ch, state); ok {
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
