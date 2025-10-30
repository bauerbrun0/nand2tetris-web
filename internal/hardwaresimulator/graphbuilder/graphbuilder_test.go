package graphbuilder

import (
	"testing"

	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/lexer"
	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/parser"
	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/resolver"
	"github.com/k0kubun/pp/v3"
	"github.com/stretchr/testify/assert"
)

func TestBuildGraph(t *testing.T) {
	tests := []struct {
		name          string
		chipFileName  string
		hdls          map[string]string
		expectedError string
		after         func(t *testing.T, g *Graph)
		printDebug    bool
	}{
		{
			name:         "Expect error in feedback loop with combinational logic",
			chipFileName: "CustomChip",
			hdls: map[string]string{
				"CustomChip": `CHIP CustomChip {
                    IN a, b;
                    OUT out;

                    PARTS:
                    Nand(a=nandout2, b=b, out=nandout1);
                    Nand(a=a, b=nandout1, out=nandout2);
                }`,
			},
			expectedError: "Simulation error: Graph has cycles, cannot determine topological order",
		},
		{
			name:         "Does not expect error in feedback loop with sequential logic",
			chipFileName: "CustomChip",
			hdls: map[string]string{
				"CustomChip": `CHIP CustomChip {
                    IN a, b;
                    OUT out;

                    PARTS:
                    Nand(a=nandout2, b=b, out=nandout1);
                    DFF(in=nandout1, out=nandout2);
                }`,
			},
			after: func(t *testing.T, g *Graph) {
				assert.Equal(t, 2, len(g.Nodes), "expected 2 nodes in the graph")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chd, chds := mustLexParseAndResolve(t, tt.hdls, tt.chipFileName)
			chds[chd.Name] = chd

			gb := New(chds)
			graph, err := gb.BuildGraph(tt.chipFileName)

			if tt.expectedError != "" {
				if err == nil {
					t.Fatalf("expected error: %s, got nil", tt.expectedError)
				}
				assert.Equal(t, tt.expectedError, err.Error())
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.after != nil {
				tt.after(t, graph)
			}

			if tt.printDebug {
				pp.Println(graph.Nodes)
			}
		})
	}
}

func mustLexParseAndResolve(t *testing.T, hdls map[string]string, chipFileName string) (
	*resolver.ResolvedChipDefinition,
	map[string]*resolver.ResolvedChipDefinition,
) {
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
	resolvedChipDef, resolvedUsedChipDefs, err := r.Resolve([]string{}, []string{})
	if err != nil {
		t.Fatal(err)
	}

	return resolvedChipDef, resolvedUsedChipDefs
}
