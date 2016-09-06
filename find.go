/***************************************************
 * Concurrent find implemention in Go
 * Copyright (C) 2016, Heiko Koehler
 ***************************************************/

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type Options struct {
	path  string // top-level directory
	print bool   // print all file paths
	du    bool   // disk usage
	num   bool   // number of files
}

func find(opts Options) error {
	dir, err := os.Open(opts.path)
	if err != nil {
		return err
	}
	for {
		files, err := dir.Readdir(512)
		switch err {
		case io.EOF:
			return nil
		case nil:
		default:
			log.Fatal(err)
			return err
		}
		for _, file := range files {
			filePath := filepath.Join(opts.path, file.Name())
			// descent into directory
			if file.Mode().IsDir() {
				newOpts := opts
				newOpts.path = filePath
				err = find(newOpts)
				if err != nil {
					return err
				}
			}
			if opts.print {
				fmt.Println(filePath)
			}
		}
	}
}

func main() {
	var opts Options

	flag.BoolVar(&opts.print, "p", false, "Print all file paths")
	flag.BoolVar(&opts.du, "d", false, "Calculate disk usage")
	flag.BoolVar(&opts.num, "n", false, "Count number of files")
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		var err error
		opts.path, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		opts.path = args[0]
	}
	err := find(opts)
	if err != nil {
		log.Fatal(err)
	}
}
