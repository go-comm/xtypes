package messageid

import (
	"encoding/hex"
)

const (
	// Size length of id
	Size = 8
)

type ID int64

func (id ID) Compare(o ID) int {
	if id < o {
		return -1
	}
	if id == o {
		return 0
	}
	return 1
}

func (id ID) Hex() string {
	var h [16]byte
	var b [8]byte

	b[0] = byte(id >> 56)
	b[1] = byte(id >> 48)
	b[2] = byte(id >> 40)
	b[3] = byte(id >> 32)
	b[4] = byte(id >> 24)
	b[5] = byte(id >> 16)
	b[6] = byte(id >> 8)
	b[7] = byte(id)

	hex.Encode(h[:], b[:])
	return string(h[:])
}
