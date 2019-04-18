/*
go language port of pyBigList

MIT License

Copyright (c) 2019 Adam Koch

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	scriptPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	homeDir, _ := os.UserHomeDir()
	blockListGz := fmt.Sprintf("%s/biglist.p2p.gz", scriptPath)
	blockListURL := "http://john.bitsurge.net/public/biglist.p2p.gz"
	finalList := fmt.Sprintf("%s/Documents/biglist.p2p", homeDir)

	// removing old files
	if _, err := os.Stat(blockListGz); err == nil {
		os.Remove(blockListGz)
	}
	if _, err := os.Stat(finalList); err == nil {
		os.Remove(finalList)
	}

	// download file using net/http
	fmt.Printf("Downloading %s. . \n", blockListURL)
	if err := downloadList(blockListGz, blockListURL); err != nil {
		panic(err)
	}

	// open the gzip file
	fmt.Printf("Decompressing %s. . .\n", blockListGz)
	fo, err := os.Open(blockListGz)
	if err != nil {
		panic(err)
	}

	// initialize a new gzip reader
	zr, err := gzip.NewReader(fo)
	if err != nil {
		panic(err)
	}

	// create the decompressed file
	fw, err := os.Create(finalList)
	if err != nil {
		panic(err)
	}

	// copy the contents of the gzip stream into the new file
	if _, err := io.Copy(fw, zr); err != nil {
		panic(err)
	}

	// close the file handlers
	defer fo.Close()
	defer zr.Close()
	defer fw.Close()

	// clean up
	fmt.Printf("Cleaning up. . .\n")
	if _, err := os.Stat(blockListGz); err == nil {
		os.Remove(blockListGz)
	}

	fmt.Printf("Done!\n")
}

// downloads the file using net/http
func downloadList(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
