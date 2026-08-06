package module

import (
	kubernetesprovider "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes"
	kubernetesmysqlv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetesmysql/v1alpha1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apiextensions"
	kubernetesmeta "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// createCluster renders the pxc.percona.com/v1 PerconaXtraDBCluster as an
// UNTYPED CustomResource. The typed crd2pulumi tree is structurally unable
// to carry this CR: the CRD models spec.backup.storages as a map of
// storage-name → nested object (s3/azure/volume blocks), but crd2pulumi
// flattens it to map[string]map[string]string — a shape that cannot hold
// nested objects, so any backup-enabled cluster is unrepresentable through
// the typed path.
//
// The spec body built here is the exact twin of the Terraform module's
// local.mysql_manifest (locals.tf) — same keys rendered and omitted,
// numbers as ints, booleans as booleans. An unset optional is simply never
// inserted into the map (the Go twin of TF's null-prune), so the apiserver
// applies the CRD's own defaults. Shape errors still fail loudly without
// compile-time typing: the operator validates the applied spec against its
// schema, and the kind-cluster E2E lanes exercise the rendered arms live.
func createCluster(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependencies []pulumi.ResourceOption,
) (pulumi.Resource, error) {
	spec := locals.Spec

	specBody := map[string]interface{}{
		"crVersion":      vars.CrVersion,
		"secretsName":    locals.RootPasswordSecretName,
		"updateStrategy": resolveUpdateStrategy(spec.GetUpdateStrategy()),
		"upgradeOptions": map[string]interface{}{
			"versionServiceEndpoint": vars.VersionServiceEndpoint,
			"apply":                  vars.UpgradeApply,
		},
		"pxc":          pxcBody(locals),
		"haproxy":      haproxyBody(locals),
		"proxysql":     proxysqlBody(locals),
		"logcollector": logcollectorBody(spec),
	}

	if spec.GetPause() {
		specBody["pause"] = true
	}
	if tls := tlsBody(spec); tls != nil {
		specBody["tls"] = tls
	}
	if users := usersBody(locals); len(users) > 0 {
		specBody["users"] = users
	}
	if backup := backupBody(locals); backup != nil {
		specBody["backup"] = backup
	}
	if flags := unsafeFlagsBody(spec.GetUnsafe()); len(flags) > 0 {
		specBody["unsafeFlags"] = flags
	}

	return apiextensions.NewCustomResource(ctx, locals.ClusterName,
		&apiextensions.CustomResourceArgs{
			ApiVersion: pulumi.String("pxc.percona.com/v1"),
			Kind:       pulumi.String("PerconaXtraDBCluster"),
			Metadata: &kubernetesmeta.ObjectMetaArgs{
				Name:      pulumi.String(locals.ClusterName),
				Namespace: pulumi.String(locals.Namespace),
				Labels:    pulumi.ToStringMap(locals.Labels),
				Finalizers: pulumi.StringArray{
					pulumi.String(vars.PxcPodsFinalizer),
				},
			},
			OtherFields: kubernetes.UntypedArgs{
				"spec": specBody,
			},
		}, append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, dependencies...)...)
}

func resolveUpdateStrategy(strategy string) string {
	if strategy != "" {
		return strategy
	}
	return vars.DefaultUpdateStrategy
}

func resolveInstances(spec *kubernetesmysqlv1alpha1.KubernetesMysqlSpec) int {
	if spec.Instances != nil {
		return int(spec.GetInstances())
	}
	return vars.DefaultInstances
}

func resolvePxcImage(spec *kubernetesmysqlv1alpha1.KubernetesMysqlSpec) string {
	if spec.GetImageName() != "" {
		return spec.GetImageName()
	}
	return vars.PxcDefaultImage
}

// serviceType mirrors TF's coalesce(type, "ClusterIP"): every rendered
// expose block always carries an explicit Service type.
func serviceType(exposeType string) string {
	if exposeType != "" {
		return exposeType
	}
	return "ClusterIP"
}

