package workloadpod

import (
	kubernetesv1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// PodTemplateInputs carries everything a workload kind contributes to the pod
// template beyond the shared WorkloadPod proto: the assembled container lists,
// the labels the controller stamps on pods, and kind-specific overrides
// (a Job's restart policy, a batch active deadline).
type PodTemplateInputs struct {
	// Labels stamped on the pod template metadata — the workload's full label
	// set (selector labels included). WorkloadPod.labels are merged on top;
	// selector keys always win so user labels can never break selection.
	Labels map[string]string

	// Containers is the complete pod container list (app first, then sidecars),
	// built via BuildContainers.
	Containers corev1.ContainerArray

	// Volumes is the pod-level volume list derived from every container's
	// volume mounts, built via BuildVolumes.
	Volumes corev1.VolumeArray

	// RestartPolicy overrides the pod restart policy — Jobs and CronJobs pass
	// "Never"/"OnFailure"; long-running kinds leave it empty (Always).
	RestartPolicy string

	// ImagePullSecretNames are resolved secret names from the workload spec's
	// pod.image_pull_secrets (values already resolved from refs by the
	// orchestrator), plus any module-created pull secret.
	ImagePullSecretNames []string
}

// BuildPodTemplateSpec assembles the complete pod template from the shared
// WorkloadPod proto and the kind's own inputs. This is the single place pod
// semantics are rendered for every workload kind — scheduling, security,
// identity, DNS, and termination behavior land here.
func BuildPodTemplateSpec(pod *kubernetesv1.WorkloadPod, in PodTemplateInputs, envSecretName string) *corev1.PodTemplateSpecArgs {
	labels := map[string]string{}
	for k, v := range in.Labels {
		labels[k] = v
	}

	podSpec := &corev1.PodSpecArgs{
		Containers: in.Containers,
	}
	if len(in.Volumes) > 0 {
		podSpec.Volumes = in.Volumes
	}
	if in.RestartPolicy != "" {
		podSpec.RestartPolicy = pulumi.String(in.RestartPolicy)
	}

	templateMeta := &metav1.ObjectMetaArgs{}

	if pod != nil {
		// User pod labels merge under the controller's labels: selector keys must
		// keep the controller-derived values or the workload would orphan its pods.
		for k, v := range pod.Labels {
			if _, exists := labels[k]; !exists {
				labels[k] = v
			}
		}
		if len(pod.Annotations) > 0 {
			templateMeta.Annotations = pulumi.ToStringMap(pod.Annotations)
		}

		if pod.ServiceAccount.GetValue() != "" {
			podSpec.ServiceAccountName = pulumi.String(pod.ServiceAccount.GetValue())
		}
		if pod.AutomountServiceAccountToken != nil {
			podSpec.AutomountServiceAccountToken = pulumi.Bool(*pod.AutomountServiceAccountToken)
		}

		podSpec.InitContainers = BuildInitContainers(pod.InitContainers, envSecretName)

		if pod.Scheduling != nil {
			applyScheduling(podSpec, pod.Scheduling)
		}
		if pod.SecurityContext != nil {
			podSpec.SecurityContext = buildPodSecurityContext(pod.SecurityContext)
		}
		if pod.TerminationGracePeriodSeconds != nil {
			podSpec.TerminationGracePeriodSeconds = pulumi.Int(int(*pod.TerminationGracePeriodSeconds))
		}
		if pod.DnsPolicy != "" {
			podSpec.DnsPolicy = pulumi.String(pod.DnsPolicy)
		}
		if pod.DnsConfig != nil {
			podSpec.DnsConfig = buildDnsConfig(pod.DnsConfig)
		}
		if len(pod.HostAliases) > 0 {
			aliases := make(corev1.HostAliasArray, 0, len(pod.HostAliases))
			for _, ha := range pod.HostAliases {
				aliases = append(aliases, &corev1.HostAliasArgs{
					Ip:        pulumi.String(ha.Ip),
					Hostnames: pulumi.ToStringArray(ha.Hostnames),
				})
			}
			podSpec.HostAliases = aliases
		}
		if pod.HostNetwork {
			podSpec.HostNetwork = pulumi.Bool(true)
		}
		if pod.HostPid {
			podSpec.HostPID = pulumi.Bool(true)
		}
		if pod.PriorityClassName != "" {
			podSpec.PriorityClassName = pulumi.String(pod.PriorityClassName)
		}
		if pod.RuntimeClassName != "" {
			podSpec.RuntimeClassName = pulumi.StringPtr(pod.RuntimeClassName)
		}
	}

	// Pod-level image pull secrets: spec-listed names plus the module-created
	// registry secret. ServiceAccount-attached pull secrets need no entry here.
	if len(in.ImagePullSecretNames) > 0 {
		refs := make(corev1.LocalObjectReferenceArray, 0, len(in.ImagePullSecretNames))
		for _, name := range in.ImagePullSecretNames {
			refs = append(refs, corev1.LocalObjectReferenceArgs{Name: pulumi.String(name)})
		}
		podSpec.ImagePullSecrets = refs
	}

	templateMeta.Labels = pulumi.ToStringMap(labels)

	return &corev1.PodTemplateSpecArgs{
		Metadata: templateMeta,
		Spec:     podSpec,
	}
}

// BuildVolumes derives the pod-level volume list from every container's mounts
// (app, sidecars, and init containers), de-duplicating by volume name — two
// containers mounting the same name share one pod volume, which is exactly how
// containers share an EmptyDir. The first declaration of a name wins; validation
// keeps duplicate names with conflicting sources out of the spec.
func BuildVolumes(containers ...*kubernetesv1.WorkloadContainer) corev1.VolumeArray {
	volumes := make(corev1.VolumeArray, 0)
	seen := map[string]bool{}

	for _, c := range containers {
		if c == nil {
			continue
		}
		for _, vm := range c.VolumeMounts {
			if seen[vm.Name] {
				continue
			}

			volumeArgs := &corev1.VolumeArgs{
				Name: pulumi.String(vm.Name),
			}

			switch {
			case vm.ConfigMap != nil:
				configMapVolumeSource := &corev1.ConfigMapVolumeSourceArgs{
					Name: pulumi.String(vm.ConfigMap.Name),
				}
				if vm.ConfigMap.Key != "" {
					path := vm.ConfigMap.Path
					if path == "" {
						path = vm.ConfigMap.Key
					}
					configMapVolumeSource.Items = corev1.KeyToPathArray{
						&corev1.KeyToPathArgs{
							Key:  pulumi.String(vm.ConfigMap.Key),
							Path: pulumi.String(path),
						},
					}
				}
				if vm.ConfigMap.DefaultMode > 0 {
					configMapVolumeSource.DefaultMode = pulumi.Int(int(vm.ConfigMap.DefaultMode))
				}
				volumeArgs.ConfigMap = configMapVolumeSource

			case vm.Secret != nil:
				secretVolumeSource := &corev1.SecretVolumeSourceArgs{
					SecretName: pulumi.String(vm.Secret.Name),
				}
				if vm.Secret.Key != "" {
					path := vm.Secret.Path
					if path == "" {
						path = vm.Secret.Key
					}
					secretVolumeSource.Items = corev1.KeyToPathArray{
						&corev1.KeyToPathArgs{
							Key:  pulumi.String(vm.Secret.Key),
							Path: pulumi.String(path),
						},
					}
				}
				if vm.Secret.DefaultMode > 0 {
					secretVolumeSource.DefaultMode = pulumi.Int(int(vm.Secret.DefaultMode))
				}
				volumeArgs.Secret = secretVolumeSource

			case vm.HostPath != nil:
				hostPathVolumeSource := &corev1.HostPathVolumeSourceArgs{
					Path: pulumi.String(vm.HostPath.Path),
				}
				if vm.HostPath.Type != "" {
					hostPathVolumeSource.Type = pulumi.StringPtr(vm.HostPath.Type)
				}
				volumeArgs.HostPath = hostPathVolumeSource

			case vm.EmptyDir != nil:
				emptyDirVolumeSource := &corev1.EmptyDirVolumeSourceArgs{}
				if vm.EmptyDir.Medium != "" {
					emptyDirVolumeSource.Medium = pulumi.String(vm.EmptyDir.Medium)
				}
				if vm.EmptyDir.SizeLimit != "" {
					emptyDirVolumeSource.SizeLimit = pulumi.String(vm.EmptyDir.SizeLimit)
				}
				volumeArgs.EmptyDir = emptyDirVolumeSource

			case vm.Pvc != nil:
				volumeArgs.PersistentVolumeClaim = &corev1.PersistentVolumeClaimVolumeSourceArgs{
					ClaimName: pulumi.String(vm.Pvc.ClaimName),
					ReadOnly:  pulumi.Bool(vm.Pvc.ReadOnly),
				}

			default:
				// A mount without a source references a volume defined elsewhere —
				// for StatefulSets that is a volumeClaimTemplate name, which the
				// StatefulSet controller binds per pod without a pod-level volume.
				continue
			}

			seen[vm.Name] = true
			volumes = append(volumes, volumeArgs)
		}
	}

	return volumes
}

func applyScheduling(podSpec *corev1.PodSpecArgs, s *kubernetesv1.WorkloadScheduling) {
	if len(s.NodeSelector) > 0 {
		podSpec.NodeSelector = pulumi.ToStringMap(s.NodeSelector)
	}

	if len(s.Tolerations) > 0 {
		tolerations := make(corev1.TolerationArray, 0, len(s.Tolerations))
		for _, t := range s.Tolerations {
			tolArgs := &corev1.TolerationArgs{}
			if t.Key != "" {
				tolArgs.Key = pulumi.String(t.Key)
			}
			if t.Operator != "" {
				tolArgs.Operator = pulumi.String(t.Operator)
			}
			if t.Value != "" {
				tolArgs.Value = pulumi.String(t.Value)
			}
			if t.Effect != "" {
				tolArgs.Effect = pulumi.String(t.Effect)
			}
			if t.TolerationSeconds != nil {
				tolArgs.TolerationSeconds = pulumi.Int(int(*t.TolerationSeconds))
			}
			tolerations = append(tolerations, tolArgs)
		}
		podSpec.Tolerations = tolerations
	}

	affinity := &corev1.AffinityArgs{}
	hasAffinity := false

	if s.NodeAffinity != nil {
		nodeAffinity := &corev1.NodeAffinityArgs{}
		if len(s.NodeAffinity.Required) > 0 {
			terms := make(corev1.NodeSelectorTermArray, 0, len(s.NodeAffinity.Required))
			for _, t := range s.NodeAffinity.Required {
				terms = append(terms, buildNodeSelectorTerm(t))
			}
			nodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution = &corev1.NodeSelectorArgs{
				NodeSelectorTerms: terms,
			}
		}
		if len(s.NodeAffinity.Preferred) > 0 {
			preferred := make(corev1.PreferredSchedulingTermArray, 0, len(s.NodeAffinity.Preferred))
			for _, p := range s.NodeAffinity.Preferred {
				preferred = append(preferred, &corev1.PreferredSchedulingTermArgs{
					Weight:     pulumi.Int(p.Weight),
					Preference: buildNodeSelectorTerm(p.Term),
				})
			}
			nodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution = preferred
		}
		affinity.NodeAffinity = nodeAffinity
		hasAffinity = true
	}

	if s.PodAffinity != nil {
		affinity.PodAffinity = &corev1.PodAffinityArgs{
			RequiredDuringSchedulingIgnoredDuringExecution:  buildPodAffinityTerms(s.PodAffinity.Required),
			PreferredDuringSchedulingIgnoredDuringExecution: buildWeightedPodAffinityTerms(s.PodAffinity.Preferred),
		}
		hasAffinity = true
	}

	if s.PodAntiAffinity != nil {
		affinity.PodAntiAffinity = &corev1.PodAntiAffinityArgs{
			RequiredDuringSchedulingIgnoredDuringExecution:  buildPodAffinityTerms(s.PodAntiAffinity.Required),
			PreferredDuringSchedulingIgnoredDuringExecution: buildWeightedPodAffinityTerms(s.PodAntiAffinity.Preferred),
		}
		hasAffinity = true
	}

	if hasAffinity {
		podSpec.Affinity = affinity
	}

	if len(s.TopologySpreadConstraints) > 0 {
		constraints := make(corev1.TopologySpreadConstraintArray, 0, len(s.TopologySpreadConstraints))
		for _, c := range s.TopologySpreadConstraints {
			constraintArgs := &corev1.TopologySpreadConstraintArgs{
				MaxSkew:           pulumi.Int(c.MaxSkew),
				TopologyKey:       pulumi.String(c.TopologyKey),
				WhenUnsatisfiable: pulumi.String(c.WhenUnsatisfiable),
			}
			// Empty match_labels defaults to the workload's own selector labels at the
			// kind module level (self-spreading); the module substitutes them before
			// calling here, so an empty map at this point is emitted as-is.
			if len(c.MatchLabels) > 0 {
				constraintArgs.LabelSelector = &metav1.LabelSelectorArgs{
					MatchLabels: pulumi.ToStringMap(c.MatchLabels),
				}
			}
			constraints = append(constraints, constraintArgs)
		}
		podSpec.TopologySpreadConstraints = constraints
	}

	if s.SchedulerName != "" {
		podSpec.SchedulerName = pulumi.String(s.SchedulerName)
	}
}

func buildNodeSelectorTerm(t *kubernetesv1.WorkloadNodeSelectorTerm) corev1.NodeSelectorTermArgs {
	expressions := make(corev1.NodeSelectorRequirementArray, 0, len(t.MatchExpressions))
	for _, e := range t.MatchExpressions {
		reqArgs := &corev1.NodeSelectorRequirementArgs{
			Key:      pulumi.String(e.Key),
			Operator: pulumi.String(e.Operator),
		}
		if len(e.Values) > 0 {
			reqArgs.Values = pulumi.ToStringArray(e.Values)
		}
		expressions = append(expressions, reqArgs)
	}
	return corev1.NodeSelectorTermArgs{
		MatchExpressions: expressions,
	}
}

func buildPodAffinityTerms(terms []*kubernetesv1.WorkloadPodAffinityTerm) corev1.PodAffinityTermArray {
	if len(terms) == 0 {
		return nil
	}
	result := make(corev1.PodAffinityTermArray, 0, len(terms))
	for _, t := range terms {
		result = append(result, buildPodAffinityTerm(t))
	}
	return result
}

func buildWeightedPodAffinityTerms(terms []*kubernetesv1.WorkloadWeightedPodAffinityTerm) corev1.WeightedPodAffinityTermArray {
	if len(terms) == 0 {
		return nil
	}
	result := make(corev1.WeightedPodAffinityTermArray, 0, len(terms))
	for _, t := range terms {
		result = append(result, &corev1.WeightedPodAffinityTermArgs{
			Weight:          pulumi.Int(t.Weight),
			PodAffinityTerm: buildPodAffinityTerm(t.Term),
		})
	}
	return result
}

func buildPodAffinityTerm(t *kubernetesv1.WorkloadPodAffinityTerm) corev1.PodAffinityTermArgs {
	args := corev1.PodAffinityTermArgs{
		TopologyKey: pulumi.String(t.TopologyKey),
		LabelSelector: &metav1.LabelSelectorArgs{
			MatchLabels: pulumi.ToStringMap(t.MatchLabels),
		},
	}
	if len(t.Namespaces) > 0 {
		args.Namespaces = pulumi.ToStringArray(t.Namespaces)
	}
	return args
}

func buildPodSecurityContext(sc *kubernetesv1.WorkloadPodSecurityContext) *corev1.PodSecurityContextArgs {
	args := &corev1.PodSecurityContextArgs{}
	if sc.RunAsUser != nil {
		args.RunAsUser = pulumi.Int(int(*sc.RunAsUser))
	}
	if sc.RunAsGroup != nil {
		args.RunAsGroup = pulumi.Int(int(*sc.RunAsGroup))
	}
	if sc.RunAsNonRoot != nil {
		args.RunAsNonRoot = pulumi.Bool(*sc.RunAsNonRoot)
	}
	if sc.FsGroup != nil {
		args.FsGroup = pulumi.Int(int(*sc.FsGroup))
	}
	if sc.FsGroupChangePolicy != "" {
		args.FsGroupChangePolicy = pulumi.StringPtr(sc.FsGroupChangePolicy)
	}
	if len(sc.SupplementalGroups) > 0 {
		groups := make(pulumi.IntArray, 0, len(sc.SupplementalGroups))
		for _, g := range sc.SupplementalGroups {
			groups = append(groups, pulumi.Int(int(g)))
		}
		args.SupplementalGroups = groups
	}
	if len(sc.Sysctls) > 0 {
		sysctls := make(corev1.SysctlArray, 0, len(sc.Sysctls))
		for _, s := range sc.Sysctls {
			sysctls = append(sysctls, &corev1.SysctlArgs{
				Name:  pulumi.String(s.Name),
				Value: pulumi.String(s.Value),
			})
		}
		args.Sysctls = sysctls
	}
	if sc.SeccompProfile != nil {
		args.SeccompProfile = &corev1.SeccompProfileArgs{
			Type: pulumi.String(sc.SeccompProfile.Type),
		}
		if sc.SeccompProfile.LocalhostProfile != "" {
			args.SeccompProfile = &corev1.SeccompProfileArgs{
				Type:             pulumi.String(sc.SeccompProfile.Type),
				LocalhostProfile: pulumi.StringPtr(sc.SeccompProfile.LocalhostProfile),
			}
		}
	}
	return args
}

func buildDnsConfig(dc *kubernetesv1.WorkloadPodDnsConfig) *corev1.PodDNSConfigArgs {
	args := &corev1.PodDNSConfigArgs{}
	if len(dc.Nameservers) > 0 {
		args.Nameservers = pulumi.ToStringArray(dc.Nameservers)
	}
	if len(dc.Searches) > 0 {
		args.Searches = pulumi.ToStringArray(dc.Searches)
	}
	if len(dc.Options) > 0 {
		options := make(corev1.PodDNSConfigOptionArray, 0, len(dc.Options))
		for _, o := range dc.Options {
			optArgs := &corev1.PodDNSConfigOptionArgs{
				Name: pulumi.StringPtr(o.Name),
			}
			if o.Value != "" {
				optArgs.Value = pulumi.StringPtr(o.Value)
			}
			options = append(options, optArgs)
		}
		args.Options = options
	}
	return args
}

// ResolveImagePullSecretNames flattens the pod's image_pull_secrets references
// into plain names (the orchestrator resolves refs to literals before IaC runs),
// appending any module-created secret name (e.g. the docker-config secret a kind
// materializes from its stack input).
func ResolveImagePullSecretNames(pod *kubernetesv1.WorkloadPod, moduleCreated ...string) []string {
	names := make([]string, 0)
	if pod != nil {
		for _, ref := range pod.ImagePullSecrets {
			if v := ref.GetValue(); v != "" {
				names = append(names, v)
			}
		}
	}
	for _, n := range moduleCreated {
		if n != "" {
			names = append(names, n)
		}
	}
	return names
}
