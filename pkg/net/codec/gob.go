package codec

import (
	"encoding/gob"
	"io"
)

type GobCodec struct{}

func NewGob() Codec {
	return &GobCodec{}
}

func (g *GobCodec) Encode(w io.Writer, v any) error {
	return gob.NewEncoder(w).Encode(v)
}

func (g *GobCodec) Decode(r io.Reader, v any) error {
	return gob.NewDecoder(r).Decode(v)
}
