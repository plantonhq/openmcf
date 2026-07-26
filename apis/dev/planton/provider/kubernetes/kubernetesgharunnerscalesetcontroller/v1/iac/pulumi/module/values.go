package module

import (
	"github.com/pkg/errors"
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

	// fullnameOverride pins the chart fullname (and with it the created
	// ServiceAccount's name — the scale-set discovery handle) to the
	// resource name; the exported service_account_name output is built on
	// that contract.
	values["fullnameOverride"] = locals.ReleaseName

	if spec.Replicas != nil {
		values["replicaCount"] = int(spec.GetReplicas())
	}

	// ---- image override (air-gap) -------------------------------------------
	// The chart takes the image reference COMBINED (image.repository holds
	// the full mirror path; no separate registry value — verified in the
	// chart's values + deployment template at the pin).
	if img := spec.GetImage(); img != nil && (img.GetRepo() != "" || img.GetTag() != "") {
		image := map[string]interface{}{}
		if img.GetRepo() != "" {
			image["repository"] = img.GetRepo()
		}
		if img.GetTag() != "" {
			image["tag"] = img.GetTag()
		}
		values["image"] = image
	}

	if resources := resourcesBlock(spec.GetResources()); resources != nil {
		values["resources"] = resources
	}

	// ---- image pull secrets ----------------------------------------------------
	// Joined with the image override's own pull_secret_name, deduplicated —
	// the chart passes these to the controller pod AND to every listener
	// pod it creates.
	if names := pullSecretNames(locals); len(names) > 0 {
		pullSecrets := make([]interface{}, 0, len(names))
		for _, name := range names {
			pullSecrets = append(pullSecrets, map[string]interface{}{"name": name})
		}
		values["imagePullSecrets"] = pullSecrets
	}

	// ---- flags -------------------------------------------------------------------
	if flags := spec.GetFlags(); flags != nil {
		flagsBlock := map[string]interface{}{}
		if flags.GetLogLevel() != "" {
			flagsBlock["logLevel"] = flags.GetLogLevel()
		}
		if flags.GetLogFormat() != "" {
			flagsBlock["logFormat"] = flags.GetLogFormat()
		}
		if flags.GetWatchSingleNamespace() != "" {
			flagsBlock["watchSingleNamespace"] = flags.GetWatchSingleNamespace()
		}
		if flags.RunnerMaxConcurrentReconciles != nil {
			flagsBlock["runnerMaxConcurrentReconciles"] = int(flags.GetRunnerMaxConcurrentReconciles())
		}
		if flags.GetUpdateStrategy() != "" {
			flagsBlock["updateStrategy"] = flags.GetUpdateStrategy()
		}
		if len(flags.GetExcludeLabelPropagationPrefixes()) > 0 {
			prefixes := make([]interface{}, 0, len(flags.GetExcludeLabelPropagationPrefixes()))
			for _, prefix := range flags.GetExcludeLabelPropagationPrefixes() {
				prefixes = append(prefixes, prefix)
			}
			flagsBlock["excludeLabelPropagationPrefixes"] = prefixes
		}
		if flags.K8SClientRateLimiterQps != nil {
			flagsBlock["k8sClientRateLimiterQPS"] = int(flags.GetK8SClientRateLimiterQps())
		}
		if flags.K8SClientRateLimiterBurst != nil {
			flagsBlock["k8sClientRateLimiterBurst"] = int(flags.GetK8SClientRateLimiterBurst())
		}
		if flags.GetRateLimiter() != "" {
			flagsBlock["rateLimiter"] = map[string]interface{}{
				"name": flags.GetRateLimiter(),
			}
		}
		if flags.GetHealthProbeBindAddress() != "" {
			flagsBlock["healthProbeBindAddress"] = flags.GetHealthProbeBindAddress()
		}
		if len(flagsBlock) > 0 {
			values["flags"] = flagsBlock
		}
		if flags.GetPriorityClassName() != "" {
			values["priorityClassName"] = flags.GetPriorityClassName()
		}
	}

	// ---- metrics -------------------------------------------------------------------
	// Declaring the block ENABLES metrics: the chart wires the three
	// addresses into the controller args and every listener pod; leaving
	// it absent passes empty flags (metrics disabled — the chart
	// default).
	if metrics := spec.GetMetrics(); metrics != nil {
		values["metrics"] = map[string]interface{}{
			"controllerManagerAddr": metrics.GetControllerManagerAddr(),
			"listenerAddr":          metrics.GetListenerAddr(),
			"listenerEndpoint":      metrics.GetListenerEndpoint(),
		}
	}

	// ---- scheduling ------------------------------------------------------------------
	if scheduling := spec.GetScheduling(); scheduling != nil {
		if len(scheduling.GetNodeSelector()) > 0 {
			values["nodeSelector"] = stringMapToInterface(scheduling.GetNodeSelector())
		}
		if len(scheduling.GetTolerations()) > 0 {
			values["tolerations"] = tolerationsSlice(scheduling.GetTolerations())
		}
	}

	// ---- escape hatch (merged LAST, Helm -f semantics) --------------------------------
	if spec.GetHelmValues() != "" {
		overrides := map[string]interface{}{}
		if err := yaml.Unmarshal([]byte(spec.GetHelmValues()), &overrides); err != nil {
			return nil, errors.Wrap(err, "failed to parse helm_values as YAML")
		}
		values = mergeMaps(values, overrides)
	}

	// fullnameOverride re-pinned AFTER the merge — the one deliberate
	// exception to the escape hatch's last-word contract (twin of the
	// Terraform module's third values document). The ServiceAccount name
	// output derives from the fullname; letting an override move it would
	// break the scale-set discovery handle.
	values["fullnameOverride"] = locals.ReleaseName

	return values, nil
}

// pullSecretNames joins image_pull_secrets with the image override's own
// pull_secret_name, deduplicated (Terraform twin:
// image_pull_secret_names in locals.tf).
func pullSecretNames(locals *Locals) []string {
	names := append([]string{}, locals.Spec.GetImagePullSecrets()...)
	if extra := locals.Spec.GetImage().GetPullSecretName(); extra != "" {
		seen := false
		for _, name := range names {
			if name == extra {
				seen = true
				break
			}
		}
		if !seen {
			names = append(names, extra)
		}
	}
	return names
}
