package evaluator

import (
	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/graphbuilder"
)

func getAddressFromBits(bits []*graphbuilder.BitRef) int {
	address := 0
	for i, bit := range bits {
		if bit.Bit.Value {
			address |= (1 << i)
		}
	}
	return address
}
