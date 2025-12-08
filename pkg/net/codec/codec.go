package codec

import (
	"io"
	"strings"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/log"
)

type Codec interface {
	Encode(w io.Writer, v any) error
	Decode(r io.Reader, v any) error
}

func New(w *ecs.World, c string) Codec {
	switch strings.ToLower(c) {
	case "gob":
		return newGob()
	case "json":
		return newJson()

	default:
		log.LogError(w, "unknown codec", log.F("got", c))
		return nil
	}
}
