package module

import (
	"github.com/pkg/errors"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	kubernetesschedulingv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/scheduling/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// createPriorityClass creates the scheduling.k8s.io/v1 PriorityClass.
//
// preemptionPolicy is ALWAYS sent explicitly (with the API server's default
// applied module-side) so both engines submit byte-identical objects. value
// is immutable upstream — a change forces replacement, and DeleteBeforeReplace
// avoids the name collision (PriorityClass names are cluster-unique).
func createPriorityClass(ctx *pulumi.Context, locals *Locals, provider pulumi.ProviderResource) (*kubernetesschedulingv1.PriorityClass, error) {
	spec := locals.Spec

	priorityClass, err := kubernetesschedulingv1.NewPriorityClass(
		ctx,
		locals.Name,
		&kubernetesschedulingv1.PriorityClassArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name:        pulumi.String(locals.Name),
				Labels:      pulumi.ToStringMap(locals.Labels),
				Annotations: pulumi.ToStringMap(locals.Annotations),
			},
			Value: pulumi.Int(int(spec.GetValue())),
			// The Kubernetes default is false; sending it explicitly keeps
			// both engines' submitted objects identical.
			GlobalDefault:    pulumi.Bool(spec.GetGlobalDefault()),
			Description:      pulumi.String(spec.GetDescription()),
			PreemptionPolicy: pulumi.String(locals.PreemptionPolicy),
		},
		pulumi.Provider(provider),
		pulumi.DeleteBeforeReplace(true),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create priority class %s", locals.Name)
	}

	return priorityClass, nil
}
