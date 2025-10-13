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

type ChipNotFoundError struct {
	ChipName string
}

func (e *ChipNotFoundError) Error() string {
	return fmt.Sprintf("Chip not found: %s", e.ChipName)
}

func NewChipNotFoundError(chipName string) *ChipNotFoundError {
	return &ChipNotFoundError{
		ChipName: chipName,
	}
}

type ResolutionError struct {
	Message string
	File    string
	Line    int
	Column  int
}

func (e *ResolutionError) Error() string {
	if e.Line > 0 && e.Column > 0 {
		return fmt.Sprintf("Resolution error at line %d, column %d: %s", e.Line, e.Column, e.Message)
	} else {
		return fmt.Sprintf("Resolution error: %s", e.Message)
	}
}

func NewResolutionError(message string, line, column int, file string) *ResolutionError {
	return &ResolutionError{
		Message: message,
		File:    file,
		Line:    line,
		Column:  column,
	}
}
