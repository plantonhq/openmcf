package module

import (
	"github.com/pkg/errors"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	kubernetesmeta "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// service fronts one role's Deployment with a ClusterIP Service. The
// selector carries the role-specific "app" label so the registry Service
// never routes to REST-proxy pods (same image, same namespace) and vice
// versa.
func service(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	createdDeployment *appsv1.Deployment,
	name string, port int, selectorLabels map[string]string,
) error {
	labels := mergedLabels(locals, selectorLabels)

	_, err := kubernetescorev1.NewService(ctx,
		name,
		&kubernetescorev1.ServiceArgs{
			Metadata: &kubernetesmeta.ObjectMetaArgs{
				Name:      pulumi.String(name),
				Namespace: pulumi.String(locals.Namespace),
				Labels:    pulumi.ToStringMap(labels),
			},
			Spec: &kubernetescorev1.ServiceSpecArgs{
				Type:     pulumi.String("ClusterIP"),
				Selector: pulumi.ToStringMap(selectorLabels),
				Ports: kubernetescorev1.ServicePortArray{
					&kubernetescorev1.ServicePortArgs{
						Name:       pulumi.String("http"),
						Port:       pulumi.Int(port),
						TargetPort: pulumi.Int(port),
					},
				},
			},
		},
		pulumi.Provider(kubernetesProvider),
		pulumi.DependsOn([]pulumi.Resource{createdDeployment}))
	if err != nil {
		return errors.Wrapf(err, "failed to create %s service", name)
	}
	return nil
}
