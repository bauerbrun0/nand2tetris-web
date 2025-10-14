package graphbuilder

import (
	"slices"

	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/chips"
	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/resolver"
)

type GraphBuilder struct {
	chipDefinition  *resolver.ResolvedChipDefinition
	chipDefinitions map[string]*resolver.ResolvedChipDefinition
}

func New(chipDefinitions map[string]*resolver.ResolvedChipDefinition) *GraphBuilder {
	return &GraphBuilder{
		chipDefinitions: chipDefinitions,
	}
}

func (gb *GraphBuilder) BuildGraph(chipName string, inputPins map[string]*Pin, outputPins map[string]*Pin) (*Graph, error) {
	chipDefinition := gb.chipDefinitions[chipName]
	gb.chipDefinition = chipDefinition

	internalSignals := make(map[string]*InternalSignal)
	for signalName, signal := range gb.chipDefinition.InternalSignals {
		internalSignals[signalName] = &InternalSignal{
			Signal: &Signal{Name: signalName, Values: make([]bool, signal.Width)},
		}
	}

	graph := &Graph{
		Nodes:           []*Node{},
		InputPins:       inputPins,
		OutputPins:      outputPins,
		InternalSignals: internalSignals,
		Edges:           map[*Node][]*Node{},
	}

	for _, part := range gb.chipDefinition.Parts {
		err := gb.buildNodeFromPart(graph, &part)
		if err != nil {
			return nil, err
		}
	}

	for _, internalSignal := range graph.InternalSignals {
		sourceNode := internalSignal.SourceNode
		for _, dependentNode := range internalSignal.DependentNodes {
			if !slices.Contains(graph.Edges[sourceNode], dependentNode) && internalSignal.Signal.IsSequential == false {
				graph.Edges[sourceNode] = append(graph.Edges[sourceNode], dependentNode)
			}
		}
	}

	nodesInOrder, err := graph.GetNodesInTopologicalOrder()
	if err != nil {
		return nil, err
	}
	graph.Nodes = nodesInOrder
	return graph, nil
}

func (gb *GraphBuilder) buildNodeFromPart(graph *Graph, part *resolver.Part) error {
	inputs := make(map[string]*Pin)
	outputs := make(map[string]*Pin)
	internalSignals := make(map[string]*InternalSignal)
	node := &Node{}

	// first, find the chip definition for this part
	// it can be a built-in chip or a user-defined chip
	builtinChipIO, isBuiltinChip := chips.BuiltInChips[part.Name]
	customChipDefinition, _ := gb.chipDefinitions[part.Name]

	for _, inputConn := range part.InputConnections {
		if _, ok := inputs[inputConn.Pin.Name]; !ok {
			var width int
			if isBuiltinChip {
				width = builtinChipIO.Inputs[inputConn.Pin.Name].Width
			} else {
				width = customChipDefinition.Inputs[inputConn.Pin.Name].Width
			}

			newPin := &Pin{
				Width:       width,
				Connections: []Connection{},
			}
			inputs[inputConn.Pin.Name] = newPin
		}

		pin := inputs[inputConn.Pin.Name]
		if signal, ok := graph.InternalSignals[inputConn.Signal.Name]; ok {
			connection := Connection{
				PinRange:    Range{Start: inputConn.Pin.Range.Start, End: inputConn.Pin.Range.End},
				SignalRange: Range{Start: inputConn.Signal.Range.Start, End: inputConn.Signal.Range.End},
				Signal:      signal.Signal,
			}
			signal.DependentNodes = append(signal.DependentNodes, node)
			pin.Connections = append(pin.Connections, connection)
			inputs[inputConn.Pin.Name] = pin
		} else {
			parentInput := graph.InputPins[inputConn.Signal.Name]
			for _, conn := range parentInput.Connections {
				connection := Connection{
					PinRange:    Range{Start: inputConn.Pin.Range.Start, End: inputConn.Pin.Range.End},
					SignalRange: Range{Start: inputConn.Signal.Range.Start, End: inputConn.Signal.Range.End},
					Signal:      conn.Signal,
				}
				pin.Connections = append(pin.Connections, connection)
			}
			inputs[inputConn.Pin.Name] = pin
		}
	}

	for _, outputConn := range part.OutputConnections {
		if _, ok := outputs[outputConn.Pin.Name]; !ok {
			var width int
			if isBuiltinChip {
				width = builtinChipIO.Outputs[outputConn.Pin.Name].Width
			} else {
				width = customChipDefinition.Outputs[outputConn.Pin.Name].Width
			}

			newPin := &Pin{
				Width:       width,
				Connections: []Connection{},
			}
			outputs[outputConn.Pin.Name] = newPin
		}

		pin := outputs[outputConn.Pin.Name]
		if signal, ok := graph.InternalSignals[outputConn.Signal.Name]; ok {
			if isBuiltinChip && part.Name == "DFF" {
				signal.Signal.IsSequential = true
			}
			connection := Connection{
				PinRange:    Range{Start: outputConn.Pin.Range.Start, End: outputConn.Pin.Range.End},
				SignalRange: Range{Start: outputConn.Signal.Range.Start, End: outputConn.Signal.Range.End},
				Signal:      signal.Signal,
			}
			signal.SourceNode = node
			pin.Connections = append(pin.Connections, connection)
			outputs[outputConn.Pin.Name] = pin
		} else {
			newSignal := &Signal{
				Name:         outputConn.Signal.Name,
				Values:       make([]bool, outputConn.Signal.Range.End-outputConn.Signal.Range.Start+1),
				IsSequential: isBuiltinChip && part.Name == "DFF",
			}
			connection := Connection{
				PinRange:    Range{Start: outputConn.Pin.Range.Start, End: outputConn.Pin.Range.End},
				SignalRange: Range{Start: outputConn.Signal.Range.Start, End: outputConn.Signal.Range.End},
				Signal:      newSignal,
			}

			graph.OutputPins[outputConn.Signal.Name].Connections = append(graph.OutputPins[outputConn.Signal.Name].Connections, connection)
			pin.Connections = append(pin.Connections, connection)
			outputs[outputConn.Pin.Name] = pin
		}
	}

	if !isBuiltinChip {
		subGb := New(gb.chipDefinitions)
		sg, err := subGb.BuildGraph(part.Name, inputs, outputs)
		if err != nil {
			return err
		}
		node.Subgraph = sg
	}

	node.Type = part.Name
	node.Inputs = inputs
	node.Outputs = outputs
	node.InternalSignals = internalSignals

	graph.Nodes = append(graph.Nodes, node)
	return nil
}
