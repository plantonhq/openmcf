package module

import (
	"github.com/pkg/errors"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	kubernetesmetav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// service fronts the Deployment's pods with a ClusterIP Service, one Service
// port per app-container port. No Service is created when the app container
// exposes no ports (workers, consumers).
func service(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider pulumi.ProviderResource, createdDeployment *appsv1.Deployment, namespaceDeps []pulumi.ResourceOption) error {

	appPorts := locals.KubernetesDeployment.Spec.Container.App.Ports
	if len(appPorts) == 0 {
		return nil
	}

	portsArray := make(kubernetescorev1.ServicePortArray, 0, len(appPorts))
	for _, p := range appPorts {
		// service_port defaults to the container port when unset, so the common
		// "expose as-is" case needs no extra configuration.
		servicePort := p.ServicePort
		if servicePort == 0 {
			servicePort = p.ContainerPort
		}

		portArgs := &kubernetescorev1.ServicePortArgs{
			Name:       pulumi.String(p.Name),
			Port:       pulumi.Int(servicePort),
			TargetPort: pulumi.Int(p.ContainerPort),
		}
		if p.NetworkProtocol != "" {
			portArgs.Protocol = pulumi.String(p.NetworkProtocol)
		}
		if p.AppProtocol != "" {
			portArgs.AppProtocol = pulumi.String(p.AppProtocol)
		}
		portsArray = append(portsArray, portArgs)
	}

	serviceArgs := &kubernetescorev1.ServiceArgs{
		Metadata: kubernetesmetav1.ObjectMetaArgs{
			Name:      pulumi.String(locals.KubeServiceName),
			Namespace: pulumi.String(locals.Namespace),
			Labels:    pulumi.ToStringMap(locals.Labels),
		},
		Spec: &kubernetescorev1.ServiceSpecArgs{
			Type:     pulumi.String("ClusterIP"),
			Selector: pulumi.ToStringMap(locals.SelectorLabels),
			Ports:    portsArray,
		},
	}

	svcOpts := append([]pulumi.ResourceOption{
		pulumi.Provider(kubernetesProvider),
		pulumi.DependsOn([]pulumi.Resource{createdDeployment}),
	}, namespaceDeps...)
	_, err := kubernetescorev1.NewService(ctx,
		locals.KubeServiceName,
		serviceArgs,
		svcOpts...)
	if err != nil {
		return errors.Wrap(err, "failed to create service")
	}
	return nil
}
