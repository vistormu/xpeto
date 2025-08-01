package image

import (
	"bytes"
	_ "image/jpeg"
	_ "image/png"
	"io"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func LoadImage(reader io.Reader, path string) (any, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	image := &Image{
		Img: img,
	}

	return image, nil
}
