package module

import (
	"github.com/pkg/errors"
	kuberneteskyvernov1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kuberneteskyverno/v1"
	"sigs.k8s.io/yaml"
)

// buildHelmValues renders the typed spec into the chart's values map, then
// merges the spec's helm_values escape hatch over it with Helm `-f`
// semantics (maps deep-merge with the later document winning, lists
// replace).
//
// PARITY: the Terraform module reaches the same result natively — its
// helm_release passes values = [yamlencode(typed values), helm_values] and
// the provider merges the documents in exactly this order. Keep every
// typed mapping below in lockstep with the Terraform module's locals.
func buildHelmValues(locals *Locals) (map[string]interface{}, error) {
	spec := locals.Spec
	values := map[string]interface{}{}

	// fullnameOverride pins the chart fullname to the resource name;
	// every chart-derived child name (controller Deployments, the
	// webhook Service, the runtime ConfigMap) and both exported name
	// outputs are built on that contract.
	values["fullnameOverride"] = locals.ReleaseName

	// ---- CRDs (the crds subchart) ---------------------------------------------
	// The CRDs are chart-TEMPLATED (no keep by default): they install and
	// DELETE with the release, cascade-deleting every policy on the
	// cluster. keep_on_uninstall injects the resource-policy annotation.
	if crds := crdsBlock(spec.GetCrds()); len(crds) > 0 {
		values["crds"] = crds
	}

	// ---- runtime config (resource filters, webhook selectors) -------------------
	if config := configBlock(spec.GetConfig()); len(config) > 0 {
		values["config"] = config
	}

	// ---- feature flags -----------------------------------------------------------
	if features := featuresBlock(spec.GetFeatures()); len(features) > 0 {
		values["features"] = features
	}

	// ---- the four controllers ----------------------------------------------------
	// The admission controller's container resources sit under
	// `container.resources` (it runs an init container + the main
	// container); the other three controllers take `resources` directly —
	// chart-truth from the values shapes at the pin.
	admission := map[string]interface{}{}
	if ac := spec.GetAdmissionController(); ac != nil {
		if ac.Replicas != nil {
			admission["replicas"] = int(ac.GetReplicas())
		}
		if resources := resourcesBlock(ac.GetResources()); resources != nil {
			admission["container"] = map[string]interface{}{"resources": resources}
		}
		mergeInto(admission, schedulingEntries(ac.GetScheduling()))
		if hpa := ac.GetAutoscaling(); hpa != nil {
			autoscaling := map[string]interface{}{
				"enabled":     true,
				"maxReplicas": int(hpa.GetMaxReplicas()),
			}
			if hpa.MinReplicas != nil {
				autoscaling["minReplicas"] = int(hpa.GetMinReplicas())
			}
			if hpa.TargetCpuUtilizationPercentage != nil {
				autoscaling["targetCPUUtilizationPercentage"] = int(hpa.GetTargetCpuUtilizationPercentage())
			}
			admission["autoscaling"] = autoscaling
		}
	}

	background := optionalControllerBlock(spec.GetBackgroundController())
	cleanup := optionalControllerBlock(spec.GetCleanupController())
	reports := optionalControllerBlock(spec.GetReportsController())

	// ---- webhook certificates (the cert-manager arm) -----------------------------
	// Omitted = Kyverno-managed runtime certificates with rotation (the
	// chart default; nothing to render). The cert-manager arm applies to
	// BOTH webhook servers — the admission controller AND the cleanup
	// controller (each exposes its own certManager block; missing one
	// would leave that webhook on runtime certs, a split-brain trust
	// posture).
	if cm := spec.GetCertificates().GetCertManager(); cm != nil {
		certManager := certManagerBlock(cm)
		admission["certManager"] = certManager
		cleanup["certManager"] = certManager
	}

	// ---- metrics: ServiceMonitor fan-out ------------------------------------------
	// All four controllers expose the toggle — enabling only some would
	// silently blind the others (the fan-out completeness lesson).
	if spec.GetMetrics().GetServiceMonitor() {
		serviceMonitor := map[string]interface{}{"enabled": true}
		admission["serviceMonitor"] = serviceMonitor
		background["serviceMonitor"] = serviceMonitor
		cleanup["serviceMonitor"] = serviceMonitor
		reports["serviceMonitor"] = serviceMonitor
	}

	if len(admission) > 0 {
		values["admissionController"] = admission
	}
	if len(background) > 0 {
		values["backgroundController"] = background
	}
	if len(cleanup) > 0 {
		values["cleanupController"] = cleanup
	}
	if len(reports) > 0 {
		values["reportsController"] = reports
	}

	// ---- global image registry (air-gap / mirror) ----------------------------------
	// Overrides the per-image registry (ghcr.io) across every controller
	// and hook container; repository paths and tags stay chart-managed.
	// NOTE: the CRD migration hook defaults to reg.kyverno.io — this
	// override reroutes it too (chart-truth: global registry wins over
	// per-image registries).
	if spec.GetImageRegistry() != "" {
		values["global"] = map[string]interface{}{
			"image": map[string]interface{}{
				"registry": spec.GetImageRegistry(),
			},
		}
	}

	// ---- image pull secrets ----------------------------------------------------------
	// existingImagePullSecrets takes NAMES of pre-existing secrets; the
	// chart's imagePullSecrets (a map) CREATES secrets from registry
	// credentials and is deliberately not modeled — credentials belong in
	// KubernetesSecret resources, not chart values.
	if len(spec.GetImagePullSecrets()) > 0 {
		names := make([]interface{}, 0, len(spec.GetImagePullSecrets()))
		for _, name := range spec.GetImagePullSecrets() {
			names = append(names, name)
		}
		values["existingImagePullSecrets"] = names
	}

	// ---- the pre-delete webhook cleanup hook -------------------------------------------
	// Rendered EXPLICITLY in both states: the hook is the designed
	// uninstall path for the runtime-registered webhook configurations
	// (see the spec comment for the stranding consequences of disabling).
	values["webhooksCleanup"] = map[string]interface{}{
		"enabled": spec.WebhooksCleanupEnabled == nil || spec.GetWebhooksCleanupEnabled(),
	}

	// ---- escape hatch (merged LAST, Helm -f semantics) -----------------------------------
	if spec.GetHelmValues() != "" {
		overrides := map[string]interface{}{}
		if err := yaml.Unmarshal([]byte(spec.GetHelmValues()), &overrides); err != nil {
			return nil, errors.Wrap(err, "failed to parse helm_values as YAML")
		}
		values = mergeMaps(values, overrides)
	}

	// fullnameOverride re-pinned AFTER the merge — the one deliberate
	// exception to the escape hatch's last-word contract (twin of the
	// Terraform module's third values document). Both exported name
	// outputs derive from the fullname; letting an override move it would
	// break them.
	values["fullnameOverride"] = locals.ReleaseName

	return values, nil
}

