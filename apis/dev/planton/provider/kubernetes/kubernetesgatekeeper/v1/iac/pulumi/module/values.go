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
//
// NAMING: gatekeeper's chart HARDCODES its resource names
// (gatekeeper-webhook-service, gatekeeper-webhook-server-cert, the webhook
// configuration names) — there is no fullname derivation to pin and the
// engine is a per-cluster singleton by construction.
func buildHelmValues(locals *Locals) (map[string]interface{}, error) {
	spec := locals.Spec
	values := map[string]interface{}{}

	if spec.Replicas != nil {
		values["replicas"] = int(spec.GetReplicas())
	}

	// ---- validating webhook ------------------------------------------------------
	if vw := spec.GetValidatingWebhook(); vw != nil {
		if vw.Enabled != nil {
			values["disableValidatingWebhook"] = !vw.GetEnabled()
		}
		if vw.GetFailurePolicy() != "" {
			values["validatingWebhookFailurePolicy"] = vw.GetFailurePolicy()
		}
		if vw.TimeoutSeconds != nil {
			values["validatingWebhookTimeoutSeconds"] = int(vw.GetTimeoutSeconds())
		}
		if vw.GetEnableDeleteOperations() {
			values["enableDeleteOperations"] = true
		}
		if vw.GetCheckIgnoreFailurePolicy() != "" {
			values["validatingWebhookCheckIgnoreFailurePolicy"] = vw.GetCheckIgnoreFailurePolicy()
		}
	}

	// ---- mutating webhook ---------------------------------------------------------
	if mw := spec.GetMutatingWebhook(); mw != nil {
		if mw.Enabled != nil {
			values["disableMutation"] = !mw.GetEnabled()
		}
		if mw.GetFailurePolicy() != "" {
			values["mutatingWebhookFailurePolicy"] = mw.GetFailurePolicy()
		}
		if mw.TimeoutSeconds != nil {
			values["mutatingWebhookTimeoutSeconds"] = int(mw.GetTimeoutSeconds())
		}
		if mw.GetMutationAnnotations() {
			values["mutationAnnotations"] = true
		}
	}

	// ---- audit ----------------------------------------------------------------------
	audit := map[string]interface{}{}
	if a := spec.GetAudit(); a != nil {
		if a.IntervalSeconds != nil {
			values["auditInterval"] = int(a.GetIntervalSeconds())
		}
		if a.ConstraintViolationsLimit != nil {
			values["constraintViolationsLimit"] = int(a.GetConstraintViolationsLimit())
		}
		if a.GetFromCache() {
			values["auditFromCache"] = true
		}
		if a.GetMatchKindOnly() {
			values["auditMatchKindOnly"] = true
		}
		if a.ChunkSize != nil {
			values["auditChunkSize"] = int(a.GetChunkSize())
		}
		if resources := resourcesBlock(a.GetResources()); resources != nil {
			audit["resources"] = resources
		}
	}

	// ---- namespace exemptions --------------------------------------------------------
	controllerManager := map[string]interface{}{}
	if len(spec.GetExemptNamespaces()) > 0 {
		controllerManager["exemptNamespaces"] = toInterfaceSlice(spec.GetExemptNamespaces())
	}
	if len(spec.GetExemptNamespacePrefixes()) > 0 {
		controllerManager["exemptNamespacePrefixes"] = toInterfaceSlice(spec.GetExemptNamespacePrefixes())
	}

	// ---- engine capabilities -----------------------------------------------------------
	if engine := spec.GetEngine(); engine != nil {
		if engine.EnableExternalData != nil {
			values["enableExternalData"] = engine.GetEnableExternalData()
		}
		if engine.EnableK8SNativeValidation != nil {
			values["enableK8sNativeValidation"] = engine.GetEnableK8SNativeValidation()
		}
		if engine.EnableGeneratorResourceExpansion != nil {
			values["enableGeneratorResourceExpansion"] = engine.GetEnableGeneratorResourceExpansion()
		}
		if len(engine.GetDisabledBuiltins()) > 0 {
			values["disabledBuiltins"] = toInterfaceSlice(engine.GetDisabledBuiltins())
		}
		if engine.GetLogDenies() {
			values["logDenies"] = true
		}
		if engine.GetLogLevel() != "" {
			values["logLevel"] = engine.GetLogLevel()
		}
	}

	// ---- controller-manager sizing + scheduling -------------------------------------------
	// The scheduling entries apply to BOTH deployments (controller-manager
	// and audit) — placing an engine on dedicated nodes must move the
	// audit loop with it, or audit findings and enforcement diverge.
	if resources := resourcesBlock(spec.GetResources()); resources != nil {
		controllerManager["resources"] = resources
	}
	if scheduling := spec.GetScheduling(); scheduling != nil {
		if len(scheduling.GetNodeSelector()) > 0 {
			nodeSelector := stringMapToInterface(scheduling.GetNodeSelector())
			controllerManager["nodeSelector"] = nodeSelector
			audit["nodeSelector"] = nodeSelector
		}
		if len(scheduling.GetTolerations()) > 0 {
			tolerations := tolerationsSlice(scheduling.GetTolerations())
			controllerManager["tolerations"] = tolerations
			audit["tolerations"] = tolerations
		}
	}

	// ---- external webhook certificate (the cert-manager arm) --------------------------------
	// CHART-TRUTH ASYMMETRY at the pin: with externalCertInjection
	// enabled the AUDIT deployment auto-disables its cert rotation
	// (`or .Values.audit.disableCertRotation .Values.externalCertInjection.enabled`)
	// but the CONTROLLER-MANAGER reads only its own flag — without
	// disableCertRotation=true the embedded rotator keeps overwriting the
	// injected Secret. The module sets it explicitly.
	if ec := spec.GetExternalCert(); ec != nil {
		values["externalCertInjection"] = map[string]interface{}{
			"enabled":    true,
			"secretName": ec.GetSecretName().GetValue(),
		}
		controllerManager["disableCertRotation"] = true
	}

	// ---- lifecycle hooks -----------------------------------------------------------------------
	if hooks := spec.GetHooks(); hooks != nil {
		if hooks.LabelNamespace != nil {
			values["postInstall"] = map[string]interface{}{
				"labelNamespace": map[string]interface{}{
					"enabled": hooks.GetLabelNamespace(),
				},
			}
		}
		if hooks.ProbeWebhook != nil {
			postInstall, ok := values["postInstall"].(map[string]interface{})
			if !ok {
				postInstall = map[string]interface{}{}
				values["postInstall"] = postInstall
			}
			postInstall["probeWebhook"] = map[string]interface{}{
				"enabled": hooks.GetProbeWebhook(),
			}
		}
		if hooks.UpgradeCrds != nil {
			values["upgradeCRDs"] = map[string]interface{}{
				"enabled": hooks.GetUpgradeCrds(),
			}
		}
		if hooks.GetDeleteWebhookConfigurationsOnUninstall() {
			// The extra "name" key works around a chart bug at the 3.23.0
			// pin: the hook's ClusterRoleBinding subject renders
			// .Values.preUninstall.deleteWebhookConfigurations.name — a key
			// that does not exist in the chart's own values (the SA name
			// lives under .serviceAccount.name) — so enabling the arm
			// without it fails every uninstall on the CRB's empty subject.
			// Value = the chart's SA-name default. Terraform twin renders
			// it identically.
			values["preUninstall"] = map[string]interface{}{
				"deleteWebhookConfigurations": map[string]interface{}{
					"enabled": true,
					"name":    "gatekeeper-delete-webhook-configs",
				},
			}
		}
	}

	// ---- image override (air-gap) -----------------------------------------------------------
	// The chart's image keys are repository + RELEASE (the tag key is
	// named "release" — chart-truth) with a separate crdRepository for
	// the hook containers. The typed override maps repo/tag onto
	// repository/release; the crdRepository and curl hook images ride
	// helm_values for full air-gap installs (the spec comment teaches
	// it). Each key renders independently: pull_secret_name alone is a
	// legal shape (authenticated pulls of the DEFAULT image, e.g. Docker
	// Hub rate-limit credentials) and must render image.pullSecrets even
	// with repo/tag unset — the Terraform twin prunes per key the same
	// way.
	if img := spec.GetImage(); img != nil && (img.GetRepo() != "" || img.GetTag() != "" || img.GetPullSecretName() != "") {
		image := map[string]interface{}{}
		if img.GetRepo() != "" {
			image["repository"] = img.GetRepo()
		}
		if img.GetTag() != "" {
			image["release"] = img.GetTag()
		}
		if img.GetPullSecretName() != "" {
			image["pullSecrets"] = []interface{}{
				map[string]interface{}{"name": img.GetPullSecretName()},
			}
		}
		values["image"] = image
	}

	if len(controllerManager) > 0 {
		values["controllerManager"] = controllerManager
	}
	if len(audit) > 0 {
		values["audit"] = audit
	}

	// ---- escape hatch (merged LAST, Helm -f semantics) ---------------------------------------
	if spec.GetHelmValues() != "" {
		overrides := map[string]interface{}{}
		if err := yaml.Unmarshal([]byte(spec.GetHelmValues()), &overrides); err != nil {
			return nil, errors.Wrap(err, "failed to parse helm_values as YAML")
		}
		values = mergeMaps(values, overrides)
	}

	return values, nil
}
