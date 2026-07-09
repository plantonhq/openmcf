package infrachart

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The end-to-end tests exercise the full gate against the permanent _test
// kind (TestCloudResourceGeneric), whose annotated_ref field carries the full
// FK annotation pair (default_kind = TestCloudResourceGeneric,
// default_kind_field_path = "status.outputs.id") — a hermetic fixture that
// never moves with production resource shapes.

const testChartYaml = `apiVersion: infra-hub.planton.ai/v1
kind: InfraChart
metadata:
  name: Test Chart
spec:
  selector:
    kind: platform
  description: hermetic test chart
`

const testValuesYaml = `params:
  - name: base_name
    description: base resource name
    value: node
`

func writeChart(t *testing.T, templates map[string]string) string {
	t.Helper()
	return writeChartWithValues(t, testValuesYaml, templates)
}

func writeChartWithValues(t *testing.T, valuesYaml string, templates map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	mustWrite(t, filepath.Join(dir, "Chart.yaml"), testChartYaml)
	mustWrite(t, filepath.Join(dir, "values.yaml"), valuesYaml)
	for name, content := range templates {
		mustWrite(t, filepath.Join(dir, "templates", name), content)
	}
	return dir
}

func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func mustValidate(t *testing.T, dir string, opts Options) *Report {
	t.Helper()
	report, err := Validate(dir, opts)
	if err != nil {
		t.Fatalf("Validate failed: %v", err)
	}
	return report
}

func issueMessages(report *Report, severity Severity) []string {
	var out []string
	for _, issue := range report.Issues {
		if issue.Severity == severity {
			out = append(out, issue.Message)
		}
	}
	return out
}

func requireError(t *testing.T, report *Report, substring string) {
	t.Helper()
	for _, msg := range issueMessages(report, SeverityError) {
		if strings.Contains(msg, substring) {
			return
		}
	}
	t.Fatalf("expected an error containing %q, got errors: %v", substring, issueMessages(report, SeverityError))
}

const validNode = `---
apiVersion: _test.planton.dev/v1
kind: TestCloudResourceGeneric
metadata:
  name: "{{ env }}-{{ values.base_name }}-b"
spec:
  requiredRef:
    value: literal
`

func TestValidateHappyPathWithInChartReference(t *testing.T) {
	dir := writeChart(t, map[string]string{
		"b.yaml": validNode,
		"a.yaml": `---
apiVersion: _test.planton.dev/v1
kind: TestCloudResourceGeneric
metadata:
  name: "{{ env }}-{{ values.base_name }}-a"
spec:
  requiredRef:
    value: literal
  annotatedRef:
    valueFrom:
      name: "{{ env }}-{{ values.base_name }}-b"
`,
	})
	report := mustValidate(t, dir, Options{})
	if report.HasErrors() {
		t.Fatalf("expected no errors, got: %v", issueMessages(report, SeverityError))
	}
	if len(report.Docs) != 2 {
		t.Fatalf("expected 2 docs, got %d", len(report.Docs))
	}
	if warnings := issueMessages(report, SeverityWarning); len(warnings) != 0 {
		t.Fatalf("expected no warnings (reference resolves in-chart), got: %v", warnings)
	}
	// The annotated default fieldPath was applied — the reference resolved
	// without an explicit fieldPath in the template.
	if report.Docs[0].Name != "dev-node-a" && report.Docs[1].Name != "dev-node-a" {
		t.Fatalf("rendered names unexpected: %v / %v", report.Docs[0].Name, report.Docs[1].Name)
	}
}

func TestValidateFkOverrideIsAnError(t *testing.T) {
	dir := writeChart(t, map[string]string{
		"b.yaml": validNode,
		"a.yaml": `---
apiVersion: _test.planton.dev/v1
kind: TestCloudResourceGeneric
metadata:
  name: "a"
spec:
  requiredRef:
    value: literal
  annotatedRef:
    valueFrom:
      name: "dev-node-b"
      fieldPath: status.outputs.name
`,
	})
	report := mustValidate(t, dir, Options{})
	requireError(t, report, "overrides the annotated composition key")
}

func TestValidateDanglingReferenceIsAWarning(t *testing.T) {
	dir := writeChart(t, map[string]string{
		"a.yaml": `---
apiVersion: _test.planton.dev/v1
kind: TestCloudResourceGeneric
metadata:
  name: "a"
spec:
  requiredRef:
    value: literal
  annotatedRef:
    valueFrom:
      name: "ghost"
`,
	})
	report := mustValidate(t, dir, Options{})
	if report.HasErrors() {
		t.Fatalf("dangling reference must not be an error: %v", issueMessages(report, SeverityError))
	}
	warnings := issueMessages(report, SeverityWarning)
	if len(warnings) != 1 || !strings.Contains(warnings[0], "does not define") {
		t.Fatalf("expected one dangling-reference warning, got: %v", warnings)
	}
}

