package schedule

type label struct {
	id    uint64
	stage stage
}

type labelIndex struct {
	values map[string]label
}

func newLabelIndex() *labelIndex {
	return &labelIndex{
		values: make(map[string]label),
	}
}

func (li *labelIndex) add(l string, id uint64, stage stage) bool {
	_, ok := li.values[l]
	if ok {
		return false
	}

	li.values[l] = label{
		id:    id,
		stage: stage,
	}

	return true
}

func (li *labelIndex) get(label string) (label, bool) {
	ref, ok := li.values[label]
	return ref, ok
}
