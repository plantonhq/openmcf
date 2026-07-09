package infrachart

import (
	"strings"
	"testing"
)

func render(t *testing.T, source string, ctx map[string]any) string {
	t.Helper()
	out, err := renderTemplate("test.yaml", source, ctx)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	return out
}

func TestRenderVariableBindings(t *testing.T) {
	ctx, _, err := buildContext([]Param{
		{Name: "region", Value: "us-central1"},
	}, "acme", "dev", nil)
	if err != nil {
		t.Fatal(err)
	}

	out := render(t, "{{ values.region }}/{{ region }}/{{ org }}/{{ env }}", ctx)
	if out != "us-central1/us-central1/acme/dev" {
		t.Fatalf("unexpected render: %q", out)
	}
}

func TestRenderBoolFilter(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  string
	}{
		{"typed bool true", true, "yes"},
		{"typed bool false", false, "no"},
		{"string true", "true", "yes"},
		{"string True mixed case", "True", "yes"},
		{"string false", "false", "no"},
		{"string junk", "junk", "no"},
		{"number nonzero", float64(2), "yes"},
		{"number zero", float64(0), "no"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx, _, err := buildContext([]Param{{Name: "flag", Value: tc.value}}, "acme", "dev", nil)
			if err != nil {
				t.Fatal(err)
			}
			out := render(t, "{% if values.flag | bool %}yes{% else %}no{% endif %}", ctx)
			if out != tc.want {
				t.Fatalf("flag=%v: got %q, want %q", tc.value, out, tc.want)
			}
		})
	}
}

func TestRenderB64DecodeFilter(t *testing.T) {
	ctx, _, err := buildContext([]Param{{Name: "blob", Value: "aGVsbG8="}}, "acme", "dev", nil)
	if err != nil {
		t.Fatal(err)
	}
	if out := render(t, "{{ values.blob | b64decode }}", ctx); out != "hello" {
		t.Fatalf("b64decode: got %q", out)
	}
}

func TestRenderForLoopOverListParam(t *testing.T) {
	ctx, _, err := buildContext([]Param{
		{Name: "zones", Type: "list", Value: []any{"us-central1-a", "us-central1-b"}},
	}, "acme", "dev", nil)
	if err != nil {
		t.Fatal(err)
	}
	out := render(t, "{% for z in values.zones %}[{{ z }}]{% endfor %}", ctx)
	if out != "[us-central1-a][us-central1-b]" {
		t.Fatalf("for loop: got %q", out)
	}
}

func TestRenderPlaceholderForValuelessParam(t *testing.T) {
	ctx, placeholders, err := buildContext([]Param{{Name: "bucket_name"}}, "acme", "dev", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(placeholders) != 1 || placeholders[0] != "<bucket_name>" {
		t.Fatalf("placeholders: got %v", placeholders)
	}
	if out := render(t, "{{ values.bucket_name }}", ctx); out != "<bucket_name>" {
		t.Fatalf("placeholder render: got %q", out)
	}
}

func TestScanBannedConstructs(t *testing.T) {
	findings := scanBannedConstructs(`
{% set x = 1 %}
{% include "other.yaml" %}
{{ items | map('upper') }}
`)
	if len(findings) != 3 {
		t.Fatalf("expected 3 findings, got %d: %v", len(findings), findings)
	}
	joined := strings.Join(findings, "\n")
	for _, want := range []string{"set", "include", "map"} {
		if !strings.Contains(joined, want) {
			t.Fatalf("findings missing %q: %v", want, findings)
		}
	}
}

func TestScanAllowsPermittedConstructs(t *testing.T) {
	findings := scanBannedConstructs(`
{% if values.flag | bool %}
value: {{ values.name | default("x") | upper }}
{% endif %}
{% for z in values.zones %}{{ z | trim }}{% endfor %}
{# a comment #}
`)
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %v", findings)
	}
}
