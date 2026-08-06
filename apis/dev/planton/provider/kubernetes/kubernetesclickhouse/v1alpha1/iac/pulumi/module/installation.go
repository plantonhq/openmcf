package module

import (
	kubernetesprovider "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes"
	kubernetesclickhousev1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetesclickhouse/v1alpha1"
	clickhousev1 "github.com/plantonhq/planton/pkg/kubernetes/kubernetestypes/altinityoperator/kubernetes/clickhouse/v1"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	kubernetesmeta "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// installation renders the clickhouse.altinity.com/v1 ClickHouseInstallation
// with the typed crd2pulumi SDK (field/structure drift against the pinned
// CRD fails at compile time). Unset optionals are omitted entirely so the
// operator applies its own defaults — presence discipline mirrors the
// Terraform module's null-prune rendering byte for byte.
func installation(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	createdAuthSecret *kubernetescorev1.Secret,
	dependencies []pulumi.ResourceOption,
) (pulumi.Resource, error) {
	chiSpec := clickhousev1.ClickHouseInstallationSpecArgs{
		Configuration: buildConfiguration(locals, createdAuthSecret),
		Templates:     buildTemplates(locals),
	}

	if defaults := buildDefaults(locals); defaults != nil {
		chiSpec.Defaults = defaults
	}
	// The CHI `stop` verb: "yes" scales every host StatefulSet to zero
	// keeping all PVCs; omitted (operator default "no") otherwise.
	if locals.Spec.GetStopped() {
		chiSpec.Stop = crdStringBool("yes")
	}

	return clickhousev1.NewClickHouseInstallation(ctx, locals.ChiName,
		&clickhousev1.ClickHouseInstallationArgs{
			Metadata: kubernetesmeta.ObjectMetaArgs{
				Name:      pulumi.String(locals.ChiName),
				Namespace: pulumi.String(locals.Namespace),
				Labels:    pulumi.ToStringMap(locals.Labels),
			},
			Spec: chiSpec,
		}, append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, dependencies...)...)
}

// buildConfiguration renders configuration: the single cluster's layout,
// the coordination (zookeeper) section, users/profiles/quotas and the
// settings/files pass-through maps.
func buildConfiguration(locals *Locals,
	createdAuthSecret *kubernetescorev1.Secret,
) clickhousev1.ClickHouseInstallationSpecConfigurationArgs {
	spec := locals.Spec

	cluster := clickhousev1.ClickHouseInstallationSpecConfigurationClustersArgs{
		Name: pulumi.String(locals.ClusterName),
		Layout: clickhousev1.ClickHouseInstallationSpecConfigurationClustersLayoutArgs{
			ShardsCount:   pulumi.Int(locals.Shards),
			ReplicasCount: pulumi.Int(locals.Replicas),
		},
	}

	// Operator-generated shared secret for distributed queries between
	// this cluster's own hosts (spec default true); meaningful only when
	// there is more than one host.
	autoInterNodeSecret := spec.AutoInterNodeSecret == nil || spec.GetAutoInterNodeSecret()
	if autoInterNodeSecret && locals.Shards*locals.Replicas > 1 {
		cluster.Secret = clickhousev1.ClickHouseInstallationSpecConfigurationClustersSecretArgs{
			Auto: crdStringBool("true"),
		}
	}

	// One operator-managed PDB per cluster; spec default 1 equals the
	// operator default, rendered explicitly to keep the contract visible.
	pdbMaxUnavailable := 1
	if spec.PdbMaxUnavailable != nil {
		pdbMaxUnavailable = int(spec.GetPdbMaxUnavailable())
	}
	cluster.PdbMaxUnavailable = pulumi.Int(pdbMaxUnavailable)

	configuration := clickhousev1.ClickHouseInstallationSpecConfigurationArgs{
		Clusters: clickhousev1.ClickHouseInstallationSpecConfigurationClustersArray{cluster},
	}

	if zookeeper := buildZookeeper(locals); zookeeper != nil {
		configuration.Zookeeper = zookeeper
	}
	if users := buildUsers(spec.GetUsers(), createdAuthSecret); users != nil {
		configuration.Users = users
	}
	if profiles := buildNamedSettings(spec.GetProfiles()); profiles != nil {
		configuration.Profiles = profiles
	}
	if quotas := buildNamedSettings(spec.GetQuotas()); quotas != nil {
		configuration.Quotas = quotas
	}
	if len(spec.GetSettings()) > 0 {
		configuration.Settings = stringValueMap(spec.GetSettings())
	}
	if len(spec.GetFiles()) > 0 {
		configuration.Files = stringValueMap(spec.GetFiles())
	}

	return configuration
}

