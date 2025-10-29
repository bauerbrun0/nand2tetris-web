package evaluator

import (
	"fmt"

	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/graphbuilder"
)

type Evaluator struct {
	Graph *graphbuilder.Graph
}

func New(graph *graphbuilder.Graph) *Evaluator {
	return &Evaluator{
		Graph: graph,
	}
}

func (e *Evaluator) SetInputs(inputs map[string][]bool) {
	for inputName, input := range e.Graph.InputPins {
		for i, bit := range input.Bits {
			bit.Bit.Value = inputs[inputName][i]
		}
	}
}

func (e *Evaluator) GetOutputsAndInternalPins() (map[string][]bool, map[string][]bool) {
	outputs := make(map[string][]bool)
	for outputName, output := range e.Graph.OutputPins {
		bits := make([]bool, len(output.Bits))
		for i, bit := range output.Bits {
			bits[i] = bit.Bit.Value
		}
		outputs[outputName] = bits
	}

	internals := make(map[string][]bool)
	for internalName, internal := range e.Graph.InternalPins {
		bits := make([]bool, len(internal.Bits))
		for i, bit := range internal.Bits {
			bits[i] = bit.Bit.Value
		}
		internals[internalName] = bits
	}

	return outputs, internals
}

func (e *Evaluator) Evaluate() {
	for _, node := range e.Graph.Nodes {
		e.evaluateNode(node)
	}
	// second loop is a hack until proper evaluation order is implemented
	// this intrduces performance issues for large graphs
	for _, node := range e.Graph.Nodes {
		e.evaluateNode(node)
	}

}

func (e *Evaluator) Commit() {
	for _, node := range e.Graph.Nodes {
		e.commitNode(node)
	}
}

func (e *Evaluator) evaluateNode(node *graphbuilder.Node) {
	switch node.ChipName {
	case "Nand":
		a := node.InputPins["a"].Bits[0].Bit.Value
		b := node.InputPins["b"].Bits[0].Bit.Value
		v := !(a && b)
		node.OutputPins["out"].Bits[0].Bit.Value = v
	case "And":
		a := node.InputPins["a"].Bits[0].Bit.Value
		b := node.InputPins["b"].Bits[0].Bit.Value
		v := a && b
		node.OutputPins["out"].Bits[0].Bit.Value = v
	case "Or":
		a := node.InputPins["a"].Bits[0].Bit.Value
		b := node.InputPins["b"].Bits[0].Bit.Value
		v := a || b
		node.OutputPins["out"].Bits[0].Bit.Value = v
	case "Not":
		in := node.InputPins["in"].Bits[0].Bit.Value
		v := !in
		node.OutputPins["out"].Bits[0].Bit.Value = v
	case "Xor":
		a := node.InputPins["a"].Bits[0].Bit.Value
		b := node.InputPins["b"].Bits[0].Bit.Value
		v := (a || b) && !(a && b)
		node.OutputPins["out"].Bits[0].Bit.Value = v
	case "Mux":
		a := node.InputPins["a"].Bits[0].Bit.Value
		b := node.InputPins["b"].Bits[0].Bit.Value
		sel := node.InputPins["sel"].Bits[0].Bit.Value
		var v bool
		if sel {
			v = b
		} else {
			v = a
		}
		node.OutputPins["out"].Bits[0].Bit.Value = v
	case "DFF":
		if node.State == nil {
			node.State = make(map[string][]bool)
			node.State["out"] = []bool{false} // initial state
		}

		// DFF logic: output the current state
		node.OutputPins["out"].Bits[0].Bit.Value = node.State["out"][0]
	case "RAM64":
		addressBits := node.InputPins["address"].Bits
		address := 0
		for i, bit := range addressBits {
			if bit.Bit.Value {
				address |= (1 << i)
			}
		}

		// initialize state if not exist
		if node.State == nil {
			node.State = make(map[string][]bool)
			for i := range 64 {
				node.State[fmt.Sprintf("out_%d", i)] = make([]bool, 16)
			}
		}

		// output the value at the given address
		outKey := fmt.Sprintf("out_%d", address)
		for i := range 16 {
			node.OutputPins["out"].Bits[i].Bit.Value = node.State[outKey][i]
		}
	default:
		// custom chip with subgraph
		if node.SubGraph == nil {
			// should never be nil here
			return
		}

		subEvaluator := New(node.SubGraph)
		subEvaluator.Evaluate()
	}
}

func (e *Evaluator) commitNode(node *graphbuilder.Node) {
	switch node.ChipName {
	case "DFF":
		// create state maps if not exist
		if node.State == nil {
			node.State = make(map[string][]bool)
			node.State["out"] = []bool{false} // initial state
		}

		// DFF logic: store the input value into state
		node.State["out"][0] = node.InputPins["in"].Bits[0].Bit.Value
	case "RAM64":
		addressBits := node.InputPins["address"].Bits
		address := 0
		for i, bit := range addressBits {
			if bit.Bit.Value {
				address |= (1 << i)
			}
		}

		load := node.InputPins["load"].Bits[0].Bit.Value
		if !load {
			return // do not store if load is false
		}

		// create state maps if not exist
		if node.State == nil {
			node.State = make(map[string][]bool)
			for i := range 64 {
				node.State[fmt.Sprintf("out_%d", i)] = make([]bool, 16)
			}
		}

		// store the input value into the addressed memory location
		outKey := fmt.Sprintf("out_%d", address)
		for i := range 16 {
			node.State[outKey][i] = node.InputPins["in"].Bits[i].Bit.Value
		}
	default:
		// custom chip with subgraph
		if node.SubGraph == nil {
			return
		}
		subEvaluator := New(node.SubGraph)
		subEvaluator.Commit()
	}
}
