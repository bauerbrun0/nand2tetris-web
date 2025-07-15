package crypto

import (
	"crypto/rand"
	"encoding/binary"
	"strconv"
)

func GenerateEmailVerificationCode() string {
	var result string
	for range 8 {
		result += strconv.FormatUint(uint64(GenerateRandomUint32(9)), 10)
	}
	return result
}

func GenerateRandomUint32(max uint32) uint32 {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	randUint32 := binary.BigEndian.Uint32(bytes) // Convert bytes to uint32
	return randUint32 % max
}
