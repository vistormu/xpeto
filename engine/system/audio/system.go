package audio

import (
	"github.com/vistormu/xpeto/audio"
	"github.com/vistormu/xpeto/ecs"
	"github.com/vistormu/xpeto/event"
	"github.com/vistormu/xpeto/image"
	st "github.com/vistormu/xpeto/internal/structures"
)

type action int

const (
	play action = iota
	pause
	resume
	stop
)

type request struct {
	handle audio.Handle
	action action
	loop   bool
	volume float32
}

type System struct {
	active   *st.HashSet[*audio.Audio]
	looping  *st.HashSet[*audio.Audio]
	requests *st.QueueArray[request]
}

func NewSystem() *System {
	return &System{
		active:   st.NewHashSet[*audio.Audio](),
		looping:  st.NewHashSet[*audio.Audio](),
		requests: st.NewQueueArray[request](),
	}
}

func (s *System) OnLoad(ctx *ecs.Context) {
	em, _ := ecs.GetResource[*event.Manager](ctx)

	event.Subscribe(em, func(event AudioPlay) {
		s.requests.Enqueue(request{
			handle: event.Audio,
			action: play,
			loop:   event.Loop,
			volume: event.Volume,
		})
	})

	event.Subscribe(em, func(event AudioPause) {
		s.requests.Enqueue(request{
			handle: event.Audio,
			action: pause,
		})
	})

	event.Subscribe(em, func(event AudioResume) {
		s.requests.Enqueue(request{
			handle: event.Audio,
			action: resume,
		})
	})

	event.Subscribe(em, func(event AudioStop) {
		s.requests.Enqueue(request{
			handle: event.Audio,
			action: stop,
		})
	})
}

func (s *System) OnUnload(ctx *ecs.Context) {}

func (s *System) FixedUpdate(ctx *ecs.Context, dt float32) {}

func (s *System) Draw(screen *image.Image) {}

func (s *System) Update(ctx *ecs.Context, dt float32) {
	am, _ := ecs.GetResource[*audio.Manager](ctx)

	for !s.requests.IsEmpty() {
		req, _ := s.requests.Dequeue()

		aud, ok := am.Audio(req.handle)
		if !ok {
			continue
		}

		switch req.action {
		case play:
			aud.SetVolume(float64(req.volume))
			aud.Rewind()
			aud.Play()

			s.active.Add(aud)
			if req.loop {
				s.looping.Add(aud)
			}

		case pause:
			aud.Pause()

		case resume:
			aud.Play()

		case stop:
			aud.Close()
			s.active.Remove(aud)
			s.looping.Remove(aud)
		}

	}

	for _, aud := range s.active.Values() {
		if !aud.IsPlaying() && s.looping.Contains(aud) {
			aud.Rewind()
			aud.Play()
		}

		if !aud.IsPlaying() && !s.looping.Contains(aud) {
			aud.Close()
			s.active.Remove(aud)
		}
	}
}
