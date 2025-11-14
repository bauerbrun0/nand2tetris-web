package simulator

import (
	"testing"
	"time"

	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/testutils"
	"github.com/stretchr/testify/assert"
)

func TestProject3ChipsSimulation(t *testing.T) {
	tests := []struct {
		name                        string
		chipFileName                string
		hdls                        map[string]string
		expectedInputsAfterProcess  map[string]int
		expectedOutputsAfterProcess map[string]int
		afterProcess                func(t *testing.T, hs *HardwareSimulator)
	}{
		{
			name:                        "Bit Chip",
			chipFileName:                "BitChip",
			hdls:                        testutils.ChipImplementations,
			expectedInputsAfterProcess:  map[string]int{"in": 1, "load": 1},
			expectedOutputsAfterProcess: map[string]int{"out": 1},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				in := []bool{false}
				load := []bool{false}
				outputs, _ := hs.Tick(map[string][]bool{"in": in, "load": load})
				expectedOutputs := map[string][]bool{"out": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load})
				expectedOutputs = map[string][]bool{"out": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				in = []bool{false}
				load = []bool{true}
				outputs, _ = hs.Tick(map[string][]bool{"in": in, "load": load})
				expectedOutputs = map[string][]bool{"out": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load})
				expectedOutputs = map[string][]bool{"out": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				in = []bool{true}
				load = []bool{false}
				outputs, _ = hs.Tick(map[string][]bool{"in": in, "load": load})
				expectedOutputs = map[string][]bool{"out": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load})
				expectedOutputs = map[string][]bool{"out": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				in = []bool{true}
				load = []bool{true}
				outputs, _ = hs.Tick(map[string][]bool{"in": in, "load": load})
				expectedOutputs = map[string][]bool{"out": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load})
				expectedOutputs = map[string][]bool{"out": {true}}
				assert.Equal(t, expectedOutputs, outputs)

				in = []bool{false}
				load = []bool{false}
				outputs, _ = hs.Tick(map[string][]bool{"in": in, "load": load})
				expectedOutputs = map[string][]bool{"out": {true}}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load})
				expectedOutputs = map[string][]bool{"out": {true}}
				assert.Equal(t, expectedOutputs, outputs)

				in = []bool{true}
				load = []bool{false}
				outputs, _ = hs.Tick(map[string][]bool{"in": in, "load": load})
				expectedOutputs = map[string][]bool{"out": {true}}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load})
				expectedOutputs = map[string][]bool{"out": {true}}
				assert.Equal(t, expectedOutputs, outputs)

				in = []bool{false}
				load = []bool{true}
				outputs, _ = hs.Tick(map[string][]bool{"in": in, "load": load})
				expectedOutputs = map[string][]bool{"out": {true}}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load})
				expectedOutputs = map[string][]bool{"out": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				in = []bool{true}
				load = []bool{true}
				outputs, _ = hs.Tick(map[string][]bool{"in": in, "load": load})
				expectedOutputs = map[string][]bool{"out": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load})
				expectedOutputs = map[string][]bool{"out": {true}}
				assert.Equal(t, expectedOutputs, outputs)

				in = []bool{false}
				load = []bool{false}
				for _ = range 49 {
					outputs, _ = hs.Tick(map[string][]bool{"in": in, "load": load})
					expectedOutputs = map[string][]bool{"out": {true}}
					assert.Equal(t, expectedOutputs, outputs)

					outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load})
					expectedOutputs = map[string][]bool{"out": {true}}
					assert.Equal(t, expectedOutputs, outputs)
				}

				in = []bool{false}
				load = []bool{true}
				outputs, _ = hs.Tick(map[string][]bool{"in": in, "load": load})
				expectedOutputs = map[string][]bool{"out": {true}}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load})
				expectedOutputs = map[string][]bool{"out": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				in = []bool{true}
				load = []bool{false}
				for _ = range 49 {
					outputs, _ = hs.Tick(map[string][]bool{"in": in, "load": load})
					expectedOutputs = map[string][]bool{"out": {false}}
					assert.Equal(t, expectedOutputs, outputs)

					outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load})
					expectedOutputs = map[string][]bool{"out": {false}}
					assert.Equal(t, expectedOutputs, outputs)
				}
			},
		},
		{
			name:                        "Register Chip",
			chipFileName:                "RegisterChip",
			hdls:                        testutils.ChipImplementations,
			expectedInputsAfterProcess:  map[string]int{"in": 16, "load": 1},
			expectedOutputsAfterProcess: map[string]int{"out": 16},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				in := testutils.StringToBoolArray("0000000000000000")
				load := []bool{false}
				outputs, _ := hs.Tick(map[string][]bool{"in": in, "load": load})
				expectedOutputs := map[string][]bool{"out": testutils.StringToBoolArray("0000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)

				in = testutils.StringToBoolArray("1010101010101010")
				load = []bool{true}
				outputs, _ = hs.Tick(map[string][]bool{"in": in, "load": load})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("1010101010101010")}
				assert.Equal(t, expectedOutputs, outputs)

				in = testutils.StringToBoolArray("0101010101010101")
				load = []bool{false}
				outputs, _ = hs.Tick(map[string][]bool{"in": in, "load": load})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("1010101010101010")}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("1010101010101010")}
				assert.Equal(t, expectedOutputs, outputs)

				in = testutils.StringToBoolArray("1111111111111111")
				load = []bool{true}
				outputs, _ = hs.Tick(map[string][]bool{"in": in, "load": load})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("1010101010101010")}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("1111111111111111")}
				assert.Equal(t, expectedOutputs, outputs)
			},
		},
		{
			name:                        "RAM8 Chip",
			chipFileName:                "RAM8Chip",
			hdls:                        testutils.ChipImplementations,
			expectedInputsAfterProcess:  map[string]int{"in": 16, "load": 1, "address": 3},
			expectedOutputsAfterProcess: map[string]int{"out": 16},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				in := testutils.StringToBoolArray("0000000000000001")
				load := []bool{false}
				address := testutils.StringToBoolArray("000")

				outputs, _ := hs.Tick(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs := map[string][]bool{"out": testutils.StringToBoolArray("0000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)

				load = []bool{true}
				outputs, _ = hs.Tick(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000001")}
				assert.Equal(t, expectedOutputs, outputs)

				in = testutils.StringToBoolArray("1000000000000000")
				load = []bool{true}
				address = testutils.StringToBoolArray("001")
				outputs, _ = hs.Tick(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("1000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)

				load = []bool{false}
				address = testutils.StringToBoolArray("000")
				outputs, _ = hs.Tick(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000001")}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000001")}
				assert.Equal(t, expectedOutputs, outputs)

				address = testutils.StringToBoolArray("001")
				outputs, _ = hs.Tick(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("1000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("1000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)
			},
		},
		{
			name:                        "RAM64 Chip",
			chipFileName:                "RAM64Chip",
			hdls:                        testutils.ChipImplementations,
			expectedInputsAfterProcess:  map[string]int{"in": 16, "load": 1, "address": 6},
			expectedOutputsAfterProcess: map[string]int{"out": 16},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				in := testutils.StringToBoolArray("0000000000000001")
				load := []bool{false}
				address := testutils.StringToBoolArray("000000")

				outputs, _ := hs.Tick(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs := map[string][]bool{"out": testutils.StringToBoolArray("0000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)

				load = []bool{true}
				outputs, _ = hs.Tick(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000001")}
				assert.Equal(t, expectedOutputs, outputs)

				in = testutils.StringToBoolArray("1000000000000000")
				load = []bool{true}
				address = testutils.StringToBoolArray("111111")
				outputs, _ = hs.Tick(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("1000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)

				load = []bool{false}
				address = testutils.StringToBoolArray("000000")
				outputs, _ = hs.Tick(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000001")}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000001")}
				assert.Equal(t, expectedOutputs, outputs)

				address = testutils.StringToBoolArray("111111")
				outputs, _ = hs.Tick(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("1000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("1000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)
			},
		},
		{
			name:                        "RAM512 Chip",
			chipFileName:                "RAM512Chip",
			hdls:                        testutils.ChipImplementations,
			expectedInputsAfterProcess:  map[string]int{"in": 16, "load": 1, "address": 9},
			expectedOutputsAfterProcess: map[string]int{"out": 16},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				in := testutils.StringToBoolArray("0000000000000001")
				load := []bool{false}
				address := testutils.StringToBoolArray("000000000")

				outputs, _ := hs.Tick(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs := map[string][]bool{"out": testutils.StringToBoolArray("0000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)

				load = []bool{true}
				outputs, _ = hs.Tick(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000001")}
				assert.Equal(t, expectedOutputs, outputs)

				in = testutils.StringToBoolArray("1000000000000000")
				load = []bool{true}
				address = testutils.StringToBoolArray("111111111")
				outputs, _ = hs.Tick(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("1000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)

				load = []bool{false}
				address = testutils.StringToBoolArray("000000000")
				outputs, _ = hs.Tick(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000001")}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000001")}
				assert.Equal(t, expectedOutputs, outputs)

				address = testutils.StringToBoolArray("111111111")
				outputs, _ = hs.Tick(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("1000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("1000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)
			},
		},
		{
			name:                        "RAM4K Chip",
			chipFileName:                "RAM4KChip",
			hdls:                        testutils.ChipImplementations,
			expectedInputsAfterProcess:  map[string]int{"in": 16, "load": 1, "address": 12},
			expectedOutputsAfterProcess: map[string]int{"out": 16},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				start := time.Now()
				in := testutils.StringToBoolArray("0000000000000001")
				load := []bool{false}
				address := testutils.StringToBoolArray("000000000000")

				outputs, _ := hs.Tick(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs := map[string][]bool{"out": testutils.StringToBoolArray("0000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)

				t.Logf("RAM4K first tick/tock took %v\n", time.Since(start))
				start = time.Now()

				load = []bool{true}
				outputs, _ = hs.Tick(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000001")}
				assert.Equal(t, expectedOutputs, outputs)
				t.Logf("RAM4K second tick/tock took %v\n", time.Since(start))

				in = testutils.StringToBoolArray("1000000000000000")
				load = []bool{true}
				address = testutils.StringToBoolArray("111111111111")
				outputs, _ = hs.Tick(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("1000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)

				load = []bool{false}
				address = testutils.StringToBoolArray("000000000000")
				outputs, _ = hs.Tick(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000001")}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000001")}
				assert.Equal(t, expectedOutputs, outputs)

				address = testutils.StringToBoolArray("111111111111")
				outputs, _ = hs.Tick(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("1000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("1000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)
			},
		},
		{
			name:                        "RAM16K Chip",
			chipFileName:                "RAM16KChip",
			hdls:                        testutils.ChipImplementations,
			expectedInputsAfterProcess:  map[string]int{"in": 16, "load": 1, "address": 14},
			expectedOutputsAfterProcess: map[string]int{"out": 16},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				start := time.Now()
				in := testutils.StringToBoolArray("0000000000000001")
				load := []bool{false}
				address := testutils.StringToBoolArray("00000000000000")

				outputs, _ := hs.Tick(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs := map[string][]bool{"out": testutils.StringToBoolArray("0000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)

				t.Logf("RAM16K first tick/tock took %v\n", time.Since(start))
				start = time.Now()

				load = []bool{true}
				outputs, _ = hs.Tick(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000001")}
				assert.Equal(t, expectedOutputs, outputs)
				t.Logf("RAM16K second tick/tock took %v\n", time.Since(start))

				in = testutils.StringToBoolArray("1000000000000000")
				load = []bool{true}
				address = testutils.StringToBoolArray("11111111111111")
				outputs, _ = hs.Tick(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("1000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)

				load = []bool{false}
				address = testutils.StringToBoolArray("00000000000000")
				outputs, _ = hs.Tick(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000001")}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000001")}
				assert.Equal(t, expectedOutputs, outputs)

				address = testutils.StringToBoolArray("11111111111111")
				outputs, _ = hs.Tick(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("1000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load, "address": address})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("1000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)
			},
		},
		{
			name:                        "PC Chip",
			chipFileName:                "PCChip",
			hdls:                        testutils.ChipImplementations,
			expectedInputsAfterProcess:  map[string]int{"in": 16, "load": 1, "reset": 1, "inc": 1},
			expectedOutputsAfterProcess: map[string]int{"out": 16},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				in := testutils.StringToBoolArray("0000000000100000")
				load := []bool{false}
				reset := []bool{false}
				inc := []bool{false}

				outputs, _ := hs.Tick(map[string][]bool{"in": in, "load": load, "reset": reset, "inc": inc})
				expectedOutputs := map[string][]bool{"out": testutils.StringToBoolArray("0000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load, "reset": reset, "inc": inc})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)

				inc = []bool{true}
				outputs, _ = hs.Tick(map[string][]bool{"in": in, "load": load, "reset": reset, "inc": inc})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load, "reset": reset, "inc": inc})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000001")}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tick(map[string][]bool{"in": in, "load": load, "reset": reset, "inc": inc})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000001")}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load, "reset": reset, "inc": inc})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000010")}
				assert.Equal(t, expectedOutputs, outputs)

				reset = []bool{true}
				outputs, _ = hs.Tick(map[string][]bool{"in": in, "load": load, "reset": reset, "inc": inc})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000010")}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load, "reset": reset, "inc": inc})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)

				reset = []bool{false}
				outputs, _ = hs.Tick(map[string][]bool{"in": in, "load": load, "reset": reset, "inc": inc})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000000")}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load, "reset": reset, "inc": inc})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000001")}
				assert.Equal(t, expectedOutputs, outputs)

				load = []bool{true}
				outputs, _ = hs.Tick(map[string][]bool{"in": in, "load": load, "reset": reset, "inc": inc})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000000001")}
				assert.Equal(t, expectedOutputs, outputs)

				outputs, _ = hs.Tock(map[string][]bool{"in": in, "load": load, "reset": reset, "inc": inc})
				expectedOutputs = map[string][]bool{"out": testutils.StringToBoolArray("0000000000100000")}
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
