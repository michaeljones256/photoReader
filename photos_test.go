package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func BenchmarkSingleThreadRecursive(b *testing.B) {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file: ", err)
	}
	path := os.Getenv("PATH_FILES")

	for i := 0; i < b.N; i++ {
		findPhotosSingleThread(path)
	}
}

func BenchmarkMultiThreadRecursive(b *testing.B) {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file: ", err)
	}
	path := os.Getenv("PATH_FILES")

	for i := 0; i < b.N; i++ {
		findPhotosMulitThread(path)
	}
}

func BenchmarkMulitThreadReceivingChannel(b *testing.B) {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file: ", err)
	}
	path := os.Getenv("PATH_FILES")

	for i := 0; i < b.N; i++ {
		findPhotosMulitThreadChannel(path)
	}
}

func BenchmarkRecursiveWorkerPool(b *testing.B) {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file: ", err)
	}
	path := os.Getenv("PATH_FILES")

	for i := 0; i < b.N; i++ {
		findPhotosWorkerPool(path, 3)
	}
}
