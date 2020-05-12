package objectid

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"unsafe"

	"github.com/go-comm/xtypes"
	"github.com/go-comm/xtypes/internal/codec"
)

type ID [Size]byte

func (id *ID) Compare(o xtypes.Object) int {
	if v, ok := o.(*ID); ok {
		return bytes.Compare((*id)[:], (*v)[:])
	}
	return -1
}

func (id *ID) Unmarshal(b []byte) error {
	if len(b) < 2*Size {
		return fmt.Errorf("objectid: expected len %d, not but %d", 2*Size, len(b))
	}
	_, err := codec.DecodeFromHex((*id)[:], b)
	return err
}

func (id *ID) Marshal(b []byte) ([]byte, error) {
	return codec.EncodeToHex(b, (*id)[:]), nil
}

func (id *ID) UnmarshalJSON(b []byte) error {
	if len(b) < 2*Size+2 {
		return fmt.Errorf("objectid: expected len >= %d, not but %d", 2*Size+2, len(b))
	}
	if b[0] != '"' || b[len(b)-1] != '"' {
		return fmt.Errorf("objectid: []byte invalid")
	}
	return id.Unmarshal(b[1 : len(b)-1])
}

func (id *ID) MarshalJSON() ([]byte, error) {
	dst := make([]byte, 2*Size+2)
	dst[0] = '"'
	dst[len(dst)-1] = '"'
	id.Marshal(dst[1:])
	return dst, nil
}

func (id ID) String() string {
	b, _ := id.Marshal(nil)
	return *(*string)(unsafe.Pointer(&b))
}

func (id *ID) Value() (driver.Value, error) {
	b, _ := id.Marshal(nil)
	return *(*string)(unsafe.Pointer(&b)), nil
}

func (id *ID) Scan(v interface{}) error {
	d, ok := v.([]byte)
	if !ok {
		return fmt.Errorf("objectid: scan %+v", v)
	}
	return id.Unmarshal(d)
}
