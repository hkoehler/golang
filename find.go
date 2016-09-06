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

// Result of directoy traversal
type Result struct {
	du  int64
	num uint64
}

// impose global limit on open files
var openFileSem chan int

// report result back to parent via comp channel
func find(opts Options, comp chan *Result) {
	res := new(Result)
	count := 0
	subChan := make(chan *Result)

	// acquire open file
	<-openFileSem
	dir, err := os.Open(opts.path)
	if err != nil {
		log.Fatal(err)
	}

	for {
		files, err := dir.Readdir(512)
		switch err {
		case io.EOF:
			goto done
		case nil:
			break
		default:
			log.Fatal(err)
		}

		// spawn gopher for each sub directory
		for _, file := range files {
			filePath := filepath.Join(opts.path, file.Name())
			if opts.print {
				fmt.Println(filePath)
			}
			res.du += file.Size() // XXX size not disk usage
			res.num++
			if file.Mode().IsDir() {
				newOpts := opts
				newOpts.path = filePath
				count++
				go find(newOpts, subChan)
			}
		}
	}

done:
	// wait for gophers to complete
	// update total result with results from gophers
	for count > 0 {
		subRes := <-subChan
		res.du += subRes.du
		res.num += subRes.num
		count--
	}

	// signal parent gopher
	comp <- res

	// release open file
	openFileSem <- 1
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

	// limit number of open files to 1000
	openFileSem = make(chan int, 1000)
	for i := 0; i < 1000; i++ {
		openFileSem <- 1
	}

	// go gopher go!
	wait := make(chan *Result)
	go find(opts, wait)
	res := <-wait

	if opts.num {
		fmt.Printf("%d files\n", res.num)
	}
	if opts.du {
		fmt.Printf("%d bytes\n", res.du)
	}
}
