package module

import (
	"github.com/pkg/errors"
	kubernetesprovider "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes"
	kuberneteskafkamirrormaker2v1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kuberneteskafkamirrormaker2/v1alpha1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apiextensions"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	kubernetesmeta "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// createMirrorMaker2 renders the kafka.strimzi.io/v1 KafkaMirrorMaker2
// resource as an UNTYPED CustomResource — the Strimzi CRDs type the
// cluster and connector `config` blocks with
// x-kubernetes-preserve-unknown-fields, which crd2pulumi cannot carry, so
// no generated package is shipped for the Kafka family (the same ruling
// as the KubernetesKafka module).
//
// ALIAS SEMANTICS: the target and every mirror source carry an alias that
// names the cluster in the replication flow. Under the default
// replication policy mirrored topics arrive PREFIXED with the source
// alias ("prod-msk.orders"); setting replication.policy.class to
// IdentityReplicationPolicy on a mirror's source AND checkpoint
// connectors keeps original names — the usual migration posture (the
// spec comments carry the full contract).
//
// The spec body built here is the exact twin of the Terraform module's
// local.mirrormaker2_manifest (locals.tf) — same keys rendered and
// omitted, numbers as ints, booleans as booleans. An unset optional is
// simply never inserted into the map (the Go twin of TF's null-prune).
//
// No await machinery, deliberately: engine readiness depends on the
// operator (image pulls, worker group formation, connector startup) that
// is not part of applying the resource — the same
// never-block-on-a-controller posture as the sibling operator-CR kinds.
func createMirrorMaker2(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependencies []pulumi.ResourceOption,
) (pulumi.Resource, error) {
	spec := locals.Spec

	// Twin of TF's coalesce(replicas, 1): the CRD requires a worker
	// count and the spec defaults to one.
	replicas := 1
	if spec.Replicas != nil {
		replicas = int(spec.GetReplicas())
	}

	specBody := map[string]interface{}{
		"replicas": replicas,
		"target":   targetBody(locals),
		"mirrors":  mirrorsBody(spec.GetMirrors()),
	}
	if spec.GetVersion() != "" {
		specBody["version"] = spec.GetVersion()
	}
	if resources := resourcesMap(spec.GetResources()); resources != nil {
		specBody["resources"] = resources
	}
	if jvm := jvmBody(spec.GetJvm()); jvm != nil {
		specBody["jvmOptions"] = jvm
	}
	if spec.GetRack().GetTopologyKey() != "" {
		specBody["rack"] = map[string]interface{}{
			"topologyKey": spec.GetRack().GetTopologyKey(),
		}
	}
	if spec.GetMetrics().GetEnabled() {
		// The module owns the rules ConfigMap (metricsConfigMap below);
		// the CR only points at it.
		specBody["metricsConfig"] = map[string]interface{}{
			"type": "jmxPrometheusExporter",
			"valueFrom": map[string]interface{}{
				"configMapKeyRef": map[string]interface{}{
					"name": locals.MetricsConfigMapName,
					"key":  vars.MetricsConfigKey,
				},
			},
		}
	}
	if template := podTemplateBody(spec); template != nil {
		specBody["template"] = template
	}

	return apiextensions.NewCustomResource(ctx, locals.MirrorMakerName,
		&apiextensions.CustomResourceArgs{
			ApiVersion: pulumi.String(vars.ApiVersion),
			Kind:       pulumi.String("KafkaMirrorMaker2"),
			Metadata: &kubernetesmeta.ObjectMetaArgs{
				Name:      pulumi.String(locals.MirrorMakerName),
				Namespace: pulumi.String(locals.Namespace),
				Labels:    pulumi.ToStringMap(locals.Labels),
			},
			OtherFields: kubernetes.UntypedArgs{
				"spec": specBody,
			},
		}, append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, dependencies...)...)
}

// metricsConfigMap renders the module-owned JMX Prometheus Exporter rules
// ConfigMap when spec.metrics.enabled — the canonical Strimzi connect
// rule set (metrics_rules.go) under the key the CR's metricsConfig points
// at. Returns nil when metrics are disabled. Terraform equivalent:
// kubernetes_config_map_v1 with count.
func metricsConfigMap(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependencies []pulumi.ResourceOption,
) (pulumi.Resource, error) {
	if !locals.Spec.GetMetrics().GetEnabled() {
		return nil, nil
	}

	createdConfigMap, err := kubernetescorev1.NewConfigMap(ctx, locals.MetricsConfigMapName,
		&kubernetescorev1.ConfigMapArgs{
			Metadata: kubernetesmeta.ObjectMetaPtrInput(
				&kubernetesmeta.ObjectMetaArgs{
					Name:      pulumi.String(locals.MetricsConfigMapName),
					Namespace: pulumi.String(locals.Namespace),
					Labels:    pulumi.ToStringMap(locals.Labels),
				}),
			Data: pulumi.StringMap{
				vars.MetricsConfigKey: pulumi.String(mirrorMaker2MetricsRules),
			},
		}, append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, dependencies...)...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create mirror maker 2 metrics ConfigMap")
	}

	return createdConfigMap, nil
}

