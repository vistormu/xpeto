package net

import (
	"reflect"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/log"
	"github.com/vistormu/xpeto/core/schedule"

	"github.com/vistormu/xpeto/pkg/net/codec"
	"github.com/vistormu/xpeto/pkg/net/transport"
)

// =====
// types
// =====
type Channel struct {
	Transport transport.Transport
	Codec     codec.Codec
}

type Connection struct {
	Address string
	Channel Channel
}

// =======
// helpers
// =======
func baseType(t reflect.Type) reflect.Type {
	if t == nil {
		return nil
	}
	if t.Kind() == reflect.Pointer {
		return t.Elem()
	}
	return t
}

// ===
// API
// ===
func AddChannel[T any](w *ecs.World) {
	session, ok := ecs.GetResource[session](w)
	if !ok {
		log.LogError(w, "cannot use AddChannel if net.Pkg is not included")
		return
	}

	b := new(T)
	bType := baseType(reflect.TypeFor[T]())
	bValue := reflect.ValueOf(b).Elem()

	if bType.Kind() != reflect.Struct {
		log.LogError(w, "the channel bundle must be a struct", log.F("got", bType.Kind().String()))
		return
	}

	// iterate over the bundle
	for i := range bValue.NumField() {
		fValue := bValue.Field(i)
		fType := bType.Field(i)

		protocol := fType.Tag.Get("protocol")
		encoding := fType.Tag.Get("codec")
		addr := fType.Tag.Get("listen")

		if !fValue.CanSet() {
			log.LogError(w, "field must be settable", log.F("field", bType.Field(i).Name))
			continue
		}

		if fValue.Type() != reflect.TypeFor[Channel]() {
			log.LogError(w, "all fields of the channel bundle must be of type net.Channel")
			continue
		}

		// set fields
		cod := codec.New(w, encoding)
		if cod == nil {
			log.LogError(w, "error creating the codec")
			continue
		}

		tr := transport.New(w, protocol)
		if tr == nil {
			log.LogError(w, "error creating the transport")
			continue
		}

		if addr == "" {
			addr = ":0"
		}

		err := tr.Listen(addr)
		if err != nil {
			log.LogError(w, "could not listen to specified address", log.F("address", addr))
			continue
		}

		ch := Channel{Transport: tr, Codec: cod}
		fValue.Set(reflect.ValueOf(ch))
		session.channels = append(session.channels, ch)
	}

	// add bundle to the resources
	ecs.AddResource(w, b)
}

func AddMessage[T any](sch *schedule.Scheduler, label string) {
	schedule.AddSystem(sch, schedule.PreUpdate, receive[T]).Label("net.receive." + label)
	schedule.AddSystem(sch, schedule.PostUpdate, send[T]).Label("net.send." + label)
}
