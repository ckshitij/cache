package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ckshitij/cache/pkg/cache"
)

type Person struct {
	name string
	id   int64
}

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

	keysTTL := 1 * time.Second
	sweepTime := 3 * time.Second
	ds, err := cache.NewKeyValueCache[Person](ctx, keysTTL, cache.WithSweeping(sweepTime))
	if err != nil {
		log.Fatal(err)
	}

	var i int64 = 0
	for {
		key := fmt.Sprintf("person_id-%d", i+1)

		ds.Put(key, Person{
			name: fmt.Sprintf("name_%s", i),
			id:   i,
		})

		p, _ := ds.Get(key)
		fmt.Println(p.name)

		time.Sleep(100 * time.Millisecond)
		i++

		if i > 100 {
			break
		}
	}
}
