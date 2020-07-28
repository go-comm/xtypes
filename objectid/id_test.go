package objectid

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
	t.Log(id, o, id.Compare(o))
}

func Test_Generator(t *testing.T) {
	for i := 0; i < 10; i++ {
		id, _ := New()
		t.Log(id)
	}
}

func Benchmark_Generator(b *testing.B) {
	var data = make([]byte, Size*2)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		id, _ := New()
		data, _ = xtypes.Marshal(id)
		var o ID
		xtypes.Unmarshal(data, &o)
	}
	b.StopTimer()
}
