package hardwaresimulator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var chipImplementations = map[string]string{
	"NotChip": `CHIP NotChip {
		IN in;
		OUT out;

		PARTS:
		Nand(a=in, b=in, out=out);
	}`,
	"AndChip": `CHIP AndChip {
		IN a, b;
		OUT out;

		PARTS:
		Nand(a = a, b = b, out = aNandB);
		NotChip(in = aNandB, out = out);
	}`,
	"OrChip": `CHIP OrChip {
		IN a, b;
		OUT out;

		PARTS:
		NotChip(in = a, out = notA);
		NotChip(in = b, out = notB);
		Nand(a = notA, b = notB, out = out);
	}`,
	"XorChip": `CHIP XorChip {
		IN a, b;
		OUT out;

		PARTS:
		OrChip(a = a, b = b, out = AOrB);
		Nand(a = a, b = b, out = ANandB);
		AndChip(a = AOrB, b = ANandB, out = out);
	}`,
	"MuxChip": `CHIP MuxChip {
		IN a, b, sel;
		OUT out;

		PARTS:
		NotChip(in = sel, out = notSel);
		AndChip(a = notSel, b = a, out = notSelAndA);
		AndChip(a = sel, b = b, out = selAndB);
		OrChip(a = notSelAndA, b = selAndB, out = out);
	}`,
	"DMuxChip": `CHIP DMuxChip {
		IN in, sel;
		OUT a, b;

		PARTS:
		NotChip(in = sel, out = notSel);
		AndChip(a = notSel, b = in, out = a);
		AndChip(a = sel, b = in, out = b);
	}`,
	"Not16Chip": `CHIP Not16Chip {
		IN in[16];
		OUT out[16];

		PARTS:
		NotChip(in = in[0], out = out[0]);
		NotChip(in = in[1], out = out[1]);
		NotChip(in = in[2], out = out[2]);
		NotChip(in = in[3], out = out[3]);
		NotChip(in = in[4], out = out[4]);
		NotChip(in = in[5], out = out[5]);
		NotChip(in = in[6], out = out[6]);
		NotChip(in = in[7], out = out[7]);
		NotChip(in = in[8], out = out[8]);
		NotChip(in = in[9], out = out[9]);
		NotChip(in = in[10], out = out[10]);
		NotChip(in = in[11], out = out[11]);
		NotChip(in = in[12], out = out[12]);
		NotChip(in = in[13], out = out[13]);
		NotChip(in = in[14], out = out[14]);
		NotChip(in = in[15], out = out[15]);
	}`,
	"And16Chip": `CHIP And16Chip {
		IN a[16], b[16];
		OUT out[16];

		PARTS:
		AndChip(a = a[0], b = b[0], out = out[0]);
    	AndChip(a = a[1], b = b[1], out = out[1]);
    	AndChip(a = a[2], b = b[2], out = out[2]);
    	AndChip(a = a[3], b = b[3], out = out[3]);
    	AndChip(a = a[4], b = b[4], out = out[4]);
     	AndChip(a = a[5], b = b[5], out = out[5]);
      	AndChip(a = a[6], b = b[6], out = out[6]);
      	AndChip(a = a[7], b = b[7], out = out[7]);
      	AndChip(a = a[8], b = b[8], out = out[8]);
      	AndChip(a = a[9], b = b[9], out = out[9]);
      	AndChip(a = a[10], b = b[10], out = out[10]);
      	AndChip(a = a[11], b = b[11], out = out[11]);
      	AndChip(a = a[12], b = b[12], out = out[12]);
      	AndChip(a = a[13], b = b[13], out = out[13]);
      	AndChip(a = a[14], b = b[14], out = out[14]);
      	AndChip(a = a[15], b = b[15], out = out[15]);
    }`,
	"Or16Chip": `CHIP Or16Chip {
    	IN a[16], b[16];
     	OUT out[16];

		PARTS:
		OrChip(a = a[0], b = b[0], out = out[0]);
    	OrChip(a = a[1], b = b[1], out = out[1]);
    	OrChip(a = a[2], b = b[2], out = out[2]);
    	OrChip(a = a[3], b = b[3], out = out[3]);
    	OrChip(a = a[4], b = b[4], out = out[4]);
    	OrChip(a = a[5], b = b[5], out = out[5]);
    	OrChip(a = a[6], b = b[6], out = out[6]);
    	OrChip(a = a[7], b = b[7], out = out[7]);
    	OrChip(a = a[8], b = b[8], out = out[8]);
    	OrChip(a = a[9], b = b[9], out = out[9]);
    	OrChip(a = a[10], b = b[10], out = out[10]);
    	OrChip(a = a[11], b = b[11], out = out[11]);
    	OrChip(a = a[12], b = b[12], out = out[12]);
    	OrChip(a = a[13], b = b[13], out = out[13]);
    	OrChip(a = a[14], b = b[14], out = out[14]);
    	OrChip(a = a[15], b = b[15], out = out[15]);
	}`,
	"Mux16Chip": `CHIP Mux16Chip {
		IN a[16], b[16], sel;
		OUT out[16];

		PARTS:
		MuxChip(a = a[0], b = b[0], sel = sel, out = out[0]);
		MuxChip(a = a[1], b = b[1], sel = sel, out = out[1]);
		MuxChip(a = a[2], b = b[2], sel = sel, out = out[2]);
		MuxChip(a = a[3], b = b[3], sel = sel, out = out[3]);
		MuxChip(a = a[4], b = b[4], sel = sel, out = out[4]);
		MuxChip(a = a[5], b = b[5], sel = sel, out = out[5]);
		MuxChip(a = a[6], b = b[6], sel = sel, out = out[6]);
		MuxChip(a = a[7], b = b[7], sel = sel, out = out[7]);
		MuxChip(a = a[8], b = b[8], sel = sel, out = out[8]);
		MuxChip(a = a[9], b = b[9], sel = sel, out = out[9]);
		MuxChip(a = a[10], b = b[10], sel = sel, out = out[10]);
		MuxChip(a = a[11], b = b[11], sel = sel, out = out[11]);
		MuxChip(a = a[12], b = b[12], sel = sel, out = out[12]);
		MuxChip(a = a[13], b = b[13], sel = sel, out = out[13]);
		MuxChip(a = a[14], b = b[14], sel = sel, out = out[14]);
		MuxChip(a = a[15], b = b[15], sel = sel, out = out[15]);
	}`,
	"Or8WayChip": `CHIP Or8WayChip {
		IN in[8];
		OUT out;

		PARTS:
		OrChip(a = in[0], b = in[1], out = or01);
		OrChip(a = or01, b = in[2], out = or012);
		OrChip(a = or012, b = in[3], out = or0123);
		OrChip(a = or0123, b = in[4], out = or01234);
		OrChip(a = or01234, b = in[5], out = or012345);
		OrChip(a = or012345, b = in[6], out = or0123456);
		OrChip(a = or0123456, b = in[7], out = out);
	}`,
	"Mux4Way16Chip": `CHIP Mux4Way16Chip {
		IN a[16], b[16], c[16], d[16], sel[2];
    	OUT out[16];

     	PARTS:
      	Mux16Chip(a = a, b = b, sel = sel[0], out = aMuxb);
		Mux16Chip(a = c, b = d, sel = sel[0], out = cMuxd);
		Mux16Chip(a = aMuxb, b = cMuxd, sel = sel[1], out = out);
	}`,
	"Mux8Way16Chip": `CHIP Mux8Way16Chip {
		IN a[16], b[16], c[16], d[16], e[16], f[16], g[16], h[16], sel[3];
		OUT out[16];

		PARTS:
		Mux4Way16Chip(a = a, b = b, c = c, d = d, sel[0..1] = sel[0..1], out = muxabcd);
		Mux4Way16Chip(a = e, b = f, c = g, d = h, sel[0..1] = sel[0..1], out = muxefgh);
		Mux16Chip(a = muxabcd, b = muxefgh, sel=sel[2], out = out);
	}`,
	"DMux4WayChip": `CHIP DMux4WayChip {
		IN in, sel[2];
		OUT a, b, c, d;

		PARTS:
		DMuxChip(in = in, sel = sel[1], a = dmuxa, b = dmuxb);
		DMuxChip(in = dmuxa, sel = sel[0], a = a, b = b);
		DMuxChip(in = dmuxb, sel = sel[0], a = c, b = d);
	}`,
	"DMux8WayChip": `CHIP DMux8WayChip {
		IN in, sel[3];
		OUT a, b, c, d, e, f, g, h;

		PARTS:
		DMuxChip(in = in, sel = sel[2], a = dmuxa, b = dmuxb);
		DMux4WayChip(in = dmuxa, sel[0..1] = sel[0..1], a = a, b = b, c = c, d = d);
		DMux4WayChip(in = dmuxb, sel[0..1] = sel[0..1], a = e, b = f, c = g, d = h);
	}`,
	"HalfAdderChip": `CHIP HalfAdderChip {
		IN a, b;
		OUT sum, carry;

		PARTS:
		XorChip(a = a, b = b, out = sum);
    	AndChip(a = a, b = b, out = carry);
	}`,
	"FullAdderChip": `CHIP FullAdderChip {
		IN a, b, c;
		OUT sum, carry;

		PARTS:
		HalfAdderChip(a = a, b = b, sum = sumab, carry = carryab);
		HalfAdderChip(a = sumab, b = c, sum = sum, carry = carrySumabc);
		OrChip(a = carryab, b = carrySumabc, out = carry);
	}`,
	"Add16Chip": `CHIP Add16Chip {
		IN a[16], b[16];
		OUT out[16];

		PARTS:
		HalfAdderChip(a = a[0], b = b[0], carry = carry0, sum = out[0]);
		FullAdderChip(a = a[1], b = b[1], c = carry0, carry = carry1, sum = out[1]);
		FullAdderChip(a = a[2], b = b[2], c = carry1, carry = carry2, sum = out[2]);
		FullAdderChip(a = a[3], b = b[3], c = carry2, carry = carry3, sum = out[3]);
		FullAdderChip(a = a[4], b = b[4], c = carry3, carry = carry4, sum = out[4]);
		FullAdderChip(a = a[5], b = b[5], c = carry4, carry = carry5, sum = out[5]);
		FullAdderChip(a = a[6], b = b[6], c = carry5, carry = carry6, sum = out[6]);
		FullAdderChip(a = a[7], b = b[7], c = carry6, carry = carry7, sum = out[7]);
		FullAdderChip(a = a[8], b = b[8], c = carry7, carry = carry8, sum = out[8]);
		FullAdderChip(a = a[9], b = b[9], c = carry8, carry = carry9, sum = out[9]);
		FullAdderChip(a = a[10], b = b[10], c = carry9, carry = carry10, sum = out[10]);
		FullAdderChip(a = a[11], b = b[11], c = carry10, carry = carry11, sum = out[11]);
		FullAdderChip(a = a[12], b = b[12], c = carry11, carry = carry12, sum = out[12]);
		FullAdderChip(a = a[13], b = b[13], c = carry12, carry = carry13, sum = out[13]);
		FullAdderChip(a = a[14], b = b[14], c = carry13, carry = carry14, sum = out[14]);
		FullAdderChip(a = a[15], b = b[15], c = carry14, carry = carry15, sum = out[15]);
	}`,
	"Inc16Chip": `CHIP Inc16Chip {
		IN in[16];
		OUT out[16];

		PARTS:
		HalfAdderChip(a = in[0], b = true, carry = carry1, sum = out[0]);
		HalfAdderChip(a = carry1, b = in[1], carry = carry2, sum  = out[1]);
		HalfAdderChip(a = carry2, b = in[2], carry = carry3, sum  = out[2]);
		HalfAdderChip(a = carry3, b = in[3], carry = carry4, sum  = out[3]);
		HalfAdderChip(a = carry4, b = in[4], carry = carry5, sum  = out[4]);
		HalfAdderChip(a = carry5, b = in[5], carry = carry6, sum  = out[5]);
		HalfAdderChip(a = carry6, b = in[6], carry = carry7, sum  = out[6]);
		HalfAdderChip(a = carry7, b = in[7], carry = carry8, sum  = out[7]);
		HalfAdderChip(a = carry8, b = in[8], carry = carry9, sum  = out[8]);
		HalfAdderChip(a = carry9, b = in[9], carry = carry10, sum = out[9]);
		HalfAdderChip(a = carry10, b = in[10], carry = carry11, sum = out[10]);
		HalfAdderChip(a = carry11, b = in[11], carry = carry12, sum = out[11]);
		HalfAdderChip(a = carry12, b = in[12], carry = carry13, sum = out[12]);
		HalfAdderChip(a = carry13, b = in[13], carry = carry14, sum = out[13]);
		HalfAdderChip(a = carry14, b = in[14], carry = carry15, sum = out[14]);
		HalfAdderChip(a = carry15, b = in[15], carry = carry16, sum = out[15]);
	}`,
	"ALUChip": `CHIP ALUChip {
		IN
			x[16], y[16],  // 16-bit inputs
         	zx, // zero the x input?
         	nx, // negate the x input?
         	zy, // zero the y input?
         	ny, // negate the y input?
         	f,  // compute (out = x + y) or (out = x & y)?
         	no; // negate the out output?
        OUT
        	out[16], // 16-bit output
         	zr,      // if (out == 0) equals 1, else 0
          	ng;      // if (out < 0)  equals 1, else 0

         PARTS:
         // zx / zy
         Mux16Chip(a = x, b[0..15] = false, sel = zx, out = outzx);
         Mux16Chip(a = y, b[0..15] = false, sel = zy, out = outzy);

         // nx / ny
         Not16Chip(in = outzx, out = notOutzx);
         Not16Chip(in = outzy, out = notOutzy);
         Mux16Chip(a = outzx, b = notOutzx, sel = nx, out = outnx);
         Mux16Chip(a = outzy, b = notOutzy, sel = ny, out = outny);

         // f
         And16Chip(a = outnx, b = outny, out = nxAndny);
         Add16Chip(a = outnx, b = outny, out = nxAddny);
         Mux16Chip(a = nxAndny, b = nxAddny, sel = f, out = outf);

         // no
         Not16Chip(in = outf, out = notoutf);
         Mux16Chip(a = outf, b = notoutf, sel = no, out[15] = ng, out[0..7] = out1, out[8..15] = out2, out = out);

         // zr
         Or8WayChip(in = out1, out = or1);
         Or8WayChip(in = out2, out = or2);
         OrChip(a = or1, b = or2, out = or12);
         NotChip(in = or12, out = zr);
	}`,
	"DoubleDFFChip": `CHIP DoubleDFFChip {
		IN in;
		OUT dff1, dff2;

		PARTS:
		DFF(in = in, out = dff1, out = dff1internal);
		DFF(in = dff1internal, out = dff2);
	}`,
	"TestChip": `CHIP TestChip {
		IN in;
		OUT out, outnot;

		PARTS:
		DFF(in = in, out = dffout, out = out);
		NotChip(in = dffout, out = outnot);
		// Nand(a = in, b = in, out = nandout); // TODO: this panics
	}`,
	"BitChip": `CHIP BitChip {
		IN in, load;
    	OUT out;

     	PARTS:
     	MuxChip(a = dffout, b = in, sel = load, out = muxout);
      	DFF(in = muxout, out = dffout, out = out);
	}`,
}

