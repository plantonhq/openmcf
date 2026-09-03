// Package keptcrds applies a Helm chart's CustomResourceDefinitions as
// module-owned, kept resources ahead of the chart's release: the Pulumi half of
// the catalog's derive-branch primitive for charts that carry CRDs.
//
// A module on the derive branch does three things, and this package does the
// first two for it: derive the CRD set from the pinned version (through
// pkg/kubernetes/helmcrds, in-process), apply each CRD keyed by its own name
// as a kept resource that re-adopts on reinstall and refuses a schema
// downgrade, and then install the release with CRDs skipped and the chart's
// own CRD switch pinned off (the module's release code). The Terraform twin is
// the generated helm_crds.tf every module carries byte-identically.
//
// Why the mechanics look the way they do, each verified live:
//
//   - Keep-on-uninstall rides a retainOnDelete resource TRANSFORMATION on a
//     classic-yaml ConfigGroup. Neither yaml package forwards ordinary options
//     to the CRD children (classic passes only parent/version/pluginDownloadURL;
//     yaml/v2 creates children provider-side, beyond any SDK option), and a
//     transformation is the one mechanism the SDK propagates down the parent
//     chain to in-process children.
//   - Re-adoption rides a dedicated upsert provider (server-side apply with
//     upsertExistingObjects and force): a destroy leaves kept CRDs on the
//     cluster by design, so the next install finds them there and a plain
//     create fails AlreadyExists. Only the CRDs use that provider; the release
//     keeps create-conflict semantics.
//   - The never-downgrade check reads the cluster in-process, before any
//     resource registers, through the same kubeconfig the provider uses.
package keptcrds

import (
	"context"
	"os"
	"time"

	"github.com/pkg/errors"
	kubernetesprovider "github.com/plantonhq/planton/catalog/kubernetes"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/kubernetes/pulumikubernetesprovider"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/plantonhq/planton/pkg/kubernetes/execcredential"
	"github.com/plantonhq/planton/pkg/kubernetes/helmcrds"
	"github.com/plantonhq/planton/pkg/kubernetes/kubeconfig"
	pulumiyaml "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

// Args describes one module's CRD set and how to keep it.
type Args struct {
	// Source is where the CRDs come from at the pinned version. Its Values
	// MUST be the release's own values documents in helm -f order, with only
	// the CRD switch turned on in CRDOverride: templated CRDs depend on
	// other values (fullname, webhook trust, ports), and a render from a
	// minimal set produces CRDs that point at the wrong release.
	Source helmcrds.Source
	// ReleaseName and Namespace are the identity the render runs under.
	ReleaseName string
	Namespace   string
	// Install false is the bring-your-own-CRDs arm: nothing is derived or
	// applied and Apply returns no resources. The release must still skip
	// CRDs; with the CRDs absent the operator cannot start, so this is never
	// a lighter install.
	Install bool
	// KeepOnUninstall true (the default posture) retains every CRD on
	// destroy so an operator uninstall never cascade-deletes the custom
	// resources built on it. False lets destroy delete them: the arm e2e
	// lanes use to leave a cluster clean.
	KeepOnUninstall bool
	// ProviderConfig is the module's Kubernetes connection; nil means the
	// host kubeconfig (the local workflow), exactly as the provider getter
	// treats it.
	ProviderConfig *kubernetesprovider.KubernetesProviderConfig
	// ProviderName names the dedicated upsert provider resource; modules
	// keep it distinct from their plain provider ("kubernetes-crd-upsert").
	ProviderName string
}

// liveReadTimeout bounds the never-downgrade read. A cluster that cannot
// answer a list in this time will not take an apply either.
const liveReadTimeout = 30 * time.Second

// Apply derives, checks, and registers the CRD set, returning the resources
// the release must depend on. Errors are helmcrds.Failure values wherever the
// cause is one the primitive anticipates.
func Apply(ctx *pulumi.Context, args Args) ([]pulumi.Resource, error) {
	if !args.Install {
		return nil, nil
	}
	if args.ProviderName == "" {
		args.ProviderName = "kubernetes-crd-upsert"
	}

	// A destroy re-runs the program only so delete hooks can fire; every
	// resource in state is deleted whatever is registered here. Deriving
	// can fail for reasons that have nothing to do with the teardown (the
	// manifest currently pins a version that was never published, the
	// repository is unreachable, a downgrade is pending), and none of them
	// may stand between a user and deleting their stack: register nothing
	// and let the destroy proceed.
	if stackinput.IsDestroy() {
		return nil, nil
	}

	// The read runs during preview too: a downgrade is refused at plan time
	// in both engines, before anything is touched.
	if err := refuseDowngrade(args); err != nil {
		return nil, err
	}

	crds, err := helmcrds.Derive(context.Background(), args.Source, args.ReleaseName, args.Namespace)
	if err != nil {
		return nil, err
	}

	upsertProvider, err := pulumikubernetesprovider.GetWithKubernetesProviderConfigUpsert(ctx,
		args.ProviderConfig, args.ProviderName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create the CRD upsert kubernetes provider")
	}

	options := []pulumi.ResourceOption{pulumi.Provider(upsertProvider)}
	if args.KeepOnUninstall {
		options = append(options, pulumi.Transformations([]pulumi.ResourceTransformation{
			func(targs *pulumi.ResourceTransformationArgs) *pulumi.ResourceTransformationResult {
				return &pulumi.ResourceTransformationResult{
					Props: targs.Props,
					Opts:  append(targs.Opts, pulumi.RetainOnDelete(true)),
				}
			},
		}))
	}

	// One ConfigGroup per CRD, named by the CRD's own metadata.name (helmcrds
	// sorts them), so state addresses stay stable across chart reorderings
	// and match the Terraform twin's for_each keys.
	resources := make([]pulumi.Resource, 0, len(crds))
	for _, crd := range crds {
		group, err := pulumiyaml.NewConfigGroup(ctx, crd.Name,
			&pulumiyaml.ConfigGroupArgs{YAML: []string{crd.YAML}}, options...)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to apply CRD %s", crd.Name)
		}
		resources = append(resources, group)
	}
	return resources, nil
}

