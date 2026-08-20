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

// customResourceDefinitions applies the module-owned
// plantonplatforms.planton.ai CRD from the file staged at
// vars.CrdDirectory — extracted from the published chart package at
// vars.DefaultChartVersion, so the schema always matches the operator the
// default install schedules. One ConfigGroup per file, keyed by the CRD's
// OWN metadata.name (parsed from the file — never a positional index), so
// state addresses stay stable.
//
// WHY MODULE-OWNED: the chart ships its CRD in the crds/ directory —
// Helm's install-once posture (created on first install, never upgraded,
// never removed). Owning the CRD here upgrades it deliberately with every
// chart_version bump (the staged file re-stages with the pin) and makes
// keep-on-uninstall a guarantee instead of an accident. The release always
// installs with SkipCrds (main.go), so the two owners never meet.
//
// KEEP-ON-UNINSTALL (the load-bearing option): the CRD applies with
// retainOnDelete — on destroy Pulumi drops it from state WITHOUT deleting
// it from the cluster, so destroying the operator never cascade-deletes
// PlantonPlatform resources (and the platforms behind them). This is the
// exact semantic twin of the Terraform module's kubectl_manifest
// apply_only = true.
//
// The option must reach the ConfigGroup's CHILDREN (the actual CRD
// resource), and neither yaml package forwards ordinary resource options
// to them — hence the CLASSIC yaml package with retainOnDelete delivered
// through a resource TRANSFORMATION, the one mechanism the SDK propagates
// down the parent chain to in-process children.
//
// RE-ADOPTION: a destroy leaves the CRD on the cluster by design, so the
// next install finds it already there — which is why the caller routes
// this function through the UPSERT provider (server-side apply adopts the
// existing object instead of failing AlreadyExists). Ordinary resources
// keep create-conflict semantics through the plain provider.
func customResourceDefinitions(ctx *pulumi.Context, locals *Locals,
	upsertProvider pulumi.ProviderResource) ([]pulumi.Resource, error) {

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
			pulumi.Provider(upsertProvider),
			// The keep-on-uninstall mechanism — see the function comment.
			retainOnDelete,
		)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to apply CRD %s", manifest.Metadata.Name)
		}
		crds = append(crds, crd)
	}

	// FAIL LOUDLY when the staged CRD file did not travel with the module:
	// an empty directory would silently apply ZERO CRDs and the operator
	// would run against whatever CRD happens to exist. One is the staged
	// count at chart 0.7.1 — re-stage ../crds and update this count
	// together with DefaultChartVersion. Twin of the Terraform module's
	// precondition.
	if len(crds) != vars.ExpectedCrdCount {
		return nil, errors.Errorf("the staged CRD directory %s carries %d CRDs, expected %d — the module owns the CRD lifecycle and cannot install without its full staged set", vars.CrdDirectory, len(crds), vars.ExpectedCrdCount)
	}

	return crds, nil
}
