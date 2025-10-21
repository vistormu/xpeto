//go:build !headless

package xp

import (
	"github.com/vistormu/xpeto/pkg"

	"github.com/vistormu/xpeto/pkg/asset"
	"github.com/vistormu/xpeto/pkg/font"
	"github.com/vistormu/xpeto/pkg/image"
	"github.com/vistormu/xpeto/pkg/input"
	"github.com/vistormu/xpeto/pkg/sprite"
	"github.com/vistormu/xpeto/pkg/text"
	"github.com/vistormu/xpeto/pkg/vector"
)

// ################
// DEFAULT PACKAGES
// ################

var DefaultPkgs = pkg.DefaultPkgs

// =====
// asset
// =====

// ---------
// condition
// ---------

// description
func IsAssetLoaded[B any]() ConditionFn {
	return asset.IsAssetLoaded[B]()
}

// -----
// event
// -----

// ------
// handle
// ------

// description
type Handle = asset.Handle

// ------
// loader
// ------

// ------
// server
// ------

// description
var AddFileSystem = asset.AddFileSystem

// description
var AddAssetLoader = asset.AddAssetLoader

// description
func AddAssets[T, B any](w *World) {
	asset.AddAssets[T, B](w)
}

// description
func GetAsset[T any](w *World, handle Handle) (T, bool) {
	return asset.GetAsset[T](w, handle)
}

// =====
// audio
// =====

// ====
// font
// ====

// description
type Font = font.Font

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

// ======
// vector
// ======

// description
type Circle = vector.Circle

// description
type Rect = vector.Rect
