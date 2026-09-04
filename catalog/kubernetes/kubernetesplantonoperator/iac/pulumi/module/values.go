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
// the provider merges the documents in exactly this order. Keep every typed
// mapping below in lockstep with the Terraform module's locals.
//
// No fullnameOverride and no nameOverride, deliberately: the chart's name
// helpers feed the `app.kubernetes.io/name: planton-operator` label the
// operator's OWN one-per-cluster startup guard matches on — renaming would
// take the Deployment out of the guard's view and make a second install
// possible (and both would then fight over reconciles). With the release
// name fixed to the chart name the Deployment renders as
// "planton-operator" regardless.
func buildHelmValues(locals *Locals) (map[string]interface{}, error) {
	spec := locals.Spec

	values := map[string]interface{}{}

	// ---- CRDs -----------------------------------------------------------------
	// The chart owns its two definitions as release resources behind these
	// two values. Planton default: install them with the release and keep
	// them on uninstall (kept definitions preserve every PlantonPlatform
	// declaration and the platforms behind them; a later install under the
	// fixed release name adopts them). Rendered unconditionally so the
	// release's values always state the posture, whichever way the dials
	// were left.
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

	// ---- sizing ---------------------------------------------------------------
	if spec.Replicas != nil {
		values["replicaCount"] = int(spec.GetReplicas())
	}
	// leaderElection.enabled matches the chart's own default (true) —
	// rendered only on explicit opt-out (single-replica dev clusters).
	if spec.LeaderElection != nil && !spec.GetLeaderElection() {
		values["leaderElection"] = map[string]interface{}{
			"enabled": false,
		}
	}
	// Rendered values deep-merge OVER the chart's own resource defaults
	// (requests 10m/256Mi, limits 500m/512Mi) — a partial spec keeps the
	// untouched halves.
	if r := resourcesMap(spec.GetResources()); r != nil {
		values["resources"] = r
	}

	// ---- service account --------------------------------------------------------
	serviceAccount := map[string]interface{}{}
	if spec.GetServiceAccount() != nil {
		if spec.GetServiceAccount().Create != nil && !spec.GetServiceAccount().GetCreate() {
			serviceAccount["create"] = false
		}
		if spec.GetServiceAccount().GetName() != "" {
			serviceAccount["name"] = spec.GetServiceAccount().GetName()
		}
		if len(spec.GetServiceAccount().GetAnnotations()) > 0 {
			serviceAccount["annotations"] = stringMapToInterface(spec.GetServiceAccount().GetAnnotations())
		}
	}
	if len(serviceAccount) > 0 {
		values["serviceAccount"] = serviceAccount
	}

	// ---- labels / scheduling ------------------------------------------------------
	if len(spec.GetCommonLabels()) > 0 {
		values["commonLabels"] = stringMapToInterface(spec.GetCommonLabels())
	}
	if len(spec.GetPodAnnotations()) > 0 {
		values["podAnnotations"] = stringMapToInterface(spec.GetPodAnnotations())
	}
	if len(spec.GetNodeSelector()) > 0 {
		values["nodeSelector"] = stringMapToInterface(spec.GetNodeSelector())
	}
	if len(spec.GetTolerations()) > 0 {
		values["tolerations"] = tolerationsSlice(spec.GetTolerations())
	}

	// ---- image ----------------------------------------------------------------
	// Pull secrets are name references in the chart's values
	// ([{name: ...}], toYaml'd verbatim into the pod spec); the image
	// override renders only the halves that are set — an empty tag keeps
	// the chart's appVersion default.
	if len(spec.GetImagePullSecrets()) > 0 {
		pullSecrets := make([]interface{}, 0, len(spec.GetImagePullSecrets()))
		for _, name := range spec.GetImagePullSecrets() {
			pullSecrets = append(pullSecrets, map[string]interface{}{"name": name})
		}
		values["imagePullSecrets"] = pullSecrets
	}
	image := map[string]interface{}{}
	if spec.GetImage().GetRepository() != "" {
		image["repository"] = spec.GetImage().GetRepository()
	}
	if spec.GetImage().GetTag() != "" {
		image["tag"] = spec.GetImage().GetTag()
	}
	if len(image) > 0 {
		values["image"] = image
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
