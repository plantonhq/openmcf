package module

import (
	kubernetesprovider "github.com/plantonhq/planton/catalog/kubernetes"
	kubernetesrabbitmqv1alpha1 "github.com/plantonhq/planton/catalog/kubernetes/kubernetesrabbitmq/v1alpha1"
	rabbitmqv1beta1 "github.com/plantonhq/planton/pkg/kubernetes/kubernetestypes/rabbitmqoperator/kubernetes/rabbitmq/v1beta1"
	kubernetesmeta "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// cluster renders the rabbitmq.com/v1beta1 RabbitmqCluster with the typed
// crd2pulumi SDK (field/structure drift against the pinned CRD fails at
// compile time). Unset optionals are omitted entirely so the operator
// applies its own defaults — presence discipline mirrors the Terraform
// module's null-prune rendering byte for byte.
//
// The CR declares BACKGROUND deletion propagation: the operator's own
// deletion finalizer is the cascade, and foreground propagation deadlocks
// against operators that keep reconciling children during deletion
// (verified live on sibling operator-owned CRs; Terraform twin:
// delete_cascade = "Background" on kubectl_manifest).
func cluster(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependencies []pulumi.ResourceOption,
) (pulumi.Resource, error) {
	spec := locals.Spec

	clusterSpec := rabbitmqv1beta1.RabbitmqClusterSpecArgs{
		Replicas: pulumi.Int(locals.Replicas),
	}

	if image := resolveImage(spec.GetImage()); image != "" {
		clusterSpec.Image = pulumi.String(image)
	}
	if pullSecrets := imagePullSecretsArray(pullSecretNames(spec)); pullSecrets != nil {
		clusterSpec.ImagePullSecrets = pullSecrets
	}
	if service := buildService(spec.GetService()); service != nil {
		clusterSpec.Service = service
	}
	clusterSpec.Persistence = buildPersistence(spec)
	if resources := buildResources(spec.GetResources()); resources != nil {
		clusterSpec.Resources = resources
	}
	if rabbitmq := buildRabbitmqConfiguration(spec.GetConfiguration()); rabbitmq != nil {
		clusterSpec.Rabbitmq = rabbitmq
	}
	if tls := buildTls(spec.GetTls()); tls != nil {
		clusterSpec.Tls = tls
	}
	if tolerations := buildTolerations(spec.GetTolerations()); tolerations != nil {
		clusterSpec.Tolerations = tolerations
	}
	if affinity := buildAffinity(locals); affinity != nil {
		clusterSpec.Affinity = affinity
	}
	// Value-based rendering (not presence-based) for the two knobs whose
	// spec defaults equal the operator defaults: Terraform's tfvars
	// contract flattens proto presence, so both engines render these only
	// when they DIFFER from the operator defaults — the common contract
	// that keeps the CR bodies byte-identical.
	if spec.TerminationGracePeriodSeconds != nil && spec.GetTerminationGracePeriodSeconds() != 604800 {
		clusterSpec.TerminationGracePeriodSeconds = pulumi.Int(int(spec.GetTerminationGracePeriodSeconds()))
	}
	if spec.DelayStartSeconds != nil && spec.GetDelayStartSeconds() != 30 {
		clusterSpec.DelayStartSeconds = pulumi.Int(int(spec.GetDelayStartSeconds()))
	}
	if spec.GetSkipPostDeploySteps() {
		clusterSpec.SkipPostDeploySteps = pulumi.Bool(true)
	}
	if spec.GetAutoEnableAllFeatureFlags() {
		clusterSpec.AutoEnableAllFeatureFlags = pulumi.Bool(true)
	}
	if secretBackend := buildSecretBackend(spec.GetSecretBackend()); secretBackend != nil {
		clusterSpec.SecretBackend = secretBackend
	}

	return rabbitmqv1beta1.NewRabbitmqCluster(ctx, locals.ClusterName,
		&rabbitmqv1beta1.RabbitmqClusterArgs{
			Metadata: kubernetesmeta.ObjectMetaArgs{
				Name:      pulumi.String(locals.ClusterName),
				Namespace: pulumi.String(locals.Namespace),
				Labels:    pulumi.ToStringMap(locals.Labels),
				Annotations: pulumi.StringMap{
					"pulumi.com/deletionPropagationPolicy": pulumi.String("background"),
				},
			},
			Spec: clusterSpec,
		}, append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, dependencies...)...)
}

