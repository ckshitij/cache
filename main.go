package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ckshitij/cache/pkg/cache"
)

/*
Demo to how to consume the inmemory key value datastore
*/
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGINT,
		os.Kill,
	) // we can add more sycalls.SIGQUIT etc.
	defer cancel()

	ds := cache.NewKeyValueCache[string](1 * time.Second)
	go ds.Sweep(ctx, 3*time.Second)

	var i int64 = 0
	for {
		key := fmt.Sprintf(" key-%d", i+1)
		value := fmt.Sprintf(" value-%d", i+1)
		ds.Put(key, value)

		time.Sleep(100 * time.Millisecond)
		i++
	}
}
