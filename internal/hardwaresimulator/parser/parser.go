package parser

import (
	"fmt"

	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/errors"
	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/lexer"
	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/token"
)

type Parser struct {
	ts   lexer.TokenStream
	chip *ParsedChipDefinition
}

func New(ts lexer.TokenStream) *Parser {
	chip := &ParsedChipDefinition{}
	return &Parser{ts: ts, chip: chip}
}

func (p *Parser) ParseChipDefinition() (*ParsedChipDefinition, error) {
	err := p.parseChipName()
	if err != nil {
		return nil, err
	}

	err = p.parseChipIO()
	if err != nil {
		return nil, err
	}

	err = p.parseChipParts()
	if err != nil {
		return nil, err
	}

	if !p.curTokenIs(token.RBRACE) {
		message := fmt.Sprintf("expected '}', got [%s] => %s", p.ts.Current().TokenType, p.ts.Current().Literal)
		return nil, newError(message, p.ts.Current().Line, p.ts.Current().Column)
	}

	p.ts.Next()

	if !p.curTokenIs(token.EOF) {
		message := fmt.Sprintf("expected EOF after '}', got [%s] => %s", p.ts.Current().TokenType, p.ts.Current().Literal)
		return nil, newError(message, p.ts.Current().Line, p.ts.Current().Column)
	}

	return p.chip, nil
}

func (p *Parser) parseChipName() error {
	if !p.curTokenIs(token.CHIP) {
		message := fmt.Sprintf("expected CHIP keyword, got [%s] => %s", p.ts.Current().TokenType, p.ts.Current().Literal)
		return newError(message, p.ts.Current().Line, p.ts.Current().Column)
	}
	p.ts.Next()

	if !p.curTokenIs(token.IDENTIFIER) {
		message := fmt.Sprintf("expected chip name, got [%s] => %s", p.ts.Current().TokenType, p.ts.Current().Literal)
		return newError(message, p.ts.Current().Line, p.ts.Current().Column)
	}

	p.chip.ChipName.Name = p.ts.Current().Literal
	p.chip.ChipName.Loc = getLoc(p.ts.Current())

	p.ts.Next()

	if !p.curTokenIs(token.LBRACE) {
		message := fmt.Sprintf("expected '{', got [%s] => %s", p.ts.Current().TokenType, p.ts.Current().Literal)
		return newError(message, p.ts.Current().Line, p.ts.Current().Column)
	}

	p.ts.Next()
	return nil
}

func (p *Parser) parseChipIO() error {
	if !p.curTokenIs(token.IN) && !p.curTokenIs(token.OUT) {
		message := fmt.Sprintf("expected IN or OUT, got [%s] => %s", p.ts.Current().TokenType, p.ts.Current().Literal)
		return newError(message, p.ts.Current().Line, p.ts.Current().Column)
	}

	if p.curTokenIs(token.IN) {
		// parse inputs, then outputs
		p.ts.Next()
		ioList, err := p.parseIOList()
		if err != nil {
			return err
		}
		p.chip.Inputs = ioList

		if !p.curTokenIs(token.OUT) {
			message := fmt.Sprintf("expected OUT, got [%s] => %s", p.ts.Current().TokenType, p.ts.Current().Literal)
			return newError(message, p.ts.Current().Line, p.ts.Current().Column)
		}
		p.ts.Next()
		ioList, err = p.parseIOList()
		if err != nil {
			return err
		}
		p.chip.Outputs = ioList
	} else {
		// parse outputs, then inputs
		p.ts.Next()
		ioList, err := p.parseIOList()
		if err != nil {
			return err
		}
		p.chip.Outputs = ioList

		if !p.curTokenIs(token.IN) {
			message := fmt.Sprintf("expected IN, got [%s] => %s", p.ts.Current().TokenType, p.ts.Current().Literal)
			return newError(message, p.ts.Current().Line, p.ts.Current().Column)
		}
		p.ts.Next()
		ioList, err = p.parseIOList()
		if err != nil {
			return err
		}
		p.chip.Inputs = ioList
	}

	return nil
}

