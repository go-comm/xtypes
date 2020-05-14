package objectid

import (
	"bytes"
	"database/sql/driver"
	"encoding/base32"
	"fmt"
	"unsafe"
)

const (
	// Size length of id
	Size = 12
)

var (
	nilID ID

	encoding = base32.HexEncoding.WithPadding(base32.NoPadding)
)

type ID [Size]byte

func (id ID) Compare(o ID) int {
	return bytes.Compare(id[:], o[:])
}

func (id ID) IsNil() bool {
	return id == nilID
}

func NilID() ID {
	return nilID
}

func (id *ID) Unmarshal(b []byte) error {
	encLen := encoding.EncodedLen(Size)
	if encLen != len(b) {
		return fmt.Errorf("objectid: expected len %d, not but %d", encLen, len(b))
	}
	_, err := encoding.Decode((*id)[:], b)
	return err
}

func (id ID) MarshalWithBuffer(b []byte) ([]byte, error) {
	encLen := encoding.EncodedLen(Size)
	if len(b) != encLen {
		b = make([]byte, encLen)
	}
	encoding.Encode(b, id[:])
	return b[:encLen], nil
}

func (id ID) Marshal() ([]byte, error) {
	encLen := encoding.EncodedLen(Size)
	b := make([]byte, encLen)
	encoding.Encode(b, id[:])
	return b[:encLen], nil
}

func (id *ID) UnmarshalJSON(b []byte) error {
	encLen := encoding.EncodedLen(Size) + 2
	if encLen != len(b) {
		return fmt.Errorf("objectid: expected len %d, not but %d", encLen, len(b))
	}
	if b[0] != '"' || b[len(b)-1] != '"' {
		return fmt.Errorf("objectid: []byte invalid")
	}
	return id.Unmarshal(b[1 : len(b)-1])
}

func (id ID) MarshalJSON() ([]byte, error) {
	dst := make([]byte, encoding.EncodedLen(Size)+2)
	dst[0] = '"'
	dst[len(dst)-1] = '"'
	id.MarshalWithBuffer(dst[1 : len(dst)-1])
	return dst, nil
}

func (id ID) String() string {
	b, _ := id.Marshal()
	return *(*string)(unsafe.Pointer(&b))
}

func (id ID) Value() (driver.Value, error) {
	if id.IsNil() {
		return nil, nil
	}
	b, err := id.Marshal()
	return *(*string)(unsafe.Pointer(&b)), err
}

func (id *ID) Scan(v interface{}) error {
	switch b := v.(type) {
	case []byte:
		return id.Unmarshal(b)
	case string:
		return id.Unmarshal([]byte(b))
	case nil:
		*id = nilID
		return nil
	default:
		return fmt.Errorf("objectid: scanning unsupported type: %T", b)
	}
}
