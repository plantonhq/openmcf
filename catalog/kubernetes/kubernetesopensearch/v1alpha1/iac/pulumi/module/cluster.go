package module

import (
	"slices"
	"strconv"

	kubernetesprovider "github.com/plantonhq/planton/catalog/kubernetes"
	kubernetesopensearchv1alpha1 "github.com/plantonhq/planton/catalog/kubernetes/kubernetesopensearch/v1alpha1"
	opensearchv1 "github.com/plantonhq/planton/pkg/kubernetes/kubernetestypes/opensearchoperator/kubernetes/opensearch/v1"
	kubernetesmeta "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// createCluster renders the opensearch.opster.io/v1 OpenSearchCluster
// resource with the typed crd2pulumi SDK (field/structure drift against the
// pinned CRD fails at compile time). Unset optionals are omitted entirely so
// the apiserver applies the CRD's own defaults — presence discipline mirrors
// the Terraform module's null-prune rendering byte for byte.
func createCluster(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependencies []pulumi.ResourceOption,
) (pulumi.Resource, error) {
	spec := locals.Spec

	clusterSpec := opensearchv1.OpenSearchClusterSpecArgs{
		General:   buildGeneral(locals),
		NodePools: buildNodePools(spec.GetNodePools()),
	}

	if bootstrap := buildBootstrap(spec.GetBootstrap()); bootstrap != nil {
		clusterSpec.Bootstrap = bootstrap
	}
	if security := buildSecurity(spec.GetSecurity()); security != nil {
		clusterSpec.Security = security
	}
	if dashboards := buildDashboards(locals); dashboards != nil {
		clusterSpec.Dashboards = dashboards
	}

	return opensearchv1.NewOpenSearchCluster(ctx, locals.ClusterName,
		&opensearchv1.OpenSearchClusterArgs{
			Metadata: kubernetesmeta.ObjectMetaArgs{
				Name:      pulumi.String(locals.ClusterName),
				Namespace: pulumi.String(locals.Namespace),
				Labels:    pulumi.ToStringMap(locals.Labels),
			},
			Spec: clusterSpec,
		}, append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, dependencies...)...)
}

