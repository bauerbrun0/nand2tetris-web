package simulator

import (
	"testing"

	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/testutils"
	"github.com/stretchr/testify/assert"
)

func TestProject2ChipsSimulation(t *testing.T) {
	tests := []struct {
		name                        string
		chipFileName                string
		hdls                        map[string]string
		expectedInputsAfterProcess  map[string]int
		expectedOutputsAfterProcess map[string]int
		afterProcess                func(t *testing.T, hs *HardwareSimulator)
	}{
		{
			name:                        "HalfAdder Chip",
			chipFileName:                "HalfAdderChip",
			hdls:                        testutils.ChipImplementations,
			expectedInputsAfterProcess:  map[string]int{"a": 1, "b": 1},
			expectedOutputsAfterProcess: map[string]int{"sum": 1, "carry": 1},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				a := []bool{false}
				b := []bool{false}
				outputs, _ := hs.Evaluate(map[string][]bool{"a": a, "b": b})
				expectedOutputs := map[string][]bool{"sum": {false}, "carry": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				a = []bool{false}
				b = []bool{true}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b})
				expectedOutputs = map[string][]bool{"sum": {true}, "carry": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				a = []bool{true}
				b = []bool{false}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b})
				expectedOutputs = map[string][]bool{"sum": {true}, "carry": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				a = []bool{true}
				b = []bool{true}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b})
				expectedOutputs = map[string][]bool{"sum": {false}, "carry": {true}}
				assert.Equal(t, expectedOutputs, outputs)
			},
		},
		{
			name:                        "FullAdder Chip",
			chipFileName:                "FullAdderChip",
			hdls:                        testutils.ChipImplementations,
			expectedInputsAfterProcess:  map[string]int{"a": 1, "b": 1, "c": 1},
			expectedOutputsAfterProcess: map[string]int{"sum": 1, "carry": 1},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				a := []bool{false}
				b := []bool{false}
				c := []bool{false}
				outputs, _ := hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c})
				expectedOutputs := map[string][]bool{"sum": {false}, "carry": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				a = []bool{false}
				b = []bool{false}
				c = []bool{true}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c})
				expectedOutputs = map[string][]bool{"sum": {true}, "carry": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				a = []bool{false}
				b = []bool{true}
				c = []bool{false}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c})
				expectedOutputs = map[string][]bool{"sum": {true}, "carry": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				a = []bool{false}
				b = []bool{true}
				c = []bool{true}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c})
				expectedOutputs = map[string][]bool{"sum": {false}, "carry": {true}}
				assert.Equal(t, expectedOutputs, outputs)

				a = []bool{true}
				b = []bool{false}
				c = []bool{false}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c})
				expectedOutputs = map[string][]bool{"sum": {true}, "carry": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				a = []bool{true}
				b = []bool{false}
				c = []bool{true}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c})
				expectedOutputs = map[string][]bool{"sum": {false}, "carry": {true}}
				assert.Equal(t, expectedOutputs, outputs)

				a = []bool{true}
				b = []bool{true}
				c = []bool{false}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c})
				expectedOutputs = map[string][]bool{"sum": {false}, "carry": {true}}
				assert.Equal(t, expectedOutputs, outputs)

				a = []bool{true}
				b = []bool{true}
				c = []bool{true}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c})
				expectedOutputs = map[string][]bool{"sum": {true}, "carry": {true}}
				assert.Equal(t, expectedOutputs, outputs)
			},
		},
		{
			name:                        "Add16 Chip",
			chipFileName:                "Add16Chip",
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
				b = testutils.RepeatBool(true, 16)
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b})
				expectedOutput := testutils.StringToBoolArray("1111111111111110")
				assert.Equal(t, map[string][]bool{"out": expectedOutput}, outputs)

				a = testutils.StringToBoolArray("1010101010101010")
				b = testutils.StringToBoolArray("0101010101010101")
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b})
				expectedOutput = testutils.StringToBoolArray("1111111111111111")
				assert.Equal(t, map[string][]bool{"out": expectedOutput}, outputs)

				a = testutils.StringToBoolArray("0011110011000011")
				b = testutils.StringToBoolArray("0000111111110000")
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b})
				expectedOutput = testutils.StringToBoolArray("0100110010110011")
				assert.Equal(t, map[string][]bool{"out": expectedOutput}, outputs)

				a = testutils.StringToBoolArray("0001001000110100")
				b = testutils.StringToBoolArray("1001100001110110")
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b})
				expectedOutput = testutils.StringToBoolArray("1010101010101010")
				assert.Equal(t, map[string][]bool{"out": expectedOutput}, outputs)
			},
		},
		{
			name:                        "Inc16 Chip",
			chipFileName:                "Inc16Chip",
			hdls:                        testutils.ChipImplementations,
			expectedInputsAfterProcess:  map[string]int{"in": 16},
			expectedOutputsAfterProcess: map[string]int{"out": 16},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				in := testutils.RepeatBool(false, 16)
				outputs, _ := hs.Evaluate(map[string][]bool{"in": in})
				expectedOutput := testutils.StringToBoolArray("0000000000000001")
				assert.Equal(t, map[string][]bool{"out": expectedOutput}, outputs)

				in = testutils.StringToBoolArray("1111111111111111")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in})
				expectedOutput = testutils.StringToBoolArray("0000000000000000")
				assert.Equal(t, map[string][]bool{"out": expectedOutput}, outputs)

				in = testutils.StringToBoolArray("0000000000000101")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in})
				expectedOutput = testutils.StringToBoolArray("0000000000000110")
				assert.Equal(t, map[string][]bool{"out": expectedOutput}, outputs)

				in = testutils.StringToBoolArray("1111111111111011")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in})
				expectedOutput = testutils.StringToBoolArray("1111111111111100")
				assert.Equal(t, map[string][]bool{"out": expectedOutput}, outputs)
			},
		},
		{
			name:                        "ALU Chip",
			chipFileName:                "ALUChip",
			hdls:                        testutils.ChipImplementations,
			expectedInputsAfterProcess:  map[string]int{"x": 16, "y": 16, "zx": 1, "nx": 1, "zy": 1, "ny": 1, "f": 1, "no": 1},
			expectedOutputsAfterProcess: map[string]int{"out": 16, "zr": 1, "ng": 1},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				x := testutils.StringToBoolArray("0000000000000000")
				y := testutils.StringToBoolArray("1111111111111111")

				zx, nx, zy, ny, f, no := getALUFlagInputs("101010")
				outputs, _ := hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput := testutils.StringToBoolArray("0000000000000000")
				expectedZR := []bool{true}
				expectedNG := []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("111111")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = testutils.StringToBoolArray("0000000000000001")
				expectedZR = []bool{false}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("111010")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = testutils.StringToBoolArray("1111111111111111")
				expectedZR = []bool{false}
				expectedNG = []bool{true}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("001100")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = testutils.StringToBoolArray("0000000000000000")
				expectedZR = []bool{true}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("110000")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = testutils.StringToBoolArray("1111111111111111")
				expectedZR = []bool{false}
				expectedNG = []bool{true}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("001101")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = testutils.StringToBoolArray("1111111111111111")
				expectedZR = []bool{false}
				expectedNG = []bool{true}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("110001")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = testutils.StringToBoolArray("0000000000000000")
				expectedZR = []bool{true}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("001111")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = testutils.StringToBoolArray("0000000000000000")
				expectedZR = []bool{true}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("110011")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = testutils.StringToBoolArray("0000000000000001")
				expectedZR = []bool{false}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("011111")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = testutils.StringToBoolArray("0000000000000001")
				expectedZR = []bool{false}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("110111")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = testutils.StringToBoolArray("0000000000000000")
				expectedZR = []bool{true}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("001110")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = testutils.StringToBoolArray("1111111111111111")
				expectedZR = []bool{false}
				expectedNG = []bool{true}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("110010")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = testutils.StringToBoolArray("1111111111111110")
				expectedZR = []bool{false}
				expectedNG = []bool{true}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("000010")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = testutils.StringToBoolArray("1111111111111111")
				expectedZR = []bool{false}
				expectedNG = []bool{true}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("010011")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = testutils.StringToBoolArray("0000000000000001")
				expectedZR = []bool{false}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("000111")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = testutils.StringToBoolArray("1111111111111111")
				expectedZR = []bool{false}
				expectedNG = []bool{true}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("000000")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = testutils.StringToBoolArray("0000000000000000")
				expectedZR = []bool{true}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("010101")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = testutils.StringToBoolArray("1111111111111111")
				expectedZR = []bool{false}
				expectedNG = []bool{true}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				x = testutils.StringToBoolArray("0000000000010001")
				y = testutils.StringToBoolArray("0000000000000011")

				zx, nx, zy, ny, f, no = getALUFlagInputs("101010")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = testutils.StringToBoolArray("0000000000000000")
				expectedZR = []bool{true}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("111111")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = testutils.StringToBoolArray("0000000000000001")
				expectedZR = []bool{false}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("111010")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = testutils.StringToBoolArray("1111111111111111")
				expectedZR = []bool{false}
				expectedNG = []bool{true}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("001100")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = testutils.StringToBoolArray("0000000000010001")
				expectedZR = []bool{false}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("110000")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = testutils.StringToBoolArray("0000000000000011")
				expectedZR = []bool{false}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("001101")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = testutils.StringToBoolArray("1111111111101110")
				expectedZR = []bool{false}
				expectedNG = []bool{true}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("110001")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = testutils.StringToBoolArray("1111111111111100")
				expectedZR = []bool{false}
				expectedNG = []bool{true}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("001111")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = testutils.StringToBoolArray("1111111111101111")
				expectedZR = []bool{false}
				expectedNG = []bool{true}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("110011")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = testutils.StringToBoolArray("1111111111111101")
				expectedZR = []bool{false}
				expectedNG = []bool{true}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("011111")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = testutils.StringToBoolArray("0000000000010010")
				expectedZR = []bool{false}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("110111")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = testutils.StringToBoolArray("0000000000000100")
				expectedZR = []bool{false}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("001110")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = testutils.StringToBoolArray("0000000000010000")
				expectedZR = []bool{false}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("110010")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = testutils.StringToBoolArray("0000000000000010")
				expectedZR = []bool{false}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("000010")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = testutils.StringToBoolArray("0000000000010100")
				expectedZR = []bool{false}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("010011")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = testutils.StringToBoolArray("0000000000001110")
				expectedZR = []bool{false}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("000111")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = testutils.StringToBoolArray("1111111111110010")
				expectedZR = []bool{false}
				expectedNG = []bool{true}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("000000")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = testutils.StringToBoolArray("0000000000000001")
				expectedZR = []bool{false}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("010101")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = testutils.StringToBoolArray("0000000000010011")
				expectedZR = []bool{false}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			hs := New()
			hs.SetChipHDLs(tt.hdls)
			inputs, outputs, err := hs.Process(tt.chipFileName)
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
