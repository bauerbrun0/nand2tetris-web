package lexer

import "github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/token"

type TokenStream struct {
	tokens  []token.Token
	pos     int
	current *token.Token
	peek    *token.Token
}

func NewTokenStream(tokens []token.Token) TokenStream {
	ts := TokenStream{
		tokens: tokens,
		pos:    -1,
	}
	ts.Next()
	return ts
}

func (ts *TokenStream) Current() *token.Token {
	return ts.current
}

func (ts *TokenStream) Peek() *token.Token {
	return ts.peek
}

func (ts *TokenStream) Next() {
	ts.pos++

	if ts.pos >= len(ts.tokens) {
		ts.current = nil
		return
	}

	ts.current = &ts.tokens[ts.pos]

	if ts.pos+1 >= len(ts.tokens) {
		ts.peek = nil
		return
	}

	ts.peek = &ts.tokens[ts.pos+1]
}
