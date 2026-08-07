package module

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	kuberneteslocustv1alpha1 "github.com/plantonhq/planton/catalog/kubernetes/kuberneteslocust/v1alpha1"
	"sigs.k8s.io/yaml"
)

// buildHelmValues renders the typed spec into the chart's values map,
// then merges the spec's helm_values escape hatch over it with Helm `-f`
// semantics (maps deep-merge with the later document winning, lists
// replace), then RE-PINS the security-critical values.
//
// SECRET DISCIPLINE (load-bearing): nothing rendered here carries
// credential material. The web-UI login credential and the Flask
// session key ride Secret-projected files (mount_external_secret); test
// credentials ride Secret references (environment_external_secret /
// environment_load_from_secrets). The chart's `environment_secret`
// block — which renders a Secret FROM VALUES — is force-emptied after
// the escape-hatch merge so it can never engage.
//
// PARITY: the Terraform module reaches the same result natively — its
// helm_release passes values = [yamlencode(typed values), helm_values,
// yamlencode(re-pins)] and the provider merges the documents in exactly
// this order (its re-pin document nulls environment_secret, the Helm
// null-deletion twin of the force-empty here). Keep every typed mapping
// below in lockstep with the Terraform module's locals.
func buildHelmValues(locals *Locals) (map[string]interface{}, error) {
	spec := locals.Spec
	loadTest := spec.GetLoadTest()

	// ------------------------- name-budget fail-loud ----------------------
	// The longest derived child name is the module's own
	// `<name>-locustfile` ConfigMap (11-char suffix); every chart child
	// (`-master`, `-worker`, `-web-auth`, `-auth`) is shorter. All must
	// fit the 63-char DNS bound.
	if len(locals.ReleaseName) > vars.NameBudget {
		return nil, errors.Errorf(
			"metadata.name %q is %d characters — the module derives `<name>-locustfile` (11-char suffix), so the name must be at most %d characters",
			locals.ReleaseName, len(locals.ReleaseName), vars.NameBudget)
	}

	// -------------------------- login tag floor ---------------------------
	imageTag := defaultString(spec.GetImage().GetTag(), vars.DefaultImageTag)
	if locals.WebLoginEnabled && !imageTagAllowsWebLogin(imageTag) {
		return nil, errors.Errorf(
			"image.tag %q cannot prove Locust >= %d.%d — below that the chart renders the login credential as a literal pod argument, which this module refuses; use a numeric tag at or above %d.%d.0, or disable web_ui_auth",
			imageTag, vars.AuthMinMajor, vars.AuthMinMinor, vars.AuthMinMajor, vars.AuthMinMinor)
	}

	values := map[string]interface{}{
		// Deterministic child names (`<name>-master`, `<name>-worker`,
		// the bare `<name>` Service) — the release name never
		// double-prefixes and the import map stays exact (the fullname
		// re-pin discipline).
		"fullnameOverride": locals.ReleaseName,
		// The COMBINED image form; the tag always pinned explicitly so
		// the deployed Locust version is declared, never inherited.
		"image": map[string]interface{}{
			"repository": defaultString(spec.GetImage().GetRepository(), vars.DefaultImageRepository),
			"tag":        imageTag,
		},
	}

	if len(spec.GetImagePullSecrets()) > 0 {
		pullSecrets := make([]interface{}, 0, len(spec.GetImagePullSecrets()))
		for _, name := range spec.GetImagePullSecrets() {
			pullSecrets = append(pullSecrets, map[string]interface{}{"name": name})
		}
		values["imagePullSecrets"] = pullSecrets
	}

	values["loadtest"] = loadtestBlock(locals, loadTest)
	values["master"] = masterBlock(locals, spec.GetMaster())
	values["worker"] = workerBlock(locals, spec.GetWorkers())

	// ------------------------------ service -------------------------------
	// Ports pinned to the chart defaults so the rendered Service (and
	// the exported endpoints) never drift with a chart-default change.
	serviceBlock := map[string]interface{}{
		"type":       "ClusterIP",
		"port":       vars.WebPort,
		"targetPort": vars.WebPort,
	}
	if service := spec.GetService(); service != nil {
		if service.GetType() != "" {
			serviceBlock["type"] = service.GetType()
		}
		if len(service.GetAnnotations()) > 0 {
			serviceBlock["annotations"] = toInterfaceMap(service.GetAnnotations())
		}
	}
	values["service"] = serviceBlock

	// ------------------- extraConfigMaps (login backend) ------------------
	if locals.WebLoginEnabled {
		values["extraConfigMaps"] = map[string]interface{}{
			locals.WebAuthCodeName: vars.WebAuthCodeMountPath,
		}
	}

	// -------------------------- helm_values merge --------------------------
	if spec.GetHelmValues() != "" {
		var overrides map[string]interface{}
		if err := yaml.Unmarshal([]byte(spec.GetHelmValues()), &overrides); err != nil {
			return nil, errors.Wrap(err, "failed to parse helm_values as YAML")
		}
		values = mergeMaps(values, overrides)
	}

	// ------------------------------ re-pins --------------------------------
	// The deliberate exceptions to the escape hatch's last-word
	// contract (twin of the Terraform module's third values document):
	// the deterministic names, the script wiring, the login wiring —
	// and the chart's values-rendered-Secret path, force-emptied so
	// credentials can never ride rendered values.
	values["fullnameOverride"] = locals.ReleaseName
	loadtestValues, ok := values["loadtest"].(map[string]interface{})
	if !ok {
		return nil, errors.New("helm_values replaced the loadtest block with a non-map value")
	}
	loadtestValues["locust_locustfile_configmap"] = locals.LocustfileConfigMap
	loadtestValues["locust_locustfile"] = locals.LocustfileName
	loadtestValues["locust_locustfile_path"] = vars.LocustfileMountPath
	loadtestValues["locust_lib_configmap"] = locals.LibConfigMap
	loadtestValues["environment_secret"] = map[string]interface{}{}
	masterValues, ok := values["master"].(map[string]interface{})
	if !ok {
		return nil, errors.New("helm_values replaced the master block with a non-map value")
	}
	if locals.WebLoginEnabled {
		masterValues["auth"] = map[string]interface{}{"enabled": true}
		masterValues["args"] = masterArgs(locals)
		loadtestValues["mount_external_secret"] = webAuthSecretMount(locals)
	} else {
		masterValues["auth"] = map[string]interface{}{"enabled": false}
	}

	return values, nil
}

