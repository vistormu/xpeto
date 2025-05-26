package image

import (
	"io/fs"
	"os"
	"testing"
	"testing/fstest"
)

var ImageTest Handle

func TestImage(t *testing.T) {
	data, err := os.ReadFile("../assets/default.png")
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	// create a new in-memory filesystem
	mockFS := fstest.MapFS{
		"assets/default.png": &fstest.MapFile{
			Data: data,
			Mode: fs.ModePerm,
		},
	}

	m := NewManager(mockFS)

	ImageTest = m.Register("assets/default.png")
	m.Load(ImageTest)

	image := m.Image(ImageTest)
	if image == nil {
		t.Fatal("image is nil")
	}

	h := image.Bounds().Dy()
	w := image.Bounds().Dx()

	if h != 2048 || w != 2048 {
		t.Fatalf("image size is incorrect: expected %dx%d, got %dx%d", 2048, 2048, w, h)
	}

	m.Unload(ImageTest)
}
