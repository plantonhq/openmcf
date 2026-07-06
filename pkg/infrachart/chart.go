// Package infrachart loads and validates infra-charts entirely offline: it renders each
// chart's Jinja templates with the defaults declared in values.yaml and checks every
// rendered manifest against the local kind registry -- protovalidate on the spec plus
// resolution of every valueFrom reference. This is the same render -> parse -> validate
// pipeline the control plane runs when a chart is published, rebuilt from the pieces the
// CLI already owns, so chart breakage is caught at authoring time instead of after a
// release ships.
//
// The control plane remains the authoritative validator (its renderer is the engine of
// record and it validates against the platform's kind registry); this package is the
// shippability gate that runs with no backend at all.
package infrachart

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"
	goyaml "gopkg.in/yaml.v3"
)

const (
	chartFileName  = "Chart.yaml"
	valuesFileName = "values.yaml"
	templatesDir   = "templates"
)

// ParamType mirrors the chart param type vocabulary: string (the default), number, bool,
// and list.
type ParamType string

const (
	ParamTypeString ParamType = "string"
	ParamTypeNumber ParamType = "number"
	ParamTypeBool   ParamType = "bool"
	ParamTypeList   ParamType = "list"
)

// Param is one entry of the values.yaml params list -- the chart's public input surface.
type Param struct {
	Name        string    `yaml:"name"`
	Description string    `yaml:"description"`
	Type        ParamType `yaml:"type"`
	Value       any       `yaml:"value"`
}

// Chart is a loaded infra-chart directory: its params and its template files, keyed by
// path relative to templates/ (the key order is made deterministic by Templates()).
type Chart struct {
	// Dir is the chart directory the chart was loaded from.
	Dir string
	// Params are the declared inputs with their default values.
	Params []Param

	templates map[string]string
}

// IsChartDir reports whether dir holds an infra-chart (a Chart.yaml is the marker,
// matching how chart bundles are discovered everywhere else).
func IsChartDir(dir string) bool {
	info, err := os.Stat(filepath.Join(dir, chartFileName))
	return err == nil && !info.IsDir()
}

// Load reads a chart directory: params from values.yaml and every template under
// templates/ (recursively -- charts nest templates in subdirectories).
func Load(dir string) (*Chart, error) {
	if !IsChartDir(dir) {
		return nil, errors.Errorf("%s is not a chart directory (no %s)", dir, chartFileName)
	}

	valuesBytes, err := os.ReadFile(filepath.Join(dir, valuesFileName))
	if err != nil {
		return nil, errors.Wrapf(err, "reading %s", valuesFileName)
	}
	var values struct {
		Params []Param `yaml:"params"`
	}
	if err := goyaml.Unmarshal(valuesBytes, &values); err != nil {
		return nil, errors.Wrapf(err, "parsing %s", valuesFileName)
	}

	templates := map[string]string{}
	templatesRoot := filepath.Join(dir, templatesDir)
	err = filepath.WalkDir(templatesRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		if !strings.HasSuffix(d.Name(), ".yaml") && !strings.HasSuffix(d.Name(), ".yml") {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return errors.Wrapf(err, "reading template %s", path)
		}
		rel, err := filepath.Rel(templatesRoot, path)
		if err != nil {
			return err
		}
		templates[rel] = string(content)
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "walking %s", templatesRoot)
	}
	if len(templates) == 0 {
		return nil, errors.Errorf("chart %s has no templates under %s", dir, templatesDir)
	}

	return &Chart{Dir: dir, Params: values.Params, templates: templates}, nil
}

// Templates returns the template file names in deterministic order.
func (c *Chart) Templates() []string {
	names := make([]string, 0, len(c.templates))
	for name := range c.templates {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Template returns one template's source by its Templates() name.
func (c *Chart) Template(name string) string {
	return c.templates[name]
}

// BoolParams returns the names of the chart's bool-typed params, in declaration order.
// These are the branch toggles validation flips to reach conditional manifests.
func (c *Chart) BoolParams() []string {
	var names []string
	for _, p := range c.Params {
		if p.Type == ParamTypeBool {
			names = append(names, p.Name)
		}
	}
	return names
}