func TestValidateDependencyCycleIsAnError(t *testing.T) {
	node := func(name, refName string) string {
		return `---
apiVersion: _test.planton.dev/v1
kind: TestCloudResourceGeneric
metadata:
  name: "` + name + `"
spec:
  requiredRef:
    value: literal
  annotatedRef:
    valueFrom:
      name: "` + refName + `"
`
	}
	dir := writeChart(t, map[string]string{
		"a.yaml": node("a", "b"),
		"b.yaml": node("b", "a"),
	})
	report := mustValidate(t, dir, Options{})
	requireError(t, report, "dependency cycle")
}

func TestValidateUnresolvableFieldPathIsAnError(t *testing.T) {
	dir := writeChart(t, map[string]string{
		"a.yaml": `---
apiVersion: _test.planton.dev/v1
kind: TestCloudResourceGeneric
metadata:
  name: "a"
spec:
  requiredRef:
    value: literal
  optionalRef:
    valueFrom:
      kind: TestCloudResourceGeneric
      name: "a"
      fieldPath: status.outputs.no_such_output
`,
	})
	report := mustValidate(t, dir, Options{})
	requireError(t, report, "does not resolve")
}

func TestValidateReferenceWithoutKindIsAnError(t *testing.T) {
	dir := writeChart(t, map[string]string{
		"a.yaml": `---
apiVersion: _test.planton.dev/v1
kind: TestCloudResourceGeneric
metadata:
  name: "a"
spec:
  requiredRef:
    value: literal
  optionalRef:
    valueFrom:
      name: "something"
`,
	})
	report := mustValidate(t, dir, Options{})
	requireError(t, report, "does not name a kind")
}

func TestValidateSpecViolationIsAnError(t *testing.T) {
	dir := writeChart(t, map[string]string{
		"a.yaml": `---
apiVersion: _test.planton.dev/v1
kind: TestCloudResourceGeneric
metadata:
  name: "a"
spec:
  stringField: hello
`,
	})
	report := mustValidate(t, dir, Options{})
	requireError(t, report, "spec validation failed")
}

func TestValidateMissingMetadataNameIsAnError(t *testing.T) {
	dir := writeChart(t, map[string]string{
		"a.yaml": `---
apiVersion: _test.planton.dev/v1
kind: TestCloudResourceGeneric
metadata:
  id: not-a-name
spec:
  requiredRef:
    value: literal
`,
	})
	report := mustValidate(t, dir, Options{})
	requireError(t, report, "metadata.name is required")
}

func TestValidateUnknownFieldIsAnError(t *testing.T) {
	dir := writeChart(t, map[string]string{
		"a.yaml": `---
apiVersion: _test.planton.dev/v1
kind: TestCloudResourceGeneric
metadata:
  name: "a"
spec:
  requiredRef:
    value: literal
  fieldThatDoesNotExist: true
`,
	})
	report := mustValidate(t, dir, Options{})
	if !report.HasErrors() {
		t.Fatal("unknown spec field must fail strict loading")
	}
}

func TestValidateToggleRendersBothArms(t *testing.T) {
	values := testValuesYaml + `  - name: extraEnabled
    description: deploy the extra node
    type: bool
    value: false
`
	templates := map[string]string{
		"a.yaml": validNode,
		"extra.yaml": `{% if values.extraEnabled | bool %}
---
apiVersion: _test.planton.dev/v1
kind: TestCloudResourceGeneric
metadata:
  name: "extra"
spec:
  requiredRef:
    value: literal
{% endif %}
`,
	}
	dir := writeChartWithValues(t, values, templates)

	// Default arm: toggle off — the conditional file renders to nothing.
	report := mustValidate(t, dir, Options{})
	if report.HasErrors() || len(report.Docs) != 1 {
		t.Fatalf("toggle-off: errors=%v docs=%d", issueMessages(report, SeverityError), len(report.Docs))
	}

	// Non-default arm via --set.
	report = mustValidate(t, dir, Options{Set: map[string]string{"extraEnabled": "true"}})
	if report.HasErrors() || len(report.Docs) != 2 {
		t.Fatalf("toggle-on: errors=%v docs=%d", issueMessages(report, SeverityError), len(report.Docs))
	}
}

func TestValidateBannedTagIsAnError(t *testing.T) {
	dir := writeChart(t, map[string]string{
		"a.yaml": "{% set x = 1 %}\n" + validNode,
	})
	report := mustValidate(t, dir, Options{})
	requireError(t, report, "disabled by the platform's sandboxed renderer")
}

func TestValidateDuplicateIdentityIsAnError(t *testing.T) {
	dir := writeChart(t, map[string]string{
		"a.yaml": validNode,
		"b.yaml": validNode,
	})
	report := mustValidate(t, dir, Options{})
	requireError(t, report, "duplicate resource identity")
}

func TestValidateValuelessParamWarnsWhenUsed(t *testing.T) {
	values := `params:
  - name: base_name
    description: has no default
`
	dir := writeChartWithValues(t, values, map[string]string{
		"a.yaml": `---
apiVersion: _test.planton.dev/v1
kind: TestCloudResourceGeneric
metadata:
  name: "{{ values.base_name }}"
spec:
  requiredRef:
    value: literal
`,
	})
	report := mustValidate(t, dir, Options{})
	found := false
	for _, msg := range issueMessages(report, SeverityWarning) {
		if strings.Contains(msg, "has no default value") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected a no-default warning, got: %+v", report.Issues)
	}
}
