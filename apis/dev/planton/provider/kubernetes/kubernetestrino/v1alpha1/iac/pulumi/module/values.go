package module

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	kubernetesprovider "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes"
	kubernetestrinov1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetestrino/v1alpha1"
	"sigs.k8s.io/yaml"
)

// buildHelmValues renders the typed spec into the chart's values map, then
// merges the spec's helm_values escape hatch over it with Helm `-f`
// semantics (maps deep-merge with the later document winning, lists
// replace), then RE-PINS the security-critical values.
//
// SECRET DISCIPLINE (load-bearing): every properties surface in this
// chart — config.properties, every catalog, resource groups, the
// exchange manager — renders into a ConfigMap. Nothing rendered here
// carries credential material: catalog passwords and the internal
// shared secret are `${ENV:VAR}` references (Trino's own secrets
// substitution, docs/security/secrets.md at the pin) resolved from
// Secret-sourced environment variables at process start.
//
// PARITY: the Terraform module reaches the same result natively — its
// helm_release passes values = [yamlencode(typed values), helm_values,
// yamlencode(re-pins)] and the provider merges the documents in exactly
// this order. Keep every typed mapping below in lockstep with the
// Terraform module's locals.
func buildHelmValues(locals *Locals) (map[string]interface{}, error) {
	spec := locals.Spec

	// ------------------------- name-budget fail-loud ----------------------
	// Chart truth at the pin: `<fullname>-schemas-volume-coordinator`
	// (27-char suffix) renders unconditionally; the resource-groups
	// ConfigMap suffix is 36 chars but renders only when resource
	// groups are declared. Both must fit the 63-char DNS bound.
	if spec.GetResourceGroupsConfig() != "" && len(locals.ReleaseName) > vars.NameBudgetResourceGroups {
		return nil, errors.Errorf(
			"metadata.name %q is %d characters — with resource_groups_config set the chart derives `<name>-resource-groups-volume-coordinator` (36-char suffix), so the name must be at most %d characters",
			locals.ReleaseName, len(locals.ReleaseName), vars.NameBudgetResourceGroups)
	}
	if len(locals.ReleaseName) > vars.NameBudget {
		return nil, errors.Errorf(
			"metadata.name %q is %d characters — the chart derives `<name>-schemas-volume-coordinator` (27-char suffix), so the name must be at most %d characters",
			locals.ReleaseName, len(locals.ReleaseName), vars.NameBudget)
	}

	values := map[string]interface{}{
		// Deterministic child names (`<name>-coordinator`, …) — the
		// release name never double-prefixes and the import map stays
		// exact (the fullname re-pin discipline).
		"fullnameOverride": locals.ReleaseName,
	}

	// ------------------------------- image --------------------------------
	// The chart's SPLIT image form: registry (empty = Docker Hub) +
	// repository + tag (empty = appVersion). The module always pins
	// the tag explicitly so the deployed engine version is declared,
	// never inherited.
	imageBlock := map[string]interface{}{
		"repository": defaultString(spec.GetImage().GetRepository(), "trinodb/trino"),
		"tag":        defaultString(spec.GetImage().GetTag(), "480"),
	}
	if spec.GetImage().GetRegistry() != "" {
		imageBlock["registry"] = spec.GetImage().GetRegistry()
	}
	values["image"] = imageBlock

	if len(spec.GetImagePullSecrets()) > 0 {
		pullSecrets := make([]interface{}, 0, len(spec.GetImagePullSecrets()))
		for _, name := range spec.GetImagePullSecrets() {
			pullSecrets = append(pullSecrets, map[string]interface{}{"name": name})
		}
		values["imagePullSecrets"] = pullSecrets
	}

	// ------------------------------- server -------------------------------
	workerReplicas := 2
	if workers := spec.GetWorkers(); workers != nil && workers.Replicas != nil {
		workerReplicas = int(workers.GetReplicas())
	}

	serverConfigBlock := map[string]interface{}{
		"query": map[string]interface{}{
			"maxMemory": defaultString(spec.GetMaxQueryMemory(), "4GB"),
		},
	}
	if locals.AuthEnabled {
		serverConfigBlock["authenticationType"] = "PASSWORD"
	}
	httpsEnabled := spec.GetHttps().GetEnabled()
	if httpsEnabled {
		httpsPort := 8443
		if spec.GetHttps().Port != nil {
			httpsPort = int(spec.GetHttps().GetPort())
		}
		serverConfigBlock["https"] = map[string]interface{}{
			"enabled": true,
			"port":    httpsPort,
			"keystore": map[string]interface{}{
				// The keystore Secret mounts through secretMounts
				// below; the chart wires this path into
				// config.properties.
				"path": "/etc/trino/keystore/" + defaultString(spec.GetHttps().GetKeystoreSecret().GetSecretKey(), "keystore.jks"),
			},
		}
	}

	serverBlock := map[string]interface{}{
		"workers": workerReplicas,
		"node": map[string]interface{}{
			"environment": defaultString(spec.GetNodeEnvironment(), "production"),
		},
		"log": map[string]interface{}{
			"trino": map[string]interface{}{
				"level": defaultString(spec.GetLogLevel(), "INFO"),
			},
		},
		"config": serverConfigBlock,
	}

	// Fault-tolerant execution: the exchange manager spools exchange
	// data durably; the retry policy itself rides
	// additionalConfigProperties (the chart's documented pairing).
	if fte := spec.GetFaultTolerantExecution(); fte != nil {
		baseDirs := make([]interface{}, 0, len(fte.GetExchangeManager().GetBaseDirectories()))
		for _, dir := range fte.GetExchangeManager().GetBaseDirectories() {
			baseDirs = append(baseDirs, dir)
		}
		serverBlock["exchangeManager"] = map[string]interface{}{
			"name":    "filesystem",
			"baseDir": baseDirs,
		}
		if extra := fte.GetExchangeManager().GetAdditionalProperties(); len(extra) > 0 {
			values["additionalExchangeManagerProperties"] = toInterfaceSlice(extra)
		}
	}

	// Worker autoscaling — HPA XOR KEDA (the spec's oneof).
	if hpa := spec.GetWorkers().GetHpa(); hpa != nil {
		autoscalingBlock := map[string]interface{}{
			"enabled":     true,
			"maxReplicas": int(hpa.GetMaxReplicas()),
		}
		// The chart disables a metric when its target is an EMPTY
		// STRING — 0 in the spec maps to that disable contract.
		cpuTarget := 50
		if hpa.TargetCpuUtilizationPercent != nil {
			cpuTarget = int(hpa.GetTargetCpuUtilizationPercent())
		}
		if cpuTarget == 0 {
			autoscalingBlock["targetCPUUtilizationPercentage"] = ""
		} else {
			autoscalingBlock["targetCPUUtilizationPercentage"] = cpuTarget
		}
		memoryTarget := 80
		if hpa.TargetMemoryUtilizationPercent != nil {
			memoryTarget = int(hpa.GetTargetMemoryUtilizationPercent())
		}
		if memoryTarget == 0 {
			autoscalingBlock["targetMemoryUtilizationPercentage"] = ""
		} else {
			autoscalingBlock["targetMemoryUtilizationPercentage"] = memoryTarget
		}
		serverBlock["autoscaling"] = autoscalingBlock
	} else if keda := spec.GetWorkers().GetKeda(); keda != nil {
		var triggers []interface{}
		if err := yaml.Unmarshal([]byte(keda.GetTriggers()), &triggers); err != nil {
			return nil, errors.Wrap(err, "workers.keda.triggers is not valid YAML (expected the `triggers:` array content — a list of KEDA trigger objects)")
		}
		kedaBlock := map[string]interface{}{
			"enabled":         true,
			"maxReplicaCount": int(keda.GetMaxReplicas()),
			"triggers":        triggers,
		}
		if keda.MinReplicas != nil {
			kedaBlock["minReplicaCount"] = int(keda.GetMinReplicas())
		}
		if keda.PollingIntervalSeconds != nil {
			kedaBlock["pollingInterval"] = int(keda.GetPollingIntervalSeconds())
		}
		if keda.CooldownPeriodSeconds != nil {
			kedaBlock["cooldownPeriod"] = int(keda.GetCooldownPeriodSeconds())
		}
		serverBlock["keda"] = kedaBlock
	}
	values["server"] = serverBlock

	// ------------------- additional config properties ---------------------
	// Module-owned lines FIRST (the security spine), then the spec's
	// escape-hatch lines. This list is re-pinned after the helm_values
	// merge — additional_config_properties is the supported extension
	// point, never helm_values.
	configProperties := []interface{}{}
	if locals.AuthEnabled {
		// Trino REQUIRES a shared secret for internal communication
		// once authentication is on (internal-communication.md at the
		// pin) — delivered via env, never rendered.
		configProperties = append(configProperties,
			fmt.Sprintf("internal-communication.shared-secret=${ENV:%s}", vars.SharedSecretEnvVar))
		if !httpsEnabled {
			// Password auth engages ONLY on secure requests. Verified in
			// the server's AuthenticationFilter at the pin:
			// `allow-insecure-over-http` does NOT run password auth over
			// HTTP — it routes plain-HTTP requests to the username-trust
			// authenticator, so the password file would guard nothing.
			// `process-forwarded` is upstream's TLS-terminating-proxy
			// recipe: requests arriving with X-Forwarded-Proto: https
			// (what composed exposure kinds send) are treated secure and
			// the PASSWORD authenticator ENFORCES the file, while plain
			// HTTP data-plane requests fail CLOSED (403). Health probes
			// are unaffected (/v1/info and /v1/status are PUBLIC routes).
			configProperties = append(configProperties,
				"http-server.process-forwarded=true")
		}
	}
	if fte := spec.GetFaultTolerantExecution(); fte != nil {
		configProperties = append(configProperties,
			fmt.Sprintf("retry-policy=%s", fte.GetRetryPolicy()))
	}
	for _, line := range spec.GetAdditionalConfigProperties() {
		configProperties = append(configProperties, line)
	}
	if len(configProperties) > 0 {
		values["additionalConfigProperties"] = configProperties
	}

	// -------------------------------- auth --------------------------------
	if locals.AuthEnabled {
		values["auth"] = buildAuthBlock(locals)
	}
	if httpsEnabled {
		// Mount the keystore Secret onto all nodes at the path wired
		// into config.properties above.
		values["secretMounts"] = []interface{}{
			map[string]interface{}{
				"name":       "trino-keystore",
				"secretName": spec.GetHttps().GetKeystoreSecret().GetSecretName(),
				"path":       "/etc/trino/keystore",
			},
		}
	}

	// ------------------------------ catalogs ------------------------------
	catalogsBlock := map[string]interface{}{}
	for _, catalog := range spec.GetCatalogs().GetPostgres() {
		port := 5432
		if catalog.Port != nil {
			port = int(catalog.GetPort())
		}
		lines := []string{
			"connector.name=postgresql",
			fmt.Sprintf("connection-url=jdbc:postgresql://%s:%d/%s",
				catalog.GetHost().GetValue(), port, catalog.GetDatabase()),
			fmt.Sprintf("connection-user=%s", defaultString(catalog.GetUsername(), "app")),
			fmt.Sprintf("connection-password=${ENV:%s}", catalogPasswordEnvVar(catalog.GetName())),
		}
		lines = append(lines, catalog.GetAdditionalProperties()...)
		catalogsBlock[catalog.GetName()] = strings.Join(lines, "\n")
	}
	for _, catalog := range spec.GetCatalogs().GetMysql() {
		port := 3306
		if catalog.Port != nil {
			port = int(catalog.GetPort())
		}
		// MySQL exposes databases as Trino schemas — the JDBC URL
		// carries no database segment (connector truth at the pin).
		lines := []string{
			"connector.name=mysql",
			fmt.Sprintf("connection-url=jdbc:mysql://%s:%d",
				catalog.GetHost().GetValue(), port),
			fmt.Sprintf("connection-user=%s", defaultString(catalog.GetUsername(), "root")),
			fmt.Sprintf("connection-password=${ENV:%s}", catalogPasswordEnvVar(catalog.GetName())),
		}
		lines = append(lines, catalog.GetAdditionalProperties()...)
		catalogsBlock[catalog.GetName()] = strings.Join(lines, "\n")
	}
	for name, properties := range spec.GetCatalogs().GetCustom() {
		catalogsBlock[name] = properties
	}
	sampleCatalogsEnabled := true
	if catalogs := spec.GetCatalogs(); catalogs != nil && catalogs.SampleCatalogsEnabled != nil {
		sampleCatalogsEnabled = catalogs.GetSampleCatalogsEnabled()
	}
	if !sampleCatalogsEnabled {
		// Helm null-deletes the chart's default map entries.
		catalogsBlock["tpch"] = nil
		catalogsBlock["tpcds"] = nil
	}
	if len(catalogsBlock) > 0 {
		values["catalogs"] = catalogsBlock
	}

	// ----------------- secret-sourced environment variables ---------------
	// The delivery vehicle for every ${ENV:...} reference rendered
	// above: catalog passwords, the internal shared secret, and the
	// user's extra entries — Secret NAMES only ever render here.
	envEntries := []interface{}{}
	if locals.AuthEnabled {
		envEntries = append(envEntries, map[string]interface{}{
			"name": vars.SharedSecretEnvVar,
			"valueFrom": map[string]interface{}{
				"secretKeyRef": map[string]interface{}{
					"name": locals.InternalSecretName,
					"key":  vars.SharedSecretKey,
				},
			},
		})
	}
	for _, ref := range locals.SecretEnvRefs {
		envEntries = append(envEntries, map[string]interface{}{
			"name": ref.EnvName,
			"valueFrom": map[string]interface{}{
				"secretKeyRef": map[string]interface{}{
					"name": ref.SecretName,
					"key":  ref.SecretKey,
				},
			},
		})
	}
	for name, value := range spec.GetExtraEnv() {
		envEntries = append(envEntries, map[string]interface{}{
			"name":  name,
			"value": value,
		})
	}
	if len(envEntries) > 0 {
		values["env"] = envEntries
	}

	// -------------------- coordinator / worker blocks ---------------------
	values["coordinator"] = buildNodeBlock(
		spec.GetCoordinator().GetJvm(),
		defaultString(spec.GetCoordinator().GetMaxQueryMemoryPerNode(), "1GB"),
		spec.GetCoordinator().GetHeapHeadroomPerNode(),
		spec.GetCoordinator().GetResources(),
		spec.GetCoordinator().GetScheduling(),
		map[string]interface{}{
			"config": map[string]interface{}{
				"nodeScheduler": map[string]interface{}{
					"includeCoordinator": spec.GetCoordinator().GetIncludeInScheduling(),
				},
			},
		})

	workerExtras := map[string]interface{}{}
	if shutdown := spec.GetWorkers().GetGracefulShutdown(); shutdown.GetEnabled() {
		gracePeriod := 120
		if shutdown.GracePeriodSeconds != nil {
			gracePeriod = int(shutdown.GetGracePeriodSeconds())
		}
		workerExtras["gracefulShutdown"] = map[string]interface{}{
			"enabled":            true,
			"gracePeriodSeconds": gracePeriod,
		}
		// Chart rule: the pod termination budget must be at least
		// TWICE the drain window — set exactly 2× so the drain always
		// fits.
		workerExtras["terminationGracePeriodSeconds"] = gracePeriod * 2
	}
	values["worker"] = buildNodeBlock(
		spec.GetWorkers().GetJvm(),
		defaultString(spec.GetWorkers().GetMaxQueryMemoryPerNode(), "1GB"),
		spec.GetWorkers().GetHeapHeadroomPerNode(),
		spec.GetWorkers().GetResources(),
		spec.GetWorkers().GetScheduling(),
		workerExtras)

	// --------------------------- config documents -------------------------
	if rules := spec.GetAccessControlRules(); rules != "" {
		values["accessControl"] = map[string]interface{}{
			"type":       "configmap",
			"configFile": "rules.json",
			"rules": map[string]interface{}{
				"rules.json": rules,
			},
		}
	}
	if groups := spec.GetResourceGroupsConfig(); groups != "" {
		values["resourceGroups"] = map[string]interface{}{
			"type":                 "configmap",
			"resourceGroupsConfig": groups,
		}
	}
	if session := spec.GetSessionPropertiesConfig(); session != "" {
		values["sessionProperties"] = map[string]interface{}{
			"type":                    "configmap",
			"sessionPropertiesConfig": session,
		}
	}
	if listeners := spec.GetEventListenerProperties(); len(listeners) > 0 {
		values["eventListenerProperties"] = toInterfaceSlice(listeners)
	}

	// ------------------------------- metrics ------------------------------
	if metrics := spec.GetMetrics(); metrics.GetEnabled() {
		// The standalone JMX exporter FATALS without a hostPort/jmxUrl
		// in its config (verified live: "you must configure 'jmxUrl'
		// or 'hostPort'"), and the chart's default configProperties is
		// EMPTY — enabling the sidecar without composing the config
		// ships a crash-loop. The module renders the chart's own
		// documented pairing; the `tpl` reference keeps the port
		// single-sourced from the chart's jmx.registryPort.
		exporterBlock := map[string]interface{}{
			"enabled": true,
			"configProperties": "hostPort: localhost:{{- .Values.jmx.registryPort }}\n" +
				"startDelaySeconds: 0\n" +
				"ssl: false\n",
		}
		if metrics.GetExporterImage() != "" {
			exporterBlock["image"] = metrics.GetExporterImage()
		}
		values["jmx"] = map[string]interface{}{
			"enabled":  true,
			"exporter": exporterBlock,
		}
		if metrics.GetServiceMonitorEnabled() {
			values["serviceMonitor"] = map[string]interface{}{
				"enabled": true,
			}
		}
	}

	// --------------------------- network policy ---------------------------
	if spec.GetNetworkPolicyEnabled() {
		values["networkPolicy"] = map[string]interface{}{"enabled": true}
	}

	// ------------------------------- service ------------------------------
	serviceBlock := map[string]interface{}{
		"type": defaultString(spec.GetService().GetType(), "ClusterIP"),
	}
	if annotations := spec.GetService().GetAnnotations(); len(annotations) > 0 {
		serviceBlock["annotations"] = toInterfaceMap(annotations)
	}
	values["service"] = serviceBlock

	// -------------------------- helm_values merge -------------------------
	if spec.GetHelmValues() != "" {
		var userValues map[string]interface{}
		if err := yaml.Unmarshal([]byte(spec.GetHelmValues()), &userValues); err != nil {
			return nil, errors.Wrap(err, "spec.helm_values is not valid YAML")
		}
		values = mergeMaps(values, userValues)

		// RE-PINS: the security spine survives any merge. The
		// deterministic names, the authentication wiring and the
		// module-owned config-properties list (which carries the
		// shared secret and the process-forwarded pairing) cannot be
		// silently disabled from the escape hatch.
		values["fullnameOverride"] = locals.ReleaseName
		if locals.AuthEnabled {
			serverValues, ok := values["server"].(map[string]interface{})
			if !ok {
				serverValues = map[string]interface{}{}
			}
			configValues, ok := serverValues["config"].(map[string]interface{})
			if !ok {
				configValues = map[string]interface{}{}
			}
			configValues["authenticationType"] = "PASSWORD"
			serverValues["config"] = configValues
			values["server"] = serverValues
			values["auth"] = buildAuthBlock(locals)
		}
		if len(configProperties) > 0 {
			values["additionalConfigProperties"] = configProperties
		}
	}

	return values, nil
}

