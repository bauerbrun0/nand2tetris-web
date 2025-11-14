package evaluator

import (
	"strconv"

	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/graphbuilder"
)

// Committer functions update the internal state of the chip based on the input pins.
var BuiltinChipComitterFns = map[string]func(node *graphbuilder.Node){
	"DFF": func(node *graphbuilder.Node) {
		inputValue := node.InputPins["in"].Bits[0].Bit.Value
		node.State["out"][0] = inputValue
	},
	"Bit": func(node *graphbuilder.Node) {
		load := node.InputPins["load"].Bits[0].Bit.Value
		if !load {
			return // do not store if load is false
		}

		inputValue := node.InputPins["in"].Bits[0].Bit.Value
		node.State["out"][0] = inputValue
	},
	"Register": func(node *graphbuilder.Node) {
		load := node.InputPins["load"].Bits[0].Bit.Value
		if !load {
			return
		}

		for i := range 16 {
			node.State["out"][i] = node.InputPins["in"].Bits[i].Bit.Value
		}
	},
	"PC": func(node *graphbuilder.Node) {
		load := node.InputPins["load"].Bits[0].Bit.Value
		inc := node.InputPins["inc"].Bits[0].Bit.Value
		reset := node.InputPins["reset"].Bits[0].Bit.Value

		if reset {
			for i := range 16 {
				node.State["out"][i] = false
			}
			return
		}
		if load {
			for i := range 16 {
				node.State["out"][i] = node.InputPins["in"].Bits[i].Bit.Value
			}
			return
		}
		if inc {
			carry := true
			for i := range 16 {
				in := node.State["out"][i]
				sum := (in != carry)
				carry = in && carry
				node.State["out"][i] = sum
			}
			return
		}
	},
	"RAM8": func(node *graphbuilder.Node) {
		addressBits := node.InputPins["address"].Bits
		address := getAddressFromBits(addressBits)

		load := node.InputPins["load"].Bits[0].Bit.Value
		if !load {
			return
		}

		outKey := "out_" + strconv.Itoa(address)
		for i := range 16 {
			node.State[outKey][i] = node.InputPins["in"].Bits[i].Bit.Value
		}
	},
	"RAM64": func(node *graphbuilder.Node) {
		addressBits := node.InputPins["address"].Bits
		address := getAddressFromBits(addressBits)

		load := node.InputPins["load"].Bits[0].Bit.Value
		if !load {
			return
		}

		outKey := "out_" + strconv.Itoa(address)
		for i := range 16 {
			node.State[outKey][i] = node.InputPins["in"].Bits[i].Bit.Value
		}
	},
	"RAM512": func(node *graphbuilder.Node) {
		addressBits := node.InputPins["address"].Bits
		address := getAddressFromBits(addressBits)

		load := node.InputPins["load"].Bits[0].Bit.Value
		if !load {
			return
		}

		outKey := "out_" + strconv.Itoa(address)
		for i := range 16 {
			node.State[outKey][i] = node.InputPins["in"].Bits[i].Bit.Value
		}
	},
	"RAM4K": func(node *graphbuilder.Node) {
		addressBits := node.InputPins["address"].Bits
		address := getAddressFromBits(addressBits)

		load := node.InputPins["load"].Bits[0].Bit.Value
		if !load {
			return
		}

		outKey := "out_" + strconv.Itoa(address)
		for i := range 16 {
			node.State[outKey][i] = node.InputPins["in"].Bits[i].Bit.Value
		}
	},
	"RAM16K": func(node *graphbuilder.Node) {
		addressBits := node.InputPins["address"].Bits
		address := getAddressFromBits(addressBits)

		load := node.InputPins["load"].Bits[0].Bit.Value
		if !load {
			return
		}

		outKey := "out_" + strconv.Itoa(address)
		for i := range 16 {
			node.State[outKey][i] = node.InputPins["in"].Bits[i].Bit.Value
		}
	},
}
