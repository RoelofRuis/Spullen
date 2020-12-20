package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Finder struct {
	root string
}

func (f *Finder) FindDatabases() ([]string, error) {
	var path = "*.db"
	if f.root != "" {
		path = fmt.Sprintf("%s/*.db", f.root)
	}

	println(path)

	files, err := filepath.Glob(path)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, f := range files {
		names = append(names, strings.TrimSuffix(f, ".db"))
	}

	return names, nil
}