// buildAuthBlock renders the chart's auth block — Secret NAMES only.
func buildAuthBlock(locals *Locals) map[string]interface{} {
	authBlock := map[string]interface{}{
		"passwordAuthSecret": locals.PasswordDbSecretName,
	}
	if locals.GroupsSecretName != "" {
		authBlock["groupsAuthSecret"] = locals.GroupsSecretName
	}
	return authBlock
}

// buildNodeBlock renders the shared coordinator/worker shape: JVM heap,
// per-node query memory, resources and scheduling, plus the
// caller-provided extras (deep-merged last).
func buildNodeBlock(
	jvm *kubernetestrinov1alpha1.KubernetesTrinoJvm,
	maxQueryMemoryPerNode string,
	heapHeadroom string,
	resources *kubernetesprovider.ContainerResources,
	scheduling *kubernetestrinov1alpha1.KubernetesTrinoScheduling,
	extras map[string]interface{},
) map[string]interface{} {
	block := map[string]interface{}{}

	// JVM heap: percent-based sizing only works when the fixed -Xmx is
	// UNSET (chart truth). The chart's 8G default is disabled with an
	// EMPTY STRING, never nil: the chart guards the -Xmx line with
	// `{{- if .Values.<node>.jvm.maxHeapSize }}` and "" is falsy
	// there, while a nil does NOT survive to Helm through this
	// engine's release seam (verified live: the chart's -Xmx8G
	// default silently overrode the percent sizing inside the
	// container limit).
	if jvm != nil {
		jvmBlock := map[string]interface{}{}
		if jvm.MaxHeapPercent != nil {
			jvmBlock["maxHeapSize"] = ""
			jvmBlock["maxHeapPercent"] = int(jvm.GetMaxHeapPercent())
		} else if jvm.GetMaxHeapSize() != "" {
			jvmBlock["maxHeapSize"] = jvm.GetMaxHeapSize()
		}
		if len(jvmBlock) > 0 {
			block["jvm"] = jvmBlock
		}
	}

	configBlock := map[string]interface{}{
		"query": map[string]interface{}{
			"maxMemoryPerNode": maxQueryMemoryPerNode,
		},
	}
	if heapHeadroom != "" {
		configBlock["memory"] = map[string]interface{}{
			"heapHeadroomPerNode": heapHeadroom,
		}
	}
	block["config"] = configBlock

	if resourcesValues := resourcesBlock(resources); resourcesValues != nil {
		block["resources"] = resourcesValues
	}
	if scheduling != nil {
		if len(scheduling.GetNodeSelector()) > 0 {
			block["nodeSelector"] = toInterfaceMap(scheduling.GetNodeSelector())
		}
		if tolerations := tolerationsBlock(scheduling.GetTolerations()); tolerations != nil {
			block["tolerations"] = tolerations
		}
	}

	return mergeMaps(block, extras)
}
