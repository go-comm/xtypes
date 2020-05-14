package xtypes

import "encoding/json"

type MarshalerWithBuffer interface {
	MarshalWithBuffer([]byte) ([]byte, error)
}

type Marshaler interface {
	Marshal() ([]byte, error)
}

type Unmarshaler interface {
	Unmarshal([]byte) error
}

func Marshal(v interface{}) ([]byte, error) {
	if m, ok := v.(Marshaler); ok {
		return m.Marshal()
	}
	return json.Marshal(v)
}

func MarshalWithBuffer(data []byte, v interface{}) ([]byte, error) {
	if m, ok := v.(MarshalerWithBuffer); ok {
		return m.MarshalWithBuffer(data)
	}
	return json.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	if m, ok := v.(Unmarshaler); ok {
		return m.Unmarshal(data)
	}
	return json.Unmarshal(data, v)
}
