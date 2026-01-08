package render

type RenderStage uint8

const (
	Opaque RenderStage = iota
	Transparent
	Ui
	PostFx
)

var stagesOrder = []RenderStage{
	Opaque,
	Transparent,
	Ui,
	PostFx,
}