// pxcBody is the twin of TF's pxc_body.
func pxcBody(locals *Locals) map[string]interface{} {
	spec := locals.Spec

	body := map[string]interface{}{
		"size":  resolveInstances(spec),
		"image": resolvePxcImage(spec),
		"volumeSpec": map[string]interface{}{
			"persistentVolumeClaim": pvcBody(
				spec.GetStorage().GetStorageClass().GetValue(),
				spec.GetStorage().GetSize()),
		},
	}

	if spec.GetMysqlConfig() != "" {
		body["configuration"] = spec.GetMysqlConfig()
	}

	// auto_recovery defaults true upstream; only an explicit false renders.
	if spec.AutoRecovery != nil && !spec.GetAutoRecovery() {
		body["autoRecovery"] = false
	}

	if resources := containerResourcesBody(spec.GetResources()); resources != nil {
		body["resources"] = resources
	}

	scheduling := spec.GetScheduling()
	if scheduling.GetAntiAffinityTopologyKey() != "" {
		body["affinity"] = map[string]interface{}{
			"antiAffinityTopologyKey": scheduling.GetAntiAffinityTopologyKey(),
		}
	}
	if len(scheduling.GetNodeSelector()) > 0 {
		body["nodeSelector"] = scheduling.GetNodeSelector()
	}
	if tolerations := tolerationsBody(scheduling.GetTolerations()); len(tolerations) > 0 {
		body["tolerations"] = tolerations
	}
	if scheduling.GetPriorityClassName() != "" {
		body["priorityClassName"] = scheduling.GetPriorityClassName()
	}

	pdb := spec.GetPodDisruptionBudget()
	if pdb.GetMaxUnavailable() > 0 || pdb.GetMinAvailable() > 0 {
		pdbBody := map[string]interface{}{}
		if pdb.GetMaxUnavailable() > 0 {
			pdbBody["maxUnavailable"] = int(pdb.GetMaxUnavailable())
		}
		if pdb.GetMinAvailable() > 0 {
			pdbBody["minAvailable"] = int(pdb.GetMinAvailable())
		}
		body["podDisruptionBudget"] = pdbBody
	}

	if pullSecrets := imagePullSecretsBody(spec.GetImagePullSecrets()); pullSecrets != nil {
		body["imagePullSecrets"] = pullSecrets
	}

	return body
}

// pvcBody renders a persistentVolumeClaim block for the pxc and proxysql
// volumeSpec (the backup filesystem storage carries an extra accessModes
// key and is rendered inline in backupStoragesBody).
func pvcBody(storageClass, size string) map[string]interface{} {
	body := map[string]interface{}{
		"resources": map[string]interface{}{
			"requests": map[string]interface{}{
				"storage": size,
			},
		},
	}
	if storageClass != "" {
		body["storageClassName"] = storageClass
	}
	return body
}

// containerResourcesBody is the twin of the TF resources sub-object: the
// limits/requests keys render whenever their block is present in the spec —
// even when they prune to an empty object — and the whole body renders
// whenever the resources block is present.
func containerResourcesBody(resources *kubernetesprovider.ContainerResources) map[string]interface{} {
	if resources == nil {
		return nil
	}
	body := map[string]interface{}{}
	if limits := resources.GetLimits(); limits != nil {
		body["limits"] = cpuMemoryBody(limits.GetCpu(), limits.GetMemory())
	}
	if requests := resources.GetRequests(); requests != nil {
		body["requests"] = cpuMemoryBody(requests.GetCpu(), requests.GetMemory())
	}
	return body
}

func cpuMemoryBody(cpu, memory string) map[string]interface{} {
	body := map[string]interface{}{}
	if cpu != "" {
		body["cpu"] = cpu
	}
	if memory != "" {
		body["memory"] = memory
	}
	return body
}

func tolerationsBody(tolerations []*kubernetesprovider.WorkloadToleration) []interface{} {
	var result []interface{}
	for _, toleration := range tolerations {
		body := map[string]interface{}{}
		if toleration.GetKey() != "" {
			body["key"] = toleration.GetKey()
		}
		if toleration.GetOperator() != "" {
			body["operator"] = toleration.GetOperator()
		}
		if toleration.GetValue() != "" {
			body["value"] = toleration.GetValue()
		}
		if toleration.GetEffect() != "" {
			body["effect"] = toleration.GetEffect()
		}
		if toleration.TolerationSeconds != nil {
			body["tolerationSeconds"] = int(toleration.GetTolerationSeconds())
		}
		result = append(result, body)
	}
	return result
}

