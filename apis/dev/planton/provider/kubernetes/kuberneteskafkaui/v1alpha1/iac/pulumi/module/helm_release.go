package module

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	kubernetesprovider "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes"
	kuberneteskafkauiv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kuberneteskafkaui/v1alpha1"
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"sigs.k8s.io/yaml"
)

// THE PLACEHOLDER / secretMappings MECHANISM (why no secret ever lands in
// rendered values):
//
// The chart writes yamlApplicationConfig verbatim into a ConfigMap
// (config.yml) — anything rendered there is world-readable to anyone who
// can read ConfigMaps, ends up in Helm release history, and in both
// engines' state files. So the module NEVER renders a credential value
// into the application config. Instead every password position carries a
// ${ENV_VAR} placeholder, and the chart's envs.secretMappings values wire
// each of those env vars to a Kubernetes Secret key (the deployment
// template renders them as valueFrom.secretKeyRef entries). The kafbat UI
// is a Spring Boot app: Spring's property resolution expands ${ENV_VAR}
// placeholders in the mounted config.yml against the container's
// environment at startup — the credential exists only inside the running
// container.
//
// Two kinds of credentials feed the mappings:
//   - REFERENCED credentials (cluster sasl, schema registry, Connect
//     password_secret entries): already live in a Secret; the mapping
//     points straight at that Secret/key — the module copies nothing.
//   - LITERAL credentials (the declared console login password): the
//     module materializes them into ONE Secret, "<name>-secrets"
//     (secret.go), and maps from it.
//
// Env var naming is deterministic and index-based so both engines emit the
// identical placeholders:
//   KAFKA_CLUSTER_<i>_PASSWORD                     — cluster i sasl
//   KAFKA_CLUSTER_<i>_SCHEMA_REGISTRY_PASSWORD     — cluster i registry
//   KAFKA_CLUSTER_<i>_CONNECT_<j>_PASSWORD         — cluster i, connect j
//   KAFKA_UI_USER_PASSWORD                         — console login

// consolePasswordEnvVar carries the console login password (LOGIN_FORM)
// into spring.security.user.password via its ${...} placeholder.
const consolePasswordEnvVar = "KAFKA_UI_USER_PASSWORD"

// consolePasswordSecretKey is the key inside the module-materialized
// "<name>-secrets" Secret holding the console login password — must match
// secret.go.
const consolePasswordSecretKey = "console-user-password"

func clusterPasswordEnvVar(i int) string {
	return fmt.Sprintf("KAFKA_CLUSTER_%d_PASSWORD", i)
}

func schemaRegistryPasswordEnvVar(i int) string {
	return fmt.Sprintf("KAFKA_CLUSTER_%d_SCHEMA_REGISTRY_PASSWORD", i)
}

func connectPasswordEnvVar(i, j int) string {
	return fmt.Sprintf("KAFKA_CLUSTER_%d_CONNECT_%d_PASSWORD", i, j)
}

// caMountPath is where cluster i's CA-certificate Secret volume mounts —
// must match the volumes/volumeMounts rendering below.
func caMountPath(i int) string {
	return fmt.Sprintf("/etc/kafkaui/cluster-%d-ca", i)
}

