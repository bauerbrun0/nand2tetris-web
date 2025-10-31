package simulator

import (
	"testing"

	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/testutils"
	"github.com/stretchr/testify/assert"
)

func TestProject5ChipsSimulation(t *testing.T) {
	tests := []struct {
		name                        string
		chipFileName                string
		hdls                        map[string]string
		expectedInputsAfterProcess  map[string]int
		expectedOutputsAfterProcess map[string]int
		afterProcess                func(t *testing.T, hs *HardwareSimulator)
	}{
		{
			name:                        "CPU Chip",
			chipFileName:                "CPUChip",
			hdls:                        testutils.ChipImplementations,
			expectedInputsAfterProcess:  map[string]int{"inM": 16, "instruction": 16, "reset": 1},
			expectedOutputsAfterProcess: map[string]int{"outM": 16, "writeM": 1, "addressM": 15, "pc": 15},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				inM := testutils.RepeatBool(false, 16)
				instruction := testutils.StringToBoolArray("0011000000111001")
				reset := []bool{false}
				outputs, internalPins := hs.Tick(map[string][]bool{"inM": inM, "instruction": instruction, "reset": reset})

				expectedOutM := testutils.StringToBoolArray("0000000000000000")
				expectedWriteM := []bool{false}
				expectedAddressM := testutils.StringToBoolArray("000000000000000")
				expectedPC := testutils.StringToBoolArray("000000000000000")
				expectedRegDOut := testutils.StringToBoolArray("0000000000000000")

				assert.Equal(t, expectedOutM, outputs["outM"], "outM output mismatch")
				assert.Equal(t, expectedWriteM, outputs["writeM"], "writeM output mismatch")
				assert.Equal(t, expectedAddressM, outputs["addressM"], "addressM output mismatch")
				assert.Equal(t, expectedPC, outputs["pc"], "pc output mismatch")
				assert.Equal(t, expectedRegDOut, internalPins["regDOut"], "regDOut internal pin mismatch")

				outputs, internalPins = hs.Tock(map[string][]bool{"inM": inM, "instruction": instruction, "reset": reset})

				// expected addressM is 12345
				expectedAddressM = testutils.StringToBoolArray("011000000111001")
				expectedPC = testutils.StringToBoolArray("000000000000001")

				assert.Equal(t, expectedOutM, outputs["outM"], "outM output mismatch")
				assert.Equal(t, expectedWriteM, outputs["writeM"], "writeM output mismatch")
				assert.Equal(t, expectedAddressM, outputs["addressM"], "addressM output mismatch")
				assert.Equal(t, expectedPC, outputs["pc"], "pc output mismatch")
				assert.Equal(t, expectedRegDOut, internalPins["regDOut"], "regDOut internal pin mismatch")

				// further tests can be added based on the Nand2Tetris CPU specification/test script
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