// resolveImage joins repo:tag when either is set; empty when neither is —
// the operator (or its fleet-wide default_rabbitmq_image) then picks the
// image.
func resolveImage(image *kubernetesprovider.ContainerImage) string {
	if image.GetRepo() == "" && image.GetTag() == "" {
		return ""
	}
	repo := image.GetRepo()
	tag := image.GetTag()
	if repo == "" || tag == "" {
		// A lone repo or lone tag cannot resolve against the operator's
		// compiled-in default — render whatever was given verbatim so
		// the API surface stays honest (the CRD takes one image string).
		return repo + tag
	}
	return repo + ":" + tag
}

// pullSecretNames joins image_pull_secrets with image's own
// pull_secret_name, deduplicated — a private image override naturally
// travels with its own pull secret (Terraform twin:
// image_pull_secret_names in locals.tf).
func pullSecretNames(spec *kubernetesrabbitmqv1alpha1.KubernetesRabbitMqSpec) []string {
	names := append([]string{}, spec.GetImagePullSecrets()...)
	if extra := spec.GetImage().GetPullSecretName(); extra != "" {
		seen := false
		for _, name := range names {
			if name == extra {
				seen = true
				break
			}
		}
		if !seen {
			names = append(names, extra)
		}
	}
	return names
}

// imagePullSecretsArray renders the LocalObjectReference list; nil when
// empty.
func imagePullSecretsArray(names []string) rabbitmqv1beta1.RabbitmqClusterSpecImagePullSecretsArrayInput {
	if len(names) == 0 {
		return nil
	}
	out := rabbitmqv1beta1.RabbitmqClusterSpecImagePullSecretsArray{}
	for _, name := range names {
		out = append(out, rabbitmqv1beta1.RabbitmqClusterSpecImagePullSecretsArgs{
			Name: pulumi.String(name),
		})
	}
	return out
}

// buildService renders the client-Service shape; nil when the spec block
// is absent or carries only defaults (the operator's own ClusterIP default
// applies).
func buildService(service *kubernetesrabbitmqv1alpha1.KubernetesRabbitMqService) rabbitmqv1beta1.RabbitmqClusterSpecServicePtrInput {
	if service == nil {
		return nil
	}

	out := rabbitmqv1beta1.RabbitmqClusterSpecServiceArgs{}
	rendered := false

	switch service.GetType() {
	case kubernetesrabbitmqv1alpha1.KubernetesRabbitMqService_load_balancer:
		out.Type = pulumi.String("LoadBalancer")
		rendered = true
	case kubernetesrabbitmqv1alpha1.KubernetesRabbitMqService_node_port:
		out.Type = pulumi.String("NodePort")
		rendered = true
	}
	if len(service.GetAnnotations()) > 0 {
		out.Annotations = pulumi.ToStringMap(service.GetAnnotations())
		rendered = true
	}
	if len(service.GetLabels()) > 0 {
		out.Labels = pulumi.ToStringMap(service.GetLabels())
		rendered = true
	}
	switch service.GetIpFamilyPolicy() {
	case kubernetesrabbitmqv1alpha1.KubernetesRabbitMqService_single_stack:
		out.IpFamilyPolicy = pulumi.String("SingleStack")
		rendered = true
	case kubernetesrabbitmqv1alpha1.KubernetesRabbitMqService_prefer_dual_stack:
		out.IpFamilyPolicy = pulumi.String("PreferDualStack")
		rendered = true
	case kubernetesrabbitmqv1alpha1.KubernetesRabbitMqService_require_dual_stack:
		out.IpFamilyPolicy = pulumi.String("RequireDualStack")
		rendered = true
	}

	if !rendered {
		return nil
	}
	return out
}

// buildPersistence renders the storage contract. The ephemeral posture is
// the CR's own storage-0 + emptyDir mechanism; otherwise the resolved disk
// size (spec default 10Gi) and the optional storage class render as the
// persistent volume claim template.
func buildPersistence(spec *kubernetesrabbitmqv1alpha1.KubernetesRabbitMqSpec) rabbitmqv1beta1.RabbitmqClusterSpecPersistencePtrInput {
	if spec.GetEphemeral() {
		return rabbitmqv1beta1.RabbitmqClusterSpecPersistenceArgs{
			Storage:  pulumi.String("0"),
			EmptyDir: rabbitmqv1beta1.RabbitmqClusterSpecPersistenceEmptyDirArgs{},
		}
	}

	diskSize := spec.GetDiskSize()
	if diskSize == "" {
		diskSize = "10Gi"
	}
	persistence := rabbitmqv1beta1.RabbitmqClusterSpecPersistenceArgs{
		Storage: pulumi.String(diskSize),
	}
	if storageClass := spec.GetStorageClass().GetValue(); storageClass != "" {
		persistence.StorageClassName = pulumi.String(storageClass)
	}
	return persistence
}

