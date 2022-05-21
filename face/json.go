package face

import (
	json "github.com/hkloudou/nrpc/face/internal/json"
)

type jsonCodec struct {
}

var defaultCodec Codec = &jsonCodec{}

func (m *jsonCodec) Check(val interface{}) bool {
	return true
}

func (m *jsonCodec) Marshal(val interface{}) ([]byte, error) {
	return json.Marshal(val)
}

func (m *jsonCodec) Unmarshal(data []byte, val interface{}) error {
	return json.Unmarshal(data, val)
}
