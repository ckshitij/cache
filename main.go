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

type Person struct {
	Name string
	ID   int64
}

/*
Demo to how to consume the inmemory key value datastore
*/
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGHUP,
	) // we can add more sycalls.SIGQUIT etc.
	defer cancel()

	keysTTL := 1 * time.Second
	sweepTime := 3 * time.Second
	ds, err := cache.NewKeyValueCache[Person](ctx, keysTTL, cache.WithSweeping(sweepTime))
	if err != nil {
		fmt.Println("failed to initialize cache with error : ", err.Error())
	}

	var i int64 = 0
	for {
		key := fmt.Sprintf("person_id-%d", i+1)

		ds.Put(key, Person{
			Name: fmt.Sprintf("name_%d", i),
			ID:   i,
		})

		data, _ := ds.Get(key)
		fmt.Println(data.Value.Name)

		time.Sleep(100 * time.Millisecond)
		i++

		if i > 100 {
			break
		}
	}
}
