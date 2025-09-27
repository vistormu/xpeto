package asset

import (
	"io"
	"reflect"

	"github.com/vistormu/xpeto/internal/core"
)

// ======
// handle
// ======
type Handle = core.Id

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
type LoaderFn func(reader io.Reader) (any, error)

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
