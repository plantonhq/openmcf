package module

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
)

// buildHelmValues renders the typed spec into the chart's values map, then
// merges the spec's helm_values escape hatch over it with Helm `-f`
// semantics (maps deep-merge with the later document winning, lists
// replace), then RE-PINS the security-critical values.
//
// SECRET DISCIPLINE (load-bearing): nothing credential-bearing renders
// into values. The chart's own env Secret is OFF (secretEnv.create=false)
// — the module composes `<name>-env` (secrets.go); referenced credentials
// (the database/cache passwords) arrive as extraEnvRaw secretKeyRef
// entries (Secret NAMES only; explicit env beats envFrom — the chart's
// own bring-your-own mechanism); the rendered superset_config.py builds
// every connection FROM ENVIRONMENT. The chart's literal-credential
// paths — database.password, cache.password, init.adminUser.password —
// are never set: the admin user is created by the module's init-command
// override reading ADMIN_PASSWORD from environment, and the
// password-bearing config blocks (RESULTS_BACKEND, the async-queries
// backends) are replaced by module configOverrides snippets that read
// environment too.
//
// PARITY: the Terraform module reaches the same result natively — its
// helm_release passes values = [yamlencode(typed values), helm_values,
// yamlencode(re-pins)] and the provider merges the documents in exactly
// this order. Keep every typed mapping below in lockstep with the
// Terraform module's locals.
func buildHelmValues(locals *Locals) (map[string]interface{}, error) {
	spec := locals.Spec

	// ------------------------- name-budget fail-loud ----------------------
	// Chart truth at the pin: the longest derived child is
	// `<name>-celerybeat` (11-char suffix) — every child must fit the
	// 63-char DNS bound.
	if len(locals.ReleaseName) > vars.NameBudget {
		return nil, errors.Errorf(
			"metadata.name %q is %d characters — the chart derives `<name>-celerybeat` (11-char suffix), so the name must be at most %d characters",
			locals.ReleaseName, len(locals.ReleaseName), vars.NameBudget)
	}

	values := map[string]interface{}{
		// Deterministic child names (`<name>`, `<name>-worker`, the
		// `<name>-env`/`<name>-config` Secrets) — the release name
		// never double-prefixes and the import map stays exact.
		"fullnameOverride": locals.ReleaseName,
		// The module composes the environment Secret — the chart's
		// copy (which renders credentials from values) never ships.
		"secretEnv":     map[string]interface{}{"create": false},
		"envFromSecret": locals.EnvSecretName,
		// The bundled subcharts ride frozen legacy image lines and
		// never ship from this kind — the metadata database and the
		// cache are ALWAYS external (composition-first).
		"postgresql": map[string]interface{}{"enabled": false},
		"redis":      map[string]interface{}{"enabled": false},
	}

	// ------------------------------- image --------------------------------
	values["image"] = map[string]interface{}{
		"repository": defaultString(spec.GetImage().GetRepository(), "apachesuperset.docker.scarf.sh/apache/superset"),
		"tag":        defaultString(spec.GetImage().GetTag(), "6.1.0"),
	}
	if len(spec.GetImagePullSecrets()) > 0 {
		pullSecrets := make([]interface{}, 0, len(spec.GetImagePullSecrets()))
		for _, name := range spec.GetImagePullSecrets() {
			pullSecrets = append(pullSecrets, map[string]interface{}{"name": name})
		}
		values["imagePullSecrets"] = pullSecrets
	}

	// --------------------------- database facts ---------------------------
	// NON-SECRET connection facts only — the password never rides
	// values (DB_PASS arrives via extraEnvRaw). The ssl block feeds the
	// rendered config's sslmode parameters; host/port feed the chart's
	// helpers (init-container waits, results-backend host).
	databaseBlock := map[string]interface{}{
		"host": locals.EnvPlain["DB_HOST"],
		"port": locals.EnvPlain["DB_PORT"],
		"user": locals.EnvPlain["DB_USER"],
		"name": locals.EnvPlain["DB_NAME"],
	}
	if ssl := spec.GetMetadataDatabase().GetSsl(); ssl.GetEnabled() {
		databaseBlock["ssl"] = map[string]interface{}{
			"enabled": true,
			"mode":    defaultString(ssl.GetMode(), "require"),
		}
	}
	values["database"] = databaseBlock

	// ------------------------------- cache --------------------------------
	if locals.CacheEnabled {
		cache := spec.GetCache()
		cacheBlock := map[string]interface{}{
			"enabled": true,
			"host":    locals.EnvPlain["REDIS_HOST"],
			"port":    locals.EnvPlain["REDIS_PORT"],
		}
		if cache.GetUsername() != "" {
			cacheBlock["user"] = cache.GetUsername()
		}
		if cache.CacheDb != nil {
			cacheBlock["cacheDb"] = int(cache.GetCacheDb())
		}
		if cache.CeleryDb != nil {
			cacheBlock["celeryDb"] = int(cache.GetCeleryDb())
		}
		values["cache"] = cacheBlock
	} else {
		// Web-only Superset: synchronous queries, no Celery machinery.
		values["cache"] = map[string]interface{}{"enabled": false}
	}

	// ----------------------------- components -----------------------------
	webBlock := map[string]interface{}{}
	web := spec.GetWeb()
	if hpa := web.GetHpa(); hpa != nil {
		autoscaling := map[string]interface{}{
			"enabled":     true,
			"maxReplicas": int(hpa.GetMaxReplicas()),
		}
		if hpa.MinReplicas != nil {
			autoscaling["minReplicas"] = int(hpa.GetMinReplicas())
		}
		if hpa.TargetCpuUtilizationPercent != nil {
			autoscaling["targetCPUUtilizationPercentage"] = int(hpa.GetTargetCpuUtilizationPercent())
		}
		webBlock["autoscaling"] = autoscaling
	} else {
		replicas := 1
		if web != nil && web.Replicas != nil {
			replicas = int(web.GetReplicas())
		}
		webBlock["replicas"] = map[string]interface{}{
			"enabled":      true,
			"replicaCount": replicas,
		}
	}
	if resources := resourcesBlock(web.GetResources()); resources != nil {
		webBlock["resources"] = resources
	}
	values["supersetNode"] = webBlock

	worker := spec.GetWorker()
	workerBlock := map[string]interface{}{}
	if locals.WorkerEnabled {
		if hpa := worker.GetHpa(); hpa != nil {
			autoscaling := map[string]interface{}{
				"enabled":     true,
				"maxReplicas": int(hpa.GetMaxReplicas()),
			}
			if hpa.MinReplicas != nil {
				autoscaling["minReplicas"] = int(hpa.GetMinReplicas())
			}
			if hpa.TargetCpuUtilizationPercent != nil {
				autoscaling["targetCPUUtilizationPercentage"] = int(hpa.GetTargetCpuUtilizationPercent())
			}
			workerBlock["autoscaling"] = autoscaling
		} else {
			replicas := 1
			if worker != nil && worker.Replicas != nil {
				replicas = int(worker.GetReplicas())
			}
			workerBlock["replicas"] = map[string]interface{}{
				"enabled":      true,
				"replicaCount": replicas,
			}
		}
		if resources := resourcesBlock(worker.GetResources()); resources != nil {
			workerBlock["resources"] = resources
		}
	} else {
		// No cache (or explicitly disabled): the worker Deployment
		// never renders — it would crash-loop without a broker.
		workerBlock["replicas"] = map[string]interface{}{"enabled": false}
	}
	values["supersetWorker"] = workerBlock

	values["supersetCeleryBeat"] = map[string]interface{}{"enabled": locals.BeatEnabled}
	values["supersetCeleryFlower"] = map[string]interface{}{"enabled": locals.FlowerEnabled}

	wsBlock := map[string]interface{}{"enabled": locals.WebsocketsEnabled}
	if locals.WebsocketsEnabled {
		ws := spec.GetWebsockets()
		wsBlock["image"] = map[string]interface{}{
			"repository": defaultString(ws.GetImage().GetRepository(), "oneacrefund/superset-websocket"),
			"tag":        defaultString(ws.GetImage().GetTag(), "latest"),
			"pullPolicy": "IfNotPresent",
		}
		replicas := 1
		if ws.Replicas != nil {
			replicas = int(ws.GetReplicas())
		}
		wsBlock["replicaCount"] = replicas
	}
	values["supersetWebsockets"] = wsBlock

	values["supersetMcp"] = map[string]interface{}{"enabled": locals.McpEnabled}

	// -------------------------------- init --------------------------------
	// createAdmin stays FALSE so the chart's literal-password rendering
	// path (and its config-template fail gate) never engage; the
	// module's command override appends an idempotent
	// create-admin-from-env step after the chart's own rendered init
	// script (schema migration + role init).
	values["init"] = map[string]interface{}{
		"createAdmin":  false,
		"loadExamples": spec.GetInit().GetLoadExamples(),
		"command":      initCommand(),
	}

	if spec.GetBootstrapScript() != "" {
		values["bootstrapScript"] = spec.GetBootstrapScript()
	}

	// --------------------------- configOverrides --------------------------
	configOverrides := moduleConfigOverrides(locals)
	for name, snippet := range spec.GetConfigOverrides() {
		configOverrides[name] = snippet
	}
	if len(configOverrides) > 0 {
		values["configOverrides"] = configOverrides
	}

	// ------------------------------ env wiring ----------------------------
	if len(spec.GetExtraEnv()) > 0 {
		values["extraEnv"] = toInterfaceMap(spec.GetExtraEnv())
	}
	extraEnvRaw := make([]interface{}, 0, len(locals.SecretEnvRefs))
	for _, ref := range locals.SecretEnvRefs {
		extraEnvRaw = append(extraEnvRaw, map[string]interface{}{
			"name": ref.EnvName,
			"valueFrom": map[string]interface{}{
				"secretKeyRef": map[string]interface{}{
					"name": ref.SecretName,
					"key":  ref.SecretKey,
				},
			},
		})
	}
	values["extraEnvRaw"] = extraEnvRaw

	// ---------------------------- feature flags ---------------------------
	if len(spec.GetFeatureFlags()) > 0 {
		flags := map[string]interface{}{}
		for name, enabled := range spec.GetFeatureFlags() {
			flags[name] = enabled
		}
		values["featureFlags"] = flags
	}

	// ------------------------------- service ------------------------------
	serviceBlock := map[string]interface{}{
		"type": defaultString(spec.GetService().GetType(), "ClusterIP"),
		"port": vars.HttpPort,
	}
	if annotations := spec.GetService().GetAnnotations(); len(annotations) > 0 {
		serviceBlock["annotations"] = toInterfaceMap(annotations)
	}
	values["service"] = serviceBlock

	// ------------------------------ scheduling ----------------------------
	if scheduling := spec.GetScheduling(); scheduling != nil {
		if len(scheduling.GetNodeSelector()) > 0 {
			values["nodeSelector"] = toInterfaceMap(scheduling.GetNodeSelector())
		}
		if tolerations := tolerationsBlock(scheduling.GetTolerations()); tolerations != nil {
			values["tolerations"] = tolerations
		}
	}

	// -------------------------- helm_values merge -------------------------
	if spec.GetHelmValues() != "" {
		var userValues map[string]interface{}
		if err := yaml.Unmarshal([]byte(spec.GetHelmValues()), &userValues); err != nil {
			return nil, errors.Wrap(err, "spec.helm_values is not valid YAML")
		}
		values = mergeMaps(values, userValues)

		// RE-PINS: the security spine survives any merge — the
		// deterministic names, the module-owned env Secret contract,
		// the dead bundled subcharts, and the env-driven admin
		// bootstrap cannot be silently re-enabled or redirected from
		// the escape hatch.
		values["fullnameOverride"] = locals.ReleaseName
		values["secretEnv"] = map[string]interface{}{"create": false}
		values["envFromSecret"] = locals.EnvSecretName
		values["postgresql"] = map[string]interface{}{"enabled": false}
		values["redis"] = map[string]interface{}{"enabled": false}
		initValues, ok := values["init"].(map[string]interface{})
		if !ok {
			initValues = map[string]interface{}{}
		}
		initValues["createAdmin"] = false
		initValues["command"] = initCommand()
		values["init"] = initValues
		mergedOverrides, ok := values["configOverrides"].(map[string]interface{})
		if !ok {
			mergedOverrides = map[string]interface{}{}
		}
		for name, snippet := range moduleConfigOverrides(locals) {
			mergedOverrides[name] = snippet
		}
		values["configOverrides"] = mergedOverrides
	}

	return values, nil
}