// crdsBlock renders the crds subchart values. install and the migration
// hook are rendered EXPLICITLY whenever the block is declared (the chart
// defaults both true; explicit rendering keeps the manifest's intent
// visible in the diff).
func crdsBlock(crds *kuberneteskyvernov1.KubernetesKyvernoCrds) map[string]interface{} {
	if crds == nil {
		return nil
	}
	block := map[string]interface{}{
		"install": crds.Install == nil || crds.GetInstall(),
		"migration": map[string]interface{}{
			"enabled": crds.MigrationEnabled == nil || crds.GetMigrationEnabled(),
		},
	}
	if crds.GetKeepOnUninstall() {
		block["annotations"] = map[string]interface{}{
			"helm.sh/resource-policy": "keep",
		}
	}
	return block
}

// configBlock renders the engine's runtime ConfigMap edits. Declaring
// webhook_exclude_namespaces REPLACES the chart's webhooks value, so the
// chart's own kube-system exclusion is re-included here by construction —
// dropping it would put the control plane on the policy path.
func configBlock(config *kuberneteskyvernov1.KubernetesKyvernoConfig) map[string]interface{} {
	if config == nil {
		return nil
	}
	block := map[string]interface{}{}
	if namespaces := config.GetWebhookExcludeNamespaces(); len(namespaces) > 0 {
		excluded := []interface{}{"kube-system"}
		for _, namespace := range namespaces {
			if namespace == "kube-system" {
				continue
			}
			excluded = append(excluded, namespace)
		}
		block["webhooks"] = map[string]interface{}{
			"namespaceSelector": map[string]interface{}{
				"matchExpressions": []interface{}{
					map[string]interface{}{
						"key":      "kubernetes.io/metadata.name",
						"operator": "NotIn",
						"values":   excluded,
					},
				},
			},
		}
	}
	if filters := config.GetResourceFiltersInclude(); len(filters) > 0 {
		block["resourceFiltersInclude"] = toInterfaceSlice(filters)
	}
	if filters := config.GetResourceFiltersExclude(); len(filters) > 0 {
		block["resourceFiltersExclude"] = toInterfaceSlice(filters)
	}
	if groups := config.GetExcludeGroups(); len(groups) > 0 {
		block["excludeGroups"] = toInterfaceSlice(groups)
	}
	if usernames := config.GetExcludeUsernames(); len(usernames) > 0 {
		block["excludeUsernames"] = toInterfaceSlice(usernames)
	}
	if config.GetDefaultRegistry() != "" {
		block["defaultRegistry"] = config.GetDefaultRegistry()
	}
	if config.EnableDefaultRegistryMutation != nil {
		block["enableDefaultRegistryMutation"] = config.GetEnableDefaultRegistryMutation()
	}
	return block
}

