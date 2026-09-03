// Package cataloglogo is the machine-enforced guard on every component's
// logo.svg. The division of labor with its siblings is deliberate:
// pkg/anatomy checks that logo.svg EXISTS, cataloglogo checks what the file
// IS, and pkg/catalogbundle projects the logo's versionless URL without ever
// gating. Three rules, one per thing that has actually gone wrong in this
// tree:
//
//   - malformed-svg: the file must be a browser-renderable SVG document.
//     Browsers parse image/svg+xml strictly -- a raw control byte, a "--"
//     inside a comment, or a duplicate attribute makes the ENTIRE image render
//     as a broken glyph -- and every one of those defect classes has shipped
//     here, invisible to every other gate. This rule has no baseline.
//   - missing-provenance: the file must say what it is, in a <desc> whose
//     first words name one of three classes: an official provider product
//     icon, a Planton-drawn glyph, or the Planton brand mark. The class is
//     how the update flow finds the drawn glyphs an official icon may one day
//     replace, how this gate knows which sharing is intended, and how the
//     next reader learns why a drawing looks the way it does.
//   - shared-glyph: two kinds of one provider never wear one glyph. On a
//     diagram the icon is the one identity signal nothing else supplies --
//     color says the family, shape says the nature, only the glyph says
//     WHICH thing -- so a chart built from one product's family (an instance,
//     its database, its user) must not read as a row of identical marks. The
//     one intended sharing is Planton's own kinds wearing the Planton brand
//     mark, and the files say so themselves.
//
// The law behind the last two rules, in the words the forge and update flows
// use: one kind, one glyph. A kind that IS a provider product wears the
// provider's official mark, unaltered. Planton's own kinds wear Planton's
// mark. Every other kind wears a Planton-drawn glyph that says what that kind
// is, in the provider's icon language, never containing or modifying a
// provider mark. When a provider publishes an official icon for a kind wearing
// a drawn glyph, the official one replaces it.
//
// The walk is keyed off the kind registry (crkreflect), never directory
// globs, mirroring pkg/anatomy: the gate sees exactly the components the
// product serves. A component without a logo file is anatomy's finding, not
// this gate's -- absence is skipped here so one defect reports in one place.
//
// The CI lane is .github/workflows/lint.catalog-logo.yaml. Provenance and
// sharing carry a burn-down baseline (baseline.yaml, anatomy's shape) for the
// providers whose logo sets have not yet been judged; malformed-svg never does.
package cataloglogo

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/plantonhq/planton/pkg/crkreflect"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
)

// Rule identifiers. Stable: baseline.yaml entries reference them.
const (
	RuleMalformedSVG      = "malformed-svg"
	RuleMissingProvenance = "missing-provenance"
	RuleSharedGlyph       = "shared-glyph"
)

// Provenance classes, matched against the first words of the <desc>. A logo
// that opens its description with one of these has said what it is.
const (
	ProvenanceOfficial = "Official "      // e.g. "Official Google Cloud product icon: <library>/<file>"
	ProvenanceDrawn    = "Planton-drawn " // e.g. "Planton-drawn glyph: <what it depicts>"
	ProvenanceBrand    = "Planton brand mark"
)

// Violation is one logo defect at one path.
type Violation struct {
	// Path is repo-root-relative.
	Path string
	// Rule is the stable rule identifier (baseline entries key on it).
	Rule string
	// Detail says what is broken and what to do, in plain language.
	Detail string
}

// ID is the stable baseline identifier: "<path>:<rule>".
func (v Violation) ID() string { return v.Path + ":" + v.Rule }

// logo is one registered kind's logo as the gate read it.
type logo struct {
	provider string
	kind     string
	path     string
	hash     string
	class    string // one of the Provenance* prefixes, or "" when the file names none
}

// Check validates every registered kind's logo.svg under repoRoot. Missing
// files are skipped (anatomy owns existence); everything present must be a
// strict-XML SVG document, must name its provenance, and must not wear a
// glyph another kind of the same provider wears -- unless both are the
// Planton brand mark.
func Check(repoRoot string) ([]Violation, error) {
	var vs []Violation
	var logos []logo
	metaByKind := crkreflect.KindToKindMetaMap()
	for _, kind := range crkreflect.KindsList() {
		if kind == cloudresourcekind.CloudResourceKind_unspecified {
			continue
		}
		meta := metaByKind[kind]
		if meta == nil {
			continue
		}
		providerDir := crkreflect.ProviderDirName(meta.GetProvider())
		kindDir := strings.ToLower(kind.String())
		rel := filepath.Join("catalog", providerDir, kindDir, "logo.svg")
		data, err := os.ReadFile(filepath.Join(repoRoot, rel))
		if err != nil {
			if os.IsNotExist(err) {
				continue // anatomy's finding, not ours
			}
			return nil, err
		}
		if detail := validateSVG(data); detail != "" {
			vs = append(vs, Violation{Path: rel, Rule: RuleMalformedSVG, Detail: detail})
			continue // a file browsers cannot draw has no glyph to compare
		}
		sum := sha256.Sum256(data)
		l := logo{provider: providerDir, kind: kindDir, path: rel, hash: hex.EncodeToString(sum[:]), class: provenanceClass(data)}
		if l.class == "" {
			vs = append(vs, Violation{Path: rel, Rule: RuleMissingProvenance, Detail: "no provenance: open the <desc> with one of \"Official <provider> product icon: <library>/<file>\", \"Planton-drawn glyph: <what it depicts>\", or \"Planton brand mark\", so the update flow, this gate, and the next reader know what this file is"})
		}
		logos = append(logos, l)
	}
	vs = append(vs, sharedGlyphViolations(logos)...)
	sort.Slice(vs, func(i, j int) bool { return vs[i].ID() < vs[j].ID() })
	return vs, nil
}

