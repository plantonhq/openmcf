package module

import (
	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
)

// buildHelmValues renders the typed spec into the flink-kubernetes-operator
// chart's values map, then merges the spec's helm_values escape hatch over
// it with Helm `-f` semantics (maps deep-merge with the later document
// winning, lists replace).
//
// PARITY: the Terraform module reaches the same result natively — its
// helm_release passes values = [yamlencode(typed values), helm_values,
// yamlencode(keystore re-pin)] and the provider merges the documents in
// exactly this order. Keep every typed mapping below in lockstep with the
// Terraform module's locals.
//
// Chart-default-matching values render only on divergence, so the
// rendered values stay minimal on both engines — with the deliberate
// always-rendered exceptions called out inline (nameOverride,
// image.tag, the keystore wiring whenever the webhook is on) and ONE
// key re-pinned AFTER the escape-hatch merge (see the bottom of this
// function).
func buildHelmValues(locals *Locals) (map[string]interface{}, error) {
	spec := locals.Spec

	values := map[string]interface{}{
		// nameOverride is THIS chart's identity pin: the operator
		// Deployment (and its selector labels) render from the
		// `flink-operator.name` helper (default .Chart.Name |
		// nameOverride) — fullnameOverride alone is a no-op for the
		// Deployment and leaves it at the chart-constant
		// `flink-kubernetes-operator` (verified live: the pinned name
		// was NotFound while the chart-named Deployment served). Both
		// keys are set so any fullname-derived fallback also hangs off
		// the resource name.
		"nameOverride":     locals.ReleaseName,
		"fullnameOverride": locals.ReleaseName,
	}

	// ---- operator image ---------------------------------------------------
	// tag is ALWAYS pinned to the chart version: the chart's own default
	// tag is the unpinned "latest" — the one values.yaml default that
	// must never stand. image_registry replaces ONLY the registry part
	// (chart default `ghcr.io/apache/flink-kubernetes-operator`); this
	// never rewrites the Flink images deployments run — those ride each
	// KubernetesFlinkDeployment's own image field.
	image := map[string]interface{}{
		"tag": locals.ChartVersion,
	}
	if spec.GetImageRegistry() != "" {
		image["repository"] = spec.GetImageRegistry() + "/" + vars.OperatorImagePath
	}
	values["image"] = image

	// ---- watch scope ------------------------------------------------------
	// Scopes RBAC AND the admission webhook's namespaceSelector to
	// exactly these namespaces (template-verified: the chart wires
	// kubernetes.operator.watched.namespaces from the same list).
	if len(spec.GetWatchNamespaces()) > 0 {
		watchNamespaces := make([]interface{}, 0, len(spec.GetWatchNamespaces()))
		for _, ns := range spec.GetWatchNamespaces() {
			watchNamespaces = append(watchNamespaces, ns)
		}
		values["watchNamespaces"] = watchNamespaces
	}

	// ---- webhook ------------------------------------------------------------
	// Enabled: the module-owned keystore password replaces the chart's
	// hardcoded-public default. Disabled: webhook.create=false removes
	// the webhook, the certificate machinery, and the cert-manager
	// dependency.
	if locals.WebhookEnabled {
		values["webhook"] = map[string]interface{}{
			"keystore": map[string]interface{}{
				"useDefaultPassword": false,
				"passwordSecretRef": map[string]interface{}{
					"name": locals.WebhookKeystoreSecretName,
					"key":  "password",
				},
			},
		}
	} else {
		values["webhook"] = map[string]interface{}{"create": false}
	}

	// Rendered on presence — an explicit 1 re-states the chart default
	// harmlessly; >1 pairs with the leader-election keys rendered into
	// locals.OperatorConfig (the chart REFUSES multi-replica installs
	// without leader election, by design).
	if spec.Replicas != nil {
		values["replicas"] = spec.GetReplicas()
	}

	// ---- operator configuration --------------------------------------------
	// create:true and append:true ARE the chart defaults — rendered
	// explicitly alongside the file for self-documentation (append keeps
	// the chart's built-in conf underneath; ours layers over).
	if len(locals.OperatorConfig) > 0 {
		values["defaultConfiguration"] = map[string]interface{}{
			"create":          true,
			"append":          true,
			"flink-conf.yaml": locals.FlinkConfFile,
		}
	}

	// Rendered only on divergence from the chart default "flink". The
	// chart keeps the service account past uninstall
	// (helm.sh/resource-policy: keep) — running jobs never lose their
	// identity.
	if locals.JobServiceAccount != "flink" {
		values["jobServiceAccount"] = map[string]interface{}{
			"name": locals.JobServiceAccount,
		}
	}

	// ---- operator pod (resources + scheduling) --------------------------------
	// priorityClassName lives under operatorPod in this chart's values
	// (template-verified) — alongside nodeSelector/tolerations.
	operatorPod := map[string]interface{}{}
	// The chart ships NO default requests/limits for the operator
	// container — the resources key renders only when the spec sets
	// them. Helm deep-merges per key.
	if r := resourcesMap(spec.GetResources()); r != nil {
		operatorPod["resources"] = r
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
	if len(operatorPod) > 0 {
		values["operatorPod"] = operatorPod
	}

	// ---- pull secrets ---------------------------------------------------------
	// Raw Kubernetes object list, piped into the pod spec.
	if len(spec.GetImagePullSecrets()) > 0 {
		pullSecrets := make([]interface{}, 0, len(spec.GetImagePullSecrets()))
		for _, name := range spec.GetImagePullSecrets() {
			pullSecrets = append(pullSecrets, map[string]interface{}{"name": name})
		}
		values["imagePullSecrets"] = pullSecrets
	}

	// ---- escape hatch (merged LAST, helm -f semantics) -----------------------
	if spec.GetHelmValues() != "" {
		overrides := map[string]interface{}{}
		if err := yaml.Unmarshal([]byte(spec.GetHelmValues()), &overrides); err != nil {
			return nil, errors.Wrap(err, "failed to parse helm_values as a YAML document")
		}
		values = mergeMaps(values, overrides)
	}

	// The keystore re-pin, applied AFTER the escape-hatch merge — the
	// deliberate exception to helm -f semantics (twin of the Terraform
	// module's third values document): whenever the webhook is enabled,
	// webhook.keystore.useDefaultPassword pins false so the chart's
	// hardcoded-password default can never resurface through
	// helm_values.
	if locals.WebhookEnabled {
		webhook, ok := values["webhook"].(map[string]interface{})
		if !ok {
			webhook = map[string]interface{}{}
			values["webhook"] = webhook
		}
		keystore, ok := webhook["keystore"].(map[string]interface{})
		if !ok {
			keystore = map[string]interface{}{}
			webhook["keystore"] = keystore
		}
		keystore["useDefaultPassword"] = false
	}

	return values, nil
}
