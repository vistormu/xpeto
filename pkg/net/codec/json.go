package codec

import (
	"encoding/json"
	"io"
)

type JsonCodec struct{}

func newJson() Codec {
	return &JsonCodec{}
}

func (g *JsonCodec) Encode(w io.Writer, v any) error {
	return json.NewEncoder(w).Encode(v)
}

func (g *JsonCodec) Decode(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(v)
}