// targetBody renders the TARGET cluster — where mirrored data lands and
// where the Connect-style engine keeps its state. The group-identity
// values arrive pre-resolved on locals (spec overrides with
// metadata.name-derived fallbacks).
func targetBody(locals *Locals) map[string]interface{} {
	target := locals.Spec.GetTarget()

	body := map[string]interface{}{
		"alias":              locals.TargetAlias,
		"bootstrapServers":   target.BootstrapServers.GetValue(),
		"groupId":            locals.GroupId,
		"configStorageTopic": locals.ConfigStorageTopic,
		"statusStorageTopic": locals.StatusStorageTopic,
		"offsetStorageTopic": locals.OffsetStorageTopic,
	}
	if tls := clientTlsBody(target.GetTls()); tls != nil {
		body["tls"] = tls
	}
	if authentication := clientAuthenticationBody(target.GetAuthentication()); authentication != nil {
		body["authentication"] = authentication
	}
	if len(target.GetConfig()) > 0 {
		body["config"] = stringMapToInterface(target.GetConfig())
	}
	return body
}

// mirrorsBody renders one entry per declared mirror: the source cluster
// connection, the replication scope patterns, and the per-mirror
// MirrorSourceConnector / MirrorCheckpointConnector tuning.
func mirrorsBody(mirrors []*kuberneteskafkamirrormaker2v1alpha1.KubernetesKafkaMirrorMaker2Mirror) []interface{} {
	out := make([]interface{}, 0, len(mirrors))
	for _, mirror := range mirrors {
		body := map[string]interface{}{
			"source": sourceBody(mirror.GetSource()),
		}
		if mirror.GetTopicsPattern() != "" {
			body["topicsPattern"] = mirror.GetTopicsPattern()
		}
		if mirror.GetTopicsExcludePattern() != "" {
			body["topicsExcludePattern"] = mirror.GetTopicsExcludePattern()
		}
		if mirror.GetGroupsPattern() != "" {
			body["groupsPattern"] = mirror.GetGroupsPattern()
		}
		if mirror.GetGroupsExcludePattern() != "" {
			body["groupsExcludePattern"] = mirror.GetGroupsExcludePattern()
		}
		if connector := connectorBody(mirror.GetSourceConnector()); connector != nil {
			body["sourceConnector"] = connector
		}
		if connector := connectorBody(mirror.GetCheckpointConnector()); connector != nil {
			body["checkpointConnector"] = connector
		}
		out = append(out, body)
	}
	return out
}

func sourceBody(source *kuberneteskafkamirrormaker2v1alpha1.KubernetesKafkaMirrorMaker2Source) map[string]interface{} {
	body := map[string]interface{}{
		"alias":            source.GetAlias(),
		"bootstrapServers": source.BootstrapServers.GetValue(),
	}
	if tls := clientTlsBody(source.GetTls()); tls != nil {
		body["tls"] = tls
	}
	if authentication := clientAuthenticationBody(source.GetAuthentication()); authentication != nil {
		body["authentication"] = authentication
	}
	if len(source.GetConfig()) > 0 {
		body["config"] = stringMapToInterface(source.GetConfig())
	}
	return body
}

func connectorBody(connector *kuberneteskafkamirrormaker2v1alpha1.KubernetesKafkaMirrorMaker2Connector) map[string]interface{} {
	if connector == nil {
		return nil
	}
	body := map[string]interface{}{}
	if connector.TasksMax != nil {
		body["tasksMax"] = int(connector.GetTasksMax())
	}
	if len(connector.GetConfig()) > 0 {
		body["config"] = stringMapToInterface(connector.GetConfig())
	}
	if autoRestart := connector.GetAutoRestart(); autoRestart != nil {
		autoRestartBody := map[string]interface{}{
			"enabled": autoRestart.GetEnabled(),
		}
		if autoRestart.MaxRestarts != nil {
			autoRestartBody["maxRestarts"] = int(autoRestart.GetMaxRestarts())
		}
		body["autoRestart"] = autoRestartBody
	}
	return body
}

// clientTlsBody renders the shared StrimziKafkaClientTls message: the
// certificates the CLIENT trusts when verifying brokers, each naming a
// Secret with either one file (certificate) or a glob (pattern) — the
// spec enforces exactly one of the two.
func clientTlsBody(tls *kubernetesprovider.StrimziKafkaClientTls) map[string]interface{} {
	if tls == nil {
		return nil
	}
	certificates := make([]interface{}, 0, len(tls.GetTrustedCertificates()))
	for _, certificate := range tls.GetTrustedCertificates() {
		entry := map[string]interface{}{
			"secretName": certificate.GetSecretName().GetValue(),
		}
		if certificate.GetCertificate() != "" {
			entry["certificate"] = certificate.GetCertificate()
		}
		if certificate.GetPattern() != "" {
			entry["pattern"] = certificate.GetPattern()
		}
		certificates = append(certificates, entry)
	}
	return map[string]interface{}{
		"trustedCertificates": certificates,
	}
}

