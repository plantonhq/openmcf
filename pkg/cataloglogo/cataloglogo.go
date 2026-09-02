// Package cataloglogo is the machine-enforced guard that every component's
// logo.svg is a browser-renderable SVG document. The division of labor with
// its siblings is deliberate: pkg/anatomy checks that logo.svg EXISTS,
// cataloglogo checks that an existing logo actually PARSES as strict XML
// with an <svg> root, and pkg/catalogbundle projects the logo's versionless
// URL without ever gating. Browsers parse image/svg+xml strictly -- a raw
// control byte, a "--" inside a comment, or a duplicate attribute makes the
// ENTIRE image render as a broken glyph -- and every one of those defect
// classes has shipped in this tree, invisible to every other gate.
//
// The walk is keyed off the kind registry (crkreflect), never directory
// globs, mirroring pkg/anatomy: the gate sees exactly the components the
// product serves. A component without a logo file is anatomy's finding, not
// this gate's -- absence is skipped here so one defect reports in one place.
//
// The CI lane is .github/workflows/lint.catalog-logo.yaml. There is no
// baseline: the tree was brought to zero violations when the gate was
// introduced, and it must stay there.
package cataloglogo

import (
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

// Violation is one logo defect at one path.
type Violation struct {
	// Path is repo-root-relative.
	Path string
	// Detail says what is broken, in plain language.
	Detail string
}

// Check validates every registered kind's logo.svg under repoRoot. Missing
// files are skipped (anatomy owns existence); everything present must be a
// strict-XML SVG document.
func Check(repoRoot string) ([]Violation, error) {
	var vs []Violation
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
			vs = append(vs, Violation{Path: rel, Detail: detail})
		}
	}
	sort.Slice(vs, func(i, j int) bool { return vs[i].Path < vs[j].Path })
	return vs, nil
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
