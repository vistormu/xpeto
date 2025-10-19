package debug

import (
	"image/color"
)

type Settings struct {
	Enabled          bool
	DrawAABBs        bool
	DrawVelocities   bool
	DrawContacts     bool
	DrawOccupiedGrid bool

	Layer uint16
	Order uint16

	AABBStroke      color.Color
	VelocityColor   color.Color
	ContactFill     color.Color
	GridStroke      color.Color
	LineWidthPx     float32
	VelocityScale   float32 // pixels per (unit/sec)
	ContactRadiusPx float32
}

func defaultSettings() Settings {
	return Settings{
		Enabled:          true,
		DrawAABBs:        true,
		DrawVelocities:   true,
		DrawContacts:     true,
		DrawOccupiedGrid: true,
		Layer:            65000,
		Order:            0,

		AABBStroke:      color.RGBA{0, 255, 0, 200},
		VelocityColor:   color.RGBA{0, 150, 255, 200},
		ContactFill:     color.RGBA{255, 50, 50, 220},
		GridStroke:      color.RGBA{255, 255, 0, 160},
		LineWidthPx:     1.5,
		VelocityScale:   0.05,
		ContactRadiusPx: 3,
	}
}
