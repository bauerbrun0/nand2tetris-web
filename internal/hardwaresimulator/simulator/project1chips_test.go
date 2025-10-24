package simulator

import (
	"testing"

	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/testutils"
	"github.com/stretchr/testify/assert"
)

func TestProject1ChipsSimulation(t *testing.T) {
	tests := []struct {
		name                        string
		chipFileName                string
		hdls                        map[string]string
		expectedInputsAfterProcess  map[string]int
		expectedOutputsAfterProcess map[string]int
		afterProcess                func(t *testing.T, hs *HardwareSimulator)
	}{
		{
			name:         "Not Chip",
			chipFileName: "NotChip",
			hdls: map[string]string{
				"NotChip": `CHIP NotChip {
					IN in;
					OUT out;

					PARTS:
					Nand(a=in, b=in, out=out, out=outinternal);
				}`,
			},
			expectedInputsAfterProcess:  map[string]int{"in": 1},
			expectedOutputsAfterProcess: map[string]int{"out": 1},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				outputs, internalPins := hs.Evaluate(map[string][]bool{"in": {false}})
				assert.Equal(t, map[string][]bool{"out": {true}}, outputs, "expected out to be true")
				assert.Equal(t, map[string][]bool{"outinternal": {true}}, internalPins, "expected outinternal to be true")

				outputs, internalPins = hs.Evaluate(map[string][]bool{"in": {true}})
				assert.Equal(t, map[string][]bool{"out": {false}}, outputs, "expected out to be false")
				assert.Equal(t, map[string][]bool{"outinternal": {false}}, internalPins, "expected outinternal to be false")

				outputs, internalPins = hs.Tick(map[string][]bool{"in": {false}})
				assert.Equal(t, map[string][]bool{"out": {true}}, outputs, "expected out to be true")
				assert.Equal(t, map[string][]bool{"outinternal": {true}}, internalPins, "expected outinternal to be true")

				outputs, internalPins = hs.Tock(map[string][]bool{"in": {true}})
				assert.Equal(t, map[string][]bool{"out": {false}}, outputs, "expected out to be false")
				assert.Equal(t, map[string][]bool{"outinternal": {false}}, internalPins, "expected outinternal to be false")
			},
		},
		{
			name:         "And Chip",
			chipFileName: "AndChip",
			hdls: map[string]string{
				"AndChip": `CHIP AndChip {
					IN a, b;
					OUT out;

					PARTS:
					Nand(a = a, b = b, out = aNandB);
					NotChip(in = aNandB, out = out);
				}`,
				"NotChip": testutils.ChipImplementations["NotChip"],
			},
			expectedInputsAfterProcess:  map[string]int{"a": 1, "b": 1},
			expectedOutputsAfterProcess: map[string]int{"out": 1},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				outputs, _ := hs.Evaluate(map[string][]bool{"a": {false}, "b": {false}})
				assert.Equal(t, map[string][]bool{"out": {false}}, outputs)

				outputs, _ = hs.Evaluate(map[string][]bool{"a": {false}, "b": {true}})
				assert.Equal(t, map[string][]bool{"out": {false}}, outputs)

				outputs, _ = hs.Evaluate(map[string][]bool{"a": {true}, "b": {false}})
				assert.Equal(t, map[string][]bool{"out": {false}}, outputs)

				outputs, _ = hs.Evaluate(map[string][]bool{"a": {true}, "b": {true}})
				assert.Equal(t, map[string][]bool{"out": {true}}, outputs)
			},
		},
		{
			name:         "Or Chip",
			chipFileName: "OrChip",
			hdls: map[string]string{
				"OrChip": `CHIP OrChip {
					IN a, b;
					OUT out;

					PARTS:
					NotChip(in = a, out = notA);
					NotChip(in = b, out = notB);
					Nand(a = notA, b = notB, out = out);
				}`,
				"NotChip": testutils.ChipImplementations["NotChip"],
			},
			expectedInputsAfterProcess:  map[string]int{"a": 1, "b": 1},
			expectedOutputsAfterProcess: map[string]int{"out": 1},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				outputs, _ := hs.Evaluate(map[string][]bool{"a": {false}, "b": {false}})
				assert.Equal(t, map[string][]bool{"out": {false}}, outputs)

				outputs, _ = hs.Evaluate(map[string][]bool{"a": {false}, "b": {true}})
				assert.Equal(t, map[string][]bool{"out": {true}}, outputs)

				outputs, _ = hs.Evaluate(map[string][]bool{"a": {true}, "b": {false}})
				assert.Equal(t, map[string][]bool{"out": {true}}, outputs)

				outputs, _ = hs.Evaluate(map[string][]bool{"a": {true}, "b": {true}})
				assert.Equal(t, map[string][]bool{"out": {true}}, outputs)
			},
		},
		{
			name:         "Xor Chip",
			chipFileName: "XorChip",
			hdls: map[string]string{
				"XorChip": `CHIP XorChip {
					IN a, b;
					OUT out;

					PARTS:
					OrChip(a = a, b = b, out = AOrB);
					Nand(a = a, b = b, out = ANandB);
					AndChip(a = AOrB, b = ANandB, out = out);
				}`,
				"OrChip":  testutils.ChipImplementations["OrChip"],
				"AndChip": testutils.ChipImplementations["AndChip"],
				"NotChip": testutils.ChipImplementations["NotChip"],
			},
			expectedInputsAfterProcess:  map[string]int{"a": 1, "b": 1},
			expectedOutputsAfterProcess: map[string]int{"out": 1},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				outputs, _ := hs.Evaluate(map[string][]bool{"a": {false}, "b": {false}})
				assert.Equal(t, map[string][]bool{"out": {false}}, outputs)

				outputs, _ = hs.Evaluate(map[string][]bool{"a": {false}, "b": {true}})
				assert.Equal(t, map[string][]bool{"out": {true}}, outputs)

				outputs, _ = hs.Evaluate(map[string][]bool{"a": {true}, "b": {false}})
				assert.Equal(t, map[string][]bool{"out": {true}}, outputs)

				outputs, _ = hs.Evaluate(map[string][]bool{"a": {true}, "b": {true}})
				assert.Equal(t, map[string][]bool{"out": {false}}, outputs)
			},
		},
		{
			name:                        "Mux Chip",
			chipFileName:                "MuxChip",
			hdls:                        testutils.ChipImplementations,
			expectedInputsAfterProcess:  map[string]int{"a": 1, "b": 1, "sel": 1},
			expectedOutputsAfterProcess: map[string]int{"out": 1},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				outputs, _ := hs.Evaluate(map[string][]bool{"a": {false}, "b": {false}, "sel": {false}})
				assert.Equal(t, map[string][]bool{"out": {false}}, outputs)

				outputs, _ = hs.Evaluate(map[string][]bool{"a": {false}, "b": {false}, "sel": {true}})
				assert.Equal(t, map[string][]bool{"out": {false}}, outputs)

				outputs, _ = hs.Evaluate(map[string][]bool{"a": {false}, "b": {true}, "sel": {false}})
				assert.Equal(t, map[string][]bool{"out": {false}}, outputs)

				outputs, _ = hs.Evaluate(map[string][]bool{"a": {false}, "b": {true}, "sel": {true}})
				assert.Equal(t, map[string][]bool{"out": {true}}, outputs)

				outputs, _ = hs.Evaluate(map[string][]bool{"a": {true}, "b": {false}, "sel": {false}})
				assert.Equal(t, map[string][]bool{"out": {true}}, outputs)

				outputs, _ = hs.Evaluate(map[string][]bool{"a": {true}, "b": {false}, "sel": {true}})
				assert.Equal(t, map[string][]bool{"out": {false}}, outputs)

				outputs, _ = hs.Evaluate(map[string][]bool{"a": {true}, "b": {true}, "sel": {false}})
				assert.Equal(t, map[string][]bool{"out": {true}}, outputs)

				outputs, _ = hs.Evaluate(map[string][]bool{"a": {true}, "b": {true}, "sel": {true}})
				assert.Equal(t, map[string][]bool{"out": {true}}, outputs)
			},
		},
		{
			name:                        "DMux Chip",
			chipFileName:                "DMuxChip",
			hdls:                        testutils.ChipImplementations,
			expectedInputsAfterProcess:  map[string]int{"in": 1, "sel": 1},
			expectedOutputsAfterProcess: map[string]int{"a": 1, "b": 1},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				outputs, _ := hs.Evaluate(map[string][]bool{"in": {false}, "sel": {false}})
				assert.Equal(t, map[string][]bool{"a": {false}, "b": {false}}, outputs)

				outputs, _ = hs.Evaluate(map[string][]bool{"in": {false}, "sel": {true}})
				assert.Equal(t, map[string][]bool{"a": {false}, "b": {false}}, outputs)

				outputs, _ = hs.Evaluate(map[string][]bool{"in": {true}, "sel": {false}})
				assert.Equal(t, map[string][]bool{"a": {true}, "b": {false}}, outputs)

				outputs, _ = hs.Evaluate(map[string][]bool{"in": {true}, "sel": {true}})
				assert.Equal(t, map[string][]bool{"a": {false}, "b": {true}}, outputs)
			},
		},
		{
			name:                        "Not16 Chip",
			chipFileName:                "Not16Chip",
			hdls:                        testutils.ChipImplementations,
			expectedInputsAfterProcess:  map[string]int{"in": 16},
			expectedOutputsAfterProcess: map[string]int{"out": 16},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				input := testutils.RepeatBool(false, 16)
				outputs, _ := hs.Evaluate(map[string][]bool{"in": input})
				assert.Equal(t, map[string][]bool{"out": testutils.RepeatBool(true, 16)}, outputs)

				input = testutils.RepeatBool(true, 16)
				outputs, _ = hs.Evaluate(map[string][]bool{"in": input})
				assert.Equal(t, map[string][]bool{"out": testutils.RepeatBool(false, 16)}, outputs)

				input = []bool{true, false, true, false, true, false, true, false, true, false, true, false, true, false, true, false}
				outputs, _ = hs.Evaluate(map[string][]bool{"in": input})
				expectedOutput := []bool{false, true, false, true, false, true, false, true, false, true, false, true, false, true, false, true}
				assert.Equal(t, map[string][]bool{"out": expectedOutput}, outputs)

				input = []bool{false, false, true, true, true, true, false, false, true, true, false, false, false, false, true, true}
				outputs, _ = hs.Evaluate(map[string][]bool{"in": input})
				expectedOutput = []bool{true, true, false, false, false, false, true, true, false, false, true, true, true, true, false, false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput}, outputs)

				input = []bool{false, false, false, true, false, false, true, false, false, false, true, true, false, true, false, false}
				outputs, _ = hs.Evaluate(map[string][]bool{"in": input})
				expectedOutput = []bool{true, true, true, false, true, true, false, true, true, true, false, false, true, false, true, true}
				assert.Equal(t, map[string][]bool{"out": expectedOutput}, outputs)
			},
		},
		{
			name:                        "And16 Chip",
			chipFileName:                "And16Chip",
			hdls:                        testutils.ChipImplementations,
			expectedInputsAfterProcess:  map[string]int{"a": 16, "b": 16},
			expectedOutputsAfterProcess: map[string]int{"out": 16},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				a := testutils.RepeatBool(false, 16)
				b := testutils.RepeatBool(false, 16)
				outputs, _ := hs.Evaluate(map[string][]bool{"a": a, "b": b})
				assert.Equal(t, map[string][]bool{"out": testutils.RepeatBool(false, 16)}, outputs)

				a = testutils.RepeatBool(false, 16)
				b = testutils.RepeatBool(true, 16)
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b})
				assert.Equal(t, map[string][]bool{"out": testutils.RepeatBool(false, 16)}, outputs)

				a = testutils.RepeatBool(true, 16)
				b = testutils.RepeatBool(false, 16)
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b})
				assert.Equal(t, map[string][]bool{"out": testutils.RepeatBool(false, 16)}, outputs)

				a = testutils.RepeatBool(true, 16)
				b = testutils.RepeatBool(true, 16)
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b})
				assert.Equal(t, map[string][]bool{"out": testutils.RepeatBool(true, 16)}, outputs)

				a = []bool{true, false, true, false, true, false, true, false, true, false, true, false, true, false, true, false}
				b = []bool{false, true, false, true, false, true, false, true, false, true, false, true, false, true, false, true}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b})
				expectedOutput := testutils.RepeatBool(false, 16)
				assert.Equal(t, map[string][]bool{"out": expectedOutput}, outputs)

				// 0011110011000011
				a = []bool{false, false, true, true, true, true, false, false, true, true, false, false, false, false, true, true}
				// 0000111111110000
				b = []bool{false, false, false, false, true, true, true, true, true, true, true, true, false, false, false, false}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b})
				// 0000110011000000
				expectedOutput = []bool{false, false, false, false, true, true, false, false, true, true, false, false, false, false, false, false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput}, outputs)

				// 0001001000110100
				a = []bool{false, false, false, true, false, false, true, false, false, false, true, true, false, true, false, false}
				// 1001100001110110
				b = []bool{true, false, false, true, true, false, false, false, false, true, true, true, false, true, true, false}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b})
				// 0001000000110100
				expectedOutput = []bool{false, false, false, true, false, false, false, false, false, false, true, true, false, true, false, false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput}, outputs)
			},
		},
		{
			name:                        "Or16 Chip",
			chipFileName:                "Or16Chip",
			hdls:                        testutils.ChipImplementations,
			expectedInputsAfterProcess:  map[string]int{"a": 16, "b": 16},
			expectedOutputsAfterProcess: map[string]int{"out": 16},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				a := testutils.RepeatBool(false, 16)
				b := testutils.RepeatBool(false, 16)
				outputs, _ := hs.Evaluate(map[string][]bool{"a": a, "b": b})
				assert.Equal(t, map[string][]bool{"out": testutils.RepeatBool(false, 16)}, outputs)

				a = testutils.RepeatBool(false, 16)
				b = testutils.RepeatBool(true, 16)
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b})
				assert.Equal(t, map[string][]bool{"out": testutils.RepeatBool(true, 16)}, outputs)

				a = testutils.RepeatBool(true, 16)
				b = testutils.RepeatBool(false, 16)
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b})
				assert.Equal(t, map[string][]bool{"out": testutils.RepeatBool(true, 16)}, outputs)

				a = testutils.RepeatBool(true, 16)
				b = testutils.RepeatBool(true, 16)
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b})
				assert.Equal(t, map[string][]bool{"out": testutils.RepeatBool(true, 16)}, outputs)

				a = []bool{true, false, true, false, true, false, true, false, true, false, true, false, true, false, true, false}
				b = []bool{false, true, false, true, false, true, false, true, false, true, false, true, false, true, false, true}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b})
				expectedOutput := testutils.RepeatBool(true, 16)
				assert.Equal(t, map[string][]bool{"out": expectedOutput}, outputs)

				// 0011110011000011
				a = []bool{false, false, true, true, true, true, false, false, true, true, false, false, false, false, true, true}
				// 0000111111110000
				b = []bool{false, false, false, false, true, true, true, true, true, true, true, true, false, false, false, false}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b})
				// 0011111111110011
				expectedOutput = []bool{false, false, true, true, true, true, true, true, true, true, true, true, false, false, true, true}
				assert.Equal(t, map[string][]bool{"out": expectedOutput}, outputs)

				// 0001001000110100
				a = []bool{false, false, false, true, false, false, true, false, false, false, true, true, false, true, false, false}
				// 1001100001110110
				b = []bool{true, false, false, true, true, false, false, false, false, true, true, true, false, true, true, false}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b})
				// 1001101001110110
				expectedOutput = []bool{true, false, false, true, true, false, true, false, false, true, true, true, false, true, true, false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput}, outputs)
			},
		},
		{
			name:                        "Mux16 Chip",
			chipFileName:                "Mux16Chip",
			hdls:                        testutils.ChipImplementations,
			expectedInputsAfterProcess:  map[string]int{"a": 16, "b": 16, "sel": 1},
			expectedOutputsAfterProcess: map[string]int{"out": 16},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				a := testutils.RepeatBool(false, 16)
				b := testutils.RepeatBool(false, 16)
				sel := []bool{false}
				outputs, _ := hs.Evaluate(map[string][]bool{"a": a, "b": b, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": testutils.RepeatBool(false, 16)}, outputs)

				sel = []bool{true}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": testutils.RepeatBool(false, 16)}, outputs)

				a = testutils.RepeatBool(false, 16)
				b = []bool{false, false, false, true, false, false, true, false, false, false, true, true, false, true, false, false}
				sel = []bool{false}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": testutils.RepeatBool(false, 16)}, outputs)

				sel = []bool{true}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": b}, outputs)

				// 1001100001110110
				a = []bool{true, false, false, true, true, false, false, false, false, true, true, true, false, true, true, false}
				b = testutils.RepeatBool(false, 16)
				sel = []bool{false}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": a}, outputs)

				sel = []bool{true}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": testutils.RepeatBool(false, 16)}, outputs)

				// 1010101010101010
				a = []bool{true, false, true, false, true, false, true, false, true, false, true, false, true, false, true, false}
				// 0101010101010101
				b = []bool{false, true, false, true, false, true, false, true, false, true, false, true, false, true, false, true}
				sel = []bool{false}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": a}, outputs)

				sel = []bool{true}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": b}, outputs)
			},
		},
		{
			name:                        "Or8Way Chip",
			chipFileName:                "Or8WayChip",
			hdls:                        testutils.ChipImplementations,
			expectedInputsAfterProcess:  map[string]int{"in": 8},
			expectedOutputsAfterProcess: map[string]int{"out": 1},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				input := testutils.RepeatBool(false, 8)
				outputs, _ := hs.Evaluate(map[string][]bool{"in": input})
				assert.Equal(t, map[string][]bool{"out": {false}}, outputs)

				input = testutils.RepeatBool(true, 8)
				outputs, _ = hs.Evaluate(map[string][]bool{"in": input})
				assert.Equal(t, map[string][]bool{"out": {true}}, outputs)

				input = []bool{false, false, false, true, false, false, false, false}
				outputs, _ = hs.Evaluate(map[string][]bool{"in": input})
				assert.Equal(t, map[string][]bool{"out": {true}}, outputs)

				input = []bool{false, false, false, false, false, false, false, true}
				outputs, _ = hs.Evaluate(map[string][]bool{"in": input})
				assert.Equal(t, map[string][]bool{"out": {true}}, outputs)

				input = []bool{false, false, true, false, false, true, true, false}
				outputs, _ = hs.Evaluate(map[string][]bool{"in": input})
				assert.Equal(t, map[string][]bool{"out": {true}}, outputs)
			},
		},
		{
			name:                        "Mux4Way16 Chip",
			chipFileName:                "Mux4Way16Chip",
			hdls:                        testutils.ChipImplementations,
			expectedInputsAfterProcess:  map[string]int{"a": 16, "b": 16, "c": 16, "d": 16, "sel": 2},
			expectedOutputsAfterProcess: map[string]int{"out": 16},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				a := testutils.RepeatBool(false, 16)
				b := testutils.RepeatBool(false, 16)
				c := testutils.RepeatBool(false, 16)
				d := testutils.RepeatBool(false, 16)
				sel := []bool{false, false}
				outputs, _ := hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": testutils.RepeatBool(false, 16)}, outputs)

				sel = []bool{false, true}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": testutils.RepeatBool(false, 16)}, outputs)

				sel = []bool{true, false}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": testutils.RepeatBool(false, 16)}, outputs)

				sel = []bool{true, true}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": testutils.RepeatBool(false, 16)}, outputs)

				// 0001001000110100
				a = []bool{false, false, false, true, false, false, true, false, false, false, true, true, false, true, false, false}
				// 1001100001110110
				b = []bool{true, false, false, true, true, false, false, false, false, true, true, true, false, true, true, false}
				// 1010101010101010
				c = []bool{true, false, true, false, true, false, true, false, true, false, true, false, true, false, true, false}
				// 0101010101010101
				d = []bool{false, true, false, true, false, true, false, true, false, true, false, true, false, true, false, true}

				sel = []bool{false, false}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": a}, outputs)

				sel = []bool{true, false}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": b}, outputs)

				sel = []bool{false, true}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": c}, outputs)

				sel = []bool{true, true}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": d}, outputs)
			},
		},
		{
			name:                        "Mux8Way16 Chip",
			chipFileName:                "Mux8Way16Chip",
			hdls:                        testutils.ChipImplementations,
			expectedInputsAfterProcess:  map[string]int{"a": 16, "b": 16, "c": 16, "d": 16, "e": 16, "f": 16, "g": 16, "h": 16, "sel": 3},
			expectedOutputsAfterProcess: map[string]int{"out": 16},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				a := testutils.RepeatBool(false, 16)
				b := testutils.RepeatBool(false, 16)
				c := testutils.RepeatBool(false, 16)
				d := testutils.RepeatBool(false, 16)
				e := testutils.RepeatBool(false, 16)
				f := testutils.RepeatBool(false, 16)
				g := testutils.RepeatBool(false, 16)
				h := testutils.RepeatBool(false, 16)
				sel := []bool{false, false, false}
				outputs, _ := hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "e": e, "f": f, "g": g, "h": h, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": testutils.RepeatBool(false, 16)}, outputs)

				sel = []bool{true, false, false}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "e": e, "f": f, "g": g, "h": h, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": testutils.RepeatBool(false, 16)}, outputs)

				sel = []bool{false, true, false}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "e": e, "f": f, "g": g, "h": h, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": testutils.RepeatBool(false, 16)}, outputs)

				sel = []bool{true, true, false}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "e": e, "f": f, "g": g, "h": h, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": testutils.RepeatBool(false, 16)}, outputs)

				sel = []bool{false, false, true}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "e": e, "f": f, "g": g, "h": h, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": testutils.RepeatBool(false, 16)}, outputs)

				sel = []bool{true, false, true}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "e": e, "f": f, "g": g, "h": h, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": testutils.RepeatBool(false, 16)}, outputs)

				sel = []bool{false, true, true}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "e": e, "f": f, "g": g, "h": h, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": testutils.RepeatBool(false, 16)}, outputs)

				sel = []bool{true, true, true}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "e": e, "f": f, "g": g, "h": h, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": testutils.RepeatBool(false, 16)}, outputs)

				a = testutils.StringToBoolArray("0001001000110100")
				b = testutils.StringToBoolArray("0010001101000101")
				c = testutils.StringToBoolArray("0011010001010110")
				d = testutils.StringToBoolArray("0100010101100111")
				e = testutils.StringToBoolArray("0101011001111000")
				f = testutils.StringToBoolArray("0110011110001001")
				g = testutils.StringToBoolArray("0111100010011010")
				h = testutils.StringToBoolArray("1000100110101011")

				sel = testutils.StringToBoolArray("000")
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "e": e, "f": f, "g": g, "h": h, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": a}, outputs)

				sel = testutils.StringToBoolArray("001")
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "e": e, "f": f, "g": g, "h": h, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": b}, outputs)

				sel = testutils.StringToBoolArray("010")
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "e": e, "f": f, "g": g, "h": h, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": c}, outputs)

				sel = testutils.StringToBoolArray("011")
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "e": e, "f": f, "g": g, "h": h, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": d}, outputs)

				sel = testutils.StringToBoolArray("100")
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "e": e, "f": f, "g": g, "h": h, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": e}, outputs)

				sel = testutils.StringToBoolArray("101")
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "e": e, "f": f, "g": g, "h": h, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": f}, outputs)

				sel = testutils.StringToBoolArray("110")
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "e": e, "f": f, "g": g, "h": h, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": g}, outputs)

				sel = testutils.StringToBoolArray("111")
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "e": e, "f": f, "g": g, "h": h, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": h}, outputs)
			},
		},
		{
			name:                        "DMux4Way Chip",
			chipFileName:                "DMux4WayChip",
			hdls:                        testutils.ChipImplementations,
			expectedInputsAfterProcess:  map[string]int{"in": 1, "sel": 2},
			expectedOutputsAfterProcess: map[string]int{"a": 1, "b": 1, "c": 1, "d": 1},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				in := []bool{false}
				sel := testutils.StringToBoolArray("00")
				outputs, _ := hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				expectedOutputs := map[string][]bool{"a": {false}, "b": {false}, "c": {false}, "d": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				sel = testutils.StringToBoolArray("01")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				assert.Equal(t, expectedOutputs, outputs)

				sel = testutils.StringToBoolArray("10")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				assert.Equal(t, expectedOutputs, outputs)

				sel = testutils.StringToBoolArray("11")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				assert.Equal(t, expectedOutputs, outputs)

				in = []bool{true}

				sel = testutils.StringToBoolArray("00")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				expectedOutputs = map[string][]bool{"a": {true}, "b": {false}, "c": {false}, "d": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				sel = testutils.StringToBoolArray("01")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				expectedOutputs = map[string][]bool{"a": {false}, "b": {true}, "c": {false}, "d": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				sel = testutils.StringToBoolArray("10")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				expectedOutputs = map[string][]bool{"a": {false}, "b": {false}, "c": {true}, "d": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				sel = testutils.StringToBoolArray("11")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				expectedOutputs = map[string][]bool{"a": {false}, "b": {false}, "c": {false}, "d": {true}}
				assert.Equal(t, expectedOutputs, outputs)
			},
		},
		{
			name:                        "DMux8Way Chip",
			chipFileName:                "DMux8WayChip",
			hdls:                        testutils.ChipImplementations,
			expectedInputsAfterProcess:  map[string]int{"in": 1, "sel": 3},
			expectedOutputsAfterProcess: map[string]int{"a": 1, "b": 1, "c": 1, "d": 1, "e": 1, "f": 1, "g": 1, "h": 1},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				in := []bool{false}
				sel := testutils.StringToBoolArray("000")
				outputs, _ := hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				expectedOutputs := map[string][]bool{"a": {false}, "b": {false}, "c": {false}, "d": {false}, "e": {false}, "f": {false}, "g": {false}, "h": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				sel = testutils.StringToBoolArray("001")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				assert.Equal(t, expectedOutputs, outputs)

				sel = testutils.StringToBoolArray("010")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				assert.Equal(t, expectedOutputs, outputs)

				sel = testutils.StringToBoolArray("011")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				assert.Equal(t, expectedOutputs, outputs)

				sel = testutils.StringToBoolArray("100")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				assert.Equal(t, expectedOutputs, outputs)

				sel = testutils.StringToBoolArray("101")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				assert.Equal(t, expectedOutputs, outputs)

				sel = testutils.StringToBoolArray("110")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				assert.Equal(t, expectedOutputs, outputs)

				sel = testutils.StringToBoolArray("111")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				assert.Equal(t, expectedOutputs, outputs)

				in = []bool{true}
				sel = testutils.StringToBoolArray("000")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				expectedOutputs = map[string][]bool{"a": {true}, "b": {false}, "c": {false}, "d": {false}, "e": {false}, "f": {false}, "g": {false}, "h": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				sel = testutils.StringToBoolArray("001")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				expectedOutputs = map[string][]bool{"a": {false}, "b": {true}, "c": {false}, "d": {false}, "e": {false}, "f": {false}, "g": {false}, "h": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				sel = testutils.StringToBoolArray("010")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				expectedOutputs = map[string][]bool{"a": {false}, "b": {false}, "c": {true}, "d": {false}, "e": {false}, "f": {false}, "g": {false}, "h": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				sel = testutils.StringToBoolArray("011")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				expectedOutputs = map[string][]bool{"a": {false}, "b": {false}, "c": {false}, "d": {true}, "e": {false}, "f": {false}, "g": {false}, "h": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				sel = testutils.StringToBoolArray("100")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				expectedOutputs = map[string][]bool{"a": {false}, "b": {false}, "c": {false}, "d": {false}, "e": {true}, "f": {false}, "g": {false}, "h": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				sel = testutils.StringToBoolArray("101")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				expectedOutputs = map[string][]bool{"a": {false}, "b": {false}, "c": {false}, "d": {false}, "e": {false}, "f": {true}, "g": {false}, "h": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				sel = testutils.StringToBoolArray("110")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				expectedOutputs = map[string][]bool{"a": {false}, "b": {false}, "c": {false}, "d": {false}, "e": {false}, "f": {false}, "g": {true}, "h": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				sel = testutils.StringToBoolArray("111")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				expectedOutputs = map[string][]bool{"a": {false}, "b": {false}, "c": {false}, "d": {false}, "e": {false}, "f": {false}, "g": {false}, "h": {true}}
				assert.Equal(t, expectedOutputs, outputs)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			hs := New()
			hs.SetChipHDLs(tt.hdls)
			inputs, outputs, _, err := hs.Process(tt.chipFileName)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			assert.Equal(t, len(tt.expectedInputsAfterProcess), len(inputs), "number of inputs mismatch")
			for inputName, width := range tt.expectedInputsAfterProcess {
				assert.NotNil(t, inputs[inputName], "expected input %s to be present", inputName)
				assert.Equal(t, width, inputs[inputName], "expected input %s to have width %d", inputName, width)
			}

			assert.Equal(t, len(tt.expectedOutputsAfterProcess), len(outputs), "number of outputs mismatch")
			for outputName, width := range tt.expectedOutputsAfterProcess {
				assert.NotNil(t, outputs[outputName], "expected output %s to be present", outputName)
				assert.Equal(t, width, outputs[outputName], "expected output %s to have width %d", outputName, width)
			}

			if tt.afterProcess != nil {
				tt.afterProcess(t, hs)
			}
		})
	}
}
