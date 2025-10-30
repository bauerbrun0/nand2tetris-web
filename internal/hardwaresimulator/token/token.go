package token

type TokenType string

type Token struct {
	TokenType TokenType
	Literal   string
	Line      int
	Column    int
}

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	IDENTIFIER TokenType = "IDENTIFIER"
	NUMBER     TokenType = "NUMBER"

	ASSIGN TokenType = "="

	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"
	COLON     TokenType = ":"

	LBRACE   TokenType = "{"
	RBRACE   TokenType = "}"
	LPAREN   TokenType = "("
	RPAREN   TokenType = ")"
	LBRACKET TokenType = "["
	RBRACKET TokenType = "]"

	RANGE TokenType = ".."

	CHIP  TokenType = "CHIP"
	IN    TokenType = "IN"
	OUT   TokenType = "OUT"
	PARTS TokenType = "PARTS"
	TRUE  TokenType = "TRUE"
	FALSE TokenType = "FALSE"

	LINE_COMMENT  TokenType = "LINE_COMMENT"
	BLOCK_COMMENT TokenType = "BLOCK_COMMENT"
)

var keywords = map[string]TokenType{
	"CHIP":  CHIP,
	"IN":    IN,
	"OUT":   OUT,
	"PARTS": PARTS,
	"true":  TRUE,
	"false": FALSE,
}

func LookupTokenType(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENTIFIER
}