// initCommand is the module's init-Job command: the chart's own rendered
// init script (schema migration + role init — createAdmin=false keeps it
// admin-free) followed by an idempotent create-admin step that reads the
// admin identity and password FROM ENVIRONMENT (the env Secret /
// extraEnvRaw) — the chart's literal-password rendering path is never
// used. Keep byte-identical with the Terraform module's locals.
func initCommand() []interface{} {
	script := strings.Join([]string{
		fmt.Sprintf(". %s/superset_bootstrap.sh", vars.ConfigMountPath),
		fmt.Sprintf(". %s/superset_init.sh", vars.ConfigMountPath),
		`if superset fab list-users 2>/dev/null | grep -qF "username:${ADMIN_USERNAME}"; then echo "Admin user already exists, skipping."; else superset fab create-admin --username "$ADMIN_USERNAME" --firstname Superset --lastname Admin --email "$ADMIN_EMAIL" --password "$ADMIN_PASSWORD"; fi`,
	}, "; ")
	return []interface{}{"/bin/sh", "-c", script}
}

// moduleConfigOverrides are the module-owned superset_config.py snippets —
// the env-indirection replacements for the chart's password-from-values
// config blocks. Keep byte-identical with the Terraform module's locals.
func moduleConfigOverrides(locals *Locals) map[string]interface{} {
	overrides := map[string]interface{}{}

	// The chart renders cache.password LITERALLY into RESULTS_BACKEND
	// and the async-queries backends — this kind never sets it, so on
	// authed stores those blocks must be redefined reading environment.
	if locals.CacheEnabled && locals.CachePasswordSecret != "" {
		overrides["planton_redis_auth"] = strings.Join([]string{
			"# Authed cache: the chart-rendered RESULTS_BACKEND and async-queries",
			"# backends carry no password (it never rides values) — redefine them",
			"# reading the environment the pods already carry.",
			"from flask_caching.backends.rediscache import RedisCache as _PlantonRedisCache",
			"RESULTS_BACKEND = _PlantonRedisCache(",
			"    host=env('REDIS_HOST'),",
			"    port=int(env('REDIS_PORT', '6379')),",
			"    password=env('REDIS_PASSWORD') or None,",
			"    key_prefix='superset_results',",
			")",
			"GLOBAL_ASYNC_QUERIES_CACHE_BACKEND = {",
			"    'CACHE_TYPE': 'RedisCache',",
			"    'CACHE_REDIS_HOST': env('REDIS_HOST'),",
			"    'CACHE_REDIS_PORT': int(env('REDIS_PORT', '6379')),",
			"    'CACHE_REDIS_PASSWORD': env('REDIS_PASSWORD', ''),",
			"    'CACHE_REDIS_DB': int(env('REDIS_DB', '1')),",
			"    'CACHE_KEY_PREFIX': 'qc-',",
			"    'CACHE_DEFAULT_TIMEOUT': 86400,",
			"}",
			"GLOBAL_ASYNC_QUERIES_RESULTS_BACKEND = {",
			"    'backend': 'redis',",
			"    'host': env('REDIS_HOST'),",
			"    'port': int(env('REDIS_PORT', '6379')),",
			"    'password': env('REDIS_PASSWORD', ''),",
			"    'db': int(env('REDIS_DB', '1')),",
			"    'prefix': 'qc-',",
			"}",
		}, "\n")
	}

	// The websocket JWT: the ws server reads JWT_SECRET from its
	// environment natively; Superset's side reads the same variable —
	// the chart's values-borne jwtSecret path never carries material.
	if locals.WebsocketsEnabled {
		overrides["planton_ws_jwt"] = strings.Join([]string{
			"# The async-queries JWT shared with the websocket server — both sides",
			"# read the same module-generated environment variable.",
			fmt.Sprintf("GLOBAL_ASYNC_QUERIES_JWT_SECRET = env('%s')", vars.JwtSecretEnvVar),
		}, "\n")
	}

	return overrides
}
