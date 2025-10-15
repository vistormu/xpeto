package schedule

type StageFn = func(*Scheduler, *Schedule) Stage
type Stage uint32

const (
	empty Stage = iota
	preStartup
	startup
	postStartup

	first
	preUpdate

	stateTransition

	fixedFirst
	fixedPreUpdate
	fixedUpdate
	fixedPostUpdate
	fixedLast

	update
	postUpdate
	last

	preDraw
	draw
	postDraw

	exit
)

// startup
func PreStartup(*Scheduler, *Schedule) Stage  { return preStartup }
func Startup(*Scheduler, *Schedule) Stage     { return startup }
func PostStartup(*Scheduler, *Schedule) Stage { return postStartup }

// update
func First(*Scheduler, *Schedule) Stage     { return first }
func PreUpdate(*Scheduler, *Schedule) Stage { return preUpdate }

func StateTransition(*Scheduler, *Schedule) Stage { return stateTransition }

func FixedFirst(*Scheduler, *Schedule) Stage      { return fixedFirst }
func FixedPreUpdate(*Scheduler, *Schedule) Stage  { return fixedPreUpdate }
func FixedUpdate(*Scheduler, *Schedule) Stage     { return fixedUpdate }
func FixedPostUpdate(*Scheduler, *Schedule) Stage { return fixedPostUpdate }
func FixedLast(*Scheduler, *Schedule) Stage       { return fixedLast }

func Update(*Scheduler, *Schedule) Stage     { return update }
func PostUpdate(*Scheduler, *Schedule) Stage { return postUpdate }
func Last(*Scheduler, *Schedule) Stage       { return last }

// draw
func PreDraw(*Scheduler, *Schedule) Stage  { return preDraw }
func Draw(*Scheduler, *Schedule) Stage     { return draw }
func PostDraw(*Scheduler, *Schedule) Stage { return postDraw }

// exit
func Exit(*Scheduler, *Schedule) Stage { return exit }
