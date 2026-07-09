package infrachart

import (
	"fmt"
	"io"
	"os"
	"strings"

	"buf.build/go/protovalidate"
	"github.com/pkg/errors"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/internal/manifest"
	"github.com/plantonhq/planton/pkg/crkreflect"
	"github.com/plantonhq/planton/pkg/reflection/metadatareflect"
	"google.golang.org/protobuf/proto"
	yamlv3 "gopkg.in/yaml.v3"
)

// Severity of a validation issue.
type Severity string

const (
	// SeverityError fails the gate.
	SeverityError Severity = "error"
	// SeverityWarning is surfaced but does not fail the gate.
	SeverityWarning Severity = "warning"
)

// Issue is one problem found while validating a chart, attributed to the
// template file it came from (the honest granularity — rendering merges
// templates, so line numbers do not survive).
type Issue struct {
	Severity     Severity
	File         string
	ResourceKind string
	ResourceName string
	Message      string
}

// Doc is one rendered cloud-resource document.
type Doc struct {
	File string
	Kind string
	Name string
	YAML []byte
	Msg  proto.Message

	refUses []refUse
}

// Report is the outcome of validating one chart.
type Report struct {
	ChartName string
	Docs      []Doc
	Issues    []Issue
}

// HasErrors reports whether any issue is error-severity.
func (r *Report) HasErrors() bool {
	for _, issue := range r.Issues {
		if issue.Severity == SeverityError {
			return true
		}
	}
	return false
}

// Options configures a validation run.
type Options struct {
	// Org and Env are bound as the reserved org/env template variables (the
	// platform injects the real ones at render time). Sensible synthetic
	// defaults are applied when empty.
	Org string
	Env string

	// Set overrides param values by name (and may override org/env), letting
	// the gate exercise feature toggles in their non-default positions.
	Set map[string]string
}

// Validate loads the chart at dir, renders it with its default values (plus
// any overrides), and validates every rendered document and reference. It
// returns an error only when the chart cannot be processed at all; rendering
// and validation problems are returned as Issues in the Report.
func Validate(dir string, opts Options) (*Report, error) {
	chart, err := LoadDir(dir)
	if err != nil {
		return nil, err
	}

	org := opts.Org
	if org == "" {
		org = "acme"
	}
	env := opts.Env
	if env == "" {
		env = "dev"
	}

	ctx, placeholders, err := buildContext(chart.Params, org, env, opts.Set)
	if err != nil {
		return nil, err
	}

	report := &Report{ChartName: chart.Name}

	for _, tpl := range chart.Templates {
		report.validateTemplate(tpl, ctx, placeholders)
	}

	report.checkReferences()
	return report, nil
}

// validateTemplate renders one template file and validates every document in
// it, appending issues and docs to the report.
func (r *Report) validateTemplate(tpl TemplateFile, ctx map[string]any, placeholders []string) {
	addError := func(format string, args ...any) {
		r.Issues = append(r.Issues, Issue{Severity: SeverityError, File: tpl.Name, Message: fmt.Sprintf(format, args...)})
	}

	for _, finding := range scanBannedConstructs(tpl.Content) {
		addError("%s", finding)
	}

	rendered, err := renderTemplate(tpl.Name, tpl.Content, ctx)
	if err != nil {
		addError("template rendering failed: %v", err)
		return
	}

	// A value-less param renders as its "<name>" placeholder — legal for the
	// platform (the value arrives from the InfraProject), but this gate's
	// promise is that DEFAULTS render deployable, so depending on one is
	// worth a warning even before schema validation rejects it.
	for _, placeholder := range placeholders {
		if strings.Contains(rendered, placeholder) {
			r.Issues = append(r.Issues, Issue{
				Severity: SeverityWarning, File: tpl.Name,
				Message: fmt.Sprintf("template uses param %s, which has no default value — supply one in values.yaml or via --set", placeholder),
			})
		}
	}

	docs, err := splitDocs(rendered)
	if err != nil {
		addError("rendered output is not valid YAML: %v", err)
		return
	}

	for _, docYaml := range docs {
		r.validateDoc(tpl.Name, docYaml)
	}
}

// validateDoc validates one rendered document: strict schema load, presence
// of metadata.name, the spec's full protovalidate/CEL rule set, and collects
// its valueFrom references for the cross-document pass.
func (r *Report) validateDoc(file string, docYaml []byte) {
	msg, err := loadRenderedDoc(docYaml)
	if err != nil {
		r.Issues = append(r.Issues, Issue{Severity: SeverityError, File: file, Message: compactError(err)})
		return
	}

	doc := Doc{
		File: file,
		Kind: string(msg.ProtoReflect().Descriptor().Name()),
		YAML: docYaml,
		Msg:  msg,
	}
	if meta := metadatareflect.ExtractMetadata(msg); meta != nil {
		doc.Name = meta.GetName()
	}

	addIssue := func(severity Severity, format string, args ...any) {
		r.Issues = append(r.Issues, Issue{
			Severity: severity, File: file,
			ResourceKind: doc.Kind, ResourceName: doc.Name,
			Message: fmt.Sprintf(format, args...),
		})
	}

	if doc.Name == "" {
		addIssue(SeverityError, "metadata.name is required (it is the resource's identity in the chart's dependency graph)")
	}

	spec, err := manifest.ExtractSpec(msg)
	if err != nil {
		addIssue(SeverityError, "manifest has no spec: %v", err)
	} else {
		validator, err := protovalidate.New(protovalidate.WithDisableLazy(), protovalidate.WithMessages(spec))
		if err != nil {
			addIssue(SeverityError, "failed to initialize validator: %v", err)
		} else if validationErr := validator.Validate(spec); validationErr != nil {
			addIssue(SeverityError, "spec validation failed: %s", compactError(validationErr))
		}
	}

	doc.refUses = collectRefUses(msg)
	r.Docs = append(r.Docs, doc)
}

