package cataloglogo

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestCatalogLogoGate is the gate: every registered kind's logo.svg must be a
// strict-XML SVG document (never baselined), must name its provenance, and
// must not share its glyph with another kind of the same provider (both
// carried by baseline.yaml for the providers whose logo sets predate the
// law). On failure, either fix the logo (the detail says how) or -- for a
// provider whose set has not been judged yet -- regenerate the baseline with
// PLANTON_REGEN_LOGO_BASELINE=1 and justify the growth in review.
func TestCatalogLogoGate(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate repo root")
	}
	repoRoot := filepath.Dir(filepath.Dir(filepath.Dir(thisFile)))
	vs, err := Check(repoRoot)
	if err != nil {
		t.Fatal(err)
	}
	baselinePath := filepath.Join(filepath.Dir(thisFile), "baseline.yaml")
	if os.Getenv("PLANTON_REGEN_LOGO_BASELINE") == "1" {
		if err := WriteBaseline(baselinePath, vs); err != nil {
			t.Fatalf("write baseline: %v", err)
		}
		t.Logf("baseline regenerated -- review the diff before committing (a baseline only ever shrinks)")
	}
	baseline, err := LoadBaseline(baselinePath)
	if err != nil {
		t.Fatalf("load baseline: %v", err)
	}
	res := Gate(vs, baseline)
	for _, v := range res.Malformed {
		t.Errorf("%s: %s", v.Path, v.Detail)
	}
	for _, v := range res.NewViolations {
		t.Errorf("%s [%s]: %s", v.Path, v.Rule, v.Detail)
	}
	for _, id := range res.StaleEntries {
		t.Errorf("stale baseline entry (no longer a violation): %s -- remove it from baseline.yaml", id)
	}
	if len(res.Malformed) > 0 {
		t.Fatalf("%d logo(s) will render as broken images in a browser -- fix the SVGs (this rule has no baseline)", len(res.Malformed))
	}
}

// judgedProviders are the providers whose logo sets have been judged glyph by
// glyph under the law: every kind wears its own glyph and names its
// provenance, with no baseline entry. A provider joins this list in the same
// change that empties its baseline entries, and never leaves it. Adding a
// provider's directory name here is the whole act of pinning it.
var judgedProviders = []string{"gcp", "cloudflare", "kubernetes", "auth0", "openfga"}

// TestJudgedProviderLogoSets pins every judged provider at zero: no violation
// under its catalog directory, no entry for it in baseline.yaml. An entry
// reappearing for a judged provider is the wave reopening, not debt -- fix the
// logo, never re-baseline it.
func TestJudgedProviderLogoSets(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate repo root")
	}
	repoRoot := filepath.Dir(filepath.Dir(filepath.Dir(thisFile)))
	vs, err := Check(repoRoot)
	if err != nil {
		t.Fatal(err)
	}
	baseline, err := LoadBaseline(filepath.Join(filepath.Dir(thisFile), "baseline.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	for _, provider := range judgedProviders {
		prefix := "catalog/" + provider + "/"
		t.Run(provider, func(t *testing.T) {
			for _, v := range vs {
				if strings.HasPrefix(v.Path, prefix) {
					t.Errorf("%s [%s]: %s", v.Path, v.Rule, v.Detail)
				}
			}
			for id := range baseline {
				if strings.HasPrefix(id, prefix) {
					t.Errorf("baseline carries a %s entry: %s -- %s's logo set is judged; fix the logo, never re-baseline it", provider, id, provider)
				}
			}
		})
	}
}

