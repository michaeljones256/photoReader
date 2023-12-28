package main

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/joho/godotenv"
)

func BenchmarkRecursive(b *testing.B) {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file: ", err)
	}
	path := os.Getenv("PATH_FILES")
	stringSlice := make([]string, 0)
	for i := 0; i < b.N; i++ {
		recurseDirectories(&stringSlice, path)
	}
}

func BenchmarkRecursiveMulti(b *testing.B) {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file: ", err)
	}
	path := os.Getenv("PATH_FILES")
	stringSlice := make([]string, 0)
	wg := new(sync.WaitGroup)
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		recurseDirectoriesMultiThread(wg, &stringSlice, path)
		wg.Wait()
	}
}

func BenchmarkRecursiveMultiChannel(b *testing.B) {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file: ", err)
	}
	path := os.Getenv("PATH_FILES")
	stringSlice := make([]string, 0)
	wg := new(sync.WaitGroup)
	for i := 0; i < b.N; i++ {
		result := make(chan PhotoFile)
		wg.Add(1)
		recurseDirectoriesMultiThreadChannel(wg, result, path)
		go func() {
			wg.Wait()
			close(result)
		}()
		reciever(&stringSlice, result)
	}
}