func (p *Parser) parseIOList() ([]IO, error) {
	var ioList []IO
	var currentIO IO

	for {
		if !p.curTokenIs(token.IDENTIFIER) {
			message := fmt.Sprintf("expected identifier, got [%s] => %s", p.ts.Current().TokenType, p.ts.Current().Literal)
			return nil, newError(message, p.ts.Current().Line, p.ts.Current().Column)
		}

		currentIO = IO{}
		currentIO.Name = p.ts.Current().Literal
		currentIO.Width = 1
		currentIO.Loc = getLoc(p.ts.Current())

		p.ts.Next()

		if p.curTokenIs(token.LBRACKET) {
			// parse the width
			p.ts.Next()
			if !p.curTokenIs(token.NUMBER) {
				message := fmt.Sprintf("expected number for width, got [%s] => %s", p.ts.Current().TokenType, p.ts.Current().Literal)
				return nil, newError(message, p.ts.Current().Line, p.ts.Current().Column)
			}

			var width int
			_, err := fmt.Sscanf(p.ts.Current().Literal, "%d", &width)
			if err != nil {
				message := fmt.Sprintf("invalid number for width: %v", p.ts.Current().Literal)
				return nil, newError(message, p.ts.Current().Line, p.ts.Current().Column)
			}
			currentIO.Width = width
			p.ts.Next()

			if !p.curTokenIs(token.RBRACKET) {
				message := fmt.Sprintf("expected ']', got [%s] => %s", p.ts.Current().TokenType, p.ts.Current().Literal)
				return nil, newError(message, p.ts.Current().Line, p.ts.Current().Column)
			}
			p.ts.Next()
		}

		ioList = append(ioList, currentIO)

		if p.curTokenIs(token.COMMA) {
			p.ts.Next()
			continue
		}

		if p.curTokenIs(token.SEMICOLON) {
			p.ts.Next()
			break
		}

		message := fmt.Sprintf("expected ',' or ';', got [%s] => %s", p.ts.Current().TokenType, p.ts.Current().Literal)
		return nil, newError(message, p.ts.Current().Line, p.ts.Current().Column)
	}

	return ioList, nil
}

func (p *Parser) parseChipParts() error {
	if !p.curTokenIs(token.PARTS) {
		message := fmt.Sprintf("expected PARTS keyword, got [%s] => %s", p.ts.Current().TokenType, p.ts.Current().Literal)
		return newError(message, p.ts.Current().Line, p.ts.Current().Column)
	}
	p.ts.Next()

	if !p.curTokenIs(token.COLON) {
		message := fmt.Sprintf("expected ':', got [%s] => %s", p.ts.Current().TokenType, p.ts.Current().Literal)
		return newError(message, p.ts.Current().Line, p.ts.Current().Column)
	}
	p.ts.Next()

	var parts []Part
	var currentPart Part

	for {
		if !p.curTokenIs(token.IDENTIFIER) {
			message := fmt.Sprintf("expected part name, got [%s] => %s", p.ts.Current().TokenType, p.ts.Current().Literal)
			return newError(message, p.ts.Current().Line, p.ts.Current().Column)
		}
		currentPart = Part{}
		currentPart.Name = p.ts.Current().Literal
		currentPart.Loc = getLoc(p.ts.Current())
		p.ts.Next()

		if !p.curTokenIs(token.LPAREN) {
			message := fmt.Sprintf("expected '(', got [%s] => %s", p.ts.Current().TokenType, p.ts.Current().Literal)
			return newError(message, p.ts.Current().Line, p.ts.Current().Column)
		}
		p.ts.Next()

		err := p.parsePartConnections(&currentPart)
		if err != nil {
			return err
		}

		if !p.curTokenIs(token.RPAREN) {
			message := fmt.Sprintf("expected ')', got [%s] => %s", p.ts.Current().TokenType, p.ts.Current().Literal)
			return newError(message, p.ts.Current().Line, p.ts.Current().Column)
		}
		p.ts.Next()

		if !p.curTokenIs(token.SEMICOLON) {
			message := fmt.Sprintf("expected ';', got [%s] => %s", p.ts.Current().TokenType, p.ts.Current().Literal)
			return newError(message, p.ts.Current().Line, p.ts.Current().Column)
		}
		parts = append(parts, currentPart)
		p.ts.Next()

		if p.curTokenIs(token.RBRACE) {
			break
		}

		if !p.curTokenIs(token.IDENTIFIER) && !p.curTokenIs(token.RBRACE) {
			message := fmt.Sprintf("expected part name or '}', got [%s] => %s", p.ts.Current().TokenType, p.ts.Current().Literal)
			return newError(message, p.ts.Current().Line, p.ts.Current().Column)
		}
	}
	p.chip.Parts = parts
	return nil
}

