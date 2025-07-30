package asset

import (
	"io"
	"reflect"
)

// ======
// handle
// ======
type Handle struct {
	Id      uint32
	Version uint32
}

// ==========
// load state
// ==========
type LoadState uint8

const (
	NotFound LoadState = iota
	Loading
	Loaded
	Failed
)

// ======
// loader
// ======
type LoaderFn func(reader io.Reader, path string) (any, error)

// ============
// load request
// ============
type LoadRequest struct {
	Path       string
	Bundle     any
	BundleType reflect.Type
	Handle     Handle
	AssetType  reflect.Type
	LoaderFn   LoaderFn
}
