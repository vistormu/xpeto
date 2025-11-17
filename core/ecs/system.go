package ecs

type systemInfo struct {
	id    uint64
	label string
}

type System = func(w *World)

func GetSystemId(w *World) uint64 {
	sys, _ := GetResource[systemInfo](w)
	return sys.id
}

func GetSystemLabel(w *World) string {
	sys, _ := GetResource[systemInfo](w)
	return sys.label
}

// THIS FUNCTION SHOULDN'T BE USED
//
// it should only be used by the scheduler to assign the current running system
func SetSystemInfo(w *World, id uint64, label string) {
	sys, _ := GetResource[systemInfo](w)
	sys.id = id
	sys.label = label
}
