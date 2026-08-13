package module

import (
	"sort"

	awsbatchjobdefinitionv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsbatchjobdefinition/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/batch"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// buildEksProperties renders the Batch-on-EKS pod arm as the provider's
// typed eksProperties block. Field-for-field send parity with the
// Terraform module: presence-typed spec fields (host_network, security
// context run_as_user/run_as_group/allow_privilege_escalation) are sent
// only when set so AWS's own defaults apply; plain bools are always sent
// (state-pinned) when their enclosing block renders.
//
// The provider duplicates the container shape into distinct Container and
// InitContainer types, hence the twin builders below -- keep them
// member-for-member identical.
func buildEksProperties(eks *awsbatchjobdefinitionv1alpha1.AwsBatchJobDefinitionEks) *batch.JobDefinitionEksPropertiesArgs {
	pod := &batch.JobDefinitionEksPropertiesPodPropertiesArgs{
		Containers:            buildEksContainers(eks.Containers),
		ShareProcessNamespace: pulumi.BoolPtr(eks.ShareProcessNamespace),
	}

	// Unset means AWS's default, which is TRUE for Batch pods.
	if eks.HostNetwork != nil {
		pod.HostNetwork = pulumi.BoolPtr(eks.GetHostNetwork())
	}
	if eks.DnsPolicy != "" {
		pod.DnsPolicy = pulumi.StringPtr(eks.DnsPolicy)
	}
	if eks.ServiceAccountName != "" {
		pod.ServiceAccountName = pulumi.StringPtr(eks.ServiceAccountName)
	}
	if len(eks.PodLabels) > 0 {
		pod.Metadata = &batch.JobDefinitionEksPropertiesPodPropertiesMetadataArgs{
			Labels: pulumi.ToStringMap(eks.PodLabels),
		}
	}
	if len(eks.InitContainers) > 0 {
		pod.InitContainers = buildEksInitContainers(eks.InitContainers)
	}
	if len(eks.ImagePullSecretNames) > 0 {
		var secrets batch.JobDefinitionEksPropertiesPodPropertiesImagePullSecretArray
		for _, name := range eks.ImagePullSecretNames {
			secrets = append(secrets, &batch.JobDefinitionEksPropertiesPodPropertiesImagePullSecretArgs{
				Name: pulumi.String(name),
			})
		}
		pod.ImagePullSecrets = secrets
	}
	if len(eks.Volumes) > 0 {
		var volumes batch.JobDefinitionEksPropertiesPodPropertiesVolumeArray
		for _, volume := range eks.Volumes {
			entry := &batch.JobDefinitionEksPropertiesPodPropertiesVolumeArgs{
				Name: pulumi.StringPtr(volume.Name),
			}
			if volume.EmptyDir != nil {
				emptyDir := &batch.JobDefinitionEksPropertiesPodPropertiesVolumeEmptyDirArgs{
					SizeLimit: pulumi.String(volume.EmptyDir.SizeLimit),
				}
				// Unset medium means node-backed storage (AWS default "").
				if volume.EmptyDir.Medium != "" {
					emptyDir.Medium = pulumi.StringPtr(volume.EmptyDir.Medium)
				}
				entry.EmptyDir = emptyDir
			}
			if volume.HostPath != "" {
				entry.HostPath = &batch.JobDefinitionEksPropertiesPodPropertiesVolumeHostPathArgs{
					Path: pulumi.String(volume.HostPath),
				}
			}
			if volume.Secret != nil {
				entry.Secret = &batch.JobDefinitionEksPropertiesPodPropertiesVolumeSecretArgs{
					SecretName: pulumi.String(volume.Secret.SecretName),
					Optional:   pulumi.BoolPtr(volume.Secret.Optional),
				}
			}
			volumes = append(volumes, entry)
		}
		pod.Volumes = volumes
	}

	return &batch.JobDefinitionEksPropertiesArgs{PodProperties: pod}
}

// sortedEnvNames returns the env map's keys in sorted order so both
// engines render the same deterministic list.
func sortedEnvNames(env map[string]string) []string {
	names := make([]string, 0, len(env))
	for name := range env {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func buildEksContainers(containers []*awsbatchjobdefinitionv1alpha1.AwsBatchJobDefinitionEksContainer) batch.JobDefinitionEksPropertiesPodPropertiesContainerArray {
	var result batch.JobDefinitionEksPropertiesPodPropertiesContainerArray
	for _, container := range containers {
		entry := &batch.JobDefinitionEksPropertiesPodPropertiesContainerArgs{
			Image: pulumi.String(container.Image),
		}
		if container.Name != "" {
			entry.Name = pulumi.StringPtr(container.Name)
		}
		if len(container.Command) > 0 {
			entry.Commands = pulumi.ToStringArray(container.Command)
		}
		if len(container.Args) > 0 {
			entry.Args = pulumi.ToStringArray(container.Args)
		}
		if container.ImagePullPolicy != "" {
			entry.ImagePullPolicy = pulumi.StringPtr(container.ImagePullPolicy)
		}
		if len(container.Env) > 0 {
			var envs batch.JobDefinitionEksPropertiesPodPropertiesContainerEnvArray
			for _, name := range sortedEnvNames(container.Env) {
				envs = append(envs, &batch.JobDefinitionEksPropertiesPodPropertiesContainerEnvArgs{
					Name:  pulumi.String(name),
					Value: pulumi.String(container.Env[name]),
				})
			}
			entry.Envs = envs
		}
		if container.Resources != nil {
			resources := &batch.JobDefinitionEksPropertiesPodPropertiesContainerResourcesArgs{}
			if len(container.Resources.Limits) > 0 {
				resources.Limits = pulumi.ToStringMap(container.Resources.Limits)
			}
			if len(container.Resources.Requests) > 0 {
				resources.Requests = pulumi.ToStringMap(container.Resources.Requests)
			}
			entry.Resources = resources
		}
		if sc := container.SecurityContext; sc != nil {
			securityContext := &batch.JobDefinitionEksPropertiesPodPropertiesContainerSecurityContextArgs{
				RunAsNonRoot:           pulumi.BoolPtr(sc.RunAsNonRoot),
				Privileged:             pulumi.BoolPtr(sc.Privileged),
				ReadOnlyRootFileSystem: pulumi.BoolPtr(sc.ReadOnlyRootFileSystem),
			}
			// 0 (root) is a legal explicit UID/GID -- presence decides.
			if sc.RunAsUser != nil {
				securityContext.RunAsUser = pulumi.IntPtr(int(sc.GetRunAsUser()))
			}
			if sc.RunAsGroup != nil {
				securityContext.RunAsGroup = pulumi.IntPtr(int(sc.GetRunAsGroup()))
			}
			if sc.AllowPrivilegeEscalation != nil {
				securityContext.AllowPrivilegeEscalation = pulumi.BoolPtr(sc.GetAllowPrivilegeEscalation())
			}
			entry.SecurityContext = securityContext
		}
		if len(container.VolumeMounts) > 0 {
			var mounts batch.JobDefinitionEksPropertiesPodPropertiesContainerVolumeMountArray
			for _, mount := range container.VolumeMounts {
				mounts = append(mounts, &batch.JobDefinitionEksPropertiesPodPropertiesContainerVolumeMountArgs{
					Name:      pulumi.String(mount.Name),
					MountPath: pulumi.String(mount.MountPath),
					ReadOnly:  pulumi.BoolPtr(mount.ReadOnly),
				})
			}
			entry.VolumeMounts = mounts
		}
		result = append(result, entry)
	}
	return result
}

func buildEksInitContainers(containers []*awsbatchjobdefinitionv1alpha1.AwsBatchJobDefinitionEksContainer) batch.JobDefinitionEksPropertiesPodPropertiesInitContainerArray {
	var result batch.JobDefinitionEksPropertiesPodPropertiesInitContainerArray
	for _, container := range containers {
		entry := &batch.JobDefinitionEksPropertiesPodPropertiesInitContainerArgs{
			Image: pulumi.String(container.Image),
		}
		if container.Name != "" {
			entry.Name = pulumi.StringPtr(container.Name)
		}
		if len(container.Command) > 0 {
			entry.Commands = pulumi.ToStringArray(container.Command)
		}
		if len(container.Args) > 0 {
			entry.Args = pulumi.ToStringArray(container.Args)
		}
		if container.ImagePullPolicy != "" {
			entry.ImagePullPolicy = pulumi.StringPtr(container.ImagePullPolicy)
		}
		if len(container.Env) > 0 {
			var envs batch.JobDefinitionEksPropertiesPodPropertiesInitContainerEnvArray
			for _, name := range sortedEnvNames(container.Env) {
				envs = append(envs, &batch.JobDefinitionEksPropertiesPodPropertiesInitContainerEnvArgs{
					Name:  pulumi.String(name),
					Value: pulumi.String(container.Env[name]),
				})
			}
			entry.Envs = envs
		}
		if container.Resources != nil {
			resources := &batch.JobDefinitionEksPropertiesPodPropertiesInitContainerResourcesArgs{}
			if len(container.Resources.Limits) > 0 {
				resources.Limits = pulumi.ToStringMap(container.Resources.Limits)
			}
			if len(container.Resources.Requests) > 0 {
				resources.Requests = pulumi.ToStringMap(container.Resources.Requests)
			}
			entry.Resources = resources
		}
		if sc := container.SecurityContext; sc != nil {
			securityContext := &batch.JobDefinitionEksPropertiesPodPropertiesInitContainerSecurityContextArgs{
				RunAsNonRoot:           pulumi.BoolPtr(sc.RunAsNonRoot),
				Privileged:             pulumi.BoolPtr(sc.Privileged),
				ReadOnlyRootFileSystem: pulumi.BoolPtr(sc.ReadOnlyRootFileSystem),
			}
			if sc.RunAsUser != nil {
				securityContext.RunAsUser = pulumi.IntPtr(int(sc.GetRunAsUser()))
			}
			if sc.RunAsGroup != nil {
				securityContext.RunAsGroup = pulumi.IntPtr(int(sc.GetRunAsGroup()))
			}
			if sc.AllowPrivilegeEscalation != nil {
				securityContext.AllowPrivilegeEscalation = pulumi.BoolPtr(sc.GetAllowPrivilegeEscalation())
			}
			entry.SecurityContext = securityContext
		}
		if len(container.VolumeMounts) > 0 {
			var mounts batch.JobDefinitionEksPropertiesPodPropertiesInitContainerVolumeMountArray
			for _, mount := range container.VolumeMounts {
				mounts = append(mounts, &batch.JobDefinitionEksPropertiesPodPropertiesInitContainerVolumeMountArgs{
					Name:      pulumi.String(mount.Name),
					MountPath: pulumi.String(mount.MountPath),
					ReadOnly:  pulumi.BoolPtr(mount.ReadOnly),
				})
			}
			entry.VolumeMounts = mounts
		}
		result = append(result, entry)
	}
	return result
}
