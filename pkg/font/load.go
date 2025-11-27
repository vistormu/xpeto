package font

import (
	"bytes"

	ebitext "github.com/hajimehoshi/ebiten/v2/text/v2"
)

func load(data []byte, path string) (*Font, error) {
	face, err := ebitext.NewGoTextFaceSource(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	font := &Font{
		face: face,
	}

	return font, nil
}
