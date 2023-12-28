package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

var m sync.Mutex
var m1 sync.Mutex

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
func recurseDirectoriesMultiThreadChannel(wg *sync.WaitGroup, result chan PhotoFile, directory string) {
	defer wg.Done()
	items, _ := os.ReadDir(directory)
	for _, item := range items {
		if item.IsDir() {
			wg.Add(1)
			go recurseDirectoriesMultiThreadChannel(wg, result, directory+"/"+item.Name())
		} else {
			// handle file
			result <- PhotoFile{Path: directory + "/" + item.Name()}
		}
	}
}
func reciever(photoPaths *[]string, results chan PhotoFile) {
	for i := range results {
		// m1.Lock()
		*photoPaths = append(*photoPaths, i.Path)
		// m1.Unlock()
	}
}

func main() {
	// Find .env file
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file: ", err)
	}
	path := os.Getenv("PATH_FILES")
	fmt.Println(path)

	//-----------mutli thread no channel-----------
	stringSlice2 := make([]string, 0) // slice
	wg := new(sync.WaitGroup)
	wg.Add(1)
	start := time.Now()
	recurseDirectoriesMultiThread(wg, &stringSlice2, path)
	wg.Wait()
	multiTheadTime := time.Since(start).String()
	fmt.Println("Mulit thread time " + multiTheadTime)
	fmt.Println("Number of files " + strconv.Itoa(len(stringSlice2)))

	//-----------single thread recursive-----------
	stringSlice1 := make([]string, 0) // slice
	start = time.Now()
	recurseDirectories(&stringSlice1, path)
	singleThreadTime := time.Since(start).String()
	fmt.Println("Single thread time " + singleThreadTime)
	fmt.Println("Number of files " + strconv.Itoa(len(stringSlice1)))

	//-----------mulit thread recieving channel-----------
	result := make(chan PhotoFile)
	wg.Add(1)
	start = time.Now()
	recurseDirectoriesMultiThreadChannel(wg, result, path)

	stringSlice := make([]string, 0)
	go func() {
		wg.Wait()
		close(result)
		fmt.Println("Multi thread channel time: " + time.Since(start).String())
		fmt.Println("Number of files " + strconv.Itoa(len(stringSlice)))
	}()
	reciever(&stringSlice, result)

}
