package net

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/pkg/log"
)

func logTransportErrors(w *ecs.World) {
	session, ok := ecs.GetResource[session](w)
	if !ok || len(session.channels) == 0 {
		return
	}

	for _, ch := range session.channels {
		errs := ch.Transport.FlushErrors()

		for _, err := range errs {
			log.LogError(w, "net: transport error", log.F("error", err))
		}
	}
}
