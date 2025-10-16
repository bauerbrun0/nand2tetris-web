package evaluator

import (
	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/graphbuilder"
)

type Evaluator struct {
	graph *graphbuilder.Graph
}

func New(graph *graphbuilder.Graph) *Evaluator {
	return &Evaluator{
		graph: graph,
	}
}

func (e *Evaluator) Evaluate() {
	for _, node := range e.graph.Nodes {
		e.evaluateNode(node)
	}
}

// Step advances the state of sequential elements (like DFF)
// Must be called after Evaluate()
func (e *Evaluator) Step() {
	for _, node := range e.graph.Nodes {
		e.stepNode(node)
	}
}

func (e *Evaluator) stepNode(node *graphbuilder.Node) {
	switch node.ChipName {
	case "DFF":
		in := node.InputPins["in"].Bits[0].Value
		// set next state for "out" pin
		node.State["out"][0] = in
	}
}

func (e *Evaluator) evaluateNode(node *graphbuilder.Node) {
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
		}

		// initialize state for "out" pin if not exist
		if _, exists := node.State["out"]; !exists {
			node.State["out"] = []bool{false} // initial state
		}

		// DFF logic: output the current state
		node.OutputPins["out"].Bits[0].Value = node.State["out"][0]
	default:
		// custom chip with subgraph
		for _, n := range node.SubGraph.Nodes {
			e.evaluateNode(n)
		}
	}
}