func (p *Parser) parsePartConnections(part *Part) error {
	var connections []Connection
	var currentConnection Connection
	for {
		if !p.curTokenIs(token.IDENTIFIER) {
			message := fmt.Sprintf("expected connection name, got [%s] => %s", p.ts.Current().TokenType, p.ts.Current().Literal)
			return newError(message, p.ts.Current().Line, p.ts.Current().Column)
		}

		currentConnection = Connection{}
		loc := getLoc(p.ts.Current())
		currentConnection.Loc = loc
		currentConnection.Pin = Pin{
			Name:  p.ts.Current().Literal,
			Range: Range{Start: 0, End: 0, Loc: loc},
			Loc:   loc,
		}

		p.ts.Next()
		if p.curTokenIs(token.LBRACKET) {
			p.ts.Next()
			pinRange := Range{}
			err := p.parseRange(&pinRange)
			if err != nil {
				return err
			}
			pinRange.IsSpecified = true
			currentConnection.Pin.Range = pinRange

			if !p.curTokenIs(token.RBRACKET) {
				message := fmt.Sprintf("expected ']' or '..', got [%s] => %s", p.ts.Current().TokenType, p.ts.Current().Literal)
				return newError(message, p.ts.Current().Line, p.ts.Current().Column)
			}
			p.ts.Next()
		}

		if !p.curTokenIs(token.ASSIGN) {
			message := fmt.Sprintf("expected '=', got [%s] => %s", p.ts.Current().TokenType, p.ts.Current().Literal)
			return newError(message, p.ts.Current().Line, p.ts.Current().Column)
		}

		p.ts.Next()

		if !p.curTokenIs(token.IDENTIFIER) {
			message := fmt.Sprintf("expected signal name, got [%s] => %s", p.ts.Current().TokenType, p.ts.Current().Literal)
			return newError(message, p.ts.Current().Line, p.ts.Current().Column)
		}

		currentConnection.Signal = Signal{
			Name:  p.ts.Current().Literal,
			Range: Range{Start: 0, End: 0, Loc: getLoc(p.ts.Current())},
			Loc:   getLoc(p.ts.Current()),
		}
		p.ts.Next()

		if p.curTokenIs(token.LBRACKET) {
			p.ts.Next()
			signalRange := Range{}
			err := p.parseRange(&signalRange)
			if err != nil {
				return err
			}
			signalRange.IsSpecified = true
			currentConnection.Signal.Range = signalRange

			if !p.curTokenIs(token.RBRACKET) {
				message := fmt.Sprintf("expected ']' or '..', got [%s] => %s", p.ts.Current().TokenType, p.ts.Current().Literal)
				return newError(message, p.ts.Current().Line, p.ts.Current().Column)
			}
			p.ts.Next()
		}

		if p.curTokenIs(token.RPAREN) {
			connections = append(connections, currentConnection)
			break
		}

		if p.curTokenIs(token.COMMA) {
			connections = append(connections, currentConnection)
			p.ts.Next()
			continue
		}

		message := fmt.Sprintf("expected ')', ',' or '[', got [%s] => %s", p.ts.Current().TokenType, p.ts.Current().Literal)
		return newError(message, p.ts.Current().Line, p.ts.Current().Column)
	}

	if len(connections) < 2 {
		message := fmt.Sprintf("expected connection name, got [%s] => %s", p.ts.Current().TokenType, p.ts.Current().Literal)
		return newError(message, p.ts.Current().Line, p.ts.Current().Column)
	}
	part.Connections = connections
	return nil
}

func (p *Parser) parseRange(r *Range) error {
	if !p.curTokenIs(token.NUMBER) {
		message := fmt.Sprintf("expected number for range, got [%s] => %s", p.ts.Current().TokenType, p.ts.Current().Literal)
		return newError(message, p.ts.Current().Line, p.ts.Current().Column)
	}

	var start int
	var end int

	_, err := fmt.Sscanf(p.ts.Current().Literal, "%d", &start)
	if err != nil {
		message := fmt.Sprintf("invalid number for range: %v", p.ts.Current().Literal)
		return newError(message, p.ts.Current().Line, p.ts.Current().Column)
	}
	r.Loc = getLoc(p.ts.Current())

	p.ts.Next() // advance to next token which can be '..' or ']'
	// later will be checked by the caller

	if p.curTokenIs(token.RANGE) {
		p.ts.Next()
		if !p.curTokenIs(token.NUMBER) {
			message := fmt.Sprintf("expected number for range end, got [%s] => %s", p.ts.Current().TokenType, p.ts.Current().Literal)
			return newError(message, p.ts.Current().Line, p.ts.Current().Column)
		}

		_, err := fmt.Sscanf(p.ts.Current().Literal, "%d", &end)
		if err != nil {
			message := fmt.Sprintf("invalid number for range end: %v", p.ts.Current().Literal)
			return newError(message, p.ts.Current().Line, p.ts.Current().Column)
		}
		p.ts.Next()
	} else {
		end = start
	}

	r.Start = start
	r.End = end
	return nil
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.ts.Current() != nil && p.ts.Current().TokenType == t
}

func newError(message string, line, column int) error {
	return errors.NewParsingError(message, line, column)
}

func getLoc(t *token.Token) Loc {
	return Loc{Line: t.Line, Column: t.Column}
}
