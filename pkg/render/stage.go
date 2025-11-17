package render

type RenderStage uint8

const (
	Transparent RenderStage = iota
	Opaque
	Ui
	PostFx
)

var stagesOrder = []RenderStage{
	Transparent,
	Opaque,
	Ui,
	PostFx,
}
