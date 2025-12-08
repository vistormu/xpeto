package font

import (
	"bytes"

	ebitext "github.com/hajimehoshi/ebiten/v2/text/v2"
)

func load(data []byte, path string) (*font, error) {
	face, err := ebitext.NewGoTextFaceSource(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	font := &font{face}

	return font, nil
}