// buildGeneral renders spec.general. Module-owned constants:
// serviceName is ALWAYS the cluster name (the operator names the main
// Service after it — the exported service_name contract depends on this)
// and vendor is always "opensearch".
func buildGeneral(locals *Locals) opensearchv1.OpenSearchClusterSpecGeneralArgs {
	spec := locals.Spec

	general := opensearchv1.OpenSearchClusterSpecGeneralArgs{
		HttpPort:    pulumi.Int(locals.HttpPort),
		ServiceName: pulumi.String(locals.ClusterName),
		Vendor:      pulumi.String(vars.Vendor),
		Version:     pulumi.String(spec.GetVersion()),
	}

	// The CRD/operator default for setVMMaxMapCount is FALSE (no init
	// container); the spec default is true — so true is rendered
	// explicitly and false is omitted.
	if spec.SetVmMaxMapCount == nil || spec.GetSetVmMaxMapCount() {
		general.SetVMMaxMapCount = pulumi.Bool(true)
	}
	if spec.GetDrainDataNodes() {
		general.DrainDataNodes = pulumi.Bool(true)
	}
	if len(spec.GetAdditionalConfig()) > 0 {
		general.AdditionalConfig = pulumi.ToStringMap(spec.GetAdditionalConfig())
	}
	if len(spec.GetServiceAnnotations()) > 0 {
		general.Annotations = pulumi.ToStringMap(spec.GetServiceAnnotations())
	}
	if len(spec.GetPluginsList()) > 0 {
		general.PluginsList = pulumi.ToStringArray(spec.GetPluginsList())
	}

	// The CRD's ImageSpec takes ONE image string (not repo/tag fields) —
	// the shared ContainerImage folds into `repo:tag`.
	if repo := spec.GetImage().GetRepo(); repo != "" {
		image := repo
		if tag := spec.GetImage().GetTag(); tag != "" {
			image = repo + ":" + tag
		}
		general.Image = pulumi.String(image)
	}
	// image.pull_secret_name joins image_pull_secrets (deduplicated) — a
	// private image override naturally travels with its own pull secret,
	// so the shared ContainerImage's third field is honored, never dead.
	pullSecretNames := spec.GetImagePullSecrets()
	if name := spec.GetImage().GetPullSecretName(); name != "" && !slices.Contains(pullSecretNames, name) {
		pullSecretNames = append(pullSecretNames, name)
	}
	if len(pullSecretNames) > 0 {
		pullSecrets := opensearchv1.OpenSearchClusterSpecGeneralImagePullSecretsArray{}
		for _, name := range pullSecretNames {
			pullSecrets = append(pullSecrets,
				opensearchv1.OpenSearchClusterSpecGeneralImagePullSecretsArgs{
					Name: pulumi.String(name),
				})
		}
		general.ImagePullSecrets = pullSecrets
	}

	if len(spec.GetKeystore()) > 0 {
		keystore := opensearchv1.OpenSearchClusterSpecGeneralKeystoreArray{}
		for _, entry := range spec.GetKeystore() {
			entryArgs := opensearchv1.OpenSearchClusterSpecGeneralKeystoreArgs{
				Secret: opensearchv1.OpenSearchClusterSpecGeneralKeystoreSecretArgs{
					Name: pulumi.String(entry.GetSecret().GetValue()),
				},
			}
			if len(entry.GetKeyMappings()) > 0 {
				entryArgs.KeyMappings = pulumi.ToStringMap(entry.GetKeyMappings())
			}
			keystore = append(keystore, entryArgs)
		}
		general.Keystore = keystore
	}

	if len(spec.GetSnapshotRepositories()) > 0 {
		repos := opensearchv1.OpenSearchClusterSpecGeneralSnapshotRepositoriesArray{}
		for _, repo := range spec.GetSnapshotRepositories() {
			repoArgs := opensearchv1.OpenSearchClusterSpecGeneralSnapshotRepositoriesArgs{
				Name: pulumi.String(repo.GetName()),
				Type: pulumi.String(repo.GetType()),
			}
			if len(repo.GetSettings()) > 0 {
				repoArgs.Settings = pulumi.ToStringMap(repo.GetSettings())
			}
			repos = append(repos, repoArgs)
		}
		general.SnapshotRepositories = repos
	}

	if monitoring := spec.GetMonitoring(); monitoring.GetEnabled() {
		monitoringArgs := opensearchv1.OpenSearchClusterSpecGeneralMonitoringArgs{
			Enable: pulumi.Bool(true),
		}
		if monitoring.GetScrapeInterval() != "" {
			monitoringArgs.ScrapeInterval = pulumi.String(monitoring.GetScrapeInterval())
		}
		if monitoring.GetMonitoringUserSecret().GetValue() != "" {
			monitoringArgs.MonitoringUserSecret = pulumi.String(monitoring.GetMonitoringUserSecret().GetValue())
		}
		if monitoring.GetPluginUrl() != "" {
			monitoringArgs.PluginUrl = pulumi.String(monitoring.GetPluginUrl())
		}
		general.Monitoring = monitoringArgs
	}

	if len(spec.GetAdditionalVolumes()) > 0 {
		volumes := opensearchv1.OpenSearchClusterSpecGeneralAdditionalVolumesArray{}
		for _, volume := range spec.GetAdditionalVolumes() {
			volumeArgs := opensearchv1.OpenSearchClusterSpecGeneralAdditionalVolumesArgs{
				Name: pulumi.String(volume.GetName()),
				Path: pulumi.String(volume.GetPath()),
			}
			if volume.GetSubPath() != "" {
				volumeArgs.SubPath = pulumi.String(volume.GetSubPath())
			}
			if volume.GetSecretName() != "" {
				volumeArgs.Secret = opensearchv1.OpenSearchClusterSpecGeneralAdditionalVolumesSecretArgs{
					SecretName: pulumi.String(volume.GetSecretName()),
				}
			}
			if volume.GetConfigMapName() != "" {
				volumeArgs.ConfigMap = opensearchv1.OpenSearchClusterSpecGeneralAdditionalVolumesConfigMapArgs{
					Name: pulumi.String(volume.GetConfigMapName()),
				}
			}
			if volume.GetRestartPods() {
				volumeArgs.RestartPods = pulumi.Bool(true)
			}
			volumes = append(volumes, volumeArgs)
		}
		general.AdditionalVolumes = volumes
	}

	return general
}

func buildBootstrap(bootstrap *kubernetesopensearchv1alpha1.KubernetesOpenSearchBootstrap) opensearchv1.OpenSearchClusterSpecBootstrapPtrInput {
	if bootstrap == nil {
		return nil
	}

	bootstrapArgs := opensearchv1.OpenSearchClusterSpecBootstrapArgs{}
	hasAny := false

	if limits, requests := resourceMaps(bootstrap.GetResources()); limits != nil || requests != nil {
		bootstrapArgs.Resources = opensearchv1.OpenSearchClusterSpecBootstrapResourcesArgs{
			Limits:   limits,
			Requests: requests,
		}
		hasAny = true
	}
	if bootstrap.GetJvm() != "" {
		bootstrapArgs.Jvm = pulumi.String(bootstrap.GetJvm())
		hasAny = true
	}
	if len(bootstrap.GetAdditionalConfig()) > 0 {
		bootstrapArgs.AdditionalConfig = pulumi.ToStringMap(bootstrap.GetAdditionalConfig())
		hasAny = true
	}

	if !hasAny {
		return nil
	}
	return bootstrapArgs
}