// buildZookeeper renders the coordination section. The managed Keeper is
// wired through the CHI's native keeper reference — name only, namespace
// omitted (same namespace), endpoints resolved by the operator itself.
// External coordination enumerates the ensemble nodes from the spec.
func buildZookeeper(locals *Locals) clickhousev1.ClickHouseInstallationSpecConfigurationZookeeperPtrInput {
	if locals.DeployKeeper {
		return clickhousev1.ClickHouseInstallationSpecConfigurationZookeeperArgs{
			Keeper: clickhousev1.ClickHouseInstallationSpecConfigurationZookeeperKeeperArgs{
				Name: pulumi.String(locals.KeeperName),
			},
		}
	}

	coordination := locals.Spec.GetCoordination()
	coordinationType := coordination.GetType()
	if coordinationType != kubernetesclickhousev1alpha1.KubernetesClickHouseCoordination_external_keeper &&
		coordinationType != kubernetesclickhousev1alpha1.KubernetesClickHouseCoordination_external_zookeeper {
		return nil
	}

	external := coordination.GetExternal()
	nodes := clickhousev1.ClickHouseInstallationSpecConfigurationZookeeperNodesArray{}
	for _, node := range external.GetNodes() {
		port := 2181
		if node.Port != nil && node.GetPort() > 0 {
			port = int(node.GetPort())
		}
		nodes = append(nodes, clickhousev1.ClickHouseInstallationSpecConfigurationZookeeperNodesArgs{
			Host: pulumi.String(node.GetHost()),
			Port: pulumi.Int(port),
		})
	}

	zookeeper := clickhousev1.ClickHouseInstallationSpecConfigurationZookeeperArgs{
		Nodes: nodes,
	}
	if external.GetRoot() != "" {
		zookeeper.Root = pulumi.String(external.GetRoot())
	}
	if external.GetIdentity() != "" {
		zookeeper.Identity = pulumi.ToSecret(pulumi.String(external.GetIdentity())).(pulumi.StringOutput)
	}
	return zookeeper
}

// buildUsers renders the CHI users section — path-keyed entries per the
// upstream model. Passwords reference the auth Secret (key = user name) via
// valueFrom.secretKeyRef; the Secret's resource output carries the
// dependency so the CHI never reconciles before its Secret exists.
func buildUsers(users []*kubernetesclickhousev1alpha1.KubernetesClickHouseUser,
	createdAuthSecret *kubernetescorev1.Secret,
) pulumi.MapInput {
	if len(users) == 0 {
		return nil
	}

	entries := pulumi.Map{}
	for _, user := range users {
		name := user.GetName()
		entries[name+"/password"] = pulumi.Map{
			"valueFrom": pulumi.Map{
				"secretKeyRef": pulumi.Map{
					"name": createdAuthSecret.Metadata.Name(),
					"key":  pulumi.String(name),
				},
			},
		}
		if user.GetProfile() != "" {
			entries[name+"/profile"] = pulumi.String(user.GetProfile())
		}
		if user.GetQuota() != "" {
			entries[name+"/quota"] = pulumi.String(user.GetQuota())
		}
		if len(user.GetNetworks()) > 0 {
			entries[name+"/networks/ip"] = pulumi.ToStringArray(user.GetNetworks())
		}
		if len(user.GetGrants()) > 0 {
			entries[name+"/grants/query"] = pulumi.ToStringArray(user.GetGrants())
		}
		if user.GetAccessManagement() {
			entries[name+"/access_management"] = pulumi.Int(1)
		}
		for path, value := range user.GetSettings() {
			entries[name+"/"+path] = pulumi.String(value)
		}
	}
	return entries
}

// buildNamedSettings flattens named bundles into the CHI's path-keyed
// profiles/quotas shape: "<bundle-name>/<path>" = value.
func buildNamedSettings(bundles []*kubernetesclickhousev1alpha1.KubernetesClickHouseNamedSettings) pulumi.MapInput {
	if len(bundles) == 0 {
		return nil
	}

	entries := pulumi.Map{}
	for _, bundle := range bundles {
		for path, value := range bundle.GetSettings() {
			entries[bundle.GetName()+"/"+path] = pulumi.String(value)
		}
	}
	return entries
}