// refuseDowngrade lists the CRDs this source has ever stamped on the cluster
// and refuses when any carries a higher source version than the one about to
// be applied. A cluster with none of them (first install) passes trivially.
func refuseDowngrade(args Args) error {
	existing, err := listStampedCRDs(args)
	if err != nil {
		return err
	}
	return helmcrds.CheckNoDowngrade(existing, args.Source.Version)
}

var crdGVR = schema.GroupVersionResource{
	Group: "apiextensions.k8s.io", Version: "v1", Resource: "customresourcedefinitions",
}

func listStampedCRDs(args Args) ([]helmcrds.ExistingCRD, error) {
	restConfig, err := kubeconfig.RESTConfig(args.ProviderConfig, os.Getenv(execcredential.CommandPathEnvVar))
	if err != nil {
		return nil, errors.Wrap(err, "failed to resolve the cluster connection for the CRD version check")
	}
	client, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create the cluster client for the CRD version check")
	}

	ctx, cancel := context.WithTimeout(context.Background(), liveReadTimeout)
	defer cancel()
	list, err := client.Resource(crdGVR).List(ctx, metav1.ListOptions{
		LabelSelector: helmcrds.LabelSelector(args.Source),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list the CRDs stamped by %s for the version check", args.Source.SourceDescription())
	}
	return toExisting(list.Items), nil
}

func toExisting(items []unstructured.Unstructured) []helmcrds.ExistingCRD {
	existing := make([]helmcrds.ExistingCRD, 0, len(items))
	for _, item := range items {
		existing = append(existing, helmcrds.ExistingCRD{
			Name:    item.GetName(),
			Version: item.GetAnnotations()[helmcrds.AnnotationSourceVersion],
		})
	}
	return existing
}
