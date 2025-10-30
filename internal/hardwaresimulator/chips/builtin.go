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
	"And": {
		Inputs: map[string]IO{
			"a": {Width: 1},
			"b": {Width: 1},
		},
		Outputs: map[string]IO{
			"out": {Width: 1},
		},
	},
	"Or": {
		Inputs: map[string]IO{
			"a": {Width: 1},
			"b": {Width: 1},
		},
		Outputs: map[string]IO{
			"out": {Width: 1},
		},
	},
	"Not": {
		Inputs: map[string]IO{
			"in": {Width: 1},
		},
		Outputs: map[string]IO{
			"out": {Width: 1},
		},
	},
	"Xor": {
		Inputs: map[string]IO{
			"a": {Width: 1},
			"b": {Width: 1},
		},
		Outputs: map[string]IO{
			"out": {Width: 1},
		},
	},
	"Mux": {
		Inputs: map[string]IO{
			"a":   {Width: 1},
			"b":   {Width: 1},
			"sel": {Width: 1},
		},
		Outputs: map[string]IO{
			"out": {Width: 1},
		},
	},
	"DMux": {
		Inputs: map[string]IO{
			"in":  {Width: 1},
			"sel": {Width: 1},
		},
		Outputs: map[string]IO{
			"a": {Width: 1},
			"b": {Width: 1},
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
	"RAM64": {
		Inputs: map[string]IO{
			"in":      {Width: 16},
			"load":    {Width: 1},
			"address": {Width: 6},
		},
		Outputs: map[string]IO{
			"out": {Width: 16},
		},
	},
}

func IsSequentialBit(chipName string, signalName string, bitIndex int) bool {
	if _, exists := BuiltInChips[chipName]; !exists {
		return false
	}

	switch chipName {
	case "DFF":
		return signalName == "out" && bitIndex == 0
	case "RAM64":
		return signalName == "out" && bitIndex >= 0 && bitIndex < 16
	default:
		return false
	}
}
