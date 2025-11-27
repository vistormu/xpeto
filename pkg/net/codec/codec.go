package codec

import (
	"io"
)

type Codec interface {
	Encode(w io.Writer, v any) error
	Decode(r io.Reader, v any) error
}
