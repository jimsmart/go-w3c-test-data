//go:build ignore

package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// const rdfTestsRepoZipURL = "https://github.com/w3c/rdf-tests/archive/refs/heads/main.zip"

const rdfTestsRepoZipURL = "https://github.com/jimsmart/rdf-tests/archive/refs/heads/main.zip"
const rdfTestsZipFile = "testdata/rdf-tests.zip"

const n3RepoZipURL = "https://github.com/w3c/N3/archive/refs/heads/master.zip"
const n3ZipFile = "testdata/n3.zip"

const outputFolder = "testdata/github.com/w3c/"

func main() {
	var err error

	err = fetch(rdfTestsRepoZipURL, rdfTestsZipFile)
	if err != nil {
		exit(err)
	}
	err = unzip(rdfTestsZipFile, outputFolder, "rdf-tests-main/", "rdf-tests/")
	if err != nil {
		exit(err)
	}
	err = os.Remove(rdfTestsZipFile)
	if err != nil {
		exit(err)
	}

	err = fetch(n3RepoZipURL, n3ZipFile)
	if err != nil {
		exit(err)
	}
	err = unzip(n3ZipFile, outputFolder, "N3-master/", "N3/")
	if err != nil {
		exit(err)
	}
	err = os.Remove(n3ZipFile)
	if err != nil {
		exit(err)
	}
}

func exit(err error) {
	fmt.Println("error:", err)
	os.Exit(1)
}

func fetch(url, dst string) error {
	fmt.Printf("Fetching %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()

	n, err := io.Copy(f, resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf("Saved to %s (%d bytes)\n", dst, n)
	return nil
}

func unzip(src, dst, prefixSrc, prefixDst string) error {
	fmt.Printf("Unzip %s to %s\n", src, dst)
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	n := 0
	for _, zf := range r.File {

		name := zf.Name
		if !strings.HasPrefix(name, prefixSrc) {
			fmt.Println("prefix does not match", name)
			continue
		}
		name = name[len(prefixSrc):]
		name = filepath.Join(prefixDst, name)

		if !include(name) {
			continue
		}

		path := filepath.Join(dst, name)

		// Prevent directory traversal.
		if !strings.HasPrefix(path, dst) {
			return fmt.Errorf("invalid output path: %s", path)
		}

		if zf.FileInfo().IsDir() {
			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				return err
			}
			continue
		}

		r, err := zf.Open()
		if err != nil {
			return err
		}
		defer r.Close()

		// Ensure the folder exists where
		// this file is to be written -
		// because we are quite heavy-handed
		// with our 'include' filter.
		err = os.MkdirAll(filepath.Dir(path), os.ModePerm)
		if err != nil {
			return err
		}

		w, err := os.Create(path)
		if err != nil {
			return err
		}
		defer w.Close()

		_, err = io.Copy(w, r)
		if err != nil {
			return err
		}

		n++
	}
	fmt.Printf("Unzipped %d files\n", n)
	return nil
}

func include(path string) bool {
	switch {
	case strings.HasSuffix(path, ".zip"),
		strings.HasSuffix(path, ".tar.gz"):
		return false
	case strings.Contains(path, "LICENSE"),
		strings.Contains(path, "README"),
		strings.HasPrefix(path, "rdf-tests/rdf"),
		strings.HasPrefix(path, "N3/tests/manifest.ttl"),
		strings.HasPrefix(path, "N3/tests/TurtleTests"),
		strings.HasPrefix(path, "N3/tests/N3Tests"):
		return true
	default:
		return false
	}
}