// checkReferences runs the cross-document pass: per-reference integrity
// (annotations, field-path resolution), in-chart target resolution, and
// dependency-cycle detection.
func (r *Report) checkReferences() {
	// Index the chart's own resources by (kind, name). A duplicate identity
	// is an error — the platform's dependency graph could not tell the two
	// resources apart.
	index := map[refTarget]int{}
	for i, doc := range r.Docs {
		if doc.Name == "" {
			continue
		}
		key := refTarget{kind: crkreflect.KindFromString(doc.Kind), name: doc.Name}
		if prev, dup := index[key]; dup {
			r.Issues = append(r.Issues, Issue{
				Severity: SeverityError, File: doc.File, ResourceKind: doc.Kind, ResourceName: doc.Name,
				Message: fmt.Sprintf("duplicate resource identity: %s %q is also defined in %s", doc.Kind, doc.Name, r.Docs[prev].File),
			})
			continue
		}
		index[key] = i
	}

	edges := make([][]int, len(r.Docs))
	for i, doc := range r.Docs {
		for _, use := range doc.refUses {
			target, problems := checkRef(use)
			for _, p := range problems {
				r.Issues = append(r.Issues, Issue{
					Severity: SeverityError, File: doc.File, ResourceKind: doc.Kind, ResourceName: doc.Name,
					Message: p,
				})
			}
			if target.kind == cloudresourcekind.CloudResourceKind_unspecified || target.name == "" {
				continue
			}
			// A reference that names an env explicitly points outside this
			// chart's deployment by design — never an in-chart edge.
			if use.ref.GetEnv() != "" {
				continue
			}
			if targetIdx, ok := index[target]; ok {
				edges[i] = append(edges[i], targetIdx)
			} else {
				r.Issues = append(r.Issues, Issue{
					Severity: SeverityWarning, File: doc.File, ResourceKind: doc.Kind, ResourceName: doc.Name,
					Message: fmt.Sprintf("%s: references %s %q, which this chart does not define — it must already exist in the target environment at deploy time", use.fieldPath, target.kind, target.name),
				})
			}
		}
	}

	if cycle := findCycle(edges); cycle != nil {
		var parts []string
		for _, idx := range cycle {
			parts = append(parts, fmt.Sprintf("%s %q", r.Docs[idx].Kind, r.Docs[idx].Name))
		}
		r.Issues = append(r.Issues, Issue{
			Severity: SeverityError, File: r.Docs[cycle[0]].File,
			ResourceKind: r.Docs[cycle[0]].Kind, ResourceName: r.Docs[cycle[0]].Name,
			Message: "dependency cycle: " + strings.Join(parts, " -> "),
		})
	}
}

// loadRenderedDoc writes a rendered document to a temp file and loads it
// through the CLI's canonical manifest loader — the same strict path (kind
// registry, protojson with unknown fields rejected, proto field defaults)
// every other planton command uses.
func loadRenderedDoc(docYaml []byte) (proto.Message, error) {
	tmp, err := os.CreateTemp("", "planton-chart-doc-*.yaml")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create temp file")
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(docYaml); err != nil {
		tmp.Close()
		return nil, errors.Wrap(err, "failed to write temp file")
	}
	if err := tmp.Close(); err != nil {
		return nil, errors.Wrap(err, "failed to close temp file")
	}
	return manifest.LoadManifest(tmp.Name())
}

// splitDocs splits rendered multi-document YAML into per-document byte
// slices, skipping empty documents (a template whose conditionals all
// evaluate false legitimately renders to nothing).
func splitDocs(rendered string) ([][]byte, error) {
	var out [][]byte
	decoder := yamlv3.NewDecoder(strings.NewReader(rendered))
	for {
		var node yamlv3.Node
		err := decoder.Decode(&node)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		if isNullNode(&node) {
			continue
		}
		docBytes, err := yamlv3.Marshal(&node)
		if err != nil {
			return nil, err
		}
		out = append(out, docBytes)
	}
	return out, nil
}

func isNullNode(node *yamlv3.Node) bool {
	if node == nil || node.Kind == 0 {
		return true
	}
	if node.Kind == yamlv3.DocumentNode {
		if len(node.Content) == 0 {
			return true
		}
		return node.Content[0].Tag == "!!null"
	}
	return node.Tag == "!!null"
}

// compactError flattens a (possibly multi-line, CLI-decorated) error into a
// single readable message for the issue report.
func compactError(err error) string {
	msg := err.Error()
	msg = strings.TrimPrefix(msg, "validation error:")
	msg = strings.ReplaceAll(msg, "\n", " ")
	return strings.Join(strings.Fields(msg), " ")
}
