package image

import (
	"bytes"
	_ "image/jpeg"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Image struct {
	Img *ebiten.Image
}

func loadImage(data []byte, path string) (*Image, error) {
	img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	image := &Image{
		Img: img,
	}

	return image, nil
}
