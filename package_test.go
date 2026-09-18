package hier

import (
	"fmt"
	"testing"
	"testing/fstest"
)

var (
	default_yaml = []byte(`
exampled:
 pkg_version: 1.0.0
`)
	env_dev_yaml = []byte(`
exampled:
 pkg_version: 2.0.0
`)
)

func TestHier(t *testing.T) {

	testfs := fstest.MapFS{
		"env/dev.yaml": &fstest.MapFile{
			Data: env_dev_yaml,
		},
		"default.yaml": &fstest.MapFile{
			Data: default_yaml,
		},
	}

	hspec := NewHierSpec().
		// WithFS(os.DirFS("./data/")).
		WithFS(testfs).
		WithHierSpecPath("by-host", "host/{{ .Inputs.env }}/{{ .Inputs.host }}.yaml").
		WithHierSpecPath("by-env", "env/{{ .Inputs.env }}.yaml").
		WithHierSpecPath("default", "default.yaml")

	hspec.WithInput("host", "host1.example.net")
	hspec.WithInput("env", "dev")

	res, err := Lookup(hspec, "exampled.pkg_version")

	fmt.Printf("Lookup %+v\n", res)
	fmt.Printf("Lookup error %+v\n", err)
}
