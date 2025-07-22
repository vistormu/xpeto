package audio

import (
	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/ecs"
	"github.com/vistormu/xpeto/internal/event"
)

type action int

const (
	play action = iota
	pause
	resume
	stop
)

type request struct {
	audio  Audio
	action action
	loop   bool
	volume float32
}

type System struct {
	active   *core.HashSet[*Player]
	looping  *core.HashSet[*Player]
	requests *core.QueueArray[request]
}

func NewSystem() *System {
	return &System{
		active:   core.NewHashSet[*Player](),
		looping:  core.NewHashSet[*Player](),
		requests: core.NewQueueArray[request](),
	}
}

func (s *System) OnLoad(ctx *ecs.Context) {
	em, _ := ecs.GetResource[*event.Manager](ctx)

	event.Subscribe(em, func(event AudioPlay) {
		s.requests.Enqueue(request{
			audio:  event.Audio,
			action: play,
			loop:   event.Loop,
			volume: event.Volume,
		})
	})

	event.Subscribe(em, func(event AudioPause) {
		s.requests.Enqueue(request{
			audio:  event.Audio,
			action: pause,
		})
	})

	event.Subscribe(em, func(event AudioResume) {
		s.requests.Enqueue(request{
			audio:  event.Audio,
			action: resume,
		})
	})

	event.Subscribe(em, func(event AudioStop) {
		s.requests.Enqueue(request{
			audio:  event.Audio,
			action: stop,
		})
	})
}

func (s *System) Update(ctx *ecs.Context, dt float32) {
	am, _ := ecs.GetResource[*Manager](ctx)

	for !s.requests.IsEmpty() {
		req, _ := s.requests.Dequeue()

		player, ok := am.Player(req.audio)
		if !ok {
			continue
		}

		switch req.action {
		case play:
			player.SetVolume(float64(req.volume))
			player.Rewind()
			player.Play()

			s.active.Add(player)
			if req.loop {
				s.looping.Add(player)
			}

		case pause:
			player.Pause()

		case resume:
			player.Play()

		case stop:
			player.Close()
			s.active.Remove(player)
			s.looping.Remove(player)
		}

	}

	for _, player := range s.active.Values() {
		if !player.IsPlaying() && s.looping.Contains(player) {
			player.Rewind()
			player.Play()
		}

		if !player.IsPlaying() && !s.looping.Contains(player) {
			player.Close()
			s.active.Remove(player)
		}
	}
}
