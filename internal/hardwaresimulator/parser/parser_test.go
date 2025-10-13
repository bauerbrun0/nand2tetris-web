package parser

import (
	"testing"

	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/lexer"
	"github.com/stretchr/testify/assert"
)

func TestParseChipDefinition(t *testing.T) {

	tests := []struct {
		name          string
		input         string
		expectedError string
		expectedChip  ParsedChipDefinition
	}{
		{
			name: "Valid chip definition",
			input: `CHIP HalfAdder {
    IN a[16], b;
    OUT out;

    PARTS:
    nand(a[1..2]=a[2], b=b);
}`,
			expectedChip: ParsedChipDefinition{
				ChipName: ChipName{
					Name: "HalfAdder",
					Loc:  Loc{Line: 1, Column: 6},
				},
				Inputs: []IO{
					{Name: "a", Width: 16, Loc: Loc{Line: 2, Column: 8}},
					{Name: "b", Width: 1, Loc: Loc{Line: 2, Column: 15}},
				},
				Outputs: []IO{
					{Name: "out", Width: 1, Loc: Loc{Line: 3, Column: 9}},
				},
				Parts: []Part{
					{
						Name: "nand",
						Loc:  Loc{Line: 6, Column: 5},
						Connections: []Connection{
							{
								Pin: Pin{
									Name:  "a",
									Range: Range{Start: 1, End: 2, Loc: Loc{Line: 6, Column: 12}, IsSpecified: true},
									Loc:   Loc{Line: 6, Column: 10},
								},
								Signal: Signal{
									Name:  "a",
									Range: Range{Start: 2, End: 2, Loc: Loc{Line: 6, Column: 20}, IsSpecified: true},
									Loc:   Loc{Line: 6, Column: 18},
								},
								Loc: Loc{Line: 6, Column: 10},
							},
							{
								Pin: Pin{
									Name:  "b",
									Range: Range{Start: 0, End: 0, Loc: Loc{Line: 6, Column: 24}},
									Loc:   Loc{Line: 6, Column: 24},
								},
								Signal: Signal{
									Name:  "b",
									Range: Range{Start: 0, End: 0, Loc: Loc{Line: 6, Column: 26}},
									Loc:   Loc{Line: 6, Column: 26},
								},
								Loc: Loc{Line: 6, Column: 24},
							},
						},
					},
				},
			},
		},
		{
			name:          "Missing CHIP keyword",
			input:         `HalfAdder {`,
			expectedError: "Parser error at line 1, column 1: expected CHIP keyword, got [IDENTIFIER] => HalfAdder",
		},
		{
			name:          "Missing chip name",
			input:         `CHIP {`,
			expectedError: "Parser error at line 1, column 6: expected chip name, got [{] => {",
		},
		{
			name:          "Missing '{' after chip name",
			input:         `CHIP HalfAdder IN `,
			expectedError: "Parser error at line 1, column 16: expected '{', got [IN] => IN",
		},
		{
			name:          "Missing IN keyword",
			input:         `CHIP HalfAdder { OUT out; PARTS: nand(a=a, b=b); }`,
			expectedError: "Parser error at line 1, column 27: expected IN, got [PARTS] => PARTS",
		},
		{
			name:          "Missing OUT keyword",
			input:         `CHIP HalfAdder { IN in; PARTS: nand(a=a, b=b); }`,
			expectedError: "Parser error at line 1, column 25: expected OUT, got [PARTS] => PARTS",
		},
		{
			name:          "Missing input list",
			input:         `CHIP HalfAdder { IN ;`,
			expectedError: "Parser error at line 1, column 21: expected identifier, got [;] => ;",
		},
		{
			name:          "Missing output list",
			input:         `CHIP HalfAdder { IN a; OUT ;`,
			expectedError: "Parser error at line 1, column 28: expected identifier, got [;] => ;",
		},
		{
			name:          "Wrong token type as input name",
			input:         `CHIP HalfAdder { IN 123a;`,
			expectedError: "Parser error at line 1, column 21: expected identifier, got [NUMBER] => 123",
		},
		{
			name:          "Missing ']' in input width",
			input:         `CHIP HalfAdder { IN a[16, b`,
			expectedError: "Parser error at line 1, column 25: expected ']', got [,] => ,",
		},
		{
			name:          "Not a number in input width",
			input:         `CHIP HalfAdder { IN a[xx]`,
			expectedError: "Parser error at line 1, column 23: expected number for width, got [IDENTIFIER] => xx",
		},
		{
			name:          "Wrong token type as output name",
			input:         `CHIP HalfAdder { OUT 123a;`,
			expectedError: "Parser error at line 1, column 22: expected identifier, got [NUMBER] => 123",
		},
		{
			name:          "Missing ']' in output width",
			input:         `CHIP HalfAdder { OUT a[16, b`,
			expectedError: "Parser error at line 1, column 26: expected ']', got [,] => ,",
		},
		{
			name:          "Not a number in output width",
			input:         `CHIP HalfAdder { OUT a[xx]`,
			expectedError: "Parser error at line 1, column 24: expected number for width, got [IDENTIFIER] => xx",
		},
		{
			name:          "Missing ';' after input list",
			input:         `CHIP HalfAdder { IN a, b OUT out;`,
			expectedError: "Parser error at line 1, column 26: expected ',' or ';', got [OUT] => OUT",
		},
		{
			name:          "Missing ';' after output list",
			input:         `CHIP HalfAdder { IN a, b; OUT out PARTS`,
			expectedError: "Parser error at line 1, column 35: expected ',' or ';', got [PARTS] => PARTS",
		},

		{
			name:          "Missing PARTS keyword",
			input:         `CHIP HalfAdder { IN a, b; OUT out; nand(a=a, b=b); }`,
			expectedError: "Parser error at line 1, column 36: expected PARTS keyword, got [IDENTIFIER] => nand",
		},
		{
			name:          "Missing ':' after PARTS keyword",
			input:         `CHIP HalfAdder { IN a, b; OUT out; PARTS nand(a=a, b=b); }`,
			expectedError: "Parser error at line 1, column 42: expected ':', got [IDENTIFIER] => nand",
		},
		{
			name:          "Missing part name",
			input:         `CHIP HalfAdder { IN a, b; OUT out; PARTS: (a=a, b=b); }`,
			expectedError: "Parser error at line 1, column 43: expected part name, got [(] => (",
		},
		{
			name:          "Missing '(' after part name",
			input:         `CHIP HalfAdder { IN a, b; OUT out; PARTS: nand a=a, b=b); }`,
			expectedError: "Parser error at line 1, column 48: expected '(', got [IDENTIFIER] => a",
		},
		{
			name:          "Missing connections in part",
			input:         `CHIP HalfAdder { IN a, b; OUT out; PARTS: nand(); }`,
			expectedError: "Parser error at line 1, column 48: expected connection name, got [)] => )",
		},
		{
			name:          "Missing '=' in connection",
			input:         `CHIP HalfAdder { IN a, b; OUT out; PARTS: nand(a, b=b); }`,
			expectedError: "Parser error at line 1, column 49: expected '=', got [,] => ,",
		},
		{
			name:          "Missing pin name in connection",
			input:         `CHIP HalfAdder { IN a, b; OUT out; PARTS: nand(=a, b=b); }`,
			expectedError: "Parser error at line 1, column 48: expected connection name, got [=] => =",
		},
		{
			name:          "Missing signal name in connection",
			input:         `CHIP HalfAdder { IN a, b; OUT out; PARTS: nand(a=, b=b); }`,
			expectedError: "Parser error at line 1, column 50: expected signal name, got [,] => ,",
		},
		{
			name:          "Insufficient number of connections",
			input:         `CHIP HalfAdder { IN a, b; OUT out; PARTS: nand(a=a); }`,
			expectedError: "Parser error at line 1, column 51: expected connection name, got [)] => )",
		},
		{
			name:          "Character ')' after ',' in connections",
			input:         `CHIP HalfAdder { IN a, b; OUT out; PARTS: nand(a=a,); }`,
			expectedError: "Parser error at line 1, column 52: expected connection name, got [)] => )",
		},
		{
			name:          "Missing number in pin range",
			input:         `CHIP HalfAdder { IN a, b; OUT out; PARTS: nand(a[=a, b=b); }`,
			expectedError: "Parser error at line 1, column 50: expected number for range, got [=] => =",
		},
		{
			name:          "Missing ']' in pin range",
			input:         `CHIP HalfAdder { IN a, b; OUT out; PARTS: nand(a[1=a, b=b); }`,
			expectedError: "Parser error at line 1, column 51: expected ']' or '..', got [=] => =",
		},
		{
			name:          "Missing number in pin range end",
			input:         `CHIP HalfAdder { IN a, b; OUT out; PARTS: nand(a[1..=a, b=b); }`,
			expectedError: "Parser error at line 1, column 53: expected number for range end, got [=] => =",
		},
		{
			name:          "Missing number in signal range",
			input:         `CHIP HalfAdder { IN a, b; OUT out; PARTS: nand(a=a[, b=b); }`,
			expectedError: "Parser error at line 1, column 52: expected number for range, got [,] => ,",
		},
		{
			name:          "Missing ']' in signal range",
			input:         `CHIP HalfAdder { IN a, b; OUT out; PARTS: nand(a=a[1, b=b); }`,
			expectedError: "Parser error at line 1, column 53: expected ']' or '..', got [,] => ,",
		},
		{
			name:          "Missing number in signal range end",
			input:         `CHIP HalfAdder { IN a, b; OUT out; PARTS: nand(a=a[1.., b=b); }`,
			expectedError: "Parser error at line 1, column 55: expected number for range end, got [,] => ,",
		},
		{
			name:          "Missing closing ')' in part",
			input:         `CHIP HalfAdder { IN a, b; OUT out; PARTS: nand(a=a, b=b; }`,
			expectedError: "Parser error at line 1, column 56: expected ')', ',' or '[', got [;] => ;",
		},
		{
			name:          "Missing ';' after parts section",
			input:         `CHIP HalfAdder { IN a, b; OUT out; PARTS: nand(a=a, b=b) }`,
			expectedError: "Parser error at line 1, column 58: expected ';', got [}] => }",
		},
		{
			name:          "Missing closing '}' in chip definition",
			input:         `CHIP HalfAdder { IN a, b; OUT out; PARTS: nand(a=a, b=b);`,
			expectedError: "Parser error at line 1, column 58: expected part name or '}', got [EOF] => ",
		},
		{
			name:          "Extra tokens after closing '}'",
			input:         `CHIP HalfAdder { IN a, b; OUT out; PARTS: nand(a=a, b=b); } extra`,
			expectedError: "Parser error at line 1, column 61: expected EOF after '}', got [IDENTIFIER] => extra",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			l := lexer.New(tt.input)
			ts, err := l.Tokenize()

			if err != nil {
				t.Fatalf("Failed to tokenize input: %v", err)
			}

			parser := New(ts)
			chip, err := parser.ParseChipDefinition()
			if tt.expectedError != "" {
				if err == nil {
					t.Fatalf("Expected error '%s', got nil", tt.expectedError)
				}

				assert.Equal(t, tt.expectedError, err.Error(), "Error message mismatch")
				return
			}

			if err != nil {
				t.Fatalf("Failed to parse chip definition: %v", err)
			}

			assert.Equal(t, tt.expectedChip.ChipName.Name, chip.ChipName.Name, "Chip name mismatch")
			locsMustEqual(t, "Chip name", chip.ChipName.Loc, tt.expectedChip.ChipName.Loc)
			assert.Equal(t, len(tt.expectedChip.Inputs), len(chip.Inputs), "Number of inputs mismatch")

			for i, input := range chip.Inputs {
				expectedInput := tt.expectedChip.Inputs[i]

				assert.Equal(t, expectedInput.Name, input.Name, "Input name mismatch")
				assert.Equal(t, expectedInput.Width, input.Width, "Input width mismatch")
				locsMustEqual(t, "Input "+input.Name, input.Loc, expectedInput.Loc)
			}

			assert.Equal(t, len(tt.expectedChip.Outputs), len(chip.Outputs), "Number of outputs mismatch")

			for i, output := range chip.Outputs {
				expectedOutput := tt.expectedChip.Outputs[i]

				assert.Equal(t, expectedOutput.Name, output.Name, "Output name mismatch")
				assert.Equal(t, expectedOutput.Width, output.Width, "Output width mismatch")
				locsMustEqual(t, "Output "+output.Name, output.Loc, tt.expectedChip.Outputs[i].Loc)
			}

			assert.Equal(t, len(tt.expectedChip.Parts), len(chip.Parts), "Number of parts mismatch")

			for i, part := range chip.Parts {
				expectedPart := tt.expectedChip.Parts[i]

				assert.Equal(t, expectedPart.Name, part.Name, "Part name mismatch")
				locsMustEqual(t, "Part "+part.Name, part.Loc, tt.expectedChip.Parts[i].Loc)
				assert.Equal(t, len(expectedPart.Connections), len(part.Connections), "Number of connections mismatch")

				for j, conn := range part.Connections {
					expectedConn := tt.expectedChip.Parts[i].Connections[j]

					// check pin
					assert.Equal(t, expectedConn.Pin.Name, conn.Pin.Name, "Pin name mismatch")
					rangesMustEqual(t, "Pin "+conn.Pin.Name, conn.Pin.Range, expectedConn.Pin.Range)
					locsMustEqual(t, "Pin "+conn.Pin.Name, conn.Pin.Loc, expectedConn.Pin.Loc)

					// check signal
					assert.Equal(t, expectedConn.Signal.Name, conn.Signal.Name, "Signal name mismatch")
					rangesMustEqual(t, "Signal "+conn.Signal.Name, conn.Signal.Range, expectedConn.Signal.Range)
					locsMustEqual(t, "Signal "+conn.Signal.Name, conn.Signal.Loc, expectedConn.Signal.Loc)

					locsMustEqual(t, "Connection "+conn.Pin.Name+"="+conn.Signal.Name, conn.Loc, expectedConn.Loc)
				}
			}
		})
	}
}

func locsMustEqual(t *testing.T, context string, got, expected Loc) {
	assert.Equal(t, expected.Line, got.Line, context+" line mismatch")
	assert.Equal(t, expected.Column, got.Column, context+" column mismatch")
}

func rangesMustEqual(t *testing.T, context string, got, expected Range) {
	assert.Equal(t, expected.Start, got.Start, context+" range start mismatch")
	assert.Equal(t, expected.End, got.End, context+" range end mismatch")
	assert.Equal(t, expected.IsSpecified, got.IsSpecified, context+" range IsSpecified mismatch")
	locsMustEqual(t, context+" range", got.Loc, expected.Loc)
}
