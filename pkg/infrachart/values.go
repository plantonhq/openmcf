package infrachart

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
)

// Param is one parameter declared in values.yaml. The field set mirrors the
// platform's InfraChartParam contract; types beyond it are rejected so a
// chart that validates offline is guaranteed to parse on the platform.
type Param struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Type        string   `json:"type"`
	Value       any      `json:"value"`
	EnumValues  []string `json:"enumValues"`
	Hidden      bool     `json:"hidden"`
	FileContent bool     `json:"fileContent"`
}

// Param types, mirroring the platform's ParamType enum. An empty type means
// string.
const (
	paramTypeString     = "string"
	paramTypeNumber     = "number"
	paramTypeBool       = "bool"
	paramTypeList       = "list"
	paramTypeStringEnum = "string_enum"
)

type valuesYaml struct {
	Params []Param `json:"params"`
}

func loadValuesFile(path string) ([]Param, error) {
	valuesBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read %s", path)
	}
	var vy valuesYaml
	if err := yaml.Unmarshal(valuesBytes, &vy); err != nil {
		return nil, errors.Wrap(err, "failed to parse values.yaml")
	}
	seen := map[string]bool{}
	for i := range vy.Params {
		p := &vy.Params[i]
		if p.Name == "" {
			return nil, errors.Errorf("values.yaml params[%d] has no name", i)
		}
		if seen[p.Name] {
			return nil, errors.Errorf("values.yaml declares param %q more than once", p.Name)
		}
		seen[p.Name] = true
		if err := validateParam(p); err != nil {
			return nil, err
		}
	}
	return vy.Params, nil
}

// validateParam mirrors the platform's param-type validation: a declared type
// must match the supplied default value's YAML type. Params without values
// are legal (supplied at deployment time), but the offline gate will surface
// them as placeholders if templates depend on them.
func validateParam(p *Param) error {
	mismatch := func(want string) error {
		return errors.Errorf("param %q is declared %s but its value is not a %s (got %T)",
			p.Name, p.Type, want, p.Value)
	}
	switch p.Type {
	case "", paramTypeString:
		if p.Value == nil {
			return nil
		}
		if _, ok := p.Value.(string); !ok {
			return mismatch("string")
		}
	case paramTypeNumber:
		if p.Value == nil {
			return nil
		}
		switch p.Value.(type) {
		case float64, int, int64: // YAML numbers arrive as float64 via JSON; be liberal on ints
		default:
			return mismatch("number")
		}
	case paramTypeBool:
		if p.Value == nil {
			return nil
		}
		if _, ok := p.Value.(bool); !ok {
			return mismatch("bool")
		}
	case paramTypeList:
		if p.Value == nil {
			return nil
		}
		items, ok := p.Value.([]any)
		if !ok {
			return mismatch("list of strings")
		}
		for _, it := range items {
			if _, ok := it.(string); !ok {
				return mismatch("list of strings")
			}
		}
	case paramTypeStringEnum:
		if len(p.EnumValues) == 0 {
			return errors.Errorf("param %q is string_enum but declares no enumValues", p.Name)
		}
		if p.Value == nil {
			return nil
		}
		s, ok := p.Value.(string)
		if !ok {
			return mismatch("string")
		}
		for _, allowed := range p.EnumValues {
			if s == allowed {
				return nil
			}
		}
		return errors.Errorf("param %q value %q is not one of its enumValues %v", p.Name, s, p.EnumValues)
	default:
		return errors.Errorf("param %q has unknown type %q (valid: string, number, bool, list, string_enum)", p.Name, p.Type)
	}
	return nil
}

// buildContext converts the chart params (plus org/env and any --set
// overrides) into the template rendering context, mirroring the platform's
// renderer exactly: every param is bound both as a top-level variable and
// under "values", org/env are always bound (overwriting any params of those
// names), and a param without a value renders as the placeholder "<name>".
// The returned placeholders are those literal "<name>" strings, so the
// validator can flag templates that depend on value-less params.
func buildContext(params []Param, org, env string, setOverrides map[string]string) (map[string]any, []string, error) {
	values := map[string]any{}
	byName := map[string]*Param{}
	for i := range params {
		p := &params[i]
		byName[p.Name] = p
		values[p.Name] = templateValue(p)
	}

	for key, raw := range setOverrides {
		switch key {
		case "org":
			org = raw
			continue
		case "env":
			env = raw
			continue
		}
		p, ok := byName[key]
		if !ok {
			return nil, nil, errors.Errorf("--set %s: chart declares no param named %q", key, key)
		}
		coerced, err := coerceOverride(p, raw)
		if err != nil {
			return nil, nil, err
		}
		values[key] = coerced
	}

	// org/env always win — the platform overwrites them at render time.
	values["org"] = org
	values["env"] = env

	var placeholders []string
	for name, v := range values {
		if s, ok := v.(string); ok && s == "<"+name+">" {
			placeholders = append(placeholders, s)
		}
	}

	ctx := map[string]any{}
	for k, v := range values {
		ctx[k] = v
	}
	ctx["values"] = values
	ctx["org"] = org
	ctx["env"] = env
	return ctx, placeholders, nil
}

// templateValue converts a param's default into the value templates see. The
// "<name>" placeholder for value-less params mirrors the platform renderer,
// so a chart whose defaults don't fully render fails validation loudly (the
// placeholder never survives schema validation) instead of silently.
func templateValue(p *Param) any {
	if p.Value == nil {
		return "<" + p.Name + ">"
	}
	if items, ok := p.Value.([]any); ok {
		strs := make([]string, 0, len(items))
		for _, it := range items {
			strs = append(strs, fmt.Sprintf("%v", it))
		}
		return strs
	}
	return p.Value
}

// coerceOverride parses a --set string into the param's declared type.
func coerceOverride(p *Param, raw string) (any, error) {
	switch p.Type {
	case paramTypeBool:
		b, err := strconv.ParseBool(strings.TrimSpace(raw))
		if err != nil {
			return nil, errors.Errorf("--set %s: %q is not a bool", p.Name, raw)
		}
		return b, nil
	case paramTypeNumber:
		f, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
		if err != nil {
			return nil, errors.Errorf("--set %s: %q is not a number", p.Name, raw)
		}
		return f, nil
	case paramTypeList:
		if raw == "" {
			return []string{}, nil
		}
		return strings.Split(raw, ","), nil
	case paramTypeStringEnum:
		for _, allowed := range p.EnumValues {
			if raw == allowed {
				return raw, nil
			}
		}
		return nil, errors.Errorf("--set %s: %q is not one of %v", p.Name, raw, p.EnumValues)
	default:
		return raw, nil
	}
}

// ParseSetFlags parses repeated --set key=value flags.
func ParseSetFlags(flags []string) (map[string]string, error) {
	out := map[string]string{}
	for _, f := range flags {
		key, value, found := strings.Cut(f, "=")
		if !found || key == "" {
			return nil, errors.Errorf("--set %q is not in key=value form", f)
		}
		out[key] = value
	}
	return out, nil
}