// loadtestBlock renders the load-test surface. The script ConfigMaps are
// named EXPLICITLY (empty string when absent — the chart's bundled
// example defaults are a fragile literal-string coupling, never
// engaged).
func loadtestBlock(locals *Locals, loadTest *kuberneteslocustv1alpha1.KubernetesLocustLoadTest) map[string]interface{} {
	block := map[string]interface{}{
		"name":                           locals.LoadTestName,
		"locust_locustfile":              locals.LocustfileName,
		"locust_locustfile_path":         vars.LocustfileMountPath,
		"locust_locustfile_configmap":    locals.LocustfileConfigMap,
		"locust_lib_configmap":           locals.LibConfigMap,
		"pip_requirementsfile_configmap": loadTest.GetPipRequirementsConfigMap(),
		"headless":                       loadTest.GetHeadless(),
	}

	// The target host — resolved literal (a reference resolves before
	// the module runs). Empty stays empty: the locustfile must then
	// declare its own host.
	block["locust_host"] = loadTest.GetTargetHost().GetValue()

	if len(loadTest.GetPipPackages()) > 0 {
		block["pip_packages"] = toInterfaceSlice(loadTest.GetPipPackages())
	}
	if len(loadTest.GetEnvironment()) > 0 {
		block["environment"] = toInterfaceMap(loadTest.GetEnvironment())
	}
	if len(loadTest.GetEnvFromSecrets()) > 0 {
		block["environment_load_from_secrets"] = toInterfaceSlice(loadTest.GetEnvFromSecrets())
	}
	if len(loadTest.GetEnvFromSecretKeys()) > 0 {
		externalSecrets := map[string]interface{}{}
		for _, entry := range loadTest.GetEnvFromSecretKeys() {
			externalSecrets[entry.GetSecretName()] = toInterfaceSlice(entry.GetKeys())
		}
		block["environment_external_secret"] = externalSecrets
	}
	// The chart takes tags as ONE space-joined string (rendered into
	// `--tags`/`--exclude-tags`); the CEL forbids whitespace inside a
	// tag, so the join is unambiguous.
	if len(loadTest.GetTags()) > 0 {
		block["tags"] = strings.Join(loadTest.GetTags(), " ")
	}
	if len(loadTest.GetExcludeTags()) > 0 {
		block["excludeTags"] = strings.Join(loadTest.GetExcludeTags(), " ")
	}
	if locals.WebLoginEnabled {
		block["mount_external_secret"] = webAuthSecretMount(locals)
	}

	return block
}

