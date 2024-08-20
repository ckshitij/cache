package main

import (
	"fmt"
	"time"

	"github.com/ckshitij/data-store/pkg/datastore"
)

/*
Datastore consumtion code
*/
func main() {
	ds := datastore.NewKeyValueDataStore(time.Duration(2 * time.Second))
	ds.Put("1", 10)
	ds.Put("2", "shubham")
	ds.Put("3", 7845)
	ds.Put("4", []string{"kshitij", "chaurasiya"})

	time.Sleep(1 * time.Second)
	ds.Put("6", 784.05)
	ds.Put("7", []int{2, 3})
	fmt.Printf("Records present in Datastore \n%+v\n", ds.GetAllKeyValues())

	time.Sleep(1 * time.Second)
	fmt.Printf("Records present in Datastore \n%+v\n", ds.GetAllKeyValues())

}
