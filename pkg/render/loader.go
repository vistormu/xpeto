package render

import (
	"bytes"
	_ "image/jpeg"
	_ "image/png"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/vistormu/go-dsa/errors"
)

func NewRenderer(fsys fs.FS, path string) (*Renderer, error) {
	data, err := fs.ReadFile(fsys, path)
	if err != nil {
		return nil, errors.New(ImagePathNotFound).With("path", path).Wrap(err)
	}

	img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(data))
	if err != nil {
		return nil, errors.New(ImageLoadError).With("path", path).Wrap(err)
	}

	return img, nil
}