// helmRelease installs the kafka-ui chart as a real Helm release with the
// typed spec rendered into chart values and the helm_values escape hatch
// merged LAST (Helm -f semantics) — the exact semantic twin of the
// Terraform module's helm_release with values = [typed, helm_values].
func helmRelease(ctx *pulumi.Context,
	locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependsOn []pulumi.ResourceOption,
) error {
	mergedValues, err := buildHelmValues(locals)
	if err != nil {
		return errors.Wrap(err, "failed to build helm values")
	}

	releaseArgs := &helmv3.ReleaseArgs{
		Name:      pulumi.String(locals.ReleaseName),
		Namespace: pulumi.String(locals.Namespace),
		Chart:     pulumi.String(vars.HelmChartName),
		Version:   pulumi.String(locals.ChartVersion),
		RepositoryOpts: &helmv3.RepositoryOptsArgs{
			Repo: pulumi.String(vars.HelmChartRepo),
		},
		Values: pulumi.ToMap(mergedValues),
		// The module owns namespace creation (create_namespace flag).
		CreateNamespace: pulumi.Bool(false),
		// Wait for the console to become Ready — a console that never
		// starts (bad image, unschedulable pod, unresolvable cluster
		// config) should fail THIS deploy, not the first browser hit.
		// SkipAwait false is Helm --wait, stated explicitly to mirror
		// the Terraform twin's `wait = true`.
		SkipAwait:     pulumi.Bool(false),
		Atomic:        pulumi.Bool(true),
		CleanupOnFail: pulumi.Bool(true),
		Timeout:       pulumi.Int(600),
	}

	opts := append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, dependsOn...)

	_, err = helmv3.NewRelease(ctx, locals.ReleaseName, releaseArgs, opts...)
	if err != nil {
		return errors.Wrap(err, "failed to install kafka-ui helm release")
	}
	return nil
}

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

	values := map[string]interface{}{
		// Always rendered with resolved defaults so both engines emit the
		// identical documents whether or not the platform's defaulting
		// middleware ran.
		// Pins the chart's fullname to the resource name (see
		// locals.ServiceName).
		"fullnameOverride": locals.ReleaseName,
		"replicaCount":     locals.Replicas,
		"service": map[string]interface{}{
			"type": locals.ServiceType,
			"port": locals.ServicePort,
		},
		"yamlApplicationConfig": yamlApplicationConfig(locals),
	}

	if r := resourcesMap(spec.GetResources()); r != nil {
		values["resources"] = r
	}
	if len(spec.GetNodeSelector()) > 0 {
		values["nodeSelector"] = stringMapToInterface(spec.GetNodeSelector())
	}
	if len(spec.GetTolerations()) > 0 {
		values["tolerations"] = tolerationsSlice(spec.GetTolerations())
	}
	// Chart default is ghcr.io/kafbat/kafka-ui at the chart appVersion —
	// only the registry seam is typed (air-gapped mirrors).
	if spec.GetImageRegistry() != "" {
		values["image"] = map[string]interface{}{
			"registry": spec.GetImageRegistry(),
		}
	}

	if mappings := secretMappings(locals); len(mappings) > 0 {
		values["envs"] = map[string]interface{}{
			"secretMappings": mappings,
		}
	}

	// One secret volume per TLS cluster, mounted where the rendered
	// ssl.truststore.location points (the chart passes volumes /
	// volumeMounts through to the Deployment verbatim).
	if volumes := tlsVolumes(spec); len(volumes) > 0 {
		values["volumes"] = volumes
		values["volumeMounts"] = tlsVolumeMounts(spec)
	}

	// ---- escape hatch (merged LAST, helm -f semantics) -------------------
	if spec.GetHelmValues() != "" {
		overrides := map[string]interface{}{}
		if err := yaml.Unmarshal([]byte(spec.GetHelmValues()), &overrides); err != nil {
			return nil, errors.Wrap(err, "failed to parse helm_values as a YAML document")
		}
		values = mergeMaps(values, overrides)
	}

	return values, nil
}

// yamlApplicationConfig renders the typed cluster wiring into the app's
// config document (ClustersProperties shape: kafka.clusters[], auth.type,
// spring.security.user). The chart mounts it as /kafka-ui/config.yml and
// points SPRING_CONFIG_ADDITIONAL-LOCATION at it.
func yamlApplicationConfig(locals *Locals) map[string]interface{} {
	spec := locals.Spec

	config := map[string]interface{}{
		"kafka": map[string]interface{}{
			"clusters": clusterEntries(spec),
		},
	}

	if locals.AuthEnabled {
		config["auth"] = map[string]interface{}{"type": "LOGIN_FORM"}
		// LOGIN_FORM rides Spring Boot's DEFAULT security user — the app
		// (io.kafbat.ui.config.auth.BasicAuthSecurityConfig) registers no
		// user store of its own, so exactly ONE account exists:
		// spring.security.user.name/password — which is why the spec
		// models a single `user` (multi-user login needs LDAP/OAuth2
		// through helm_values). The password is the ${...} placeholder —
		// the literal lives in the "<name>-secrets" Secret (secret.go).
		config["spring"] = map[string]interface{}{
			"security": map[string]interface{}{
				"user": map[string]interface{}{
					"name":     spec.GetAuth().GetUser().GetUsername(),
					"password": "${" + consolePasswordEnvVar + "}",
				},
			},
		}
	} else {
		config["auth"] = map[string]interface{}{"type": "DISABLED"}
	}

	return config
}

