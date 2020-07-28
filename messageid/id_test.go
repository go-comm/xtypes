package messageid

import (
	"testing"

	"github.com/go-comm/xtypes"
)

func Test_New(t *testing.T) {
	var err error
	var b []byte
	id, _ := New()
	b, err = xtypes.Marshal(id)
	if err != nil {
		t.Error(err)
	}
	var o ID
	err = xtypes.Unmarshal(b, &o)
	if err != nil {
		t.Error(err)
	}
	t.Log(id, o, id.Compare(o), id.Hex())
}
