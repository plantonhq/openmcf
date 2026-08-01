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

// customResourceDefinitions applies the module-owned OpenSearch CRDs from
// the files staged at vars.CrdDirectory (the chart 2.8.0 CRD set), one
// ConfigGroup per CRD, keyed by the CRD's OWN metadata.name (parsed from
// the file — never a positional index), so state addresses stay stable
// across file reorderings.
//
// KEEP-ON-UNINSTALL (the load-bearing option): every CRD applies with
// retainOnDelete — on destroy Pulumi drops the CRDs from state WITHOUT
// deleting them from the cluster, so destroying the operator never
// cascade-deletes OpenSearchCluster resources and their data. This is the
// exact semantic twin of the Terraform module's kubectl_manifest
// apply_only = true ("When true, Delete is a no-op" in the provider
// source).
//
// The option must reach the ConfigGroup's CHILDREN (the actual CRD
// resources), and neither yaml package forwards ordinary resource
// options to them: the classic yaml SDK's GetChildOptions passes only
// parent/version/pluginDownloadURL, and yaml/v2 is a REMOTE component
// whose children are created provider-side, beyond the reach of ANY
// SDK-side option or transformation (both verified in the pinned
// pulumi-kubernetes source). Hence the CLASSIC yaml package here, with
// retainOnDelete delivered through a resource TRANSFORMATION — the one
// mechanism the SDK propagates down the parent chain to in-process
// children.
//
// The chart itself never touches these CRDs: buildHelmValues pins
// installCRDs: false unconditionally.
func customResourceDefinitions(ctx *pulumi.Context,
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

		var manifest crdManifest
		if err := yaml.Unmarshal(content, &manifest); err != nil {
			return nil, errors.Wrapf(err, "failed to parse staged CRD file %s", crdFilePath)
		}
		if manifest.Metadata.Name == "" {
			return nil, errors.Errorf("staged CRD file %s carries no metadata.name", crdFilePath)
		}

		crd, err := pulumiyaml.NewConfigGroup(ctx, manifest.Metadata.Name,
			&pulumiyaml.ConfigGroupArgs{
				YAML: []string{string(content)},
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
	// module: an empty or partial directory would silently apply too few
	// CRDs and the operator would run against whatever CRDs happen to
	// exist (the class was caught live elsewhere: a lane "passed" riding
	// a previous install's retained CRDs). Ten is the staged count at
	// chart 2.8.0 — restage ../crds and update this count together with
	// DefaultChartVersion. Twin of the Terraform module's precondition.
	if len(crds) != 10 {
		return nil, errors.Errorf("the staged CRD directory %s carries %d CRDs, expected 10 — the module owns the CRD lifecycle and cannot install without its full staged set", vars.CrdDirectory, len(crds))
	}

	return crds, nil
}
