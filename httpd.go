/***********************************************************
 * Simplistic HTTP server for serving files and directories
 * Copyright (C) 2016, Heiko Koehler
 ***********************************************************/

package main

import (
    "fmt"
    "os"
    "io"
    "flag"
    "log"
    "net/http"
    "path/filepath"
)

var (
    // root path
    path string
)

type Page struct {
    resp http.ResponseWriter
    file *os.File
    filePath string
    urlPath string
    fileInfo os.FileInfo
}

func OpenPage(w http.ResponseWriter, urlPath string) (*Page, error) {
    var err error
    
    log.Printf("Open %s\n", urlPath)
    page := new(Page)
    page.resp = w
    page.filePath = filepath.Join(path, urlPath)
    page.urlPath = urlPath
    page.file, err = os.Open(page.filePath)
    if err != nil {
        return nil, err
    }
    page.fileInfo, err = page.file.Stat()
    if err != nil {
        return nil, err
    }
    return page, nil
}

func (page *Page) ShowDir() {
    dir := page.file
    w := page.resp

    fmt.Fprintf(w, "<html>")
    fmt.Fprintf(w, "<head>")
    fmt.Fprintf(w, "<title> Directory %s </title>", page.urlPath)
    fmt.Fprintf(w, "</head>")
    
    fmt.Fprintf(w, "<body>")
    fmt.Fprintf(w, "<h1> Directory %s </h1>", page.urlPath)
    for {
        files, err := dir.Readdir(64)
        if err == nil {
            for _, f := range files {
                href := filepath.Join(page.urlPath, f.Name())
                fmt.Fprintf(w, "<a href=\"%s\"> %s </a> <br>\n", href, f.Name())
            }
        } else if err == io.EOF {
            break
        } else {
            fmt.Fprintf(w, "Error: %s", err)
            break
        }
    }
    fmt.Fprintf(w, "</body>")
    fmt.Fprintf(w, "</html>")
}

func (page *Page) ShowFile() {
    io.Copy(page.resp, page.file)
}

func (page *Page) Show() {
    if page.fileInfo.IsDir() {
        page.ShowDir()
    } else {
        page.ShowFile()
    }
}

func (page *Page) Close() {
    page.file.Close()
}

func handler(w http.ResponseWriter, r *http.Request) {
    page, err := OpenPage(w, r.URL.Path)
    if err != nil {
        fmt.Fprintf(w, "Error opening %s: %s", r.URL.Path, err)
        return
    }
    defer page.Close()
    page.Show()
}

func main() {
    flag.StringVar(&path, "p", ".", "path to export")
    flag.Parse()

    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}
