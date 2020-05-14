package objectid

import (
	"encoding/json"
	"testing"

	"github.com/go-comm/xtypes"
)

func Test_New(t *testing.T) {
	var err error
	var b []byte
	id := New()
	b, err = xtypes.Marshal(id)
	if err != nil {
		t.Error(err)
	}
	var o ID
	err = xtypes.Unmarshal(b, &o)
	if err != nil {
		t.Error(err)
	}
	t.Log(id, o, id.Compare(o), len(o.String()))

	b, err = json.Marshal(id)
	if err != nil {
		t.Error(err)
	}
	var o2 ID
	err = json.Unmarshal(b, &o2)
	if err != nil {
		t.Error(err)
	}
	t.Log(id, &o2, id.Compare(o2), len(o2.String()))
}

func Test_Generator(t *testing.T) {
	for i := 0; i < 10; i++ {
		id := New()
		t.Log(id)
	}
}

func Benchmark_Generator(b *testing.B) {
	var data = make([]byte, Size*2)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		id := New()
		data, _ = xtypes.Marshal(id)
		var o ID
		xtypes.Unmarshal(data, &o)
	}
	b.StopTimer()
}
