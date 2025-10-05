package lexer

import (
	"testing"

	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/token"
)

func TestNextToken(t *testing.T) {

	type expectedToken struct {
		TokenType token.TokenType
		Literal   string
		Line      int
		Column    int
	}

	tests := []struct {
		name           string
		input          string
		expectedTokens []expectedToken
	}{
		{
			name:  "Single character tokens",
			input: `=,;:{}()[]`,
			expectedTokens: []expectedToken{
				{token.ASSIGN, "=", 1, 1},
				{token.COMMA, ",", 1, 2},
				{token.SEMICOLON, ";", 1, 3},
				{token.COLON, ":", 1, 4},
				{token.LBRACE, "{", 1, 5},
				{token.RBRACE, "}", 1, 6},
				{token.LPAREN, "(", 1, 7},
				{token.RPAREN, ")", 1, 8},
				{token.LBRACKET, "[", 1, 9},
				{token.RBRACKET, "]", 1, 10},
				{token.EOF, "", 1, 11},
			},
		},
		{
			name:  "Keywords",
			input: `CHIP IN OUT PARTS true false`,
			expectedTokens: []expectedToken{
				{token.CHIP, "CHIP", 1, 1},
				{token.IN, "IN", 1, 6},
				{token.OUT, "OUT", 1, 9},
				{token.PARTS, "PARTS", 1, 13},
				{token.TRUE, "true", 1, 19},
				{token.FALSE, "false", 1, 24},
				{token.EOF, "", 1, 29},
			},
		},
		{
			name:  "Identifiers and Numbers",
			input: `foo Foo fooBar foo_bar _foo fooBar12 1 16`,
			expectedTokens: []expectedToken{
				{token.IDENTIFIER, "foo", 1, 1},
				{token.IDENTIFIER, "Foo", 1, 5},
				{token.IDENTIFIER, "fooBar", 1, 9},
				{token.IDENTIFIER, "foo_bar", 1, 16},
				{token.IDENTIFIER, "_foo", 1, 24},
				{token.IDENTIFIER, "fooBar12", 1, 29},
				{token.NUMBER, "1", 1, 38},
				{token.NUMBER, "16", 1, 40},
				{token.EOF, "", 1, 42},
			},
		},
		{
			name:  "Illegal tokens",
			input: `@$*/`,
			expectedTokens: []expectedToken{
				{token.ILLEGAL, "@", 1, 1},
				{token.ILLEGAL, "$", 1, 2},
				{token.ILLEGAL, "*", 1, 3},
				{token.ILLEGAL, "/", 1, 4},
				{token.EOF, "", 1, 5},
			},
		},
		{
			name:  "Range token",
			input: `.. 1..12`,
			expectedTokens: []expectedToken{
				{token.RANGE, "..", 1, 1},
				{token.NUMBER, "1", 1, 4},
				{token.RANGE, "..", 1, 5},
				{token.NUMBER, "12", 1, 7},
				{token.EOF, "", 1, 9},
			},
		},
		{
			name: "Whitespace handling",
			input: `    foo
    bar
foo    bar`,
			expectedTokens: []expectedToken{
				{token.IDENTIFIER, "foo", 1, 5},
				{token.IDENTIFIER, "bar", 2, 5},
				{token.IDENTIFIER, "foo", 3, 1},
				{token.IDENTIFIER, "bar", 3, 8},
				{token.EOF, "", 3, 11},
			},
		},
		{
			name: "Line comments",
			input: `// This is a line comment
CHIP And {
// Another comment
    // Yet another comment
    IN// Last comment`,
			expectedTokens: []expectedToken{
				{token.LINE_COMMENT, "", 1, 1},
				{token.CHIP, "CHIP", 2, 1},
				{token.IDENTIFIER, "And", 2, 6},
				{token.LBRACE, "{", 2, 10},
				{token.LINE_COMMENT, "", 3, 1},
				{token.LINE_COMMENT, "", 4, 5},
				{token.IN, "IN", 5, 5},
				{token.LINE_COMMENT, "", 5, 7},
				{token.EOF, "", 5, 23},
			},
		},
		{
			name: "Block comments",
			input: `/* This is a
block comment */
CHIP And /* Another comment */ {
    IN a, b; /* Yet another comment */
    OUT out;
    /* Last comment */
}`,
			expectedTokens: []expectedToken{
				{token.BLOCK_COMMENT, "", 1, 1},
				{token.CHIP, "CHIP", 3, 1},
				{token.IDENTIFIER, "And", 3, 6},
				{token.BLOCK_COMMENT, "", 3, 10},
				{token.LBRACE, "{", 3, 32},
				{token.IN, "IN", 4, 5},
				{token.IDENTIFIER, "a", 4, 8},
				{token.COMMA, ",", 4, 9},
				{token.IDENTIFIER, "b", 4, 11},
				{token.SEMICOLON, ";", 4, 12},
				{token.BLOCK_COMMENT, "", 4, 14},
				{token.OUT, "OUT", 5, 5},
				{token.IDENTIFIER, "out", 5, 9},
				{token.SEMICOLON, ";", 5, 12},
				{token.BLOCK_COMMENT, "", 6, 5},
				{token.RBRACE, "}", 7, 1},
				{token.EOF, "", 7, 2},
			},
		},
		{
			name:  "Unclosed block comment",
			input: `/* This is an unclosed block comment`,
			expectedTokens: []expectedToken{
				{token.ILLEGAL, "EOF", 1, 1},
				{token.EOF, "", 1, 37},
			},
		},
		{
			name:  "Immediate EOF after block comment start",
			input: `/*`,
			expectedTokens: []expectedToken{
				{token.ILLEGAL, "EOF", 1, 1},
				{token.EOF, "", 1, 3},
			},
		},
		{
			name: "Complete HDL snippet",
			input: `// This is a line comment
/* This is a
 * block comment
 */
CHIP And16 {
    IN a[16], b;
    OUT out;

    PARTS:
    Nand(a = a, b = b, out = aNandB);
    Not(in = aNandB, out = out);
    Not(in = a[0..1] b = true, c = false);
}`,
			expectedTokens: []expectedToken{
				{token.LINE_COMMENT, "", 1, 1},
				{token.BLOCK_COMMENT, "", 2, 1},
				{token.CHIP, "CHIP", 5, 1},
				{token.IDENTIFIER, "And16", 5, 6},
				{token.LBRACE, "{", 5, 12},
				{token.IN, "IN", 6, 5},
				{token.IDENTIFIER, "a", 6, 8},
				{token.LBRACKET, "[", 6, 9},
				{token.NUMBER, "16", 6, 10},
				{token.RBRACKET, "]", 6, 12},
				{token.COMMA, ",", 6, 13},
				{token.IDENTIFIER, "b", 6, 15},
				{token.SEMICOLON, ";", 6, 16},
				{token.OUT, "OUT", 7, 5},
				{token.IDENTIFIER, "out", 7, 9},
				{token.SEMICOLON, ";", 7, 12},
				{token.PARTS, "PARTS", 9, 5},
				{token.COLON, ":", 9, 10},
				{token.IDENTIFIER, "Nand", 10, 5},
				{token.LPAREN, "(", 10, 9},
				{token.IDENTIFIER, "a", 10, 10},
				{token.ASSIGN, "=", 10, 12},
				{token.IDENTIFIER, "a", 10, 14},
				{token.COMMA, ",", 10, 15},
				{token.IDENTIFIER, "b", 10, 17},
				{token.ASSIGN, "=", 10, 19},
				{token.IDENTIFIER, "b", 10, 21},
				{token.COMMA, ",", 10, 22},
				{token.IDENTIFIER, "out", 10, 24},
				{token.ASSIGN, "=", 10, 28},
				{token.IDENTIFIER, "aNandB", 10, 30},
				{token.RPAREN, ")", 10, 36},
				{token.SEMICOLON, ";", 10, 37},
				{token.IDENTIFIER, "Not", 11, 5},
				{token.LPAREN, "(", 11, 8},
				{token.IDENTIFIER, "in", 11, 9},
				{token.ASSIGN, "=", 11, 12},
				{token.IDENTIFIER, "aNandB", 11, 14},
				{token.COMMA, ",", 11, 20},
				{token.IDENTIFIER, "out", 11, 22},
				{token.ASSIGN, "=", 11, 26},
				{token.IDENTIFIER, "out", 11, 28},
				{token.RPAREN, ")", 11, 31},
				{token.SEMICOLON, ";", 11, 32},
				{token.IDENTIFIER, "Not", 12, 5},
				{token.LPAREN, "(", 12, 8},
				{token.IDENTIFIER, "in", 12, 9},
				{token.ASSIGN, "=", 12, 12},
				{token.IDENTIFIER, "a", 12, 14},
				{token.LBRACKET, "[", 12, 15},
				{token.NUMBER, "0", 12, 16},
				{token.RANGE, "..", 12, 17},
				{token.NUMBER, "1", 12, 19},
				{token.RBRACKET, "]", 12, 20},
				{token.IDENTIFIER, "b", 12, 22},
				{token.ASSIGN, "=", 12, 24},
				{token.TRUE, "true", 12, 26},
				{token.COMMA, ",", 12, 30},
				{token.IDENTIFIER, "c", 12, 32},
				{token.ASSIGN, "=", 12, 34},
				{token.FALSE, "false", 12, 36},
				{token.RPAREN, ")", 12, 41},
				{token.SEMICOLON, ";", 12, 42},
				{token.RBRACE, "}", 13, 1},
				{token.EOF, "", 13, 2},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			lexer := New(tt.input)

			for i, expected := range tt.expectedTokens {
				tok := lexer.NextToken()

				if tok.TokenType != expected.TokenType {
					t.Fatalf("At test[%d]: wrong TokenType, got %q; expected %q", i, tok.TokenType, expected.TokenType)
				}

				if tok.Literal != expected.Literal {
					t.Fatalf("At test[%d]: wrong Literal, got %q; expected %q", i, tok.Literal, expected.Literal)
				}

				if tok.Line != expected.Line {
					t.Fatalf("At test[%d]: wrong Line, got %d; expected %d", i, tok.Line, expected.Line)
				}

				if tok.Column != expected.Column {
					t.Fatalf("At test[%d]: wrong Column, got %d; expected %d", i, tok.Column, expected.Column)
				}
			}
		})
	}
}

func TestTokenize(t *testing.T) {
	input := `$`

	lexer := New(input)
	_, err := lexer.Tokenize()

	if err == nil {
		t.Fatalf("Expected an error for illegal token, but got nil")
	}

	if err.Error() != "Lexer error at line 1, column 1: illegal token '$'" {
		t.Fatalf("Unexpected error message: %s", err.Error())
	}
}
