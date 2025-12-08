package game

import (
	"image/color"
	"strconv"

	"github.com/vistormu/xpeto"
)

type Fonts struct {
	Regular xp.Asset `path:"default/font.ttf"`
}

// ====
// menu
// ====
type menuUi struct{}

func createMenuUi(w *xp.World) {
	fonts, _ := xp.GetResource[Fonts](w)
	ww, wh := xp.GetVirtualWindowSize[float64](w)

	e := xp.AddEntity(w)
	xp.AddComponent(w, e, xp.Text{
		Font:    fonts.Regular,
		Content: "press enter to start",
		Align:   xp.AlignCenter,
		Color:   color.RGBA{221, 199, 161, 255},
		Size:    24,
		Layer:   2,
	})
	xp.AddComponent(w, e, xp.Transform{
		X: ww / 2,
		Y: wh / 2,
	})
	xp.AddComponent(w, e, menuUi{})
}

func hideMenuUi(w *xp.World) {
	b, _ := xp.Query1[menuUi](w).First()
	xp.RemoveEntity(w, b.Entity())
}

// =====
// score
// =====
type scoreUi struct {
	IsLeft bool
}

func createScoreUi(w *xp.World) {
	fonts, _ := xp.GetResource[Fonts](w)
	score, _ := xp.GetResource[Score](w)

	ww, wh := xp.GetVirtualWindowSize[float64](w)

	e := xp.AddEntity(w)
	xp.AddComponent(w, e, xp.Text{
		Font:    fonts.Regular,
		Content: strconv.Itoa(score.Left),
		Align:   xp.AlignEnd,
		Color:   color.RGBA{221, 199, 161, 255},
		Size:    32,
		Layer:   2,
	})
	xp.AddComponent(w, e, xp.Transform{
		X: ww/2 - 30,
		Y: wh / 2,
	})
	xp.AddComponent(w, e, scoreUi{
		IsLeft: true,
	})

	e = xp.AddEntity(w)
	xp.AddComponent(w, e, xp.Text{
		Font:    fonts.Regular,
		Content: strconv.Itoa(score.Left),
		Align:   xp.AlignStart,
		Color:   color.RGBA{221, 199, 161, 255},
		Size:    32,
		Layer:   2,
	})
	xp.AddComponent(w, e, xp.Transform{
		X: ww/2 + 30,
		Y: wh / 2,
	})
	xp.AddComponent(w, e, scoreUi{
		IsLeft: false,
	})
}

func updateScoreUi(w *xp.World) {
	_, ok := xp.GetEvents[ScoreEvent](w)
	if !ok {
		return
	}

	score, _ := xp.GetResource[Score](w)

	q := xp.Query2[xp.Text, scoreUi](w)
	for _, b := range q.Iter() {
		t, ui := b.Components()

		if ui.IsLeft {
			t.Content = strconv.Itoa(score.Left)
		} else {
			t.Content = strconv.Itoa(score.Right)
		}
	}
}

func hideScoreUi(w *xp.World) {
	q := xp.Query1[scoreUi](w)
	for _, b := range q.Iter() {
		xp.RemoveEntity(w, b.Entity())
	}
}

func uiMiniPkg(w *xp.World, sch *xp.Scheduler) {
	xp.AddAsset[Fonts](w)

	xp.AddSystem(sch, xp.OnEnter(stateGameOver), createMenuUi)
	xp.AddSystem(sch, xp.OnEnter(statePlaying), hideMenuUi)

	xp.AddSystem(sch, xp.OnEnter(statePlaying), createScoreUi)
	xp.AddSystem(sch, xp.Update, updateScoreUi)
	xp.AddSystem(sch, xp.OnEnter(stateGameOver), hideScoreUi)
}
