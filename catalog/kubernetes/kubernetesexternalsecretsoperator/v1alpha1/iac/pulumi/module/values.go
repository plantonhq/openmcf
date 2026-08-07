package module

import (
	"sort"

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
	// installCRDs matches the chart's own default (true). keep_on_uninstall
	// has NO chart knob — the chart templates its CRDs and Helm would
	// DELETE them on uninstall, cascading to every ESO object cluster-wide.
	// Planton default keeps them via the standard Helm resource-policy
	// annotation, which the chart forwards onto the CRDs (crds.annotations).
	crdsInstall := true
	if spec.GetCrds() != nil && spec.GetCrds().Install != nil {
		crdsInstall = spec.GetCrds().GetInstall()
	}
	crdsKeep := true
	if spec.GetCrds() != nil && spec.GetCrds().KeepOnUninstall != nil {
		crdsKeep = spec.GetCrds().GetKeepOnUninstall()
	}
	values["installCRDs"] = crdsInstall
	if crdsInstall && crdsKeep {
		values["crds"] = map[string]interface{}{
			"annotations": map[string]interface{}{"helm.sh/resource-policy": "keep"},
		}
	}

	// ---- controller ----------------------------------------------------
	if spec.Replicas != nil {
		values["replicaCount"] = int(spec.GetReplicas())
	}
	if spec.GetLeaderElect() {
		values["leaderElect"] = true
	}
	if spec.Concurrent != nil {
		values["concurrent"] = int(spec.GetConcurrent())
	}
	if spec.GetControllerClass() != "" {
		values["controllerClass"] = spec.GetControllerClass()
	}
	if spec.GetScopedNamespace() != "" {
		values["scopedNamespace"] = spec.GetScopedNamespace()
	}
	if spec.GetScopedRbac() {
		values["scopedRBAC"] = true
	}
	if spec.LogLevel != nil {
		values["log"] = map[string]interface{}{"level": spec.GetLogLevel()}
	}
	if r := resourcesMap(spec.GetResources()); r != nil {
		values["resources"] = r
	}

	// ---- workload identity ----------------------------------------------
	// The chart creates the controller ServiceAccount; the module pins its
	// name (deterministic identity subject) and rides ambient-identity
	// annotations on it. Per-store identities (store auth blocks) need
	// nothing here.
	serviceAccount := map[string]interface{}{"name": locals.ControllerServiceAccount}
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
			// The azure-workload-identity webhook only injects the
			// federated token volume into pods carrying this label.
			values["podLabels"] = map[string]interface{}{"azure.workload.identity/use": "true"}
		}
		if len(annotations) > 0 {
			serviceAccount["annotations"] = annotations
		}
	}
	values["serviceAccount"] = serviceAccount

	// ---- scheduling --------------------------------------------------------
	if len(spec.GetNodeSelector()) > 0 {
		values["nodeSelector"] = stringMapToInterface(spec.GetNodeSelector())
	}
	if len(spec.GetTolerations()) > 0 {
		values["tolerations"] = tolerationsSlice(spec.GetTolerations())
	}
	if spec.GetPriorityClassName() != "" {
		values["priorityClassName"] = spec.GetPriorityClassName()
	}
	if spec.GetPodDisruptionBudget() {
		values["podDisruptionBudget"] = map[string]interface{}{"enabled": true, "minAvailable": 1}
	}

	// ---- observability ------------------------------------------------------
	if p := spec.GetPrometheus(); p != nil && p.GetServiceMonitor() {
		sm := map[string]interface{}{"enabled": true}
		if p.ServiceMonitorInterval != nil {
			sm["interval"] = p.GetServiceMonitorInterval()
		}
		if len(p.GetServiceMonitorLabels()) > 0 {
			sm["additionalLabels"] = stringMapToInterface(p.GetServiceMonitorLabels())
		}
		values["serviceMonitor"] = sm
	}

	// ---- webhook -------------------------------------------------------------
	if w := spec.GetWebhook(); w != nil {
		webhook := map[string]interface{}{}
		if w.Enabled != nil && !w.GetEnabled() {
			webhook["create"] = false
		}
		if w.Replicas != nil {
			webhook["replicaCount"] = int(w.GetReplicas())
		}
		if r := resourcesMap(w.GetResources()); r != nil {
			webhook["resources"] = r
		}
		if len(webhook) > 0 {
			values["webhook"] = webhook
		}
	}

	// ---- cert-controller ------------------------------------------------------
	if c := spec.GetCertController(); c != nil {
		certController := map[string]interface{}{}
		if c.Enabled != nil && !c.GetEnabled() {
			certController["create"] = false
		}
		if c.Replicas != nil {
			certController["replicaCount"] = int(c.GetReplicas())
		}
		if r := resourcesMap(c.GetResources()); r != nil {
			certController["resources"] = r
		}
		if len(certController) > 0 {
			values["certController"] = certController
		}
	}

	// ---- image ----------------------------------------------------------------
	if spec.GetImageRepository() != "" {
		values["image"] = map[string]interface{}{"repository": spec.GetImageRepository()}
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

// tolerationsSlice renders the shared WorkloadToleration list into the
// chart's tolerations shape.
func tolerationsSlice(tolerations []*kubernetesprovider.WorkloadToleration) []interface{} {
	out := make([]interface{}, 0, len(tolerations))
	for _, t := range tolerations {
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
		out = append(out, tol)
	}
	return out
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

// sortedKeys returns a map's keys sorted — kept for deterministic rendering
// of map-shaped fields that land in ordered structures.
func sortedKeys(in map[string]string) []string {
	keys := make([]string, 0, len(in))
	for k := range in {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
