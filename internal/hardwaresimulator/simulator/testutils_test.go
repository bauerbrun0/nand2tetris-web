package simulator

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
