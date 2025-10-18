package simulator

import (
	"testing"

	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/testutils"
	"github.com/stretchr/testify/assert"
)

func TestAdditionalChipsSimulation(t *testing.T) {
	tests := []struct {
		name                        string
		chipFileName                string
		hdls                        map[string]string
		expectedInputsAfterProcess  map[string]int
		expectedOutputsAfterProcess map[string]int
		afterProcess                func(t *testing.T, hs *HardwareSimulator)
	}{
		{
			name:                        "Double DFF Chip",
			chipFileName:                "DoubleDFFChip",
			hdls:                        testutils.ChipImplementations,
			expectedInputsAfterProcess:  map[string]int{"in": 1},
			expectedOutputsAfterProcess: map[string]int{"dff1": 1, "dff2": 1},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				in := []bool{true}
				outputs, internalPins := hs.Tick(map[string][]bool{"in": in})
				expectedOutputs := map[string][]bool{"dff1": {false}, "dff2": {false}}
				assert.Equal(t, expectedOutputs, outputs)
				expectedInternalPins := map[string][]bool{"dff1internal": {false}}
				assert.Equal(t, expectedInternalPins, internalPins)

				outputs, internalPins = hs.Tock(map[string][]bool{"in": in})
				expectedOutputs = map[string][]bool{"dff1": {true}, "dff2": {false}}
				assert.Equal(t, expectedOutputs, outputs)
				expectedInternalPins = map[string][]bool{"dff1internal": {true}}
				assert.Equal(t, expectedInternalPins, internalPins)

				outputs, internalPins = hs.Tick(map[string][]bool{"in": in})
				expectedOutputs = map[string][]bool{"dff1": {true}, "dff2": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, internalPins = hs.Tock(map[string][]bool{"in": in})
				expectedOutputs = map[string][]bool{"dff1": {true}, "dff2": {true}}
				assert.Equal(t, expectedOutputs, outputs)

				in = []bool{false}
				outputs, internalPins = hs.Tick(map[string][]bool{"in": in})
				expectedOutputs = map[string][]bool{"dff1": {true}, "dff2": {true}}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, internalPins = hs.Tock(map[string][]bool{"in": in})
				expectedOutputs = map[string][]bool{"dff1": {false}, "dff2": {true}}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, internalPins = hs.Tick(map[string][]bool{"in": in})
				expectedOutputs = map[string][]bool{"dff1": {false}, "dff2": {true}}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, internalPins = hs.Tock(map[string][]bool{"in": in})
				expectedOutputs = map[string][]bool{"dff1": {false}, "dff2": {false}}
				assert.Equal(t, expectedOutputs, outputs)
			},
		},
		{
			name:                        "Test Chip",
			chipFileName:                "TestChip",
			hdls:                        testutils.ChipImplementations,
			expectedInputsAfterProcess:  map[string]int{"in": 1},
			expectedOutputsAfterProcess: map[string]int{"out": 1, "outnot": 1},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				in := []bool{true}
				hs.Evaluator.SetInputs(map[string][]bool{"in": in})
				hs.Evaluator.Evaluate()
				outputs, internalPins := hs.Evaluator.GetOutputsAndInternalPins()
				expectedOutputs := map[string][]bool{"out": {false}, "outnot": {true}}
				assert.Equal(t, expectedOutputs, outputs)
				expectedInternalPins := map[string][]bool{"dffout": {false}}
				assert.Equal(t, expectedInternalPins, internalPins)

				hs.Evaluator.Commit()
				outputs, internalPins = hs.Evaluator.GetOutputsAndInternalPins()
				expectedOutputs = map[string][]bool{"out": {false}, "outnot": {true}}
				assert.Equal(t, expectedOutputs, outputs)
				expectedInternalPins = map[string][]bool{"dffout": {false}}
				assert.Equal(t, expectedInternalPins, internalPins)

				hs.Evaluator.Evaluate()
				outputs, internalPins = hs.Evaluator.GetOutputsAndInternalPins()
				expectedOutputs = map[string][]bool{"out": {true}, "outnot": {false}}
				assert.Equal(t, expectedOutputs, outputs)
				expectedInternalPins = map[string][]bool{"dffout": {true}}
				assert.Equal(t, expectedInternalPins, internalPins)

				// outputs, internalPins = hs.Tock(map[string][]bool{"in": in})
				// expectedOutputs = map[string][]bool{"out": {true}, "outnot": {false}}
				// assert.Equal(t, expectedOutputs, outputs)
				// expectedInternalPins = map[string][]bool{"dffout": {true}}
				// assert.Equal(t, expectedInternalPins, internalPins)
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
