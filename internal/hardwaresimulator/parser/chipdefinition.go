package parser

type ParsedChipDefinition struct {
	ChipName ChipName
	Inputs   []IO
	Outputs  []IO
	Parts    []Part
}

type ChipName struct {
	Name string
	Loc  Loc
}

type Loc struct {
	Line   int
	Column int
}

type IO struct {
	Name  string
	Width int
	Loc   Loc
}

type Part struct {
	Name        string
	Connections []Connection
	Loc         Loc
}

type Connection struct {
	Pin    Pin
	Signal Signal
	Loc    Loc
}

type Pin struct {
	Name  string
	Range Range
	Loc   Loc
}

type Signal struct {
	Name  string
	Range Range
	Loc   Loc
}

type Range struct {
	IsSpecified bool
	Start       int
	End         int
	Loc         Loc
}

func (chd *ParsedChipDefinition) GetUsedChipNames() []string {
	var usedChips []string
	seen := make(map[string]bool)

	for _, part := range chd.Parts {
		if !seen[part.Name] {
			usedChips = append(usedChips, part.Name)
			seen[part.Name] = true
		}
	}

	return usedChips
}
