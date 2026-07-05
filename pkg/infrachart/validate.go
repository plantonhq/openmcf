package infrachart

import (
	"fmt"
	"regexp"
	"strings"

	"buf.build/go/protovalidate"
	"github.com/pkg/errors"
	"github.com/plantonhq/planton/internal/manifest"
)

// Validation renders with fixed org/env slugs, mirroring the control plane's behavior of
// always injecting the caller's org/env over any declared params. The values are
// slug-shaped so rendered names look like real resource names.
const (
	validationOrg = "planton"
	validationEnv = "dev"
)

// DocError is one failed validation on one rendered manifest document.
type DocError struct {
	// Template is the chart-relative template file the document came from.
	Template string
	// DocIndex is the zero-based document index within that template's rendered output.
	DocIndex int
	// Kind is the manifest's kind when it could be determined.
	Kind string
	// Name is the manifest's metadata.name when it could be determined.
	Name string
	// Err is what failed: render, load (unknown kind / unknown field), protovalidate, or
	// a valueFrom reference that does not resolve.
	Err error
}

// VariantResult is the outcome of validating one render variant of a chart.
type VariantResult struct {
	// Name identifies the variant: "defaults", or "<param>=<value>" for a flipped toggle.
	Name string
	// Docs is how many manifest documents the variant rendered and checked.
	Docs int
	// Errors are the failures, empty when the variant is fully valid.
	Errors []DocError
}

// Report is the outcome of validating one chart across all its render variants.
type Report struct {
	// ChartDir is the directory the chart was loaded from.
	ChartDir string
	// Variants holds one result per render variant.
	Variants []VariantResult
}

// Valid reports whether every variant validated cleanly.
func (r *Report) Valid() bool {
	for _, v := range r.Variants {
		if len(v.Errors) > 0 {
			return false
		}
	}
	return len(r.Variants) > 0
}

// ValidateChart loads and validates a chart directory offline. Beyond the defaults
// render, every bool param is flipped once so each conditional manifest is exercised in
// both branches -- one variant per toggle keeps the run linear in the number of toggles
// while still rendering every conditional block at least once in each state.
func ValidateChart(dir string) (*Report, error) {
	chart, err := Load(dir)
	if err != nil {
		return nil, err
	}

	report := &Report{ChartDir: dir}

	variants := []struct {
		name      string
		overrides map[string]any
	}{{name: "defaults"}}
	for _, name := range chart.BoolParams() {
		flipped := !boolParamDefault(chart, name)
		variants = append(variants, struct {
			name      string
			overrides map[string]any
		}{
			name:      fmt.Sprintf("%s=%v", name, flipped),
			overrides: map[string]any{name: flipped},
		})
	}

	for _, variant := range variants {
		values := ParamValues(chart.Params, variant.overrides)
		result := VariantResult{Name: variant.name}
		for _, templateName := range chart.Templates() {
			validateTemplate(chart, templateName, values, &result)
		}
		report.Variants = append(report.Variants, result)
	}
	return report, nil
}

// validateTemplate renders one template file and validates every document it produces.
func validateTemplate(chart *Chart, templateName string, values map[string]any, result *VariantResult) {
	rendered, err := RenderTemplate(templateName, chart.Template(templateName), values, validationOrg, validationEnv)
	if err != nil {
		result.Errors = append(result.Errors, DocError{Template: templateName, Err: err})
		return
	}

	for i, doc := range splitDocs(rendered) {
		result.Docs++
		if docErr := validateDoc(doc); docErr != nil {
			result.Errors = append(result.Errors, DocError{
				Template: templateName,
				DocIndex: i,
				Kind:     scrapeTopLevel(doc, "kind"),
				Name:     scrapeScalar(doc, "name"),
				Err:      docErr,
			})
		}
	}
}

// validateDoc runs the full offline gate on one rendered manifest document: typed load
// (catches unknown kinds and unknown/renamed fields), protovalidate on the spec, and
// valueFrom reference resolution. Violations are reported compactly (one line each) --
// this is a batch gate over many documents, not the interactive single-manifest flow.
func validateDoc(doc string) error {
	loaded, err := manifest.LoadManifestBytes([]byte(doc), "rendered manifest")
	if err != nil {
		return err
	}

	spec, err := manifest.ExtractSpec(loaded)
	if err != nil {
		return errors.Wrap(err, "extracting spec")
	}
	validator, err := protovalidate.New(protovalidate.WithDisableLazy(), protovalidate.WithMessages(spec))
	if err != nil {
		return errors.Wrap(err, "initializing validator")
	}
	if err := validator.Validate(spec); err != nil {
		var ve *protovalidate.ValidationError
		if errors.As(err, &ve) {
			msgs := make([]string, len(ve.Violations))
			for i, violation := range ve.Violations {
				msgs[i] = fmt.Sprintf("spec.%s: %s",
					protovalidate.FieldPathString(violation.Proto.GetField()),
					violation.Proto.GetMessage())
			}
			return errors.New(strings.Join(msgs, "\n"))
		}
		return err
	}

	if refErrs := CheckValueFromRefs(loaded); len(refErrs) > 0 {
		msgs := make([]string, len(refErrs))
		for i, re := range refErrs {
			msgs[i] = re.Error()
		}
		return errors.Errorf("unresolvable valueFrom reference(s):\n  %s", strings.Join(msgs, "\n  "))
	}
	return nil
}

// docSeparator matches a YAML document separator line, the convention chart templates and
// the control plane's template combiner both rely on.
var docSeparator = regexp.MustCompile(`(?m)^---\s*$`)

// splitDocs splits rendered template output into manifest documents, dropping documents
// that are empty or comment-only (a conditional block whose toggle is off renders empty).
func splitDocs(rendered string) []string {
	var docs []string
	for _, doc := range docSeparator.Split(rendered, -1) {
		if isBlank(doc) {
			continue
		}
		docs = append(docs, doc)
	}
	return docs
}

// isBlank reports whether a document has no content other than whitespace and comments.
func isBlank(doc string) bool {
	for _, line := range strings.Split(doc, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.HasPrefix(trimmed, "#") {
			return false
		}
	}
	return true
}

// boolParamDefault returns a bool param's declared default, treating a missing or
// non-bool default as false.
func boolParamDefault(chart *Chart, name string) bool {
	for _, p := range chart.Params {
		if p.Name == name {
			b, _ := p.Value.(bool)
			return b
		}
	}
	return false
}

// scrapeTopLevel reads a top-level scalar key from a YAML document textually -- used only
// to label error reports, never for validation decisions.
func scrapeTopLevel(doc, key string) string {
	re := regexp.MustCompile(`(?m)^` + key + `:\s*"?([^"\n]+)"?\s*$`)
	if m := re.FindStringSubmatch(doc); m != nil {
		return strings.TrimSpace(m[1])
	}
	return ""
}

// scrapeScalar reads the first occurrence of an indented scalar key (metadata.name in
// practice) -- like scrapeTopLevel, only for labeling error reports.
func scrapeScalar(doc, key string) string {
	re := regexp.MustCompile(`(?m)^\s+` + key + `:\s*"?([^"\n]+)"?\s*$`)
	if m := re.FindStringSubmatch(doc); m != nil {
		return strings.TrimSpace(m[1])
	}
	return ""
}
