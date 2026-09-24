package hier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"strings"
	"text/template"

	"github.com/itchyny/gojq"
	"gopkg.in/yaml.v3"
)

func Lookup(spec *HierSpec, varname string) (*LookupResult, error) {
	// our lookup result
	res := &LookupResult{
		Merge: spec.Merge,
	}

	// debug print the spec
	if spec.Debug {
		slog.Debug("Spec", "spec", spec)
	}

	// go through the spec in order to find value
	// we take the first match

	results := make([]any, 0)
	found := false

	for _, h := range spec.Hierarchy {
		if spec.Debug {
			slog.Debug("process hierarchy", "name", h.Name)
		}

		// eval the path template to a real path on disk
		path, err := evalSpecPath(h.Path, spec)
		if err != nil {
			return nil, err
		}
		if !fileExistsReg(spec.FS, path) {
			if spec.Debug {
				slog.Debug("Path not found", "path", path)
			}
			continue
		}

		if spec.Debug {
			slog.Debug("Got path from template", "path", path, "template", h.Path)
		}
		fbuf, err := fs.ReadFile(spec.FS, path)
		if err != nil {
			slog.Debug("ReadFile error", "f", path, "err", err)
			continue
		}
		// load data
		var (
			data map[string]any

			// strval string
			// found  bool
		)
		err = yaml.Unmarshal(fbuf, &data)
		if err != nil {
			slog.Debug("Unmarshal error", "f", path, "err", err)
			continue
		}

		// try a lookup via jq

		val, err := execJqQuery(spec, varname, data)
		if err != nil {
			slog.Debug("execJqQuery", "err", err)
			continue
		}

		if spec.Merge {
			found = true
			results = append(results, val)
		} else {
			res.Name = h.Name
			res.Path = h.Path
			res.Value = val
			found = true
			break
		}

	}
	if found {

		if spec.Merge {
			b, _ := json.Marshal(results)
			var parsed any
			if err := json.Unmarshal([]byte(b), &parsed); err != nil {
				panic(err)
			}

			merged, err := MergeNestedArraysAny(parsed)
			if err != nil {
				slog.Debug("MergeNestedArraysAny error", "err", err)
			}

			res.Path = ""
			res.Value = merged

		}

		return res, nil
	}

	return res, LookupNotFound
}

func execJqQuery(spec *HierSpec, varname string, data any) (any, error) {
	var (
		retval any
		found  bool
		expr   string
	)
	if strings.HasPrefix(varname, ".") {
		expr = varname
	} else {
		expr = "." + varname
	}
	query, err := gojq.Parse(expr)
	if err != nil {
		slog.Debug("jq parse error", "expr", expr, "err", err)
		return nil, err
	}

	results := make([]any, 0)

	var iter gojq.Iter
	iter = query.Run(data)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			slog.Warn("jq iter", "e", err)
			break
		}
		type_id := fmt.Sprintf("%T", v)
		if spec.Debug {
			slog.Debug("iter var type", "type_id", type_id)
		}

		// todo: Corce to string
		// it could be any type here that yaml detects

		switch v.(type) {
		case bool:
			bval := v.(bool)
			retval = bval
			//strval = fmt.Sprintf("%v", bval)
		case nil:
			//strval = ""
			retval = nil
		case string:
			//strval, ok = v.(string)
			retval, ok = v.(string)
			if !ok {
				continue
			}
		default:
			// slog.Warn("Unsupported type", "type", type_id)
			//buf, _ := json.Marshal(v)
			//strval = string(buf)
			retval = v
		}

		results = append(results, retval)
		found = true
		if spec.Merge == false {
			break
		}
	}
	if !found {
		return nil, nil
	}

	if spec.Merge {
		retval = results
		return retval, nil
	}
	return retval, nil
}

type LookupResult struct {
	Merge bool
	Name  string
	Path  string
	Value any
}

func evalSpecPath(path_template string, spec *HierSpec) (string, error) {
	t, err := template.New("t").Parse(path_template)
	if err != nil {
		return "", err
	}
	var out bytes.Buffer
	err = t.Execute(&out, spec)
	if err != nil {
		return "", err
	}
	path := out.String()
	return path, nil
}

// regular file exists; if directory returns false
func fileExistsReg(fs fs.FS, path string) bool {
	f, err := fs.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return false
	}

	if st.IsDir() {
		return false
	}
	return true
}
