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

func (e *Evaluator) InitializeNodeStates() {
	if e.Graph.StatesInitialized {
		return
	}
	for _, node := range e.Graph.Nodes {
		e.initializeNodeState(node)
	}
	e.Graph.StatesInitialized = true
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

func (e *Evaluator) EvaluateAndCommit() {
	for _, node := range e.Graph.Nodes {
		e.evaluateAndCommitNode(node)
	}
}

func (e *Evaluator) Apply() {
	for _, node := range e.Graph.Nodes {
		e.applyNodeState(node)
	}
}

func (e *Evaluator) Evaluate() {
	for _, node := range e.Graph.Nodes {
		e.evaluateNode(node)
	}
}

func (e *Evaluator) initializeNodeState(node *graphbuilder.Node) {
	if initializeState, ok := BuiltinChipStateInitializerFns[node.ChipName]; ok {
		initializeState(node)
		return
	}

	if node.SubGraph != nil {
		subEvaluator := New(node.SubGraph)
		subEvaluator.InitializeNodeStates()
	}
}

func (e *Evaluator) evaluateAndCommitNode(node *graphbuilder.Node) {
	commit, committerExists := BuiltinChipComitterFns[node.ChipName]
	evaluate, evaluatorExists := BuiltinChipEvaluatorFns[node.ChipName]

	if evaluatorExists {
		evaluate(node)
	}

	if committerExists {
		commit(node)
	}

	if evaluatorExists || committerExists {
		return
	}

	if node.SubGraph != nil {
		subEvaluator := New(node.SubGraph)
		subEvaluator.EvaluateAndCommit()
	}
}

func (e *Evaluator) evaluateNode(node *graphbuilder.Node) {
	if evaluate, ok := BuiltinChipEvaluatorFns[node.ChipName]; ok {
		evaluate(node)
		return
	}

	if node.SubGraph != nil {
		subEvaluator := New(node.SubGraph)
		subEvaluator.Evaluate()
	}
}

func (e *Evaluator) applyNodeState(node *graphbuilder.Node) {
	if apply, ok := BuiltinChipApplierFns[node.ChipName]; ok {
		apply(node)
		return
	}

	if node.SubGraph != nil {
		subEvaluator := New(node.SubGraph)
		subEvaluator.Apply()
	}
}