// clientAuthenticationBody renders the shared
// StrimziKafkaClientAuthentication message. Each type carries only its
// own credential shape (the spec's CEL rules guarantee the referenced
// fields are present): tls = the client certificate the workload
// presents; the SASL trio = username + password Secret reference; custom
// = bring-your-own mechanism via sasl + config.
func clientAuthenticationBody(authentication *kubernetesprovider.StrimziKafkaClientAuthentication) map[string]interface{} {
	if authentication == nil {
		return nil
	}
	body := map[string]interface{}{
		"type": authentication.GetType(),
	}
	switch authentication.GetType() {
	case "tls":
		certificateAndKey := authentication.GetCertificateAndKey()
		// KubernetesKafkaUser credential Secrets carry user.crt/user.key;
		// the spec defaults mirror that, so unset keys resolve through
		// the fallbacks (twin of TF's coalesce).
		certificate := certificateAndKey.GetCertificate()
		if certificate == "" {
			certificate = "user.crt"
		}
		key := certificateAndKey.GetKey()
		if key == "" {
			key = "user.key"
		}
		body["certificateAndKey"] = map[string]interface{}{
			"secretName":  certificateAndKey.GetSecretName().GetValue(),
			"certificate": certificate,
			"key":         key,
		}
	case "scram-sha-512", "scram-sha-256", "plain":
		body["username"] = authentication.GetUsername()
		// KubernetesKafkaUser credential Secrets carry the password
		// under "password"; the spec default mirrors that.
		password := authentication.GetPasswordSecret().GetPassword()
		if password == "" {
			password = "password"
		}
		body["passwordSecret"] = map[string]interface{}{
			"secretName": authentication.GetPasswordSecret().GetSecretName().GetValue(),
			"password":   password,
		}
	case "custom":
		body["sasl"] = authentication.GetSasl()
		if len(authentication.GetConfig()) > 0 {
			body["config"] = stringMapToInterface(authentication.GetConfig())
		}
	}
	return body
}

// podTemplateBody renders the worker scheduling knobs into Strimzi's pod
// template. The Strimzi pod template carries affinity and tolerations but
// NO nodeSelector — a node_selector map therefore translates to a
// requiredDuringSchedulingIgnoredDuringExecution nodeAffinity with one
// matchExpressions entry per label (semantically identical for exact-match
// selection; the Terraform module renders the same translation).
func podTemplateBody(spec *kuberneteskafkamirrormaker2v1alpha1.KubernetesKafkaMirrorMaker2Spec) map[string]interface{} {
	podBody := map[string]interface{}{}

	if len(spec.GetTolerations()) > 0 {
		podBody["tolerations"] = tolerationsSlice(spec.GetTolerations())
	}

	if len(spec.GetNodeSelector()) > 0 {
		matchExpressions := make([]interface{}, 0, len(spec.GetNodeSelector()))
		// Sorted iteration keeps the rendered CR deterministic across
		// runs (Go map order is random; TF for-expressions sort keys).
		for _, key := range sortedKeys(spec.GetNodeSelector()) {
			matchExpressions = append(matchExpressions, map[string]interface{}{
				"key":      key,
				"operator": "In",
				"values":   []interface{}{spec.GetNodeSelector()[key]},
			})
		}
		podBody["affinity"] = map[string]interface{}{
			"nodeAffinity": map[string]interface{}{
				"requiredDuringSchedulingIgnoredDuringExecution": map[string]interface{}{
					"nodeSelectorTerms": []interface{}{
						map[string]interface{}{
							"matchExpressions": matchExpressions,
						},
					},
				},
			},
		}
	}

	if len(podBody) == 0 {
		return nil
	}
	return map[string]interface{}{
		"pod": podBody,
	}
}

func jvmBody(jvm *kuberneteskafkamirrormaker2v1alpha1.KubernetesKafkaMirrorMaker2Jvm) map[string]interface{} {
	if jvm == nil {
		return nil
	}
	body := map[string]interface{}{}
	if jvm.GetXms() != "" {
		body["-Xms"] = jvm.GetXms()
	}
	if jvm.GetXmx() != "" {
		body["-Xmx"] = jvm.GetXmx()
	}
	if len(body) == 0 {
		return nil
	}
	return body
}

// resourcesMap renders the shared ContainerResources message into the
// standard Kubernetes limits/requests shape. Returns nil when nothing is
// set.
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

func sortedKeys(in map[string]string) []string {
	keys := make([]string, 0, len(in))
	for k := range in {
		keys = append(keys, k)
	}
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[j] < keys[i] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}
	return keys
}

func stringMapToInterface(in map[string]string) map[string]interface{} {
	out := make(map[string]interface{}, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