// buildDefaults renders defaults: the template wiring and — only when
// retain_volumes_on_delete — the Retain PVC reclaim policy (the operator
// default is Delete, so false is omitted).
func buildDefaults(locals *Locals) clickhousev1.ClickHouseInstallationSpecDefaultsPtrInput {
	templates := clickhousev1.ClickHouseInstallationSpecDefaultsTemplatesArgs{
		PodTemplate:             pulumi.String("server"),
		DataVolumeClaimTemplate: pulumi.String("data"),
	}
	if locals.Spec.GetLogDiskSize() != "" {
		templates.LogVolumeClaimTemplate = pulumi.String("logs")
	}

	defaults := clickhousev1.ClickHouseInstallationSpecDefaultsArgs{
		Templates: templates,
	}
	if locals.Spec.GetRetainVolumesOnDelete() {
		defaults.StorageManagement = clickhousev1.ClickHouseInstallationSpecDefaultsStorageManagementArgs{
			ReclaimPolicy: pulumi.String("Retain"),
		}
	}
	return defaults
}

// buildTemplates renders the pod, volume-claim and (conditionally) service
// templates every host StatefulSet is stamped from.
func buildTemplates(locals *Locals) clickhousev1.ClickHouseInstallationSpecTemplatesArgs {
	spec := locals.Spec

	templates := clickhousev1.ClickHouseInstallationSpecTemplatesArgs{
		PodTemplates: clickhousev1.ClickHouseInstallationSpecTemplatesPodTemplatesArray{
			buildServerPodTemplate(locals),
		},
		VolumeClaimTemplates: buildVolumeClaimTemplates(spec),
	}

	// The cluster-wide client Service template exists ONLY to carry
	// service_annotations; without them the operator's own default
	// service (already ClusterIP) needs no override. The template's spec
	// is copied verbatim by the operator — no port defaulting — so the
	// standard interface ports are declared explicitly.
	if len(spec.GetServiceAnnotations()) > 0 {
		templates.ServiceTemplates = clickhousev1.ClickHouseInstallationSpecTemplatesServiceTemplatesArray{
			clickhousev1.ClickHouseInstallationSpecTemplatesServiceTemplatesArgs{
				Name:         pulumi.String("client"),
				GenerateName: pulumi.String("clickhouse-{chi}"),
				Metadata: pulumi.Map{
					"annotations": stringValueMap(spec.GetServiceAnnotations()),
				},
				Spec: pulumi.Map{
					"type": pulumi.String("ClusterIP"),
					"ports": pulumi.Array{
						pulumi.Map{
							"name": pulumi.String("http"),
							"port": pulumi.Int(vars.HttpPort),
						},
						pulumi.Map{
							"name": pulumi.String("tcp"),
							"port": pulumi.Int(vars.TcpPort),
						},
					},
				},
			},
		}
	}

	return templates
}

// buildServerPodTemplate renders the "server" pod template: the clickhouse
// container (resolved image, resources) plus placement — nodeSelector,
// tolerations, imagePullSecrets, and the ShardAntiAffinity distribution
// when spread_replicas_across_nodes.
func buildServerPodTemplate(locals *Locals) clickhousev1.ClickHouseInstallationSpecTemplatesPodTemplatesArgs {
	spec := locals.Spec

	container := pulumi.Map{
		"name":  pulumi.String("clickhouse"),
		"image": pulumi.String(locals.Image),
	}
	if resources := containerResourcesMap(spec.GetResources()); resources != nil {
		container["resources"] = resources
	}

	podSpec := pulumi.Map{
		"containers": pulumi.Array{container},
	}
	if len(spec.GetNodeSelector()) > 0 {
		podSpec["nodeSelector"] = stringValueMap(spec.GetNodeSelector())
	}
	if tolerations := tolerationsArray(spec.GetTolerations()); tolerations != nil {
		podSpec["tolerations"] = tolerations
	}
	if len(spec.GetImagePullSecrets()) > 0 {
		pullSecrets := pulumi.Array{}
		for _, name := range spec.GetImagePullSecrets() {
			pullSecrets = append(pullSecrets, pulumi.Map{"name": pulumi.String(name)})
		}
		podSpec["imagePullSecrets"] = pullSecrets
	}

	podTemplate := clickhousev1.ClickHouseInstallationSpecTemplatesPodTemplatesArgs{
		Name: pulumi.String("server"),
		Spec: podSpec,
	}
	if spec.GetSpreadReplicasAcrossNodes() {
		podTemplate.PodDistribution = clickhousev1.ClickHouseInstallationSpecTemplatesPodTemplatesPodDistributionArray{
			clickhousev1.ClickHouseInstallationSpecTemplatesPodTemplatesPodDistributionArgs{
				Type: pulumi.String("ShardAntiAffinity"),
			},
		}
	}
	return podTemplate
}