// featuresBlock renders the typed feature flags. Every chart feature is
// the nested {enabled: bool} shape (chart-truth at the pin) — a bare bool
// would break template rendering.
func featuresBlock(features *kuberneteskyvernov1.KubernetesKyvernoFeatures) map[string]interface{} {
	if features == nil {
		return nil
	}
	block := map[string]interface{}{}
	if features.GetForceFailurePolicyIgnore() {
		block["forceFailurePolicyIgnore"] = map[string]interface{}{"enabled": true}
	}
	if scan := features.GetBackgroundScan(); scan != nil {
		scanBlock := map[string]interface{}{
			"enabled": scan.Enabled == nil || scan.GetEnabled(),
		}
		if scan.Workers != nil {
			scanBlock["backgroundScanWorkers"] = int(scan.GetWorkers())
		}
		if scan.GetInterval() != "" {
			scanBlock["backgroundScanInterval"] = scan.GetInterval()
		}
		block["backgroundScan"] = scanBlock
	}
	if features.GenerateValidatingAdmissionPolicy != nil {
		block["generateValidatingAdmissionPolicy"] = map[string]interface{}{
			"enabled": features.GetGenerateValidatingAdmissionPolicy(),
		}
	}
	if features.AdmissionReports != nil {
		block["admissionReports"] = map[string]interface{}{"enabled": features.GetAdmissionReports()}
	}
	if features.AggregateReports != nil {
		block["aggregateReports"] = map[string]interface{}{"enabled": features.GetAggregateReports()}
	}
	if features.PolicyReports != nil {
		block["policyReports"] = map[string]interface{}{"enabled": features.GetPolicyReports()}
	}
	if features.GetLoggingFormat() != "" || features.LoggingVerbosity != nil {
		logging := map[string]interface{}{}
		if features.GetLoggingFormat() != "" {
			logging["format"] = features.GetLoggingFormat()
		}
		if features.LoggingVerbosity != nil {
			logging["verbosity"] = int(features.GetLoggingVerbosity())
		}
		block["logging"] = logging
	}
	if eventTypes := features.GetOmitEventTypes(); len(eventTypes) > 0 {
		block["omitEvents"] = map[string]interface{}{
			"eventTypes": toInterfaceSlice(eventTypes),
		}
	}
	return block
}

// optionalControllerBlock renders a background/cleanup/reports controller
// block (enabled, replicas, resources, scheduling). Always returns a map
// so callers can attach serviceMonitor/certManager entries.
func optionalControllerBlock(controller *kuberneteskyvernov1.KubernetesKyvernoOptionalController) map[string]interface{} {
	block := map[string]interface{}{}
	if controller == nil {
		return block
	}
	if controller.Enabled != nil {
		block["enabled"] = controller.GetEnabled()
	}
	if controller.Replicas != nil {
		block["replicas"] = int(controller.GetReplicas())
	}
	if resources := resourcesBlock(controller.GetResources()); resources != nil {
		block["resources"] = resources
	}
	mergeInto(block, schedulingEntries(controller.GetScheduling()))
	return block
}

// certManagerBlock renders the chart's certManager values from the typed
// cert-manager arm. An issuer reference switches the chart off its own
// self-signed ClusterIssuer onto the referenced issuer.
func certManagerBlock(cm *kuberneteskyvernov1.KubernetesKyvernoCertManagerCertificates) map[string]interface{} {
	block := map[string]interface{}{"enabled": true}
	if issuerName := cm.GetIssuerName().GetValue(); issuerName != "" {
		issuerKind := cm.GetIssuerKind()
		if issuerKind == "" {
			issuerKind = "ClusterIssuer"
		}
		block["createSelfSignedIssuer"] = false
		block["issuerRef"] = map[string]interface{}{
			"name":  issuerName,
			"kind":  issuerKind,
			"group": "cert-manager.io",
		}
	}
	return block
}

// schedulingEntries renders node selector + tolerations as chart-value
// entries for a controller block.
func schedulingEntries(scheduling *kuberneteskyvernov1.KubernetesKyvernoScheduling) map[string]interface{} {
	if scheduling == nil {
		return nil
	}
	entries := map[string]interface{}{}
	if len(scheduling.GetNodeSelector()) > 0 {
		entries["nodeSelector"] = stringMapToInterface(scheduling.GetNodeSelector())
	}
	if len(scheduling.GetTolerations()) > 0 {
		entries["tolerations"] = tolerationsSlice(scheduling.GetTolerations())
	}
	return entries
}
