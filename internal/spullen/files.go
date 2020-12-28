package spullen

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Finder struct {
	Root string
}

func (f *Finder) FindDatabases() ([]string, error) {
	var path = "*.db"
	if f.Root != "" {
		path = fmt.Sprintf("%s/*.db", f.Root)
	}

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
