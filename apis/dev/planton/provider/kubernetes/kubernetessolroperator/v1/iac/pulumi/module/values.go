package module

import (
	"strings"

	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
)

// buildHelmValues renders the typed spec into the solr-operator chart's
// values map, then merges the spec's helm_values escape hatch over it with
// Helm `-f` semantics (maps deep-merge with the later document winning,
// lists replace).
//
// PARITY: the Terraform module reaches the same result natively — its
// helm_release passes values = [yamlencode(typed values), helm_values] and
// the provider merges the documents in exactly this order. Keep every typed
// mapping below in lockstep with the Terraform module's locals.
//
// Chart-default-matching values render only on divergence (the watch
// scope, the true-defaulted toggles, the image source), so the rendered
// values stay minimal on both engines. The ONE always-rendered value is
// `zookeeper-operator.crd.create: false` — the module owns the
// ZookeeperCluster CRD (applied with the keep-on-uninstall posture), so
// the bundled subchart must NEVER install its own copy and put it under
// Helm's delete-on-uninstall lifecycle.
func buildHelmValues(locals *Locals) (map[string]interface{}, error) {
	spec := locals.Spec

	values := map[string]interface{}{}

	// ---- operator sizing -------------------------------------------------
	if spec.Replicas != nil {
		values["replicaCount"] = int(spec.GetReplicas())
	}
	// The chart ships NO default resources (`resources: {}` — the
	// operator is lightweight); the key renders only when the spec sets
	// it.
	if r := resourcesMap(spec.GetResources()); r != nil {
		values["resources"] = r
	}

	// ---- watch scope -----------------------------------------------------
	// The chart takes a COMMA-JOINED string (templates/_helpers.tpl
	// splits it back apart), not a YAML list. Empty = the operator
	// watches ALL namespaces (the chart default; also the chart's
	// ClusterRole-vs-Role switch).
	if len(spec.GetWatchNamespaces()) > 0 {
		values["watchNamespaces"] = strings.Join(spec.GetWatchNamespaces(), ",")
	}

	// ---- bundled zookeeper-operator ----------------------------------------
	// NOTE the dash in the values key: "zookeeper-operator" is the
	// SUBCHART name, addressed as a plain map key. crd.create is pinned
	// false unconditionally — see the function comment. install renders
	// on presence (chart default true); use renders only on divergence
	// (chart default false, and it is ignored whenever install is true).
	zookeeperOperator := map[string]interface{}{
		"crd": map[string]interface{}{
			"create": false,
		},
	}
	if zk := spec.GetZookeeperOperator(); zk != nil {
		if zk.Install != nil {
			zookeeperOperator["install"] = zk.GetInstall()
		}
		if zk.GetUseExisting() {
			zookeeperOperator["use"] = true
		}
	}
	values["zookeeper-operator"] = zookeeperOperator

	// ---- leader election / metrics -------------------------------------------
	// Both nest a single enable flag and default true in the chart.
	// Rendered on presence — an explicit true re-states the chart
	// default harmlessly, an explicit false is the actual opt-out.
	if spec.LeaderElectionEnabled != nil {
		values["leaderElection"] = map[string]interface{}{
			"enable": spec.GetLeaderElectionEnabled(),
		}
	}
	if spec.MetricsEnabled != nil {
		values["metrics"] = map[string]interface{}{
			"enable": spec.GetMetricsEnabled(),
		}
	}

	// ---- operator -> Solr mutual TLS ----------------------------------------
	// Secret names resolve through the value-or-ref (a literal name or a
	// KubernetesSecret reference the platform resolved before the module
	// ran). Scalars with chart defaults (caCertSecretKey,
	// insecureSkipVerify, watchForUpdates) render on presence.
	if m := spec.GetMtls(); m != nil {
		mtls := map[string]interface{}{}
		if m.GetClientCertSecret().GetValue() != "" {
			mtls["clientCertSecret"] = m.GetClientCertSecret().GetValue()
		}
		if m.GetCaCertSecret().GetValue() != "" {
			mtls["caCertSecret"] = m.GetCaCertSecret().GetValue()
		}
		if m.CaCertSecretKey != nil {
			mtls["caCertSecretKey"] = m.GetCaCertSecretKey()
		}
		if m.InsecureSkipVerify != nil {
			mtls["insecureSkipVerify"] = m.GetInsecureSkipVerify()
		}
		if m.WatchForUpdates != nil {
			mtls["watchForUpdates"] = m.GetWatchForUpdates()
		}
		if len(mtls) > 0 {
			values["mTLS"] = mtls
		}
	}

	// ---- scheduling ----------------------------------------------------------
	if len(spec.GetNodeSelector()) > 0 {
		values["nodeSelector"] = stringMapToInterface(spec.GetNodeSelector())
	}
	if len(spec.GetTolerations()) > 0 {
		values["tolerations"] = tolerationsSlice(spec.GetTolerations())
	}

	// ---- image source ----------------------------------------------------------
	// image.repository / image.tag are the air-gap override; the pull
	// secret is the chart's image.imagePullSecret — a SINGULAR string
	// (the chart accepts exactly one), not the usual imagePullSecrets
	// object list.
	image := map[string]interface{}{}
	if spec.GetImage().GetRepository() != "" {
		image["repository"] = spec.GetImage().GetRepository()
	}
	if spec.GetImage().GetTag() != "" {
		image["tag"] = spec.GetImage().GetTag()
	}
	if spec.GetImagePullSecret() != "" {
		image["imagePullSecret"] = spec.GetImagePullSecret()
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
