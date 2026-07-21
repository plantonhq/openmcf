package module

import (
	"strconv"

	"github.com/pkg/errors"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// createService creates the Kubernetes Service resource.
//
// Optional spec fields are set on the API object only when the user set them —
// sending an empty value is NOT the same as omitting the field for several of
// these (clusterIP, healthCheckNodePort, loadBalancerClass are immutable or
// type-gated, and the API rejects them on the wrong service type), so each one
// is guarded by the same condition the Terraform module uses.
func createService(ctx *pulumi.Context, locals *Locals, provider pulumi.ProviderResource) (*kubernetescorev1.Service, error) {
	spec := locals.Spec
	isExternalName := locals.ServiceType == "ExternalName"
	isLoadBalancer := locals.ServiceType == "LoadBalancer"
	isExternallyReachable := locals.ServiceType == "NodePort" || isLoadBalancer

	serviceSpecArgs := &kubernetescorev1.ServiceSpecArgs{
		Type: pulumi.String(locals.ServiceType),
	}

	// Ports and selector do not apply to ExternalName (a pure DNS alias).
	if !isExternalName {
		serviceSpecArgs.Ports = buildServicePorts(locals)
		if len(spec.GetSelector()) > 0 {
			serviceSpecArgs.Selector = pulumi.ToStringMap(spec.GetSelector())
		}
	}

	// Headless IS clusterIP "None"; otherwise honor an explicitly requested
	// static cluster IP. Both are create-time-only (the field is immutable).
	if spec.GetHeadless() {
		serviceSpecArgs.ClusterIP = pulumi.String("None")
	} else if spec.GetClusterIpAddress() != "" {
		serviceSpecArgs.ClusterIP = pulumi.String(spec.GetClusterIpAddress())
	}

	if isExternalName {
		serviceSpecArgs.ExternalName = pulumi.String(spec.GetExternalDnsName())
	}

	if len(spec.GetExternalIps()) > 0 {
		serviceSpecArgs.ExternalIPs = pulumi.ToStringArray(spec.GetExternalIps())
	}

	// Traffic policies are type-gated by the API: external only for
	// externally-reachable types, internal never for ExternalName.
	if isExternallyReachable {
		serviceSpecArgs.ExternalTrafficPolicy = pulumi.String(locals.ExternalTrafficPolicy)
	}
	if !isExternalName && spec.InternalTrafficPolicy != nil {
		serviceSpecArgs.InternalTrafficPolicy = pulumi.String(locals.InternalTrafficPolicy)
	}

	// trafficDistribution is a hint; only send it when the user chose one.
	// PARITY-EXCEPTION: the Terraform kubernetes provider (v3.2.x) does not
	// expose spec.trafficDistribution, so only the Pulumi engine can apply this
	// field; the Terraform module fails the plan loudly via a precondition when
	// it is set instead of silently dropping it.
	if locals.TrafficDistribution != "" {
		serviceSpecArgs.TrafficDistribution = pulumi.String(locals.TrafficDistribution)
	}

	serviceSpecArgs.SessionAffinity = pulumi.String(locals.SessionAffinity)
	if locals.SessionAffinity == "ClientIP" && spec.SessionAffinityTimeoutSeconds != nil {
		serviceSpecArgs.SessionAffinityConfig = &kubernetescorev1.SessionAffinityConfigArgs{
			ClientIP: &kubernetescorev1.ClientIPConfigArgs{
				TimeoutSeconds: pulumi.Int(int(spec.GetSessionAffinityTimeoutSeconds())),
			},
		}
	}

	// LoadBalancer-only knobs — the CEL rules guarantee they are unset for
	// other types, and the API would reject them there anyway.
	if isLoadBalancer {
		if len(spec.GetLoadBalancerSourceRanges()) > 0 {
			serviceSpecArgs.LoadBalancerSourceRanges = pulumi.ToStringArray(spec.GetLoadBalancerSourceRanges())
		}
		if spec.GetLoadBalancerClass() != "" {
			serviceSpecArgs.LoadBalancerClass = pulumi.String(spec.GetLoadBalancerClass())
		}
		if spec.AllocateLoadBalancerNodePorts != nil {
			serviceSpecArgs.AllocateLoadBalancerNodePorts = pulumi.Bool(spec.GetAllocateLoadBalancerNodePorts())
		}
		if spec.GetHealthCheckNodePort() != 0 {
			serviceSpecArgs.HealthCheckNodePort = pulumi.Int(int(spec.GetHealthCheckNodePort()))
		}
	}

	if spec.GetPublishNotReadyAddresses() {
		serviceSpecArgs.PublishNotReadyAddresses = pulumi.Bool(true)
	}

	// Dual-stack: families and policy are only sent when requested; the cluster
	// otherwise assigns from its own configuration.
	if len(locals.IpFamilies) > 0 {
		serviceSpecArgs.IpFamilies = pulumi.ToStringArray(locals.IpFamilies)
	}
	if locals.IpFamilyPolicy != "" {
		serviceSpecArgs.IpFamilyPolicy = pulumi.String(locals.IpFamilyPolicy)
	}

	service, err := kubernetescorev1.NewService(
		ctx,
		locals.Name,
		&kubernetescorev1.ServiceArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name:        pulumi.String(locals.Name),
				Namespace:   pulumi.String(locals.Namespace),
				Labels:      pulumi.ToStringMap(locals.Labels),
				Annotations: pulumi.ToStringMap(locals.Annotations),
			},
			Spec: serviceSpecArgs,
		},
		pulumi.Provider(provider),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create service %s/%s", locals.Namespace, locals.Name)
	}

	return service, nil
}

// buildServicePorts converts the protobuf port definitions to Pulumi service port arguments.
func buildServicePorts(locals *Locals) kubernetescorev1.ServicePortArray {
	var ports kubernetescorev1.ServicePortArray

	for _, p := range locals.Spec.GetPorts() {
		portArgs := &kubernetescorev1.ServicePortArgs{
			Port:     pulumi.Int(int(p.GetPort())),
			Protocol: pulumi.String(resolveProtocol(p.GetProtocol())),
		}

		if p.GetName() != "" {
			portArgs.Name = pulumi.String(p.GetName())
		}

		if p.GetAppProtocol() != "" {
			portArgs.AppProtocol = pulumi.String(p.GetAppProtocol())
		}

		// target_port is an IntOrString upstream: a numeric string routes to a
		// port number, anything else is a named container port. Omitted means
		// "same as port" (the API's identity mapping).
		if p.GetTargetPort() != "" {
			if num, err := strconv.Atoi(p.GetTargetPort()); err == nil {
				portArgs.TargetPort = pulumi.Int(num)
			} else {
				portArgs.TargetPort = pulumi.String(p.GetTargetPort())
			}
		}

		if p.GetNodePort() > 0 {
			portArgs.NodePort = pulumi.Int(int(p.GetNodePort()))
		}

		ports = append(ports, portArgs)
	}

	return ports
}
