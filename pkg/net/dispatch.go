package net

import (
	"github.com/vistormu/xpeto/core/ecs"
)

func dispatch(w *ecs.World) {
	session, _ := ecs.GetResource[session](w)
	if len(session.channels) == 0 {
		return
	}

	cache, _ := ecs.GetResource[cache](w)
	cache.packets = cache.packets[:0]
	cache.packetCodecs = cache.packetCodecs[:0]

	for _, ch := range session.channels {
		packets := ch.Transport.Flush(1000)

		for _, p := range packets {
			cache.packets = append(cache.packets, p)
			cache.packetCodecs = append(cache.packetCodecs, ch.Codec)
		}
	}
}