// masterBlock renders the master surface. The login wiring: the chart's
// auth.enabled renders `--web-login` (the tag floor guarantees the
// modern path — the legacy credentials-as-arguments path never
// renders), and the `-f` argument names the login backend next to the
// locustfile — command-line -f overrides the LOCUST_LOCUSTFILE env, so
// the master loads BOTH files while workers (default env) load the
// locustfile alone.
func masterBlock(locals *Locals, master *kuberneteslocustv1alpha1.KubernetesLocustMaster) map[string]interface{} {
	block := map[string]interface{}{
		"logLevel": defaultString(master.GetLogLevel(), "INFO"),
		// The module's content hash — pod-template annotation; see
		// configChecksum for why the chart's own checksums are not
		// enough.
		"annotations": map[string]interface{}{
			vars.ChecksumAnnotation: locals.ConfigChecksum,
		},
		"auth": map[string]interface{}{
			"enabled": locals.WebLoginEnabled,
		},
	}
	if resources := resourcesBlock(master.GetResources()); resources != nil {
		block["resources"] = resources
	}
	if master.GetPdbEnabled() {
		block["pdb"] = map[string]interface{}{"enabled": true}
	}
	if locals.WebLoginEnabled {
		block["args"] = masterArgs(locals)
	}
	if scheduling := master.GetScheduling(); scheduling != nil {
		if len(scheduling.GetNodeSelector()) > 0 {
			block["nodeSelector"] = toInterfaceMap(scheduling.GetNodeSelector())
		}
		if tolerations := tolerationsBlock(scheduling.GetTolerations()); tolerations != nil {
			block["tolerations"] = tolerations
		}
	}
	return block
}

