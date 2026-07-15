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
		Methods: exec.Methods{
			Bool:  builtins.Methods.Bool,
			Int:   builtins.Methods.Int,
			Float: builtins.Methods.Float,
			Str:   patchedStrMethods(),
			Dict:  builtins.Methods.Dict,
			List:  builtins.Methods.List,
		},
	}
}

// gonjaStrMethodNames is gonja v2.8.0's full builtin string-method surface. The
// patched method set below copies each of these from the builtins so overriding
// the broken ones does not silently drop the rest (the builtin map itself is
// unexported and cannot be merged into).
var gonjaStrMethodNames = []string{
	"capitalize", "capwords", "casefold", "center", "count", "encode",
	"endswith", "errors", "expandtabs", "find", "format", "format_map",
	"isalnum", "isalpha", "isascii", "isdecimal", "isdigit", "islower",
	"isnumeric", "isprintable", "isspace", "istitle", "isupper", "join",
	"ljust", "lower", "lstrip", "partition", "removeprefix", "removesuffix",
	"replace", "rfind", "rjust", "rpartition", "rsplit", "rstrip", "split",
	"splitlines", "startswith", "strip", "swapcase", "title", "upper", "zfill",
}

// patchedStrMethods returns the builtin string methods with strip/lstrip/rstrip
// replaced: gonja v2.8.0's implementations initialize the right-side cut index
// to 0 instead of the string length, so any input WITHOUT trailing cut
// characters strips to the empty string ("a.com".strip() == "") -- silently
// corrupting every chart that trims a split CSV param.
func patchedStrMethods() *exec.MethodSet[string] {
	patched := make(map[string]exec.Method[string], len(gonjaStrMethodNames))
	for _, name := range gonjaStrMethodNames {
		if m, ok := builtins.Methods.Str.Get(name); ok {
			patched[name] = m
		}
	}
	patched["strip"] = stripMethod(strings.Trim)
	patched["lstrip"] = stripMethod(strings.TrimLeft)
	patched["rstrip"] = stripMethod(strings.TrimRight)
	return exec.NewMethodSet(patched)
}

// stripMethod adapts a strings.Trim* function to Python's str.strip contract:
// no argument (or an empty one) trims whitespace, otherwise the argument is the
// set of characters to remove.
func stripMethod(trim func(s, cutset string) string) exec.Method[string] {
	return func(self string, _ *exec.Value, arguments *exec.VarArgs) (any, error) {
		var cutset string
		if err := arguments.Take(
			exec.PositionalArgument("cutset", exec.AsValue(""), exec.StringArgument(&cutset)),
		); err != nil {
			return nil, exec.ErrInvalidCall(err)
		}
		if cutset == "" {
			cutset = " \t\n\r\v\f"
		}
		return trim(self, cutset), nil
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
