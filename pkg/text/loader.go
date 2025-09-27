package text

import (
	"io"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func loadText(reader io.Reader) (any, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	// Parse the font. (Supports TTF and most OTF.)
	fnt, err := opentype.Parse(data)
	if err != nil {
		return nil, err
	}

	// Choose sensible defaults; you can make these configurable later.
	const (
		defaultSize = 16.0
		defaultDPI  = 72.0
	)
	const defaultHinting = font.HintingFull

	face, err := opentype.NewFace(fnt, &opentype.FaceOptions{
		Size:    defaultSize,
		DPI:     defaultDPI,
		Hinting: defaultHinting,
	})
	if err != nil {
		return nil, err
	}

	txt := &Font{
		Font:    fnt,
		Face:    face,
		Size:    defaultSize,
		DPI:     defaultDPI,
		Hinting: defaultHinting,
	}

	return txt, nil
}
