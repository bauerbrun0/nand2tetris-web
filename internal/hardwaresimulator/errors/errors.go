package errors

import "fmt"

type LexingError struct {
	Message string
	Line    int
	Column  int
}

func (e *LexingError) Error() string {
	return fmt.Sprintf("Lexer error at line %d, column %d: %s", e.Line, e.Column, e.Message)
}

func NewLexingError(message string, line, column int) *LexingError {
	return &LexingError{
		Message: message,
		Line:    line,
		Column:  column,
	}
}

type ParsingError struct {
	Message string
	Line    int
	Column  int
}

func (e *ParsingError) Error() string {
	return fmt.Sprintf("Parser error at line %d, column %d: %s", e.Line, e.Column, e.Message)
}

func NewParsingError(message string, line, column int) *ParsingError {
	return &ParsingError{
		Message: message,
		Line:    line,
		Column:  column,
	}
}
