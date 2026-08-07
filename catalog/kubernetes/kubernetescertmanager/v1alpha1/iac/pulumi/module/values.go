package module

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pkg/errors"
	kubernetesprovider "github.com/plantonhq/planton/catalog/kubernetes"
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
func buildHelmValues(locals *Locals) (map[string]interface{}, error) {
	spec := locals.Spec
	values := map[string]interface{}{}

	// ---- CRDs ---------------------------------------------------------
	// Planton default: install CRDs with the release (upstream chart
	// defaults to false and expects a separate kubectl apply — one
	// component owning both halves is strictly simpler). keep=true guards
	// cluster-wide certificate data on uninstall.
	crdsInstall := true
	if spec.GetCrds() != nil && spec.GetCrds().Install != nil {
		crdsInstall = spec.GetCrds().GetInstall()
	}
	crdsKeep := true
	if spec.GetCrds() != nil && spec.GetCrds().KeepOnUninstall != nil {
		crdsKeep = spec.GetCrds().GetKeepOnUninstall()
	}
	values["crds"] = map[string]interface{}{
		"enabled": crdsInstall,
		"keep":    crdsKeep,
	}

	// ---- controller ----------------------------------------------------
	if spec.Replicas != nil {
		values["replicaCount"] = int(spec.GetReplicas())
	}
	if r := resourcesMap(spec.GetResources()); r != nil {
		values["resources"] = r
	}
	if spec.LogLevel != nil {
		setNested(values, map[string]interface{}{"logLevel": int(spec.GetLogLevel())}, "global")
	}
	if spec.GetClusterResourceNamespace() != "" {
		values["clusterResourceNamespace"] = spec.GetClusterResourceNamespace()
	}
	if spec.LeaderElectionNamespace != nil && spec.GetLeaderElectionNamespace() != "kube-system" {
		setNested(values, map[string]interface{}{"namespace": spec.GetLeaderElectionNamespace()}, "global", "leaderElection")
	}
	if spec.GetEnableCertificateOwnerRef() {
		values["enableCertificateOwnerRef"] = true
	}

	// featureGates renders as the chart's comma-separated string, sorted
	// for determinism across engines.
	if len(spec.GetFeatureGates()) > 0 {
		gates := make([]string, 0, len(spec.GetFeatureGates()))
		keys := make([]string, 0, len(spec.GetFeatureGates()))
		for k := range spec.GetFeatureGates() {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			gates = append(gates, fmt.Sprintf("%s=%t", k, spec.GetFeatureGates()[k]))
		}
		values["featureGates"] = strings.Join(gates, ",")
	}

	// ---- DNS-01 self-check resolvers ------------------------------------
	if sc := spec.GetDns01SelfCheck(); sc != nil && len(sc.GetRecursiveNameservers()) > 0 {
		values["dns01RecursiveNameservers"] = strings.Join(sc.GetRecursiveNameservers(), ",")
		if sc.GetRecursiveNameserversOnly() {
			values["dns01RecursiveNameserversOnly"] = true
		}
	}

	if spec.MaxConcurrentChallenges != nil {
		values["maxConcurrentChallenges"] = int(spec.GetMaxConcurrentChallenges())
	}

	// ---- workload identity ----------------------------------------------
	// The chart creates the controller ServiceAccount; the identity
	// annotation rides serviceAccount.annotations. AKS additionally needs
	// the azure.workload.identity/use pod label for the webhook to inject
	// the token volume.
	if wi := spec.GetWorkloadIdentity(); wi != nil {
		annotations := map[string]interface{}{}
		if gke := wi.GetGke(); gke != nil {
			annotations["iam.gke.io/gcp-service-account"] = gke.GetServiceAccountEmail().GetValue()
		}
		if eks := wi.GetEks(); eks != nil {
			annotations["eks.amazonaws.com/role-arn"] = eks.GetRoleArn().GetValue()
		}
		if aks := wi.GetAks(); aks != nil {
			annotations["azure.workload.identity/client-id"] = aks.GetClientId().GetValue()
			if aks.TenantId != nil {
				annotations["azure.workload.identity/tenant-id"] = aks.GetTenantId()
			}
			values["podLabels"] = map[string]interface{}{"azure.workload.identity/use": "true"}
		}
		setNested(values, map[string]interface{}{"annotations": annotations}, "serviceAccount")
	}

	// ---- images ----------------------------------------------------------
	if spec.GetImageRegistry() != "" {
		values["imageRegistry"] = spec.GetImageRegistry()
	}

	// ---- prometheus -------------------------------------------------------
	if p := spec.GetPrometheus(); p != nil {
		prom := map[string]interface{}{}
		if p.Enabled != nil {
			prom["enabled"] = p.GetEnabled()
		}
		if p.GetServiceMonitor() {
			sm := map[string]interface{}{"enabled": true}
			if p.ServiceMonitorInterval != nil {
				sm["interval"] = p.GetServiceMonitorInterval()
			}
			if len(p.GetServiceMonitorLabels()) > 0 {
				sm["labels"] = stringMapToInterface(p.GetServiceMonitorLabels())
			}
			prom["servicemonitor"] = sm
		}
		if len(prom) > 0 {
			values["prometheus"] = prom
		}
	}

	// ---- scheduling --------------------------------------------------------
	if len(spec.GetNodeSelector()) > 0 {
		values["nodeSelector"] = stringMapToInterface(spec.GetNodeSelector())
	}
	if len(spec.GetTolerations()) > 0 {
		tolerations := make([]interface{}, 0, len(spec.GetTolerations()))
		for _, t := range spec.GetTolerations() {
			tol := map[string]interface{}{}
			if t.GetKey() != "" {
				tol["key"] = t.GetKey()
			}
			if t.GetOperator() != "" {
				tol["operator"] = t.GetOperator()
			}
			if t.GetValue() != "" {
				tol["value"] = t.GetValue()
			}
			if t.GetEffect() != "" {
				tol["effect"] = t.GetEffect()
			}
			if t.TolerationSeconds != nil {
				tol["tolerationSeconds"] = t.GetTolerationSeconds()
			}
			tolerations = append(tolerations, tol)
		}
		values["tolerations"] = tolerations
	}
	if spec.GetPodDisruptionBudget() {
		values["podDisruptionBudget"] = map[string]interface{}{"enabled": true}
	}

	// ---- webhook -------------------------------------------------------------
	if w := spec.GetWebhook(); w != nil {
		webhook := map[string]interface{}{}
		if w.Replicas != nil {
			webhook["replicaCount"] = int(w.GetReplicas())
		}
		if w.TimeoutSeconds != nil {
			webhook["timeoutSeconds"] = int(w.GetTimeoutSeconds())
		}
		if w.GetHostNetwork() {
			webhook["hostNetwork"] = true
		}
		if w.SecurePort != nil {
			webhook["securePort"] = int(w.GetSecurePort())
		}
		if r := resourcesMap(w.GetResources()); r != nil {
			webhook["resources"] = r
		}
		if len(webhook) > 0 {
			values["webhook"] = webhook
		}
	}

	// ---- cainjector ------------------------------------------------------------
	if c := spec.GetCainjector(); c != nil {
		cainjector := map[string]interface{}{}
		if c.Enabled != nil {
			cainjector["enabled"] = c.GetEnabled()
		}
		if c.Replicas != nil {
			cainjector["replicaCount"] = int(c.GetReplicas())
		}
		if r := resourcesMap(c.GetResources()); r != nil {
			cainjector["resources"] = r
		}
		if len(cainjector) > 0 {
			values["cainjector"] = cainjector
		}
	}

	// ---- startupapicheck ----------------------------------------------------------
	if s := spec.GetStartupapicheck(); s != nil {
		check := map[string]interface{}{}
		if s.Enabled != nil {
			check["enabled"] = s.GetEnabled()
		}
		if s.Timeout != nil {
			check["timeout"] = s.GetTimeout()
		}
		if len(check) > 0 {
			values["startupapicheck"] = check
		}
	}

	// ---- escape hatch (merged LAST, helm -f semantics) --------------------------
	if spec.GetHelmValues() != "" {
		overrides := map[string]interface{}{}
		if err := yaml.Unmarshal([]byte(spec.GetHelmValues()), &overrides); err != nil {
			return nil, errors.Wrap(err, "failed to parse helm_values as a YAML document")
		}
		values = mergeMaps(values, overrides)
	}

	return values, nil
}

