package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const sample = `---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.20.1
  name: widgets.example.com
spec:
  group: example.com
  names:
    kind: Widget
    plural: widgets
  scope: Namespaced
  versions:
  - name: v1
    schema:
      openAPIV3Schema:
        description: |-
          A description with a
          ---
          separator-looking line inside it.
        type: object
    served: true
    storage: true
`

func TestRenderKeepsControllerGenBytesVerbatim(t *testing.T) {
	out, err := Render([]byte(sample))
	if err != nil {
		t.Fatal(err)
	}
	rendered := string(out)

	if !strings.HasPrefix(rendered, "{{- /*") {
		t.Fatalf("missing template header:\n%s", rendered)
	}
	if !strings.Contains(rendered, "{{- if .Values.crds.enabled }}\n---\napiVersion") {
		t.Fatalf("the enabled guard must wrap the document:\n%s", rendered)
	}
	if !strings.HasSuffix(rendered, "{{- end }}\n") {
		t.Fatalf("the enabled guard must close the document:\n%s", rendered)
	}

	want := "  annotations:\n    {{- if .Values.crds.keep }}\n    helm.sh/resource-policy: keep\n    {{- end }}\n    controller-gen.kubebuilder.io/version: v0.20.1\n"
	if !strings.Contains(rendered, want) {
		t.Fatalf("keep block must sit inside metadata.annotations:\n%s", rendered)
	}

	// Strip what this program adds; controller-gen's bytes must be left intact.
	body := rendered[strings.Index(rendered, "---\napiVersion"):]
	body = strings.TrimSuffix(body, "{{- end }}\n")
	body = strings.Replace(body, "    {{- if .Values.crds.keep }}\n    helm.sh/resource-policy: keep\n    {{- end }}\n", "", 1)
	if body != sample {
		t.Fatalf("controller-gen bytes changed.\nwant:\n%s\ngot:\n%s", sample, body)
	}
}

func TestRenderRefusesUnexpectedShapes(t *testing.T) {
	cases := map[string]string{
		"not a CRD":            "---\napiVersion: v1\nkind: ConfigMap\nmetadata:\n  annotations:\n  name: x\n",
		"no annotations line":  "---\napiVersion: apiextensions.k8s.io/v1\nkind: CustomResourceDefinition\nmetadata:\n  name: x\n",
		"two annotations maps": strings.Replace(sample, "  scope: Namespaced\n", "  scope: Namespaced\n  annotations:\n", 1),
	}
	for name, src := range cases {
		if _, err := Render([]byte(src)); err == nil {
			t.Errorf("%s: expected an error", name)
		}
	}
}

func TestRunReplacesTheOutputDirectory(t *testing.T) {
	in := t.TempDir()
	out := t.TempDir()
	if err := os.WriteFile(filepath.Join(in, "a.yaml"), []byte(sample), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(out, "removed.yaml"), []byte("stale"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := run(in, out); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(out, "removed.yaml")); !os.IsNotExist(err) {
		t.Fatal("a definition that left the Go types must leave the chart")
	}
	if _, err := os.Stat(filepath.Join(out, "a.yaml")); err != nil {
		t.Fatal("expected a.yaml to be written")
	}
}

func TestRunRequiresInput(t *testing.T) {
	if err := run(t.TempDir(), t.TempDir()); err == nil {
		t.Fatal("an empty input directory must be an error, not an empty chart")
	}
}