func TestHardwareSimulation(t *testing.T) {

	tests := []struct {
		name                        string
		chipFileName                string
		hdls                        map[string]string
		expectedInputsAfterProcess  map[string]int
		expectedOutputsAfterProcess map[string]int
		afterProcess                func(t *testing.T, hs *HardwareSimulator)
		dontRun                     bool
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
				"NotChip": chipImplementations["NotChip"],
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
				"NotChip": chipImplementations["NotChip"],
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
				"OrChip":  chipImplementations["OrChip"],
				"AndChip": chipImplementations["AndChip"],
				"NotChip": chipImplementations["NotChip"],
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
			hdls:                        chipImplementations,
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
			hdls:                        chipImplementations,
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
			hdls:                        chipImplementations,
			expectedInputsAfterProcess:  map[string]int{"in": 16},
			expectedOutputsAfterProcess: map[string]int{"out": 16},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				input := repeatBool(false, 16)
				outputs, _ := hs.Evaluate(map[string][]bool{"in": input})
				assert.Equal(t, map[string][]bool{"out": repeatBool(true, 16)}, outputs)

				input = repeatBool(true, 16)
				outputs, _ = hs.Evaluate(map[string][]bool{"in": input})
				assert.Equal(t, map[string][]bool{"out": repeatBool(false, 16)}, outputs)

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
			hdls:                        chipImplementations,
			expectedInputsAfterProcess:  map[string]int{"a": 16, "b": 16},
			expectedOutputsAfterProcess: map[string]int{"out": 16},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				a := repeatBool(false, 16)
				b := repeatBool(false, 16)
				outputs, _ := hs.Evaluate(map[string][]bool{"a": a, "b": b})
				assert.Equal(t, map[string][]bool{"out": repeatBool(false, 16)}, outputs)

				a = repeatBool(false, 16)
				b = repeatBool(true, 16)
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b})
				assert.Equal(t, map[string][]bool{"out": repeatBool(false, 16)}, outputs)

				a = repeatBool(true, 16)
				b = repeatBool(false, 16)
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b})
				assert.Equal(t, map[string][]bool{"out": repeatBool(false, 16)}, outputs)

				a = repeatBool(true, 16)
				b = repeatBool(true, 16)
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b})
				assert.Equal(t, map[string][]bool{"out": repeatBool(true, 16)}, outputs)

				a = []bool{true, false, true, false, true, false, true, false, true, false, true, false, true, false, true, false}
				b = []bool{false, true, false, true, false, true, false, true, false, true, false, true, false, true, false, true}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b})
				expectedOutput := repeatBool(false, 16)
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
			hdls:                        chipImplementations,
			expectedInputsAfterProcess:  map[string]int{"a": 16, "b": 16},
			expectedOutputsAfterProcess: map[string]int{"out": 16},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				a := repeatBool(false, 16)
				b := repeatBool(false, 16)
				outputs, _ := hs.Evaluate(map[string][]bool{"a": a, "b": b})
				assert.Equal(t, map[string][]bool{"out": repeatBool(false, 16)}, outputs)

				a = repeatBool(false, 16)
				b = repeatBool(true, 16)
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b})
				assert.Equal(t, map[string][]bool{"out": repeatBool(true, 16)}, outputs)

				a = repeatBool(true, 16)
				b = repeatBool(false, 16)
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b})
				assert.Equal(t, map[string][]bool{"out": repeatBool(true, 16)}, outputs)

				a = repeatBool(true, 16)
				b = repeatBool(true, 16)
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b})
				assert.Equal(t, map[string][]bool{"out": repeatBool(true, 16)}, outputs)

				a = []bool{true, false, true, false, true, false, true, false, true, false, true, false, true, false, true, false}
				b = []bool{false, true, false, true, false, true, false, true, false, true, false, true, false, true, false, true}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b})
				expectedOutput := repeatBool(true, 16)
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
			hdls:                        chipImplementations,
			expectedInputsAfterProcess:  map[string]int{"a": 16, "b": 16, "sel": 1},
			expectedOutputsAfterProcess: map[string]int{"out": 16},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				a := repeatBool(false, 16)
				b := repeatBool(false, 16)
				sel := []bool{false}
				outputs, _ := hs.Evaluate(map[string][]bool{"a": a, "b": b, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": repeatBool(false, 16)}, outputs)

				sel = []bool{true}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": repeatBool(false, 16)}, outputs)

				a = repeatBool(false, 16)
				b = []bool{false, false, false, true, false, false, true, false, false, false, true, true, false, true, false, false}
				sel = []bool{false}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": repeatBool(false, 16)}, outputs)

				sel = []bool{true}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": b}, outputs)

				// 1001100001110110
				a = []bool{true, false, false, true, true, false, false, false, false, true, true, true, false, true, true, false}
				b = repeatBool(false, 16)
				sel = []bool{false}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": a}, outputs)

				sel = []bool{true}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": repeatBool(false, 16)}, outputs)

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
			hdls:                        chipImplementations,
			expectedInputsAfterProcess:  map[string]int{"in": 8},
			expectedOutputsAfterProcess: map[string]int{"out": 1},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				input := repeatBool(false, 8)
				outputs, _ := hs.Evaluate(map[string][]bool{"in": input})
				assert.Equal(t, map[string][]bool{"out": {false}}, outputs)

				input = repeatBool(true, 8)
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
			hdls:                        chipImplementations,
			expectedInputsAfterProcess:  map[string]int{"a": 16, "b": 16, "c": 16, "d": 16, "sel": 2},
			expectedOutputsAfterProcess: map[string]int{"out": 16},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				a := repeatBool(false, 16)
				b := repeatBool(false, 16)
				c := repeatBool(false, 16)
				d := repeatBool(false, 16)
				sel := []bool{false, false}
				outputs, _ := hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": repeatBool(false, 16)}, outputs)

				sel = []bool{false, true}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": repeatBool(false, 16)}, outputs)

				sel = []bool{true, false}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": repeatBool(false, 16)}, outputs)

				sel = []bool{true, true}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": repeatBool(false, 16)}, outputs)

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
			hdls:                        chipImplementations,
			expectedInputsAfterProcess:  map[string]int{"a": 16, "b": 16, "c": 16, "d": 16, "e": 16, "f": 16, "g": 16, "h": 16, "sel": 3},
			expectedOutputsAfterProcess: map[string]int{"out": 16},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				a := repeatBool(false, 16)
				b := repeatBool(false, 16)
				c := repeatBool(false, 16)
				d := repeatBool(false, 16)
				e := repeatBool(false, 16)
				f := repeatBool(false, 16)
				g := repeatBool(false, 16)
				h := repeatBool(false, 16)
				sel := []bool{false, false, false}
				outputs, _ := hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "e": e, "f": f, "g": g, "h": h, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": repeatBool(false, 16)}, outputs)

				sel = []bool{true, false, false}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "e": e, "f": f, "g": g, "h": h, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": repeatBool(false, 16)}, outputs)

				sel = []bool{false, true, false}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "e": e, "f": f, "g": g, "h": h, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": repeatBool(false, 16)}, outputs)

				sel = []bool{true, true, false}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "e": e, "f": f, "g": g, "h": h, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": repeatBool(false, 16)}, outputs)

				sel = []bool{false, false, true}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "e": e, "f": f, "g": g, "h": h, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": repeatBool(false, 16)}, outputs)

				sel = []bool{true, false, true}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "e": e, "f": f, "g": g, "h": h, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": repeatBool(false, 16)}, outputs)

				sel = []bool{false, true, true}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "e": e, "f": f, "g": g, "h": h, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": repeatBool(false, 16)}, outputs)

				sel = []bool{true, true, true}
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "e": e, "f": f, "g": g, "h": h, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": repeatBool(false, 16)}, outputs)

				a = stringToBoolArray("0001001000110100")
				b = stringToBoolArray("0010001101000101")
				c = stringToBoolArray("0011010001010110")
				d = stringToBoolArray("0100010101100111")
				e = stringToBoolArray("0101011001111000")
				f = stringToBoolArray("0110011110001001")
				g = stringToBoolArray("0111100010011010")
				h = stringToBoolArray("1000100110101011")

				sel = stringToBoolArray("000")
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "e": e, "f": f, "g": g, "h": h, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": a}, outputs)

				sel = stringToBoolArray("001")
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "e": e, "f": f, "g": g, "h": h, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": b}, outputs)

				sel = stringToBoolArray("010")
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "e": e, "f": f, "g": g, "h": h, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": c}, outputs)

				sel = stringToBoolArray("011")
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "e": e, "f": f, "g": g, "h": h, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": d}, outputs)

				sel = stringToBoolArray("100")
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "e": e, "f": f, "g": g, "h": h, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": e}, outputs)

				sel = stringToBoolArray("101")
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "e": e, "f": f, "g": g, "h": h, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": f}, outputs)

				sel = stringToBoolArray("110")
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "e": e, "f": f, "g": g, "h": h, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": g}, outputs)

				sel = stringToBoolArray("111")
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b, "c": c, "d": d, "e": e, "f": f, "g": g, "h": h, "sel": sel})
				assert.Equal(t, map[string][]bool{"out": h}, outputs)
			},
		},
		{
			name:                        "DMux4Way Chip",
			chipFileName:                "DMux4WayChip",
			hdls:                        chipImplementations,
			expectedInputsAfterProcess:  map[string]int{"in": 1, "sel": 2},
			expectedOutputsAfterProcess: map[string]int{"a": 1, "b": 1, "c": 1, "d": 1},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				in := []bool{false}
				sel := stringToBoolArray("00")
				outputs, _ := hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				expectedOutputs := map[string][]bool{"a": {false}, "b": {false}, "c": {false}, "d": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				sel = stringToBoolArray("01")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				assert.Equal(t, expectedOutputs, outputs)

				sel = stringToBoolArray("10")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				assert.Equal(t, expectedOutputs, outputs)

				sel = stringToBoolArray("11")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				assert.Equal(t, expectedOutputs, outputs)

				in = []bool{true}

				sel = stringToBoolArray("00")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				expectedOutputs = map[string][]bool{"a": {true}, "b": {false}, "c": {false}, "d": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				sel = stringToBoolArray("01")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				expectedOutputs = map[string][]bool{"a": {false}, "b": {true}, "c": {false}, "d": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				sel = stringToBoolArray("10")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				expectedOutputs = map[string][]bool{"a": {false}, "b": {false}, "c": {true}, "d": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				sel = stringToBoolArray("11")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				expectedOutputs = map[string][]bool{"a": {false}, "b": {false}, "c": {false}, "d": {true}}
				assert.Equal(t, expectedOutputs, outputs)
			},
		},
		{
			name:                        "DMux8Way Chip",
			chipFileName:                "DMux8WayChip",
			hdls:                        chipImplementations,
			expectedInputsAfterProcess:  map[string]int{"in": 1, "sel": 3},
			expectedOutputsAfterProcess: map[string]int{"a": 1, "b": 1, "c": 1, "d": 1, "e": 1, "f": 1, "g": 1, "h": 1},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				in := []bool{false}
				sel := stringToBoolArray("000")
				outputs, _ := hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				expectedOutputs := map[string][]bool{"a": {false}, "b": {false}, "c": {false}, "d": {false}, "e": {false}, "f": {false}, "g": {false}, "h": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				sel = stringToBoolArray("001")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				assert.Equal(t, expectedOutputs, outputs)

				sel = stringToBoolArray("010")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				assert.Equal(t, expectedOutputs, outputs)

				sel = stringToBoolArray("011")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				assert.Equal(t, expectedOutputs, outputs)

				sel = stringToBoolArray("100")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				assert.Equal(t, expectedOutputs, outputs)

				sel = stringToBoolArray("101")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				assert.Equal(t, expectedOutputs, outputs)

				sel = stringToBoolArray("110")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				assert.Equal(t, expectedOutputs, outputs)

				sel = stringToBoolArray("111")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				assert.Equal(t, expectedOutputs, outputs)

				in = []bool{true}
				sel = stringToBoolArray("000")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				expectedOutputs = map[string][]bool{"a": {true}, "b": {false}, "c": {false}, "d": {false}, "e": {false}, "f": {false}, "g": {false}, "h": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				sel = stringToBoolArray("001")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				expectedOutputs = map[string][]bool{"a": {false}, "b": {true}, "c": {false}, "d": {false}, "e": {false}, "f": {false}, "g": {false}, "h": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				sel = stringToBoolArray("010")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				expectedOutputs = map[string][]bool{"a": {false}, "b": {false}, "c": {true}, "d": {false}, "e": {false}, "f": {false}, "g": {false}, "h": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				sel = stringToBoolArray("011")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				expectedOutputs = map[string][]bool{"a": {false}, "b": {false}, "c": {false}, "d": {true}, "e": {false}, "f": {false}, "g": {false}, "h": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				sel = stringToBoolArray("100")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				expectedOutputs = map[string][]bool{"a": {false}, "b": {false}, "c": {false}, "d": {false}, "e": {true}, "f": {false}, "g": {false}, "h": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				sel = stringToBoolArray("101")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				expectedOutputs = map[string][]bool{"a": {false}, "b": {false}, "c": {false}, "d": {false}, "e": {false}, "f": {true}, "g": {false}, "h": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				sel = stringToBoolArray("110")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				expectedOutputs = map[string][]bool{"a": {false}, "b": {false}, "c": {false}, "d": {false}, "e": {false}, "f": {false}, "g": {true}, "h": {false}}
				assert.Equal(t, expectedOutputs, outputs)

				sel = stringToBoolArray("111")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in, "sel": sel})
				expectedOutputs = map[string][]bool{"a": {false}, "b": {false}, "c": {false}, "d": {false}, "e": {false}, "f": {false}, "g": {false}, "h": {true}}
				assert.Equal(t, expectedOutputs, outputs)
			},
		},
		{
			name:                        "HalfAdder Chip",
			chipFileName:                "HalfAdderChip",
			hdls:                        chipImplementations,
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
			hdls:                        chipImplementations,
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
			hdls:                        chipImplementations,
			expectedInputsAfterProcess:  map[string]int{"a": 16, "b": 16},
			expectedOutputsAfterProcess: map[string]int{"out": 16},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				a := repeatBool(false, 16)
				b := repeatBool(false, 16)
				outputs, _ := hs.Evaluate(map[string][]bool{"a": a, "b": b})
				assert.Equal(t, map[string][]bool{"out": repeatBool(false, 16)}, outputs)

				a = repeatBool(false, 16)
				b = repeatBool(true, 16)
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b})
				assert.Equal(t, map[string][]bool{"out": repeatBool(true, 16)}, outputs)

				a = repeatBool(true, 16)
				b = repeatBool(true, 16)
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b})
				expectedOutput := stringToBoolArray("1111111111111110")
				assert.Equal(t, map[string][]bool{"out": expectedOutput}, outputs)

				a = stringToBoolArray("1010101010101010")
				b = stringToBoolArray("0101010101010101")
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b})
				expectedOutput = stringToBoolArray("1111111111111111")
				assert.Equal(t, map[string][]bool{"out": expectedOutput}, outputs)

				a = stringToBoolArray("0011110011000011")
				b = stringToBoolArray("0000111111110000")
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b})
				expectedOutput = stringToBoolArray("0100110010110011")
				assert.Equal(t, map[string][]bool{"out": expectedOutput}, outputs)

				a = stringToBoolArray("0001001000110100")
				b = stringToBoolArray("1001100001110110")
				outputs, _ = hs.Evaluate(map[string][]bool{"a": a, "b": b})
				expectedOutput = stringToBoolArray("1010101010101010")
				assert.Equal(t, map[string][]bool{"out": expectedOutput}, outputs)
			},
		},
		{
			name:                        "Inc16 Chip",
			chipFileName:                "Inc16Chip",
			hdls:                        chipImplementations,
			expectedInputsAfterProcess:  map[string]int{"in": 16},
			expectedOutputsAfterProcess: map[string]int{"out": 16},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				in := repeatBool(false, 16)
				outputs, _ := hs.Evaluate(map[string][]bool{"in": in})
				expectedOutput := stringToBoolArray("0000000000000001")
				assert.Equal(t, map[string][]bool{"out": expectedOutput}, outputs)

				in = stringToBoolArray("1111111111111111")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in})
				expectedOutput = stringToBoolArray("0000000000000000")
				assert.Equal(t, map[string][]bool{"out": expectedOutput}, outputs)

				in = stringToBoolArray("0000000000000101")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in})
				expectedOutput = stringToBoolArray("0000000000000110")
				assert.Equal(t, map[string][]bool{"out": expectedOutput}, outputs)

				in = stringToBoolArray("1111111111111011")
				outputs, _ = hs.Evaluate(map[string][]bool{"in": in})
				expectedOutput = stringToBoolArray("1111111111111100")
				assert.Equal(t, map[string][]bool{"out": expectedOutput}, outputs)
			},
		},
		{
			name:                        "ALU Chip",
			chipFileName:                "ALUChip",
			hdls:                        chipImplementations,
			expectedInputsAfterProcess:  map[string]int{"x": 16, "y": 16, "zx": 1, "nx": 1, "zy": 1, "ny": 1, "f": 1, "no": 1},
			expectedOutputsAfterProcess: map[string]int{"out": 16, "zr": 1, "ng": 1},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				x := stringToBoolArray("0000000000000000")
				y := stringToBoolArray("1111111111111111")

				zx, nx, zy, ny, f, no := getALUFlagInputs("101010")
				outputs, _ := hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput := stringToBoolArray("0000000000000000")
				expectedZR := []bool{true}
				expectedNG := []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("111111")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = stringToBoolArray("0000000000000001")
				expectedZR = []bool{false}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("111010")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = stringToBoolArray("1111111111111111")
				expectedZR = []bool{false}
				expectedNG = []bool{true}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("001100")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = stringToBoolArray("0000000000000000")
				expectedZR = []bool{true}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("110000")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = stringToBoolArray("1111111111111111")
				expectedZR = []bool{false}
				expectedNG = []bool{true}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("001101")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = stringToBoolArray("1111111111111111")
				expectedZR = []bool{false}
				expectedNG = []bool{true}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("110001")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = stringToBoolArray("0000000000000000")
				expectedZR = []bool{true}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("001111")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = stringToBoolArray("0000000000000000")
				expectedZR = []bool{true}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("110011")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = stringToBoolArray("0000000000000001")
				expectedZR = []bool{false}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("011111")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = stringToBoolArray("0000000000000001")
				expectedZR = []bool{false}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("110111")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = stringToBoolArray("0000000000000000")
				expectedZR = []bool{true}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("001110")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = stringToBoolArray("1111111111111111")
				expectedZR = []bool{false}
				expectedNG = []bool{true}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("110010")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = stringToBoolArray("1111111111111110")
				expectedZR = []bool{false}
				expectedNG = []bool{true}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("000010")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = stringToBoolArray("1111111111111111")
				expectedZR = []bool{false}
				expectedNG = []bool{true}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("010011")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = stringToBoolArray("0000000000000001")
				expectedZR = []bool{false}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("000111")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = stringToBoolArray("1111111111111111")
				expectedZR = []bool{false}
				expectedNG = []bool{true}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("000000")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = stringToBoolArray("0000000000000000")
				expectedZR = []bool{true}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("010101")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = stringToBoolArray("1111111111111111")
				expectedZR = []bool{false}
				expectedNG = []bool{true}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				x = stringToBoolArray("0000000000010001")
				y = stringToBoolArray("0000000000000011")

				zx, nx, zy, ny, f, no = getALUFlagInputs("101010")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = stringToBoolArray("0000000000000000")
				expectedZR = []bool{true}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("111111")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = stringToBoolArray("0000000000000001")
				expectedZR = []bool{false}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("111010")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = stringToBoolArray("1111111111111111")
				expectedZR = []bool{false}
				expectedNG = []bool{true}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("001100")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = stringToBoolArray("0000000000010001")
				expectedZR = []bool{false}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("110000")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = stringToBoolArray("0000000000000011")
				expectedZR = []bool{false}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("001101")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = stringToBoolArray("1111111111101110")
				expectedZR = []bool{false}
				expectedNG = []bool{true}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("110001")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = stringToBoolArray("1111111111111100")
				expectedZR = []bool{false}
				expectedNG = []bool{true}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("001111")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = stringToBoolArray("1111111111101111")
				expectedZR = []bool{false}
				expectedNG = []bool{true}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("110011")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = stringToBoolArray("1111111111111101")
				expectedZR = []bool{false}
				expectedNG = []bool{true}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("011111")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = stringToBoolArray("0000000000010010")
				expectedZR = []bool{false}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("110111")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = stringToBoolArray("0000000000000100")
				expectedZR = []bool{false}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("001110")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = stringToBoolArray("0000000000010000")
				expectedZR = []bool{false}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("110010")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = stringToBoolArray("0000000000000010")
				expectedZR = []bool{false}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("000010")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = stringToBoolArray("0000000000010100")
				expectedZR = []bool{false}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("010011")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = stringToBoolArray("0000000000001110")
				expectedZR = []bool{false}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("000111")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = stringToBoolArray("1111111111110010")
				expectedZR = []bool{false}
				expectedNG = []bool{true}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("000000")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = stringToBoolArray("0000000000000001")
				expectedZR = []bool{false}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)

				zx, nx, zy, ny, f, no = getALUFlagInputs("010101")
				outputs, _ = hs.Evaluate(map[string][]bool{"x": x, "y": y, "zx": zx, "nx": nx, "zy": zy, "ny": ny, "f": f, "no": no})
				expectedOutput = stringToBoolArray("0000000000010011")
				expectedZR = []bool{false}
				expectedNG = []bool{false}
				assert.Equal(t, map[string][]bool{"out": expectedOutput, "zr": expectedZR, "ng": expectedNG}, outputs)
			},
		},
		{
			name:                        "Double DFF Chip",
			chipFileName:                "DoubleDFFChip",
			hdls:                        chipImplementations,
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
			hdls:                        chipImplementations,
			expectedInputsAfterProcess:  map[string]int{"in": 1},
			expectedOutputsAfterProcess: map[string]int{"out": 1, "outnot": 1},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				in := []bool{true}
				outputs, internalPins := hs.Tick(map[string][]bool{"in": in})
				expectedOutputs := map[string][]bool{"out": {false}, "outnot": {true}}
				assert.Equal(t, expectedOutputs, outputs)
				expectedInternalPins := map[string][]bool{"dffout": {false}}
				assert.Equal(t, expectedInternalPins, internalPins)

				outputs, internalPins = hs.Tock(map[string][]bool{"in": in})
				expectedOutputs = map[string][]bool{"out": {true}, "outnot": {false}}
				assert.Equal(t, expectedOutputs, outputs)
				expectedInternalPins = map[string][]bool{"dffout": {true}}
				assert.Equal(t, expectedInternalPins, internalPins)
			},
		},
		{
			dontRun:                     true, // TODO: only temporarily disable this test
			name:                        "Bit Chip",
			chipFileName:                "BitChip",
			hdls:                        chipImplementations,
			expectedInputsAfterProcess:  map[string]int{"in": 1, "load": 1},
			expectedOutputsAfterProcess: map[string]int{"out": 1},
			afterProcess: func(t *testing.T, hs *HardwareSimulator) {
				/*
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
				*/

				in := []bool{true}
				load := []bool{true}
				outputs, internalPins := hs.Tick(map[string][]bool{"in": in, "load": load})
				_ = internalPins
				expectedOutputs := map[string][]bool{"out": {false}}
				assert.Equal(t, expectedOutputs, outputs)
				// pp.Println("tick -> in: 1; load: 1")
				// pp.Println("outputs")
				// pp.Println(outputs)
				// pp.Println("internalPins")
				// pp.Println(internalPins)

				outputs, internalPins = hs.Tock(map[string][]bool{"in": in, "load": load})
				expectedOutputs = map[string][]bool{"out": {true}}
				assert.Equal(t, expectedOutputs, outputs)
				// pp.Println("tock -> in: 1; load: 1")
				// pp.Println("outputs")
				// pp.Println(outputs)
				// pp.Println("internalPins")
				// pp.Println(internalPins)

				in = []bool{false}
				load = []bool{false}
				outputs, internalPins = hs.Tick(map[string][]bool{"in": in, "load": load})
				expectedOutputs = map[string][]bool{"out": {true}}
				assert.Equal(t, expectedOutputs, outputs)
				// pp.Println("tick -> in: 0; load: 0")
				// pp.Println("outputs")
				// pp.Println(outputs)
				// pp.Println("internalPins")
				// pp.Println(internalPins)

				outputs, internalPins = hs.Tock(map[string][]bool{"in": in, "load": load})
				expectedOutputs = map[string][]bool{"out": {true}}
				// pp.Println("tock -> in: 0; load: 0")
				// pp.Println("outputs")
				// pp.Println(outputs)
				// pp.Println("internalPins")
				// pp.Println(internalPins)
				assert.Equal(t, expectedOutputs, outputs)
			},
		},
	}

	for _, tt := range tests {
		if tt.dontRun {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
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

func repeatBool(value bool, count int) []bool {
	result := make([]bool, count)
	for i := range count {
		result[i] = value
	}
	return result
}

func stringToBoolArray(s string) []bool {
	result := make([]bool, len(s))
	for i := len(s) - 1; i >= 0; i-- {
		char := s[i]
		if char == '1' {
			result[len(s)-1-i] = true
		} else {
			result[len(s)-1-i] = false
		}
	}
	return result
}

func getALUFlagInputs(s string) (zx, nx, zy, ny, f, no []bool) {
	if len(s) != 6 {
		panic("invalid ALU flag string length")
	}

	zx = []bool{s[0] == '1'}
	nx = []bool{s[1] == '1'}
	zy = []bool{s[2] == '1'}
	ny = []bool{s[3] == '1'}
	f = []bool{s[4] == '1'}
	no = []bool{s[5] == '1'}
	return
}
