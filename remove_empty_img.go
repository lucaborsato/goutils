// package remomve_empty_img to remove images (png, jpg, jpeg) with size 0.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// fileInfo holds all you need to know about a file
type fileInfo struct {
	AbsPath string
	Name    string
	Size    int64
}

// filesInfo contains the files found in path
type filesInfo []fileInfo

func main() {

	// get path where to search for empty images with default to $USER/Temp/To_Clean_Folder/sdcard
	var inPath string
	flag.StringVar(&inPath, "i", filepath.Join(os.Getenv("USER"), "Temp", "To_Clean_Folder", "sdcard"), "path where to search for empty images")

	// get absolute full path from inPath
	fullPath, err := filepath.Abs(inPath)
	// in case of error print it and exit
	if err != nil {
		err = errors.Wrapf(err, "can't get absolute path of %v", inPath)
		log.Fatal(err)
	}
	// print the full path
	fmt.Println("Input path:", fullPath)

	// store all images found in inPath
	var images filesInfo
	// recursively look for images and store their info
	err = filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
		path, _ = filepath.Abs(path)
		fileExt := filepath.Ext(path)
		if fileExt == ".png" || fileExt == ".jpg" || fileExt == ".jpeg" {
			fi, err := os.Stat(path)
			if err != nil {
				err = errors.Wrapf(err, "can't get file info for %v", path)
				return err
			}

			// store image info
			images = append(images, fileInfo{
				AbsPath: path,
				Name:    filepath.Base(path),
				Size:    fi.Size(),
			})
		}
		return nil
	})
	// if an error occurred while walking the FS from inPath
	if err != nil {
		err = errors.Wrap(err, "error while images files info")
		log.Fatal(err)
	}

	// if no image is found, exit
	if len(images) == 0 {
		fmt.Println("No file found")
		os.Exit(0)
	}

	// manage images to be trashed
	discardPath := filepath.Join(fullPath, "discarded")
	err = os.MkdirAll(discardPath, os.ModePerm)
	if err != nil {
		err = errors.Wrapf(err, "can't create folder (%v) for discarded files", discardPath)
	}

	fmt.Println("Found files to be discarded:")
	for _, f := range images {
		if f.Size == 0 {
			newPath := filepath.Join(discardPath, f.Name)
			fmt.Printf("%s with size %d bytes\n", f.AbsPath, f.Size)
			fmt.Printf("Moving to %s\n", newPath)
			err = os.Rename(f.AbsPath, newPath)
			if err != nil {
				err = errors.Wrapf(err, "can't move %v to %v", f.Name, newPath)
			}
		}
	}
}
