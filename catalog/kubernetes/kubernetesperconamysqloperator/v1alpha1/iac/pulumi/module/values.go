package module

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
)

// buildHelmValues renders the typed spec into the pxc-operator chart's
// values map, then merges the spec's helm_values escape hatch over it with
// Helm `-f` semantics (maps deep-merge with the later document winning,
// lists replace).
//
// PARITY: the Terraform module reaches the same result natively — its
// helm_release passes values = [yamlencode(typed values), helm_values] and
// the provider merges the documents in exactly this order. Keep every typed
// mapping below in lockstep with the Terraform module's locals.
//
// Chart-default-matching values render only on divergence (telemetry,
// structured logging, the cluster-wide switch, and the XtraBackup sidecar
// gate all default off upstream), so the rendered values stay minimal on
// both engines.
func buildHelmValues(locals *Locals) (map[string]interface{}, error) {
	spec := locals.Spec

	values := map[string]interface{}{}

	// ---- operator sizing -------------------------------------------------
	if spec.Replicas != nil {
		values["replicaCount"] = int(spec.GetReplicas())
	}
	// The chart SHIPS default requests/limits (requests 100m/20Mi, limits
	// 200m/500Mi) — the resources key renders only when the spec sets
	// them, so the chart defaults survive an empty spec. Helm deep-merges
	// per key, so a partial spec block overrides only the halves it
	// carries.
	if r := resourcesMap(spec.GetResources()); r != nil {
		values["resources"] = r
	}

	// ---- watch scope -----------------------------------------------------
	// The chart models the scope with two independent values:
	// watchAllNamespaces (cluster-wide RBAC, default false) and
	// watchNamespace (a comma-separated namespace list; blank = the
	// operator's own namespace). Spec CEL rules make the two spec arms
	// mutually exclusive, so at most one of the values renders.
	//
	// The chart's own `createNamespace` value (create the WATCHED
	// namespaces) is never rendered: this module owns only the
	// installation namespace; watched namespaces must already exist.
	if spec.GetWatch().GetClusterWide() {
		values["watchAllNamespaces"] = true
	}
	if len(spec.GetWatch().GetNamespaces()) > 0 {
		values["watchNamespace"] = strings.Join(spec.GetWatch().GetNamespaces(), ",")
	}

	// ---- reconciliation throughput ----------------------------------------
	// Both integers in this chart's values file (unlike the sibling
	// MongoDB operator chart, which declares its default as a string).
	if spec.MaxConcurrentReconciles != nil {
		values["maxConcurrentReconciles"] = int(spec.GetMaxConcurrentReconciles())
	}
	if spec.S3WorkersLimit != nil {
		values["s3WorkersLimit"] = int(spec.GetS3WorkersLimit())
	}

	// ---- logging -----------------------------------------------------------
	// logStructured matches the chart default (false) — rendered only when
	// on. logLevel renders whenever the spec carries a value (the chart
	// default is "INFO"; re-stating it is harmless and keeps rendering
	// presence-driven, not value-driven).
	if spec.GetLog().GetStructured() {
		values["logStructured"] = true
	}
	if spec.GetLog() != nil && spec.GetLog().Level != nil && spec.GetLog().GetLevel() != "" {
		values["logLevel"] = spec.GetLog().GetLevel()
	}

	// ---- telemetry ----------------------------------------------------------
	// Chart default false (telemetry on) — rendered only on explicit
	// opt-out.
	if spec.GetDisableTelemetry() {
		values["disableTelemetry"] = true
	}

	// ---- leader election ------------------------------------------------------
	// The chart flattens these to four top-level values. Each renders on
	// presence: the enabled flag whenever the spec carries it (the chart
	// default is true, so an explicit true is a harmless re-statement),
	// the three timing knobs whenever non-empty.
	if le := spec.GetLeaderElection(); le != nil {
		if le.Enabled != nil {
			values["leaderElectionEnabled"] = le.GetEnabled()
		}
		if le.GetLeaseDuration() != "" {
			values["leaseDuration"] = le.GetLeaseDuration()
		}
		if le.GetRenewDeadline() != "" {
			values["renewDeadline"] = le.GetRenewDeadline()
		}
		if le.GetRetryPeriod() != "" {
			values["retryPeriod"] = le.GetRetryPeriod()
		}
	}

	// ---- feature gates -----------------------------------------------------
	// The chart folds every gate into one PXCO_FEATURE_GATES environment
	// variable; xtrabackupSidecar is the only gate it declares. Rendered
	// only when on (chart default false) — Helm deep-merges the map, so
	// the chart's own entry is replaced, not duplicated.
	if spec.GetXtrabackupSidecar() {
		values["featureGates"] = map[string]interface{}{
			"xtrabackupSidecar": true,
		}
	}

	// ---- scheduling ----------------------------------------------------------
	if len(spec.GetNodeSelector()) > 0 {
		values["nodeSelector"] = stringMapToInterface(spec.GetNodeSelector())
	}
	if len(spec.GetTolerations()) > 0 {
		values["tolerations"] = tolerationsSlice(spec.GetTolerations())
	}

	// ---- image ----------------------------------------------------------------
	// Pull secrets render as the raw Kubernetes object list
	// ([{name: ...}]) — the chart pipes imagePullSecrets straight into the
	// pod spec with toYaml.
	//
	// The image override is TWO chart values with a precedence rule: the
	// chart's image helper uses `image` verbatim when non-empty, else
	// "<operatorImageRepository>:<chart app version>". A repository alone
	// therefore maps to operatorImageRepository (tag stays the chart
	// version); any custom TAG requires the full `image` override, with
	// the repository half falling back to the chart's own default
	// repository when the spec leaves it empty.
	if len(spec.GetImagePullSecrets()) > 0 {
		pullSecrets := make([]interface{}, 0, len(spec.GetImagePullSecrets()))
		for _, name := range spec.GetImagePullSecrets() {
			pullSecrets = append(pullSecrets, map[string]interface{}{"name": name})
		}
		values["imagePullSecrets"] = pullSecrets
	}
	if spec.GetImage().GetRepository() != "" {
		values["operatorImageRepository"] = spec.GetImage().GetRepository()
	}
	if spec.GetImage().GetTag() != "" {
		repository := spec.GetImage().GetRepository()
		if repository == "" {
			repository = vars.DefaultImageRepository
		}
		values["image"] = fmt.Sprintf("%s:%s", repository, spec.GetImage().GetTag())
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
