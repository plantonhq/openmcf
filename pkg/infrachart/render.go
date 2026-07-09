package infrachart

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/builtins"
	"github.com/nikolalohinski/gonja/v2/config"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/loaders"
	"github.com/pkg/errors"
)

// The platform renders charts with a SANDBOXED Jinjava engine: the tags and
// filters below are disabled server-side, so a template using them would
// render offline but fail (or worse, misbehave) on the platform. They are
// rejected here by source scan, keeping the offline gate at least as strict
// as the platform.
var (
	bannedTagRe    = regexp.MustCompile(`\{%-?\s*(set|include|import|from|macro|call|do)\b`)
	bannedFilterRe = regexp.MustCompile(`\|\s*(attr|map|sort|shuffle)\b`)
)

func init() {
	// gonja logs through logrus by default; a validation library must not
	// write to the host process's output.
	gonja.SetLoggerOutput(io.Discard)
}

// chartEnvironment carries the builtin Jinja surface plus the two custom
// filters the platform registers (bool) or that charts conventionally rely on
// (b64decode). Built once; never mutates gonja's shared globals.
var chartEnvironment = &exec.Environment{
	Context:           gonja.DefaultContext,
	Filters:           exec.NewFilterSet(map[string]exec.FilterFunction{"bool": filterBool, "b64decode": filterB64Decode}).Update(builtins.Filters),
	Tests:             builtins.Tests,
	ControlStructures: builtins.ControlStructures,
	Methods:           builtins.Methods,
}

// filterBool mirrors the platform's `bool` filter semantics (Java's
// Boolean.parseBoolean for strings): booleans pass through, the string "true"
// (case-insensitive) is true and every other string is false, numbers are
// true when non-zero. Charts use it so boolean tests behave identically
// whether a value arrived as a typed bool param or a plain string.
func filterBool(_ *exec.Evaluator, in *exec.Value, _ *exec.VarArgs) *exec.Value {
	switch {
	case in.IsBool():
		return in
	case in.IsString():
		return exec.AsValue(strings.EqualFold(strings.TrimSpace(in.String()), "true"))
	case in.IsNumber():
		return exec.AsValue(in.Float() != 0)
	default:
		return exec.AsValue(false)
	}
}

// filterB64Decode decodes a base64 string — the platform registers the same
// filter for charts that embed encoded content.
func filterB64Decode(_ *exec.Evaluator, in *exec.Value, _ *exec.VarArgs) *exec.Value {
	decoded, err := base64.StdEncoding.DecodeString(in.String())
	if err != nil {
		return exec.AsValue(errors.Wrap(err, "b64decode: input is not valid base64"))
	}
	return exec.AsValue(string(decoded))
}

// scanBannedConstructs reports the sandbox-disabled tags/filters a template
// uses, so authors learn about a server-side render failure offline.
func scanBannedConstructs(templateSource string) []string {
	var findings []string
	for _, m := range bannedTagRe.FindAllStringSubmatch(templateSource, -1) {
		findings = append(findings, fmt.Sprintf("tag {%% %s %%} is disabled by the platform's sandboxed renderer", m[1]))
	}
	for _, m := range bannedFilterRe.FindAllStringSubmatch(templateSource, -1) {
		findings = append(findings, fmt.Sprintf("filter | %s is disabled by the platform's sandboxed renderer", m[1]))
	}
	return findings
}

// renderTemplate renders one template file with the given context.
func renderTemplate(name, source string, ctx map[string]any) (string, error) {
	baseLoader, err := loaders.NewFileSystemLoader("")
	if err != nil {
		return "", errors.Wrap(err, "failed to construct template loader")
	}
	shifted, err := loaders.NewShiftedLoader(name, bytes.NewReader([]byte(source)), baseLoader)
	if err != nil {
		return "", errors.Wrap(err, "failed to construct template loader")
	}
	tpl, err := exec.NewTemplate(name, config.New(), shifted, chartEnvironment)
	if err != nil {
		return "", err
	}
	return tpl.ExecuteToString(exec.NewContext(ctx))
}
