package module

import (
	"strconv"

	"github.com/pkg/errors"
	kubernetesclusterautoscalerv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetesclusterautoscaler/v1alpha1"
	"sigs.k8s.io/yaml"
)

// buildHelmValues renders the typed spec into the chart's values map, then
// merges the spec's helm_values escape hatch over it with Helm `-f`
// semantics (maps deep-merge with the later document winning, lists
// replace).
//
// PARITY: the Terraform module reaches the same result natively — its
// helm_release passes values = [yamlencode(typed values), helm_values] and
// the provider merges the documents in exactly this order. Keep every typed
// mapping below in lockstep with the Terraform module's locals.
//
// CHART GATE (verified in templates/deployment.yaml line 1): the Deployment
// only renders when autoDiscovery.clusterName / autoDiscovery.namespace /
// autoDiscovery.labels or autoscalingGroups is set. autoscalingGroupsnamePrefix
// (the GCE contract) does NOT satisfy the gate — the chart README explicitly
// requires autoDiscovery.clusterName ("any-name"; the value is unused by the
// gce provider blocks) for GCE, and kwok has no typed gate key at all. Both
// arms therefore render autoDiscovery.clusterName = metadata.name — a
// benign, deterministic gate value (no gce/kwok template consumes it beyond
// the gate). Without it the release would "succeed" while installing NO
// autoscaler pod.
func buildHelmValues(locals *Locals) (map[string]interface{}, error) {
	spec := locals.Spec

	values := map[string]interface{}{
		// Selects the chart's per-provider command/env blocks in
		// deployment.yaml AND which keys the chart's credential Secret
		// carries (templates/secret.yaml).
		"cloudProvider": locals.CloudProvider,
	}

	// Sub-maps assembled across arms and attached at the end when
	// non-empty.
	autoDiscovery := map[string]interface{}{}
	rbac := map[string]interface{}{}
	serviceAccountAnnotations := map[string]interface{}{}

	// ---- provider arm (exactly one; proto oneof) --------------------------
	switch {
	case spec.GetAws() != nil:
		aws := spec.GetAws()
		// awsRegion also gates the chart's AWS env block (deployment.yaml
		// renders AWS_REGION and the key envs only when awsRegion != "").
		values["awsRegion"] = aws.GetRegion()
		if ad := aws.GetAutoDiscovery(); ad != nil {
			// Tag-based discovery: the chart renders
			// --node-group-auto-discovery=asg:tag=<tags> from
			// autoDiscovery.tags, whose chart default is the standard
			// k8s.io/cluster-autoscaler/enabled +
			// k8s.io/cluster-autoscaler/<clusterName> pair — tags render
			// only on explicit override.
			autoDiscovery["clusterName"] = ad.GetClusterName()
			if len(ad.GetTags()) > 0 {
				autoDiscovery["tags"] = stringSliceToInterface(ad.GetTags())
			}
		} else {
			// Static ASGs: each entry renders one
			// --nodes=<min>:<max>:<name> flag.
			values["autoscalingGroups"] = nodeGroupsSlice(aws.GetNodeGroups())
		}
		if aws.GetIrsaRoleArn() != "" {
			// The chart forwards service-account annotations verbatim
			// (templates/serviceaccount.yaml) — the EKS webhook picks the
			// role up from this well-known key.
			serviceAccountAnnotations["eks.amazonaws.com/role-arn"] = aws.GetIrsaRoleArn()
		}
		if ak := aws.GetAccessKeys(); ak != nil {
			// The chart materializes its own Secret from these two values
			// (templates/secret.yaml: AwsAccessKeyId/AwsSecretAccessKey,
			// rendered only when BOTH are set) and wires the env vars via
			// secretKeyRef — the secret key never lands in the pod spec.
			values["awsAccessKeyID"] = ak.GetAccessKeyId()
			values["awsSecretAccessKey"] = ak.GetSecretAccessKey()
		}

	case spec.GetAzure() != nil:
		azure := spec.GetAzure()
		// Both land in the chart's credential Secret and reach the pod as
		// ARM_SUBSCRIPTION_ID / ARM_RESOURCE_GROUP via secretKeyRef
		// (deployment.yaml azure env block).
		values["azureSubscriptionID"] = azure.GetSubscriptionId()
		values["azureResourceGroup"] = azure.GetResourceGroup()
		if azure.GetClusterName() != "" {
			// Renders --node-group-auto-discovery=label:...,
			// cluster-autoscaler-name=<clusterName> (deployment.yaml
			// azure block).
			autoDiscovery["clusterName"] = azure.GetClusterName()
		} else {
			values["autoscalingGroups"] = nodeGroupsSlice(azure.GetNodeGroups())
		}
		identity := azure.GetIdentity()
		switch {
		case identity.GetUseWorkloadIdentity():
			// Sets ARM_USE_WORKLOAD_IDENTITY_EXTENSION=true. VERIFIED: no
			// template adds the azure.workload.identity/use pod label —
			// clusters relying on the azure-workload-identity webhook add
			// podLabels via helm_values (the chart README documents the
			// manual extraEnv/extraVolumes alternative).
			values["azureUseWorkloadIdentityExtension"] = true
		case identity.GetUseManagedIdentity():
			values["azureUseManagedIdentityExtension"] = true
			if identity.GetUserAssignedIdentityId() != "" {
				// Lands in the chart Secret; reaches the pod as
				// ARM_USER_ASSIGNED_IDENTITY_ID only in the
				// managed-identity env branch.
				values["azureUserAssignedIdentityID"] = identity.GetUserAssignedIdentityId()
			}
		case identity.GetServicePrincipal() != nil:
			sp := identity.GetServicePrincipal()
			values["azureTenantID"] = sp.GetTenantId()
			values["azureClientID"] = sp.GetClientId()
			values["azureClientSecret"] = sp.GetClientSecret()
		}

	case spec.GetGce() != nil:
		gce := spec.GetGce()
		// Each entry renders one
		// --node-group-auto-discovery=mig:namePrefix=<name>,min=..,max=..
		// flag (deployment.yaml gce block). NOTE the chart's key really is
		// "autoscalingGroupsnamePrefix" (lowercase n — values.yaml line 64).
		values["autoscalingGroupsnamePrefix"] = nodeGroupsSlice(gce.GetInstanceGroupPrefixes())
		// Deployment render gate — see the function comment. The value is
		// not consumed by any gce template block.
		autoDiscovery["clusterName"] = locals.ResourceName
		if gce.GetWorkloadIdentityServiceAccountEmail() != "" {
			serviceAccountAnnotations["iam.gke.io/gcp-service-account"] = gce.GetWorkloadIdentityServiceAccountEmail()
		}

	case spec.GetClusterApi() != nil:
		capi := spec.GetClusterApi()
		// Matches the chart's own default ("incluster-incluster") when
		// unset — rendered only when set.
		if capi.GetMode() != "" {
			values["clusterAPIMode"] = capi.GetMode()
		}
		if capi.GetKubeconfigSecret() != "" {
			// The chart both mounts this Secret and derives the
			// --kubeconfig/--cloud-config paths from the mode
			// (deployment.yaml clusterapi block + volumes).
			values["clusterAPIKubeconfigSecret"] = capi.GetKubeconfigSecret()
		}
		if capi.GetNamespace() != "" {
			// Renders --node-group-auto-discovery=clusterapi:namespace=<ns>
			// via the capiAutodiscoveryConfig helper — and satisfies the
			// Deployment render gate. With namespace empty the gate needs
			// autoDiscovery.clusterName/labels via helm_values (upstream
			// requires at least one discovery dimension for clusterapi).
			autoDiscovery["namespace"] = capi.GetNamespace()
		}
		if capi.GetNamespaceScopedRbac() {
			// rbac.clusterScoped=false switches the chart from
			// ClusterRole/ClusterRoleBinding to namespaced Role/RoleBinding
			// — the least-privilege posture the chart documents "most
			// useful for Cluster-API".
			rbac["clusterScoped"] = false
		}

	case spec.GetCivo() != nil:
		civo := spec.GetCivo()
		// All four land in the chart's credential Secret
		// (templates/secret.yaml civo branch: api-url/api-key/cluster-id/
		// region) and reach the pod as CIVO_* env vars via secretKeyRef.
		values["civoClusterID"] = civo.GetClusterId()
		values["civoRegion"] = civo.GetRegion()
		values["civoApiKey"] = civo.GetApiKey()
		// Matches the chart's own default when unset — rendered only when
		// set.
		if civo.GetApiUrl() != "" {
			values["civoApiUrl"] = civo.GetApiUrl()
		}
		// NO gate key is injected for civo: the chart README requires
		// autoscalingGroups (the Civo node pools with size bounds), which
		// the typed spec does not model — they ride helm_values, and that
		// same document satisfies the Deployment render gate.

	case spec.GetKwok() != nil:
		kwok := spec.GetKwok()
		// Matches the chart's own default ("kwok-provider-config") when
		// unset. Reaches the pod as KWOK_PROVIDER_CONFIGMAP
		// (deployment.yaml kwok env branch) and names the ConfigMap the
		// chart itself creates for kwok (templates/configmap.yaml).
		if kwok.GetConfigMapName() != "" {
			values["kwokConfigMapName"] = kwok.GetConfigMapName()
		}
		// Deployment render gate — see the function comment. kwok node
		// groups come from the provider ConfigMap, not from values.
		autoDiscovery["clusterName"] = locals.ResourceName
	}

	// ---- scaling flags → extraArgs ----------------------------------------
	// The chart renders every extraArgs entry as --<key>=<value>
	// (deployment.yaml) — flag names carry no leading dashes. Values are
	// rendered as strings on BOTH engines (they are CLI flag text).
	//
	// PRECEDENCE (comment mirrored in the Terraform module): the typed
	// scaling block renders first, then spec.extra_args merges OVER it —
	// user entries win on key collision. The chart's own extraArgs
	// defaults (logtostderr/stderrthreshold/v) stay chart-side: Helm
	// coalesces our extraArgs map over the chart default per key, so
	// unspecified defaults survive on both engines identically.
	extraArgs := map[string]interface{}{}
	if s := spec.GetScaling(); s != nil {
		if s.GetExpander() != "" {
			extraArgs["expander"] = s.GetExpander()
		}
		// Plain bool — false is the upstream default, so only true renders.
		if s.GetBalanceSimilarNodeGroups() {
			extraArgs["balance-similar-node-groups"] = "true"
		}
		if s.GetScanInterval() != "" {
			extraArgs["scan-interval"] = s.GetScanInterval()
		}
		if s.GetMaxNodeProvisionTime() != "" {
			extraArgs["max-node-provision-time"] = s.GetMaxNodeProvisionTime()
		}
		// Presence-aware optional bools: upstream defaults are true, so an
		// explicit false MUST render — only absence stays silent.
		if s.SkipNodesWithLocalStorage != nil {
			extraArgs["skip-nodes-with-local-storage"] = strconv.FormatBool(s.GetSkipNodesWithLocalStorage())
		}
		if s.SkipNodesWithSystemPods != nil {
			extraArgs["skip-nodes-with-system-pods"] = strconv.FormatBool(s.GetSkipNodesWithSystemPods())
		}
		if sd := s.GetScaleDown(); sd != nil {
			if sd.Enabled != nil {
				extraArgs["scale-down-enabled"] = strconv.FormatBool(sd.GetEnabled())
			}
			if sd.GetUtilizationThreshold() != "" {
				extraArgs["scale-down-utilization-threshold"] = sd.GetUtilizationThreshold()
			}
			if sd.GetUnneededTime() != "" {
				extraArgs["scale-down-unneeded-time"] = sd.GetUnneededTime()
			}
			if sd.GetDelayAfterAdd() != "" {
				extraArgs["scale-down-delay-after-add"] = sd.GetDelayAfterAdd()
			}
			if sd.GetDelayAfterDelete() != "" {
				extraArgs["scale-down-delay-after-delete"] = sd.GetDelayAfterDelete()
			}
			if sd.GetDelayAfterFailure() != "" {
				extraArgs["scale-down-delay-after-failure"] = sd.GetDelayAfterFailure()
			}
		}
	}
	for k, v := range spec.GetExtraArgs() {
		// User entries win over the typed scaling block on key collision.
		extraArgs[k] = v
	}
	if len(extraArgs) > 0 {
		values["extraArgs"] = extraArgs
	}

	// ---- deployment sizing and scheduling ----------------------------------
	if d := spec.GetDeployment(); d != nil {
		if d.Replicas != nil {
			// Replicas leader-elect; extras are warm standbys.
			values["replicaCount"] = int(d.GetReplicas())
		}
		if r := resourcesMap(d.GetResources()); r != nil {
			values["resources"] = r
		}
		// Matches the chart's own default ("system-cluster-critical") when
		// unset — rendered only when set.
		if d.GetPriorityClassName() != "" {
			values["priorityClassName"] = d.GetPriorityClassName()
		}
		if len(d.GetNodeSelector()) > 0 {
			values["nodeSelector"] = stringMapToInterface(d.GetNodeSelector())
		}
		if len(d.GetTolerations()) > 0 {
			values["tolerations"] = tolerationsSlice(d.GetTolerations())
		}
	}

	// ---- own telemetry -------------------------------------------------------
	if p := spec.GetPrometheus(); p.GetServiceMonitor() {
		serviceMonitor := map[string]interface{}{"enabled": true}
		if p.GetServiceMonitorSelectorRelease() != "" {
			// serviceMonitor.selector is a plain label map rendered onto
			// the ServiceMonitor's metadata.labels; Helm's per-key
			// coalesce replaces the chart's default
			// {release: prometheus-operator} entry.
			serviceMonitor["selector"] = map[string]interface{}{
				"release": p.GetServiceMonitorSelectorRelease(),
			}
		}
		values["serviceMonitor"] = serviceMonitor
	}

	// ---- assembled sub-maps ----------------------------------------------------
	if len(serviceAccountAnnotations) > 0 {
		rbac["serviceAccount"] = map[string]interface{}{
			"annotations": serviceAccountAnnotations,
		}
	}
	if len(rbac) > 0 {
		values["rbac"] = rbac
	}
	if len(autoDiscovery) > 0 {
		values["autoDiscovery"] = autoDiscovery
	}

	// ---- escape hatch (merged LAST, helm -f semantics) ------------------------
	if spec.GetHelmValues() != "" {
		overrides := map[string]interface{}{}
		if err := yaml.Unmarshal([]byte(spec.GetHelmValues()), &overrides); err != nil {
			return nil, errors.Wrap(err, "failed to parse helm_values as a YAML document")
		}
		values = mergeMaps(values, overrides)
	}

	return values, nil
}

// nodeGroupsSlice renders the shared NodeGroup list into the chart's
// autoscalingGroups / autoscalingGroupsnamePrefix shape ({name, minSize,
// maxSize} — same keys for both, verified in values.yaml).
func nodeGroupsSlice(groups []*kubernetesclusterautoscalerv1alpha1.KubernetesClusterAutoscalerNodeGroup) []interface{} {
	out := make([]interface{}, 0, len(groups))
	for _, g := range groups {
		out = append(out, map[string]interface{}{
			"name":    g.GetName(),
			"minSize": int(g.GetMinSize()),
			"maxSize": int(g.GetMaxSize()),
		})
	}
	return out
}
