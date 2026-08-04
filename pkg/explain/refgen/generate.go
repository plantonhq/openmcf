package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/plantonhq/planton/internal/manifest"
	"github.com/plantonhq/planton/pkg/explain"
)

// providerPathPrefix roots every cloud-component proto file; the provider is
// the path segment right after it. Deriving both the provider and the output
// directory from the descriptor's source path keeps this generator agnostic
// to the version segment (v1 today, later channels tomorrow) -- never
// reconstruct these paths from kind metadata.
const providerPathPrefix = "dev/planton/provider/"

// Summary is one generation run's full account: the files to write plus
// every degradation, so a caller can render an honest report and CI can fail
// on real catalog bugs instead of silently shipping thinner pages.
type Summary struct {
	// Files maps repo-relative output paths (apis/.../reference.md) to
	// rendered content.
	Files map[string]string
	// MissingManifests lists kinds with no iac/hack/manifest.yaml; their
	// pages render without an Example section.
	MissingManifests []string
	// InvalidManifests lists kinds whose hack manifest failed loading or
	// protovalidation. Each entry is a REAL catalog bug: the page renders
	// without the broken example, and the run must be reported as failed.
	InvalidManifests []InvalidManifest
}

// InvalidManifest names one hack manifest that failed to load or validate.
type InvalidManifest struct {
	Kind string
	Err  error
}

// Generate renders the reference page for every kind under the provider
// filter (empty = every kind). It touches nothing on disk except reads;
// writing is the caller's decision so tests can byte-compare in memory.
func Generate(repoRoot, provider string) (*Summary, error) {
	engine := explain.DefaultEngine()
	summary := &Summary{Files: map[string]string{}}

	for _, name := range explain.KindNames() {
		res, err := explain.ResolveKindName(name)
		if err != nil {
			return nil, errors.Wrapf(err, "resolve kind %s", name)
		}
		protoDir := filepath.Dir(res.Message.ParentFile().Path())
		if !strings.HasPrefix(protoDir, providerPathPrefix) {
			continue
		}
		if provider != "" && !strings.HasPrefix(protoDir, providerPathPrefix+provider+"/") {
			continue
		}

		report, err := engine.Explain(res, nil)
		if err != nil {
			return nil, errors.Wrapf(err, "explain %s", name)
		}

		kindDir := filepath.Join(repoRoot, "apis", protoDir)
		opts := explain.MarkdownOptions{SeeAlso: seeAlso(kindDir)}
		switch example, err := loadValidatedExample(res, kindDir); {
		case err != nil:
			summary.InvalidManifests = append(summary.InvalidManifests, InvalidManifest{Kind: name, Err: err})
		case example == "":
			summary.MissingManifests = append(summary.MissingManifests, name)
		default:
			opts.ExampleYAML = example
		}

		summary.Files[filepath.Join("apis", protoDir, "reference.md")] = explain.RenderMarkdown(report, opts)
	}
	return summary, nil
}

// loadValidatedExample reads the kind's hack manifest and admits it as the
// page's example only after it round-trips into the typed message and passes
// protovalidate. An empty string with nil error means no manifest exists
// (graceful degradation); an error means the manifest exists and is broken.
func loadValidatedExample(res explain.Resource, kindDir string) (string, error) {
	path := filepath.Join(kindDir, "iac", "hack", "manifest.yaml")
	raw, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", errors.Wrapf(err, "read %s", path)
	}
	loaded, err := manifest.LoadManifestBytes(raw, path)
	if err != nil {
		return "", errors.Wrapf(err, "load %s", path)
	}
	if got := loaded.ProtoReflect().Descriptor().FullName(); got != res.Message.FullName() {
		return "", errors.Errorf("%s declares kind %s, expected %s", path, got, res.Message.FullName())
	}
	if err := manifest.ValidateLoaded(loaded); err != nil {
		return "", errors.Wrapf(err, "validate %s", path)
	}
	return string(raw), nil
}

// seeAlso links the page to the kind's hand-written prose when it exists.
// Reference pages link to essays, never duplicate them.
func seeAlso(kindDir string) []explain.MarkdownLink {
	var links []explain.MarkdownLink
	if _, err := os.Stat(filepath.Join(kindDir, "README.md")); err == nil {
		links = append(links, explain.MarkdownLink{Title: "Overview", Path: "./README.md"})
	}
	if _, err := os.Stat(filepath.Join(kindDir, "docs", "README.md")); err == nil {
		links = append(links, explain.MarkdownLink{Title: "Design notes", Path: "./docs/README.md"})
	}
	return links
}

// SortedPaths returns the output paths in stable order for writing and
// reporting.
func (s *Summary) SortedPaths() []string {
	paths := make([]string, 0, len(s.Files))
	for path := range s.Files {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return paths
}
