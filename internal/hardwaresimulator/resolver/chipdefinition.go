package resolver

import "github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/chips"

type ResolvedChipDefinition struct {
	Name            string
	Inputs          map[string]chips.IO
	Outputs         map[string]chips.IO
	Parts           []Part
	InternalSignals map[string]InternalSignal
}

type Part struct {
	Name              string
	InputConnections  []Connection
	OutputConnections []Connection
}

type Connection struct {
	Pin    Pin
	Signal Signal
}

type Pin struct {
	Name  string
	Range Range
}

type Signal struct {
	Name  string
	Range Range
}

type Range struct {
	Start int
	End   int
}

type InternalSignal struct {
	Width int
}