// workerBlock renders the worker surface, including the autoscaling
// arms. CHART CONTRACT: the KEDA ScaledObject reuses worker.hpa
// minReplicas/maxReplicas for its bounds while worker.hpa.enabled stays
// false, and the worker Deployment still renders `replicas` on the KEDA
// arm (the template gates it on hpa.enabled only) — the module pins
// replicas to the KEDA floor so a Helm upgrade resets scaling to the
// floor, not to an unrelated count.
func workerBlock(locals *Locals, workers *kuberneteslocustv1alpha1.KubernetesLocustWorkers) map[string]interface{} {
	block := map[string]interface{}{
		"logLevel": defaultString(workers.GetLogLevel(), "INFO"),
		"annotations": map[string]interface{}{
			vars.ChecksumAnnotation: locals.ConfigChecksum,
		},
	}
	if resources := resourcesBlock(workers.GetResources()); resources != nil {
		block["resources"] = resources
	}
	if workers.GetPdbEnabled() {
		block["pdb"] = map[string]interface{}{"enabled": true}
	}
	if scheduling := workers.GetScheduling(); scheduling != nil {
		if len(scheduling.GetNodeSelector()) > 0 {
			block["nodeSelector"] = toInterfaceMap(scheduling.GetNodeSelector())
		}
		if tolerations := tolerationsBlock(scheduling.GetTolerations()); tolerations != nil {
			block["tolerations"] = tolerations
		}
	}

	switch {
	case workers.GetHpa() != nil:
		hpa := workers.GetHpa()
		minReplicas := 1
		if hpa.MinReplicas != nil {
			minReplicas = int(hpa.GetMinReplicas())
		}
		targetCpu := 40
		if hpa.TargetCpuUtilizationPercent != nil {
			targetCpu = int(hpa.GetTargetCpuUtilizationPercent())
		}
		block["hpa"] = map[string]interface{}{
			"enabled":                        true,
			"minReplicas":                    minReplicas,
			"maxReplicas":                    int(hpa.GetMaxReplicas()),
			"targetCPUUtilizationPercentage": targetCpu,
		}
	case workers.GetKeda() != nil:
		keda := workers.GetKeda()
		minReplicas := 1
		if keda.MinReplicas != nil {
			minReplicas = int(keda.GetMinReplicas())
		}
		block["replicas"] = minReplicas
		block["hpa"] = map[string]interface{}{
			"enabled":     false,
			"minReplicas": minReplicas,
			"maxReplicas": int(keda.GetMaxReplicas()),
		}
		pollingInterval := 15
		if keda.PollingIntervalSeconds != nil {
			pollingInterval = int(keda.GetPollingIntervalSeconds())
		}
		cooldownPeriod := 30
		if keda.CooldownPeriodSeconds != nil {
			cooldownPeriod = int(keda.GetCooldownPeriodSeconds())
		}
		block["keda"] = map[string]interface{}{
			"enabled":         true,
			"pollingInterval": pollingInterval,
			"cooldownPeriod":  cooldownPeriod,
			"triggers":        kedaTriggers(locals, keda),
		}
	default:
		// A raw field access (workers.Replicas) is NOT nil-safe the
		// way getters are — the workers block itself may be absent.
		replicas := 1
		if workers != nil && workers.Replicas != nil {
			replicas = int(workers.GetReplicas())
		}
		block["replicas"] = replicas
	}

	return block
}

// kedaTriggers renders the trigger list the ScaledObject scales on. The
// default trigger reads the LIVE USER COUNT from the master's own stats
// API — rendered explicitly and pinned (never the chart's tpl'd
// default), with the spec-level CEL guaranteeing that API is reachable
// (login off, non-headless). custom_triggers replaces it wholesale.
func kedaTriggers(locals *Locals, keda *kuberneteslocustv1alpha1.KubernetesLocustWorkerKeda) string {
	if keda.GetCustomTriggers() != "" {
		return keda.GetCustomTriggers()
	}
	targetUsers := 50
	if keda.TargetUsersPerWorker != nil {
		targetUsers = int(keda.GetTargetUsersPerWorker())
	}
	return fmt.Sprintf(`- type: metrics-api
  metadata:
    activationTargetValue: "0"
    targetValue: "%d"
    url: "%s/stats/requests"
    format: json
    valueLocation: user_count
`, targetUsers, locals.WebEndpoint)
}

// masterArgs appends the module's `-f` to the master command line:
// the locustfile plus the login backend, comma-separated absolute
// paths (Locust's own multi-locustfile form).
func masterArgs(locals *Locals) []interface{} {
	return []interface{}{
		"-f",
		vars.LocustfileMountPath + "/" + locals.LocustfileName +
			"," + vars.WebAuthCodeMountPath + "/planton_auth.py",
	}
}

// webAuthSecretMount projects the `<name>-auth` Secret's three keys as
// files under the login backend's read path.
func webAuthSecretMount(locals *Locals) map[string]interface{} {
	return map[string]interface{}{
		"mountPath": vars.WebAuthSecretMountPath,
		"files": map[string]interface{}{
			locals.AuthSecretName: []interface{}{
				vars.AuthUsernameKey,
				vars.AuthPasswordKey,
				vars.AuthFlaskSecretKeyKey,
			},
		},
	}
}