// buildNodePools renders one nodePools entry per spec.node_pools pool —
// component is the pool name (StatefulSets become `<cluster>-<pool>`).
func buildNodePools(pools []*kubernetesopensearchv1alpha1.KubernetesOpenSearchNodePool) opensearchv1.OpenSearchClusterSpecNodePoolsArray {
	nodePools := opensearchv1.OpenSearchClusterSpecNodePoolsArray{}
	for _, pool := range pools {
		poolArgs := opensearchv1.OpenSearchClusterSpecNodePoolsArgs{
			Component: pulumi.String(pool.GetName()),
			Replicas:  pulumi.Int(int(pool.GetReplicas())),
			Roles:     pulumi.ToStringArray(pool.GetRoles()),
		}

		if limits, requests := resourceMaps(pool.GetResources()); limits != nil || requests != nil {
			poolArgs.Resources = opensearchv1.OpenSearchClusterSpecNodePoolsResourcesArgs{
				Limits:   limits,
				Requests: requests,
			}
		}
		if pool.GetJvm() != "" {
			poolArgs.Jvm = pulumi.String(pool.GetJvm())
		}
		if pool.GetDiskSize() != "" {
			poolArgs.DiskSize = pulumi.String(pool.GetDiskSize())
		}

		if persistence := pool.GetPersistence(); persistence != nil {
			persistenceArgs := opensearchv1.OpenSearchClusterSpecNodePoolsPersistenceArgs{}
			if pvc := persistence.GetPvc(); pvc != nil {
				// The CRD's PVCSource key is `storageClass` (not
				// storageClassName); accessModes is required by the
				// operator's PVC template — pinned ReadWriteOnce.
				pvcArgs := opensearchv1.OpenSearchClusterSpecNodePoolsPersistencePvcArgs{
					AccessModes: pulumi.ToStringArray([]string{"ReadWriteOnce"}),
				}
				if pvc.GetStorageClass().GetValue() != "" {
					pvcArgs.StorageClass = pulumi.String(pvc.GetStorageClass().GetValue())
				}
				persistenceArgs.Pvc = pvcArgs
			}
			if emptyDir := persistence.GetEmptyDir(); emptyDir != nil {
				emptyDirArgs := opensearchv1.OpenSearchClusterSpecNodePoolsPersistenceEmptyDirArgs{}
				if emptyDir.GetSizeLimit() != "" {
					emptyDirArgs.SizeLimit = pulumi.String(emptyDir.GetSizeLimit())
				}
				persistenceArgs.EmptyDir = emptyDirArgs
			}
			poolArgs.Persistence = persistenceArgs
		}

		if len(pool.GetAdditionalConfig()) > 0 {
			poolArgs.AdditionalConfig = pulumi.ToStringMap(pool.GetAdditionalConfig())
		}
		if len(pool.GetNodeSelector()) > 0 {
			poolArgs.NodeSelector = pulumi.ToStringMap(pool.GetNodeSelector())
		}
		if len(pool.GetTolerations()) > 0 {
			poolArgs.Tolerations = buildTolerations(pool.GetTolerations())
		}

		if pdb := pool.GetPdb(); pdb != nil {
			pdbArgs := opensearchv1.OpenSearchClusterSpecNodePoolsPdbArgs{}
			if pdb.GetEnable() {
				pdbArgs.Enable = pulumi.Bool(true)
			}
			if pdb.GetMinAvailable() != "" {
				pdbArgs.MinAvailable = intOrString(pdb.GetMinAvailable())
			}
			if pdb.GetMaxUnavailable() != "" {
				pdbArgs.MaxUnavailable = intOrString(pdb.GetMaxUnavailable())
			}
			poolArgs.Pdb = pdbArgs
		}

		nodePools = append(nodePools, poolArgs)
	}
	return nodePools
}

