package font

import (
	"bytes"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Font struct {
	Face *text.GoTextFaceSource
}

func loadFont(data []byte, path string) (*Font, error) {
	face, err := text.NewGoTextFaceSource(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	font := &Font{
		Face: face,
	}

	return font, nil
}
