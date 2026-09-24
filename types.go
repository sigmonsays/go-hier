package hier

import (
	"fmt"
	"io/fs"
)

type HierSpec struct {
	Debug bool

	// file system we read from
	FS fs.FS

	// merge result
	Merge     bool             `yaml:"merge"`
	Inputs    map[string]any   `yaml:"inputs"`
	Hierarchy []*HierarchySpec `yaml:"hierarchy"`
}

type HierarchySpec struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

var (
	LookupNotFound = fmt.Errorf("Not found")
)
