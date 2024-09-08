package autoreload

import (
	"context"
	"fmt"
	"time"

	"github.com/ckshitij/cache/pkg/cache"
)

type DataFunc[T any] func() (map[string]cache.CacheElement[T], error)

// Auto reload cache used for automatically
// reload a cache and maintain it in memory to
// to reduce request call.
type AutoReload[T any] struct {
	lookupTable       [2]map[string]cache.CacheElement[T]
	createdAt         time.Time
	refreshDuration   time.Duration
	keys              []string
	activeInd         int8
	isReloadInitiated bool
	name              string
	opts              cache.Options
}

func NewAutoReload[T any](
	ctx context.Context,
	cacheName string,
	loadFunc DataFunc[T],
	opts ...cache.Option,
) (AutoReload[T], error) {
	autoReloadCache := AutoReload[T]{
		lookupTable: [2]map[string]cache.CacheElement[T]{
			make(map[string]cache.CacheElement[T]),
			make(map[string]cache.CacheElement[T]),
		},
		createdAt:       time.Now().UTC(),
		refreshDuration: time.Minute,
		keys:            []string{},
		activeInd:       0,
		name:            cacheName,
	}

	err := autoReloadCache.opts.Apply(opts...)
	if err != nil {
		return AutoReload[T]{}, fmt.Errorf("options apply: %w", err)
	}

	if autoReloadCache.opts.AutoReloadInterval > 0 {
		autoReloadCache.refreshDuration *= autoReloadCache.opts.AutoReloadInterval
	}

	go autoReloadCache.triggerAutoReload(ctx, loadFunc)
	return autoReloadCache, nil
}

// Since there already all the data corresponding
// to the load function into the active table which
// always use for the read purpose.
func (ds *AutoReload[T]) Get(key string) (cache.CacheElement[T], bool) {
	data, ok := ds.lookupTable[ds.activeInd][key]
	return data, ok
}

func (ds *AutoReload[T]) GetRefreshDuration() time.Duration {
	return ds.refreshDuration
}

func (ds *AutoReload[T]) GetAllKeys() []string {
	return ds.keys
}

// triggerAutoReload function, will write the data in
// into non-active table to avoid any issue in case of
// data conflicting then once data got written
// the table got toggled using XOR operation
func (ds *AutoReload[T]) triggerAutoReload(ctx context.Context, loadFunc DataFunc[T]) {
	if ds.isReloadInitiated {
		return
	}
	ds.isReloadInitiated = true

	ticker := time.NewTicker(ds.refreshDuration)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			for retry := range 3 {
				newLookUp, err := loadFunc()
				if err == nil {
					ds.lookupTable[ds.activeInd^1] = newLookUp
					ds.activeInd ^= 1
					fmt.Printf("%s cache got reloaded\n", ds.name)
					break
				} else {
					time.Sleep(time.Duration(10*retry) * time.Millisecond)
				}
			}
		case <-ctx.Done():
			fmt.Println("auto reload function terminated due to context done")
			return
		}
	}
}
