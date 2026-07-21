package module

import (
	"github.com/pkg/errors"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	kubernetesmetav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// service creates the headless governing Service of the StatefulSet. Headless
// (clusterIP: None) is what gives each replica its stable per-pod DNS name
// (<pod>.<service>.<namespace>.svc.cluster.local) — a regular ClusterIP would
// only load-balance and could never address individual members. It must exist
// before the StatefulSet: the API requires spec.serviceName to reference an
// existing Service, and pods resolve peer DNS through it at startup.
func service(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider pulumi.ProviderResource, namespaceDeps []pulumi.ResourceOption) (*kubernetescorev1.Service, error) {

	appPorts := locals.KubernetesStatefulSet.Spec.Container.App.Ports

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

	serviceSpecArgs := &kubernetescorev1.ServiceSpecArgs{
		// "None" is the headless marker — no virtual IP, DNS resolves straight
		// to pod IPs, and each pod gets its own stable DNS record.
		ClusterIP: pulumi.String("None"),
		Selector:  pulumi.ToStringMap(locals.SelectorLabels),
		// Peers must discover each other BEFORE they can pass readiness —
		// a bootstrapping member needs DNS for pods that are themselves still
		// bootstrapping. Publishing not-ready addresses breaks that deadlock.
		PublishNotReadyAddresses: pulumi.Bool(true),
	}
	if len(portsArray) > 0 {
		serviceSpecArgs.Ports = portsArray
	}

	serviceArgs := &kubernetescorev1.ServiceArgs{
		Metadata: kubernetesmetav1.ObjectMetaArgs{
			Name:      pulumi.String(locals.KubeServiceName),
			Namespace: pulumi.String(locals.Namespace),
			Labels:    pulumi.ToStringMap(locals.Labels),
		},
		Spec: serviceSpecArgs,
	}

	svcOpts := append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, namespaceDeps...)
	createdService, err := kubernetescorev1.NewService(ctx,
		locals.KubeServiceName,
		serviceArgs,
		svcOpts...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create headless service")
	}
	return createdService, nil
}