// buildResources renders requests/limits; nil when nothing is declared so
// the operator's own defaults (1 CPU / 2Gi requests, 2 CPU / 2Gi limits)
// apply.
func buildResources(resources *kubernetesprovider.ContainerResources) rabbitmqv1beta1.RabbitmqClusterSpecResourcesPtrInput {
	if resources == nil {
		return nil
	}

	out := rabbitmqv1beta1.RabbitmqClusterSpecResourcesArgs{}
	rendered := false
	if requests := quantityMap(resources.GetRequests().GetCpu(), resources.GetRequests().GetMemory()); requests != nil {
		out.Requests = requests
		rendered = true
	}
	if limits := quantityMap(resources.GetLimits().GetCpu(), resources.GetLimits().GetMemory()); limits != nil {
		out.Limits = limits
		rendered = true
	}
	if !rendered {
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

// buildRabbitmqConfiguration renders the rabbitmq configuration block; nil
// when every field is empty.
func buildRabbitmqConfiguration(configuration *kubernetesrabbitmqv1alpha1.KubernetesRabbitMqConfiguration) rabbitmqv1beta1.RabbitmqClusterSpecRabbitmqPtrInput {
	if configuration == nil {
		return nil
	}

	out := rabbitmqv1beta1.RabbitmqClusterSpecRabbitmqArgs{}
	rendered := false
	if len(configuration.GetAdditionalPlugins()) > 0 {
		out.AdditionalPlugins = pulumi.ToStringArray(configuration.GetAdditionalPlugins())
		rendered = true
	}
	if configuration.GetAdditionalConfig() != "" {
		out.AdditionalConfig = pulumi.String(configuration.GetAdditionalConfig())
		rendered = true
	}
	if configuration.GetAdvancedConfig() != "" {
		out.AdvancedConfig = pulumi.String(configuration.GetAdvancedConfig())
		rendered = true
	}
	if configuration.GetEnvConfig() != "" {
		out.EnvConfig = pulumi.String(configuration.GetEnvConfig())
		rendered = true
	}
	if configuration.GetErlangInetConfig() != "" {
		out.ErlangInetConfig = pulumi.String(configuration.GetErlangInetConfig())
		rendered = true
	}
	if !rendered {
		return nil
	}
	return out
}

// buildTls renders the TLS block; nil when no certificate Secret is
// referenced (secret_name is required WITHIN the tls message — an absent
// tls block is the no-TLS posture, so an empty resolved value here means
// the block was never declared).
func buildTls(tls *kubernetesrabbitmqv1alpha1.KubernetesRabbitMqTls) rabbitmqv1beta1.RabbitmqClusterSpecTlsPtrInput {
	secretName := tls.GetSecretName().GetValue()
	if secretName == "" {
		return nil
	}

	out := rabbitmqv1beta1.RabbitmqClusterSpecTlsArgs{
		SecretName: pulumi.String(secretName),
	}
	if caSecretName := tls.GetCaSecretName().GetValue(); caSecretName != "" {
		out.CaSecretName = pulumi.String(caSecretName)
	}
	if tls.GetDisableNonTlsListeners() {
		out.DisableNonTLSListeners = pulumi.Bool(true)
	}
	return out
}

// buildTolerations renders the shared WorkloadToleration list into the
// CR's typed tolerations; nil when the list is empty.
func buildTolerations(tolerations []*kubernetesprovider.WorkloadToleration) rabbitmqv1beta1.RabbitmqClusterSpecTolerationsArrayInput {
	if len(tolerations) == 0 {
		return nil
	}

	out := rabbitmqv1beta1.RabbitmqClusterSpecTolerationsArray{}
	for _, toleration := range tolerations {
		entry := rabbitmqv1beta1.RabbitmqClusterSpecTolerationsArgs{}
		if toleration.GetKey() != "" {
			entry.Key = pulumi.String(toleration.GetKey())
		}
		if toleration.GetOperator() != "" {
			entry.Operator = pulumi.String(toleration.GetOperator())
		}
		if toleration.GetValue() != "" {
			entry.Value = pulumi.String(toleration.GetValue())
		}
		if toleration.GetEffect() != "" {
			entry.Effect = pulumi.String(toleration.GetEffect())
		}
		if toleration.TolerationSeconds != nil {
			entry.TolerationSeconds = pulumi.Int(int(toleration.GetTolerationSeconds()))
		}
		out = append(out, entry)
	}
	return out
}

// buildAffinity renders the placement contract:
//   - node_selector becomes REQUIRED node affinity with one In-match per
//     label (the CR has no nodeSelector field; behaviorally identical for
//     exact matches — Terraform twin renders the same shape),
//   - spread_across_nodes becomes REQUIRED pod anti-affinity on the
//     operator's own `app.kubernetes.io/name: <cluster>` pod label over
//     the hostname topology.
//
// nil when neither is asked for.
func buildAffinity(locals *Locals) rabbitmqv1beta1.RabbitmqClusterSpecAffinityPtrInput {
	spec := locals.Spec
	if len(spec.GetNodeSelector()) == 0 && !spec.GetSpreadAcrossNodes() {
		return nil
	}

	affinity := rabbitmqv1beta1.RabbitmqClusterSpecAffinityArgs{}

	if len(spec.GetNodeSelector()) > 0 {
		matchExpressions := rabbitmqv1beta1.RabbitmqClusterSpecAffinityNodeAffinityRequiredDuringSchedulingIgnoredDuringExecutionNodeSelectorTermsMatchExpressionsArray{}
		for key, value := range spec.GetNodeSelector() {
			matchExpressions = append(matchExpressions,
				rabbitmqv1beta1.RabbitmqClusterSpecAffinityNodeAffinityRequiredDuringSchedulingIgnoredDuringExecutionNodeSelectorTermsMatchExpressionsArgs{
					Key:      pulumi.String(key),
					Operator: pulumi.String("In"),
					Values:   pulumi.ToStringArray([]string{value}),
				})
		}
		affinity.NodeAffinity = rabbitmqv1beta1.RabbitmqClusterSpecAffinityNodeAffinityArgs{
			RequiredDuringSchedulingIgnoredDuringExecution: rabbitmqv1beta1.RabbitmqClusterSpecAffinityNodeAffinityRequiredDuringSchedulingIgnoredDuringExecutionArgs{
				NodeSelectorTerms: rabbitmqv1beta1.RabbitmqClusterSpecAffinityNodeAffinityRequiredDuringSchedulingIgnoredDuringExecutionNodeSelectorTermsArray{
					rabbitmqv1beta1.RabbitmqClusterSpecAffinityNodeAffinityRequiredDuringSchedulingIgnoredDuringExecutionNodeSelectorTermsArgs{
						MatchExpressions: matchExpressions,
					},
				},
			},
		}
	}

	if spec.GetSpreadAcrossNodes() {
		affinity.PodAntiAffinity = rabbitmqv1beta1.RabbitmqClusterSpecAffinityPodAntiAffinityArgs{
			RequiredDuringSchedulingIgnoredDuringExecution: rabbitmqv1beta1.RabbitmqClusterSpecAffinityPodAntiAffinityRequiredDuringSchedulingIgnoredDuringExecutionArray{
				rabbitmqv1beta1.RabbitmqClusterSpecAffinityPodAntiAffinityRequiredDuringSchedulingIgnoredDuringExecutionArgs{
					TopologyKey: pulumi.String(vars.NodeAffinityKubernetesHostnameKey),
					LabelSelector: rabbitmqv1beta1.RabbitmqClusterSpecAffinityPodAntiAffinityRequiredDuringSchedulingIgnoredDuringExecutionLabelSelectorArgs{
						MatchLabels: pulumi.StringMap{
							vars.AppLabelKey: pulumi.String(locals.ClusterName),
						},
					},
				},
			},
		}
	}

	return affinity
}

// buildSecretBackend renders the vault XOR external-secret arm; nil when
// the block is absent.
func buildSecretBackend(secretBackend *kubernetesrabbitmqv1alpha1.KubernetesRabbitMqSecretBackend) rabbitmqv1beta1.RabbitmqClusterSpecSecretBackendPtrInput {
	if secretBackend == nil {
		return nil
	}

	if vault := secretBackend.GetVault(); vault != nil {
		vaultArgs := rabbitmqv1beta1.RabbitmqClusterSpecSecretBackendVaultArgs{
			Role:            pulumi.String(vault.GetRole()),
			DefaultUserPath: pulumi.String(vault.GetDefaultUserPath()),
		}
		if len(vault.GetAnnotations()) > 0 {
			vaultArgs.Annotations = pulumi.ToStringMap(vault.GetAnnotations())
		}
		if vault.GetPkiIssuerPath() != "" {
			vaultArgs.Tls = rabbitmqv1beta1.RabbitmqClusterSpecSecretBackendVaultTlsArgs{
				PkiIssuerPath: pulumi.String(vault.GetPkiIssuerPath()),
			}
		}
		return rabbitmqv1beta1.RabbitmqClusterSpecSecretBackendArgs{
			Vault: vaultArgs,
		}
	}

	if name := secretBackend.GetExternalSecretName(); name != "" {
		return rabbitmqv1beta1.RabbitmqClusterSpecSecretBackendArgs{
			ExternalSecret: rabbitmqv1beta1.RabbitmqClusterSpecSecretBackendExternalSecretArgs{
				Name: pulumi.String(name),
			},
		}
	}

	return nil
}
