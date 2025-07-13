package audio

import (
	"io/fs"
	"os"
	"testing"
	"testing/fstest"
)

var AudioTest Audio

func TestAudio(t *testing.T) {
	data, err := os.ReadFile("../assets/default.mp3")
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	// create a new in-memory filesystem
	mockFS := fstest.MapFS{
		"assets/default.mp3": &fstest.MapFile{
			Data: data,
			Mode: fs.ModePerm,
		},
	}

	m := NewManager().WithFilesystem(mockFS)

	AudioTest = m.Register("assets/default.mp3")
	m.Load(AudioTest)

	player, ok := m.Player(AudioTest)
	if !ok {
		t.Fatal("audio is nil")
	}

	player.Play()
	for player.IsPlaying() {
		// wait for the audio to finish playing
	}

	m.Unload(AudioTest)
}