// resourcesMap renders the shared ContainerResources message into the
// chart's resources shape. Returns nil when nothing is set.
func resourcesMap(r *kubernetesprovider.ContainerResources) map[string]interface{} {
	if r == nil {
		return nil
	}
	out := map[string]interface{}{}
	if l := r.GetLimits(); l != nil && (l.GetCpu() != "" || l.GetMemory() != "") {
		limits := map[string]interface{}{}
		if l.GetCpu() != "" {
			limits["cpu"] = l.GetCpu()
		}
		if l.GetMemory() != "" {
			limits["memory"] = l.GetMemory()
		}
		out["limits"] = limits
	}
	if q := r.GetRequests(); q != nil && (q.GetCpu() != "" || q.GetMemory() != "") {
		requests := map[string]interface{}{}
		if q.GetCpu() != "" {
			requests["cpu"] = q.GetCpu()
		}
		if q.GetMemory() != "" {
			requests["memory"] = q.GetMemory()
		}
		out["requests"] = requests
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// setNested writes entries into a nested map path, creating intermediate
// maps as needed and preserving siblings already set.
func setNested(root map[string]interface{}, entries map[string]interface{}, path ...string) {
	node := root
	for _, key := range path {
		child, ok := node[key].(map[string]interface{})
		if !ok {
			child = map[string]interface{}{}
			node[key] = child
		}
		node = child
	}
	for k, v := range entries {
		node[k] = v
	}
}

// mergeMaps deep-merges b over a with Helm's `-f` semantics: nested maps
// merge recursively with b winning per key; everything else (scalars,
// lists) is replaced by b's value.
func mergeMaps(a, b map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(a))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		if bChild, ok := v.(map[string]interface{}); ok {
			if aChild, ok := out[k].(map[string]interface{}); ok {
				out[k] = mergeMaps(aChild, bChild)
				continue
			}
		}
		out[k] = v
	}
	return out
}

// stringMapToInterface converts a map[string]string into the
// map[string]interface{} YAML rendering expects.
func stringMapToInterface(in map[string]string) map[string]interface{} {
	out := make(map[string]interface{}, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
