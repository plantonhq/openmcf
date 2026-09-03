// chartcrds turns controller-gen's CustomResourceDefinition files into the
// Helm chart's CRD templates.
//
// The chart owns the lifecycle of the operator's definitions: they are ordinary
// release resources rendered behind `crds.enabled`, upgraded with the release,
// and annotated `helm.sh/resource-policy: keep` behind `crds.keep` so an
// uninstall leaves them (and every resource they define) in place. Helm's own
// `crds/` directory cannot do this -- it installs once and never upgrades -- so
// the definitions live under templates/ instead, and this program is how they
// get there without a hand-typed copy that could drift from the Go types.
//
// The transformation is textual on purpose. controller-gen's bytes are kept
// verbatim; only the template header, the `crds.enabled` guard, and the
// three-line keep block after `metadata.annotations:` are added. A rendered
// definition with the keep line removed is therefore byte-identical to
// controller-gen's file, which the chart's tests hold.
//
// Invoked by `make -C operator manifests` right after controller-gen runs.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	// crdPrologue is what every controller-gen CRD file begins with; anything
	// else is not a file this program understands.
	crdPrologue = "---\napiVersion: apiextensions.k8s.io/v1\nkind: CustomResourceDefinition\nmetadata:\n"

	// annotationsLine is controller-gen's metadata.annotations key (it always
	// emits one, carrying its own version). The keep block is inserted right
	// after it so the annotation lands inside the same map.
	annotationsLine = "  annotations:\n"

	templateHeader = `{{- /*
Generated from the operator's Go types by ` + "`make -C operator manifests`" + `; never edited by hand.
This chart owns this definition's lifecycle: ` + "`crds.enabled`" + ` installs and upgrades it with the
release, and ` + "`crds.keep`" + ` leaves it in place, together with every resource it defines, when
the release is uninstalled. A definition that already exists outside this release is handled by
templates/crds-preflight.yaml.
*/ -}}
{{- if .Values.crds.enabled }}
`

	keepBlock = "    {{- if .Values.crds.keep }}\n    helm.sh/resource-policy: keep\n    {{- end }}\n"

	templateFooter = "{{- end }}\n"
)

// Render wraps one controller-gen CRD document as a chart template.
func Render(src []byte) ([]byte, error) {
	if !bytes.HasPrefix(src, []byte(crdPrologue)) {
		return nil, fmt.Errorf("not a controller-gen CustomResourceDefinition file: expected it to begin with\n%s", crdPrologue)
	}
	if n := bytes.Count(src, []byte("\n"+annotationsLine)); n != 1 {
		return nil, fmt.Errorf("expected exactly one metadata.annotations line, found %d; controller-gen's output shape changed and this program must follow", n)
	}
	at := bytes.Index(src, []byte("\n"+annotationsLine)) + 1 + len(annotationsLine)

	var out bytes.Buffer
	out.WriteString(templateHeader)
	out.Write(src[:at])
	out.WriteString(keepBlock)
	out.Write(src[at:])
	if !bytes.HasSuffix(src, []byte("\n")) {
		out.WriteString("\n")
	}
	out.WriteString(templateFooter)
	return out.Bytes(), nil
}

func run(inDir, outDir string) error {
	sources, err := filepath.Glob(filepath.Join(inDir, "*.yaml"))
	if err != nil {
		return err
	}
	if len(sources) == 0 {
		return fmt.Errorf("no CustomResourceDefinition files in %s; run controller-gen first", inDir)
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	// A definition removed from the Go types must leave the chart too.
	stale, err := filepath.Glob(filepath.Join(outDir, "*.yaml"))
	if err != nil {
		return err
	}
	for _, f := range stale {
		if err := os.Remove(f); err != nil {
			return err
		}
	}
	for _, src := range sources {
		data, err := os.ReadFile(src)
		if err != nil {
			return err
		}
		rendered, err := Render(data)
		if err != nil {
			return fmt.Errorf("%s: %w", src, err)
		}
		dst := filepath.Join(outDir, filepath.Base(src))
		if err := os.WriteFile(dst, rendered, 0o644); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	inDir := flag.String("in", "config/crd/bases", "directory of controller-gen CustomResourceDefinition files")
	outDir := flag.String("out", "../helm/planton-operator/templates/crds", "chart templates directory to write")
	flag.Parse()
	if err := run(*inDir, *outDir); err != nil {
		fmt.Fprintln(os.Stderr, "chartcrds:", strings.TrimSpace(err.Error()))
		os.Exit(1)
	}
}