// buildVolumeClaimTemplates renders "data" (always) and "logs" (only when
// log_disk_size is set) claim templates.
func buildVolumeClaimTemplates(spec *kubernetesclickhousev1alpha1.KubernetesClickHouseSpec) clickhousev1.ClickHouseInstallationSpecTemplatesVolumeClaimTemplatesArray {
	claims := clickhousev1.ClickHouseInstallationSpecTemplatesVolumeClaimTemplatesArray{
		clickhousev1.ClickHouseInstallationSpecTemplatesVolumeClaimTemplatesArgs{
			Name: pulumi.String("data"),
			Spec: pvcSpecMap(spec.GetDiskSize(), spec.GetStorageClass().GetValue()),
		},
	}
	if spec.GetLogDiskSize() != "" {
		claims = append(claims, clickhousev1.ClickHouseInstallationSpecTemplatesVolumeClaimTemplatesArgs{
			Name: pulumi.String("logs"),
			Spec: pvcSpecMap(spec.GetLogDiskSize(), spec.GetStorageClass().GetValue()),
		})
	}
	return claims
}

// pvcSpecMap renders a PVC spec; empty storageClass is omitted so the
// cluster's default storage class applies.
func pvcSpecMap(size string, storageClass string) pulumi.Map {
	pvcSpec := pulumi.Map{
		"accessModes": pulumi.ToStringArray([]string{"ReadWriteOnce"}),
		"resources": pulumi.Map{
			"requests": pulumi.Map{
				"storage": pulumi.String(size),
			},
		},
	}
	if storageClass != "" {
		pvcSpec["storageClassName"] = pulumi.String(storageClass)
	}
	return pvcSpec
}

// containerResourcesMap translates the shared ContainerResources into the
// pod-spec requests/limits shape; nil when nothing is declared so the
// container omits the block entirely.
func containerResourcesMap(resources *kubernetesprovider.ContainerResources) pulumi.Map {
	if resources == nil {
		return nil
	}

	out := pulumi.Map{}
	if requests := quantityMap(resources.GetRequests().GetCpu(), resources.GetRequests().GetMemory()); requests != nil {
		out["requests"] = requests
	}
	if limits := quantityMap(resources.GetLimits().GetCpu(), resources.GetLimits().GetMemory()); limits != nil {
		out["limits"] = limits
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func quantityMap(cpu string, memory string) pulumi.Map {
	quantities := pulumi.Map{}
	if cpu != "" {
		quantities["cpu"] = pulumi.String(cpu)
	}
	if memory != "" {
		quantities["memory"] = pulumi.String(memory)
	}
	if len(quantities) == 0 {
		return nil
	}
	return quantities
}

// tolerationsArray renders the shared WorkloadToleration list into the
// pod-spec shape; nil when the list is empty.
func tolerationsArray(tolerations []*kubernetesprovider.WorkloadToleration) pulumi.Array {
	if len(tolerations) == 0 {
		return nil
	}

	out := pulumi.Array{}
	for _, toleration := range tolerations {
		entry := pulumi.Map{}
		if toleration.GetKey() != "" {
			entry["key"] = pulumi.String(toleration.GetKey())
		}
		if toleration.GetOperator() != "" {
			entry["operator"] = pulumi.String(toleration.GetOperator())
		}
		if toleration.GetValue() != "" {
			entry["value"] = pulumi.String(toleration.GetValue())
		}
		if toleration.GetEffect() != "" {
			entry["effect"] = pulumi.String(toleration.GetEffect())
		}
		if toleration.TolerationSeconds != nil {
			entry["tolerationSeconds"] = pulumi.Int(int(toleration.GetTolerationSeconds()))
		}
		out = append(out, entry)
	}
	return out
}

// stringValueMap lifts a proto string map into an untyped pulumi map (the
// generated CRD fields take MapInput, not StringMapInput).
func stringValueMap(values map[string]string) pulumi.Map {
	out := pulumi.Map{}
	for key, value := range values {
		out[key] = pulumi.String(value)
	}
	return out
}
