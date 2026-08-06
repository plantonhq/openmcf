package module

import (
	"strconv"

	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
)

// buildHelmValues renders the typed spec into the chart's values map, then
// merges the spec's helm_values escape hatch over it with Helm `-f`
// semantics (maps deep-merge with the later document winning, lists
// replace).
//
// SECRET DISCIPLINE (load-bearing): nothing rendered here carries
// credential material — only Secret NAMES. The database, broker,
// result-backend and log-read connections ride the chart's *SecretName
// contracts pointing at module-composed Secrets; setting those names is
// ALSO what stops the chart from rendering its own connection Secrets
// (whose split-values path would embed the password in rendered values)
// and its render-time random credentials (regenerated on every upgrade).
//
// PARITY: the Terraform module reaches the same result natively — its
// helm_release passes values = [yamlencode(typed values), helm_values,
// yamlencode(re-pins)] and the provider merges the documents in exactly
// this order. Keep every typed mapping below in lockstep with the
// Terraform module's locals.
func buildHelmValues(locals *Locals) (map[string]interface{}, error) {
	spec := locals.Spec
	values := map[string]interface{}{}

	// The bundled PostgreSQL subchart NEVER ships: upstream marks it
	// non-production and its Bitnami image line is frozen
	// (bitnamilegacy). Re-pinned after the escape-hatch merge below.
	values["postgresql"] = map[string]interface{}{"enabled": false}

	values["executor"] = locals.Executor

	// airflowVersion gates version-specific rendering;
	// defaultAirflowTag keeps the image tag in lockstep with it (the
	// chart's airflow_image helper falls back to defaultAirflowTag, NOT
	// to airflowVersion — setting only airflowVersion would deploy a
	// mismatched image).
	values["airflowVersion"] = locals.AirflowVersion
	values["defaultAirflowTag"] = locals.AirflowVersion

	// ---- connection secrets (module-composed; names only) ----------------
	dataBlock := map[string]interface{}{
		"metadataSecretName": locals.MetadataConnSecretName,
	}
	if locals.CeleryEnabled {
		// The result backend must carry the `db+` scheme, so it never
		// falls back to the metadata Secret; the broker URL Secret is
		// module-composed on the bundled/valkey arms and user-provided
		// on the existing-secret arm.
		dataBlock["resultBackendSecretName"] = locals.ResultBackendConnSecretName
		dataBlock["brokerUrlSecretName"] = locals.BrokerUrlSecretName
	}
	values["data"] = dataBlock

	// ---- security key secrets (module-owned or BYO; names only) ----------
	values["fernetKeySecretName"] = locals.FernetKeySecretName
	values["apiSecretKeySecretName"] = locals.ApiSecretKeySecretName
	values["webserverSecretKeySecretName"] = locals.WebserverSecretKeyName
	values["jwtSecretName"] = locals.JwtSecretName

	// ---- broker ----------------------------------------------------------
	if locals.BundledRedisEnabled {
		redisBlock := map[string]interface{}{
			"enabled": true,
			// The module owns the password Secret — the chart's own
			// path draws a NEW random password on every render behind
			// a pre-install hook (upstream's admitted hack).
			"passwordSecretName": locals.RedisPasswordSecretName,
		}
		bundled := spec.GetBroker().GetBundledRedis()
		persistence := map[string]interface{}{}
		if bundled.GetPersistenceSize() != "" {
			persistence["size"] = bundled.GetPersistenceSize()
		}
		if bundled.GetStorageClass() != "" {
			persistence["storageClassName"] = bundled.GetStorageClass()
		}
		if len(persistence) > 0 {
			redisBlock["persistence"] = persistence
		}
		if resources := resourcesBlock(bundled.GetResources()); resources != nil {
			redisBlock["resources"] = resources
		}
		values["redis"] = redisBlock
	} else {
		// External/no broker: the chart must never deploy its bundled
		// Redis (it would on Celery executors by default).
		values["redis"] = map[string]interface{}{"enabled": false}
	}

	// ---- components -------------------------------------------------------
	components := spec.GetComponents()

	apiServer := map[string]interface{}{}
	if c := components.GetApiServer(); c != nil {
		if c.Replicas != nil {
			apiServer["replicas"] = int(c.GetReplicas())
		}
		if resources := resourcesBlock(c.GetResources()); resources != nil {
			apiServer["resources"] = resources
		}
	}
	if len(apiServer) > 0 {
		values["apiServer"] = apiServer
	}

	scheduler := map[string]interface{}{}
	if c := components.GetScheduler(); c != nil {
		if c.Replicas != nil {
			scheduler["replicas"] = int(c.GetReplicas())
		}
		if resources := resourcesBlock(c.GetResources()); resources != nil {
			scheduler["resources"] = resources
		}
	}
	if len(scheduler) > 0 {
		values["scheduler"] = scheduler
	}

	dagProcessor := map[string]interface{}{}
	if c := components.GetDagProcessor(); c != nil {
		if c.Replicas != nil {
			dagProcessor["replicas"] = int(c.GetReplicas())
		}
		if resources := resourcesBlock(c.GetResources()); resources != nil {
			dagProcessor["resources"] = resources
		}
	}
	if len(dagProcessor) > 0 {
		values["dagProcessor"] = dagProcessor
	}

	triggerer := map[string]interface{}{}
	if c := components.GetTriggerer(); c != nil {
		if c.Enabled != nil {
			triggerer["enabled"] = c.GetEnabled()
		}
		if c.Replicas != nil {
			triggerer["replicas"] = int(c.GetReplicas())
		}
		if c.GetPersistenceSize() != "" {
			triggerer["persistence"] = map[string]interface{}{"size": c.GetPersistenceSize()}
		}
		if resources := resourcesBlock(c.GetResources()); resources != nil {
			triggerer["resources"] = resources
		}
	}
	if len(triggerer) > 0 {
		values["triggerer"] = triggerer
	}

	// Celery workers — the chart's FLAT worker keys (authoritative at
	// this pin; the workers.celery sub-grain is the chart-2.0 migration
	// surface).
	workers := map[string]interface{}{}
	if c := components.GetWorkers(); c != nil {
		if c.Replicas != nil {
			workers["replicas"] = int(c.GetReplicas())
		}
		if resources := resourcesBlock(c.GetResources()); resources != nil {
			workers["resources"] = resources
		}
		persistence := map[string]interface{}{}
		if c.PersistenceEnabled != nil {
			persistence["enabled"] = c.GetPersistenceEnabled()
		}
		if c.GetPersistenceSize() != "" {
			persistence["size"] = c.GetPersistenceSize()
		}
		if len(persistence) > 0 {
			workers["persistence"] = persistence
		}
		if keda := c.GetKeda(); keda != nil && keda.GetEnabled() {
			kedaBlock := map[string]interface{}{"enabled": true}
			if keda.MinReplicas != nil {
				kedaBlock["minReplicaCount"] = int(keda.GetMinReplicas())
			}
			if keda.MaxReplicas != nil {
				kedaBlock["maxReplicaCount"] = int(keda.GetMaxReplicas())
			}
			if keda.PollingIntervalSeconds != nil {
				kedaBlock["pollingInterval"] = int(keda.GetPollingIntervalSeconds())
			}
			if keda.CooldownPeriodSeconds != nil {
				kedaBlock["cooldownPeriod"] = int(keda.GetCooldownPeriodSeconds())
			}
			workers["keda"] = kedaBlock
		}
	}
	if len(workers) > 0 {
		values["workers"] = workers
	}

	// ---- DAG delivery ------------------------------------------------------
	if gitSync := spec.GetDags().GetGitSync(); gitSync != nil {
		gitSyncBlock := map[string]interface{}{
			"enabled": true,
			"repo":    gitSync.GetRepo(),
		}
		// The chart renders BOTH env generations UNCONDITIONALLY —
		// GITSYNC_REF from `ref` (v4) and GIT_SYNC_BRANCH from
		// `branch` (legacy) — and git-sync v4 translates the
		// deprecated --branch OVER --ref, so a ref-only rendering
		// silently syncs the chart's default branch (verified live:
		// the sidecar fetched v2-2-stable while ref carried the
		// declared value). Both keys always render the spec value —
		// INCLUDING the empty string, which neutralizes the chart's
		// v2-2-stable defaults so the spec's Empty = HEAD promise
		// holds (git-sync treats empty ref/branch as HEAD).
		gitSyncBlock["ref"] = gitSync.GetRef()
		gitSyncBlock["branch"] = gitSync.GetRef()
		gitSyncBlock["subPath"] = gitSync.GetSubPath()
		if gitSync.PeriodSeconds != nil {
			gitSyncBlock["period"] = formatSeconds(int(gitSync.GetPeriodSeconds()))
		}
		if gitSync.Depth != nil {
			gitSyncBlock["depth"] = int(gitSync.GetDepth())
		}
		if gitSync.GetCredentialsSecret() != "" {
			gitSyncBlock["credentialsSecret"] = gitSync.GetCredentialsSecret()
		}
		if gitSync.GetSshKeySecret() != "" {
			gitSyncBlock["sshKeySecret"] = gitSync.GetSshKeySecret()
		}
		if gitSync.GetKnownHosts() != "" {
			gitSyncBlock["knownHosts"] = gitSync.GetKnownHosts()
		}
		if resources := resourcesBlock(gitSync.GetResources()); resources != nil {
			gitSyncBlock["resources"] = resources
		}
		values["dags"] = map[string]interface{}{"gitSync": gitSyncBlock}
	} else if persistence := spec.GetDags().GetPersistence(); persistence != nil {
		persistenceBlock := map[string]interface{}{"enabled": true}
		if persistence.GetSize() != "" {
			persistenceBlock["size"] = persistence.GetSize()
		}
		if persistence.GetStorageClass() != "" {
			persistenceBlock["storageClassName"] = persistence.GetStorageClass()
		}
		if persistence.GetExistingClaim() != "" {
			persistenceBlock["existingClaim"] = persistence.GetExistingClaim()
		}
		values["dags"] = map[string]interface{}{"persistence": persistenceBlock}
	}

	// ---- logging -----------------------------------------------------------
	if logs := spec.GetLogging().GetPersistence(); logs != nil && logs.GetEnabled() {
		logsBlock := map[string]interface{}{"enabled": true}
		if logs.GetSize() != "" {
			logsBlock["size"] = logs.GetSize()
		}
		if logs.GetStorageClass() != "" {
			logsBlock["storageClassName"] = logs.GetStorageClass()
		}
		values["logs"] = map[string]interface{}{"persistence": logsBlock}
	}
	if locals.LogBackend != "" {
		values[locals.LogBackend] = map[string]interface{}{
			"enabled":    true,
			"secretName": locals.LogReadConnSecretName,
		}
	}

	// ---- pgbouncer -----------------------------------------------------------
	if locals.PgbouncerEnabled {
		pgb := spec.GetPgbouncer()
		pgbouncerBlock := map[string]interface{}{
			"enabled": true,
			// The module composes pgbouncer.ini + users.txt — the
			// chart's own rendering path embeds the database password
			// in Helm values and is never used.
			"configSecretName": locals.PgbouncerConfigSecretName,
			// The module-composed stats DSN for the metrics-exporter
			// sidecar (see secrets.go) — with statsSecretName set, the
			// chart skips creating its split-values stats Secret.
			"metricsExporterSidecar": map[string]interface{}{
				"statsSecretName": locals.PgbouncerStatsSecretName,
			},
		}
		if pgb.MetadataPoolSize != nil {
			pgbouncerBlock["metadataPoolSize"] = int(pgb.GetMetadataPoolSize())
		}
		if pgb.ResultBackendPoolSize != nil {
			pgbouncerBlock["resultBackendPoolSize"] = int(pgb.GetResultBackendPoolSize())
		}
		if pgb.MaxClientConnections != nil {
			pgbouncerBlock["maxClientConn"] = int(pgb.GetMaxClientConnections())
		}
		if resources := resourcesBlock(pgb.GetResources()); resources != nil {
			pgbouncerBlock["resources"] = resources
		}
		values["pgbouncer"] = pgbouncerBlock
	}

	// ---- metrics + examples ---------------------------------------------------
	statsdEnabled := true
	if spec.StatsdEnabled != nil {
		statsdEnabled = spec.GetStatsdEnabled()
	}
	values["statsd"] = map[string]interface{}{"enabled": statsdEnabled}

	if spec.GetLoadExamples() {
		// The env-var form, NOT config.core: the official image BAKES
		// AIRFLOW__CORE__LOAD_EXAMPLES=False as a container env, and
		// Airflow's precedence puts env above airflow.cfg — a cfg-only
		// True is silently defeated (verified live: examples never
		// parsed; the chart's own docs prescribe the env route).
		values["env"] = []interface{}{
			map[string]interface{}{"name": "AIRFLOW__CORE__LOAD_EXAMPLES", "value": "True"},
		}
	}

	// ---- admin bootstrap user ---------------------------------------------
	values["createUserJob"] = createUserJobBlock(locals)

	// The chart's migration Job defaults to a post-install Helm HOOK,
	// and post-install hooks only run AFTER the release wait completes
	// — while every component's wait-for-airflow-migrations init
	// container blocks on the migrations that hook would apply. Under
	// any wait-style install that is a deadlock by construction
	// (verified live: no Job existed while every init container
	// crash-looped on "unapplied migrations"). Hook-less mode makes the
	// Job an ordinary release resource applied WITH the install; the
	// chart's own ttlSecondsAfterFinished: 300 default self-deletes the
	// finished Job, so day-2 applies recreate it cleanly.
	values["migrateDatabaseJob"] = map[string]interface{}{"useHelmHooks": false}

	// ---- scheduling ---------------------------------------------------------
	if scheduling := spec.GetScheduling(); scheduling != nil {
		if len(scheduling.GetNodeSelector()) > 0 {
			values["nodeSelector"] = toInterfaceMap(scheduling.GetNodeSelector())
		}
		if len(scheduling.GetTolerations()) > 0 {
			values["tolerations"] = tolerationsBlock(scheduling.GetTolerations())
		}
	}

	// ---- images ----------------------------------------------------------------
	if images := spec.GetImages(); images != nil {
		imagesBlock := map[string]interface{}{}
		airflowImage := map[string]interface{}{}
		if images.GetAirflowRepository() != "" {
			airflowImage["repository"] = images.GetAirflowRepository()
			// The mirrored image also serves the KubernetesExecutor pod
			// template (the chart wires worker_container_repository from
			// images.airflow).
			values["defaultAirflowRepository"] = images.GetAirflowRepository()
		}
		if images.GetAirflowTag() != "" {
			airflowImage["tag"] = images.GetAirflowTag()
			values["defaultAirflowTag"] = images.GetAirflowTag()
		}
		if images.GetAirflowDigest() != "" {
			airflowImage["digest"] = images.GetAirflowDigest()
		}
		if len(airflowImage) > 0 {
			imagesBlock["airflow"] = airflowImage
		}
		if images.GetStatsdRepository() != "" {
			imagesBlock["statsd"] = map[string]interface{}{"repository": images.GetStatsdRepository()}
		}
		if images.GetRedisRepository() != "" {
			imagesBlock["redis"] = map[string]interface{}{"repository": images.GetRedisRepository()}
		}
		if images.GetPgbouncerRepository() != "" {
			imagesBlock["pgbouncer"] = map[string]interface{}{"repository": images.GetPgbouncerRepository()}
		}
		if images.GetGitSyncRepository() != "" {
			imagesBlock["gitSync"] = map[string]interface{}{"repository": images.GetGitSyncRepository()}
		}
		if len(imagesBlock) > 0 {
			values["images"] = imagesBlock
		}
	}

	// ---- escape hatch (merged LAST, helm -f semantics) ------------------------
	if spec.GetHelmValues() != "" {
		overrides := map[string]interface{}{}
		if err := yaml.Unmarshal([]byte(spec.GetHelmValues()), &overrides); err != nil {
			return nil, errors.Wrap(err, "failed to parse helm_values as a YAML document")
		}
		values = mergeMaps(values, overrides)
	}

	// Deliberate post-merge re-pins — the two exceptions to the escape
	// hatch's last-word contract (twin of the Terraform module's third
	// values document):
	//   - useStandardNaming stays false: every child name (and every
	//     exported output built from them) derives from the release
	//     name; letting an override move the naming scheme would break
	//     every output.
	//   - postgresql.enabled stays false: the bundled database is
	//     non-production by upstream's own definition and its image
	//     line is frozen.
	values["useStandardNaming"] = false
	if postgresqlBlock, ok := values["postgresql"].(map[string]interface{}); ok {
		postgresqlBlock["enabled"] = false
	} else {
		values["postgresql"] = map[string]interface{}{"enabled": false}
	}

	return values, nil
}

