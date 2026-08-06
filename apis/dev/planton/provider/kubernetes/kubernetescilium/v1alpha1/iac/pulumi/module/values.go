package module

import (
	"strconv"

	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
)

// buildHelmValues renders the typed spec into the chart's values map, then
// merges the spec's helm_values escape hatch over it with Helm `-f`
// semantics (maps deep-merge with the later document winning, lists
// replace).
//
// No fullnameOverride here (unlike sibling modules): the cilium chart names
// its workloads with FIXED names — DaemonSet "cilium", Deployment
// "cilium-operator" — regardless of the release name, so there is nothing
// to pin.
//
// PARITY: the Terraform module reaches the same result natively — its
// helm_release passes values = [yamlencode(typed values), helm_values] and
// the provider merges the documents in exactly this order. Keep every typed
// mapping below in lockstep with the Terraform module's locals.
func buildHelmValues(locals *Locals) (map[string]interface{}, error) {
	spec := locals.Spec

	values := map[string]interface{}{}

	// ---- cluster identity --------------------------------------------------
	if spec.ClusterName != nil && spec.GetClusterName() != "" {
		values["cluster"] = map[string]interface{}{"name": spec.GetClusterName()}
	}

	// ---- IPAM ---------------------------------------------------------------
	if i := spec.GetIpam(); i != nil {
		ipam := map[string]interface{}{}
		if i.Mode != nil && i.GetMode() != "" {
			ipam["mode"] = i.GetMode()
		}
		operator := map[string]interface{}{}
		if len(i.GetClusterPoolIpv4PodCidrs()) > 0 {
			operator["clusterPoolIPv4PodCIDRList"] = stringSliceToInterface(i.GetClusterPoolIpv4PodCidrs())
		}
		if i.ClusterPoolIpv4MaskSize != nil {
			operator["clusterPoolIPv4MaskSize"] = i.GetClusterPoolIpv4MaskSize()
		}
		if len(operator) > 0 {
			ipam["operator"] = operator
		}
		if len(ipam) > 0 {
			values["ipam"] = ipam
		}
	}

	// ---- routing --------------------------------------------------------------
	if r := spec.GetRouting(); r != nil {
		if r.Mode != nil && r.GetMode() != "" {
			values["routingMode"] = r.GetMode()
		}
		if r.TunnelProtocol != nil && r.GetTunnelProtocol() != "" {
			values["tunnelProtocol"] = r.GetTunnelProtocol()
		}
		if r.GetIpv4NativeRoutingCidr() != "" {
			values["ipv4NativeRoutingCIDR"] = r.GetIpv4NativeRoutingCidr()
		}
		if r.GetAutoDirectNodeRoutes() {
			values["autoDirectNodeRoutes"] = true
		}
	}

	// ---- kube-proxy replacement --------------------------------------------------
	// TRAP: kubeProxyReplacement is a STRING in the chart's values.yaml
	// (historically it took "strict"/"partial"; today "true"/"false") —
	// rendering a YAML boolean would still coerce, but the string keeps the
	// rendered document byte-identical with what the chart declares and
	// with the Terraform module. Only rendered when true (chart default is
	// "false").
	if spec.GetKubeProxyReplacement() {
		values["kubeProxyReplacement"] = "true"
	}

	if spec.GetK8SServiceHost() != "" {
		values["k8sServiceHost"] = spec.GetK8SServiceHost()
	}
	// k8sServicePort is also a string in values.yaml (default ""), so the
	// number renders as its decimal string — the Terraform twin uses
	// tostring() for the same reason.
	if spec.K8SServicePort != nil {
		values["k8sServicePort"] = strconv.Itoa(int(spec.GetK8SServicePort()))
	}

	// ---- CNI installation / chaining ------------------------------------------------
	if c := spec.GetCni(); c != nil {
		cni := map[string]interface{}{}
		if c.ChainingMode != nil && c.GetChainingMode() != "" {
			cni["chainingMode"] = c.GetChainingMode()
		}
		if c.GetChainingTarget() != "" {
			cni["chainingTarget"] = c.GetChainingTarget()
		}
		// Optional bool: presence (not truth) decides rendering — an
		// explicit false is exactly the value chaining setups must send
		// (the CEL rule enforces it), while unset keeps the chart default
		// (true).
		if c.Exclusive != nil {
			cni["exclusive"] = c.GetExclusive()
		}
		if len(cni) > 0 {
			values["cni"] = cni
		}
	}

	// ---- cloud-provider datapath integrations -------------------------------------------
	if cloud := spec.GetCloud(); cloud != nil {
		// AWS ENI datapath — pairs with ipam mode "eni" (pods draw
		// VPC-routable IPs from ENIs Cilium manages).
		if cloud.GetAwsEni() {
			values["eni"] = map[string]interface{}{"enabled": true}
		}
		if cloud.GetAksByocni() {
			values["aksbyocni"] = map[string]interface{}{"enabled": true}
		}
		if cloud.GetGke() {
			values["gke"] = map[string]interface{}{"enabled": true}
		}
	}

	// ---- Hubble observability ------------------------------------------------------------
	if h := spec.GetHubble(); h != nil {
		hubble := map[string]interface{}{}
		// Chart default is enabled=true, so only an EXPLICIT false is
		// rendered (an explicit true is the default — nothing to say).
		if h.Enabled != nil && !h.GetEnabled() {
			hubble["enabled"] = false
		}
		if h.GetRelay() {
			hubble["relay"] = map[string]interface{}{"enabled": true}
		}
		if h.GetUi() {
			hubble["ui"] = map[string]interface{}{"enabled": true}
		}
		// hubble.metrics.enabled is upstream's LIST of metric families
		// (null disables) — not a boolean despite the name.
		if len(h.GetMetrics()) > 0 {
			metrics := map[string]interface{}{
				"enabled": stringSliceToInterface(h.GetMetrics()),
			}
			if h.GetMetricsServiceMonitor() {
				metrics["serviceMonitor"] = map[string]interface{}{"enabled": true}
			}
			hubble["metrics"] = metrics
		}
		if len(hubble) > 0 {
			values["hubble"] = hubble
		}
	}

	// ---- transparent encryption --------------------------------------------------------------
	if e := spec.GetEncryption(); e != nil && e.GetEnabled() {
		encryption := map[string]interface{}{"enabled": true}
		if e.Type != nil && e.GetType() != "" {
			encryption["type"] = e.GetType()
		}
		if e.GetNodeEncryption() {
			encryption["nodeEncryption"] = true
		}
		values["encryption"] = encryption
	}

	// ---- policy enforcement ---------------------------------------------------------------------
	if spec.PolicyEnforcementMode != nil && spec.GetPolicyEnforcementMode() != "" {
		values["policyEnforcementMode"] = spec.GetPolicyEnforcementMode()
	}

	// ---- Gateway API -------------------------------------------------------------------------------
	if spec.GetGatewayApi() {
		values["gatewayAPI"] = map[string]interface{}{"enabled": true}
	}

	// ---- bandwidth manager ----------------------------------------------------------------------------
	if b := spec.GetBandwidthManager(); b != nil && b.GetEnabled() {
		bandwidthManager := map[string]interface{}{"enabled": true}
		if b.GetBbr() {
			bandwidthManager["bbr"] = true
		}
		values["bandwidthManager"] = bandwidthManager
	}

	// ---- operator sizing + agent resources ------------------------------------------------------------
	// The operator map collects keys from TWO spec arms (operator sizing
	// here, operator.prometheus below) — build it once so the arms merge
	// into ONE map instead of the later overwriting the earlier.
	operatorValues := map[string]interface{}{}
	if o := spec.GetOperator(); o != nil {
		if o.Replicas != nil {
			operatorValues["replicas"] = o.GetReplicas()
		}
		if r := resourcesMap(o.GetResources()); r != nil {
			operatorValues["resources"] = r
		}
	}

	// Top-level resources = the agent container (the cilium DaemonSet).
	if r := resourcesMap(spec.GetAgentResources()); r != nil {
		values["resources"] = r
	}

	// ---- Cilium's own telemetry -------------------------------------------------------------------------
	// One spec toggle drives BOTH components: agent metrics (top-level
	// prometheus) and operator metrics (operator.prometheus) — exposing
	// only one of the two would be a confusing half-telemetry posture.
	if p := spec.GetPrometheus(); p.GetEnabled() {
		prometheus := map[string]interface{}{"enabled": true}
		operatorPrometheus := map[string]interface{}{"enabled": true}
		if p.GetServiceMonitor() {
			prometheus["serviceMonitor"] = map[string]interface{}{"enabled": true}
			operatorPrometheus["serviceMonitor"] = map[string]interface{}{"enabled": true}
		}
		values["prometheus"] = prometheus
		operatorValues["prometheus"] = operatorPrometheus
	}

	if len(operatorValues) > 0 {
		values["operator"] = operatorValues
	}

	// ---- escape hatch (merged LAST, helm -f semantics) ------------------------------
	if spec.GetHelmValues() != "" {
		overrides := map[string]interface{}{}
		if err := yaml.Unmarshal([]byte(spec.GetHelmValues()), &overrides); err != nil {
			return nil, errors.Wrap(err, "failed to parse helm_values as a YAML document")
		}
		values = mergeMaps(values, overrides)
	}

	return values, nil
}