func imagePullSecretsBody(names []string) []interface{} {
	if len(names) == 0 {
		return nil
	}
	secrets := make([]interface{}, 0, len(names))
	for _, name := range names {
		secrets = append(secrets, map[string]interface{}{"name": name})
	}
	return secrets
}

// haproxyBody is the twin of TF's haproxy_body: the disabled arm renders
// as just {enabled: false} when ProxySQL owns the proxy role.
func haproxyBody(locals *Locals) map[string]interface{} {
	if !locals.IsHaproxy {
		return map[string]interface{}{"enabled": false}
	}

	spec := locals.Spec
	haproxy := spec.GetProxy().GetHaproxy()

	size := vars.DefaultProxyReplicas
	if haproxy != nil && haproxy.Replicas != nil {
		size = int(haproxy.GetReplicas())
	}

	body := map[string]interface{}{
		"enabled": true,
		"size":    size,
		"image":   vars.HaproxyImage,
	}

	if haproxy.GetConfig() != "" {
		body["configuration"] = haproxy.GetConfig()
	}
	if resources := containerResourcesBody(haproxy.GetResources()); resources != nil {
		body["resources"] = resources
	}
	if exposePrimary := haproxy.GetExposePrimary(); exposePrimary != nil {
		body["exposePrimary"] = proxyServiceBody(exposePrimary)
	}
	if exposeReplicas := haproxy.GetExposeReplicas(); exposeReplicas != nil {
		body["exposeReplicas"] = haproxyExposeReplicasBody(exposeReplicas)
	}
	if pullSecrets := imagePullSecretsBody(spec.GetImagePullSecrets()); pullSecrets != nil {
		body["imagePullSecrets"] = pullSecrets
	}

	return body
}

func proxyServiceBody(expose *kubernetesmysqlv1alpha1.KubernetesMysqlProxyService) map[string]interface{} {
	body := map[string]interface{}{
		"type": serviceType(expose.GetType()),
	}
	if len(expose.GetAnnotations()) > 0 {
		body["annotations"] = expose.GetAnnotations()
	}
	return body
}

func haproxyExposeReplicasBody(expose *kubernetesmysqlv1alpha1.KubernetesMysqlHaproxyReplicasService) map[string]interface{} {
	// enabled defaults true — a rendered exposeReplicas block always
	// carries it explicitly, exactly like the TF twin.
	enabled := true
	if expose.Enabled != nil {
		enabled = expose.GetEnabled()
	}
	body := map[string]interface{}{
		"enabled": enabled,
		"type":    serviceType(expose.GetType()),
	}
	if expose.GetOnlyReaders() {
		body["onlyReaders"] = true
	}
	if len(expose.GetAnnotations()) > 0 {
		body["annotations"] = expose.GetAnnotations()
	}
	return body
}

// proxysqlBody is the twin of TF's proxysql_body: the disabled arm renders
// as just {enabled: false} when HAProxy owns the proxy role.
func proxysqlBody(locals *Locals) map[string]interface{} {
	if locals.IsHaproxy {
		return map[string]interface{}{"enabled": false}
	}

	spec := locals.Spec
	proxysql := spec.GetProxy().GetProxysql()

	size := vars.DefaultProxyReplicas
	if proxysql.Replicas != nil {
		size = int(proxysql.GetReplicas())
	}

	body := map[string]interface{}{
		"enabled": true,
		"size":    size,
		"image":   vars.ProxysqlImage,
		"volumeSpec": map[string]interface{}{
			"persistentVolumeClaim": pvcBody(
				proxysql.GetStorage().GetStorageClass().GetValue(),
				proxysql.GetStorage().GetSize()),
		},
	}

	if proxysql.GetConfig() != "" {
		body["configuration"] = proxysql.GetConfig()
	}
	if resources := containerResourcesBody(proxysql.GetResources()); resources != nil {
		body["resources"] = resources
	}
	if exposePrimary := proxysql.GetExposePrimary(); exposePrimary != nil {
		body["expose"] = proxyServiceBody(exposePrimary)
	}
	if pullSecrets := imagePullSecretsBody(spec.GetImagePullSecrets()); pullSecrets != nil {
		body["imagePullSecrets"] = pullSecrets
	}

	return body
}

