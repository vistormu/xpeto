package net

import (
	"bytes"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/pkg/log"
)

type Outbox[T any] struct {
	Data []T
}

func send[T any](w *ecs.World) {
	q := ecs.NewQuery2[Connection, Outbox[T]](w)

	var buf bytes.Buffer

	for _, b := range q.Iter() {
		conn, outbox := b.Components()

		if len(outbox.Data) == 0 {
			continue
		}

		tr := conn.Channel.Transport
		cd := conn.Channel.Codec

		if tr == nil || cd == nil {
			log.LogError(w, "channel not initialized")
			continue
		}

		for _, msg := range outbox.Data {
			buf.Reset()

			err := cd.Encode(&buf, msg)
			if err != nil {
				log.LogError(w, "error encoding message", log.F("error", err))
				continue
			}

			tr.Send(conn.Address, buf.Bytes())
		}

		outbox.Data = outbox.Data[:0]
	}
}
