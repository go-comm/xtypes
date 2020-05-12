package objectid

import (
	"encoding/json"
	"testing"

	"github.com/go-comm/xtypes"
)

func Test_New(t *testing.T) {
	var err error
	id := New()
	b, _ := xtypes.Marshal(nil, id)

	var o ID
	xtypes.Unmarshal(b, &o)
	t.Log(id, &o, id.Compare(&o))

	b, err = json.Marshal(id)
	if err != nil {
		t.Error(err)
	}
	var o2 ID
	err = json.Unmarshal(b, &o2)
	if err != nil {
		t.Error(err)
	}
	t.Log(id, &o2, id.Compare(&o2))
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
		data, _ = xtypes.Marshal(data, id)
		var o ID
		xtypes.Unmarshal(data, &o)
	}
	b.StopTimer()
}