// tlsBody is the twin of TF's tls_body: omitted entirely when the spec
// block is absent (operator default = enabled with self-generated
// certificates); a rendered block always carries enabled explicitly.
func tlsBody(spec *kubernetesmysqlv1alpha1.KubernetesMysqlSpec) map[string]interface{} {
	tls := spec.GetTls()
	if tls == nil {
		return nil
	}

	enabled := true
	if tls.Enabled != nil {
		enabled = tls.GetEnabled()
	}
	body := map[string]interface{}{"enabled": enabled}

	if len(tls.GetSans()) > 0 {
		body["SANs"] = tls.GetSans()
	}
	if tls.GetIssuer().GetValue() != "" {
		kind := tls.GetIssuerKind()
		if kind == "" {
			kind = vars.DefaultIssuerKind
		}
		body["issuerConf"] = map[string]interface{}{
			"name":  tls.GetIssuer().GetValue(),
			"kind":  kind,
			"group": vars.CertManagerIssuerGroup,
		}
	}

	return body
}

// usersBody is the twin of TF's users local: declared passwords travel
// only as passwordSecretRef pointers at the deterministic Secret names.
func usersBody(locals *Locals) []interface{} {
	var users []interface{}
	for _, user := range locals.Spec.GetUsers() {
		body := map[string]interface{}{"name": user.GetName()}
		if len(user.GetDbs()) > 0 {
			body["dbs"] = user.GetDbs()
		}
		if len(user.GetHosts()) > 0 {
			body["hosts"] = user.GetHosts()
		}
		if len(user.GetGrants()) > 0 {
			body["grants"] = user.GetGrants()
		}
		if user.GetWithGrantOption() {
			body["withGrantOption"] = true
		}
		if user.GetPassword() != "" {
			body["passwordSecretRef"] = map[string]interface{}{
				"name": locals.ClusterName + "-user-" + user.GetName(),
				"key":  "password",
			}
		}
		users = append(users, body)
	}
	return users
}

// backupBody is the twin of TF's backup_body: nil when the spec has no
// backup block; storages always renders (possibly empty) once it does.
func backupBody(locals *Locals) map[string]interface{} {
	backup := locals.Spec.GetBackup()
	if backup == nil {
		return nil
	}

	body := map[string]interface{}{
		"image":    vars.BackupImage,
		"storages": backupStoragesBody(locals),
	}

	if schedules := backupSchedulesBody(backup); len(schedules) > 0 {
		body["schedule"] = schedules
	}
	if pitr := pitrBody(backup); pitr != nil {
		body["pitr"] = pitr
	}
	if pullSecrets := imagePullSecretsBody(locals.Spec.GetImagePullSecrets()); pullSecrets != nil {
		body["imagePullSecrets"] = pullSecrets
	}

	return body
}

// backupStoragesBody is the twin of TF's backup_storages: a map of
// storage-name → nested object. Credentials never render inline — s3 and
// azure blocks carry only the deterministic credentialsSecret name.
func backupStoragesBody(locals *Locals) map[string]interface{} {
	storages := map[string]interface{}{}
	for _, storage := range locals.Spec.GetBackup().GetStorages() {
		credentialsSecret := locals.ClusterName + "-backup-" + storage.GetName()
		entry := map[string]interface{}{}

		switch {
		case storage.GetS3() != nil:
			entry["type"] = "s3"
			s3 := storage.GetS3()
			s3Body := map[string]interface{}{
				"bucket":            s3.GetBucket(),
				"credentialsSecret": credentialsSecret,
			}
			if s3.GetRegion() != "" {
				s3Body["region"] = s3.GetRegion()
			}
			if s3.GetPrefix() != "" {
				s3Body["prefix"] = s3.GetPrefix()
			}
			if s3.GetEndpointUrl() != "" {
				s3Body["endpointUrl"] = s3.GetEndpointUrl()
			}
			if s3.GetForcePathStyle() {
				s3Body["forcePathStyle"] = true
			}
			entry["s3"] = s3Body
		case storage.GetAzure() != nil:
			entry["type"] = "azure"
			azure := storage.GetAzure()
			azureBody := map[string]interface{}{
				"container":         azure.GetContainer(),
				"credentialsSecret": credentialsSecret,
			}
			if azure.GetPrefix() != "" {
				azureBody["prefix"] = azure.GetPrefix()
			}
			if azure.GetEndpointUrl() != "" {
				azureBody["endpointUrl"] = azure.GetEndpointUrl()
			}
			entry["azure"] = azureBody
		default:
			entry["type"] = "filesystem"
			if pvc := storage.GetPvc(); pvc != nil {
				volume := pvc.GetVolume()
				pvcBlock := map[string]interface{}{
					"accessModes": []string{"ReadWriteOnce"},
					"resources": map[string]interface{}{
						"requests": map[string]interface{}{
							"storage": volume.GetSize(),
						},
					},
				}
				if volume.GetStorageClass().GetValue() != "" {
					pvcBlock["storageClassName"] = volume.GetStorageClass().GetValue()
				}
				entry["volume"] = map[string]interface{}{
					"persistentVolumeClaim": pvcBlock,
				}
			}
		}

		if storage.VerifyTls != nil {
			entry["verifyTLS"] = storage.GetVerifyTls()
		}

		storages[storage.GetName()] = entry
	}
	return storages
}

