package evaluator

import (
	"testing"

	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/graphbuilder"
	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/lexer"
	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/parser"
	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/resolver"
	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/testutils"
	"github.com/stretchr/testify/assert"
)

func TestEvaluate(t *testing.T) {
	tests := []struct {
		name            string
		chipFileName    string
		hdls            map[string]string
		expectedError   string
		afterGraphBuild func(t *testing.T, g *graphbuilder.Graph)
	}{
		{
			name:         "A simple NotGate using a single Built-in Nand gate",
			chipFileName: "NotGate",
			hdls: map[string]string{
				"NotGate": `CHIP NotGate {
					IN in;
					OUT out;

					PARTS:
					Nand(a=in, b=in, out=out);
				}`,
			},
			afterGraphBuild: func(t *testing.T, g *graphbuilder.Graph) {
				e := New(g)
				// first: in = 0 -> out = 1
				e.SetInputs(map[string][]bool{"in": {false}})
				e.Evaluate()
				outputs, _ := e.GetOutputsAndInternalPins()
				expected := map[string][]bool{"out": {true}}
				assert.Equal(t, expected, outputs, "expected out to be true when in is false")

				// second: in = 1 -> out = 0
				e.SetInputs(map[string][]bool{"in": {true}})
				e.Evaluate()
				outputs, _ = e.GetOutputsAndInternalPins()
				expected = map[string][]bool{"out": {false}}
				assert.Equal(t, expected, outputs, "expected out to be false when in is true")

				// third: in = 0 -> out = 1 (again)
				e.SetInputs(map[string][]bool{"in": {false}})
				e.Evaluate()
				outputs, _ = e.GetOutputsAndInternalPins()
				assert.Equal(t, map[string][]bool{"out": {true}}, outputs, "expected out to be true when in is false")
			},
		},
		{
			name:         "A NotNot16Gate using a custom Not16Gate",
			chipFileName: "NotNot16Gate",
			hdls: map[string]string{
				"NotNot16Gate": `CHIP NotNot16Gate {
					IN in[16];
					OUT out[16];
					PARTS:
					Not16Gate(in=in, out=not1out);
					Not16Gate(in=not1out, out=out);
				}`,
				"Not16Gate": `CHIP Not16Gate {
					IN in[16];
					OUT out[16];
					PARTS:
					Nand(a=in[0], b=in[0], out=out[0]);
					Nand(a=in[1], b=in[1], out=out[1]);
					Nand(a=in[2], b=in[2], out=out[2]);
					Nand(a=in[3], b=in[3], out=out[3]);
					Nand(a=in[4], b=in[4], out=out[4]);
					Nand(a=in[5], b=in[5], out=out[5]);
					Nand(a=in[6], b=in[6], out=out[6]);
					Nand(a=in[7], b=in[7], out=out[7]);
					Nand(a=in[8], b=in[8], out=out[8]);
					Nand(a=in[9], b=in[9], out=out[9]);
					Nand(a=in[10], b=in[10], out=out[10]);
					Nand(a=in[11], b=in[11], out=out[11]);
					Nand(a=in[12], b=in[12], out=out[12]);
					Nand(a=in[13], b=in[13], out=out[13]);
					Nand(a=in[14], b=in[14], out=out[14]);
					Nand(a=in[15], b=in[15], out=out[15]);
				}`,
			},
			afterGraphBuild: func(t *testing.T, g *graphbuilder.Graph) {
				e := New(g)
				// first: in[16] = 0 -> out[16] = 0
				e.SetInputs(map[string][]bool{"in": make([]bool, 16)}) // all false
				e.Evaluate()
				outputs, _ := e.GetOutputsAndInternalPins()
				assert.Equal(t, make([]bool, 16), outputs["out"], "expected: in = 0000000000000000 -> out = 0000000000000000")

				// second: in[16] = 1 -> out[16] = 1
				e.SetInputs(map[string][]bool{"in": testutils.RepeatBool(true, 16)})
				e.Evaluate()
				outputs, _ = e.GetOutputsAndInternalPins()
				assert.Equal(t, testutils.RepeatBool(true, 16), outputs["out"], "expected: in = 1111111111111111 -> out = 1111111111111111")

				// third: in[16] = 0 -> out[16] = 0 (again)
				e.SetInputs(map[string][]bool{"in": make([]bool, 16)})
				e.Evaluate()
				outputs, _ = e.GetOutputsAndInternalPins()
				assert.Equal(t, make([]bool, 16), outputs["out"], "expected: in = 0000000000000000 -> out = 0000000000000000")
			},
		},
		{
			name:         "A NotNotGate using a custom NotGate",
			chipFileName: "NotNotGate",
			hdls: map[string]string{
				"NotNotGate": `CHIP NotNotGate {
					IN in;
					OUT out;
					PARTS:
					NotGate(in=in, out=not1out);
					NotGate(in=not1out, out=out);
				}`,
				"NotGate": `CHIP NotGate {
					IN in;
					OUT out;
					PARTS:
					Nand(a=in, b=in, out=out);
				}`,
			},
			afterGraphBuild: func(t *testing.T, g *graphbuilder.Graph) {
				e := New(g)
				// first: in = 0 -> out = 0
				e.SetInputs(map[string][]bool{"in": {false}})
				e.Evaluate()
				outputs, _ := e.GetOutputsAndInternalPins()
				assert.Equal(t, map[string][]bool{"out": {false}}, outputs, "expected: in = false -> out = false")

				// second: in = 1 -> out = 1
				e.SetInputs(map[string][]bool{"in": {true}})
				e.Evaluate()
				outputs, _ = e.GetOutputsAndInternalPins()
				assert.Equal(t, map[string][]bool{"out": {true}}, outputs, "expected: in = true -> out = true")

				// third: in = 0 -> out = 0 (again)
				e.SetInputs(map[string][]bool{"in": {false}})
				e.Evaluate()
				outputs, _ = e.GetOutputsAndInternalPins()
				assert.Equal(t, map[string][]bool{"out": {false}}, outputs, "expected: in = false -> out = false")
			},
		},
		{
			name:         "A DFFGate using a single Built-in DFF gate",
			chipFileName: "DFFGate",
			hdls: map[string]string{
				"DFFGate": `CHIP DFFGate {
					IN in;
					OUT out;

					PARTS:
					DFF(in=in, out=out);
				}`,
			},
			afterGraphBuild: func(t *testing.T, g *graphbuilder.Graph) {
				e := New(g)
				// set the input to true
				e.SetInputs(map[string][]bool{"in": {true}})

				e.Evaluate()
				outputs, _ := e.GetOutputsAndInternalPins()
				assert.Equal(t, map[string][]bool{"out": {false}}, outputs, "after first evaluate, expected out to be false")

				e.Evaluate()
				outputs, _ = e.GetOutputsAndInternalPins()
				assert.Equal(t, map[string][]bool{"out": {false}}, outputs, "after second evaluate, expected out to still be false")

				e.Commit()
				e.Evaluate()
				outputs, _ = e.GetOutputsAndInternalPins()
				assert.Equal(t, map[string][]bool{"out": {true}}, outputs, "after tick-tock, expected out to be true")

				e.Evaluate()
				outputs, _ = e.GetOutputsAndInternalPins()
				assert.Equal(t, map[string][]bool{"out": {true}}, outputs, "after another evaluate, expected out to still be true")

				// set the input to false
				e.SetInputs(map[string][]bool{"in": {false}})

				e.Evaluate()
				outputs, _ = e.GetOutputsAndInternalPins()
				assert.Equal(t, map[string][]bool{"out": {true}}, outputs, "after setting input to false and evaluate, expected out to still be true")

				e.Commit()   // tick
				e.Evaluate() // tock
				outputs, _ = e.GetOutputsAndInternalPins()
				assert.Equal(t, map[string][]bool{"out": {false}}, outputs, "after tick-tock with input false, expected out to be false again")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graph := mustBuildGraph(t, tt.hdls, tt.chipFileName)

			if tt.afterGraphBuild != nil {
				tt.afterGraphBuild(t, graph)
			}
		})
	}
}

func mustBuildGraph(t *testing.T, hdls map[string]string, chipFileName string) *graphbuilder.Graph {
	t.Helper()

	l := lexer.New(hdls[chipFileName])
	ts, err := l.Tokenize()
	if err != nil {
		t.Fatal(err)
	}

	p := parser.New(ts)
	chd, err := p.ParseChipDefinition()
	if err != nil {
		t.Fatal(err)
	}

	r := resolver.New(chd, chipFileName, hdls)
	rchd, rchds, err := r.Resolve([]string{}, []string{})
	if err != nil {
		t.Fatal(err)
	}
	rchds[rchd.Name] = rchd

	gb := graphbuilder.New(rchds)
	graph, err := gb.BuildGraph(chipFileName)
	if err != nil {
		t.Fatal(err)
	}

	return graph
}
