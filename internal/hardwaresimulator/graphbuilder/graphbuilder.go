package graphbuilder

import (
	"fmt"
	"slices"

	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/chips"
	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/errors"
	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/resolver"
)

type GraphBuilder struct {
	chipDefinition  *resolver.ResolvedChipDefinition
	chipDefinitions map[string]*resolver.ResolvedChipDefinition
	graph           *Graph
}

func New(chipDefinitions map[string]*resolver.ResolvedChipDefinition) *GraphBuilder {
	return &GraphBuilder{
		chipDefinitions: chipDefinitions,
	}
}

func (gb *GraphBuilder) BuildGraph(chipName string) (*Graph, error) {
	chd := gb.chipDefinitions[chipName]
	gb.chipDefinition = chd

	inputPins := make(map[string]*Pin)
	outputPins := make(map[string]*Pin)
	internalPins := make(map[string]*InternalPin)

	for inputName, input := range chd.Inputs {
		inputPins[inputName] = &Pin{
			Name: inputName,
			Bits: createBitsArray(input.Width),
		}
	}

	for outputName, output := range chd.Outputs {
		outputPins[outputName] = &Pin{
			Name: outputName,
			Bits: createBitsArray(output.Width),
		}
	}

	for internalName, internal := range chd.InternalSignals {
		internalPins[internalName] = &InternalPin{
			Name:           internalName,
			Bits:           createBitsArray(internal.Width),
			DependentNodes: make(map[*Node][]int),
		}
	}

	gb.graph = &Graph{
		Nodes:        []*Node{},
		InputPins:    inputPins,
		OutputPins:   outputPins,
		InternalPins: internalPins,
		Edges:        map[*Node][]*Node{},
	}

	for _, part := range chd.Parts {
		err := gb.buildNodeFromPart(&part)
		if err != nil {
			return nil, err
		}
	}

	for _, internalPin := range gb.graph.InternalPins {
		sourceNode := internalPin.SourceNode
		for node, indexes := range internalPin.DependentNodes {
			for i := range indexes {
				if internalPin.Bits[i].IsSequential == false && !slices.Contains(gb.graph.Edges[sourceNode], node) {
					gb.graph.Edges[sourceNode] = append(gb.graph.Edges[sourceNode], node)
				}
			}
		}
	}

	nodesInOrder, err := getNodesInTopologicalOrder(gb.graph)
	if err != nil {
		return nil, err
	}
	gb.graph.Nodes = nodesInOrder
	g := gb.graph
	fmt.Println("Input pin pointers")
	for name, pin := range g.InputPins {
		fmt.Printf("\tpin [%s] bit pointer: %p\n", name, pin.Bits[0])
	}

	fmt.Println("Output pin pointers")
	for name, pin := range g.OutputPins {
		fmt.Printf("\tpin [%s] bit pointer: %p\n", name, pin.Bits[0])
	}

	if len(g.InternalPins) > 0 {
		fmt.Println("Internal pin pointers")
	}
	for name, pin := range g.InternalPins {
		fmt.Printf("\tpin [%s] bit pointer: %p\n", name, pin.Bits[0])
	}

	for _, node := range gb.graph.Nodes {
		fmt.Printf("Node [%s]:\n", node.ChipName)
		fmt.Println("\tInput pin pointers")
		for name, pin := range node.InputPins {
			fmt.Printf("\t\tpin [%s] bit pointer: %p\n", name, pin.Bits[0])
		}

		fmt.Println("\tOutput pin pointers")
		for name, pin := range node.OutputPins {
			fmt.Printf("\t\tpin [%s] bit pointer: %p\n", name, pin.Bits[0])
		}

		if node.SubGraph == nil {
			continue
		}
		fmt.Println("\tInternal pin pointers")
		for name, pin := range node.SubGraph.InternalPins {
			fmt.Printf("\t\tpin [%s] bit pointer: %p\n", name, pin.Bits[0])
		}
	}

	return gb.graph, nil
}

func (gb *GraphBuilder) BuildGraphWithExistingIOPins(chipName string, inputPins, outputPins map[string]*Pin) (*Graph, error) {
	chd := gb.chipDefinitions[chipName]
	gb.chipDefinition = chd

	internalPins := make(map[string]*InternalPin)

	for internalName, internal := range chd.InternalSignals {
		internalPins[internalName] = &InternalPin{
			Name:           internalName,
			Bits:           createBitsArray(internal.Width),
			DependentNodes: make(map[*Node][]int),
		}
	}

	gb.graph = &Graph{
		Nodes:        []*Node{},
		InputPins:    inputPins,
		OutputPins:   outputPins,
		InternalPins: internalPins,
		Edges:        map[*Node][]*Node{},
	}

	for _, part := range chd.Parts {
		err := gb.buildNodeFromPart(&part)
		if err != nil {
			return nil, err
		}
	}

	for _, internalPin := range gb.graph.InternalPins {
		sourceNode := internalPin.SourceNode
		for node, indexes := range internalPin.DependentNodes {
			for i := range indexes {
				if internalPin.Bits[i].IsSequential == false && !slices.Contains(gb.graph.Edges[sourceNode], node) {
					gb.graph.Edges[sourceNode] = append(gb.graph.Edges[sourceNode], node)
				}
			}
		}
	}

	nodesInOrder, err := getNodesInTopologicalOrder(gb.graph)
	if err != nil {
		return nil, err
	}
	gb.graph.Nodes = nodesInOrder

	return gb.graph, nil
}

