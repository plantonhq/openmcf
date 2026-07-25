package module

import (
	"github.com/pkg/errors"
	kubernetessignozv1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetessignoz/v1"
	"sigs.k8s.io/yaml"
)

// buildHelmValues renders the typed spec into the chart-values map, merges
// the helm_values escape hatch over it with Helm -f semantics, and
// re-pins the two fullnames LAST — the exact semantic twin of the
// Terraform module's values = [typed, helm_values, re-pin] document list.
//
// SECRET DISCIPLINE (load-bearing): the bundled ClickHouse admin password
// is deliberately ABSENT here — Resources injects it as a Pulumi secret
// Output after this merge (the set_sensitive twin), so it never rides a
// plain values map and never appears in a preview diff. KNOW THIS about
// the upstream grain: the chart renders that password into the
// ClickHouseInstallation object's user section and into literal container
// env — anyone with read access on the namespace's CHI/pods can see it.
// What the module guarantees is that the chart's publicly-documented
// default password NEVER ships and the credential is unique per install.
// The external arm is fully secret-native (existingSecret → secretKeyRef),
// and the SMTP password rides a valueFrom secretKeyRef env entry.
func buildHelmValues(locals *Locals) (map[string]interface{}, error) {
	spec := locals.Spec

	values := map[string]interface{}{
		"fullnameOverride": locals.ReleaseName,
	}

	// ---- globals ---------------------------------------------------------
	// Every chart image here is the SPLIT registry+repository form and its
	// registry key defers to global.imageRegistry — one override reaches
	// the SigNoz server, the collector, ClickHouse, the bundled operator
	// and its metrics exporter, and ZooKeeper. The bundled arm's UDF init
	// container image (alpine) follows its own chart key — helm_values
	// territory, documented on the spec field.
	global := map[string]interface{}{}
	if spec.ImageRegistry != "" {
		global["imageRegistry"] = spec.ImageRegistry
	}
	if spec.ClusterName != "" {
		global["clusterName"] = spec.ClusterName
		values["clusterName"] = spec.ClusterName
	}
	if len(spec.ImagePullSecrets) > 0 {
		pullSecrets := make([]interface{}, 0, len(spec.ImagePullSecrets))
		for _, s := range spec.ImagePullSecrets {
			pullSecrets = append(pullSecrets, s)
		}
		global["imagePullSecrets"] = pullSecrets
	}
	if len(global) > 0 {
		values["global"] = global
	}

	// ---- database arm ------------------------------------------------------
	if locals.IsExternal {
		values["clickhouse"] = map[string]interface{}{"enabled": false}
		values["externalClickhouse"] = externalClickHouseValues(spec.GetExternalClickhouse())
	} else {
		values["clickhouse"] = managedClickHouseValues(locals, spec.GetManagedClickhouse())
	}

	// ---- scheduling (server + collector + migrator) -------------------------
	var nodeSelector map[string]interface{}
	var tolerations []interface{}
	priorityClass := ""
	if scheduling := spec.GetScheduling(); scheduling != nil {
		if len(scheduling.NodeSelector) > 0 {
			nodeSelector = stringMapToInterface(scheduling.NodeSelector)
		}
		if len(scheduling.Tolerations) > 0 {
			tolerations = tolerationsSlice(scheduling.Tolerations)
		}
		priorityClass = scheduling.PriorityClassName
	}

	// ---- the signoz server ---------------------------------------------------
	values["signoz"] = serverValues(locals, nodeSelector, tolerations, priorityClass)

	// ---- the ingestion collector ----------------------------------------------
	values["otelCollector"] = collectorValues(spec.GetOtelCollector(), nodeSelector, tolerations, priorityClass)

	// ---- the schema migrator ----------------------------------------------------
	migrator := map[string]interface{}{}
	if nodeSelector != nil {
		migrator["nodeSelector"] = nodeSelector
	}
	if tolerations != nil {
		migrator["tolerations"] = tolerations
	}
	if len(migrator) > 0 {
		values["telemetryStoreMigrator"] = migrator
	}

	// ---- escape hatch + re-pin -----------------------------------------------
	if spec.HelmValues != "" {
		overrides := map[string]interface{}{}
		if err := yaml.Unmarshal([]byte(spec.HelmValues), &overrides); err != nil {
			return nil, errors.Wrap(err, "helm_values is not a valid YAML document")
		}
		values = mergeMaps(values, overrides)
	}

	// Re-pin the fullnames LAST — the one deliberate exception to the
	// escape hatch's last-word contract. Every child name — and the
	// exported outputs built from them — derives from these fullnames;
	// letting an override move them would break every output.
	values["fullnameOverride"] = locals.ReleaseName
	if !locals.IsExternal {
		clickhouse, ok := values["clickhouse"].(map[string]interface{})
		if !ok {
			clickhouse = map[string]interface{}{}
			values["clickhouse"] = clickhouse
		}
		clickhouse["fullnameOverride"] = locals.ClickHouseFullname
	}

	return values, nil
}

