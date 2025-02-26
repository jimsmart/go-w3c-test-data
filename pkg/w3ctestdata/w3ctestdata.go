// Package w3ctestdata bundles an embed.FS containing various W3C test suites.
package w3ctestdata

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

//go:generate go run fetch.go

// TODO(js) Perhaps we can be more selective?
//
// // go:embed testdata/github.com/w3c/rdf-tests/ns
// // go:embed testdata/github.com/w3c/rdf-tests/rdf/rdf11
// // go:embed testdata/github.com/w3c/rdf-tests/rdf/rdf12

//

// Files embeds the contents of:
//
//	testdata/github.com/w3c/rdf-tests
//	testdata/github.com/w3c/N3/tests
//
//go:embed testdata/github.com/w3c/rdf-tests
//go:embed testdata/github.com/w3c/N3/tests
var Files embed.FS

func Get(url string) (*bytes.Reader, error) {

	url, err := remapURL(url)
	if err != nil {
		return nil, err
	}

	r, err := Files.Open(url)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}

// Map these url to the filesystem.
// https://w3c.github.io/rdf-tests/rdf/rdf11/rdf-turtle/manifest.ttl
// https://w3c.github.io/N3/tests/manifest.ttl

const rdfTestsPrefix = "https://w3c.github.io/rdf-tests/"
const n3TestsPrefix = "https://w3c.github.io/N3/tests/"

func remapURL(url string) (string, error) {

	if strings.HasPrefix(url, rdfTestsPrefix) {
		url = url[len(rdfTestsPrefix):]
		url = filepath.Join("testdata/github.com/w3c/rdf-tests", url)
		return url, nil
	}

	if strings.HasPrefix(url, n3TestsPrefix) {
		url = url[len(n3TestsPrefix):]
		url = filepath.Join("testdata/github.com/w3c/N3/tests/", url)
		return url, nil
	}

	return "", fmt.Errorf("cannot map url to filesystem: %s", url)
}
