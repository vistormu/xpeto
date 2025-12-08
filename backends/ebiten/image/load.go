package image

import (
	"bytes"
	_ "image/jpeg"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func load(data []byte, path string) (*image, error) {
	img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	image := &image{img}

	return image, nil
}
