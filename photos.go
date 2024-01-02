package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

var m sync.Mutex
var m1 sync.Mutex

type PhotoFile struct {
	Path string
}

type PhotoList struct {
	Photos *[]string
}

func regexMatchJPG(file string) bool {

	pattern := `.*\.(JPG|ARW)$`
	regexpPattern := regexp.MustCompile(pattern)
	if regexpPattern.MatchString(file) {
		return true
	} else {
		return false
	}
}
func findPhotosSingleThread(path string) {
	photoPaths := make([]string, 0)
	start := time.Now()
	recurseDirectories(&photoPaths, path)
	singleThreadTime := time.Since(start).String()
	fmt.Println("Single thread time " + singleThreadTime)
	fmt.Println("Number of files " + strconv.Itoa(len(photoPaths)))
}

func findPhotosMulitThread(path string) {
	photoPaths := make([]string, 0)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	start := time.Now()
	recurseDirectoriesMultiThread(wg, &photoPaths, path)
	wg.Wait()
	multiTheadTime := time.Since(start).String()
	fmt.Println("Mulit thread time " + multiTheadTime)
	fmt.Println("Number of files " + strconv.Itoa(len(photoPaths)))
}
func findPhotosMulitThreadChannel(path string) {
	results := make(chan PhotoFile)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	start := time.Now()
	recurseDirectoriesMultiThreadChannel(wg, results, path)

	photoPaths := make([]string, 0)
	go func() {
		wg.Wait()
		close(results)
		fmt.Println("Multi thread channel time: " + time.Since(start).String())
		fmt.Println("Number of files " + strconv.Itoa(len(photoPaths)))
	}()
	// receive results - waits until results closes
	for i := range results {
		photoPaths = append(photoPaths, i.Path)
	}
}
func findPhotosWorkerPool(path string, workers int) {
	directoryJobs := make(chan string)
	results := make(chan PhotoList)
	start := time.Now()
	go createWorkers(workers, directoryJobs, results)
	go createJobs(path, directoryJobs)
	photoPaths := make([]string, 0)
	// receive results - waits until results closes
	for i := range results {
		photoPaths = append(photoPaths, *i.Photos...)
	}
	fmt.Println("Multi worker time: " + time.Since(start).String())
	fmt.Println("Number of files " + strconv.Itoa(len(photoPaths)))
}
func recurseDirectories(photoPaths *[]string, directory string) {
	if regexMatchJPG(directory) {
		*photoPaths = append(*photoPaths, directory)
	} else {
		items, _ := os.ReadDir(directory)
		for _, item := range items {
			recurseDirectories(photoPaths, directory+"/"+item.Name())
		}
	}

}

func recurseDirectoriesMultiThread(wg *sync.WaitGroup, photoPaths *[]string, directory string) {
	defer wg.Done()
	items, _ := os.ReadDir(directory)
	for _, item := range items {
		if item.IsDir() {
			wg.Add(1)
			go recurseDirectoriesMultiThread(wg, photoPaths, directory+"/"+item.Name())
		} else {
			// handle file
			if regexMatchJPG(item.Name()) {
				m.Lock()
				*photoPaths = append(*photoPaths, directory+"/"+item.Name())
				m.Unlock()
			}
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
			if regexMatchJPG(item.Name()) {
				result <- PhotoFile{Path: directory + "/" + item.Name()}
			}
		}
	}
}
func recurseDirectoriesOld(photoPaths *[]string, directory string) {
	items, _ := os.ReadDir(directory)
	for _, item := range items {
		if item.IsDir() {
			recurseDirectoriesOld(photoPaths, directory+"/"+item.Name())
		} else {
			// handle file
			if regexMatchJPG(item.Name()) {
				*photoPaths = append(*photoPaths, directory+"/"+item.Name())
			}
		}
	}
}

func findFilesJob(directoryJobs <-chan string, results chan PhotoList, wg *sync.WaitGroup, workerId int) {
	defer wg.Done()
	for directory := range directoryJobs {
		photoPaths2 := make([]string, 0)
		// fmt.Println("worker: " + strconv.Itoa(workerId) + " doing job " + directory)
		recurseDirectories(&photoPaths2, directory)
		results <- PhotoList{Photos: &photoPaths2}
	}
}
func createJobs(intialPath string, directoryJobs chan string) {
	items, _ := os.ReadDir(intialPath)
	for _, item := range items {
		if item.IsDir() {
			directoryJobs <- intialPath + item.Name()
		} else {
			directoryJobs <- intialPath + "/" + item.Name()
		}
	}
	close(directoryJobs)
}

func createWorkers(numWorkers int, directoryJobs chan string, results chan PhotoList) {
	wg := new(sync.WaitGroup)
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go findFilesJob(directoryJobs, results, wg, w)
	}
	wg.Wait()
	close(results)

}

func main() {
	// Find .env file
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file: ", err)
	}
	path := os.Getenv("PHOTOS_PATH")

	findPhotosSingleThread(path)

	findPhotosMulitThread(path)

	findPhotosMulitThreadChannel(path)

	workers := 3
	findPhotosWorkerPool(path, workers)
}
