package ecs

type systemId struct {
	id uint64
}

type System = func(w *World)

func GetSystemId(w *World) uint64 {
	sys, _ := GetResource[systemId](w)
	return sys.id
}

// THIS FUNCTION SHOULDN'T BE USED
//
// it should only be used by the scheduler to assign the current running system
func SetSystemId(w *World, id uint64) {
	sys, _ := GetResource[systemId](w)
	sys.id = id
}