// TestSharedGlyphAndProvenanceDeliberateReds proves the two law rules on a
// scratch catalog: a copied glyph fails at both kinds by name, two brand marks
// may share, a drawn glyph without its <desc> is unlabelled, and a baseline
// entry naming a fixed defect is stale.
func TestSharedGlyphAndProvenanceDeliberateReds(t *testing.T) {
	const drawnA = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><desc>Planton-drawn glyph: a bell</desc><path d="M0 0h24v24H0z"/></svg>`
	const drawnB = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><desc>Planton-drawn glyph: a key</desc><path d="M0 0h12v12H0z"/></svg>`
	const brand = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 48 48"><desc>Planton brand mark</desc><rect width="48" height="48"/></svg>`
	const unlabelled = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M0 0h6v6H0z"/></svg>`
	const official = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><desc>Official Google Cloud product icon: legacy/cloud_sql.svg</desc><path d="M1 1h1v1H1z"/></svg>`

	mk := func(provider, kind, content string) logo {
		sum := sha256.Sum256([]byte(content))
		return logo{provider: provider, kind: kind, path: "catalog/" + provider + "/" + kind + "/logo.svg", hash: hex.EncodeToString(sum[:]), class: provenanceClass([]byte(content))}
	}

	// Provenance classes read from the file.
	if c := provenanceClass([]byte(drawnA)); c != ProvenanceDrawn {
		t.Fatalf("drawn class: got %q", c)
	}
	if c := provenanceClass([]byte(brand)); c != ProvenanceBrand {
		t.Fatalf("brand class: got %q", c)
	}
	if c := provenanceClass([]byte(official)); c != ProvenanceOfficial {
		t.Fatalf("official class: got %q", c)
	}
	if c := provenanceClass([]byte(unlabelled)); c != "" {
		t.Fatalf("unlabelled class: got %q", c)
	}

	// A glyph copied onto a second kind fails at BOTH kinds, naming the other.
	vs := sharedGlyphViolations([]logo{mk("gcp", "gcpa", drawnA), mk("gcp", "gcpb", drawnA), mk("gcp", "gcpc", drawnB)})
	if len(vs) != 2 {
		t.Fatalf("expected 2 shared-glyph violations, got %d: %+v", len(vs), vs)
	}
	for _, v := range vs {
		if v.Rule != RuleSharedGlyph || !strings.Contains(v.Detail, "one kind, one glyph") {
			t.Fatalf("unexpected violation %+v", v)
		}
	}
	// The same glyph on two PROVIDERS is not sharing: the gate is per provider.
	if vs := sharedGlyphViolations([]logo{mk("gcp", "gcpa", drawnA), mk("aws", "awsa", drawnA)}); len(vs) != 0 {
		t.Fatalf("cross-provider sharing must not fire, got %+v", vs)
	}
	// Planton's own kinds share Planton's mark by design.
	if vs := sharedGlyphViolations([]logo{mk("kubernetes", "kubernetesplantonrunner", brand), mk("kubernetes", "kubernetesplantonoperator", brand)}); len(vs) != 0 {
		t.Fatalf("brand marks may share, got %+v", vs)
	}
	// ...but a brand mark copied onto a non-Planton kind is still sharing.
	if vs := sharedGlyphViolations([]logo{mk("kubernetes", "kubernetesplantonrunner", brand), mk("kubernetes", "kubernetesredis", strings.Replace(brand, "Planton brand mark", "Planton-drawn glyph: not really", 1))}); len(vs) != 0 {
		// different content, different hash: no sharing -- the point is the next case
		t.Fatalf("unexpected %+v", vs)
	}
	brandOnRedis := logo{provider: "kubernetes", kind: "kubernetesredis", path: "catalog/kubernetes/kubernetesredis/logo.svg", hash: mk("kubernetes", "x", brand).hash, class: ProvenanceDrawn}
	if vs := sharedGlyphViolations([]logo{mk("kubernetes", "kubernetesplantonrunner", brand), brandOnRedis}); len(vs) != 2 {
		t.Fatalf("a brand mark on a non-Planton kind must fire at both kinds, got %+v", vs)
	}

	// The baseline gate: baselined debt passes, a fixed entry left behind is stale,
	// new debt fails, and malformed is reported regardless of the baseline.
	shared := Violation{Path: "catalog/azure/azurea/logo.svg", Rule: RuleSharedGlyph, Detail: "x"}
	res := Gate([]Violation{shared}, map[string]bool{shared.ID(): true})
	if !res.OK() {
		t.Fatalf("expected baselined debt to pass, got %+v", res)
	}
	res = Gate(nil, map[string]bool{shared.ID(): true})
	if len(res.StaleEntries) != 1 {
		t.Fatalf("expected one stale entry, got %+v", res)
	}
	res = Gate([]Violation{shared}, map[string]bool{})
	if len(res.NewViolations) != 1 {
		t.Fatalf("expected one new violation, got %+v", res)
	}
	malformed := Violation{Path: "catalog/aws/awsa/logo.svg", Rule: RuleMalformedSVG, Detail: "x"}
	res = Gate([]Violation{malformed}, map[string]bool{malformed.ID(): true})
	if len(res.Malformed) != 1 || res.OK() {
		t.Fatalf("malformed-svg must never be baselined, got %+v", res)
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
