package graphbuilder

import (
	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/errors"
)

type Graph struct {
	Nodes           []*Node // the chips/parts
	InputPins       map[string]*Pin
	OutputPins      map[string]*Pin
	InternalSignals map[string]*InternalSignal

	Edges map[*Node][]*Node // adjacency, node -> dependent nodes
}

func (g *Graph) GetNodesInTopologicalOrder() ([]*Node, error) {
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

type Node struct {
	Type            string // the type of chip/part: built-in chip name like "Nand" or user-defined chip name
	Inputs          map[string]*Pin
	Outputs         map[string]*Pin
	InternalSignals map[string]*InternalSignal // to access internal signals from outside

	// for sequential logic
	State []bool
	Next  []bool

	// for custom chips
	Subgraph *Graph
}

type Pin struct {
	Width       int
	Connections []Connection
}

type Connection struct {
	PinRange    Range
	SignalRange Range
	Signal      *Signal
}

type Range struct {
	Start int
	End   int
}

type InternalSignal struct {
	SourceNode     *Node
	Signal         *Signal
	DependentNodes []*Node
}

type Signal struct {
	Name         string
	Values       []bool
	IsSequential bool
}
