package evaluator

import (
	"strconv"

	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/graphbuilder"
)

// Applier functions set the output pins of the chip based on the currently committed state.
var BuiltinChipApplierFns = map[string]func(node *graphbuilder.Node){
	"DFF": func(node *graphbuilder.Node) {
		outputValue := node.State["out"][0]
		node.OutputPins["out"].Bits[0].Bit.Value = outputValue
	},
	"Bit": func(node *graphbuilder.Node) {
		outputValue := node.State["out"][0]
		node.OutputPins["out"].Bits[0].Bit.Value = outputValue
	},
	"Register": func(node *graphbuilder.Node) {
		for i := range 16 {
			node.OutputPins["out"].Bits[i].Bit.Value = node.State["out"][i]
		}
	},
	"PC": func(node *graphbuilder.Node) {
		for i := range 16 {
			node.OutputPins["out"].Bits[i].Bit.Value = node.State["out"][i]
		}
	},
	"RAM8": func(node *graphbuilder.Node) {
		addressBits := node.InputPins["address"].Bits
		address := getAddressFromBits(addressBits)

		outKey := "out_" + strconv.Itoa(address)
		for i := range 16 {
			node.OutputPins["out"].Bits[i].Bit.Value = node.State[outKey][i]
		}
	},
	"RAM64": func(node *graphbuilder.Node) {
		addressBits := node.InputPins["address"].Bits
		address := getAddressFromBits(addressBits)

		outKey := "out_" + strconv.Itoa(address)
		for i := range 16 {
			node.OutputPins["out"].Bits[i].Bit.Value = node.State[outKey][i]
		}
	},
	"RAM512": func(node *graphbuilder.Node) {
		addressBits := node.InputPins["address"].Bits
		address := getAddressFromBits(addressBits)

		outKey := "out_" + strconv.Itoa(address)
		for i := range 16 {
			node.OutputPins["out"].Bits[i].Bit.Value = node.State[outKey][i]
		}
	},
	"RAM4K": func(node *graphbuilder.Node) {
		addressBits := node.InputPins["address"].Bits
		address := getAddressFromBits(addressBits)

		outKey := "out_" + strconv.Itoa(address)
		for i := range 16 {
			node.OutputPins["out"].Bits[i].Bit.Value = node.State[outKey][i]
		}
	},
	"RAM16K": func(node *graphbuilder.Node) {
		addressBits := node.InputPins["address"].Bits
		address := getAddressFromBits(addressBits)

		outKey := "out_" + strconv.Itoa(address)
		for i := range 16 {
			node.OutputPins["out"].Bits[i].Bit.Value = node.State[outKey][i]
		}
	},
}
