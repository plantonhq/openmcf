package helmcrds

import (
	"context"
	"strings"
)

// Stamp keys written onto every derived CRD. They are the one contract both
// engines and every verifier read; the Terraform template quotes them as
// literals and a test in the generators package holds the literals equal to
// these constants.
const (
	// AnnotationSourceChart records "<repository>/<chart>" the CRD was derived
	// from (or the bundle URL template for bundle-branch kinds).
	AnnotationSourceChart = "planton.ai/crd-source-chart"
	// AnnotationSourceVersion records the chart version the CRD was derived
	// at. The never-downgrade check compares against it.
	AnnotationSourceVersion = "planton.ai/crd-source-version"
	// LabelSource carries the chart name so a module can select the CRDs it
	// owns with one label selector (label values cannot carry "/" or ":",
	// which is why the full source lives in the annotation).
	LabelSource = "planton.ai/crd-source"
)

// Source names where a kind's CRD set comes from at a pinned version. Exactly
// one of the two branches is used: the render branch when Chart is set, the
// bundle branch when BundleURL is set.
type Source struct {
	// Repository is the chart repository: an https:// index URL or an
	// oci://<registry>/<path> reference. Render branch.
	Repository string
	// Chart is the chart name inside Repository. Render branch.
	Chart string
	// Version is the exact chart version to render. Both branches use it for
	// the stamp; the bundle branch also substitutes it into BundleURL.
	Version string
	// Values are the release's OWN values documents, in helm -f order (later
	// documents override earlier ones), exactly as the release resource
	// receives them. The render must see what the install will see.
	Values []string
	// CRDOverride is one more YAML document merged LAST, carrying only the
	// chart's CRD switch turned on (for example "crds:\n  create: true\n").
	// The release itself installs with the switch off and skip_crds set, so
	// this override never reaches the cluster through Helm.
	CRDOverride string
	// APIVersions feeds .Capabilities.APIVersions for charts that gate
	// templates on them (a chart that emits cert-manager annotations only when
	// cert-manager.io/v1 is served). Omitting a version the chart checks makes
	// the render silently drop or reshape a CRD.
	APIVersions []string
	// KubeVersion feeds .Capabilities.KubeVersion; empty uses Helm's default.
	KubeVersion string

	// BundleURL is the upstream CRD bundle for charts that ship no CRDs at
	// all. It may contain "{{version}}", replaced with Version. Bundle
	// branch. Use the URL form that survives the next upstream release
	// (Apache's dist tree holds only the current release; archive.apache.org
	// holds every release).
	BundleURL string
}

// IsBundle reports whether the source is the bundle branch.
func (s Source) IsBundle() bool { return s.BundleURL != "" }

// ResolvedBundleURL substitutes the version into the bundle URL template.
func (s Source) ResolvedBundleURL() string {
	return strings.ReplaceAll(s.BundleURL, "{{version}}", s.Version)
}

// SourceDescription is what AnnotationSourceChart records: the chart's
// coordinates for the render branch, the URL template for the bundle branch.
func (s Source) SourceDescription() string {
	if s.IsBundle() {
		return s.BundleURL
	}
	return strings.TrimSuffix(s.Repository, "/") + "/" + s.Chart
}

// CRD is one derived CustomResourceDefinition, stamped and ready to apply.
type CRD struct {
	// Name is metadata.name (for example "opentelemetrycollectors.opentelemetry.io"),
	// the key every engine uses for the resource address so state addresses
	// stay stable across reorderings.
	Name string
	// YAML is the single-document manifest with the stamps applied.
	YAML string
}

// Derive produces the stamped CRD set for a source: render or fetch, split on
// document separators, keep only CustomResourceDefinitions, stamp each one.
// The release name and namespace are the identity the render runs under, so
// release-derived values inside a CRD (a conversion-webhook service, an
// inject-ca-from annotation) come out pointing at the real release.
func Derive(ctx context.Context, src Source, releaseName, namespace string) ([]CRD, error) {
	var documents []string
	var err error
	if src.IsBundle() {
		documents, err = fetchBundle(ctx, src)
	} else {
		documents, err = render(src, releaseName, namespace)
	}
	if err != nil {
		return nil, err
	}

	crds, err := filterCRDs(documents)
	if err != nil {
		return nil, err
	}
	if len(crds) == 0 {
		return nil, zeroCRDsFailure(src)
	}

	stamped := make([]CRD, 0, len(crds))
	for _, crd := range crds {
		s, err := stamp(crd, src)
		if err != nil {
			return nil, err
		}
		stamped = append(stamped, s)
	}
	return stamped, nil
}
