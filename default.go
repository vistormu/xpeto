//go:build !headless

package xp

import (
	"github.com/vistormu/xpeto/pkg"

	"github.com/vistormu/xpeto/pkg/asset"
	"github.com/vistormu/xpeto/pkg/font"
	"github.com/vistormu/xpeto/pkg/image"
	"github.com/vistormu/xpeto/pkg/input"
	"github.com/vistormu/xpeto/pkg/shape"
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

const (
	AlignStart  = font.AlignStart
	AlignCenter = font.AlignCenter
	AlignEnd    = font.AlignEnd
)

type Text = font.Text

// ========
// geometry
// ========

// =====
// image
// =====

// description
type Image = image.Image

// description
type Sprite = image.Sprite

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

// description
type EllipeShape = shape.Ellipse

var NewCircleShape = shape.NewCircle

// description
type RectShape = shape.Rect

var NewRectShape = shape.NewRect

// description
type PathShape = shape.Path

var NewPathShape = shape.NewPath

// =========
// transform
// =========

type Transform = transform.Transform
