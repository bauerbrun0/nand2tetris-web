package evaluator

import (
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
			bit.Value = inputs[inputName][i]
		}
	}
}

func (e *Evaluator) GetOutputsAndInternalPins() (map[string][]bool, map[string][]bool) {
	outputs := make(map[string][]bool)
	for outputName, output := range e.Graph.OutputPins {
		bits := make([]bool, len(output.Bits))
		for i, bit := range output.Bits {
			bits[i] = bit.Value
		}
		outputs[outputName] = bits
	}

	internals := make(map[string][]bool)
	for internalName, internal := range e.Graph.InternalPins {
		bits := make([]bool, len(internal.Bits))
		for i, bit := range internal.Bits {
			bits[i] = bit.Value
		}
		internals[internalName] = bits
	}

	return outputs, internals
}

func (e *Evaluator) Evaluate(isTick bool) {
	for _, node := range e.Graph.Nodes {
		e.evaluateNode(node, isTick)
	}
}

func (e *Evaluator) evaluateNode(node *graphbuilder.Node, isTick bool) {
	switch node.ChipName {
	case "Nand":
		a := node.InputPins["a"].Bits[0].Value
		b := node.InputPins["b"].Bits[0].Value
		v := !(a && b)
		node.OutputPins["out"].Bits[0].Value = v
	case "DFF":
		// create state maps if not exist
		if node.State == nil {
			node.State = make(map[string][]bool)
			node.State["out"] = []bool{false} // initial state
		}

		// DFF logic: output the current state
		node.OutputPins["out"].Bits[0].Value = node.State["out"][0]
		if isTick {
			node.State["out"][0] = node.InputPins["in"].Bits[0].Value
		}
	default:
		// custom chip with subgraph
		for _, n := range node.SubGraph.Nodes {
			e.evaluateNode(n, isTick)
		}
	}
}
