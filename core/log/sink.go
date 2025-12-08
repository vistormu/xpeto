package log

import (
	"fmt"

	"github.com/vistormu/go-dsa/ansi"
)

type sink interface {
	write(frame uint64, records []record)
}

type debugSink struct{}

func (s *debugSink) write(frame uint64, records []record) {
	if len(records) == 0 {
		return
	}

	msg := fmt.Sprintf("\n=== frame %d ===\n", frame)

	for _, r := range records {
		msg += fmt.Sprintf("%s%s %s%s\n",
			levelToColor[r.level],
			levelToString[r.level],
			r.message,
			ansi.Reset,
		)

		msg += fmt.Sprintf("   |> system: %s (id: %d)\n", r.systemLabel, r.systemId)
		msg += fmt.Sprintf("   |> time:   %s\n", r.time.String())

		for _, f := range r.fields {
			msg += fmt.Sprintf("   |> %s: %v\n", f.key, f.value)
		}
	}

	fmt.Println(msg)
}
