package lexer

import (
	"fmt"

	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/errors"
	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/token"
)

type Lexer struct {
	input           string
	currentPosition int // current position in input - points to currentChar
	readPosition    int // current reading position in input - points after currentChar
	currentChar     byte
	line            int // line number for error messages
	column          int // column number for error messages
}

func New(input string) *Lexer {
	lexer := &Lexer{input: input}
	lexer.readChar()
	lexer.line = 1
	lexer.column = 1
	return lexer
}

func (l *Lexer) Tokenize() (TokenStream, error) {
	var tokens []token.Token

	for {
		tok := l.NextToken()
		if tok.TokenType == token.ILLEGAL {
			message := fmt.Sprintf("illegal token '%s'", tok.Literal)
			return NewTokenStream([]token.Token{}),
				errors.NewLexingError(message, tok.Line, tok.Column)
		}

		if tok.TokenType != token.LINE_COMMENT && tok.TokenType != token.BLOCK_COMMENT {
			tokens = append(tokens, tok)
		}

		if tok.TokenType == token.EOF {
			return NewTokenStream(tokens), nil
		}
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.currentChar {
	case '=':
		tok = newToken(token.ASSIGN, l.currentChar, l.line, l.column)
	case ',':
		tok = newToken(token.COMMA, l.currentChar, l.line, l.column)
	case ';':
		tok = newToken(token.SEMICOLON, l.currentChar, l.line, l.column)
	case ':':
		tok = newToken(token.COLON, l.currentChar, l.line, l.column)
	case '{':
		tok = newToken(token.LBRACE, l.currentChar, l.line, l.column)
	case '}':
		tok = newToken(token.RBRACE, l.currentChar, l.line, l.column)
	case '(':
		tok = newToken(token.LPAREN, l.currentChar, l.line, l.column)
	case ')':
		tok = newToken(token.RPAREN, l.currentChar, l.line, l.column)
	case '[':
		tok = newToken(token.LBRACKET, l.currentChar, l.line, l.column)
	case ']':
		tok = newToken(token.RBRACKET, l.currentChar, l.line, l.column)
	case '.':
		if l.peekChar() == '.' {
			l.readChar() // read the first '.'
			// the second '.' will be read at the end of the function
			tok = token.Token{TokenType: token.RANGE, Literal: "..", Line: l.line, Column: l.column - 1}
		} else {
			tok = newToken(token.ILLEGAL, l.currentChar, l.line, l.column)
		}
	case '/':
		if l.peekChar() == '/' {
			// handle line comment
			starterColumn := l.column
			for l.currentChar != '\n' && l.currentChar != 0 {
				l.readChar()
			}
			tok = token.Token{TokenType: token.LINE_COMMENT, Literal: "", Line: l.line, Column: starterColumn}
			if l.currentChar == '\n' {
				l.line++
				l.column = 0
			}
		} else if l.peekChar() == '*' {
			starterColumn := l.column
			starterLine := l.line
			l.readChar() // read the first '*'

			// till '*/' or EOF read everyting
			for {
				if l.peekChar() != 0 {
					l.readChar()
				}
				if l.peekChar() == 0 {
					// EOF reached without closing '*/'
					tok = token.Token{TokenType: token.ILLEGAL, Literal: "EOF", Line: starterLine, Column: starterColumn}
					break
				}

				if l.currentChar == '*' && l.peekChar() == '/' {
					l.readChar()
					tok = token.Token{TokenType: token.BLOCK_COMMENT, Literal: "", Line: starterLine, Column: starterColumn}
					break
				}

				if l.currentChar == '\n' {
					l.line++
					l.column = 0
				}
			}
		} else {
			tok = newToken(token.ILLEGAL, l.currentChar, l.line, l.column)
		}
	case 0:
		tok.Literal = ""
		tok.TokenType = token.EOF
		tok.Line = l.line
		tok.Column = l.column
	default:
		if isLetter(l.currentChar) {
			tok.Literal = l.readIdentifier()
			tok.TokenType = token.LookupTokenType(tok.Literal)
			tok.Line = l.line
			tok.Column = l.column - len(tok.Literal)
			return tok
		} else if isDigit(l.currentChar) {
			tok.Literal = l.readNumber()
			tok.TokenType = token.NUMBER
			tok.Line = l.line
			tok.Column = l.column - len(tok.Literal)
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.currentChar, l.line, l.column)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.currentChar = 0
	} else {
		l.currentChar = l.input[l.readPosition]
	}
	l.currentPosition = l.readPosition
	l.readPosition++
	l.column++
}

func (l *Lexer) readIdentifier() string {
	starterPosition := l.currentPosition
	for isLetter(l.currentChar) || isDigit(l.currentChar) {
		l.readChar()
	}
	return l.input[starterPosition:l.currentPosition]
}

func (l *Lexer) readNumber() string {
	starterPosition := l.currentPosition
	for isDigit(l.currentChar) {
		l.readChar()
	}
	return l.input[starterPosition:l.currentPosition]
}

func (l *Lexer) skipWhitespace() {
	for l.currentChar == ' ' || l.currentChar == '\t' || l.currentChar == '\n' || l.currentChar == '\r' {
		switch l.currentChar {
		case '\n':
			l.line++
			l.column = 0
		case '\t':
			l.column += 3
		}

		l.readChar()
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func newToken(tokenType token.TokenType, literal byte, line int, column int) token.Token {
	return token.Token{TokenType: tokenType, Literal: string(literal), Line: line, Column: column}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
