package infrachart

import (
	"strings"
	"testing"
)

// These cases pin the render semantics for every Jinja construct the chart fleet uses.
// The control plane renders with a different engine implementation, so any drift between
// engines must surface here as a red test, never as a chart that validates offline but
// renders differently when published.
func TestRenderTemplateConformance(t *testing.T) {
	values := map[string]any{
		"name":         "demo",
		"enabled":      true,
		"disabled":     false,
		"enabled_str":  "true",
		"disabled_str": "false",
		"count":        "3",
		"csv":          "a,b,c",
		"items":        []any{"x", "y"},
		"padded":       "  padded  ",
		"long":         "abcdefghijklmnopqrstuvwxy",
	}

	cases := []struct {
		name     string
		template string
		want     string
	}{
		{"interpolation", `v: {{ values.name }}`, `v: demo`},
		{"bare org and env", `v: {{ org }}/{{ env }}`, `v: planton/dev`},
		{"values org and env", `v: {{ values.org }}/{{ values.env }}`, `v: planton/dev`},
		{"if bool true", `{% if values.enabled | bool %}yes{% else %}no{% endif %}`, `yes`},
		{"if bool false", `{% if values.disabled | bool %}yes{% else %}no{% endif %}`, `no`},
		{"if bool true string", `{% if values.enabled_str | bool %}yes{% else %}no{% endif %}`, `yes`},
		{"if bool false string", `{% if values.disabled_str | bool %}yes{% else %}no{% endif %}`, `no`},
		{"elif", `{% if values.disabled | bool %}a{% elif values.enabled | bool %}b{% else %}c{% endif %}`, `b`},
		{"for range int filter", `{% for i in range(values.count | int) %}{{ i }}{% endfor %}`, `012`},
		{"for split", `{% for s in values.csv.split(',') %}[{{ s }}]{% endfor %}`, `[a][b][c]`},
		{"for list param", `{% for s in values.items %}[{{ s }}]{% endfor %}`, `[x][y]`},
		{"replace", `{{ values.name | replace('e', 'a') }}`, `damo`},
		{"trim", `[{{ values.padded | trim }}]`, `[padded]`},
		{"upper", `{{ values.name | upper }}`, `DEMO`},
		{"lower", `{{ "DEMO" | lower }}`, `demo`},
		// Mirrors the fleet's only truncate call shape: truncate(n, true, ''). Jinja
		// leaves strings within the default 5-char leeway untouched, so only inputs
		// longer than n+5 are actually cut.
		{"truncate", `{{ values.long | truncate(15, true, '') }}`, `abcdefghijklmno`},
		{"truncate within leeway", `{{ "abcdefghij" | truncate(8, true, '') }}`, `abcdefghij`},
		{"length", `{{ values.csv.split(',') | length }}`, `3`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := RenderTemplate(tc.name, tc.template, values, "planton", "dev")
			if err != nil {
				t.Fatalf("render failed: %v", err)
			}
			if got != tc.want {
				t.Fatalf("rendered %q, want %q", got, tc.want)
			}
		})
	}
}

func TestRenderTemplateRejectsSandboxedConstructs(t *testing.T) {
	cases := []struct {
		name     string
		template string
		wantIn   string
	}{
		{"set tag", `{% set x = 1 %}{{ x }}`, "'set' tag"},
		{"include tag", `{% include "other.yaml" %}`, "'include' tag"},
		{"macro tag", `{% macro m() %}{% endmacro %}`, "'macro' tag"},
		{"map filter", `{{ values.items | map('upper') }}`, "'map' filter"},
		{"sort filter", `{{ values.items | sort }}`, "'sort' filter"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := RenderTemplate(tc.name, tc.template, map[string]any{"items": []any{"a"}}, "planton", "dev")
			if err == nil {
				t.Fatalf("expected rejection, got none")
			}
			if !strings.Contains(err.Error(), tc.wantIn) {
				t.Fatalf("error %q does not mention %q", err.Error(), tc.wantIn)
			}
		})
	}
}

func TestBoolFilterRejectsUnparseableStrings(t *testing.T) {
	_, err := RenderTemplate("bad-bool", `{% if values.v | bool %}x{% endif %}`,
		map[string]any{"v": "maybe"}, "planton", "dev")
	if err == nil {
		t.Fatal("expected an error coercing 'maybe' to bool")
	}
}
