package sink

import (
	"fmt"
	"strings"

	"github.com/vistormu/go-dsa/terminal"

	"github.com/vistormu/xpeto/core/log"
)

var levelToColor = map[log.Level]string{
	log.Debug:   terminal.FgBlue,
	log.Info:    terminal.FgGreen,
	log.Warning: terminal.FgYellow,
	log.Error:   terminal.FgRed,
	log.Fatal:   terminal.BgRed,
}

var levelToString = map[log.Level]string{
	log.Debug:   "  [debug]",
	log.Info:    "   [info]",
	log.Warning: "[warning]",
	log.Error:   "  [error]",
	log.Fatal:   "  [fatal]",
}

type TerminalSink struct{}

func (s *TerminalSink) Write(frame uint64, records []log.Record) {
	if len(records) == 0 {
		return
	}

	var b strings.Builder
	b.Grow(256 + len(records)*128)

	fmt.Fprintf(&b, "\n=== frame %d ===\n", frame)

	for _, r := range records {
		level := min(r.Level, log.MaxLevel)

		fmt.Fprintf(&b, "%s%s %s%s\n",
			levelToColor[level],
			levelToString[level],
			r.Message,
			terminal.StyleReset,
		)

		fmt.Fprintf(&b, "   |> system: %s (id: %d)\n", r.SystemLabel, r.SystemId)
		fmt.Fprintf(&b, "   |> time:   %s\n", r.Time.String())

		if r.Caller.File != "" {
			fmt.Fprintf(&b, "   |> caller: %s:%d (%s)\n", r.Caller.File, r.Caller.Line, r.Caller.Func)
		}

		for _, f := range r.Fields {
			fmt.Fprintf(&b, "   |> %s\n", f)
		}
	}

	fmt.Print(b.String())
}

func (s *TerminalSink) Flush() error { return nil }

func (s *TerminalSink) Sync() error { return nil }
