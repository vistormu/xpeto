package font

import (
	"image/color"
	"math"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/vistormu/xpeto/backends/ebiten/shared"
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/log"
	"github.com/vistormu/xpeto/core/window"
	"github.com/vistormu/xpeto/pkg/asset"
	"github.com/vistormu/xpeto/pkg/render"
	xptext "github.com/vistormu/xpeto/pkg/text"
	"github.com/vistormu/xpeto/pkg/transform"
)

// ==========
// renderable
// ==========
type renderable struct {
	face     text.Face
	content  string
	align    xptext.Align
	wrap     xptext.WrapMode
	maxWidth float64
	anchor   render.AnchorType
	color    color.Color
	x, y     float64
	orderKey render.OrderKey
	snap     bool
}

type renderableBuffer struct {
	renderables []renderable
}

func newRenderableBuffer() renderableBuffer {
	return renderableBuffer{
		renderables: make([]renderable, 0),
	}
}

// =======
// helpers
// =======
var xpToEbiAlign = map[xptext.Align]text.Align{
	xptext.AlignStart:  text.AlignStart,
	xptext.AlignCenter: text.AlignCenter,
	xptext.AlignEnd:    text.AlignEnd,
}

func wrapContent(content string, face text.Face, maxWidth float64, mode xptext.WrapMode) string {
	if maxWidth <= 0 || mode == xptext.WrapNone || content == "" {
		return content
	}

	// preserve explicit newlines and wrap each paragraph independently.
	parts := strings.Split(content, "\n")
	for i, p := range parts {
		parts[i] = wrapLine(p, face, maxWidth, mode)
	}
	return strings.Join(parts, "\n")
}

func wrapLine(line string, face text.Face, maxWidth float64, mode xptext.WrapMode) string {
	line = strings.TrimRight(line, "\r")
	if line == "" {
		return ""
	}

	if mode == xptext.WrapRune {
		return wrapRunes(line, face, maxWidth)
	}

	// word wrap (collapses whitespace).
	words := strings.Fields(line)
	if len(words) == 0 {
		return ""
	}

	lines := make([]string, 0, 4)
	cur := ""
	for _, w := range words {
		next := w
		if cur != "" {
			next = cur + " " + w
		}

		mw, _ := text.Measure(next, face, 0)
		if mw <= maxWidth {
			cur = next
			continue
		}

		if cur != "" {
			lines = append(lines, cur)
			cur = ""
		}

		// single word longer than max width: fall back to rune wrapping for this word.
		sw, _ := text.Measure(w, face, 0)
		if sw > maxWidth {
			wrapped := wrapRunes(w, face, maxWidth)
			if wrapped != "" {
				lines = append(lines, strings.Split(wrapped, "\n")...)
			}
			continue
		}

		cur = w
	}
	if cur != "" {
		lines = append(lines, cur)
	}

	return strings.Join(lines, "\n")
}

func wrapRunes(s string, face text.Face, maxWidth float64) string {
	r := []rune(s)
	lines := make([]string, 0, 4)
	cur := make([]rune, 0, len(r))
	for _, ch := range r {
		next := string(append(cur, ch))
		mw, _ := text.Measure(next, face, 0)
		if mw <= maxWidth || len(cur) == 0 {
			cur = append(cur, ch)
			continue
		}

		lines = append(lines, string(cur))
		cur = []rune{ch}
	}
	if len(cur) != 0 {
		lines = append(lines, string(cur))
	}
	return strings.Join(lines, "\n")
}

// ======
// render
// ======
func extractText(w *ecs.World) []renderable {
	buf := ecs.EnsureResource(w, newRenderableBuffer)
	buf.renderables = buf.renderables[:0]

	sc, ok := ecs.GetResource[window.Scaling](w)
	if !ok {
		return buf.renderables
	}

	q := ecs.NewQuery2[xptext.Text, transform.Transform](w)

	for _, b := range q.Iter() {
		txt, tr := b.Components()
		e := b.Entity()

		txtFont := txt.Font
		if txtFont == asset.Asset(0) {
			df, ok := xptext.GetDefaultFont(w)
			if !ok {
				continue
			}
			txtFont = df
		}

		fnt, ok := asset.GetAsset[font](w, txtFont)
		if !ok || fnt == nil {
			log.LogError(w, "tried to load a missing font",
				log.F("function", "backends/ebiten/font/render.go:extractText"),
				log.F("asset", uint64(txtFont)),
			)
			continue
		}

		anchor := render.AnchorCenter
		an, ok := ecs.GetComponent[render.Anchor](w, e)
		if ok {
			anchor = an.Type
		}

		face := fnt.Face(txt.Size)

		content := txt.Content
		if txt.Wrap != xptext.WrapNone && txt.MaxWidth > 0 {
			content = wrapContent(txt.Content, face, txt.MaxWidth, txt.Wrap)
		}

		buf.renderables = append(buf.renderables, renderable{
			face:     face,
			content:  content,
			align:    txt.Align,
			wrap:     txt.Wrap,
			maxWidth: txt.MaxWidth,
			anchor:   anchor,
			color:    txt.Color,
			x:        tr.X,
			y:        tr.Y,
			orderKey: txt.OrderKey,
			snap:     sc.SnapPixels,
		})
	}

	return buf.renderables
}

func sortText(t renderable) uint64 {
	return uint64(t.orderKey)
}

func drawText(screen *ebiten.Image, t renderable) {
	op := &text.DrawOptions{}

	op.ColorScale.ScaleWithColor(t.color)
	op.PrimaryAlign = xpToEbiAlign[t.align]
	op.SecondaryAlign = xpToEbiAlign[t.align]

	w, h := text.Measure(t.content, t.face, 0)
	dx, dy := shared.Offset(w, h, t.anchor)

	x := t.x + dx
	y := t.y + dy

	if t.snap {
		x = math.Round(x)
		y = math.Round(y)
	}

	op.GeoM.Translate(x, y)
	text.Draw(screen, t.content, t.face, op)
}
