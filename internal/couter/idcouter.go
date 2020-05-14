package couter

import (
	"crypto/rand"
	"encoding/binary"
)

func Couter() uint32 {
	b := make([]byte, 4)
	rand.Read(b)
	return binary.BigEndian.Uint32(b)
}
