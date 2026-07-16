package infrachart

import (
	"strings"
	"testing"
)

func TestValidateParamTypeMismatch(t *testing.T) {
	if err := validateParam(&Param{Name: "flag", Type: "bool", Value: "true"}); err == nil {
		t.Fatal("bool param with string value must be rejected")
	}
	if err := validateParam(&Param{Name: "count", Type: "number", Value: "3"}); err == nil {
		t.Fatal("number param with string value must be rejected")
	}
	if err := validateParam(&Param{Name: "name", Type: "string", Value: true}); err == nil {
		t.Fatal("string param with bool value must be rejected")
	}
	if err := validateParam(&Param{Name: "zones", Type: "list", Value: []any{"a", 2}}); err == nil {
		t.Fatal("list param with non-string item must be rejected")
	}
}

func TestValidateParamStringEnum(t *testing.T) {
	if err := validateParam(&Param{Name: "tier", Type: "string_enum", Value: "basic"}); err == nil {
		t.Fatal("string_enum without enumValues must be rejected")
	}
	p := &Param{Name: "tier", Type: "string_enum", Value: "gold", EnumValues: []string{"basic", "standard"}}
	if err := validateParam(p); err == nil {
		t.Fatal("string_enum value outside enumValues must be rejected")
	}
	p.Value = "standard"
	if err := validateParam(p); err != nil {
		t.Fatalf("valid string_enum rejected: %v", err)
	}
}

func TestValidateParamUnknownType(t *testing.T) {
	err := validateParam(&Param{Name: "x", Type: "object", Value: "y"})
	if err == nil || !strings.Contains(err.Error(), "unknown type") {
		t.Fatalf("unknown param type must be rejected, got: %v", err)
	}
}

func TestBuildContextSetOverrides(t *testing.T) {
	params := []Param{
		{Name: "flag", Type: "bool", Value: false},
		{Name: "size", Type: "number", Value: float64(1)},
		{Name: "zones", Type: "list", Value: []any{"a"}},
		{Name: "name", Value: "original"},
	}
	ctx, _, err := buildContext(params, "acme", "dev", map[string]string{
		"flag":  "true",
		"size":  "2.5",
		"zones": "x,y",
		"name":  "overridden",
		"env":   "prod",
	})
	if err != nil {
		t.Fatal(err)
	}
	values := ctx["values"].(map[string]any)
	if values["flag"] != true {
		t.Fatalf("bool override: got %v", values["flag"])
	}
	if values["size"] != 2.5 {
		t.Fatalf("number override: got %v", values["size"])
	}
	if zones := values["zones"].([]string); len(zones) != 2 || zones[0] != "x" {
		t.Fatalf("list override: got %v", values["zones"])
	}
	if values["name"] != "overridden" {
		t.Fatalf("string override: got %v", values["name"])
	}
	if ctx["env"] != "prod" {
		t.Fatalf("env override: got %v", ctx["env"])
	}
}

func TestBuildContextRejectsUnknownSetKey(t *testing.T) {
	_, _, err := buildContext(nil, "acme", "dev", map[string]string{"ghost": "1"})
	if err == nil {
		t.Fatal("--set on an undeclared param must be rejected")
	}
}

func TestBuildContextOrgEnvAlwaysWin(t *testing.T) {
	// A chart declaring org/env params still renders with the injected values
	// — mirroring the platform, which overwrites them at render time.
	ctx, _, err := buildContext([]Param{{Name: "org", Value: "chart-org"}}, "real-org", "dev", nil)
	if err != nil {
		t.Fatal(err)
	}
	values := ctx["values"].(map[string]any)
	if values["org"] != "real-org" || ctx["org"] != "real-org" {
		t.Fatalf("org must be overwritten: got %v / %v", values["org"], ctx["org"])
	}
}

func TestParseSetFlags(t *testing.T) {
	out, err := ParseSetFlags([]string{"a=1", "b=x=y"})
	if err != nil {
		t.Fatal(err)
	}
	if out["a"] != "1" || out["b"] != "x=y" {
		t.Fatalf("unexpected parse: %v", out)
	}
	if _, err := ParseSetFlags([]string{"no-equals"}); err == nil {
		t.Fatal("flag without '=' must be rejected")
	}
}
