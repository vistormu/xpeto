package shape

import (
	"image"
	"image/color"
	"math"
	"sort"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/vistormu/xpeto/backends/ebiten/shared"
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/window"
	"github.com/vistormu/xpeto/pkg/render"
	"github.com/vistormu/xpeto/pkg/shape"
	"github.com/vistormu/xpeto/pkg/transform"
)

type path struct {
	shape.Path
	transform.Transform

	snap      bool
	antialias bool
}

func extractPath(w *ecs.World) []path {
	q := ecs.NewQuery2[shape.Path, transform.Transform](w)

	sc, _ := ecs.GetResource[window.Scaling](w)
	rw, _ := ecs.GetResource[window.RealWindow](w)

	out := make([]path, 0)
	for _, b := range q.Iter() {
		p, t := b.Components()
		if !p.Visible || len(p.Points) < 2 {
			continue
		}
		out = append(out, path{
			Path:      *p,
			Transform: *t,
			snap:      sc.SnapPixels,
			antialias: rw.AntiAliasing,
		})
	}
	return out
}

func sortPath(p path) uint64 {
	return (uint64(p.Layer) << 16) | uint64(p.Order)
}

func drawPaths(screen *ebiten.Image, w *ecs.World) {
	paths := extractPath(w)
	sort.Slice(paths, func(i, j int) bool { return sortPath(paths[i]) < sortPath(paths[j]) })

	for _, p := range paths {
		drawPath(screen, p)
	}
}

var (
	whiteOnce sync.Once
	whiteImg  *ebiten.Image
)

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

type v2 struct {
	X, Y float64
}

func drawPath(screen *ebiten.Image, p path) {
	pts := p.Points
	if len(pts) < 2 {
		return
	}

	world := make([]v2, 0, len(pts))

	if p.Anchor == render.AnchorTopLeft {
		tlx := p.X
		tly := p.Y
		for _, pt := range pts {
			world = append(world, v2{
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

		ax, ay := shared.Offset(bw, bh, p.Anchor)

		tlx := p.X + ax
		tly := p.Y + ay

		for _, pt := range pts {
			world = append(world, v2{
				X: tlx + float64(pt.X-minX),
				Y: tly + float64(pt.Y-minY),
			})
		}
	}

	if p.snap {
		for i := range world {
			world[i].X = math.Round(world[i].X)
			world[i].Y = math.Round(world[i].Y)
		}
	}
	world = dedupConsecutive(world)
	if len(world) < 2 {
		return
	}

	for _, s := range p.Stroke {
		if !s.Visible || s.Width <= 0 {
			continue
		}
		drawStrokeFlushed(screen, world, s, p.antialias)
	}
}

func dedupConsecutive(pts []v2) []v2 {
	if len(pts) < 2 {
		return pts
	}

	out := make([]v2, 0, len(pts))
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
	pts []v2,
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
