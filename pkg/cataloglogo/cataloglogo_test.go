package cataloglogo

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestEveryComponentLogoIsRenderableSVG is the gate: every registered
// kind's logo.svg must be a strict-XML SVG document. There is no baseline
// -- the tree holds zero violations and must stay there.
func TestEveryComponentLogoIsRenderableSVG(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate repo root")
	}
	repoRoot := filepath.Dir(filepath.Dir(filepath.Dir(thisFile)))
	vs, err := Check(repoRoot)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range vs {
		t.Errorf("%s: %s", v.Path, v.Detail)
	}
	if len(vs) > 0 {
		t.Fatalf("%d logo(s) will render as broken images in a browser -- fix the SVGs (this gate has no baseline)", len(vs))
	}
}

// TestValidateSVGDeliberateReds proves each rule fires on the defect
// classes that have actually shipped in this tree, and that the healthy
// shape passes.
func TestValidateSVGDeliberateReds(t *testing.T) {
	cases := []struct {
		name    string
		svg     string
		wantHit string // substring of the expected detail; "" means must pass
	}{
		{
			name:    "healthy minimal svg",
			svg:     `<?xml version="1.0"?><svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"><path d="M0 0h16v16H0z"/></svg>`,
			wantHit: "",
		},
		{
			// The kubernetes six: content and closing tag, no root.
			name:    "missing svg root",
			svg:     "<g><path d=\"M0 0\"/></g>\n",
			wantHit: "root element is <g>",
		},
		{
			// The awsbackupreportplan case: fatal to browsers, invisible to Go.
			name:    "duplicate attribute",
			svg:     `<svg xmlns="http://www.w3.org/2000/svg"><path fill-rule="evenodd" d="M0 0" fill-rule="nonzero"/></svg>`,
			wantHit: "duplicate attribute",
		},
		{
			// The configmap/cilium/karapace/kafkaconnector/opensearch case:
			// a raw control byte where an em dash was intended.
			name:    "control byte in comment",
			svg:     "<svg xmlns=\"http://www.w3.org/2000/svg\"><!-- broken \x14 dash --><path d=\"M0 0\"/></svg>",
			wantHit: "not well-formed XML",
		},
		{
			// The udproute case: "--" is illegal inside XML comments.
			name:    "double hyphen in comment",
			svg:     `<svg xmlns="http://www.w3.org/2000/svg"><!-- arrows -- connectionless --><path d="M0 0"/></svg>`,
			wantHit: "not well-formed XML",
		},
		{
			name:    "empty file",
			svg:     "",
			wantHit: "no <svg> root",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			detail := validateSVG([]byte(c.svg))
			if c.wantHit == "" {
				if detail != "" {
					t.Fatalf("expected pass, got: %s", detail)
				}
				return
			}
			if !strings.Contains(detail, c.wantHit) {
				t.Fatalf("expected detail containing %q, got: %q", c.wantHit, detail)
			}
		})
	}
}
