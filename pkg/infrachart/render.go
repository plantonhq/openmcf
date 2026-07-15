package infrachart

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"regexp"
	"strings"

	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/builtins"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/loaders"
	"github.com/pkg/errors"
)

// The control plane renders chart templates with a sandboxed Jinja engine: the tags and
// filters below are disabled there, so a template using them would behave differently (or
// fail) when published. Rejecting them here keeps offline validation strict-or-stricter
// than the engine of record -- a template this renderer accepts is always renderable by
// the control plane.
var (
	disabledTagPattern    = regexp.MustCompile(`\{%-?\s*(import|include|macro|call|set|do|from)\b`)
	disabledFilterPattern = regexp.MustCompile(`\|\s*(attr|map|sort|shuffle)\b`)
)

// renderEnvironment carries the Jinja builtins plus the chart-specific `bool` filter.
// Built once; gonja environments are safe to share across templates.
var renderEnvironment = newRenderEnvironment()

func newRenderEnvironment() *exec.Environment {
	filters := exec.NewFilterSet(map[string]exec.FilterFunction{
		"bool": filterBool,
	}).Update(builtins.Filters)
	return &exec.Environment{
		Context:           gonja.DefaultContext,
		Filters:           filters,
		Tests:             builtins.Tests,
		ControlStructures: builtins.ControlStructures,
		Methods:           builtins.Methods,
	}
}

// filterBool coerces a value to boolean the way chart authors use it on toggle params:
// booleans pass through, "true"/"false" (any case) and "1"/"0" strings parse, numbers are
// non-zero, and nil is false. Anything else is an error rather than a silent false, so a
// template never renders differently here than on the control plane's engine.
func filterBool(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if err := params.Take(); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}
	switch {
	case in.IsNil():
		return exec.AsValue(false)
	case in.IsBool():
		return exec.AsValue(in.Bool())
	case in.IsNumber():
		return exec.AsValue(in.Integer() != 0)
	case in.IsString():
		switch strings.ToLower(strings.TrimSpace(in.String())) {
		case "true", "1":
			return exec.AsValue(true)
		case "false", "0", "":
			return exec.AsValue(false)
		}
		return exec.AsValue(exec.ErrInvalidCall(fmt.Errorf("cannot coerce %q to bool", in.String())))
	}
	return exec.AsValue(exec.ErrInvalidCall(fmt.Errorf("cannot coerce %s to bool", in.String())))
}

// RenderTemplate renders one template source with the given param values. The rendering
// context mirrors the control plane's: every param is addressable both bare and through
// `values.*`, and org/env are always present (injected into both scopes, exactly like the
// server-side overwriter, so templates reading `values.env` or bare `env` behave the same).
func RenderTemplate(name, source string, values map[string]any, org, env string) (string, error) {
	if m := disabledTagPattern.FindStringSubmatch(source); m != nil {
		return "", errors.Errorf("template %s uses the '%s' tag, which the chart renderer disables", name, m[1])
	}
	if m := disabledFilterPattern.FindStringSubmatch(source); m != nil {
		return "", errors.Errorf("template %s uses the '%s' filter, which the chart renderer disables", name, m[1])
	}

	scoped := make(map[string]any, len(values)+2)
	for k, v := range values {
		scoped[k] = v
	}
	scoped["org"] = org
	scoped["env"] = env

	ctx := make(map[string]any, len(scoped)+1)
	for k, v := range scoped {
		ctx[k] = v
	}
	ctx["values"] = scoped

	template, err := templateFromString(source)
	if err != nil {
		return "", errors.Wrapf(err, "parsing template %s", name)
	}

	var out bytes.Buffer
	if err := template.Execute(&out, exec.NewContext(ctx)); err != nil {
		return "", errors.Wrapf(err, "rendering template %s", name)
	}
	return out.String(), nil
}

// templateFromString parses a template against the chart render environment. It mirrors
// gonja.FromString, which is hard-wired to the default environment.
func templateFromString(source string) (*exec.Template, error) {
	rootID := fmt.Sprintf("root-%x", sha256.Sum256([]byte(source)))
	fsLoader, err := loaders.NewFileSystemLoader("")
	if err != nil {
		return nil, err
	}
	shiftedLoader, err := loaders.NewShiftedLoader(rootID, strings.NewReader(source), fsLoader)
	if err != nil {
		return nil, err
	}
	return exec.NewTemplate(rootID, gonja.DefaultConfig, shiftedLoader, renderEnvironment)
}

// ParamValues converts the chart's declared params into a render context value map,
// applying the overrides on top of declared defaults.
func ParamValues(params []Param, overrides map[string]any) map[string]any {
	values := make(map[string]any, len(params))
	for _, p := range params {
		values[p.Name] = p.Value
	}
	for k, v := range overrides {
		values[k] = v
	}
	return values
}
