package module

import (
	awssagemakerdomainv1alpha1 "github.com/plantonhq/planton/catalog/aws/awssagemakerdomain/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/sagemaker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// buildDefaultSpaceSettings assembles the shared-space baseline. The spec
// reuses the same app-settings messages as the user plane, but the SDK models
// every nesting path as its own Go type, so the space-plane conversions live
// here as their own builders.
func buildDefaultSpaceSettings(dss *awssagemakerdomainv1alpha1.AwsSagemakerDomainDefaultSpaceSettings) *sagemaker.DomainDefaultSpaceSettingsArgs {
	settings := &sagemaker.DomainDefaultSpaceSettingsArgs{
		ExecutionRole: pulumi.String(dss.ExecutionRoleArn.GetValue()),
	}

	if len(dss.SecurityGroupIds) > 0 {
		var sgs pulumi.StringArray
		for _, sg := range dss.SecurityGroupIds {
			sgs = append(sgs, pulumi.String(sg.GetValue()))
		}
		settings.SecurityGroups = sgs
	}

	if jl := dss.JupyterLabAppSettings; jl != nil {
		settings.JupyterLabAppSettings = buildSpaceJupyterLabAppSettings(jl)
	}

	if js := dss.JupyterServerAppSettings; js != nil {
		jsArgs := &sagemaker.DomainDefaultSpaceSettingsJupyterServerAppSettingsArgs{}
		if rs := js.DefaultResourceSpec; rs != nil {
			specArgs := &sagemaker.DomainDefaultSpaceSettingsJupyterServerAppSettingsDefaultResourceSpecArgs{}
			if rs.InstanceType != "" {
				specArgs.InstanceType = pulumi.StringPtr(rs.InstanceType)
			}
			if rs.LifecycleConfigArn != "" {
				specArgs.LifecycleConfigArn = pulumi.StringPtr(rs.LifecycleConfigArn)
			}
			if rs.SagemakerImageArn != "" {
				specArgs.SagemakerImageArn = pulumi.StringPtr(rs.SagemakerImageArn)
			}
			if rs.SagemakerImageVersionAlias != "" {
				specArgs.SagemakerImageVersionAlias = pulumi.StringPtr(rs.SagemakerImageVersionAlias)
			}
			if rs.SagemakerImageVersionArn != "" {
				specArgs.SagemakerImageVersionArn = pulumi.StringPtr(rs.SagemakerImageVersionArn)
			}
			jsArgs.DefaultResourceSpec = specArgs
		}
		if len(js.LifecycleConfigArns) > 0 {
			jsArgs.LifecycleConfigArns = pulumi.ToStringArray(js.LifecycleConfigArns)
		}
		if len(js.CodeRepositories) > 0 {
			var repos sagemaker.DomainDefaultSpaceSettingsJupyterServerAppSettingsCodeRepositoryArray
			for _, repo := range js.CodeRepositories {
				repos = append(repos, &sagemaker.DomainDefaultSpaceSettingsJupyterServerAppSettingsCodeRepositoryArgs{
					RepositoryUrl: pulumi.String(repo.RepositoryUrl),
				})
			}
			jsArgs.CodeRepositories = repos
		}
		settings.JupyterServerAppSettings = jsArgs
	}

	if kg := dss.KernelGatewayAppSettings; kg != nil {
		kgArgs := &sagemaker.DomainDefaultSpaceSettingsKernelGatewayAppSettingsArgs{}
		if rs := kg.DefaultResourceSpec; rs != nil {
			specArgs := &sagemaker.DomainDefaultSpaceSettingsKernelGatewayAppSettingsDefaultResourceSpecArgs{}
			if rs.InstanceType != "" {
				specArgs.InstanceType = pulumi.StringPtr(rs.InstanceType)
			}
			if rs.LifecycleConfigArn != "" {
				specArgs.LifecycleConfigArn = pulumi.StringPtr(rs.LifecycleConfigArn)
			}
			if rs.SagemakerImageArn != "" {
				specArgs.SagemakerImageArn = pulumi.StringPtr(rs.SagemakerImageArn)
			}
			if rs.SagemakerImageVersionAlias != "" {
				specArgs.SagemakerImageVersionAlias = pulumi.StringPtr(rs.SagemakerImageVersionAlias)
			}
			if rs.SagemakerImageVersionArn != "" {
				specArgs.SagemakerImageVersionArn = pulumi.StringPtr(rs.SagemakerImageVersionArn)
			}
			kgArgs.DefaultResourceSpec = specArgs
		}
		if len(kg.LifecycleConfigArns) > 0 {
			kgArgs.LifecycleConfigArns = pulumi.ToStringArray(kg.LifecycleConfigArns)
		}
		if len(kg.CustomImages) > 0 {
			var images sagemaker.DomainDefaultSpaceSettingsKernelGatewayAppSettingsCustomImageArray
			for _, img := range kg.CustomImages {
				imgArgs := &sagemaker.DomainDefaultSpaceSettingsKernelGatewayAppSettingsCustomImageArgs{
					AppImageConfigName: pulumi.String(img.AppImageConfigName),
					ImageName:          pulumi.String(img.ImageName),
				}
				if img.ImageVersionNumber != nil {
					imgArgs.ImageVersionNumber = pulumi.IntPtr(int(img.GetImageVersionNumber()))
				}
				images = append(images, imgArgs)
			}
			kgArgs.CustomImages = images
		}
		settings.KernelGatewayAppSettings = kgArgs
	}

	// The spec flattens the SDK's single-purpose default_ebs_storage_settings
	// wrapper; it is reconstructed here.
	if sss := dss.SpaceStorageSettings; sss != nil {
		settings.SpaceStorageSettings = &sagemaker.DomainDefaultSpaceSettingsSpaceStorageSettingsArgs{
			DefaultEbsStorageSettings: &sagemaker.DomainDefaultSpaceSettingsSpaceStorageSettingsDefaultEbsStorageSettingsArgs{
				DefaultEbsVolumeSizeInGb: pulumi.Int(int(sss.DefaultEbsVolumeSizeInGb)),
				MaximumEbsVolumeSizeInGb: pulumi.Int(int(sss.MaximumEbsVolumeSizeInGb)),
			},
		}
	}

	if len(dss.CustomFileSystemConfigs) > 0 {
		var configs sagemaker.DomainDefaultSpaceSettingsCustomFileSystemConfigArray
		for _, cfs := range dss.CustomFileSystemConfigs {
			efs := cfs.EfsFileSystemConfig
			configs = append(configs, &sagemaker.DomainDefaultSpaceSettingsCustomFileSystemConfigArgs{
				EfsFileSystemConfig: &sagemaker.DomainDefaultSpaceSettingsCustomFileSystemConfigEfsFileSystemConfigArgs{
					FileSystemId:   pulumi.String(efs.FileSystemId.GetValue()),
					FileSystemPath: pulumi.String(efs.FileSystemPath),
				},
			})
		}
		settings.CustomFileSystemConfigs = configs
	}

	if posix := dss.CustomPosixUserConfig; posix != nil {
		settings.CustomPosixUserConfig = &sagemaker.DomainDefaultSpaceSettingsCustomPosixUserConfigArgs{
			Uid: pulumi.Int(int(posix.Uid)),
			Gid: pulumi.Int(int(posix.Gid)),
		}
	}

	return settings
}

func buildSpaceJupyterLabAppSettings(jl *awssagemakerdomainv1alpha1.AwsSagemakerDomainJupyterLabAppSettings) *sagemaker.DomainDefaultSpaceSettingsJupyterLabAppSettingsArgs {
	jlArgs := &sagemaker.DomainDefaultSpaceSettingsJupyterLabAppSettingsArgs{}

	if rs := jl.DefaultResourceSpec; rs != nil {
		specArgs := &sagemaker.DomainDefaultSpaceSettingsJupyterLabAppSettingsDefaultResourceSpecArgs{}
		if rs.InstanceType != "" {
			specArgs.InstanceType = pulumi.StringPtr(rs.InstanceType)
		}
		if rs.LifecycleConfigArn != "" {
			specArgs.LifecycleConfigArn = pulumi.StringPtr(rs.LifecycleConfigArn)
		}
		if rs.SagemakerImageArn != "" {
			specArgs.SagemakerImageArn = pulumi.StringPtr(rs.SagemakerImageArn)
		}
		if rs.SagemakerImageVersionAlias != "" {
			specArgs.SagemakerImageVersionAlias = pulumi.StringPtr(rs.SagemakerImageVersionAlias)
		}
		if rs.SagemakerImageVersionArn != "" {
			specArgs.SagemakerImageVersionArn = pulumi.StringPtr(rs.SagemakerImageVersionArn)
		}
		jlArgs.DefaultResourceSpec = specArgs
	}

	if len(jl.LifecycleConfigArns) > 0 {
		jlArgs.LifecycleConfigArns = pulumi.ToStringArray(jl.LifecycleConfigArns)
	}

	if jl.BuiltInLifecycleConfigArn != "" {
		jlArgs.BuiltInLifecycleConfigArn = pulumi.StringPtr(jl.BuiltInLifecycleConfigArn)
	}

	if len(jl.CustomImages) > 0 {
		var images sagemaker.DomainDefaultSpaceSettingsJupyterLabAppSettingsCustomImageArray
		for _, img := range jl.CustomImages {
			imgArgs := &sagemaker.DomainDefaultSpaceSettingsJupyterLabAppSettingsCustomImageArgs{
				AppImageConfigName: pulumi.String(img.AppImageConfigName),
				ImageName:          pulumi.String(img.ImageName),
			}
			if img.ImageVersionNumber != nil {
				imgArgs.ImageVersionNumber = pulumi.IntPtr(int(img.GetImageVersionNumber()))
			}
			images = append(images, imgArgs)
		}
		jlArgs.CustomImages = images
	}

	if len(jl.CodeRepositories) > 0 {
		var repos sagemaker.DomainDefaultSpaceSettingsJupyterLabAppSettingsCodeRepositoryArray
		for _, repo := range jl.CodeRepositories {
			repos = append(repos, &sagemaker.DomainDefaultSpaceSettingsJupyterLabAppSettingsCodeRepositoryArgs{
				RepositoryUrl: pulumi.String(repo.RepositoryUrl),
			})
		}
		jlArgs.CodeRepositories = repos
	}

	if idle := jl.IdleSettings; idle != nil {
		jlArgs.AppLifecycleManagement = &sagemaker.DomainDefaultSpaceSettingsJupyterLabAppSettingsAppLifecycleManagementArgs{
			IdleSettings: &sagemaker.DomainDefaultSpaceSettingsJupyterLabAppSettingsAppLifecycleManagementIdleSettingsArgs{
				LifecycleManagement:     pulumi.String(lifecycleManagement(idle)),
				IdleTimeoutInMinutes:    pulumi.IntPtr(int(idle.IdleTimeoutInMinutes)),
				MinIdleTimeoutInMinutes: pulumi.IntPtr(int(idle.MinIdleTimeoutInMinutes)),
				MaxIdleTimeoutInMinutes: pulumi.IntPtr(int(idle.MaxIdleTimeoutInMinutes)),
			},
		}
	}

	if emr := jl.EmrSettings; emr != nil {
		emrArgs := &sagemaker.DomainDefaultSpaceSettingsJupyterLabAppSettingsEmrSettingsArgs{}
		if len(emr.AssumableRoleArns) > 0 {
			var arns pulumi.StringArray
			for _, arn := range emr.AssumableRoleArns {
				arns = append(arns, pulumi.String(arn.GetValue()))
			}
			emrArgs.AssumableRoleArns = arns
		}
		if len(emr.ExecutionRoleArns) > 0 {
			var arns pulumi.StringArray
			for _, arn := range emr.ExecutionRoleArns {
				arns = append(arns, pulumi.String(arn.GetValue()))
			}
			emrArgs.ExecutionRoleArns = arns
		}
		jlArgs.EmrSettings = emrArgs
	}

	return jlArgs
}