// clusterEntries renders spec.clusters into the app's kafka.clusters list.
func clusterEntries(spec *kuberneteskafkauiv1alpha1.KubernetesKafkaUiSpec) []interface{} {
	entries := make([]interface{}, 0, len(spec.GetClusters()))
	for i, cluster := range spec.GetClusters() {
		entry := map[string]interface{}{
			"name":             cluster.GetName(),
			"bootstrapServers": cluster.GetBootstrapServers().GetValue(),
		}

		// readOnly hides every mutating console action for this cluster
		// (topic create/delete, message produce, config edits) — an
		// app-side switch, not a Kafka ACL: the right posture for
		// production clusters on a shared console. Rendered only when
		// true (the app default is false).
		if cluster.GetReadOnly() {
			entry["readOnly"] = true
		}

		if props := clusterProperties(i, cluster); len(props) > 0 {
			entry["properties"] = props
		}

		if sr := cluster.GetSchemaRegistry(); sr != nil {
			entry["schemaRegistry"] = sr.GetUrl().GetValue()
			auth := map[string]interface{}{}
			if sr.GetUsername() != "" {
				auth["username"] = sr.GetUsername()
			}
			if sr.GetPasswordSecret() != nil {
				auth["password"] = "${" + schemaRegistryPasswordEnvVar(i) + "}"
			}
			if len(auth) > 0 {
				entry["schemaRegistryAuth"] = auth
			}
		}

		if len(cluster.GetKafkaConnect()) > 0 {
			connects := make([]interface{}, 0, len(cluster.GetKafkaConnect()))
			for j, kc := range cluster.GetKafkaConnect() {
				connect := map[string]interface{}{
					"name":    kc.GetName(),
					"address": kc.GetAddress().GetValue(),
				}
				if kc.GetUsername() != "" {
					connect["username"] = kc.GetUsername()
				}
				if kc.GetPasswordSecret() != nil {
					connect["password"] = "${" + connectPasswordEnvVar(i, j) + "}"
				}
				connects = append(connects, connect)
			}
			entry["kafkaConnect"] = connects
		}

		entries = append(entries, entry)
	}
	return entries
}

// clusterProperties merges the user's Kafka client properties with the
// module-owned security properties derived from the typed tls/sasl blocks
// (module-owned entries win — the spec forbids credentials in properties).
//
// THE PEM TRUSTSTORE TRICK: Kafka clients since KIP-651 accept
// ssl.truststore.type=PEM with ssl.truststore.location pointing at a plain
// PEM certificate file — no JKS/PKCS12 conversion, no truststore password.
// The module mounts the CA Secret as-is (tlsVolumes) and points the
// truststore at the mounted key, so a Strimzi cluster-CA Secret works
// directly.
func clusterProperties(i int, cluster *kuberneteskafkauiv1alpha1.KubernetesKafkaUiCluster) map[string]interface{} {
	props := map[string]interface{}{}
	for k, v := range cluster.GetProperties() {
		props[k] = v
	}

	tls := cluster.GetTls()
	sasl := cluster.GetSasl()

	if tls != nil {
		securityProtocol := "SSL"
		if sasl != nil {
			securityProtocol = "SASL_SSL"
		}
		caCertificate := tls.GetCaCertificate()
		if caCertificate == "" {
			caCertificate = "ca.crt"
		}
		props["security.protocol"] = securityProtocol
		props["ssl.truststore.type"] = "PEM"
		props["ssl.truststore.location"] = caMountPath(i) + "/" + caCertificate
	}

	if sasl != nil {
		securityProtocol := "SASL_PLAINTEXT"
		if tls != nil {
			securityProtocol = "SASL_SSL"
		}
		props["security.protocol"] = securityProtocol
		props["sasl.mechanism"] = sasl.GetMechanism()
		props["sasl.jaas.config"] = saslJaasConfig(i, sasl)
	}

	return props
}

// saslJaasConfig renders the JAAS login-module line: ScramLoginModule for
// the SCRAM-* mechanisms, PlainLoginModule for PLAIN (the spec CEL rule
// admits nothing else). The username is inline (not sensitive); the
// password is the ${...} placeholder Spring resolves from the
// Secret-backed env var — never a literal.
func saslJaasConfig(i int, sasl *kuberneteskafkauiv1alpha1.KubernetesKafkaUiClusterSasl) string {
	loginModule := "org.apache.kafka.common.security.plain.PlainLoginModule"
	if strings.HasPrefix(sasl.GetMechanism(), "SCRAM") {
		loginModule = "org.apache.kafka.common.security.scram.ScramLoginModule"
	}
	return fmt.Sprintf(`%s required username="%s" password="${%s}";`,
		loginModule, sasl.GetUsername(), clusterPasswordEnvVar(i))
}