// createUserJobBlock renders the admin bootstrap Job wiring. The chart's
// default renders the admin password as a LITERAL POD ARGUMENT — this
// block replaces the password argument with a job-scoped env var read
// from the admin Secret, so the credential never appears in a rendered
// pod spec.
func createUserJobBlock(locals *Locals) map[string]interface{} {
	// useHelmHooks false for the same reason as migrateDatabaseJob: a
	// post-install hook only fires after the release wait, which never
	// completes on a fresh database (the migration-hook deadlock class).
	if !locals.AdminCreate {
		return map[string]interface{}{"enabled": false, "useHelmHooks": false}
	}
	return map[string]interface{}{
		"enabled":      true,
		"useHelmHooks": false,
		"env": []interface{}{
			map[string]interface{}{
				"name": "ADMIN_PASSWORD",
				"valueFrom": map[string]interface{}{
					"secretKeyRef": map[string]interface{}{
						"name": locals.AdminSecretName,
						"key":  locals.AdminSecretKey,
					},
				},
			},
		},
		// The args mirror the chart's own default shape (bash -c +
		// `airflow users create "$@"` + flag list) with ONE change: the
		// password flag reads the env var above instead of a rendered
		// literal. firstName/lastName keep the chart's defaults.
		"args": []interface{}{
			"bash",
			"-c",
			"exec \\\nairflow users create \"$@\"",
			"--",
			"-r", "Admin",
			"-u", locals.AdminUsername,
			"-e", locals.AdminEmail,
			"-f", "admin",
			"-l", "user",
			"-p", "$(ADMIN_PASSWORD)",
		},
	}
}

// formatSeconds renders a go-style duration string for the git-sync
// period ("5s" — git-sync v4 takes durations, not bare integers).
func formatSeconds(seconds int) string {
	return strconv.Itoa(seconds) + "s"
}