// managedClickHouseValues renders the bundled arm (capacity and topology;
// SigNoz owns everything inside it). The admin password is injected by
// Resources as a secret Output — never here.
func managedClickHouseValues(locals *Locals, managed *kubernetessignozv1.KubernetesSignozManagedClickHouse) map[string]interface{} {
	shards := int32(1)
	replicas := int32(1)
	diskSize := vars.DefaultClickHouseDiskSize
	zookeeperReplicas := int32(1)

	clickhouse := map[string]interface{}{
		"enabled":          true,
		"fullnameOverride": locals.ClickHouseFullname,
	}

	persistence := map[string]interface{}{
		"enabled": true,
	}

	if managed != nil {
		if managed.Shards != nil {
			shards = managed.GetShards()
		}
		if managed.Replicas != nil {
			replicas = managed.GetReplicas()
		}
		if managed.GetDiskSize() != "" {
			diskSize = managed.GetDiskSize()
		}
		if managed.StorageClass.GetValue() != "" {
			persistence["storageClass"] = managed.StorageClass.GetValue()
		}
		if resources := resourcesMap(managed.GetResources()); resources != nil {
			clickhouse["resources"] = resources
		}
		if len(managed.AllowedNetworkIps) > 0 {
			ips := make([]interface{}, 0, len(managed.AllowedNetworkIps))
			for _, ip := range managed.AllowedNetworkIps {
				ips = append(ips, ip)
			}
			clickhouse["allowedNetworkIps"] = ips
		}
		if zk := managed.GetZookeeper(); zk != nil {
			if zk.Replicas != nil {
				zookeeperReplicas = zk.GetReplicas()
			}
		}
		if cold := coldStorageValues(managed.GetColdStorage()); cold != nil {
			clickhouse["coldStorage"] = cold
		}
	}

	persistence["size"] = diskSize
	clickhouse["persistence"] = persistence
	clickhouse["layout"] = map[string]interface{}{
		"shardsCount":   shards,
		"replicasCount": replicas,
	}

	zookeeper := map[string]interface{}{
		"replicaCount": zookeeperReplicas,
	}
	if managed != nil {
		if zk := managed.GetZookeeper(); zk != nil {
			if resources := resourcesMap(zk.GetResources()); resources != nil {
				zookeeper["resources"] = resources
			}
		}
	}
	clickhouse["zookeeper"] = zookeeper

	return clickhouse
}

// coldStorageValues renders the cold-storage arm: keyless IRSA (role
// annotations on the ClickHouse service account) XOR declared keys
// (rendered by the chart into ClickHouse's storage configuration — the
// upstream grain, taught on the spec field).
func coldStorageValues(cold *kubernetessignozv1.KubernetesSignozColdStorage) map[string]interface{} {
	if cold == nil {
		return nil
	}
	if s3 := cold.GetS3(); s3 != nil {
		out := map[string]interface{}{
			"enabled":  true,
			"type":     "s3",
			"endpoint": s3.Endpoint,
		}
		if s3.IrsaRoleArn != "" {
			out["role"] = map[string]interface{}{
				"enabled": true,
				"annotations": map[string]interface{}{
					"eks.amazonaws.com/role-arn": s3.IrsaRoleArn,
				},
			}
		} else {
			out["accessKey"] = s3.AccessKey
			out["secretAccess"] = s3.SecretKey
		}
		return out
	}
	if gcs := cold.GetGcs(); gcs != nil {
		return map[string]interface{}{
			"enabled":      true,
			"type":         "gcs",
			"endpoint":     gcs.Endpoint,
			"accessKey":    gcs.AccessKey,
			"secretAccess": gcs.SecretKey,
		}
	}
	return nil
}