// backupSchedulesBody is the twin of TF's backup_schedules.
func backupSchedulesBody(backup *kubernetesmysqlv1alpha1.KubernetesMysqlBackup) []interface{} {
	var schedules []interface{}
	for _, schedule := range backup.GetSchedules() {
		body := map[string]interface{}{
			"name":        schedule.GetName(),
			"schedule":    schedule.GetSchedule(),
			"storageName": schedule.GetStorageName(),
		}
		if schedule.Keep != nil {
			deleteFromStorage := true
			if schedule.DeleteFromStorage != nil {
				deleteFromStorage = schedule.GetDeleteFromStorage()
			}
			body["retention"] = map[string]interface{}{
				"type":              "count",
				"count":             int(schedule.GetKeep()),
				"deleteFromStorage": deleteFromStorage,
			}
		}
		schedules = append(schedules, body)
	}
	return schedules
}

// pitrBody is the twin of TF's pitr_body: rendered only when the spec
// enables point-in-time recovery.
func pitrBody(backup *kubernetesmysqlv1alpha1.KubernetesMysqlBackup) map[string]interface{} {
	pitr := backup.GetPitr()
	if pitr == nil || !pitr.GetEnabled() {
		return nil
	}
	timeBetweenUploads := vars.DefaultPitrTimeBetweenUploads
	if pitr.TimeBetweenUploads != nil {
		timeBetweenUploads = int(pitr.GetTimeBetweenUploads())
	}
	return map[string]interface{}{
		"enabled":            true,
		"storageName":        pitr.GetStorageName(),
		"timeBetweenUploads": timeBetweenUploads,
	}
}

// logcollectorBody is the twin of TF's logcollector_body: always rendered
// — an absent spec block still yields enabled=true plus the pinned image
// (the upstream default posture).
func logcollectorBody(spec *kubernetesmysqlv1alpha1.KubernetesMysqlSpec) map[string]interface{} {
	logCollector := spec.GetLogCollector()

	enabled := true
	if logCollector != nil && logCollector.Enabled != nil {
		enabled = logCollector.GetEnabled()
	}
	body := map[string]interface{}{
		"enabled": enabled,
		"image":   vars.LogCollectorImage,
	}
	if resources := containerResourcesBody(logCollector.GetResources()); resources != nil {
		body["resources"] = resources
	}
	return body
}

// unsafeFlagsBody is the twin of TF's unsafe_flags: only flags that are
// true render; the caller omits the whole block when none are.
func unsafeFlagsBody(unsafe *kubernetesmysqlv1alpha1.KubernetesMysqlUnsafe) map[string]interface{} {
	flags := map[string]interface{}{}
	if unsafe.GetClusterSize() {
		flags["pxcSize"] = true
	}
	if unsafe.GetTls() {
		flags["tls"] = true
	}
	if unsafe.GetProxySize() {
		flags["proxySize"] = true
	}
	if unsafe.GetBackupIfUnhealthy() {
		flags["backupIfUnhealthy"] = true
	}
	return flags
}
