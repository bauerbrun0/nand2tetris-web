package hardwaresimulator

import (
	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/errors"
	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/evaluator"
	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/graphbuilder"
	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/lexer"
	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/parser"
	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/resolver"
)

type HardwareSimulator struct {
	hdls      map[string]string
	evaluator *evaluator.Evaluator
}

func New() *HardwareSimulator {
	return &HardwareSimulator{}
}

func (hs *HardwareSimulator) SetChipHDLs(hdls map[string]string) {
	hs.hdls = hdls
}

func (hs *HardwareSimulator) Process(chipName string) (map[string]int, map[string]int, error) {
	// first, get the HDL for the chip
	// then, lex and parse it
	// then, get the used chips names and check if
	// 	first, whether they are built-in chips
	//  second, whether we have their HDL (custom chips)
	hdl, ok := hs.hdls[chipName]
	if !ok {
		return nil, nil, errors.NewChipNotFoundError(chipName)
	}

	l := lexer.New(hdl)
	ts, err := l.Tokenize()
	if err != nil {
		return nil, nil, err
	}

	p := parser.New(ts)
	chd, err := p.ParseChipDefinition()
	if err != nil {
		return nil, nil, err
	}

	r := resolver.New(chd, chipName, hs.hdls)
	rchd, rchds, err := r.Resolve([]string{}, []string{})
	if err != nil {
		return nil, nil, err
	}
	rchds[rchd.Name] = rchd

	inputs := make(map[string]int)
	for inputName, input := range rchd.Inputs {
		inputs[inputName] = input.Width
	}
	outputs := make(map[string]int)
	for outputName, output := range rchd.Outputs {
		outputs[outputName] = output.Width
	}

	gb := graphbuilder.New(rchds)
	g, err := gb.BuildGraph(rchd.Name)
	if err != nil {
		return nil, nil, err
	}

	e := evaluator.New(g)
	hs.evaluator = e

	return inputs, outputs, nil
}

func (hs *HardwareSimulator) Evaluate(inputs map[string][]bool) (map[string][]bool, map[string][]bool) {
	hs.evaluator.SetInputs(inputs)
	hs.evaluator.Evaluate(false)
	outputs, internalPins := hs.evaluator.GetOutputsAndInternalPins()
	return outputs, internalPins
}

func (hs *HardwareSimulator) Tick(inputs map[string][]bool) (map[string][]bool, map[string][]bool) {
	hs.evaluator.SetInputs(inputs)
	hs.evaluator.Evaluate(true)
	outputs, internalPins := hs.evaluator.GetOutputsAndInternalPins()
	return outputs, internalPins
}

func (hs *HardwareSimulator) Tock(inputs map[string][]bool) (map[string][]bool, map[string][]bool) {
	return hs.Evaluate(inputs)
}
