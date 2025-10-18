package testutils

var ChipImplementations = map[string]string{
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
		// NotChip(in = dffout, out = outnot);
		Nand(a = dffout, b = dffout, out = outnot);
		DFF(in = in, out = dffout, out = out);
	}`,
	"BitChip": `CHIP BitChip {
		IN in, load;
    	OUT out;

     	PARTS:
     	MuxChip(a = dffout, b = in, sel = load, out = muxout);
      	DFF(in = muxout, out = dffout, out = out);
	}`,
	"RegisterChip": `CHIP RegisterChip {
		IN in[16], load;
		OUT out[16];

		PARTS:
		BitChip(in = in[0], load = load, out = out[0]);
		BitChip(in = in[1], load = load, out = out[1]);
		BitChip(in = in[2], load = load, out = out[2]);
		BitChip(in = in[3], load = load, out = out[3]);
		BitChip(in = in[4], load = load, out = out[4]);
		BitChip(in = in[5], load = load, out = out[5]);
		BitChip(in = in[6], load = load, out = out[6]);
		BitChip(in = in[7], load = load, out = out[7]);
		BitChip(in = in[8], load = load, out = out[8]);
		BitChip(in = in[9], load = load, out = out[9]);
		BitChip(in = in[10], load = load, out = out[10]);
		BitChip(in = in[11], load = load, out = out[11]);
		BitChip(in = in[12], load = load, out = out[12]);
		BitChip(in = in[13], load = load, out = out[13]);
		BitChip(in = in[14], load = load, out = out[14]);
		BitChip(in = in[15], load = load, out = out[15]);
	}`,
	"RAM8Chip": `CHIP RAM8Chip {
		IN in[16], load, address[3];
		OUT out[16];

		PARTS:
		DMux8WayChip(in = load, sel = address, a = a, b = b, c = c, d = d, e = e, f = f, g = g, h = h);
		RegisterChip(in = in, load = a, out = rega);
		RegisterChip(in = in, load = b, out = regb);
		RegisterChip(in = in, load = c, out = regc);
		RegisterChip(in = in, load = d, out = regd);
		RegisterChip(in = in, load = e, out = rege);
		RegisterChip(in = in, load = f, out = regf);
		RegisterChip(in = in, load = g, out = regg);
		RegisterChip(in = in, load = h, out = regh);
		Mux8Way16Chip(a = rega, b = regb, c = regc, d = regd, e = rege, f = regf, g = regg, h = regh, sel = address, out = out);
	}`,
	"RAM64Chip": `CHIP RAM64Chip {
		IN in[16], load, address[6];
    	OUT out[16];

     	PARTS:
     	DMux8WayChip(in = load, sel = address[0..2], a = da, b = db, c = dc, d = dd, e = de, f = df, g = dg, h = dh);
      	RAM8Chip(in = in, load = da, address = address[3..5], out = ra);
      	RAM8Chip(in = in, load = db, address = address[3..5], out = rb);
       	RAM8Chip(in = in, load = dc, address = address[3..5], out = rc);
        RAM8Chip(in = in, load = dd, address = address[3..5], out = rd);
        RAM8Chip(in = in, load = de, address = address[3..5], out = re);
        RAM8Chip(in = in, load = df, address = address[3..5], out = rf);
        RAM8Chip(in = in, load = dg, address = address[3..5], out = rg);
        RAM8Chip(in = in, load = dh, address = address[3..5], out = rh);
        Mux8Way16Chip(a = ra, b = rb, c = rc, d = rd, e = re, f = rf, g = rg, h = rh, sel = address[0..2], out = out);
	}`,
	"RAM512Chip": `CHIP RAM512Chip {
		IN in[16], load, address[9];
		OUT out[16];

		PARTS:
		DMux8WayChip(in = load, sel = address[0..2], a = da, b = db, c = dc, d = dd, e = de, f = df, g = dg, h = dh);
		// I will not use RAM64Chip here because of performance, I will use the built-in RAM64 directly
		RAM64(in = in, load = da, address = address[3..8], out = ra);
		RAM64(in = in, load = db, address = address[3..8], out = rb);
		RAM64(in = in, load = dc, address = address[3..8], out = rc);
		RAM64(in = in, load = dd, address = address[3..8], out = rd);
		RAM64(in = in, load = de, address = address[3..8], out = re);
		RAM64(in = in, load = df, address = address[3..8], out = rf);
		RAM64(in = in, load = dg, address = address[3..8], out = rg);
		RAM64(in = in, load = dh, address = address[3..8], out = rh);
		Mux8Way16Chip(a = ra, b = rb, c = rc, d = rd, e = re, f = rf, g = rg, h = rh, sel = address[0..2], out = out);
	}`,
	"RAM4KChip": `CHIP RAM4KChip {
		IN in[16], load, address[12];
		OUT out[16];

		PARTS:
		DMux8WayChip(in = load, sel = address[0..2], a = da, b = db, c = dc, d = dd, e = de, f = df, g = dg, h = dh);
		RAM512Chip(in = in, load = da, address = address[3..11], out = ra);
		RAM512Chip(in = in, load = db, address = address[3..11], out = rb);
		RAM512Chip(in = in, load = dc, address = address[3..11], out = rc);
		RAM512Chip(in = in, load = dd, address = address[3..11], out = rd);
		RAM512Chip(in = in, load = de, address = address[3..11], out = re);
		RAM512Chip(in = in, load = df, address = address[3..11], out = rf);
		RAM512Chip(in = in, load = dg, address = address[3..11], out = rg);
		RAM512Chip(in = in, load = dh, address = address[3..11], out = rh);
		Mux8Way16Chip(a = ra, b = rb, c = rc, d = rd, e = re, f = rf, g = rg, h = rh, sel = address[0..2], out = out);
	}`,
	"RAM16KChip": `CHIP RAM16KChip {
		IN in[16], load, address[14];
		OUT out[16];

		PARTS:
		DMux4WayChip(in = load, sel = address[0..1], a = da, b = db, c = dc, d = dd);
		RAM4KChip(in = in, load = da, address = address[2..13], out = ra);
		RAM4KChip(in = in, load = db, address = address[2..13], out = rb);
		RAM4KChip(in = in, load = dc, address = address[2..13], out = rc);
		RAM4KChip(in = in, load = dd, address = address[2..13], out = rd);
    	Mux4Way16Chip(a = ra, b = rb, c = rc, d = rd, sel = address[0..1], out = out);
	}`,
	"PCChip": `CHIP PCChip {
		IN in[16], reset, load, inc;
    	OUT out[16];

     	PARTS:
      	Mux16Chip(a = reg, b = reginc, sel = inc, out = muxa);
       	Mux16Chip(a = muxa, b = in, sel = load, out = muxb);
        Mux16Chip(a = muxb, b = false, sel = reset, out = muxc);

        RegisterChip(in = muxc, load = true, out = reg, out = out);
        Inc16Chip(in = reg, out = reginc);
	}`,
	"CPUChip": `CHIP CPUChip {
		IN  inM[16],         // M value input  (M = contents of RAM[A])
        	instruction[16], // Instruction for execution
         reset;           // Signals whether to re-start the current
                         // program (reset==1) or continue executing
                         // the current program (reset==0).

        OUT outM[16],        // M value output
        	writeM,          // Write to M?
         	addressM[15],    // Address in data memory (of M)
          	pc[15];          // address of next instruction

        PARTS:
    	Nand(a=instruction[15], b=instruction[5], out=isCNandDestA);
		// if C instruction and destination is A, then mux out is ALU out
		Mux16Chip(a=outALU, b=instruction, sel=isCNandDestA, out=muxAinOut);

		NotChip(in=instruction[15], out=isA);
		OrChip(a=isA, b=instruction[5], out=regALoad);
		// load if instruction is an A instruction or destination is register A
		RegisterChip(in=muxAinOut, load=regALoad, out=regAOut, out[0..14]=addressM);

		// check if ALU y should be fed from register A or inM
		Mux16Chip(a=regAOut, b=inM, sel=instruction[12], out=y);

		AndChip(a=instruction[15], b=instruction[4], out=regDLoad);
		RegisterChip(in=outALU, load=regDLoad, out=regDOut);

		ALUChip(
			x=regDOut,
			y=y,
			zx=instruction[11],
			nx=instruction[10],
			zy=instruction[9],
			ny=instruction[8],
			f=instruction[7],
			no=instruction[6],
			out=outALU, out=outM, zr=zr, ng=ng
		);
		AndChip(a=instruction[15], b=instruction[3], out=writeM);

		// comparisons
		NotChip(in=zr, out=notZr);
		NotChip(in=ng, out=notNg);

		AndChip(a=notZr, b=notNg, out=GT);
		// zr == EQ
		OrChip(a=zr, b=notNg, out=GE);
		AndChip(a=notZr, b=ng, out=LT);
		// notZr == NE
		OrChip(a=zr, b=ng, out=LE);

		// jump bits
		NotChip(in=instruction[0], out=notJ0);
		NotChip(in=instruction[1], out=notJ1);
		NotChip(in=instruction[2], out=notJ2);

		AndChip(a=instruction[0], b=notJ1, out=J0AndNotJ1);
		AndChip(a=J0AndNotJ1, b=notJ2, out=JGT);

		AndChip(a=notJ0, b=instruction[1], out=notJ0AndJ1);
		AndChip(a=notJ0AndJ1, b=notJ2, out=JEQ);

		AndChip(a=instruction[0], b=instruction[1], out=J0AndJ1);
		AndChip(a=J0AndJ1, b=notJ2, out=JGE);

		AndChip(a=notJ0, b=notJ1, out=notJ0AndNotJ1);
		AndChip(a=notJ0AndNotJ1, b=instruction[2], out=JLT);

		AndChip(a=J0AndNotJ1, b=instruction[2], out=JNE);

		AndChip(a=notJ0AndJ1, b=instruction[2], out=JLE);

		AndChip(a=J0AndJ1, b = instruction[2], out=JMP);

		// jumps
		AndChip(a=JGT, b=GT, out=JJGT);
		AndChip(a=JEQ, b=zr, out=JJEQ);
		AndChip(a=JGE, b=GE, out=JJGE);
		AndChip(a=JLT, b=LT, out=JJLT);
		AndChip(a=JNE, b=notZr, out=JJNE);
		AndChip(a=JLE, b=LE, out=JJLE);

		Or8WayChip(in[0]=JJGT, in[1]=JJEQ, in[2]=JJGE, in[3]=JJLT, in[4]=JJNE, in[5]=JJLE, in[6]=JMP, in[7]=JMP, out=JUMP);
		AndChip(a=JUMP, b=instruction[15], out=pcLoad);

		PCChip(in=regAOut, load=pcLoad, inc=true, reset=reset, out[0..14]=pc);
	}`,
}
