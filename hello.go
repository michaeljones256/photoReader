package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

var m sync.Mutex

func recurseDirectoriesMultiThread(wg *sync.WaitGroup, photoPaths *[]string, directory string) {
	defer wg.Done()
	items, _ := os.ReadDir(directory)
	for _, item := range items {
		if item.IsDir() {
			wg.Add(1)
			go recurseDirectoriesMultiThread(wg, photoPaths, directory+"/"+item.Name())
		} else {
			// handle file
			m.Lock()
			*photoPaths = append(*photoPaths, directory+"/"+item.Name())
			m.Unlock()
		}
	}
}
func recurseDirectories(photoPaths *[]string, directory string) {
	items, _ := os.ReadDir(directory)
	for _, item := range items {
		if item.IsDir() {
			recurseDirectories(photoPaths, directory+"/"+item.Name())
		} else {
			// handle file
			*photoPaths = append(*photoPaths, directory+"/"+item.Name())
		}
	}
}

func main() {
	args := os.Args
	fmt.Println(args[1])
	path := args[1]

	stringSlice1 := make([]string, 0) // slice
	start := time.Now()
	recurseDirectories(&stringSlice1, path)
	fmt.Println("Single thread time " + time.Since(start).String())
	fmt.Println("Number of files " + strconv.Itoa(len(stringSlice1)))

	stringSlice2 := make([]string, 0) // slice
	wg := new(sync.WaitGroup)
	wg.Add(1)
	start = time.Now()
	recurseDirectoriesMultiThread(wg, &stringSlice2, path)
	wg.Wait()
	fmt.Println("Mulit thread time " + time.Since(start).String())
	fmt.Println("Number of files " + strconv.Itoa(len(stringSlice2)))
}
