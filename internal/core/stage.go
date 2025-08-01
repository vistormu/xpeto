package core

type Stage uint32

const (
	PreStartup Stage = iota
	Startup
	PostStartup

	First
	PreUpdate

	FixedFirst
	FixedPreUpdate
	FixedUpdate
	FixedPostUpdate
	FixedLast

	Update
	PostUpdate
	Last

	Draw
)

func StartupStages() []Stage {
	return []Stage{
		PreStartup,
		Startup,
		PostStartup,
	}
}

func UpdateStages() []Stage {
	return []Stage{
		First,
		PreUpdate,
		FixedFirst,
		FixedPreUpdate,
		FixedUpdate,
		FixedPostUpdate,
		FixedLast,
		Update,
		PostUpdate,
		Last,
	}
}

func DrawStages() []Stage {
	return []Stage{
		Draw,
	}
}
