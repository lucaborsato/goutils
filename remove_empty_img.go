// package remomve_empty_img to remove images (png, jpg, jpeg) with size 0.
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {

	var err error
// 	store all files
	var files []string
// 	store file size
	var sizeFiles []int64
// 	at the moment it uses hardcoded path...I am learning :D
	var inPath string = "/home/borsato/Temp/To_Clean_Folder/sdcard"
	var fullPath string
	var discardPath string
	var basePath string
	var newPath string

	// get absolute full path from inPath
	fullPath, err = filepath.Abs(inPath)
	// if it returns error it prints err and close the program
	if err != nil {
		fmt.Println("CHECKING ABSOLUTE PATH:")
		log.Fatal(err)
	}
	// print the full path
	fmt.Println("INPUT PATH:", fullPath)
	
	// walk through subfolders and list files
	err = filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
		path, _ = filepath.Abs(path)
		fileExt := filepath.Ext(path)
		if fileExt == ".png" || fileExt == ".jpg" || fileExt == ".jpeg" {
			files = append(files, path)
			fileInfo, errs := os.Stat(path)
			if errs != nil {log.Fatal(errs)}
			sizeFiles = append(sizeFiles, fileInfo.Size())
		}
		return nil
	})
	// error check as line 22
	if err != nil {
		fmt.Println("ERROR GETTING FILE PATH:")
		log.Fatal(err)
	}
	// print files found by Walk
	if len(files) == 0 {
		fmt.Println("IMG FILES NOT FOUND IN FOLDER!")
	} else {
		discardPath = filepath.Join(fullPath, "discarded")
		os.MkdirAll(discardPath, os.ModePerm)
		fmt.Println("FOUND FILES:")
		for i := 0; i < len(files); i++ {
			if sizeFiles[i] == 0 {
				basePath = filepath.Base(files[i])
				newPath = filepath.Join(discardPath, basePath)
				fmt.Printf("%s with size %d bytes\n", files[i], sizeFiles[i])
				fmt.Printf("MOVING TO %s\n", newPath)
				os.Rename(files[i], newPath)
// 				fmt.Printf("basePath: %s\n", basePath)
				
			}
		}
	}
}
