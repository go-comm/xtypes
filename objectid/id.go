package objectid

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"unsafe"
)

const (
	// Size length of id
	Size = 12
)

var (
	nilID  ID
	hex    = []byte("0123456789ABCDEF")
	hexMap [256]byte
)

func init() {
	for k, v := range hex {
		hexMap[v] = byte(k)
	}
}

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
	if len(b) < Size*2 {
		return fmt.Errorf("objectid: expected len %d, not but %d", Size*2, len(b))
	}
	id[0] = hexMap[b[0]]<<4 | hexMap[b[1]]
	id[1] = hexMap[b[2]]<<4 | hexMap[b[3]]
	id[2] = hexMap[b[4]]<<4 | hexMap[b[5]]
	id[3] = hexMap[b[6]]<<4 | hexMap[b[7]]
	id[4] = hexMap[b[8]]<<4 | hexMap[b[9]]
	id[5] = hexMap[b[10]]<<4 | hexMap[b[11]]
	id[6] = hexMap[b[12]]<<4 | hexMap[b[13]]
	id[7] = hexMap[b[14]]<<4 | hexMap[b[15]]
	id[8] = hexMap[b[16]]<<4 | hexMap[b[17]]
	id[9] = hexMap[b[18]]<<4 | hexMap[b[19]]
	id[10] = hexMap[b[20]]<<4 | hexMap[b[21]]
	id[11] = hexMap[b[22]]<<4 | hexMap[b[23]]
	return nil
}

func (id ID) MarshalWithBuffer(b []byte) ([]byte, error) {
	if len(b) < Size*2 {
		return nil, fmt.Errorf("objectid: expected len %d, not but %d", Size*2, len(b))
	}
	b[0] = hex[(id[0]>>4)&0x0F]
	b[1] = hex[id[0]&0x0F]
	b[2] = hex[(id[1]>>4)&0x0F]
	b[3] = hex[id[1]&0x0F]
	b[4] = hex[(id[2]>>4)&0x0F]
	b[5] = hex[id[2]&0x0F]
	b[6] = hex[(id[3]>>4)&0x0F]
	b[7] = hex[id[3]&0x0F]
	b[8] = hex[(id[4]>>4)&0x0F]
	b[9] = hex[id[4]&0x0F]
	b[10] = hex[(id[5]>>4)&0x0F]
	b[11] = hex[id[5]&0x0F]
	b[12] = hex[(id[6]>>4)&0x0F]
	b[13] = hex[id[6]&0x0F]
	b[14] = hex[(id[7]>>4)&0x0F]
	b[15] = hex[id[7]&0x0F]
	b[16] = hex[(id[8]>>4)&0x0F]
	b[17] = hex[id[8]&0x0F]
	b[18] = hex[(id[9]>>4)&0x0F]
	b[19] = hex[id[9]&0x0F]
	b[20] = hex[(id[10]>>4)&0x0F]
	b[21] = hex[id[10]&0x0F]
	b[22] = hex[(id[11]>>4)&0x0F]
	b[23] = hex[id[11]&0x0F]
	return b[:24], nil
}

func (id ID) Marshal() ([]byte, error) {
	var b [24]byte
	id.MarshalWithBuffer(b[:])
	return b[:], nil
}

func (id *ID) UnmarshalJSON(b []byte) error {
	if b[0] != '"' || b[len(b)-1] != '"' {
		return fmt.Errorf("objectid: []byte invalid")
	}
	return id.Unmarshal(b[1 : len(b)-1])
}

func (id ID) MarshalJSON() ([]byte, error) {
	dst := make([]byte, 24+2)
	dst[0] = '"'
	dst[len(dst)-1] = '"'
	id.MarshalWithBuffer(dst[1:])
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
