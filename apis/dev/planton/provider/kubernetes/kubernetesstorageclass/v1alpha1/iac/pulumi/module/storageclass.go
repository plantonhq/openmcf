package module

import (
	"github.com/pkg/errors"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	kubernetesstoragev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/storage/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// createStorageClass creates the storage.k8s.io/v1 StorageClass resource.
//
// reclaimPolicy and volumeBindingMode are ALWAYS sent explicitly (with the
// API server's defaults applied module-side) so both engines submit
// byte-identical objects for the same manifest. provisioner and parameters
// are immutable upstream — Kubernetes rejects in-place changes, and both
// engines replace the class instead.
func createStorageClass(ctx *pulumi.Context, locals *Locals, provider pulumi.ProviderResource) (*kubernetesstoragev1.StorageClass, error) {
	spec := locals.Spec

	storageClassArgs := &kubernetesstoragev1.StorageClassArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:        pulumi.String(locals.Name),
			Labels:      pulumi.ToStringMap(locals.Labels),
			Annotations: pulumi.ToStringMap(locals.Annotations),
		},
		Provisioner:       pulumi.String(spec.GetProvisioner()),
		ReclaimPolicy:     pulumi.String(locals.ReclaimPolicy),
		VolumeBindingMode: pulumi.String(locals.VolumeBindingMode),
		// Sent unconditionally: the Kubernetes default is false, so an
		// explicit false is identical to omission — and the explicit form
		// keeps the two engines' submitted objects identical.
		AllowVolumeExpansion: pulumi.Bool(spec.GetAllowVolumeExpansion()),
	}

	if len(spec.GetParameters()) > 0 {
		storageClassArgs.Parameters = pulumi.ToStringMap(spec.GetParameters())
	}

	if len(spec.GetMountOptions()) > 0 {
		storageClassArgs.MountOptions = pulumi.ToStringArray(spec.GetMountOptions())
	}

	// PARITY-EXCEPTION: the terraform kubernetes provider models
	// allowed_topologies as a SINGLE selector term (max_items = 1); its
	// module fails the plan with a precondition when the spec lists several.
	// This engine sends every term (the API ORs them).
	if len(spec.GetAllowedTopologies()) > 0 {
		// The topology selector types are core/v1 types in the SDK even
		// though the StorageClass itself is storage.k8s.io/v1.
		topologyArray := kubernetescorev1.TopologySelectorTermArray{}
		for _, term := range spec.GetAllowedTopologies() {
			requirementArray := kubernetescorev1.TopologySelectorLabelRequirementArray{}
			for _, req := range term.GetMatchLabelExpressions() {
				requirementArray = append(requirementArray, &kubernetescorev1.TopologySelectorLabelRequirementArgs{
					Key:    pulumi.String(req.GetKey()),
					Values: pulumi.ToStringArray(req.GetValues()),
				})
			}
			topologyArray = append(topologyArray, &kubernetescorev1.TopologySelectorTermArgs{
				MatchLabelExpressions: requirementArray,
			})
		}
		storageClassArgs.AllowedTopologies = topologyArray
	}

	storageClass, err := kubernetesstoragev1.NewStorageClass(
		ctx,
		locals.Name,
		storageClassArgs,
		pulumi.Provider(provider),
		// provisioner/parameters/reclaimPolicy/volumeBindingMode are immutable
		// upstream; deleteBeforeReplace avoids a name collision when a change
		// forces replacement (StorageClass names are cluster-unique).
		pulumi.DeleteBeforeReplace(true),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create storage class %s", locals.Name)
	}

	return storageClass, nil
}
