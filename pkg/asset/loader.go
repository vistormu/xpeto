package asset

import (
	"fmt"
	"io"
	"log"
	"reflect"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/pkg/event"
)

type LoadState uint8

const (
	NotFound LoadState = iota
	Loading
	Loaded
	Failed
)

type loaderFn = func(reader io.Reader) (any, error)

type loadRequest struct {
	path       string
	bundle     any
	bundleType reflect.Type
	handle     Handle
	assetType  reflect.Type
	loaderFn   loaderFn
}

func update(w *ecs.World) {
	as, _ := ecs.GetResource[Server](w)

	for !as.pending.IsEmpty() {
		req, _ := as.pending.Dequeue()

		// immediately add the asset to the context
		ecs.AddResourceByType(w, req.bundle, req.bundleType)

		// load asset asynchronously
		go func(r loadRequest) {
			defer func() {
				if rec := recover(); rec != nil {
					as.completed <- loadResult{r, nil,
						fmt.Errorf("loader panic: %v", rec)}
				}
			}()

			file, err := as.fsys.Open(r.path)
			if err != nil {
				as.completed <- loadResult{r, nil, err}
				return
			}

			content, err := r.loaderFn(file)
			as.completed <- loadResult{r, content, err}
		}(req)
	}

	// handle completed load results
	for {
		select {
		case res := <-as.completed:
			if res.err != nil {
				as.loadStates[res.req.path] = Failed
				as.registered.RemoveByKey(res.req.path)
				log.Printf("asset with path %s could not be loaded", res.req.path)

				return
			}

			// add the asset to the store
			storeAssetByType(as, res.req.assetType, res.req.handle, res.content)

			// notify the state of the new asset
			event.AddEvent(w, EventAssetAdded{
				Handle: res.req.handle,
			})

		default:
			return
		}
	}
}
