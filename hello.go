package main

import (
	"fmt"
	"os"
	"time"
)

func recurseDirectories(photoPaths *[]string, directory string) {
	items, _ := os.ReadDir(directory)
	for _, item := range items {
		if item.IsDir() {
			fmt.Println(item.Name())
			recurseDirectories(photoPaths, directory+"/"+item.Name())
		} else {
			// handle file there
			*photoPaths = append(*photoPaths, directory+"/"+item.Name())
		}
	}
}

func main() {
	args := os.Args
	fmt.Printf("Type of Args = %T\n", args)
	fmt.Println(args[0], args[1])
	path := args[1]
	// items, _ := ioutil.ReadDir(path)
	// for _, item := range items {
	//     if item.IsDir() {
	//         subitems, _ := ioutil.ReadDir(path+"/"+item.Name())
	//         for _, subitem := range subitems {
	//             fmt.Println("here3")
	//             if !subitem.IsDir() {
	//                 // file
	//                 fmt.Println(item.Name() + "/" + subitem.Name())
	//             }else{
	//                 //dir
	// 			}
	//         }
	//     } else {
	//         // handle file there
	//         fmt.Println(item.Name())
	//     }
	// }

	stringSlice := make([]string, 0) // make creates slices
	start := time.Now()
	recurseDirectories(&stringSlice, path)
	fmt.Println(time.Since(start))

	fmt.Println(len(stringSlice))
}
