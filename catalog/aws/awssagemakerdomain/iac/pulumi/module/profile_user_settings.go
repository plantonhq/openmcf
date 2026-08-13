package module

import (
	awssagemakerdomainv1alpha1 "github.com/plantonhq/planton/catalog/aws/awssagemakerdomain/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/sagemaker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// The per-profile builders below mirror user_settings.go exactly: the spec
// shares ONE AwsSagemakerDomainUserSettings message between the domain's
// default_user_settings and each profile's user_settings (as AWS shares one
// UserSettings API type), but the Pulumi SDK generates resource-specific arg
// types, so the rendering is duplicated with the UserProfileUserSettings*
// types. Keep in sync with user_settings.go when the tree grows.

func buildProfileUserSettings(dus *awssagemakerdomainv1alpha1.AwsSagemakerDomainUserSettings) *sagemaker.UserProfileUserSettingsArgs {
	settings := &sagemaker.UserProfileUserSettingsArgs{
		ExecutionRole: pulumi.String(dus.ExecutionRoleArn.GetValue()),
	}

	if len(dus.SecurityGroupIds) > 0 {
		var sgs pulumi.StringArray
		for _, sg := range dus.SecurityGroupIds {
			sgs = append(sgs, pulumi.String(sg.GetValue()))
		}
		settings.SecurityGroups = sgs
	}

	if dus.DefaultLandingUri != "" {
		settings.DefaultLandingUri = pulumi.StringPtr(dus.DefaultLandingUri)
	}

	if dus.StudioWebPortal != nil {
		settings.StudioWebPortal = pulumi.StringPtr(dus.GetStudioWebPortal())
	}

	if dus.AutoMountHomeEfs != nil {
		settings.AutoMountHomeEfs = pulumi.StringPtr(dus.GetAutoMountHomeEfs())
	}

	if jl := dus.JupyterLabAppSettings; jl != nil {
		settings.JupyterLabAppSettings = buildProfileJupyterLabAppSettings(jl)
	}

	if js := dus.JupyterServerAppSettings; js != nil {
		settings.JupyterServerAppSettings = buildProfileJupyterServerAppSettings(js)
	}

	if kg := dus.KernelGatewayAppSettings; kg != nil {
		settings.KernelGatewayAppSettings = buildProfileKernelGatewayAppSettings(kg)
	}

	if ce := dus.CodeEditorAppSettings; ce != nil {
		settings.CodeEditorAppSettings = buildProfileCodeEditorAppSettings(ce)
	}

	if tb := dus.TensorBoardAppSettings; tb != nil {
		tbArgs := &sagemaker.UserProfileUserSettingsTensorBoardAppSettingsArgs{}
		if rs := tb.DefaultResourceSpec; rs != nil {
			specArgs := &sagemaker.UserProfileUserSettingsTensorBoardAppSettingsDefaultResourceSpecArgs{}
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
			tbArgs.DefaultResourceSpec = specArgs
		}
		settings.TensorBoardAppSettings = tbArgs
	}

	if rsess := dus.RSessionAppSettings; rsess != nil {
		rsessArgs := &sagemaker.UserProfileUserSettingsRSessionAppSettingsArgs{}
		if rs := rsess.DefaultResourceSpec; rs != nil {
			specArgs := &sagemaker.UserProfileUserSettingsRSessionAppSettingsDefaultResourceSpecArgs{}
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
			rsessArgs.DefaultResourceSpec = specArgs
		}
		if len(rsess.CustomImages) > 0 {
			var images sagemaker.UserProfileUserSettingsRSessionAppSettingsCustomImageArray
			for _, img := range rsess.CustomImages {
				imgArgs := &sagemaker.UserProfileUserSettingsRSessionAppSettingsCustomImageArgs{
					AppImageConfigName: pulumi.String(img.AppImageConfigName),
					ImageName:          pulumi.String(img.ImageName),
				}
				if img.ImageVersionNumber != nil {
					imgArgs.ImageVersionNumber = pulumi.IntPtr(int(img.GetImageVersionNumber()))
				}
				images = append(images, imgArgs)
			}
			rsessArgs.CustomImages = images
		}
		settings.RSessionAppSettings = rsessArgs
	}

	if rsp := dus.RStudioServerProAppSettings; rsp != nil {
		rspArgs := &sagemaker.UserProfileUserSettingsRStudioServerProAppSettingsArgs{}
		if rsp.AccessStatus != "" {
			rspArgs.AccessStatus = pulumi.StringPtr(rsp.AccessStatus)
		}
		// Only meaningful when access_status is ENABLED; the spec's CEL
		// rejects the dead combination before it reaches here.
		if rsp.UserGroup != "" {
			rspArgs.UserGroup = pulumi.StringPtr(rsp.UserGroup)
		}
		settings.RStudioServerProAppSettings = rspArgs
	}

	if canvas := dus.CanvasAppSettings; canvas != nil {
		settings.CanvasAppSettings = buildProfileCanvasAppSettings(canvas)
	}

	if ss := dus.SharingSettings; ss != nil {
		ssArgs := &sagemaker.UserProfileUserSettingsSharingSettingsArgs{}
		if ss.NotebookOutputOption != nil {
			ssArgs.NotebookOutputOption = pulumi.StringPtr(ss.GetNotebookOutputOption())
		}
		if ss.S3KmsKeyId.GetValue() != "" {
			ssArgs.S3KmsKeyId = pulumi.StringPtr(ss.S3KmsKeyId.GetValue())
		}
		if ss.S3OutputPath != "" {
			ssArgs.S3OutputPath = pulumi.StringPtr(ss.S3OutputPath)
		}
		settings.SharingSettings = ssArgs
	}

	// The spec flattens the provider's single-purpose
	// default_ebs_storage_settings wrapper; it is reconstructed here.
	if sss := dus.SpaceStorageSettings; sss != nil {
		settings.SpaceStorageSettings = &sagemaker.UserProfileUserSettingsSpaceStorageSettingsArgs{
			DefaultEbsStorageSettings: &sagemaker.UserProfileUserSettingsSpaceStorageSettingsDefaultEbsStorageSettingsArgs{
				DefaultEbsVolumeSizeInGb: pulumi.Int(int(sss.DefaultEbsVolumeSizeInGb)),
				MaximumEbsVolumeSizeInGb: pulumi.Int(int(sss.MaximumEbsVolumeSizeInGb)),
			},
		}
	}

	if len(dus.CustomFileSystemConfigs) > 0 {
		var configs sagemaker.UserProfileUserSettingsCustomFileSystemConfigArray
		for _, cfs := range dus.CustomFileSystemConfigs {
			efs := cfs.EfsFileSystemConfig
			// The bridge shapes this variant as a one-element array field
			// (the Domain variant is a singular block) — same wire result.
			configs = append(configs, &sagemaker.UserProfileUserSettingsCustomFileSystemConfigArgs{
				EfsFileSystemConfigs: sagemaker.UserProfileUserSettingsCustomFileSystemConfigEfsFileSystemConfigArray{
					&sagemaker.UserProfileUserSettingsCustomFileSystemConfigEfsFileSystemConfigArgs{
						FileSystemId:   pulumi.String(efs.FileSystemId.GetValue()),
						FileSystemPath: pulumi.String(efs.FileSystemPath),
					},
				},
			})
		}
		settings.CustomFileSystemConfigs = configs
	}

	if posix := dus.CustomPosixUserConfig; posix != nil {
		settings.CustomPosixUserConfig = &sagemaker.UserProfileUserSettingsCustomPosixUserConfigArgs{
			Uid: pulumi.Int(int(posix.Uid)),
			Gid: pulumi.Int(int(posix.Gid)),
		}
	}

	if wps := dus.StudioWebPortalSettings; wps != nil {
		wpsArgs := &sagemaker.UserProfileUserSettingsStudioWebPortalSettingsArgs{}
		if len(wps.HiddenAppTypes) > 0 {
			wpsArgs.HiddenAppTypes = pulumi.ToStringArray(wps.HiddenAppTypes)
		}
		if len(wps.HiddenInstanceTypes) > 0 {
			wpsArgs.HiddenInstanceTypes = pulumi.ToStringArray(wps.HiddenInstanceTypes)
		}
		if len(wps.HiddenMlTools) > 0 {
			wpsArgs.HiddenMlTools = pulumi.ToStringArray(wps.HiddenMlTools)
		}
		settings.StudioWebPortalSettings = wpsArgs
	}

	return settings
}

func buildProfileJupyterLabAppSettings(jl *awssagemakerdomainv1alpha1.AwsSagemakerDomainJupyterLabAppSettings) *sagemaker.UserProfileUserSettingsJupyterLabAppSettingsArgs {
	jlArgs := &sagemaker.UserProfileUserSettingsJupyterLabAppSettingsArgs{}

	if rs := jl.DefaultResourceSpec; rs != nil {
		specArgs := &sagemaker.UserProfileUserSettingsJupyterLabAppSettingsDefaultResourceSpecArgs{}
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
		var images sagemaker.UserProfileUserSettingsJupyterLabAppSettingsCustomImageArray
		for _, img := range jl.CustomImages {
			imgArgs := &sagemaker.UserProfileUserSettingsJupyterLabAppSettingsCustomImageArgs{
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
		var repos sagemaker.UserProfileUserSettingsJupyterLabAppSettingsCodeRepositoryArray
		for _, repo := range jl.CodeRepositories {
			repos = append(repos, &sagemaker.UserProfileUserSettingsJupyterLabAppSettingsCodeRepositoryArgs{
				RepositoryUrl: pulumi.String(repo.RepositoryUrl),
			})
		}
		jlArgs.CodeRepositories = repos
	}

	// The spec folds idle_settings directly under the app settings and makes
	// block presence the enable switch; the SDK nests them inside a
	// single-purpose app_lifecycle_management wrapper, both reconstructed
	// here. lifecycle_management defaults to ENABLED when the block is
	// present; an explicit DISABLED keeps the timeouts as published
	// guardrails without enforcing auto-shutdown. All three timeouts are required
	// by the live API whenever the block is sent (absent members transmit as
	// 0 and AWS rejects them), so they pass through unconditionally.
	if idle := jl.IdleSettings; idle != nil {
		jlArgs.AppLifecycleManagement = &sagemaker.UserProfileUserSettingsJupyterLabAppSettingsAppLifecycleManagementArgs{
			IdleSettings: &sagemaker.UserProfileUserSettingsJupyterLabAppSettingsAppLifecycleManagementIdleSettingsArgs{
				LifecycleManagement:     pulumi.String(lifecycleManagement(idle)),
				IdleTimeoutInMinutes:    pulumi.IntPtr(int(idle.IdleTimeoutInMinutes)),
				MinIdleTimeoutInMinutes: pulumi.IntPtr(int(idle.MinIdleTimeoutInMinutes)),
				MaxIdleTimeoutInMinutes: pulumi.IntPtr(int(idle.MaxIdleTimeoutInMinutes)),
			},
		}
	}

	if emr := jl.EmrSettings; emr != nil {
		emrArgs := &sagemaker.UserProfileUserSettingsJupyterLabAppSettingsEmrSettingsArgs{}
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

func buildProfileJupyterServerAppSettings(js *awssagemakerdomainv1alpha1.AwsSagemakerDomainJupyterServerAppSettings) *sagemaker.UserProfileUserSettingsJupyterServerAppSettingsArgs {
	jsArgs := &sagemaker.UserProfileUserSettingsJupyterServerAppSettingsArgs{}

	if rs := js.DefaultResourceSpec; rs != nil {
		specArgs := &sagemaker.UserProfileUserSettingsJupyterServerAppSettingsDefaultResourceSpecArgs{}
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
		var repos sagemaker.UserProfileUserSettingsJupyterServerAppSettingsCodeRepositoryArray
		for _, repo := range js.CodeRepositories {
			repos = append(repos, &sagemaker.UserProfileUserSettingsJupyterServerAppSettingsCodeRepositoryArgs{
				RepositoryUrl: pulumi.String(repo.RepositoryUrl),
			})
		}
		jsArgs.CodeRepositories = repos
	}

	return jsArgs
}

func buildProfileKernelGatewayAppSettings(kg *awssagemakerdomainv1alpha1.AwsSagemakerDomainKernelGatewayAppSettings) *sagemaker.UserProfileUserSettingsKernelGatewayAppSettingsArgs {
	kgArgs := &sagemaker.UserProfileUserSettingsKernelGatewayAppSettingsArgs{}

	if rs := kg.DefaultResourceSpec; rs != nil {
		specArgs := &sagemaker.UserProfileUserSettingsKernelGatewayAppSettingsDefaultResourceSpecArgs{}
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
		var images sagemaker.UserProfileUserSettingsKernelGatewayAppSettingsCustomImageArray
		for _, img := range kg.CustomImages {
			imgArgs := &sagemaker.UserProfileUserSettingsKernelGatewayAppSettingsCustomImageArgs{
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

	return kgArgs
}

func buildProfileCodeEditorAppSettings(ce *awssagemakerdomainv1alpha1.AwsSagemakerDomainCodeEditorAppSettings) *sagemaker.UserProfileUserSettingsCodeEditorAppSettingsArgs {
	ceArgs := &sagemaker.UserProfileUserSettingsCodeEditorAppSettingsArgs{}

	if rs := ce.DefaultResourceSpec; rs != nil {
		specArgs := &sagemaker.UserProfileUserSettingsCodeEditorAppSettingsDefaultResourceSpecArgs{}
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
		ceArgs.DefaultResourceSpec = specArgs
	}

	if len(ce.LifecycleConfigArns) > 0 {
		ceArgs.LifecycleConfigArns = pulumi.ToStringArray(ce.LifecycleConfigArns)
	}

	if ce.BuiltInLifecycleConfigArn != "" {
		ceArgs.BuiltInLifecycleConfigArn = pulumi.StringPtr(ce.BuiltInLifecycleConfigArn)
	}

	if len(ce.CustomImages) > 0 {
		var images sagemaker.UserProfileUserSettingsCodeEditorAppSettingsCustomImageArray
		for _, img := range ce.CustomImages {
			imgArgs := &sagemaker.UserProfileUserSettingsCodeEditorAppSettingsCustomImageArgs{
				AppImageConfigName: pulumi.String(img.AppImageConfigName),
				ImageName:          pulumi.String(img.ImageName),
			}
			if img.ImageVersionNumber != nil {
				imgArgs.ImageVersionNumber = pulumi.IntPtr(int(img.GetImageVersionNumber()))
			}
			images = append(images, imgArgs)
		}
		ceArgs.CustomImages = images
	}

	if idle := ce.IdleSettings; idle != nil {
		ceArgs.AppLifecycleManagement = &sagemaker.UserProfileUserSettingsCodeEditorAppSettingsAppLifecycleManagementArgs{
			IdleSettings: &sagemaker.UserProfileUserSettingsCodeEditorAppSettingsAppLifecycleManagementIdleSettingsArgs{
				LifecycleManagement:     pulumi.String(lifecycleManagement(idle)),
				IdleTimeoutInMinutes:    pulumi.IntPtr(int(idle.IdleTimeoutInMinutes)),
				MinIdleTimeoutInMinutes: pulumi.IntPtr(int(idle.MinIdleTimeoutInMinutes)),
				MaxIdleTimeoutInMinutes: pulumi.IntPtr(int(idle.MaxIdleTimeoutInMinutes)),
			},
		}
	}

	return ceArgs
}

// buildProfileCanvasAppSettings assembles the Canvas (no-code ML) block. The spec
// flattens the SDK's two single-field wrappers (direct_deploy_settings.status
// and kendra_settings.status) to plain status scalars and Bedrock's
// generative_ai_settings wrapper to a single role ref; the wrappers are
// reconstructed here.
func buildProfileCanvasAppSettings(canvas *awssagemakerdomainv1alpha1.AwsSagemakerDomainCanvasAppSettings) *sagemaker.UserProfileUserSettingsCanvasAppSettingsArgs {
	canvasArgs := &sagemaker.UserProfileUserSettingsCanvasAppSettingsArgs{}

	if canvas.DirectDeployStatus != nil {
		canvasArgs.DirectDeploySettings = &sagemaker.UserProfileUserSettingsCanvasAppSettingsDirectDeploySettingsArgs{
			Status: pulumi.StringPtr(canvas.GetDirectDeployStatus()),
		}
	}

	if emr := canvas.EmrServerlessSettings; emr != nil {
		emrArgs := &sagemaker.UserProfileUserSettingsCanvasAppSettingsEmrServerlessSettingsArgs{}
		if emr.ExecutionRoleArn.GetValue() != "" {
			emrArgs.ExecutionRoleArn = pulumi.StringPtr(emr.ExecutionRoleArn.GetValue())
		}
		if emr.Status != "" {
			emrArgs.Status = pulumi.StringPtr(emr.Status)
		}
		canvasArgs.EmrServerlessSettings = emrArgs
	}

	// Setting the Bedrock role is what enables Canvas generative AI.
	if canvas.GenerativeAiBedrockRoleArn.GetValue() != "" {
		canvasArgs.GenerativeAiSettings = &sagemaker.UserProfileUserSettingsCanvasAppSettingsGenerativeAiSettingsArgs{
			AmazonBedrockRoleArn: pulumi.StringPtr(canvas.GenerativeAiBedrockRoleArn.GetValue()),
		}
	}

	if len(canvas.IdentityProviderOauthSettings) > 0 {
		var oauth sagemaker.UserProfileUserSettingsCanvasAppSettingsIdentityProviderOauthSettingArray
		for _, ip := range canvas.IdentityProviderOauthSettings {
			ipArgs := &sagemaker.UserProfileUserSettingsCanvasAppSettingsIdentityProviderOauthSettingArgs{
				SecretArn: pulumi.String(ip.SecretArn),
			}
			if ip.DataSourceName != "" {
				ipArgs.DataSourceName = pulumi.StringPtr(ip.DataSourceName)
			}
			if ip.Status != "" {
				ipArgs.Status = pulumi.StringPtr(ip.Status)
			}
			oauth = append(oauth, ipArgs)
		}
		canvasArgs.IdentityProviderOauthSettings = oauth
	}

	if canvas.KendraSettingsStatus != nil {
		canvasArgs.KendraSettings = &sagemaker.UserProfileUserSettingsCanvasAppSettingsKendraSettingsArgs{
			Status: pulumi.StringPtr(canvas.GetKendraSettingsStatus()),
		}
	}

	if mr := canvas.ModelRegisterSettings; mr != nil {
		mrArgs := &sagemaker.UserProfileUserSettingsCanvasAppSettingsModelRegisterSettingsArgs{}
		if mr.CrossAccountModelRegisterRoleArn != "" {
			mrArgs.CrossAccountModelRegisterRoleArn = pulumi.StringPtr(mr.CrossAccountModelRegisterRoleArn)
		}
		if mr.Status != "" {
			mrArgs.Status = pulumi.StringPtr(mr.Status)
		}
		canvasArgs.ModelRegisterSettings = mrArgs
	}

	if tsf := canvas.TimeSeriesForecastingSettings; tsf != nil {
		tsfArgs := &sagemaker.UserProfileUserSettingsCanvasAppSettingsTimeSeriesForecastingSettingsArgs{}
		if tsf.AmazonForecastRoleArn.GetValue() != "" {
			tsfArgs.AmazonForecastRoleArn = pulumi.StringPtr(tsf.AmazonForecastRoleArn.GetValue())
		}
		if tsf.Status != "" {
			tsfArgs.Status = pulumi.StringPtr(tsf.Status)
		}
		canvasArgs.TimeSeriesForecastingSettings = tsfArgs
	}

	if ws := canvas.WorkspaceSettings; ws != nil {
		wsArgs := &sagemaker.UserProfileUserSettingsCanvasAppSettingsWorkspaceSettingsArgs{}
		if ws.S3ArtifactPath != "" {
			wsArgs.S3ArtifactPath = pulumi.StringPtr(ws.S3ArtifactPath)
		}
		if ws.S3KmsKeyId.GetValue() != "" {
			wsArgs.S3KmsKeyId = pulumi.StringPtr(ws.S3KmsKeyId.GetValue())
		}
		canvasArgs.WorkspaceSettings = wsArgs
	}

	return canvasArgs
}
