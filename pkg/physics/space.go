package physics

import (
	"math"

	"github.com/vistormu/xpeto/core/ecs"
)

// ========
// contacts
// ========
type ContactPair struct {
	A, B             ecs.Entity
	PenX, PenY       float64
	NormalX, NormalY float64
	Depth            float64
}

type pair struct {
	A, B ecs.Entity
}

// ====
// cell
// ====
type Cell struct {
	I        int
	J        int
	Entities []ecs.Entity
}

// =====
// space
// =====
type Space struct {
	Height float64
	Width  float64

	CellWidth  float64
	CellHeight float64
	rows, cols int
	Cells      []*Cell // or map[[2]int]Cell?

	candidates []pair
	Contacts   []ContactPair
}

func (s *Space) AddEntity(e ecs.Entity, i, j int) bool {
	inBounds := i >= 0 && i < s.cols && j >= 0 && j < s.rows
	if !inBounds {
		return false
	}

	index := j*s.cols + i
	s.Cells[index].Entities = append(s.Cells[index].Entities, e)

	return true
}

func (s *Space) GetCell(i, j int) (*Cell, bool) {
	if i < 0 || i >= s.cols || j < 0 || j >= s.rows {
		return nil, false
	}

	index := j*s.cols + i

	return s.Cells[index], true
}

func (s *Space) IsEmpty(i, j int) bool {
	if i < 0 || i >= s.cols || j < 0 || j >= s.rows {
		return true
	}

	cell, ok := s.GetCell(i, j)
	if !ok {
		return true
	}

	return len(cell.Entities) == 0
}

func (s *Space) Clear() {
	for i := range s.Cells {
		s.Cells[i].Entities = s.Cells[i].Entities[:0]
	}

	s.Contacts = s.Contacts[:0]
	s.candidates = s.candidates[:0]
}

type lastSpaceSize struct {
	width      float64
	height     float64
	cellWidth  float64
	cellHeight float64
}

const defaultCellEntitiesCap = 8

func resizeSpace(w *ecs.World) {
	s, _ := ecs.GetResource[Space](w)
	if s == nil {
		return
	}
	ls, _ := ecs.GetResource[lastSpaceSize](w)

	changed := func() bool {
		// no previous snapshot → force build if sizes are valid
		if ls == nil {
			return true
		}
		return ls.width != s.Width ||
			ls.height != s.Height ||
			ls.cellWidth != s.CellWidth ||
			ls.cellHeight != s.CellHeight
	}()

	// if any size is non-positive, disable the grid
	if s.Width <= 0 || s.Height <= 0 || s.CellWidth <= 0 || s.CellHeight <= 0 {
		s.rows, s.cols = 0, 0
		s.Cells = s.Cells[:0]
		if ls == nil {
			ecs.AddResource(w, lastSpaceSize{
				width:      s.Width,
				height:     s.Height,
				cellWidth:  s.CellWidth,
				cellHeight: s.CellHeight,
			})
		} else {
			ls.width, ls.height = s.Width, s.Height
			ls.cellWidth, ls.cellHeight = s.CellWidth, s.CellHeight
		}
		return
	}

	if !changed {
		return // nothing to do
	}

	cols := int(math.Ceil(s.Width / s.CellWidth))
	rows := int(math.Ceil(s.Height / s.CellHeight))
	if cols < 1 {
		cols = 1
	}
	if rows < 1 {
		rows = 1
	}

	// (re)build cells
	newCount := rows * cols
	if cap(s.Cells) < newCount {
		s.Cells = make([]*Cell, newCount)
	} else {
		s.Cells = s.Cells[:newCount]
	}

	for j := 0; j < rows; j++ {
		for i := 0; i < cols; i++ {
			idx := j*cols + i
			c := s.Cells[idx]
			if c == nil {
				c = &Cell{I: i, J: j, Entities: make([]ecs.Entity, 0, defaultCellEntitiesCap)}
				s.Cells[idx] = c
			} else {
				// recycle existing cell object; reset indices and clear entities
				c.I, c.J = i, j
				c.Entities = c.Entities[:0]
			}
		}
	}

	s.rows, s.cols = rows, cols
	// also clear runtime buffers since the grid layout changed
	s.candidates = s.candidates[:0]
	s.Contacts = s.Contacts[:0]

	// update or create snapshot
	if ls == nil {
		ecs.AddResource(w, lastSpaceSize{
			width:      s.Width,
			height:     s.Height,
			cellWidth:  s.CellWidth,
			cellHeight: s.CellHeight,
		})
	} else {
		ls.width, ls.height = s.Width, s.Height
		ls.cellWidth, ls.cellHeight = s.CellWidth, s.CellHeight
	}
}
