package input

// ====
// mode
// ====
type AnalogMode uint8

const (
	AnalogAbsolute AnalogMode = iota
	AnalogTransient
)

// =====
// input
// =====
type AnalogInput struct {
	mode     AnalogMode
	value    float64
	previous float64
	delta    float64
}

func newAnalogInput(mode AnalogMode) AnalogInput {
	return AnalogInput{
		mode: mode,
	}
}

func (ai *AnalogInput) begin() {
	if ai.mode == AnalogTransient {
		ai.value = 0
		ai.delta = 0
		ai.previous = 0
	}
}

func (ai *AnalogInput) compute() {
	switch ai.mode {
	case AnalogAbsolute:
		ai.delta = ai.value - ai.previous
		ai.previous = ai.value
	case AnalogTransient:
		ai.delta = ai.value
		ai.previous = 0
	}
}

func (ai *AnalogInput) end() {
	if ai.mode == AnalogTransient {
		ai.value = 0
	}
}

func (ai *AnalogInput) add(v float64) {
	ai.value += v
}

func (ai *AnalogInput) set(v float64) {
	ai.value = v
}

func (ai *AnalogInput) reset() {
	switch ai.mode {
	case AnalogAbsolute:
		ai.delta = 0
		ai.previous = ai.value
	case AnalogTransient:
		ai.delta = 0
		ai.previous = 0
		ai.value = 0
	}
}

// ===
// API
// ===
func (ai *AnalogInput) Value() float64 {
	return ai.value
}

func (ai *AnalogInput) Delta() float64 {
	return ai.delta
}

func (ai *AnalogInput) Previous() float64 {
	return ai.previous
}