// externalClickHouseValues renders the bring-your-own arm — fully
// secret-native: the chart wires the referenced Secret as a secretKeyRef.
func externalClickHouseValues(external *kubernetessignozv1.KubernetesSignozExternalClickHouse) map[string]interface{} {
	clusterName := external.ClusterName.GetValue()
	if clusterName == "" {
		clusterName = "cluster"
	}
	tcpPort := int32(9000)
	if external.TcpPort != nil {
		tcpPort = external.GetTcpPort()
	}
	httpPort := int32(8123)
	if external.HttpPort != nil {
		httpPort = external.GetHttpPort()
	}

	out := map[string]interface{}{
		"host":                      external.Host.GetValue(),
		"cluster":                   clusterName,
		"tcpPort":                   tcpPort,
		"httpPort":                  httpPort,
		"user":                      external.Username,
		"existingSecret":            external.PasswordSecret.SecretName.GetValue(),
		"existingSecretPasswordKey": external.PasswordSecret.SecretKey,
	}
	if external.Secure {
		out["secure"] = true
	}
	if external.Verify {
		out["verify"] = true
	}
	return out
}

// serverValues renders the SigNoz server tier: state volume, resources,
// and the env map carrying the typed alerting/SMTP configuration. Typed
// entries WIN over the user's advanced env map (the spec's documented
// contract). Env keys follow SigNoz's own derivation
// (signoz_<section>_<key>, embedded underscores doubled).
func serverValues(locals *Locals, nodeSelector map[string]interface{}, tolerations []interface{}, priorityClass string) map[string]interface{} {
	spec := locals.Spec
	server := spec.GetServer()

	diskSize := vars.DefaultServerDiskSize
	persistence := map[string]interface{}{
		"enabled": true,
	}

	out := map[string]interface{}{}

	env := map[string]interface{}{}
	if server != nil {
		if server.GetDiskSize() != "" {
			diskSize = server.GetDiskSize()
		}
		if server.StorageClass.GetValue() != "" {
			persistence["storageClass"] = server.StorageClass.GetValue()
		}
		if resources := resourcesMap(server.GetResources()); resources != nil {
			out["resources"] = resources
		}
		for k, v := range server.Env {
			env[k] = v
		}
		if server.ExternalUrl != "" {
			env["signoz_alertmanager_signoz_external__url"] = server.ExternalUrl
		}
		if smtp := server.GetSmtp(); smtp != nil {
			env["signoz_emailing_enabled"] = "true"
			env["signoz_emailing_smtp_address"] = smtp.Address
			env["signoz_emailing_smtp_from"] = smtp.From
			if smtp.Username != "" {
				env["signoz_emailing_smtp_auth_username"] = smtp.Username
			}
			if smtp.TlsEnabled {
				env["signoz_emailing_smtp_tls_enabled"] = "true"
			}
			if passwordSecret := smtp.GetPasswordSecret(); passwordSecret != nil {
				// The chart's flexible env structure renders this as a
				// proper env-from-secret — a reference, never material.
				env["signoz_emailing_smtp_auth_password"] = map[string]interface{}{
					"valueFrom": map[string]interface{}{
						"secretKeyRef": map[string]interface{}{
							"name": passwordSecret.Name,
							"key":  passwordSecret.Key,
						},
					},
				}
			}
		}
	}

	persistence["size"] = diskSize
	out["persistence"] = persistence
	if len(env) > 0 {
		out["env"] = env
	}
	if nodeSelector != nil {
		out["nodeSelector"] = nodeSelector
	}
	if tolerations != nil {
		out["tolerations"] = tolerations
	}
	if priorityClass != "" {
		out["priorityClassName"] = priorityClass
	}
	return out
}

