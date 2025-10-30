package evaluator

import (
	"fmt"

	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/graphbuilder"
)

type Evaluator struct {
	Graph *graphbuilder.Graph
}

func New(graph *graphbuilder.Graph) *Evaluator {
	return &Evaluator{
		Graph: graph,
	}
}

func (e *Evaluator) SetInputs(inputs map[string][]bool) {
	for inputName, input := range e.Graph.InputPins {
		for i, bit := range input.Bits {
			bit.Bit.Value = inputs[inputName][i]
		}
	}
}

func (e *Evaluator) GetOutputsAndInternalPins() (map[string][]bool, map[string][]bool) {
	outputs := make(map[string][]bool)
	for outputName, output := range e.Graph.OutputPins {
		bits := make([]bool, len(output.Bits))
		for i, bit := range output.Bits {
			bits[i] = bit.Bit.Value
		}
		outputs[outputName] = bits
	}

	internals := make(map[string][]bool)
	for internalName, internal := range e.Graph.InternalPins {
		bits := make([]bool, len(internal.Bits))
		for i, bit := range internal.Bits {
			bits[i] = bit.Bit.Value
		}
		internals[internalName] = bits
	}

	return outputs, internals
}

func (e *Evaluator) Evaluate() {
	for _, node := range e.Graph.Nodes {
		e.evaluateNode(node)
	}
	// second loop is a hack until proper evaluation order is implemented
	// this intrduces performance issues for large graphs
	for _, node := range e.Graph.Nodes {
		e.evaluateNode(node)
	}

}

func (e *Evaluator) Commit() {
	for _, node := range e.Graph.Nodes {
		e.commitNode(node)
	}
}

