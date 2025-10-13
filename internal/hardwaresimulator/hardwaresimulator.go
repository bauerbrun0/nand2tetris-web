package hardwaresimulator

import (
	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/errors"
	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/lexer"
	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/parser"
	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/resolver"
)

type HardwareSimulator struct {
	hdls map[string]string
}

func New() *HardwareSimulator {
	return &HardwareSimulator{}
}

func (hs *HardwareSimulator) SetChipHDLs(hdls map[string]string) {
	hs.hdls = hdls
}

func (hs *HardwareSimulator) Process(chipName string) error {
	// first, get the HDL for the chip
	// then, lex and parse it
	// then, get the used chips names and check if
	// 	first, whether they are built-in chips
	//  second, whether we have their HDL (custom chips)
	hdl, ok := hs.hdls[chipName]
	if !ok {
		return errors.NewChipNotFoundError(chipName)
	}

	l := lexer.New(hdl)
	ts, err := l.Tokenize()
	if err != nil {
		return err
	}

	p := parser.New(ts)
	chd, err := p.ParseChipDefinition()
	if err != nil {
		return err
	}

	r := resolver.New(chd, chipName, hs.hdls)
	_, _, err = r.Resolve([]string{}, []string{})
	if err != nil {
		return err
	}

	return nil
}
