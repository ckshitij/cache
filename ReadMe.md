## Key-Value datastore library in go-lang

In-memory datastore library in Go which can be used to store the key-value pair in memory and also provide the auto cleanup functionality using time-to-live duration.

### Installation

- Use the below command to install the library
    ```sh
    go get github.com/ckshitij/data-store
    ```

### Code consumption example

- Below is the main file example, to consume the key-value datastore library.

    ```go
        package main

        import (
            "fmt"
            "os"
            "os/signal"
            "syscall"
            "time"

            "github.com/ckshitij/data-store/pkg/datastore"
        )

        func multiSignalHandler(signal os.Signal, done chan bool) {

            fmt.Println("Started the multiSignal Handler")
            switch signal {
            case syscall.SIGHUP:
                fmt.Println("Signal: syscall.SIGHUP ", signal.String())
                done <- true
                time.Sleep(1 * time.Second)
                close(done)
                os.Exit(0)
            case syscall.SIGINT:
                fmt.Println("Signal: syscall.SIGINT ", signal.String())
                done <- true
                time.Sleep(1 * time.Second)
                close(done)
                os.Exit(0)
            case syscall.SIGTERM:
                fmt.Println("Signal: syscall.SIGTERM ", signal.String())
                done <- true
                time.Sleep(1 * time.Second)
                close(done)
                os.Exit(0)
            default:
                fmt.Println("Unhandled/unknown signal")
            }
        }

        /*
        Demo to how to consume the inmemory key value datastore
        */
        func main() {

            sigchnl := make(chan os.Signal, 1)
            signal.Notify(
                sigchnl,
                syscall.SIGHUP,
                syscall.SIGINT,
                syscall.SIGTERM,
            ) //we can add more sycalls.SIGQUIT etc.

            done := make(chan bool)
            go func() {
                for {
                    s := <-sigchnl
                    multiSignalHandler(s, done)
                }
            }()

            ds := inmemds.NewKeyValueDataStore(1 * time.Second)

            go ds.AutoCleanUp(3*time.Second, done)

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