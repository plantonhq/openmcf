package module

import (
	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
)

// buildHelmValues renders the typed spec into the strimzi-kafka-operator
// chart's values map, then merges the spec's helm_values escape hatch over
// it with Helm `-f` semantics (maps deep-merge with the later document
// winning, lists replace).
//
// PARITY: the Terraform module reaches the same result natively — its
// helm_release passes values = [yamlencode(typed values), helm_values] and
// the provider merges the documents in exactly this order. Keep every typed
// mapping below in lockstep with the Terraform module's locals.
//
// Chart-default-matching values render only on divergence (the watch
// scope, the true-defaulted toggles, the image source), so the rendered
// values stay minimal on both engines.
func buildHelmValues(locals *Locals) (map[string]interface{}, error) {
	spec := locals.Spec

	values := map[string]interface{}{}

	// ---- operator sizing -------------------------------------------------
	if spec.Replicas != nil {
		values["replicas"] = int(spec.GetReplicas())
	}
	// The chart SHIPS default requests/limits (requests 200m/384Mi,
	// limits 1000m/384Mi) — the resources key renders only when the spec
	// sets them, so the chart defaults survive an empty spec. Helm
	// deep-merges per key, so a partial spec block overrides only the
	// halves it carries.
	if r := resourcesMap(spec.GetResources()); r != nil {
		values["resources"] = r
	}

	// ---- watch scope -----------------------------------------------------
	// The chart models the scope with two independent values:
	// watchAnyNamespace (cluster-wide RBAC, default false) and
	// watchNamespaces (a LIST of extra namespaces; the installation
	// namespace is always watched). Spec CEL rules make the two arms
	// mutually exclusive, so at most one renders.
	if spec.GetWatch().GetAnyNamespace() {
		values["watchAnyNamespace"] = true
	}
	if len(spec.GetWatch().GetNamespaces()) > 0 {
		namespaces := make([]interface{}, 0, len(spec.GetWatch().GetNamespaces()))
		for _, ns := range spec.GetWatch().GetNamespaces() {
			namespaces = append(namespaces, ns)
		}
		values["watchNamespaces"] = namespaces
	}

	// ---- reconciliation timing ---------------------------------------------
	// Both integers in the chart's values file; rendered on presence so
	// the chart defaults (120000 / 300000) survive an empty spec.
	if spec.FullReconciliationIntervalMs != nil {
		values["fullReconciliationIntervalMs"] = int(spec.GetFullReconciliationIntervalMs())
	}
	if spec.OperationTimeoutMs != nil {
		values["operationTimeoutMs"] = int(spec.GetOperationTimeoutMs())
	}

	// ---- logging / gates / DNS domain ---------------------------------------
	// logLevel renders whenever the spec carries a value (the chart
	// default is INFO via an env-substituted expression; a literal level
	// replaces it cleanly).
	if spec.LogLevel != nil && spec.GetLogLevel() != "" {
		values["logLevel"] = spec.GetLogLevel()
	}
	if spec.GetFeatureGates() != "" {
		values["featureGates"] = spec.GetFeatureGates()
	}
	if spec.KubernetesServiceDnsDomain != nil && spec.GetKubernetesServiceDnsDomain() != "" {
		values["kubernetesServiceDnsDomain"] = spec.GetKubernetesServiceDnsDomain()
	}

	// ---- leader election ------------------------------------------------------
	// The chart nests a single enable flag. Rendered on presence — an
	// explicit true re-states the chart default harmlessly, an explicit
	// false is the actual opt-out.
	if spec.LeaderElectionEnabled != nil {
		values["leaderElection"] = map[string]interface{}{
			"enable": spec.GetLeaderElectionEnabled(),
		}
	}

	// ---- operand policy generation toggles --------------------------------------
	// Both default true in the chart; rendered on presence so an explicit
	// opt-out reaches the chart and an empty spec changes nothing.
	if spec.GenerateNetworkPolicy != nil {
		values["generateNetworkPolicy"] = spec.GetGenerateNetworkPolicy()
	}
	if spec.GeneratePodDisruptionBudget != nil {
		values["generatePodDisruptionBudget"] = spec.GetGeneratePodDisruptionBudget()
	}

	// ---- cluster-scoped RBAC ownership -----------------------------------------
	// createGlobalResources defaults true; false is the second-install
	// posture (the fixed-name ClusterRoles are owned by the first
	// release).
	if spec.CreateGlobalResources != nil {
		values["createGlobalResources"] = spec.GetCreateGlobalResources()
	}

	// ---- scheduling ----------------------------------------------------------
	if len(spec.GetNodeSelector()) > 0 {
		values["nodeSelector"] = stringMapToInterface(spec.GetNodeSelector())
	}
	if len(spec.GetTolerations()) > 0 {
		values["tolerations"] = tolerationsSlice(spec.GetTolerations())
	}

	// ---- image source ----------------------------------------------------------
	// defaultImageRegistry/Repository/Tag steer EVERY Strimzi image (the
	// operator and all operand images it deploys) — the air-gap path. Pull
	// secrets ride the chart's image.imagePullSecrets list (raw
	// Kubernetes object list, piped into the pod spec).
	if spec.GetImage().GetRegistry() != "" {
		values["defaultImageRegistry"] = spec.GetImage().GetRegistry()
	}
	if spec.GetImage().GetRepository() != "" {
		values["defaultImageRepository"] = spec.GetImage().GetRepository()
	}
	if spec.GetImage().GetTag() != "" {
		values["defaultImageTag"] = spec.GetImage().GetTag()
	}
	if len(spec.GetImagePullSecrets()) > 0 {
		pullSecrets := make([]interface{}, 0, len(spec.GetImagePullSecrets()))
		for _, name := range spec.GetImagePullSecrets() {
			pullSecrets = append(pullSecrets, map[string]interface{}{"name": name})
		}
		values["image"] = map[string]interface{}{
			"imagePullSecrets": pullSecrets,
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