// collectorValues renders the ingestion collector tier. The pipeline
// receiver LISTS are always rendered from the toggles (lists replace
// under Helm merge — rendering them from one derivation is what keeps the
// Service ports and the collector pipelines in agreement by
// construction).
func collectorValues(collector *kubernetessignozv1.KubernetesSignozOtelCollector, nodeSelector map[string]interface{}, tolerations []interface{}, priorityClass string) map[string]interface{} {
	replicas := int32(1)
	jaegerEnabled := true
	zipkinEnabled := false
	httpLogsEnabled := true

	out := map[string]interface{}{}

	if collector != nil {
		if collector.Replicas != nil {
			replicas = collector.GetReplicas()
		}
		if collector.JaegerReceiverEnabled != nil {
			jaegerEnabled = collector.GetJaegerReceiverEnabled()
		}
		zipkinEnabled = collector.ZipkinReceiverEnabled
		if collector.HttpLogsReceiversEnabled != nil {
			httpLogsEnabled = collector.GetHttpLogsReceiversEnabled()
		}
		if resources := resourcesMap(collector.GetResources()); resources != nil {
			out["resources"] = resources
		}
		if collector.LowCardinalityExceptionGrouping {
			out["lowCardinalityExceptionGrouping"] = true
		}
		if autoscaling := collector.GetAutoscaling(); autoscaling != nil && autoscaling.Enabled {
			minReplicas := int32(1)
			if autoscaling.MinReplicas != nil {
				minReplicas = autoscaling.GetMinReplicas()
			}
			maxReplicas := int32(11)
			if autoscaling.MaxReplicas != nil {
				maxReplicas = autoscaling.GetMaxReplicas()
			}
			hpa := map[string]interface{}{
				"enabled":     true,
				"minReplicas": minReplicas,
				"maxReplicas": maxReplicas,
			}
			if autoscaling.TargetCpuUtilizationPercent != nil {
				hpa["targetCPUUtilizationPercentage"] = autoscaling.GetTargetCpuUtilizationPercent()
			}
			if autoscaling.TargetMemoryUtilizationPercent != nil {
				hpa["targetMemoryUtilizationPercentage"] = autoscaling.GetTargetMemoryUtilizationPercent()
			}
			out["autoscaling"] = hpa
		}
	}

	out["replicaCount"] = replicas

	// Receiver toggles: jaeger + http-logs default ON (the chart's
	// grain), zipkin defaults OFF. Only diverging Service-port toggles
	// render; the pipeline lists always render.
	ports := map[string]interface{}{}
	if zipkinEnabled {
		ports["zipkin"] = map[string]interface{}{"enabled": true}
	}
	if !jaegerEnabled {
		ports["jaeger-thrift"] = map[string]interface{}{"enabled": false}
		ports["jaeger-grpc"] = map[string]interface{}{"enabled": false}
	}
	if !httpLogsEnabled {
		ports["logsheroku"] = map[string]interface{}{"enabled": false}
		ports["logsjson"] = map[string]interface{}{"enabled": false}
	}
	if len(ports) > 0 {
		out["ports"] = ports
	}

	tracesReceivers := []interface{}{"otlp"}
	if jaegerEnabled {
		tracesReceivers = append(tracesReceivers, "jaeger")
	}
	if zipkinEnabled {
		tracesReceivers = append(tracesReceivers, "zipkin")
	}
	logsReceivers := []interface{}{"otlp"}
	if httpLogsEnabled {
		logsReceivers = append(logsReceivers, "httplogreceiver/heroku", "httplogreceiver/json")
	}

	config := map[string]interface{}{
		"service": map[string]interface{}{
			"pipelines": map[string]interface{}{
				"traces": map[string]interface{}{"receivers": tracesReceivers},
				"logs":   map[string]interface{}{"receivers": logsReceivers},
			},
		},
	}
	if zipkinEnabled {
		config["receivers"] = map[string]interface{}{
			"zipkin": map[string]interface{}{"endpoint": "0.0.0.0:9411"},
		}
	}
	out["config"] = config

	if nodeSelector != nil {
		out["nodeSelector"] = nodeSelector
	}
	if tolerations != nil {
		out["tolerations"] = tolerations
	}
	if priorityClass != "" {
		out["priorityClassName"] = priorityClass
	}
	return out
}
