//go:build !headless

package xp

import (
	c "github.com/vistormu/go-dsa/constraints"

	"github.com/vistormu/xpeto/pkg"

	"github.com/vistormu/xpeto/pkg/asset"
	"github.com/vistormu/xpeto/pkg/font"
	"github.com/vistormu/xpeto/pkg/geometry"
	"github.com/vistormu/xpeto/pkg/image"
	"github.com/vistormu/xpeto/pkg/input"
	"github.com/vistormu/xpeto/pkg/shape"
	"github.com/vistormu/xpeto/pkg/sprite"
	"github.com/vistormu/xpeto/pkg/text"
	"github.com/vistormu/xpeto/pkg/transform"
)

// ################
// DEFAULT PACKAGES
// ################

var DefaultPkgs = pkg.DefaultPkgs

// =====
// asset
// =====

// -----
// asset
// -----

// description
type Asset = asset.Asset

// ------
// loader
// ------

type LoaderFn[T any] = asset.LoaderFn[T]

// ------
// server
// ------

// description
var AddStaticFS = asset.AddStaticFS

// description
func AddLoaderFn[T any](w *World, fn asset.LoaderFn[T], extensions ...string) {
	asset.AddLoaderFn(w, fn, extensions...)
}

// description
func AddAsset[T any](w *World) {
	asset.AddAsset[T](w)
}

// description
func GetAsset[T any](w *World, a Asset) (*T, bool) {
	return asset.GetAsset[T](w, a)
}

// description
func RemoveAsset[T any](w *World, a Asset) bool {
	return asset.RemoveAsset[T](w, a)
}

// =====
// audio
// =====

// ====
// font
// ====

// description
type Font = font.Font

// ========
// geometry
// ========

// description
type Geometry[T c.Number] = geometry.Geometry[T]

// -------
// ellipse
// -------

// description
type Ellipse[T c.Number] = geometry.Ellipse[T]

// description
func NewCircle[T c.Number](r T) Geometry[T] {
	return geometry.NewCircle(r)
}

// ----
// rect
// ----

// description
func NewRect[T c.Number](w, h T) Geometry[T] {
	return geometry.NewRect(w, h)
}

// description

// =====
// image
// =====

// description
type Image = image.Image

// =====
// input
// =====

// ------
// events
// ------

// -------
// gamepad
// -------

// --------
// keyboard
// --------

// description
type Key = input.Key

const (
	KeyA = input.KeyA
	KeyD = input.KeyD
	KeyS = input.KeyS
	KeyW = input.KeyW

	KeyEnter = input.KeyEnter

	KeyArrowDown = input.KeyArrowDown
	KeyArrowUp   = input.KeyArrowUp
)

// description
type Keyboard = input.Keyboard

// -----
// mouse
// -----

// description
type Mouse = input.Mouse

// ======
// render
// ======

// =====
// shape
// =====

// description
type Shape = shape.Shape

// ======
// sprite
// ======

// description
type Sprite = sprite.Sprite

// ====
// text
// ====

const (
	AlignStart  = text.AlignStart
	AlignCenter = text.AlignCenter
	AlignEnd    = text.AlignEnd
)

type Text = text.Text

// =========
// transform
// =========

type Transform = transform.Transform
