package module

import (
	"sort"
	"strings"

	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
)

// buildHelmValues renders the typed spec into the spark-kubernetes-operator
// chart's values map, then merges the spec's helm_values escape hatch over
// it with Helm `-f` semantics (maps deep-merge with the later document
// winning, lists replace).
//
// PARITY: the Terraform module reaches the same result natively — its
// helm_release passes values = [yamlencode(typed values), helm_values] and
// the provider merges the documents in exactly this order. Keep every
// typed mapping below in lockstep with the Terraform module's
// locals.typed_values.
//
// Chart-default-matching values render only on divergence, so the rendered
// values stay minimal on both engines — with the deliberate always-rendered
// exceptions called out inline (the RBAC name re-pins and the workload
// surface). There is NO re-pin document after the escape-hatch merge: this
// chart has no release-owned CRDs and no webhook machinery whose keys an
// escape-hatch value could weaponize — the RBAC name re-pins live in the
// typed rendering and an operator deliberately overriding them owns the
// collision consciously (twin of the Terraform module's values comment).
func buildHelmValues(locals *Locals) (map[string]interface{}, error) {
	spec := locals.Spec

	values := map[string]interface{}{
		// nameOverride is THIS chart's identity pin: every named object
		// (the operator Deployment, PDB selector, NetworkPolicy)
		// renders from the `spark-operator.name` helper (default
		// .Chart.Name | nameOverride) — the chart defines a fullname
		// helper but NO template consumes it, so a fullnameOverride pin
		// is a no-op and the Deployment would keep the chart's constant
		// name `spark-kubernetes-operator` (verified live: the pinned
		// name was NotFound while the chart-named Deployment served).
		"nameOverride": locals.ReleaseName,
	}

	// ---- operator container (jvmArgs + resources) -------------------------
	// The chart ships REAL defaults for resources (1 CPU / 2Gi, requests =
	// limits) and a tuned jvmArgs default — both render only when the spec
	// diverges, so the upstream-tested sizing stands otherwise. Helm
	// deep-merges per key: a partial resources block overrides only the
	// halves it carries.
	operatorContainer := map[string]interface{}{}
	if spec.GetJvmArgs() != "" {
		operatorContainer["jvmArgs"] = spec.GetJvmArgs()
	}
	if r := resourcesMap(spec.GetResources()); r != nil {
		operatorContainer["resources"] = r
	}

	// ---- operator pod (scheduling + container) ----------------------------
	operatorPod := map[string]interface{}{}
	if len(operatorContainer) > 0 {
		operatorPod["operatorContainer"] = operatorContainer
	}
	if sched := spec.GetScheduling(); sched != nil {
		if len(sched.GetNodeSelector()) > 0 {
			operatorPod["nodeSelector"] = stringMapToInterface(sched.GetNodeSelector())
		}
		if len(sched.GetTolerations()) > 0 {
			operatorPod["tolerations"] = tolerationsSlice(sched.GetTolerations())
		}
		if sched.GetPriorityClassName() != "" {
			operatorPod["priorityClassName"] = sched.GetPriorityClassName()
		}
	}

	operatorDeployment := map[string]interface{}{}
	// Rendered on presence — an explicit 1 re-states the chart default
	// harmlessly; >1 pairs with the leader-election property rendered
	// into operatorConfiguration below (the chart REFUSES multi-replica
	// installs without it, by design).
	if spec.Replicas != nil {
		operatorDeployment["replicas"] = spec.GetReplicas()
	}
	if len(operatorPod) > 0 {
		operatorDeployment["operatorPod"] = operatorPod
	}
	if len(operatorDeployment) > 0 {
		values["operatorDeployment"] = operatorDeployment
	}

	// ---- operator RBAC name re-pins (ALWAYS rendered) ----------------------
	// The chart hardcodes every cluster-scoped RBAC name as a plain value
	// ("spark-operator-clusterrole", …), which makes a second install
	// anywhere on the cluster collide by construction. Deriving the names
	// from the release identity makes instances coexist — the same defense
	// as the fullname pin, applied to the chart's values-borne names.
	values["operatorRbac"] = map[string]interface{}{
		"serviceAccount": map[string]interface{}{
			"name": locals.ReleaseName,
		},
		"clusterRole": map[string]interface{}{
			"name": locals.ReleaseName + "-clusterrole",
		},
		"clusterRoleBinding": map[string]interface{}{
			"name": locals.ReleaseName + "-clusterrolebinding",
		},
		"configManagement": map[string]interface{}{
			"roleName":        locals.ReleaseName + "-config-monitor",
			"roleBindingName": locals.ReleaseName + "-config-monitor-binding",
		},
	}

	// ---- workload surface (ALWAYS rendered) --------------------------------
	// The watch scope and the workload RBAC are ONE chart surface, decided
	// together. Cluster-wide (not fenced): the workload ClusterRole (chart
	// default) with a release-derived name. Fenced: per-namespace Roles
	// replace it, the chart CREATES the listed namespaces, and
	// overrideWatchedNamespaces (chart default true) wires the operator's
	// watched-namespaces property from the same list — one value, one
	// truth. The workload SERVICE ACCOUNT name deliberately stays the
	// upstream contract ("spark" unless overridden): SparkApplications
	// reference it by that conventional name.
	workloadResources := map[string]interface{}{
		"serviceAccount": map[string]interface{}{
			"name": locals.WorkloadServiceAccount,
		},
		"clusterRole": map[string]interface{}{
			"create": !locals.WorkloadFenced,
			"name":   locals.ReleaseName + "-workload-clusterrole",
		},
		"role": map[string]interface{}{
			"create": locals.WorkloadFenced,
			"name":   locals.ReleaseName + "-workload-role",
		},
		// The chart derives this binding's roleRef ITSELF from
		// clusterRole.create (ClusterRole when true, Role when false) —
		// only the name is ours to pin (template-verified).
		"roleBinding": map[string]interface{}{
			"name": locals.ReleaseName + "-workload-rolebinding",
		},
	}
	if locals.WorkloadFenced {
		namespaceData := make([]interface{}, 0, len(locals.WorkloadNamespaces))
		for _, ns := range locals.WorkloadNamespaces {
			namespaceData = append(namespaceData, ns)
		}
		workloadResources["namespaces"] = map[string]interface{}{
			"create": true,
			"data":   namespaceData,
		}
	}
	values["workloadResources"] = workloadResources

	// ---- operator properties (chart: operatorConfiguration) ----------------
	// The operator is properties-file configured. The chart APPENDS this
	// document over its built-in defaults (operatorConfiguration.append,
	// chart default true — kept). Leader election is module-owned: any
	// replica count beyond 1 REQUIRES it (the chart's own contract), so
	// the property renders exactly when replicas > 1 — never a spec knob
	// that could drift from the replica count.
	operatorConfiguration := map[string]interface{}{}
	operatorProperties := map[string]string{}
	for k, v := range spec.GetOperatorProperties() {
		operatorProperties[k] = v
	}
	if spec.Replicas != nil && *spec.Replicas > 1 {
		operatorProperties["spark.kubernetes.operator.leaderElection.enabled"] = "true"
	}
	if len(operatorProperties) > 0 {
		keys := make([]string, 0, len(operatorProperties))
		for k := range operatorProperties {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		lines := make([]string, 0, len(keys))
		for _, k := range keys {
			lines = append(lines, k+"="+operatorProperties[k])
		}
		operatorConfiguration["spark-operator.properties"] = strings.Join(lines, "\n")
	}
	// Dynamic (hot-reload) properties — rendered only when enabled: the
	// chart creates the ConfigMap (dynamicConfig.create) and the RBAC
	// that lets the operator watch it (operatorRbac.configManagement,
	// chart default true — kept).
	if dc := spec.GetDynamicConfig(); dc.GetEnabled() {
		operatorConfiguration["dynamicConfig"] = map[string]interface{}{
			"enable": true,
			"create": true,
			"data":   stringMapToInterface(dc.GetProperties()),
		}
	}
	if len(operatorConfiguration) > 0 {
		values["operatorConfiguration"] = operatorConfiguration
	}

	// ---- operator image (air-gap/private-mirror registry replacement) ------
	// image_registry replaces ONLY the registry part of the operator image
	// (chart default `apache/spark-kubernetes-operator`, Docker Hub
	// implied); the tag stays the chart's appVersion-locked default. Spark
	// WORKLOAD images ride each SparkApplication's own image field — this
	// never rewrites those.
	if spec.GetImageRegistry() != "" {
		values["image"] = map[string]interface{}{
			"repository": spec.GetImageRegistry() + "/" + vars.OperatorImagePath,
		}
	}

	// ---- pull secrets --------------------------------------------------------
	if len(spec.GetImagePullSecrets()) > 0 {
		pullSecrets := make([]interface{}, 0, len(spec.GetImagePullSecrets()))
		for _, name := range spec.GetImagePullSecrets() {
			pullSecrets = append(pullSecrets, map[string]interface{}{"name": name})
		}
		values["imagePullSecrets"] = pullSecrets
	}

	// ---- escape hatch (merged LAST, helm -f semantics) -----------------------
	// Deliberately NOTHING is re-pinned after this merge — see the
	// function comment for why this chart needs no re-pin document.
	if spec.GetHelmValues() != "" {
		overrides := map[string]interface{}{}
		if err := yaml.Unmarshal([]byte(spec.GetHelmValues()), &overrides); err != nil {
			return nil, errors.Wrap(err, "failed to parse helm_values as a YAML document")
		}
		values = mergeMaps(values, overrides)
	}

	return values, nil
}
