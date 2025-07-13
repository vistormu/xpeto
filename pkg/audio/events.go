package audio

type AudioPause struct {
	Audio Audio
}

type AudioPlay struct {
	Audio  Audio
	Loop   bool
	Volume float32
}

type AudioResume struct {
	Audio Audio
}

type AudioStop struct {
	Audio Audio
}
