package module

import (
	"github.com/pkg/errors"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	kubernetesmeta "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// watchNamespaces creates every namespace listed in
// spec.watch_namespaces. The chart plants job-identity RBAC INTO those
// namespaces and does not create them — they must exist before the
// Helm release (verified live: install fails with
// "namespaces \"…\" not found"). Twin of the Terraform module's
// kubernetes_namespace_v1.watch for_each.
func watchNamespaces(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
) ([]pulumi.Resource, error) {
	watch := locals.Spec.GetWatchNamespaces()
	if len(watch) == 0 {
		return nil, nil
	}
	created := make([]pulumi.Resource, 0, len(watch))
	for _, ns := range watch {
		res, err := kubernetescorev1.NewNamespace(ctx,
			"watch-"+ns,
			&kubernetescorev1.NamespaceArgs{
				Metadata: kubernetesmeta.ObjectMetaPtrInput(
					&kubernetesmeta.ObjectMetaArgs{
						Name:   pulumi.String(ns),
						Labels: pulumi.ToStringMap(locals.Labels),
					}),
			}, pulumi.Provider(kubernetesProvider))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create watch namespace %q", ns)
		}
		created = append(created, res)
	}
	return created, nil
}
