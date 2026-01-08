package input

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

func Pkg(w *ecs.World, sch *schedule.Scheduler) {
	// keyboard
	ecs.AddResource(w, newKeyboardState())
	schedule.AddSystem(sch, schedule.PreUpdate, emitKeyboard,
		schedule.SystemOpt.Label("ebiten.input.emitKeyboard"),
		schedule.SystemOpt.RunBefore("input.applyEvents"),
	)

	// mouse
	schedule.AddSystem(sch, schedule.PreUpdate, emitMouse,
		schedule.SystemOpt.Label("ebiten.input.emitMouse"),
		schedule.SystemOpt.RunBefore("input.applyEvents"),
	)

	// gamepad
	ecs.AddResource(w, newGamepadState())
	schedule.AddSystem(sch, schedule.PreUpdate, emitGamepads,
		schedule.SystemOpt.Label("ebiten.input.emitGamepads"),
		schedule.SystemOpt.RunBefore("input.applyEvents"),
	)

	// text input
	schedule.AddSystem(sch, schedule.PreUpdate, emitText,
		schedule.SystemOpt.Label("ebiten.input.emitText"),
		schedule.SystemOpt.RunBefore("input.applyEvents"),
	)
}
