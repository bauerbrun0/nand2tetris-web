package chips

var BuiltInChips = map[string]Chip{
	"Nand": {
		Inputs: map[string]IO{
			"a": {Width: 1},
			"b": {Width: 1},
		},
		Outputs: map[string]IO{
			"out": {Width: 1},
		},
	},
	"DFF": {
		Inputs: map[string]IO{
			"in": {Width: 1},
		},
		Outputs: map[string]IO{
			"out": {Width: 1},
		},
	},
}
