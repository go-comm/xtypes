package codec

import (
	"encoding/hex"
)

func EncodeToHex(dst, src []byte) []byte {
	encLen := hex.EncodedLen(len(src))
	if len(dst) < encLen {
		dst = make([]byte, encLen)
	}
	hex.Encode(dst, src)
	return dst[:encLen]
}

func DecodeFromHex(dst, src []byte) (int, error) {
	n, err := hex.Decode(dst, src)
	return n, err
}
