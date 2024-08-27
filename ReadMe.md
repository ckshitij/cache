## In-Memory Cache

In-memory cache library in go-lang which can be used to store the key-value pair in memory and also provide the auto cleanup functionality using time-to-live duration.

### Installation

- Use the below command to install the library
    ```sh
    go get github.com/ckshitij/cache
    ```

### Code consumption example

- Below is the main file example, to consume the key-value datastore library.

    ```go
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

    ```
