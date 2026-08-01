package module

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	pulumiyaml "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"sigs.k8s.io/yaml"
)

// crdManifest is the minimal shape needed to read a CRD's identity out of
// its manifest — the resource key below.
type crdManifest struct {
	Metadata struct {
		Name string `json:"name"`
	} `json:"metadata"`
}

// customResourceDefinitions applies the module-owned opentelemetry.io CRDs
// from the files staged at vars.CrdDirectory (the chart 0.120.0 CRD set),
// one ConfigGroup per CRD, keyed by the CRD's OWN metadata.name (parsed
// from the file — never a positional index), so state addresses stay
// stable across file reorderings.
//
// THE STAGED FILES ARE RENDERED, TOKENIZED CHART TEMPLATES: unlike
// plain-YAML CRD bundles, this chart TEMPLATES its CRDs — the collector
// CRD carries the cert-manager.io/inject-ca-from annotation and a
// version-conversion webhook clientConfig, both derived from the
// RELEASE's identity. The staged files were rendered from the pinned
// chart with those release-derived values replaced by
// __PLANTON_RELEASE_NAME__ / __PLANTON_NAMESPACE__ tokens, substituted
// here (and identically in the Terraform module's crd_manifests) — so
// the kept CRDs always point at THIS release's webhook Service and
// cert-manager Certificate.
//
// KEEP-ON-UNINSTALL (the load-bearing option): every CRD applies with
// retainOnDelete — on destroy Pulumi drops the CRDs from state WITHOUT
// deleting them from the cluster, so destroying the operator never
// cascade-deletes OpenTelemetryCollector resources. This is the exact
// semantic twin of the Terraform module's kubectl_manifest
// apply_only = true ("When true, Delete is a no-op" in the provider
// source).
//
// The option must reach the ConfigGroup's CHILDREN (the actual CRD
// resources), and neither yaml package forwards ordinary resource options
// to them: the classic yaml SDK's GetChildOptions passes only
// parent/version/pluginDownloadURL, and yaml/v2 is a REMOTE component
// whose children are created provider-side, beyond the reach of ANY
// SDK-side option or transformation (both verified in the pinned
// pulumi-kubernetes source). Hence the CLASSIC yaml package here, with
// retainOnDelete delivered through a resource TRANSFORMATION — the one
// mechanism the SDK propagates down the parent chain to in-process
// children.
//
// The chart itself never touches these CRDs: buildHelmValues pins
// crds.create false unconditionally.
func customResourceDefinitions(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider pulumi.ProviderResource) ([]pulumi.Resource, error) {

	retainOnDelete := pulumi.Transformations([]pulumi.ResourceTransformation{
		func(args *pulumi.ResourceTransformationArgs) *pulumi.ResourceTransformationResult {
			return &pulumi.ResourceTransformationResult{
				Props: args.Props,
				Opts:  append(args.Opts, pulumi.RetainOnDelete(true)),
			}
		},
	})

	entries, err := os.ReadDir(vars.CrdDirectory)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read staged CRD directory %s", vars.CrdDirectory)
	}

	var crds []pulumi.Resource
	// os.ReadDir returns entries sorted by filename — deterministic
	// resource registration order.
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}

		crdFilePath := filepath.Join(vars.CrdDirectory, entry.Name())
		content, err := os.ReadFile(crdFilePath)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read staged CRD file %s", crdFilePath)
		}

		// The token substitution — twin of the Terraform module's
		// replace() pair. Tokens appear only in the collector CRD; the
		// substitution is a harmless no-op on the other files.
		substituted := strings.ReplaceAll(string(content), "__PLANTON_RELEASE_NAME__", locals.ReleaseName)
		substituted = strings.ReplaceAll(substituted, "__PLANTON_NAMESPACE__", locals.Namespace)

		var manifest crdManifest
		if err := yaml.Unmarshal([]byte(substituted), &manifest); err != nil {
			return nil, errors.Wrapf(err, "failed to parse staged CRD file %s", crdFilePath)
		}
		if manifest.Metadata.Name == "" {
			return nil, errors.Errorf("staged CRD file %s carries no metadata.name", crdFilePath)
		}

		crd, err := pulumiyaml.NewConfigGroup(ctx, manifest.Metadata.Name,
			&pulumiyaml.ConfigGroupArgs{
				YAML: []string{substituted},
			},
			pulumi.Provider(kubernetesProvider),
			// The keep-on-uninstall mechanism — see the function comment.
			retainOnDelete,
		)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to apply CRD %s", manifest.Metadata.Name)
		}
		crds = append(crds, crd)
	}

	// FAIL LOUDLY when the staged CRD files did not travel with the
	// module: an empty directory would silently apply ZERO CRDs and the
	// operator would run against whatever CRDs happen to exist. Four is
	// the staged count at chart 0.120.0 — restage ../crds and update
	// this count together with DefaultChartVersion. Twin of the
	// Terraform module's precondition.
	if len(crds) != 4 {
		return nil, errors.Errorf("the staged CRD directory %s carries %d CRDs, expected 4 — the module owns the CRD lifecycle and cannot install without its full staged set", vars.CrdDirectory, len(crds))
	}

	return crds, nil
}
