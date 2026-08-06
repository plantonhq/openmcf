package module

import (
	"github.com/pkg/errors"
	kubernetespersistentvolumeclaimv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetespersistentvolumeclaim/v1alpha1"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// createPersistentVolumeClaim creates the core/v1 PersistentVolumeClaim.
//
// The claim is created WITHOUT awaiting Bound (skipAwait annotation, set in
// locals): under a WaitForFirstConsumer StorageClass a claim is correctly
// Pending until a pod consumes it, and awaiting would hang every such deploy.
// The Terraform module's wait_until_bound=false is the same decision.
func createPersistentVolumeClaim(ctx *pulumi.Context, locals *Locals, provider pulumi.ProviderResource) (*kubernetescorev1.PersistentVolumeClaim, error) {
	spec := locals.Spec

	resources := &kubernetescorev1.VolumeResourceRequirementsArgs{
		Requests: pulumi.ToStringMap(map[string]string{"storage": spec.GetStorageRequest()}),
	}
	if spec.GetStorageLimit() != "" {
		resources.Limits = pulumi.ToStringMap(map[string]string{"storage": spec.GetStorageLimit()})
	}

	claimSpecArgs := &kubernetescorev1.PersistentVolumeClaimSpecArgs{
		AccessModes: pulumi.ToStringArray(locals.AccessModes),
		Resources:   resources,
		// Filesystem is the server default, but sending it explicitly keeps
		// both engines' submitted objects identical.
		VolumeMode: pulumi.String(locals.VolumeMode),
	}

	// nil = absent (cluster default class); "" = explicitly no dynamic
	// provisioning. The distinction is why this is a pointer.
	if locals.StorageClassName != nil {
		claimSpecArgs.StorageClassName = pulumi.String(*locals.StorageClassName)
	}

	if spec.GetVolumeName() != "" {
		claimSpecArgs.VolumeName = pulumi.String(spec.GetVolumeName())
	}

	if selector := spec.GetSelector(); selector != nil {
		claimSpecArgs.Selector = buildLabelSelector(selector)
	}

	// PARITY-EXCEPTION: the terraform kubernetes provider's PVC resource
	// cannot express spec.dataSource/dataSourceRef (clone a PVC / restore a
	// VolumeSnapshot); its module fails the plan with a precondition when the
	// field is set. This engine sends the data source natively.
	if dataSource := spec.GetDataSource(); dataSource != nil {
		claimSpecArgs.DataSource = buildDataSource(dataSource)
	}

	persistentVolumeClaim, err := kubernetescorev1.NewPersistentVolumeClaim(
		ctx,
		locals.Name,
		&kubernetescorev1.PersistentVolumeClaimArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name:        pulumi.String(locals.Name),
				Namespace:   pulumi.String(locals.Namespace),
				Labels:      pulumi.ToStringMap(locals.Labels),
				Annotations: pulumi.ToStringMap(locals.Annotations),
			},
			Spec: claimSpecArgs,
		},
		pulumi.Provider(provider),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create persistent volume claim %s/%s", locals.Namespace, locals.Name)
	}

	return persistentVolumeClaim, nil
}

// buildDataSource maps the typed data-source variants onto the API's
// TypedLocalObjectReference: PVC clones live in the core group (apiGroup
// omitted); VolumeSnapshot restores name the snapshot.storage.k8s.io group.
func buildDataSource(ds *kubernetespersistentvolumeclaimv1alpha1.KubernetesPersistentVolumeClaimDataSource) *kubernetescorev1.TypedLocalObjectReferenceArgs {
	args := &kubernetescorev1.TypedLocalObjectReferenceArgs{
		Name: pulumi.String(ds.GetName()),
	}
	if ds.GetKind() == kubernetespersistentvolumeclaimv1alpha1.KubernetesPersistentVolumeClaimDataSource_volume_snapshot {
		args.Kind = pulumi.String("VolumeSnapshot")
		args.ApiGroup = pulumi.String("snapshot.storage.k8s.io")
	} else {
		args.Kind = pulumi.String("PersistentVolumeClaim")
	}
	return args
}

// buildLabelSelector converts the proto label selector into Pulumi args.
func buildLabelSelector(s *kubernetespersistentvolumeclaimv1alpha1.KubernetesPersistentVolumeClaimLabelSelector) *metav1.LabelSelectorArgs {
	selectorArgs := &metav1.LabelSelectorArgs{}
	if len(s.GetMatchLabels()) > 0 {
		selectorArgs.MatchLabels = pulumi.ToStringMap(s.GetMatchLabels())
	}
	if len(s.GetMatchExpressions()) > 0 {
		var exprArray metav1.LabelSelectorRequirementArray
		for _, e := range s.GetMatchExpressions() {
			exprArgs := &metav1.LabelSelectorRequirementArgs{
				Key:      pulumi.String(e.GetKey()),
				Operator: pulumi.String(e.GetOperator()),
			}
			if len(e.GetValues()) > 0 {
				exprArgs.Values = pulumi.ToStringArray(e.GetValues())
			}
			exprArray = append(exprArray, exprArgs)
		}
		selectorArgs.MatchExpressions = exprArray
	}
	return selectorArgs
}
