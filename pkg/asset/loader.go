package asset

import (
	"fmt"
	"log"

	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/event"
)

func Update(ctx *core.Context) {
	as, ok := core.GetResource[*Server](ctx)
	if !ok {
		return
	}

	eb := core.MustResource[*event.Bus](ctx)

	for !as.pending.IsEmpty() {
		req, _ := as.pending.Dequeue()

		// immediately add the asset to the context
		core.AddResourceByType(ctx, req.Bundle, req.BundleType)

		// load asset asynchronously
		go func(r LoadRequest) {
			defer func() {
				if rec := recover(); rec != nil {
					as.completed <- loadResult{r, nil,
						fmt.Errorf("loader panic: %v", rec)}
				}
			}()

			file, err := as.fsys.Open(r.Path)
			if err != nil {
				as.completed <- loadResult{r, nil, err}
				return
			}

			content, err := r.LoaderFn(file, r.Path)
			as.completed <- loadResult{r, content, err}
		}(req)
	}

	// handle completed load results
	for {
		select {
		case res := <-as.completed:
			if res.err != nil {
				as.loadStates[res.req.Path] = Failed
				as.registered.RemoveByKey(res.req.Path)
				log.Printf("asset with path %s could not be loaded", res.req.Path)

				return
			}

			// add the asset to the store
			StoreAssetByType(as, res.req.AssetType, res.req.Handle, res.content)

			// notify the state of the new asset
			event.Publish(eb, AssetEvent{
				Handle: res.req.Handle,
				Kind:   Added,
			})

		default:
			return
		}
	}
}
