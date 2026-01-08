package font

import (
	"bytes"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func load(data []byte, path string) (*font, error) {
	face, err := text.NewGoTextFaceSource(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	font := &font{
		GoTextFaceSource: face,
		faces:            make(map[float64]text.Face),
	}

	return font, nil
}
