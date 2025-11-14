package evaluator

import (
	"strconv"

	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/graphbuilder"
)

// Evaluator functions set the output pins of the chip based on the input pins
// or in case of sequential chips, based on the committed state.
var BuiltinChipEvaluatorFns = map[string]func(node *graphbuilder.Node){
	"Nand": func(node *graphbuilder.Node) {
		a := node.InputPins["a"].Bits[0].Bit.Value
		b := node.InputPins["b"].Bits[0].Bit.Value
		v := !(a && b)
		node.OutputPins["out"].Bits[0].Bit.Value = v
	},
	"And": func(node *graphbuilder.Node) {
		a := node.InputPins["a"].Bits[0].Bit.Value
		b := node.InputPins["b"].Bits[0].Bit.Value
		v := a && b
		node.OutputPins["out"].Bits[0].Bit.Value = v
	},
	"Or": func(node *graphbuilder.Node) {
		a := node.InputPins["a"].Bits[0].Bit.Value
		b := node.InputPins["b"].Bits[0].Bit.Value
		v := a || b
		node.OutputPins["out"].Bits[0].Bit.Value = v
	},
	"Not": func(node *graphbuilder.Node) {
		in := node.InputPins["in"].Bits[0].Bit.Value
		v := !in
		node.OutputPins["out"].Bits[0].Bit.Value = v
	},
	"Xor": func(node *graphbuilder.Node) {
		a := node.InputPins["a"].Bits[0].Bit.Value
		b := node.InputPins["b"].Bits[0].Bit.Value
		v := (a || b) && !(a && b)
		node.OutputPins["out"].Bits[0].Bit.Value = v
	},
	"Mux": func(node *graphbuilder.Node) {
		a := node.InputPins["a"].Bits[0].Bit.Value
		b := node.InputPins["b"].Bits[0].Bit.Value
		sel := node.InputPins["sel"].Bits[0].Bit.Value
		var v bool
		if sel {
			v = b
		} else {
			v = a
		}
		node.OutputPins["out"].Bits[0].Bit.Value = v
	},
	"DMux": func(node *graphbuilder.Node) {
		in := node.InputPins["in"].Bits[0].Bit.Value
		sel := node.InputPins["sel"].Bits[0].Bit.Value
		if sel {
			node.OutputPins["b"].Bits[0].Bit.Value = in
			node.OutputPins["a"].Bits[0].Bit.Value = false
		} else {
			node.OutputPins["a"].Bits[0].Bit.Value = in
			node.OutputPins["b"].Bits[0].Bit.Value = false
		}
	},
	"DMux4Way": func(node *graphbuilder.Node) {
		in := node.InputPins["in"].Bits[0].Bit.Value
		sel := node.InputPins["sel"].Bits
		address := 0
		for i, bit := range sel {
			if bit.Bit.Value {
				address |= (1 << i)
			}
		}

		node.OutputPins["a"].Bits[0].Bit.Value = false
		node.OutputPins["b"].Bits[0].Bit.Value = false
		node.OutputPins["c"].Bits[0].Bit.Value = false
		node.OutputPins["d"].Bits[0].Bit.Value = false

		switch address {
		case 0:
			node.OutputPins["a"].Bits[0].Bit.Value = in
		case 1:
			node.OutputPins["b"].Bits[0].Bit.Value = in
		case 2:
			node.OutputPins["c"].Bits[0].Bit.Value = in
		case 3:
			node.OutputPins["d"].Bits[0].Bit.Value = in
		}
	},
	"DMux8Way": func(node *graphbuilder.Node) {
		in := node.InputPins["in"].Bits[0].Bit.Value
		sel := node.InputPins["sel"].Bits
		address := 0
		for i, bit := range sel {
			if bit.Bit.Value {
				address |= (1 << i)
			}
		}

		node.OutputPins["a"].Bits[0].Bit.Value = false
		node.OutputPins["b"].Bits[0].Bit.Value = false
		node.OutputPins["c"].Bits[0].Bit.Value = false
		node.OutputPins["d"].Bits[0].Bit.Value = false
		node.OutputPins["e"].Bits[0].Bit.Value = false
		node.OutputPins["f"].Bits[0].Bit.Value = false
		node.OutputPins["g"].Bits[0].Bit.Value = false
		node.OutputPins["h"].Bits[0].Bit.Value = false

		switch address {
		case 0:
			node.OutputPins["a"].Bits[0].Bit.Value = in
		case 1:
			node.OutputPins["b"].Bits[0].Bit.Value = in
		case 2:
			node.OutputPins["c"].Bits[0].Bit.Value = in
		case 3:
			node.OutputPins["d"].Bits[0].Bit.Value = in
		case 4:
			node.OutputPins["e"].Bits[0].Bit.Value = in
		case 5:
			node.OutputPins["f"].Bits[0].Bit.Value = in
		case 6:
			node.OutputPins["g"].Bits[0].Bit.Value = in
		case 7:
			node.OutputPins["h"].Bits[0].Bit.Value = in
		}
	},
	"And16": func(node *graphbuilder.Node) {
		aBits := node.InputPins["a"].Bits
		bBits := node.InputPins["b"].Bits
		for i := range 16 {
			a := aBits[i].Bit.Value
			b := bBits[i].Bit.Value
			node.OutputPins["out"].Bits[i].Bit.Value = a && b
		}
	},
	"Or16": func(node *graphbuilder.Node) {
		aBits := node.InputPins["a"].Bits
		bBits := node.InputPins["b"].Bits
		for i := range 16 {
			a := aBits[i].Bit.Value
			b := bBits[i].Bit.Value
			node.OutputPins["out"].Bits[i].Bit.Value = a || b
		}
	},
	"Not16": func(node *graphbuilder.Node) {
		inBits := node.InputPins["in"].Bits
		for i := range 16 {
			in := inBits[i].Bit.Value
			node.OutputPins["out"].Bits[i].Bit.Value = !in
		}
	},
	"Or8Way": func(node *graphbuilder.Node) {
		inBits := node.InputPins["in"].Bits
		result := false
		for i := range 8 {
			in := inBits[i].Bit.Value
			result = result || in
		}
		node.OutputPins["out"].Bits[0].Bit.Value = result
	},
	"Mux16": func(node *graphbuilder.Node) {
		aBits := node.InputPins["a"].Bits
		bBits := node.InputPins["b"].Bits
		sel := node.InputPins["sel"].Bits[0].Bit.Value

		for i := range 16 {
			var v bool
			if sel {
				v = bBits[i].Bit.Value
			} else {
				v = aBits[i].Bit.Value
			}
			node.OutputPins["out"].Bits[i].Bit.Value = v
		}
	},
	"Mux4Way16": func(node *graphbuilder.Node) {
		aBits := node.InputPins["a"].Bits
		bBits := node.InputPins["b"].Bits
		cBits := node.InputPins["c"].Bits
		dBits := node.InputPins["d"].Bits
		selBits := node.InputPins["sel"].Bits
		address := 0
		for i, bit := range selBits {
			if bit.Bit.Value {
				address |= (1 << i)
			}
		}

		for i := range 16 {
			var v bool
			switch address {
			case 0:
				v = aBits[i].Bit.Value
			case 1:
				v = bBits[i].Bit.Value
			case 2:
				v = cBits[i].Bit.Value
			case 3:
				v = dBits[i].Bit.Value
			}
			node.OutputPins["out"].Bits[i].Bit.Value = v
		}
	},
	"Mux8Way16": func(node *graphbuilder.Node) {
		aBits := node.InputPins["a"].Bits
		bBits := node.InputPins["b"].Bits
		cBits := node.InputPins["c"].Bits
		dBits := node.InputPins["d"].Bits
		eBits := node.InputPins["e"].Bits
		fBits := node.InputPins["f"].Bits
		gBits := node.InputPins["g"].Bits
		hBits := node.InputPins["h"].Bits
		selBits := node.InputPins["sel"].Bits

		address := 0
		for i, bit := range selBits {
			if bit.Bit.Value {
				address |= (1 << i)
			}
		}

		for i := range 16 {
			var v bool
			switch address {
			case 0:
				v = aBits[i].Bit.Value
			case 1:
				v = bBits[i].Bit.Value
			case 2:
				v = cBits[i].Bit.Value
			case 3:
				v = dBits[i].Bit.Value
			case 4:
				v = eBits[i].Bit.Value
			case 5:
				v = fBits[i].Bit.Value
			case 6:
				v = gBits[i].Bit.Value
			case 7:
				v = hBits[i].Bit.Value
			}
			node.OutputPins["out"].Bits[i].Bit.Value = v
		}
	},
	"HalfAdder": func(node *graphbuilder.Node) {
		a := node.InputPins["a"].Bits[0].Bit.Value
		b := node.InputPins["b"].Bits[0].Bit.Value

		sum := a != b
		carry := a && b
		node.OutputPins["sum"].Bits[0].Bit.Value = sum
		node.OutputPins["carry"].Bits[0].Bit.Value = carry
	},
	"FullAdder": func(node *graphbuilder.Node) {
		a := node.InputPins["a"].Bits[0].Bit.Value
		b := node.InputPins["b"].Bits[0].Bit.Value
		c := node.InputPins["c"].Bits[0].Bit.Value

		sum := (a != b) != c
		carry := (a && b) || (c && (a != b))
		node.OutputPins["sum"].Bits[0].Bit.Value = sum
		node.OutputPins["carry"].Bits[0].Bit.Value = carry
	},
	"Inc16": func(node *graphbuilder.Node) {
		inBits := node.InputPins["in"].Bits
		carry := true
		for i := range 16 {
			in := inBits[i].Bit.Value

			sum := (in != carry)
			carry = in && carry
			node.OutputPins["out"].Bits[i].Bit.Value = sum
		}
	},
	"Add16": func(node *graphbuilder.Node) {
		aBits := node.InputPins["a"].Bits
		bBits := node.InputPins["b"].Bits

		carry := false
		for i := range 16 {
			a := aBits[i].Bit.Value
			b := bBits[i].Bit.Value

			sum := (a != b) != carry
			carry = (a && b) || (carry && (a != b))

			node.OutputPins["out"].Bits[i].Bit.Value = sum
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
