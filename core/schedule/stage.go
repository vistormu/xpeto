package schedule

type Stage = func(*storage, uint64) stage
type stage uint32

const (
	empty stage = iota
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

func (s stage) String() string {
	switch s {
	case preStartup:
		return "PreStartup"
	case startup:
		return "Startup"
	case postStartup:
		return "PostStartup"
	case first:
		return "First"
	case preUpdate:
		return "PreUpdate"
	case stateTransition:
		return "StateTransition"
	case fixedFirst:
		return "FixedFirst"
	case fixedPreUpdate:
		return "FixedPreUpdate"
	case fixedUpdate:
		return "FixedUpdate"
	case fixedPostUpdate:
		return "FixedPostUpdate"
	case fixedLast:
		return "FixedLast"
	case update:
		return "Update"
	case postUpdate:
		return "PostUpdate"
	case last:
		return "Last"
	case preDraw:
		return "PreDraw"
	case draw:
		return "Draw"
	case postDraw:
		return "PostDraw"
	case exit:
		return "Exit"
	default:
		return "Unknown"
	}
}

// startup
func PreStartup(*storage, uint64) stage  { return preStartup }
func Startup(*storage, uint64) stage     { return startup }
func PostStartup(*storage, uint64) stage { return postStartup }

// update
func First(*storage, uint64) stage     { return first }
func PreUpdate(*storage, uint64) stage { return preUpdate }

// states
func stateTransitionStage(*storage, uint64) stage { return stateTransition }

func OnExit[T comparable](from T) Stage {
	return func(store *storage, id uint64) stage {
		sm, ok := getStateMachine[T](store)
		if !ok {
			return empty
		}

		sm.add(&from, nil, id)

		return empty
	}
}
func OnTransition[T comparable](from, to T) Stage {
	return func(store *storage, id uint64) stage {
		sm, ok := getStateMachine[T](store)
		if !ok {
			return empty
		}

		sm.add(&from, &to, id)

		return empty
	}
}
func OnEnter[T comparable](to T) Stage {
	return func(store *storage, id uint64) stage {
		sm, ok := getStateMachine[T](store)
		if !ok {
			return empty
		}

		sm.add(nil, &to, id)

		return empty
	}
}

func FixedFirst(*storage, uint64) stage      { return fixedFirst }
func FixedPreUpdate(*storage, uint64) stage  { return fixedPreUpdate }
func FixedUpdate(*storage, uint64) stage     { return fixedUpdate }
func FixedPostUpdate(*storage, uint64) stage { return fixedPostUpdate }
func FixedLast(*storage, uint64) stage       { return fixedLast }

func Update(*storage, uint64) stage     { return update }
func PostUpdate(*storage, uint64) stage { return postUpdate }
func Last(*storage, uint64) stage       { return last }

// draw
func PreDraw(*storage, uint64) stage  { return preDraw }
func Draw(*storage, uint64) stage     { return draw }
func PostDraw(*storage, uint64) stage { return postDraw }

// exit
func Exit(*storage, uint64) stage { return exit }
