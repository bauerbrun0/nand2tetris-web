package resolver

import (
	"strconv"
	"strings"
	"testing"

	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/lexer"
	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/parser"
	"github.com/stretchr/testify/assert"
)

func TestResolve(t *testing.T) {
	tests := []struct {
		name                         string
		chipFileName                 string
		hdls                         map[string]string
		expectedError                string
		expectedResolvedChipDef      ResolvedChipDefinition
		expectedResolvedUsedChipDefs map[string]ResolvedChipDefinition
		printDebug                   bool
	}{
		{
			name:         "And chip with CustomNot part",
			chipFileName: "And",
			hdls: map[string]string{
				"And": `CHIP And {
    IN a, b;
    OUT out;

    PARTS:
    Nand(a=a, b=b, out=nandOut);
    CustomNot(in=nandOut, out=out);
}`,
				"CustomNot": `CHIP CustomNot {
    IN in;
    OUT out;

    PARTS:
    Nand(a=in, b=in, out=out);
}`,
			},
			expectedResolvedChipDef: ResolvedChipDefinition{
				Name: "And",
			},
			expectedResolvedUsedChipDefs: map[string]ResolvedChipDefinition{
				"CustomNot": {
					Name: "CustomNot",
				},
			},
		},
		{
			name:         "File name does not match chip name",
			chipFileName: "And",
			hdls: map[string]string{
				"And": `CHIP And2 {
    IN a, b;
    OUT out;

    PARTS:
    Nand(a=a, b=b, out=nandOut);
    CustomNot(in=nandOut, out=out);
}`,
			},
			expectedError: "Resolution error at line 1, column 6: File name does not match the chip name",
		},
		{
			name:         "Number of inputs exceeds maximum",
			chipFileName: "And",
			hdls: map[string]string{
				"And": `CHIP And {
	IN ` + strings.Join(generateStrings(MAX_NUMBER_OF_IOS+1), ", ") + `;` + `
	OUT out;

	PARTS:
	Nand(a=a1, b=a2, out=nandOut1);
}`,
			},
			expectedError: "Resolution error at line 2, column 399: Number of inputs exceeds maximum allowed",
		},
		{
			name:         "Number of outputs exceeds maximum",
			chipFileName: "And",
			hdls: map[string]string{
				"And": `CHIP And {
	IN in;
	OUT ` + strings.Join(generateStrings(MAX_NUMBER_OF_IOS+1), ", ") + `;` + `

	PARTS:
	Nand(a=a1, b=a2, out=nandOut1);
}`,
			},
			expectedError: "Resolution error at line 3, column 400: Number of outputs exceeds maximum allowed",
		},
		{
			name:         "Input width is 0",
			chipFileName: "And",
			hdls: map[string]string{
				"And": `CHIP And {
    IN a[0], b;
    OUT out;

    PARTS:
    Nand(a=a, b=b, out=nandOut);
		}`,
			},
			expectedError: "Resolution error at line 2, column 8: Input 'a' width out of bounds",
		},
		{
			name:         "Output width is 0",
			chipFileName: "And",
			hdls: map[string]string{
				"And": `CHIP And {
    IN a, b;
    OUT out[0];

    PARTS:
    Nand(a=a, b=b, out=nandOut);
		}`,
			},
			expectedError: "Resolution error at line 3, column 9: Output 'out' width out of bounds",
		},
		{
			name:         "Input width exceeds maximum",
			chipFileName: "And",
			hdls: map[string]string{
				"And": `CHIP And {
    IN a[` + strconv.Itoa(MAX_IO_WIDTH+1) + `], b;
    OUT out;

    PARTS:
    Nand(a=a, b=b, out=nandOut);
		}`,
			},
			expectedError: "Resolution error at line 2, column 8: Input 'a' width out of bounds",
		},
		{
			name:         "Output width exceeds maximum",
			chipFileName: "And",
			hdls: map[string]string{
				"And": `CHIP And {
    IN a, b;
    OUT out[` + strconv.Itoa(MAX_IO_WIDTH+1) + `];

    PARTS:
    Nand(a=a, b=b, out=nandOut);
		}`,
			},
			expectedError: "Resolution error at line 3, column 9: Output 'out' width out of bounds",
		},
		{
			name:         "Duplicate input names",
			chipFileName: "And",
			hdls: map[string]string{
				"And": `CHIP And {
    IN a, a;
    OUT out;

    PARTS:
    Nand(a=a, b=b, out=nandOut);
		}`,
			},
			expectedError: "Resolution error at line 2, column 11: Duplicate input name 'a'",
		},
		{
			name:         "Duplicate output names",
			chipFileName: "And",
			hdls: map[string]string{
				"And": `CHIP And {
    IN a, b;
    OUT out, out;

    PARTS:
    Nand(a=a, b=b, out=nandOut);
		}`,
			},
			expectedError: "Resolution error at line 3, column 14: Duplicate output name 'out'",
		},
		{
			name:         "Number of parts exceeds maximum",
			chipFileName: "And",
			hdls: map[string]string{
				"And": `CHIP And {
    IN a, b;
    OUT out;

    PARTS:
    ` + strings.Repeat("Nand(a=a, b=b, out=nandOut);\n", MAX_NUMBER_OF_PARTS+1) + `
		}`,
			},
			expectedError: "Resolution error at line 106, column 1: Number of parts exceeds maximum allowed",
		},
		{
			name:         "Unknown chip in parts section",
			chipFileName: "And",
			hdls: map[string]string{
				"And": `CHIP And {
    IN a, b;
    OUT out;

    PARTS:
    Nand(a=a, b=b, out=nandOut);
    CustomNot(in=nandOut, out=out);
    Unknown(in=a, out=out);
}`,
				"CustomNot": `CHIP CustomNot {
    IN in;
    OUT out;

    PARTS:
    Nand(a=in, b=in, out=out);
}`,
			},
			expectedError: "Resolution error: Used chip 'Unknown' is neither a built-in chip nor a custom chip",
		},
		{
			name:         "Circular dependency",
			chipFileName: "CustomAnd",
			hdls: map[string]string{
				"CustomAnd": `CHIP CustomAnd {
    IN a, b;
    OUT out;

    PARTS:
    Nand(a=a, b=b, out=nandOut);
    CustomNot(in=nandOut, out=out);
}`,
				"CustomNot": `CHIP CustomNot {
    IN in;
    OUT out;

    PARTS:
    CustomAnd(a=in, b=in, out=out);
}`,
			},
			expectedError: "Resolution error: Circular dependency detected: [CustomAnd CustomNot CustomAnd]",
		},
		{
			name:         "Non-existent pin in a part",
			chipFileName: "And",
			hdls: map[string]string{
				"And": `CHIP And {
    IN a, b;
    OUT out;

    PARTS:
    Nand(a=a, b=b, out=nandOut);
    CustomNot(foobar=nandOut, out=out);
}`,
				"CustomNot": `CHIP CustomNot {
    IN in;
    OUT out;

    PARTS:
    Nand(a=in, b=in, out=out);
}`,
			},
			expectedError: "Resolution error at line 7, column 15: Pin 'foobar' not found in part 'CustomNot'",
		},
		{
			name:         "Output connection's pin width mismatch",
			chipFileName: "And",
			hdls: map[string]string{
				"And": `CHIP And {
    IN a, b;
    OUT out;

    PARTS:
    Nand(a=a, b=b, out=nandOut);
    CustomNot(in=nandOut, out[1..2]=out);
}`,
				"CustomNot": `CHIP CustomNot {
    IN in;
    OUT out;

    PARTS:
    Nand(a=in, b=in, out=out);
}`,
			},
			expectedError: "Resolution error at line 7, column 31: Pin 'out' range out of bounds for part 'CustomNot'",
		},
		{
			name:         "Trying to partially define an internal signal",
			chipFileName: "And",
			hdls: map[string]string{
				"And": `CHIP And {
    IN a, b;
    OUT out;

    PARTS:
    Nand(a=a, b=b, out=nandOut[1..2]);
    CustomNot(in=nandOut, out=out);
}`,
				"CustomNot": `CHIP CustomNot {
    IN in;
    OUT out;

    PARTS:
    Nand(a=in, b=in, out=out);
}`,
			},
			expectedError: "Resolution error at line 6, column 24: Internal output signal 'nandOut' cannot be partially defined",
		},
		{
			name:         "Partially defining an output signal with invalid range",
			chipFileName: "And",
			hdls: map[string]string{
				"And": `CHIP And {
    IN a, b;
    OUT out;

    PARTS:
    Nand(a=a, b=b, out=nandOut);
    CustomNot(in=nandOut, out=out[2..1]);
}`,
				"CustomNot": `CHIP CustomNot {
    IN in;
    OUT out;

    PARTS:
    Nand(a=in, b=in, out=out);
}`,
			},
			expectedError: "Resolution error at line 7, column 35: Signal 'out' range is invalid",
		},
		{
			name:         "Partially defining an output signal with out-of-bounds range",
			chipFileName: "And",
			hdls: map[string]string{
				"And": `CHIP And {
    IN a, b;
    OUT out;

    PARTS:
    Nand(a=a, b=b, out=nandOut);
    CustomNot(in=nandOut, out=out[1..2]);
}`,
				"CustomNot": `CHIP CustomNot {
    IN in;
    OUT out;

    PARTS:
    Nand(a=in, b=in, out=out);
}`,
			},
			expectedError: "Resolution error at line 7, column 35: Signal 'out' range out of bounds",
		},
		{
			name:         "Partially defining an output signal with pin width mismatch",
			chipFileName: "And",
			hdls: map[string]string{
				"And": `CHIP And {
    IN a, b;
    OUT out[3];

    PARTS:
    Nand(a=a, b=b, out=nandOut);
    CustomNot(in=nandOut, out=out[1..2]);
}`,
				"CustomNot": `CHIP CustomNot {
    IN in;
    OUT out;

    PARTS:
    Nand(a=in, b=in, out=out);
}`,
			},
			expectedError: "Resolution error at line 7, column 35: Signal 'out' range width does not match pin 'out' range width",
		},
		{
			name:         "Partially defining an output signal with overlapping ranges",
			chipFileName: "And",
			hdls: map[string]string{
				"And": `CHIP And {
    IN a, b;
    OUT out, outs[3];

    PARTS:
    Nand(a=a, b=b, out=nandOut);
    CustomNot(in=nandOut, out=out, outs=outs[1..2], outs=outs[1..2]);
}`,
				"CustomNot": `CHIP CustomNot {
    IN in;
    OUT out, outs[2];

    PARTS:
    Nand(a=in, b=in, out=out);
}`,
			},
			expectedError: "Resolution error at line 7, column 63: Signal 'outs' range overlaps with existing ranges",
		},
		{
			name:         "Output connection's chip output signal width not matching pin width",
			chipFileName: "And",
			hdls: map[string]string{
				"And": `CHIP And {
    IN a, b, c;
    OUT out, out2[2];

    PARTS:
    Nand(a=a, b=b, out=nandOut);
    CustomNot(in=nandOut, out=out2);
}`,
				"CustomNot": `CHIP CustomNot {
    IN in;
    OUT out, outs[2];

    PARTS:
    Nand(a=in, b=in, out=out);
}`,
			},
			expectedError: "Resolution error at line 7, column 31: Signal 'out2' width does not match pin 'out' width",
		},
		{
			name:         "Output connection's chip output signal range overlaps",
			chipFileName: "And",
			hdls: map[string]string{
				"And": `CHIP And {
    IN a, b, c;
    OUT out;

    PARTS:
    Nand(a=a, b=b, out=nandOut);
    CustomNot(in=nandOut, out=out, out=out);
}`,
				"CustomNot": `CHIP CustomNot {
    IN in;
    OUT out, outs[2];

    PARTS:
    Nand(a=in, b=in, out=out);
}`,
			},
			expectedError: "Resolution error at line 7, column 40: Signal 'out' range overlaps with existing ranges",
		},
		{
			name:         "Duplicate internal signal definition",
			chipFileName: "And",
			hdls: map[string]string{
				"And": `CHIP And {
    IN a, b, c;
    OUT out;

    PARTS:
    Nand(a=a, b=b, out=nandOut);
    CustomNot(in=nandOut, out=out, out=out1, out=out1);
}`,
				"CustomNot": `CHIP CustomNot {
    IN in;
    OUT out, outs[2];

    PARTS:
    Nand(a=in, b=in, out=out);
}`,
			},
			expectedError: "Resolution error at line 7, column 50: Internal signal 'out1' already defined",
		},
		{
			name:         "Input connection's pin width mismatch",
			chipFileName: "And",
			hdls: map[string]string{
				"And": `CHIP And {
    IN a, b, c;
    OUT out;

    PARTS:
    Nand(a[1..2]=a, b=b, out=nandOut);
    CustomNot(in=nandOut, out=out);
}`,
				"CustomNot": `CHIP CustomNot {
    IN in;
    OUT out, outs[2];

    PARTS:
    Nand(a=in, b=in, out=out);
}`,
			},
			expectedError: "Resolution error at line 6, column 12: Pin 'a' range out of bounds for part 'Nand'",
		},
		{
			name:         "Input connection's pin width overlaps",
			chipFileName: "And",
			hdls: map[string]string{
				"And": `CHIP And {
    IN a, b, c;
    OUT out;

    PARTS:
    Nand(a=a, a=a, b=b, out=nandOut);
    CustomNot(in=nandOut, out=out);
}`,
				"CustomNot": `CHIP CustomNot {
    IN in;
    OUT out, outs[2];

    PARTS:
    Nand(a=in, b=in, out=out);
}`,
			},
			expectedError: "Resolution error at line 6, column 15: Pin 'a' range overlaps with existing ranges",
		},
		{
			name:         "Input connection's pin width overlaps with specified range",
			chipFileName: "And",
			hdls: map[string]string{
				"And": `CHIP And {
    IN a, b, c;
    OUT out;

    PARTS:
    Nand(a=a, a[0..0]=a[0..0], b=b, out=nandOut);
    CustomNot(in=nandOut, out=out);
}`,
				"CustomNot": `CHIP CustomNot {
    IN in;
    OUT out, outs[2];

    PARTS:
    Nand(a=in, b=in, out=out);
}`,
			},
			expectedError: "Resolution error at line 6, column 17: Pin 'a' range overlaps with existing ranges",
		},
		{
			name:         "Input connection's signal is unknown",
			chipFileName: "And",
			hdls: map[string]string{
				"And": `CHIP And {
    IN a, b, c;
    OUT out;

    PARTS:
    Nand(a=asd, b=b, out=nandOut);
    CustomNot(in=nandOut, out=out);
}`,
				"CustomNot": `CHIP CustomNot {
    IN in;
    OUT out, outs[2];

    PARTS:
    Nand(a=in, b=in, out=out);
}`,
			},
			expectedError: "Resolution error at line 6, column 12: Signal 'asd' is neither an internal signal nor a chip input",
		},
		{
			name:         "Input connection's (internal) signal width does not match pin width",
			chipFileName: "And",
			hdls: map[string]string{
				"And": `CHIP And {
    IN a, b, c;
    OUT out;

    PARTS:
    Nand(a=a, b=b, out=nandOut);
    CustomNot(in=nandOut, out=out);
}`,
				"CustomNot": `CHIP CustomNot {
    IN in[2];
    OUT out;

    PARTS:
    Nand(a=in[0], b=in[0], out=out);
}`,
			},
			expectedError: "Resolution error at line 7, column 18: Signal 'nandOut' width does not match pin 'in' width",
		},
		{
			name:         "Input connection's (input) signal width does not match pin width",
			chipFileName: "And",
			hdls: map[string]string{
				"And": `CHIP And {
    IN a, b, c;
    OUT out;

    PARTS:
    Nand(a=a, b=b, out=nandOut);
    CustomNot(in=a, out=out);
}`,
				"CustomNot": `CHIP CustomNot {
    IN in[2];
    OUT out;

    PARTS:
    Nand(a=in[0], b=in[0], out=out);
}`,
			},
			expectedError: "Resolution error at line 7, column 18: Signal 'a' width does not match pin 'in' width",
		},
		{
			name:         "Input connection's signal range is invalid",
			chipFileName: "And",
			hdls: map[string]string{
				"And": `CHIP And {
    IN a, b, c;
    OUT out;

    PARTS:
    Nand(a=a, b=b, out=nandOut);
    CustomNot(in=a[2..1], out=out);
}`,
				"CustomNot": `CHIP CustomNot {
    IN in[2];
    OUT out;

    PARTS:
    Nand(a=in[0], b=in[0], out=out);
}`,
			},
			expectedError: "Resolution error at line 7, column 20: Signal 'a' range is invalid",
		},
		{
			name:         "Input connection's internal signal range is out of bounds",
			chipFileName: "And",
			hdls: map[string]string{
				"And": `CHIP And {
    IN a, b, c;
    OUT out;

    PARTS:
    Nand(a=a, b=b, out=nandOut);
    CustomNot(in=nandOut[1..2], out=out);
}`,
				"CustomNot": `CHIP CustomNot {
    IN in[2];
    OUT out;

    PARTS:
    Nand(a=in[0], b=in[0], out=out);
}`,
			},
			expectedError: "Resolution error at line 7, column 26: Signal 'nandOut' range out of bounds",
		},
		{
			name:         "Input connection's input signal range out of bounds",
			chipFileName: "And",
			hdls: map[string]string{
				"And": `CHIP And {
    IN a, b, c[2];
    OUT out;

    PARTS:
    Nand(a=c[1..5], b=b, out=nandOut);
    CustomNot(in=nandOut, out=out);
}`,
				"CustomNot": `CHIP CustomNot {
    IN in;
    OUT out;

    PARTS:
    Nand(a=in[0], b=in[0], out=out);
}`,
			},
			expectedError: "Resolution error at line 6, column 14: Signal 'c' range out of bounds",
		},
		{
			name:         "Input connection's input signal range does not match pin width",
			chipFileName: "And",
			hdls: map[string]string{
				"And": `CHIP And {
    IN a, b, c[2];
    OUT out;

    PARTS:
    Nand(a=c[0..1], b=b, out=nandOut);
    CustomNot(in=nandOut, out=out);
}`,
				"CustomNot": `CHIP CustomNot {
    IN in;
    OUT out;

    PARTS:
    Nand(a=in[0], b=in[0], out=out);
}`,
			},
			expectedError: "Resolution error at line 6, column 14: Signal 'c' range width does not match pin 'a' range width",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// t.Parallel()

			chd := mustLexAndParse(t, tt.hdls[tt.chipFileName])
			r := New(chd, tt.chipFileName, tt.hdls)
			resolvedChipDef, resolvedUsedChipDefs, err := r.Resolve([]string{}, []string{})
			if tt.expectedError != "" {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tt.expectedError)
				}
				assert.Equal(t, tt.expectedError, err.Error(), "error message mismatch")
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Compare resolvedChipDef
			assert.Equal(t, tt.expectedResolvedChipDef.Name, resolvedChipDef.Name, "resolvedChipDef.Name mismatch")

			// Compare resolvedUsedChipDefs
			assert.Equal(t, len(tt.expectedResolvedUsedChipDefs), len(resolvedUsedChipDefs), "resolvedUsedChipDefs length mismatch")
			for name, expectedDef := range tt.expectedResolvedUsedChipDefs {
				actualDef, exists := resolvedUsedChipDefs[name]
				if !exists {
					t.Errorf("expected resolvedUsedChipDefs to contain key %q", name)
					continue
				}
				assert.Equal(t, expectedDef.Name, actualDef.Name, "resolvedUsedChipDefs[%q].Name mismatch", name)
			}
		})
	}
}

func mustLexAndParse(t *testing.T, hdl string) *parser.ParsedChipDefinition {
	t.Helper()

	l := lexer.New(hdl)
	ts, err := l.Tokenize()
	if err != nil {
		t.Fatal(err)
	}

	p := parser.New(ts)
	chd, err := p.ParseChipDefinition()
	if err != nil {
		t.Fatal(err)
	}

	return chd
}

func generateStrings(n int) []string {
	result := make([]string, n)
	for i := 1; i <= n; i++ {
		result[i-1] = "a" + strconv.Itoa(i)
	}
	return result
}