func buildSecurity(security *kubernetesopensearchv1alpha1.KubernetesOpenSearchSecurity) opensearchv1.OpenSearchClusterSpecSecurityPtrInput {
	if security == nil {
		return nil
	}

	securityArgs := opensearchv1.OpenSearchClusterSpecSecurityArgs{}
	hasAny := false

	tlsArgs := opensearchv1.OpenSearchClusterSpecSecurityTlsArgs{}
	hasTls := false

	if transport := security.GetTransportTls(); transport != nil {
		transportArgs := opensearchv1.OpenSearchClusterSpecSecurityTlsTransportArgs{}
		// The CRD default for generate/perNode is FALSE (existing
		// certificates expected); the spec defaults are true — true is
		// rendered explicitly, false omitted.
		if transport.Generate == nil || transport.GetGenerate() {
			transportArgs.Generate = pulumi.Bool(true)
		}
		if transport.PerNode == nil || transport.GetPerNode() {
			transportArgs.PerNode = pulumi.Bool(true)
		}
		if transport.GetSecret().GetValue() != "" {
			transportArgs.Secret = opensearchv1.OpenSearchClusterSpecSecurityTlsTransportSecretArgs{
				Name: pulumi.String(transport.GetSecret().GetValue()),
			}
		}
		if transport.GetCaSecret().GetValue() != "" {
			transportArgs.CaSecret = opensearchv1.OpenSearchClusterSpecSecurityTlsTransportCaSecretArgs{
				Name: pulumi.String(transport.GetCaSecret().GetValue()),
			}
		}
		if len(transport.GetNodesDn()) > 0 {
			transportArgs.NodesDn = pulumi.ToStringArray(transport.GetNodesDn())
		}
		if len(transport.GetAdminDn()) > 0 {
			transportArgs.AdminDn = pulumi.ToStringArray(transport.GetAdminDn())
		}
		tlsArgs.Transport = transportArgs
		hasTls = true
	}

	if http := security.GetHttpTls(); http != nil {
		httpArgs := opensearchv1.OpenSearchClusterSpecSecurityTlsHttpArgs{}
		if http.Generate == nil || http.GetGenerate() {
			httpArgs.Generate = pulumi.Bool(true)
		}
		if http.GetSecret().GetValue() != "" {
			httpArgs.Secret = opensearchv1.OpenSearchClusterSpecSecurityTlsHttpSecretArgs{
				Name: pulumi.String(http.GetSecret().GetValue()),
			}
		}
		tlsArgs.Http = httpArgs
		hasTls = true
	}

	if hasTls {
		securityArgs.Tls = tlsArgs
		hasAny = true
	}

	if config := security.GetConfig(); config != nil {
		configArgs := opensearchv1.OpenSearchClusterSpecSecurityConfigArgs{}
		hasConfig := false
		if config.GetSecurityConfigSecret().GetValue() != "" {
			configArgs.SecurityConfigSecret = opensearchv1.OpenSearchClusterSpecSecurityConfigSecurityConfigSecretArgs{
				Name: pulumi.String(config.GetSecurityConfigSecret().GetValue()),
			}
			hasConfig = true
		}
		if config.GetAdminSecret().GetValue() != "" {
			configArgs.AdminSecret = opensearchv1.OpenSearchClusterSpecSecurityConfigAdminSecretArgs{
				Name: pulumi.String(config.GetAdminSecret().GetValue()),
			}
			hasConfig = true
		}
		if config.GetAdminCredentialsSecret().GetValue() != "" {
			configArgs.AdminCredentialsSecret = opensearchv1.OpenSearchClusterSpecSecurityConfigAdminCredentialsSecretArgs{
				Name: pulumi.String(config.GetAdminCredentialsSecret().GetValue()),
			}
			hasConfig = true
		}
		if hasConfig {
			securityArgs.Config = configArgs
			hasAny = true
		}
	}

	if !hasAny {
		return nil
	}
	return securityArgs
}

