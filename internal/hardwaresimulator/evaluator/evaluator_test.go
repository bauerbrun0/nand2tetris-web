package evaluator

import (
	"testing"

	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/graphbuilder"
	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/lexer"
	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/parser"
	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/resolver"
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
				g.InputPins["in"].Bits[0].Value = false
				e.Evaluate(false)
				out := g.OutputPins["out"].Bits[0].Value
				assert.True(t, out, "expected out to be true")

				// second: in = 1 -> out = 0
				g.InputPins["in"].Bits[0].Value = true
				e.Evaluate(false)
				out = g.OutputPins["out"].Bits[0].Value
				assert.False(t, out, "expected out to be false")

				// third: in = 0 -> out = 1 (again)
				g.InputPins["in"].Bits[0].Value = false
				e.Evaluate(false)
				out = g.OutputPins["out"].Bits[0].Value
				assert.True(t, out, "expected out to be true")
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
				for _, bit := range g.InputPins["in"].Bits {
					bit.Value = false
				}
				e.Evaluate(false)
				out := make([]bool, 16)
				for i, bit := range g.OutputPins["out"].Bits {
					out[i] = bit.Value
				}
				expected := make([]bool, 16) // all false
				assert.Equal(t, expected, out, "expected: in = 0000000000000000 -> out = 0000000000000000")

				// second: in[16] = 1 -> out[16] = 1
				for _, bit := range g.InputPins["in"].Bits {
					bit.Value = true
				}
				e.Evaluate(false)
				for i, bit := range g.OutputPins["out"].Bits {
					out[i] = bit.Value
				}
				expected = make([]bool, 16)
				for i := range expected {
					expected[i] = true
				}
				assert.Equal(t, expected, out, "expected: in = 1111111111111111 -> out = 1111111111111111")

				// third: in[16] = 0 -> out[16] = 0 (again)
				for _, bit := range g.InputPins["in"].Bits {
					bit.Value = false
				}
				e.Evaluate(false)
				for i, bit := range g.OutputPins["out"].Bits {
					out[i] = bit.Value
				}
				expected = make([]bool, 16) // all false
				assert.Equal(t, expected, out, "expected: in = 0000000000000000 -> out = 0000000000000000")
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
				g.InputPins["in"].Bits[0].Value = false
				e.Evaluate(false)
				out := g.OutputPins["out"].Bits[0].Value
				assert.False(t, out, "expected: in = false -> out = false")

				// second: in = 1 -> out = 1
				g.InputPins["in"].Bits[0].Value = true
				e.Evaluate(false)
				out = g.OutputPins["out"].Bits[0].Value
				assert.True(t, out, "expected in = true -> out = true")

				// third: in = 0 -> out = 0 (again)
				g.InputPins["in"].Bits[0].Value = false
				e.Evaluate(false)
				out = g.OutputPins["out"].Bits[0].Value
				assert.False(t, out, "expected: in = false -> out = false")
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

				e.Evaluate(false)
				out := g.OutputPins["out"].Bits[0].Value
				assert.False(t, out, "after first evaluate, expected out to be false")

				e.Evaluate(false)
				out = g.OutputPins["out"].Bits[0].Value
				assert.False(t, out, "after second evaluate, expected out to still be false")

				e.Evaluate(true)  // tick
				e.Evaluate(false) // tock
				out = g.OutputPins["out"].Bits[0].Value
				assert.True(t, out, "after step and evaluate, expected out to be true")

				e.Evaluate(false)
				out = g.OutputPins["out"].Bits[0].Value
				assert.True(t, out, "after another evaluate, expected out to still be true")

				// set the input to false
				e.SetInputs(map[string][]bool{"in": {false}})

				e.Evaluate(false)
				out = g.OutputPins["out"].Bits[0].Value
				assert.True(t, out, "after setting input to false and evaluate, expected out to still be true")

				e.Evaluate(true)  // tick
				e.Evaluate(false) // tock
				out = g.OutputPins["out"].Bits[0].Value
				assert.False(t, out, "after evaluate, expected out to be false")
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
