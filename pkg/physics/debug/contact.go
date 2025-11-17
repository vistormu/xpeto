package debug

// import (
// 	"image/color"

// 	"github.com/hajimehoshi/ebiten/v2"
// 	"github.com/hajimehoshi/ebiten/v2/vector"

// 	"github.com/vistormu/xpeto/core/ecs"
// 	"github.com/vistormu/xpeto/core/pkg/transform"

// 	"github.com/vistormu/xpeto/pkg/physics"
// 	"github.com/vistormu/xpeto/pkg/render"
// )

// // ====
// // item
// // ====
// type contactItem struct {
// 	x, y   float32
// 	radius float32
// 	color  color.Color
// 	key    uint64
// }

// func (li *contactItem) Draw(screen *ebiten.Image) {
// 	vector.FillCircle(screen, li.x, li.y, li.radius, li.color, false)
// }

// func (li *contactItem) SortKey() uint64 { return li.key }

// // =========
// // extractor
// // =========
// func extractContacts(w *ecs.World) []render.Renderable {
// 	s, _ := ecs.GetResource[Settings](w)
// 	sp, _ := ecs.GetResource[physics.Space](w)
// 	pairs := sp.Contacts

// 	items := make([]render.Renderable, 0)

// 	if !s.Enabled || !s.DrawContacts {
// 		return items
// 	}

// 	key := packKey(s.Layer, s.Order)

// 	for _, c := range pairs {
// 		tra, okA := ecs.GetComponent[transform.Transform](w, c.A)
// 		trb, okB := ecs.GetComponent[transform.Transform](w, c.B)
// 		if !okA || !okB {
// 			continue
// 		}

// 		mx := float32((tra.X + trb.X) * 0.5)
// 		my := float32((tra.Y + trb.Y) * 0.5)

// 		items = append(items, &contactItem{
// 			x:      mx,
// 			y:      my,
// 			radius: s.ContactRadiusPx,
// 			color:  s.ContactFill,
// 			key:    key,
// 		})
// 	}

// 	return items
// }
