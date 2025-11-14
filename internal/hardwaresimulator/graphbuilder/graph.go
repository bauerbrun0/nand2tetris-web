package graphbuilder

import (
	"fmt"
	"strings"
)

type Graph struct {
	Nodes        []*Node // the chips/parts
	InputPins    map[string]*Pin
	OutputPins   map[string]*Pin
	InternalPins map[string]*InternalPin

	Edges map[*Node][]*Node // adjacency, node -> dependent nodes

	StatesInitialized bool
}

type Node struct {
	ChipName   string
	InputPins  map[string]*Pin
	OutputPins map[string]*Pin

	// for sequential chips (eg. DFF) to keep track of the state
	State map[string][]bool // signal name -> bits

	SubGraph *Graph // nil if built-in chip
}

type Pin struct {
	Name string
	Bits []*BitRef
}

type InternalPin struct {
	Name           string
	Bits           []*BitRef
	SourceNode     *Node
	DependentNodes map[*Node][]int // node -> indexes of the bits it depends on
}

type BitRef struct {
	Bit *Bit
}

type Bit struct {
	IsSequential bool
	Value        bool
}

func (g *Graph) String() string {
	var sb strings.Builder

	sb.WriteString("Graph:\n")
	sb.WriteString("\tInputPins:\n")
	for pinName, pin := range g.InputPins {
		sb.WriteString(fmt.Sprintf("\t\t%s: %d bits\n", pinName, len(pin.Bits)))
		for i, bitRef := range pin.Bits {
			sb.WriteString(fmt.Sprintf("\t\t\tBit %d: Value=%v, pointer=%p\n", i, bitRef.Bit.Value, bitRef.Bit))
		}
	}
	sb.WriteString("\tOutputPins:\n")
	for pinName, pin := range g.OutputPins {
		sb.WriteString(fmt.Sprintf("\t\t%s: %d bits\n", pinName, len(pin.Bits)))
		for i, bitRef := range pin.Bits {
			sb.WriteString(fmt.Sprintf("\t\t\tBit %d: Value=%v, pointer=%p\n", i, bitRef.Bit.Value, bitRef.Bit))
		}
	}
	sb.WriteString("\tNodes:\n")
	for _, node := range g.Nodes {
		sb.WriteString(node.String())
	}

	return sb.String()
}

func (n *Node) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Node [%s]\n", n.ChipName))
	sb.WriteString("\tInputPins:\n")
	for pinName, pin := range n.InputPins {
		sb.WriteString(fmt.Sprintf("\t\t%s: %d bits\n", pinName, len(pin.Bits)))
		for i, bitRef := range pin.Bits {
			sb.WriteString(fmt.Sprintf("\t\t\tBit %d: Value=%v, pointer=%p\n", i, bitRef.Bit.Value, bitRef.Bit))
		}
	}
	sb.WriteString("\tOutputPins:\n")
	for pinName, pin := range n.OutputPins {
		sb.WriteString(fmt.Sprintf("\t\t%s: %d bits\n", pinName, len(pin.Bits)))
		for i, bitRef := range pin.Bits {
			sb.WriteString(fmt.Sprintf("\t\t\tBit %d: Value=%v, pointer=%p\n", i, bitRef.Bit.Value, bitRef.Bit))
		}
	}

	return sb.String()
}