// secretMappings assembles the chart's envs.secretMappings map — one entry
// per ${...} placeholder rendered anywhere in yamlApplicationConfig, each
// pointing at its source Secret/key (see the mechanism comment at the top
// of this file).
func secretMappings(locals *Locals) map[string]interface{} {
	spec := locals.Spec
	mappings := map[string]interface{}{}

	for i, cluster := range spec.GetClusters() {
		if sasl := cluster.GetSasl(); sasl != nil {
			mappings[clusterPasswordEnvVar(i)] = secretMapping(sasl.GetPasswordSecret())
		}
		if sr := cluster.GetSchemaRegistry(); sr != nil && sr.GetPasswordSecret() != nil {
			mappings[schemaRegistryPasswordEnvVar(i)] = secretMapping(sr.GetPasswordSecret())
		}
		for j, kc := range cluster.GetKafkaConnect() {
			if kc.GetPasswordSecret() != nil {
				mappings[connectPasswordEnvVar(i, j)] = secretMapping(kc.GetPasswordSecret())
			}
		}
	}

	if locals.AuthEnabled {
		mappings[consolePasswordEnvVar] = map[string]interface{}{
			"name":    locals.ConsoleSecretName,
			"keyName": consolePasswordSecretKey,
		}
	}

	return mappings
}

// secretMapping renders one password_secret reference into the chart's
// {name, keyName} mapping shape, resolving the key to the spec default
// ("password" — the KubernetesKafkaUser credential-Secret layout).
func secretMapping(passwordSecret *kuberneteskafkauiv1alpha1.KubernetesKafkaUiPasswordSecret) map[string]interface{} {
	key := passwordSecret.GetKey()
	if key == "" {
		key = "password"
	}
	return map[string]interface{}{
		"name":    passwordSecret.GetSecretName().GetValue(),
		"keyName": key,
	}
}

// tlsVolumes renders one secret volume per TLS cluster ("cluster-<i>-ca",
// index-named so entries stay stable and unique across clusters).
func tlsVolumes(spec *kuberneteskafkauiv1alpha1.KubernetesKafkaUiSpec) []interface{} {
	var volumes []interface{}
	for i, cluster := range spec.GetClusters() {
		if tls := cluster.GetTls(); tls != nil {
			volumes = append(volumes, map[string]interface{}{
				"name": fmt.Sprintf("cluster-%d-ca", i),
				"secret": map[string]interface{}{
					"secretName": tls.GetCaSecretName().GetValue(),
				},
			})
		}
	}
	return volumes
}

// tlsVolumeMounts mounts each TLS cluster's CA volume at the path the
// rendered ssl.truststore.location points into.
func tlsVolumeMounts(spec *kuberneteskafkauiv1alpha1.KubernetesKafkaUiSpec) []interface{} {
	var mounts []interface{}
	for i, cluster := range spec.GetClusters() {
		if cluster.GetTls() != nil {
			mounts = append(mounts, map[string]interface{}{
				"name":      fmt.Sprintf("cluster-%d-ca", i),
				"mountPath": caMountPath(i),
				"readOnly":  true,
			})
		}
	}
	return mounts
}

// ---- shared shape helpers (twins of the Terraform locals) ----------------

// resourcesMap renders the shared ContainerResources message into the
// chart's resources shape. Returns nil when nothing is set.
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

// tolerationsSlice renders the shared WorkloadToleration list into the
// chart's tolerations shape.
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

// mergeMaps deep-merges b over a with Helm's `-f` semantics: nested maps
// merge recursively with b winning per key; everything else (scalars,
// lists) is replaced by b's value.
func mergeMaps(a, b map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(a))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		if bChild, ok := v.(map[string]interface{}); ok {
			if aChild, ok := out[k].(map[string]interface{}); ok {
				out[k] = mergeMaps(aChild, bChild)
				continue
			}
		}
		out[k] = v
	}
	return out
}

// stringMapToInterface converts a map[string]string into the
// map[string]interface{} YAML rendering expects.
func stringMapToInterface(in map[string]string) map[string]interface{} {
	out := make(map[string]interface{}, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
