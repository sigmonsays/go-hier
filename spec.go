package hier

import (
	"io/fs"
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

func NewHierSpec() *HierSpec {
	spec := &HierSpec{
		Inputs:    make(map[string]any, 0),
		Hierarchy: make([]*HierarchySpec, 0),
	}
	return spec
}

func (me *HierSpec) WithFS(hfs fs.FS) *HierSpec {
	me.FS = hfs
	return me
}

func (me *HierSpec) WithHierSpec(spec *HierarchySpec) *HierSpec {
	me.Hierarchy = append(me.Hierarchy, spec)
	return me

}
func (me *HierSpec) WithHierSpecPath(name, path string) *HierSpec {
	s := &HierarchySpec{
		Name: name,
		Path: path,
	}
	return me.WithHierSpec(s)
}

func (me *HierSpec) WithInput(key string, value any) *HierSpec {
	me.Inputs[key] = value
	return me
}

func LoadHierSpec(cfgfile string) (*HierSpec, error) {
	// load the deployment spec
	buf, err := os.ReadFile(cfgfile)
	if err != nil {
		return nil, err
	}

	spec := NewHierSpec()
	err = yaml.Unmarshal(buf, &spec)
	if err != nil {
		return nil, err
	}

	// Set some defaults
	if len(spec.Hierarchy) == 0 {
		spec.Hierarchy = []*HierarchySpec{
			{
				Name: "by fqdn",
				Path: "data/fqdn/{{ .Inputs.fqdn }}.yaml",
			},
			{
				Name: "by domain",
				Path: "data/domain/{{ .Inputs.domain }}.yaml",
			},
			{
				Name: "by pop",
				Path: "data/pop/{{ .Inputs.pop }}.yaml",
			},
			{
				Name: "by env",
				Path: "data/env/{{ .Inputs.env }}.yaml",
			},
			{
				Name: "default",
				Path: "data/default.yaml",
			},
		}
		slog.Debug("setting default hierarchy lookup", "hier", spec.Hierarchy)
	}

	return spec, nil
}