func (e *Evaluator) evaluateNode(node *graphbuilder.Node) {
	switch node.ChipName {
	case "Nand":
		a := node.InputPins["a"].Bits[0].Bit.Value
		b := node.InputPins["b"].Bits[0].Bit.Value
		v := !(a && b)
		node.OutputPins["out"].Bits[0].Bit.Value = v
	case "And":
		a := node.InputPins["a"].Bits[0].Bit.Value
		b := node.InputPins["b"].Bits[0].Bit.Value
		v := a && b
		node.OutputPins["out"].Bits[0].Bit.Value = v
	case "Or":
		a := node.InputPins["a"].Bits[0].Bit.Value
		b := node.InputPins["b"].Bits[0].Bit.Value
		v := a || b
		node.OutputPins["out"].Bits[0].Bit.Value = v
	case "Not":
		in := node.InputPins["in"].Bits[0].Bit.Value
		v := !in
		node.OutputPins["out"].Bits[0].Bit.Value = v
	case "Xor":
		a := node.InputPins["a"].Bits[0].Bit.Value
		b := node.InputPins["b"].Bits[0].Bit.Value
		v := (a || b) && !(a && b)
		node.OutputPins["out"].Bits[0].Bit.Value = v
	case "Mux":
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
	case "DMux":
		in := node.InputPins["in"].Bits[0].Bit.Value
		sel := node.InputPins["sel"].Bits[0].Bit.Value
		if sel {
			node.OutputPins["b"].Bits[0].Bit.Value = in
			node.OutputPins["a"].Bits[0].Bit.Value = false
		} else {
			node.OutputPins["a"].Bits[0].Bit.Value = in
			node.OutputPins["b"].Bits[0].Bit.Value = false
		}
	case "DMux4Way":
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
	case "DMux8Way":
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
	case "And16":
		aBits := node.InputPins["a"].Bits
		bBits := node.InputPins["b"].Bits
		for i := 0; i < 16; i++ {
			a := aBits[i].Bit.Value
			b := bBits[i].Bit.Value
			node.OutputPins["out"].Bits[i].Bit.Value = a && b
		}
	case "Or16":
		aBits := node.InputPins["a"].Bits
		bBits := node.InputPins["b"].Bits
		for i := 0; i < 16; i++ {
			a := aBits[i].Bit.Value
			b := bBits[i].Bit.Value
			node.OutputPins["out"].Bits[i].Bit.Value = a || b
		}
	case "Not16":
		inBits := node.InputPins["in"].Bits
		for i := 0; i < 16; i++ {
			in := inBits[i].Bit.Value
			node.OutputPins["out"].Bits[i].Bit.Value = !in
		}
	case "Or8Way":
		inBits := node.InputPins["in"].Bits
		result := false
		for i := 0; i < 8; i++ {
			in := inBits[i].Bit.Value
			result = result || in
		}
		node.OutputPins["out"].Bits[0].Bit.Value = result
	case "Mux16":
		aBits := node.InputPins["a"].Bits
		bBits := node.InputPins["b"].Bits
		sel := node.InputPins["sel"].Bits[0].Bit.Value

		for i := 0; i < 16; i++ {
			var v bool
			if sel {
				v = bBits[i].Bit.Value
			} else {
				v = aBits[i].Bit.Value
			}
			node.OutputPins["out"].Bits[i].Bit.Value = v
		}
	case "Mux4Way16":
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

		for i := 0; i < 16; i++ {
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
	case "Mux8Way16":
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

		for i := 0; i < 16; i++ {
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
	case "HalfAdder":
		a := node.InputPins["a"].Bits[0].Bit.Value
		b := node.InputPins["b"].Bits[0].Bit.Value

		sum := a != b
		carry := a && b
		node.OutputPins["sum"].Bits[0].Bit.Value = sum
		node.OutputPins["carry"].Bits[0].Bit.Value = carry
	case "FullAdder":
		a := node.InputPins["a"].Bits[0].Bit.Value
		b := node.InputPins["b"].Bits[0].Bit.Value
		c := node.InputPins["c"].Bits[0].Bit.Value

		sum := (a != b) != c
		carry := (a && b) || (c && (a != b))
		node.OutputPins["sum"].Bits[0].Bit.Value = sum
		node.OutputPins["carry"].Bits[0].Bit.Value = carry
	case "Inc16":
		inBits := node.InputPins["in"].Bits
		carry := true
		for i := 0; i < 16; i++ {
			in := inBits[i].Bit.Value

			sum := (in != carry)
			carry = in && carry
			node.OutputPins["out"].Bits[i].Bit.Value = sum
		}
	case "Add16":
		aBits := node.InputPins["a"].Bits
		bBits := node.InputPins["b"].Bits

		carry := false
		for i := 0; i < 16; i++ {
			a := aBits[i].Bit.Value
			b := bBits[i].Bit.Value

			sum := (a != b) != carry
			carry = (a && b) || (carry && (a != b))

			node.OutputPins["out"].Bits[i].Bit.Value = sum
		}
	case "DFF":
		if node.State == nil {
			node.State = make(map[string][]bool)
			node.State["out"] = []bool{false} // initial state
		}

		// DFF logic: output the current state
		node.OutputPins["out"].Bits[0].Bit.Value = node.State["out"][0]
	case "Bit":
		if node.State == nil {
			node.State = make(map[string][]bool)
			node.State["out"] = []bool{false} // initial state
		}

		node.OutputPins["out"].Bits[0].Bit.Value = node.State["out"][0]
	case "Register":
		if node.State == nil {
			node.State = make(map[string][]bool)
			node.State["out"] = make([]bool, 16) // initial state
		}

		// output the current state
		for i := 0; i < 16; i++ {
			node.OutputPins["out"].Bits[i].Bit.Value = node.State["out"][i]
		}
	case "PC":
		if node.State == nil {
			node.State = make(map[string][]bool)
			node.State["out"] = make([]bool, 16) // initial state
		}

		// output the current state
		for i := 0; i < 16; i++ {
			node.OutputPins["out"].Bits[i].Bit.Value = node.State["out"][i]
		}
	case "RAM8":
		addressBits := node.InputPins["address"].Bits
		address := 0
		for i, bit := range addressBits {
			if bit.Bit.Value {
				address |= (1 << i)
			}
		}

		// initialize state if not exist
		if node.State == nil {
			node.State = make(map[string][]bool)
			for i := range 8 {
				node.State[fmt.Sprintf("out_%d", i)] = make([]bool, 16)
			}
		}

		// output the value at the given address
		outKey := fmt.Sprintf("out_%d", address)
		for i := range 16 {
			node.OutputPins["out"].Bits[i].Bit.Value = node.State[outKey][i]
		}
	case "RAM64":
		addressBits := node.InputPins["address"].Bits
		address := 0
		for i, bit := range addressBits {
			if bit.Bit.Value {
				address |= (1 << i)
			}
		}

		// initialize state if not exist
		if node.State == nil {
			node.State = make(map[string][]bool)
			for i := range 64 {
				node.State[fmt.Sprintf("out_%d", i)] = make([]bool, 16)
			}
		}

		// output the value at the given address
		outKey := fmt.Sprintf("out_%d", address)
		for i := range 16 {
			node.OutputPins["out"].Bits[i].Bit.Value = node.State[outKey][i]
		}
	case "RAM512":
		addressBits := node.InputPins["address"].Bits
		address := 0
		for i, bit := range addressBits {
			if bit.Bit.Value {
				address |= (1 << i)
			}
		}

		// initialize state if not exist
		if node.State == nil {
			node.State = make(map[string][]bool)
			for i := range 512 {
				node.State[fmt.Sprintf("out_%d", i)] = make([]bool, 16)
			}
		}

		// output the value at the given address
		outKey := fmt.Sprintf("out_%d", address)
		for i := range 16 {
			node.OutputPins["out"].Bits[i].Bit.Value = node.State[outKey][i]
		}
	case "RAM4K":
		addressBits := node.InputPins["address"].Bits
		address := 0
		for i, bit := range addressBits {
			if bit.Bit.Value {
				address |= (1 << i)
			}
		}

		// initialize state if not exist
		if node.State == nil {
			node.State = make(map[string][]bool)
			for i := range 4096 {
				node.State[fmt.Sprintf("out_%d", i)] = make([]bool, 16)
			}
		}

		// output the value at the given address
		outKey := fmt.Sprintf("out_%d", address)
		for i := range 16 {
			node.OutputPins["out"].Bits[i].Bit.Value = node.State[outKey][i]
		}
	case "RAM16K":
		addressBits := node.InputPins["address"].Bits
		address := 0
		for i, bit := range addressBits {
			if bit.Bit.Value {
				address |= (1 << i)
			}
		}

		// initialize state if not exist
		if node.State == nil {
			node.State = make(map[string][]bool)
			for i := range 16384 {
				node.State[fmt.Sprintf("out_%d", i)] = make([]bool, 16)
			}
		}

		// output the value at the given address
		outKey := fmt.Sprintf("out_%d", address)
		for i := range 16 {
			node.OutputPins["out"].Bits[i].Bit.Value = node.State[outKey][i]
		}
	default:
		// custom chip with subgraph
		if node.SubGraph == nil {
			// should never be nil here
			return
		}

		subEvaluator := New(node.SubGraph)
		subEvaluator.Evaluate()
	}
}

func (e *Evaluator) commitNode(node *graphbuilder.Node) {
	switch node.ChipName {
	case "DFF":
		// create state maps if not exist
		if node.State == nil {
			node.State = make(map[string][]bool)
			node.State["out"] = []bool{false} // initial state
		}

		// DFF logic: store the input value into state
		node.State["out"][0] = node.InputPins["in"].Bits[0].Bit.Value
	case "Bit":
		if node.State == nil {
			node.State = make(map[string][]bool)
			node.State["out"] = []bool{false} // initial state
		}

		load := node.InputPins["load"].Bits[0].Bit.Value
		if !load {
			return // do not store if load is false
		}

		node.State["out"][0] = node.InputPins["in"].Bits[0].Bit.Value
	case "Register":
		if node.State == nil {
			node.State = make(map[string][]bool)
			node.State["out"] = make([]bool, 16) // initial state
		}

		load := node.InputPins["load"].Bits[0].Bit.Value
		if !load {
			return // do not store if load is false
		}

		for i := 0; i < 16; i++ {
			node.State["out"][i] = node.InputPins["in"].Bits[i].Bit.Value
		}
	case "PC":
		if node.State == nil {
			node.State = make(map[string][]bool)
			node.State["out"] = make([]bool, 16) // initial state
		}

		load := node.InputPins["load"].Bits[0].Bit.Value
		inc := node.InputPins["inc"].Bits[0].Bit.Value
		reset := node.InputPins["reset"].Bits[0].Bit.Value

		if reset {
			for i := 0; i < 16; i++ {
				node.State["out"][i] = false
			}
			return
		}
		if load {
			for i := 0; i < 16; i++ {
				node.State["out"][i] = node.InputPins["in"].Bits[i].Bit.Value
			}
			return
		}
		if inc {
			carry := true
			for i := 0; i < 16; i++ {
				in := node.State["out"][i]
				sum := (in != carry)
				carry = in && carry
				node.State["out"][i] = sum
			}
			return
		}
	case "RAM8":
		addressBits := node.InputPins["address"].Bits
		address := 0
		for i, bit := range addressBits {
			if bit.Bit.Value {
				address |= (1 << i)
			}
		}

		load := node.InputPins["load"].Bits[0].Bit.Value
		if !load {
			return // do not store if load is false
		}

		// create state maps if not exist
		if node.State == nil {
			node.State = make(map[string][]bool)
			for i := range 8 {
				node.State[fmt.Sprintf("out_%d", i)] = make([]bool, 16)
			}
		}

		// store the input value into the addressed memory location
		outKey := fmt.Sprintf("out_%d", address)
		for i := range 16 {
			node.State[outKey][i] = node.InputPins["in"].Bits[i].Bit.Value
		}
	case "RAM64":
		addressBits := node.InputPins["address"].Bits
		address := 0
		for i, bit := range addressBits {
			if bit.Bit.Value {
				address |= (1 << i)
			}
		}

		load := node.InputPins["load"].Bits[0].Bit.Value
		if !load {
			return // do not store if load is false
		}

		// create state maps if not exist
		if node.State == nil {
			node.State = make(map[string][]bool)
			for i := range 64 {
				node.State[fmt.Sprintf("out_%d", i)] = make([]bool, 16)
			}
		}

		// store the input value into the addressed memory location
		outKey := fmt.Sprintf("out_%d", address)
		for i := range 16 {
			node.State[outKey][i] = node.InputPins["in"].Bits[i].Bit.Value
		}
	case "RAM512":
		addressBits := node.InputPins["address"].Bits
		address := 0
		for i, bit := range addressBits {
			if bit.Bit.Value {
				address |= (1 << i)
			}
		}

		load := node.InputPins["load"].Bits[0].Bit.Value
		if !load {
			return // do not store if load is false
		}

		// create state maps if not exist
		if node.State == nil {
			node.State = make(map[string][]bool)
			for i := range 512 {
				node.State[fmt.Sprintf("out_%d", i)] = make([]bool, 16)
			}
		}

		// store the input value into the addressed memory location
		outKey := fmt.Sprintf("out_%d", address)
		for i := range 16 {
			node.State[outKey][i] = node.InputPins["in"].Bits[i].Bit.Value
		}
	case "RAM4K":
		addressBits := node.InputPins["address"].Bits
		address := 0
		for i, bit := range addressBits {
			if bit.Bit.Value {
				address |= (1 << i)
			}
		}

		load := node.InputPins["load"].Bits[0].Bit.Value
		if !load {
			return // do not store if load is false
		}

		// create state maps if not exist
		if node.State == nil {
			node.State = make(map[string][]bool)
			for i := range 4096 {
				node.State[fmt.Sprintf("out_%d", i)] = make([]bool, 16)
			}
		}

		// store the input value into the addressed memory location
		outKey := fmt.Sprintf("out_%d", address)
		for i := range 16 {
			node.State[outKey][i] = node.InputPins["in"].Bits[i].Bit.Value
		}
	case "RAM16K":
		addressBits := node.InputPins["address"].Bits
		address := 0
		for i, bit := range addressBits {
			if bit.Bit.Value {
				address |= (1 << i)
			}
		}

		load := node.InputPins["load"].Bits[0].Bit.Value
		if !load {
			return // do not store if load is false
		}

		// create state maps if not exist
		if node.State == nil {
			node.State = make(map[string][]bool)
			for i := range 16384 {
				node.State[fmt.Sprintf("out_%d", i)] = make([]bool, 16)
			}
		}

		// store the input value into the addressed memory location
		outKey := fmt.Sprintf("out_%d", address)
		for i := range 16 {
			node.State[outKey][i] = node.InputPins["in"].Bits[i].Bit.Value
		}
	default:
		// custom chip with subgraph
		if node.SubGraph == nil {
			return
		}
		subEvaluator := New(node.SubGraph)
		subEvaluator.Commit()
	}
}
