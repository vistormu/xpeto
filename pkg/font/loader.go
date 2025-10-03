package font

import (
	"io"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func loadText(reader io.Reader) (any, error) {
	face, err := text.NewGoTextFaceSource(reader)
	if err != nil {
		return nil, err
	}

	font := &Font{
		Face: face,
	}

	return font, nil
}