func (gb *GraphBuilder) buildNodeFromPart(part *resolver.Part) error {
	node := &Node{}

	inputPins := make(map[string]*Pin)
	outputPins := make(map[string]*Pin)

	builtinChipDef, isBuiltin := chips.BuiltInChips[part.Name]
	customChipDef, _ := gb.chipDefinitions[part.Name]

	// we initialize the input and output pins based on the chip definitions
	// which can be either built-in or user-defined
	// we fill every bit as false initially
	if isBuiltin {
		for inputName, input := range builtinChipDef.Inputs {
			inputPins[inputName] = &Pin{
				Name: inputName,
				Bits: createBitsArray(input.Width),
			}
		}
		for outputName, output := range builtinChipDef.Outputs {
			bits := createBitsArray(output.Width)
			for i, bit := range bits {
				bit.IsSequential = chips.IsSequentialBit(part.Name, outputName, i)
			}
			outputPins[outputName] = &Pin{
				Name: outputName,
				Bits: bits,
			}
		}
	} else {
		for inputName, input := range customChipDef.Inputs {
			inputPins[inputName] = &Pin{
				Name: inputName,
				Bits: createBitsArray(input.Width),
			}
		}
		for outputName, output := range customChipDef.Outputs {
			outputPins[outputName] = &Pin{
				Name: outputName,
				Bits: createBitsArray(output.Width),
			}
		}
	}

	// now we have to connect the pins or rather modify the bits
	// based on the connections defined in the connections section
	// of the part. We use the Ranges to determine which bits to connect/change

	for _, inputConnection := range part.InputConnections {
		signalName := inputConnection.Signal.Name
		isBooleanConstant := signalName == "true" || signalName == "false"
		if isBooleanConstant {
			// create the constant bits
			bits := make([]*Bit, inputConnection.Signal.Range.End-inputConnection.Signal.Range.Start+1)
			value := signalName == "true"
			for i := range bits {
				bits[i] = &Bit{
					Value: value,
				}
			}
			// set the bits
			for i, bit := range bits {
				inputPins[inputConnection.Pin.Name].Bits[inputConnection.Pin.Range.Start+i] = bit
			}
		} else if parentInputPin, ok := gb.graph.InputPins[signalName]; ok {
			// get the needed bits from the parent input pin
			neededBits := parentInputPin.Bits[inputConnection.Signal.Range.Start : inputConnection.Signal.Range.End+1]
			// set the bits of the input pin of the node
			for i, bit := range neededBits {
				inputPins[inputConnection.Pin.Name].Bits[inputConnection.Pin.Range.Start+i] = bit
			}
		} else {
			internalPin := gb.graph.InternalPins[signalName]
			neededBits := internalPin.Bits[inputConnection.Signal.Range.Start : inputConnection.Signal.Range.End+1]
			for i, bit := range neededBits {
				inputPins[inputConnection.Pin.Name].Bits[inputConnection.Pin.Range.Start+i] = bit
			}
			for i := inputConnection.Signal.Range.Start; i <= inputConnection.Signal.Range.End; i++ {
				internalPin.DependentNodes[node] = append(internalPin.DependentNodes[node], i)
			}
		}
	}

	if !isBuiltin {
		subGraphBuilder := New(gb.chipDefinitions)
		subGraph, err := subGraphBuilder.BuildGraphWithExistingIOPins(part.Name, inputPins, outputPins)
		if err != nil {
			return err
		}
		node.SubGraph = subGraph
	}

	for _, outputConnection := range part.OutputConnections {
		signalName := outputConnection.Signal.Name
		if parentOutputPin, ok := gb.graph.OutputPins[signalName]; ok {
			neededBits := outputPins[outputConnection.Pin.Name].Bits[outputConnection.Pin.Range.Start : outputConnection.Pin.Range.End+1]
			for i, bit := range neededBits {
				parentOutputPin.Bits[outputConnection.Signal.Range.Start+i] = bit
			}
		} else {
			internalPin := gb.graph.InternalPins[signalName]
			neededBits := outputPins[outputConnection.Pin.Name].Bits[outputConnection.Pin.Range.Start : outputConnection.Pin.Range.End+1]
			internalPin.Bits = neededBits // we replace the bits here, because internal pins cannot be partially defined
			internalPin.SourceNode = node
		}
	}

	node.ChipName = part.Name
	node.InputPins = inputPins
	node.OutputPins = outputPins

	gb.graph.Nodes = append(gb.graph.Nodes, node)
	return nil
}

func createBitsArray(width int) []*Bit {
	bits := make([]*Bit, width)
	for i := range bits {
		bits[i] = &Bit{}
	}
	return bits
}

func getNodesInTopologicalOrder(g *Graph) ([]*Node, error) {
	indegrees := make(map[*Node]int)
	for _, dependentNodes := range g.Edges {
		for _, dependentNode := range dependentNodes {
			if _, exists := indegrees[dependentNode]; !exists {
				indegrees[dependentNode] = 1
			} else {
				indegrees[dependentNode]++
			}
		}
	}
	for _, node := range g.Nodes {
		if _, exists := indegrees[node]; !exists {
			indegrees[node] = 0
		}
	}

	queue := []*Node{}
	// start with nodes that have no dependencies
	for node, indegree := range indegrees {
		if indegree == 0 {
			queue = append(queue, node)
		}
	}
	// the final topological order
	order := []*Node{}

	for len(queue) > 0 {
		// pop from queue
		node := queue[0]
		queue = queue[1:]

		order = append(order, node)
		// get the nodes that depend on this node
		for _, dependentNode := range g.Edges[node] {
			indegrees[dependentNode]--
			if indegrees[dependentNode] == 0 {
				queue = append(queue, dependentNode)
			}
		}
	}

	if len(order) != len(g.Nodes) {
		return nil, errors.NewSimulationError("Graph has cycles, cannot determine topological order")
	}
	return order, nil
}
