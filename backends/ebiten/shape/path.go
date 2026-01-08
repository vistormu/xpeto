package shape

import (
	"image"
	"image/color"
	"math"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/vistormu/go-dsa/geometry"
	"github.com/vistormu/xpeto/backends/ebiten/shared"
	"github.com/vistormu/xpeto/pkg/render"
	"github.com/vistormu/xpeto/pkg/shape"
)

var (
	whiteOnce sync.Once
	whiteImg  *ebiten.Image
)

type v = geometry.Vector[float64]

func solidWhite() *ebiten.Image {
	whiteOnce.Do(func() {
		img := ebiten.NewImage(3, 3)
		img.Fill(color.White)
		whiteImg = img.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
	})
	return whiteImg
}

const (
	maxPtsPerChunk   = 256
	maxVertsPerFlush = 60000
	maxIdxPerFlush   = 90000
)

func drawPath(screen *ebiten.Image, r renderable) {
	s := r.Shape
	p := r.Shape.Path
	tr := r.Transform

	pts := p.Points
	if len(pts) < 2 {
		return
	}

	// TODO: don't allocate every frame
	world := make([]v, 0, len(pts))

	if r.anchor == render.AnchorTopLeft {
		tlx := tr.X
		tly := tr.Y
		for _, pt := range pts {
			world = append(world, v{
				X: tlx + float64(pt.X),
				Y: tly + float64(pt.Y),
			})
		}
	} else {
		minX, maxX := pts[0].X, pts[0].X
		minY, maxY := pts[0].Y, pts[0].Y
		for _, pt := range pts[1:] {
			minX = min(minX, pt.X)
			minY = min(minY, pt.Y)
			maxX = max(maxX, pt.X)
			maxY = max(maxY, pt.Y)
		}

		bw := float64(maxX - minX)
		bh := float64(maxY - minY)

		ax, ay := shared.Offset(bw, bh, r.anchor)

		tlx := tr.X + ax
		tly := tr.Y + ay

		for _, pt := range pts {
			world = append(world, v{
				X: tlx + float64(pt.X-minX),
				Y: tly + float64(pt.Y-minY),
			})
		}
	}

	if r.snap {
		for i := range world {
			world[i].X = math.Round(world[i].X)
			world[i].Y = math.Round(world[i].Y)
		}
	}
	world = dedupConsecutive(world)
	if len(world) < 2 {
		return
	}

	for _, s := range s.Strokes {
		if s.Width <= 0 {
			continue
		}
		drawStrokeFlushed(screen, world, s, r.antialias)
	}
}

func dedupConsecutive(pts []v) []v {
	if len(pts) < 2 {
		return pts
	}

	out := make([]v, 0, len(pts))
	out = append(out, pts[0])

	lastX := pts[0].X
	lastY := pts[0].Y

	for _, pt := range pts[1:] {
		if pt.X == lastX && pt.Y == lastY {
			continue
		}
		out = append(out, pt)
		lastX = pt.X
		lastY = pt.Y
	}

	return out
}

func drawStrokeFlushed(
	screen *ebiten.Image,
	pts []v,
	s shape.Stroke,
	antialias bool,
) {
	src := solidWhite()

	r16, g16, b16, a16 := color.RGBAModel.Convert(s.Color).RGBA()
	const inv = 1.0 / 65535.0
	r := float32(float64(r16) * inv)
	g := float32(float64(g16) * inv)
	b := float32(float64(b16) * inv)
	a := float32(float64(a16) * inv)

	st := &vector.StrokeOptions{
		Width:    s.Width,
		LineJoin: vector.LineJoinBevel,
	}

	draw := &ebiten.DrawTrianglesOptions{
		AntiAlias: antialias,
	}

	var verts []ebiten.Vertex
	var idx []uint16

	flush := func() {
		if len(idx) == 0 || len(verts) == 0 {
			verts = verts[:0]
			idx = idx[:0]
			return
		}

		for i := range verts {
			verts[i].ColorR = r
			verts[i].ColorG = g
			verts[i].ColorB = b
			verts[i].ColorA = a
		}

		screen.DrawTriangles(verts, idx, src, draw)

		verts = verts[:0]
		idx = idx[:0]
	}

	for start := 0; start < len(pts)-1; {
		end := start + maxPtsPerChunk
		if end > len(pts) {
			end = len(pts)
		}

		seg := pts[start:end]
		if len(seg) >= 2 {
			var vPath vector.Path
			vPath.MoveTo(float32(seg[0].X), float32(seg[0].Y))
			for _, pt := range seg[1:] {
				vPath.LineTo(float32(pt.X), float32(pt.Y))
			}

			verts, idx = vPath.AppendVerticesAndIndicesForStroke(verts, idx, st)

			if len(verts) >= maxVertsPerFlush || len(idx) >= maxIdxPerFlush {
				flush()
			}
		}

		if end == len(pts) {
			break
		}
		start = end - 1
	}

	flush()
}
