package xtypes

import "encoding/json"

func Marshal(data []byte, v interface{}) ([]byte, error) {
	if m, ok := v.(Marshaler); ok {
		return m.Marshal(data)
	}
	return json.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	if m, ok := v.(Unmarshaler); ok {
		return m.Unmarshal(data)
	}
	return json.Unmarshal(data, v)
}
