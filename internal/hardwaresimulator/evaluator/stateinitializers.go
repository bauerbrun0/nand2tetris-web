package evaluator

import (
	"strconv"

	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/graphbuilder"
)

// State initializer functions set the initial internal state of sequential chips.
var BuiltinChipStateInitializerFns = map[string]func(node *graphbuilder.Node){
	"DFF": func(node *graphbuilder.Node) {
		node.State = map[string][]bool{
			"out": {false},
		}
	},
	"Bit": func(node *graphbuilder.Node) {
		node.State = map[string][]bool{
			"out": {false},
		}
	},
	"Register": func(node *graphbuilder.Node) {
		node.State = map[string][]bool{
			"out": make([]bool, 16),
		}
	},
	"PC": func(node *graphbuilder.Node) {
		node.State = map[string][]bool{
			"out": make([]bool, 16),
		}
	},
	"RAM8": func(node *graphbuilder.Node) {
		node.State = make(map[string][]bool, 8)
		for i := range 8 {
			node.State["out_"+strconv.Itoa(i)] = make([]bool, 16)
		}
	},
	"RAM64": func(node *graphbuilder.Node) {
		node.State = make(map[string][]bool, 64)
		for i := range 64 {
			node.State["out_"+strconv.Itoa(i)] = make([]bool, 16)
		}
	},
	"RAM512": func(node *graphbuilder.Node) {
		node.State = make(map[string][]bool, 512)
		for i := range 512 {
			node.State["out_"+strconv.Itoa(i)] = make([]bool, 16)
		}
	},
	"RAM4K": func(node *graphbuilder.Node) {
		node.State = make(map[string][]bool, 4096)
		for i := range 4096 {
			node.State["out_"+strconv.Itoa(i)] = make([]bool, 16)
		}
	},
	"RAM16K": func(node *graphbuilder.Node) {
		node.State = make(map[string][]bool, 16384)
		for i := range 16384 {
			node.State["out_"+strconv.Itoa(i)] = make([]bool, 16)
		}
	},
}