func buildDashboards(locals *Locals) opensearchv1.OpenSearchClusterSpecDashboardsPtrInput {
	dashboards := locals.Spec.GetDashboards()
	if !dashboards.GetEnabled() {
		return nil
	}

	replicas := 1
	if dashboards.Replicas != nil && dashboards.GetReplicas() > 0 {
		replicas = int(dashboards.GetReplicas())
	}
	// Dashboards version defaults to the CLUSTER version (module-owned:
	// Dashboards refuses mismatched clusters, and the CRD's version field
	// is required).
	version := dashboards.GetVersion()
	if version == "" {
		version = locals.Spec.GetVersion()
	}

	dashboardsArgs := opensearchv1.OpenSearchClusterSpecDashboardsArgs{
		Enable:   pulumi.Bool(true),
		Replicas: pulumi.Int(replicas),
		Version:  pulumi.String(version),
	}

	if limits, requests := resourceMaps(dashboards.GetResources()); limits != nil || requests != nil {
		dashboardsArgs.Resources = opensearchv1.OpenSearchClusterSpecDashboardsResourcesArgs{
			Limits:   limits,
			Requests: requests,
		}
	}

	if tls := dashboards.GetTls(); tls.GetEnable() {
		tlsArgs := opensearchv1.OpenSearchClusterSpecDashboardsTlsArgs{
			Enable: pulumi.Bool(true),
		}
		if tls.Generate == nil || tls.GetGenerate() {
			tlsArgs.Generate = pulumi.Bool(true)
		}
		if tls.GetSecret().GetValue() != "" {
			tlsArgs.Secret = opensearchv1.OpenSearchClusterSpecDashboardsTlsSecretArgs{
				Name: pulumi.String(tls.GetSecret().GetValue()),
			}
		}
		dashboardsArgs.Tls = tlsArgs
	}

	if dashboards.GetBasePath() != "" {
		dashboardsArgs.BasePath = pulumi.String(dashboards.GetBasePath())
	}
	if len(dashboards.GetAdditionalConfig()) > 0 {
		dashboardsArgs.AdditionalConfig = pulumi.ToStringMap(dashboards.GetAdditionalConfig())
	}
	if dashboards.GetOpensearchCredentialsSecret().GetValue() != "" {
		dashboardsArgs.OpensearchCredentialsSecret = opensearchv1.OpenSearchClusterSpecDashboardsOpensearchCredentialsSecretArgs{
			Name: pulumi.String(dashboards.GetOpensearchCredentialsSecret().GetValue()),
		}
	}

	if service := dashboards.GetService(); service != nil {
		serviceType := "ClusterIP"
		if service.Type != nil && service.GetType() != "" {
			serviceType = service.GetType()
		}
		serviceArgs := opensearchv1.OpenSearchClusterSpecDashboardsServiceArgs{
			Type: pulumi.String(serviceType),
		}
		if len(service.GetLoadBalancerSourceRanges()) > 0 {
			serviceArgs.LoadBalancerSourceRanges = pulumi.ToStringArray(service.GetLoadBalancerSourceRanges())
		}
		dashboardsArgs.Service = serviceArgs
	}

	if len(dashboards.GetPluginsList()) > 0 {
		dashboardsArgs.PluginsList = pulumi.ToStringArray(dashboards.GetPluginsList())
	}

	return dashboardsArgs
}

// resourceMaps translates the shared ContainerResources into the CRD's
// limits/requests quantity maps; empty arms return nil so the rendered CR
// omits them.
func resourceMaps(resources *kubernetesprovider.ContainerResources) (pulumi.Map, pulumi.Map) {
	if resources == nil {
		return nil, nil
	}
	var limits, requests pulumi.Map
	if l := resources.GetLimits(); l != nil {
		limits = pulumi.Map{
			"cpu":    pulumi.String(l.GetCpu()),
			"memory": pulumi.String(l.GetMemory()),
		}
	}
	if r := resources.GetRequests(); r != nil {
		requests = pulumi.Map{
			"cpu":    pulumi.String(r.GetCpu()),
			"memory": pulumi.String(r.GetMemory()),
		}
	}
	return limits, requests
}

func buildTolerations(tolerations []*kubernetesprovider.WorkloadToleration) opensearchv1.OpenSearchClusterSpecNodePoolsTolerationsArray {
	out := opensearchv1.OpenSearchClusterSpecNodePoolsTolerationsArray{}
	for _, toleration := range tolerations {
		tolerationArgs := opensearchv1.OpenSearchClusterSpecNodePoolsTolerationsArgs{}
		if toleration.GetKey() != "" {
			tolerationArgs.Key = pulumi.String(toleration.GetKey())
		}
		if toleration.GetOperator() != "" {
			tolerationArgs.Operator = pulumi.String(toleration.GetOperator())
		}
		if toleration.GetValue() != "" {
			tolerationArgs.Value = pulumi.String(toleration.GetValue())
		}
		if toleration.GetEffect() != "" {
			tolerationArgs.Effect = pulumi.String(toleration.GetEffect())
		}
		if toleration.TolerationSeconds != nil {
			tolerationArgs.TolerationSeconds = pulumi.Int(int(toleration.GetTolerationSeconds()))
		}
		out = append(out, tolerationArgs)
	}
	return out
}

// intOrString renders a PDB bound with intstr semantics: a string that
// parses as an integer renders as a YAML number ("2" → 2), anything else
// (percentages like "25%") as a string — the Terraform twin applies
// `try(tonumber(v), v)` for the identical result.
func intOrString(value string) pulumi.Input {
	if parsed, err := strconv.Atoi(value); err == nil {
		return pulumi.Int(parsed)
	}
	return pulumi.String(value)
}