// sharedGlyphViolations reports, at each kind, the other kinds of the same
// provider wearing byte-identical logo content -- except when every wearer is
// the Planton brand mark, the law's one intended sharing.
func sharedGlyphViolations(logos []logo) []Violation {
	byProviderHash := map[string][]logo{}
	for _, l := range logos {
		key := l.provider + "\x00" + l.hash
		byProviderHash[key] = append(byProviderHash[key], l)
	}
	var vs []Violation
	for _, group := range byProviderHash {
		if len(group) < 2 {
			continue
		}
		allBrand := true
		for _, l := range group {
			if l.class != ProvenanceBrand {
				allBrand = false
				break
			}
		}
		if allBrand {
			continue
		}
		for _, l := range group {
			var others []string
			for _, o := range group {
				if o.kind != l.kind {
					others = append(others, o.kind)
				}
			}
			sort.Strings(others)
			vs = append(vs, Violation{Path: l.path, Rule: RuleSharedGlyph, Detail: fmt.Sprintf("wears the same glyph as %s; one kind, one glyph -- an official product mark only for the kind that IS the product, a Planton-drawn glyph that says what this kind is for every other (only Planton's own kinds share Planton's brand mark)", strings.Join(others, ", "))})
		}
	}
	return vs
}

// provenanceClass returns the Provenance* prefix the logo's <desc> opens
// with, or "" when the file has no <desc> or opens it with anything else.
func provenanceClass(data []byte) string {
	dec := xml.NewDecoder(strings.NewReader(string(data)))
	dec.Strict = true
	inDesc := false
	var desc strings.Builder
	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}
		switch t := tok.(type) {
		case xml.StartElement:
			if t.Name.Local == "desc" {
				inDesc = true
			}
		case xml.CharData:
			if inDesc {
				desc.Write(t)
			}
		case xml.EndElement:
			if t.Name.Local == "desc" && inDesc {
				text := strings.TrimSpace(desc.String())
				for _, prefix := range []string{ProvenanceOfficial, ProvenanceDrawn, ProvenanceBrand} {
					if strings.HasPrefix(text, prefix) {
						return prefix
					}
				}
				return ""
			}
		}
	}
	return ""
}

// validateSVG returns "" for a browser-renderable SVG document, or a
// plain-language defect description. Strict XML decoding catches the raw
// control bytes and illegal "--" comment sequences browsers refuse; the
// duplicate-attribute check exists because Go's decoder tolerates what
// browsers hard-fail on (a live logo shipped with fill-rule declared twice).
func validateSVG(data []byte) string {
	// Raw-byte pre-scan: control characters below 0x20 (except tab, LF, CR)
	// are illegal ANYWHERE in an XML 1.0 document -- comments included.
	// Go's decoder validates them in character data but tolerates them in
	// comments, and that exact gap shipped five live logos carrying a raw
	// 0x14 where an em dash was intended.
	for i, b := range data {
		if b < 0x20 && b != '\t' && b != '\n' && b != '\r' {
			return fmt.Sprintf("not well-formed XML: illegal control byte 0x%02x at offset %d -- browsers refuse the whole image", b, i)
		}
	}
	dec := xml.NewDecoder(strings.NewReader(string(data)))
	dec.Strict = true
	sawRoot := false
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Sprintf("not well-formed XML: %v", err)
		}
		start, ok := tok.(xml.StartElement)
		if !ok {
			continue
		}
		if !sawRoot {
			sawRoot = true
			if start.Name.Local != "svg" {
				return fmt.Sprintf("root element is <%s>, not <svg> -- browsers refuse the whole image", start.Name.Local)
			}
		}
		seen := map[xml.Name]bool{}
		for _, a := range start.Attr {
			if seen[a.Name] {
				return fmt.Sprintf("duplicate attribute %q on <%s> -- valid to Go's decoder, fatal to a browser's", a.Name.Local, start.Name.Local)
			}
			seen[a.Name] = true
		}
	}
	if !sawRoot {
		return "no <svg> root element -- the file has SVG content but is not an SVG document"
	}
	return ""
}
