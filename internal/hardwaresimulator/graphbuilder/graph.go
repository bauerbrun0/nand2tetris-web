package graphbuilder

type Graph struct {
	Nodes        []*Node // the chips/parts
	InputPins    map[string]*Pin
	OutputPins   map[string]*Pin
	InternalPins map[string]*InternalPin

	Edges map[*Node][]*Node // adjacency, node -> dependent nodes
}

type Node struct {
	ChipName   string
	InputPins  map[string]*Pin
	OutputPins map[string]*Pin

	SubGraph *Graph // nil if built-in chip
}

type Pin struct {
	Name string
	Bits []*Bit
}

type InternalPin struct {
	Name           string
	Bits           []*Bit
	SourceNode     *Node
	DependentNodes map[*Node][]int // node -> indexes of the bits it depends on
}

type Bit struct {
	IsSequential bool
	Value        bool
}
