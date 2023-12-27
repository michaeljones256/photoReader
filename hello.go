package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

var m sync.Mutex

func recurseDirectories(wg *sync.WaitGroup, photoPaths *[]string, directory string) {
	defer wg.Done()
	items, _ := os.ReadDir(directory)
	for _, item := range items {
		if item.IsDir() {
			fmt.Println(item.Name())
			wg.Add(1)
			go recurseDirectories(wg, photoPaths, directory+"/"+item.Name())
		} else {
			// handle file there
			m.Lock()
			*photoPaths = append(*photoPaths, directory+"/"+item.Name())
			m.Unlock()
		}
	}
}

func main() {
	args := os.Args
	fmt.Printf("Type of Args = %T\n", args)
	fmt.Println(args[0], args[1])
	path := args[1]

	stringSlice := make([]string, 0) // make creates slices

	wg := new(sync.WaitGroup)
	wg.Add(1)
	start := time.Now()
	recurseDirectories(wg, &stringSlice, path)
	wg.Wait()
	fmt.Println(time.Since(start))
	fmt.Println(len(stringSlice))
}
