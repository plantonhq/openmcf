package module

import (
	"strings"

	alicloudnetworkloadbalancerv1alpha1 "github.com/plantonhq/planton/catalog/alicloud/alicloudnetworkloadbalancer/v1alpha1"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AliCloudNetworkLoadBalancer *alicloudnetworkloadbalancerv1alpha1.AliCloudNetworkLoadBalancer
	Tags                        map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *alicloudnetworkloadbalancerv1alpha1.AliCloudNetworkLoadBalancerStackInput) *Locals {
	locals := &Locals{}
	locals.AliCloudNetworkLoadBalancer = stackInput.Target
	target := stackInput.Target

	locals.Tags = map[string]string{
		"resource":      "true",
		"resource_name": target.Metadata.Name,
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AliCloudNetworkLoadBalancer.String()),
	}

	if target.Metadata.Id != "" {
		locals.Tags["resource_id"] = target.Metadata.Id
	}

	if target.Metadata.Org != "" {
		locals.Tags["organization"] = target.Metadata.Org
	}

	if target.Metadata.Env != "" {
		locals.Tags["environment"] = target.Metadata.Env
	}

	for k, v := range target.Spec.Tags {
		locals.Tags[k] = v
	}

	return locals
}

func addressType(spec *alicloudnetworkloadbalancerv1alpha1.AliCloudNetworkLoadBalancerSpec) string {
	if spec.AddressType != nil {
		return *spec.AddressType
	}
	return "Internet"
}

func crossZoneEnabled(spec *alicloudnetworkloadbalancerv1alpha1.AliCloudNetworkLoadBalancerSpec) bool {
	if spec.CrossZoneEnabled != nil {
		return *spec.CrossZoneEnabled
	}
	return true
}

func serverGroupProtocol(sg *alicloudnetworkloadbalancerv1alpha1.AliCloudNetworkLoadBalancerServerGroup) string {
	if sg.Protocol != nil {
		return *sg.Protocol
	}
	return "TCP"
}

func serverGroupScheduler(sg *alicloudnetworkloadbalancerv1alpha1.AliCloudNetworkLoadBalancerServerGroup) string {
	if sg.Scheduler != nil {
		return *sg.Scheduler
	}
	return "Wrr"
}

func connectionDrainEnabled(sg *alicloudnetworkloadbalancerv1alpha1.AliCloudNetworkLoadBalancerServerGroup) bool {
	if sg.ConnectionDrainEnabled != nil {
		return *sg.ConnectionDrainEnabled
	}
	return false
}

func connectionDrainTimeout(sg *alicloudnetworkloadbalancerv1alpha1.AliCloudNetworkLoadBalancerServerGroup) int {
	if sg.ConnectionDrainTimeout != nil {
		return int(*sg.ConnectionDrainTimeout)
	}
	return 10
}

func preserveClientIpEnabled(sg *alicloudnetworkloadbalancerv1alpha1.AliCloudNetworkLoadBalancerServerGroup) bool {
	if sg.PreserveClientIpEnabled != nil {
		return *sg.PreserveClientIpEnabled
	}
	return true
}

func healthCheckType(hc *alicloudnetworkloadbalancerv1alpha1.AliCloudNetworkLoadBalancerHealthCheckConfig) string {
	if hc.HealthCheckType != nil {
		return *hc.HealthCheckType
	}
	return "TCP"
}

func healthCheckConnectPort(hc *alicloudnetworkloadbalancerv1alpha1.AliCloudNetworkLoadBalancerHealthCheckConfig) int {
	if hc.HealthCheckConnectPort != nil {
		return int(*hc.HealthCheckConnectPort)
	}
	return 0
}

func healthCheckConnectTimeout(hc *alicloudnetworkloadbalancerv1alpha1.AliCloudNetworkLoadBalancerHealthCheckConfig) int {
	if hc.HealthCheckConnectTimeout != nil {
		return int(*hc.HealthCheckConnectTimeout)
	}
	return 5
}

func healthCheckInterval(hc *alicloudnetworkloadbalancerv1alpha1.AliCloudNetworkLoadBalancerHealthCheckConfig) int {
	if hc.HealthCheckInterval != nil {
		return int(*hc.HealthCheckInterval)
	}
	return 10
}

func healthyThreshold(hc *alicloudnetworkloadbalancerv1alpha1.AliCloudNetworkLoadBalancerHealthCheckConfig) int {
	if hc.HealthyThreshold != nil {
		return int(*hc.HealthyThreshold)
	}
	return 2
}

func unhealthyThreshold(hc *alicloudnetworkloadbalancerv1alpha1.AliCloudNetworkLoadBalancerHealthCheckConfig) int {
	if hc.UnhealthyThreshold != nil {
		return int(*hc.UnhealthyThreshold)
	}
	return 2
}

func httpCheckMethod(hc *alicloudnetworkloadbalancerv1alpha1.AliCloudNetworkLoadBalancerHealthCheckConfig) string {
	if hc.HttpCheckMethod != nil {
		return *hc.HttpCheckMethod
	}
	return "GET"
}

func listenerIdleTimeout(l *alicloudnetworkloadbalancerv1alpha1.AliCloudNetworkLoadBalancerListener) int {
	if l.IdleTimeout != nil {
		return int(*l.IdleTimeout)
	}
	return 900
}

func listenerProxyProtocolEnabled(l *alicloudnetworkloadbalancerv1alpha1.AliCloudNetworkLoadBalancerListener) bool {
	if l.ProxyProtocolEnabled != nil {
		return *l.ProxyProtocolEnabled
	}
	return false
}

func listenerCaEnabled(l *alicloudnetworkloadbalancerv1alpha1.AliCloudNetworkLoadBalancerListener) bool {
	if l.CaEnabled != nil {
		return *l.CaEnabled
	}
	return false
}
