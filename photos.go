package main

import (
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

var m sync.Mutex
var m1 sync.Mutex

type PhotoFile struct {
	Path string
}

type PhotoList struct {
	Photos *[]string
}

func regexMatchJPG(directory string) bool {
	// Can give either a directory or a file
	directories := strings.Split(directory, "/")
	currentDirectory := directories[len(directories)-1]

	pattern := `^[a-zA-Z0-9]*.(JPG|ARW)$`
	regexpPattern := regexp.MustCompile(pattern)
	if regexpPattern.MatchString(currentDirectory) {
		return true
	} else {
		return false
	}
}

func findPhotosSingleThread(path string) FilesInfo {
	photoPaths := make([]string, 0)
	start := time.Now()
	recurseDirectories(&photoPaths, path)
	return FilesInfo{NumberFiles: len(photoPaths), Files: photoPaths, Time: time.Since(start).String()}
}

func findPhotosMulitThread(path string) FilesInfo {
	photoPaths := make([]string, 0)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	start := time.Now()
	recurseDirectoriesMultiThread(wg, &photoPaths, path)
	wg.Wait()
	return FilesInfo{NumberFiles: len(photoPaths), Files: photoPaths, Time: time.Since(start).String()}
}

func findPhotosMulitThreadChannel(path string) FilesInfo {
	results := make(chan PhotoFile)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	start := time.Now()
	recurseDirectoriesMultiThreadChannel(wg, results, path)

	photoPaths := make([]string, 0)
	end := ""
	go func(endTime *string) {
		wg.Wait()
		close(results)
		end = time.Since(start).String()
	}(&end)
	// receive results - waits until results closes
	for i := range results {
		photoPaths = append(photoPaths, i.Path)
	}
	return FilesInfo{NumberFiles: len(photoPaths), Files: photoPaths, Time: end}

}
func findPhotosWorkerPool(path string, workers int) FilesInfo {
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
	return FilesInfo{NumberFiles: len(photoPaths), Files: photoPaths, Time: time.Since(start).String()}
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
