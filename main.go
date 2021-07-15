package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

var backup_name string

func init() {
	flag.StringVar(&backup_name, "out", "backup.zip", "Archive name.")
}

func main() {
	flag.Parse()
	paths := flag.Args()

	// create .zip archive 
	backup_archive, err := os.Create(backup_name)
	if err != nil {
		log.Fatal(err)
	}
	defer backup_archive.Close()

	zip_writer := zip.NewWriter(backup_archive)
	defer zip_writer.Close()

	// add each paths provided to the archive
	for _,path := range paths {
		err := filepath.Walk(path, func(element string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", element, err)
				return err
			}
			if !info.IsDir() {
				element_src, err := os.Open(element)
				if err != nil {
					log.Fatal(err)
				}
				defer element_src.Close()
			
				element_focused, err := zip_writer.Create(element)
				if err != nil {
					log.Fatal(err)
				}
			
				_, err = io.Copy(element_focused, element_src)
				if err != nil {
					log.Fatal(err)
				}
			}
			return nil
		})
		if err != nil {
			fmt.Printf("walk error [%v]\n", err)
		}
	}
}
