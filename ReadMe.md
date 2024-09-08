## Cache

Cache library in go-lang which can be used to store the key-value pair in memory and also provide the auto cleanup functionality using time-to-live duration.

### Types

Currently there are two type of cache are supported:
 - Datastore
   - Common in-memory datastore pretty much like to store key-values, with get and set functionality and provide the clean-up based on time-to-live parameters.
 - AutoReload
   - This cache is mostly used for read-heavy functionality in one go, there is not a set method since it mainly load all the data on startup then auto reload the data after given interval.
   - e.g, mostly this could be used if user need to load all the data at once from db to memory.

### Installation

- Use the below command to install the library
    ```sh
    go get github.com/ckshitij/cache
    ```
