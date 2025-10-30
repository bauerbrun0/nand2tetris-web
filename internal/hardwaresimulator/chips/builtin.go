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
	"DMux4Way": {
		Inputs: map[string]IO{
			"in":  {Width: 1},
			"sel": {Width: 2},
		},
		Outputs: map[string]IO{
			"a": {Width: 1},
			"b": {Width: 1},
			"c": {Width: 1},
			"d": {Width: 1},
		},
	},
	"DMux8Way": {
		Inputs: map[string]IO{
			"in":  {Width: 1},
			"sel": {Width: 3},
		},
		Outputs: map[string]IO{
			"a": {Width: 1},
			"b": {Width: 1},
			"c": {Width: 1},
			"d": {Width: 1},
			"e": {Width: 1},
			"f": {Width: 1},
			"g": {Width: 1},
			"h": {Width: 1},
		},
	},
	"And16": {
		Inputs: map[string]IO{
			"a": {Width: 16},
			"b": {Width: 16},
		},
		Outputs: map[string]IO{
			"out": {Width: 16},
		},
	},
	"Or16": {
		Inputs: map[string]IO{
			"a": {Width: 16},
			"b": {Width: 16},
		},
		Outputs: map[string]IO{
			"out": {Width: 16},
		},
	},
	"Not16": {
		Inputs: map[string]IO{
			"in": {Width: 16},
		},
		Outputs: map[string]IO{
			"out": {Width: 16},
		},
	},
	"Or8Way": {
		Inputs: map[string]IO{
			"in": {Width: 8},
		},
		Outputs: map[string]IO{
			"out": {Width: 1},
		},
	},
	"Mux16": {
		Inputs: map[string]IO{
			"a":   {Width: 16},
			"b":   {Width: 16},
			"sel": {Width: 1},
		},
		Outputs: map[string]IO{
			"out": {Width: 16},
		},
	},
	"Mux4Way16": {
		Inputs: map[string]IO{
			"a":   {Width: 16},
			"b":   {Width: 16},
			"c":   {Width: 16},
			"d":   {Width: 16},
			"sel": {Width: 2},
		},
		Outputs: map[string]IO{
			"out": {Width: 16},
		},
	},
	"Mux8Way16": {
		Inputs: map[string]IO{
			"a":   {Width: 16},
			"b":   {Width: 16},
			"c":   {Width: 16},
			"d":   {Width: 16},
			"e":   {Width: 16},
			"f":   {Width: 16},
			"g":   {Width: 16},
			"h":   {Width: 16},
			"sel": {Width: 3},
		},
		Outputs: map[string]IO{
			"out": {Width: 16},
		},
	},
	"HalfAdder": {
		Inputs: map[string]IO{
			"a": {Width: 1},
			"b": {Width: 1},
		},
		Outputs: map[string]IO{
			"sum":   {Width: 1},
			"carry": {Width: 1},
		},
	},
	"FullAdder": {
		Inputs: map[string]IO{
			"a": {Width: 1},
			"b": {Width: 1},
			"c": {Width: 1},
		},
		Outputs: map[string]IO{
			"sum":   {Width: 1},
			"carry": {Width: 1},
		},
	},
	"Add16": {
		Inputs: map[string]IO{
			"a": {Width: 16},
			"b": {Width: 16},
		},
		Outputs: map[string]IO{
			"out": {Width: 16},
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
	"Bit": {
		Inputs: map[string]IO{
			"in":   {Width: 1},
			"load": {Width: 1},
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
	case "Bit":
		return signalName == "out" && bitIndex == 0
	case "RAM64":
		return signalName == "out" && bitIndex >= 0 && bitIndex < 16
	default:
		return false
	}
}
