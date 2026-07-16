package module

import (
	awssagemakerdomainv1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awssagemakerdomain/v1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/sagemaker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// buildDefaultUserSettings assembles the per-user baseline every user profile
// inherits. Each optional spec message maps to its SDK block only when
// present; the SDK models every nesting path as its own Go type, so the
// repeated default_resource_spec / custom_images / idle_settings shapes are
// converted per app type below (the spec shares one message across them).
func buildDefaultUserSettings(dus *awssagemakerdomainv1.AwsSagemakerDomainDefaultUserSettings) *sagemaker.DomainDefaultUserSettingsArgs {
	settings := &sagemaker.DomainDefaultUserSettingsArgs{
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
		settings.JupyterLabAppSettings = buildUserJupyterLabAppSettings(jl)
	}

	if js := dus.JupyterServerAppSettings; js != nil {
		settings.JupyterServerAppSettings = buildUserJupyterServerAppSettings(js)
	}

	if kg := dus.KernelGatewayAppSettings; kg != nil {
		settings.KernelGatewayAppSettings = buildUserKernelGatewayAppSettings(kg)
	}

	if ce := dus.CodeEditorAppSettings; ce != nil {
		settings.CodeEditorAppSettings = buildUserCodeEditorAppSettings(ce)
	}

	if tb := dus.TensorBoardAppSettings; tb != nil {
		tbArgs := &sagemaker.DomainDefaultUserSettingsTensorBoardAppSettingsArgs{}
		if rs := tb.DefaultResourceSpec; rs != nil {
			specArgs := &sagemaker.DomainDefaultUserSettingsTensorBoardAppSettingsDefaultResourceSpecArgs{}
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
		rsessArgs := &sagemaker.DomainDefaultUserSettingsRSessionAppSettingsArgs{}
		if rs := rsess.DefaultResourceSpec; rs != nil {
			specArgs := &sagemaker.DomainDefaultUserSettingsRSessionAppSettingsDefaultResourceSpecArgs{}
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
			var images sagemaker.DomainDefaultUserSettingsRSessionAppSettingsCustomImageArray
			for _, img := range rsess.CustomImages {
				imgArgs := &sagemaker.DomainDefaultUserSettingsRSessionAppSettingsCustomImageArgs{
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
		rspArgs := &sagemaker.DomainDefaultUserSettingsRStudioServerProAppSettingsArgs{}
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
		settings.CanvasAppSettings = buildCanvasAppSettings(canvas)
	}

	if ss := dus.SharingSettings; ss != nil {
		ssArgs := &sagemaker.DomainDefaultUserSettingsSharingSettingsArgs{}
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
		settings.SpaceStorageSettings = &sagemaker.DomainDefaultUserSettingsSpaceStorageSettingsArgs{
			DefaultEbsStorageSettings: &sagemaker.DomainDefaultUserSettingsSpaceStorageSettingsDefaultEbsStorageSettingsArgs{
				DefaultEbsVolumeSizeInGb: pulumi.Int(int(sss.DefaultEbsVolumeSizeInGb)),
				MaximumEbsVolumeSizeInGb: pulumi.Int(int(sss.MaximumEbsVolumeSizeInGb)),
			},
		}
	}

	if len(dus.CustomFileSystemConfigs) > 0 {
		var configs sagemaker.DomainDefaultUserSettingsCustomFileSystemConfigArray
		for _, cfs := range dus.CustomFileSystemConfigs {
			efs := cfs.EfsFileSystemConfig
			configs = append(configs, &sagemaker.DomainDefaultUserSettingsCustomFileSystemConfigArgs{
				EfsFileSystemConfig: &sagemaker.DomainDefaultUserSettingsCustomFileSystemConfigEfsFileSystemConfigArgs{
					FileSystemId:   pulumi.String(efs.FileSystemId.GetValue()),
					FileSystemPath: pulumi.String(efs.FileSystemPath),
				},
			})
		}
		settings.CustomFileSystemConfigs = configs
	}

	if posix := dus.CustomPosixUserConfig; posix != nil {
		settings.CustomPosixUserConfig = &sagemaker.DomainDefaultUserSettingsCustomPosixUserConfigArgs{
			Uid: pulumi.Int(int(posix.Uid)),
			Gid: pulumi.Int(int(posix.Gid)),
		}
	}

	if wps := dus.StudioWebPortalSettings; wps != nil {
		wpsArgs := &sagemaker.DomainDefaultUserSettingsStudioWebPortalSettingsArgs{}
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

func buildUserJupyterLabAppSettings(jl *awssagemakerdomainv1.AwsSagemakerDomainJupyterLabAppSettings) *sagemaker.DomainDefaultUserSettingsJupyterLabAppSettingsArgs {
	jlArgs := &sagemaker.DomainDefaultUserSettingsJupyterLabAppSettingsArgs{}

	if rs := jl.DefaultResourceSpec; rs != nil {
		specArgs := &sagemaker.DomainDefaultUserSettingsJupyterLabAppSettingsDefaultResourceSpecArgs{}
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
		var images sagemaker.DomainDefaultUserSettingsJupyterLabAppSettingsCustomImageArray
		for _, img := range jl.CustomImages {
			imgArgs := &sagemaker.DomainDefaultUserSettingsJupyterLabAppSettingsCustomImageArgs{
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
		var repos sagemaker.DomainDefaultUserSettingsJupyterLabAppSettingsCodeRepositoryArray
		for _, repo := range jl.CodeRepositories {
			repos = append(repos, &sagemaker.DomainDefaultUserSettingsJupyterLabAppSettingsCodeRepositoryArgs{
				RepositoryUrl: pulumi.String(repo.RepositoryUrl),
			})
		}
		jlArgs.CodeRepositories = repos
	}

	// The spec folds idle_settings directly under the app settings and makes
	// block presence the enable switch; the SDK nests them inside a
	// single-purpose app_lifecycle_management wrapper with an explicit
	// ENABLED flag, both reconstructed here. All three timeouts are required
	// by the live API whenever the block is sent (absent members transmit as
	// 0 and AWS rejects them), so they pass through unconditionally.
	if idle := jl.IdleSettings; idle != nil {
		jlArgs.AppLifecycleManagement = &sagemaker.DomainDefaultUserSettingsJupyterLabAppSettingsAppLifecycleManagementArgs{
			IdleSettings: &sagemaker.DomainDefaultUserSettingsJupyterLabAppSettingsAppLifecycleManagementIdleSettingsArgs{
				LifecycleManagement:     pulumi.String("ENABLED"),
				IdleTimeoutInMinutes:    pulumi.IntPtr(int(idle.IdleTimeoutInMinutes)),
				MinIdleTimeoutInMinutes: pulumi.IntPtr(int(idle.MinIdleTimeoutInMinutes)),
				MaxIdleTimeoutInMinutes: pulumi.IntPtr(int(idle.MaxIdleTimeoutInMinutes)),
			},
		}
	}

	if emr := jl.EmrSettings; emr != nil {
		emrArgs := &sagemaker.DomainDefaultUserSettingsJupyterLabAppSettingsEmrSettingsArgs{}
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

func buildUserJupyterServerAppSettings(js *awssagemakerdomainv1.AwsSagemakerDomainJupyterServerAppSettings) *sagemaker.DomainDefaultUserSettingsJupyterServerAppSettingsArgs {
	jsArgs := &sagemaker.DomainDefaultUserSettingsJupyterServerAppSettingsArgs{}

	if rs := js.DefaultResourceSpec; rs != nil {
		specArgs := &sagemaker.DomainDefaultUserSettingsJupyterServerAppSettingsDefaultResourceSpecArgs{}
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
		var repos sagemaker.DomainDefaultUserSettingsJupyterServerAppSettingsCodeRepositoryArray
		for _, repo := range js.CodeRepositories {
			repos = append(repos, &sagemaker.DomainDefaultUserSettingsJupyterServerAppSettingsCodeRepositoryArgs{
				RepositoryUrl: pulumi.String(repo.RepositoryUrl),
			})
		}
		jsArgs.CodeRepositories = repos
	}

	return jsArgs
}

func buildUserKernelGatewayAppSettings(kg *awssagemakerdomainv1.AwsSagemakerDomainKernelGatewayAppSettings) *sagemaker.DomainDefaultUserSettingsKernelGatewayAppSettingsArgs {
	kgArgs := &sagemaker.DomainDefaultUserSettingsKernelGatewayAppSettingsArgs{}

	if rs := kg.DefaultResourceSpec; rs != nil {
		specArgs := &sagemaker.DomainDefaultUserSettingsKernelGatewayAppSettingsDefaultResourceSpecArgs{}
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
		var images sagemaker.DomainDefaultUserSettingsKernelGatewayAppSettingsCustomImageArray
		for _, img := range kg.CustomImages {
			imgArgs := &sagemaker.DomainDefaultUserSettingsKernelGatewayAppSettingsCustomImageArgs{
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

func buildUserCodeEditorAppSettings(ce *awssagemakerdomainv1.AwsSagemakerDomainCodeEditorAppSettings) *sagemaker.DomainDefaultUserSettingsCodeEditorAppSettingsArgs {
	ceArgs := &sagemaker.DomainDefaultUserSettingsCodeEditorAppSettingsArgs{}

	if rs := ce.DefaultResourceSpec; rs != nil {
		specArgs := &sagemaker.DomainDefaultUserSettingsCodeEditorAppSettingsDefaultResourceSpecArgs{}
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
		var images sagemaker.DomainDefaultUserSettingsCodeEditorAppSettingsCustomImageArray
		for _, img := range ce.CustomImages {
			imgArgs := &sagemaker.DomainDefaultUserSettingsCodeEditorAppSettingsCustomImageArgs{
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
		ceArgs.AppLifecycleManagement = &sagemaker.DomainDefaultUserSettingsCodeEditorAppSettingsAppLifecycleManagementArgs{
			IdleSettings: &sagemaker.DomainDefaultUserSettingsCodeEditorAppSettingsAppLifecycleManagementIdleSettingsArgs{
				LifecycleManagement:     pulumi.String("ENABLED"),
				IdleTimeoutInMinutes:    pulumi.IntPtr(int(idle.IdleTimeoutInMinutes)),
				MinIdleTimeoutInMinutes: pulumi.IntPtr(int(idle.MinIdleTimeoutInMinutes)),
				MaxIdleTimeoutInMinutes: pulumi.IntPtr(int(idle.MaxIdleTimeoutInMinutes)),
			},
		}
	}

	return ceArgs
}

// buildCanvasAppSettings assembles the Canvas (no-code ML) block. The spec
// flattens the SDK's two single-field wrappers (direct_deploy_settings.status
// and kendra_settings.status) to plain status scalars and Bedrock's
// generative_ai_settings wrapper to a single role ref; the wrappers are
// reconstructed here.
func buildCanvasAppSettings(canvas *awssagemakerdomainv1.AwsSagemakerDomainCanvasAppSettings) *sagemaker.DomainDefaultUserSettingsCanvasAppSettingsArgs {
	canvasArgs := &sagemaker.DomainDefaultUserSettingsCanvasAppSettingsArgs{}

	if canvas.DirectDeployStatus != nil {
		canvasArgs.DirectDeploySettings = &sagemaker.DomainDefaultUserSettingsCanvasAppSettingsDirectDeploySettingsArgs{
			Status: pulumi.StringPtr(canvas.GetDirectDeployStatus()),
		}
	}

	if emr := canvas.EmrServerlessSettings; emr != nil {
		emrArgs := &sagemaker.DomainDefaultUserSettingsCanvasAppSettingsEmrServerlessSettingsArgs{}
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
		canvasArgs.GenerativeAiSettings = &sagemaker.DomainDefaultUserSettingsCanvasAppSettingsGenerativeAiSettingsArgs{
			AmazonBedrockRoleArn: pulumi.StringPtr(canvas.GenerativeAiBedrockRoleArn.GetValue()),
		}
	}

	if len(canvas.IdentityProviderOauthSettings) > 0 {
		var oauth sagemaker.DomainDefaultUserSettingsCanvasAppSettingsIdentityProviderOauthSettingArray
		for _, ip := range canvas.IdentityProviderOauthSettings {
			ipArgs := &sagemaker.DomainDefaultUserSettingsCanvasAppSettingsIdentityProviderOauthSettingArgs{
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
		canvasArgs.KendraSettings = &sagemaker.DomainDefaultUserSettingsCanvasAppSettingsKendraSettingsArgs{
			Status: pulumi.StringPtr(canvas.GetKendraSettingsStatus()),
		}
	}

	if mr := canvas.ModelRegisterSettings; mr != nil {
		mrArgs := &sagemaker.DomainDefaultUserSettingsCanvasAppSettingsModelRegisterSettingsArgs{}
		if mr.CrossAccountModelRegisterRoleArn != "" {
			mrArgs.CrossAccountModelRegisterRoleArn = pulumi.StringPtr(mr.CrossAccountModelRegisterRoleArn)
		}
		if mr.Status != "" {
			mrArgs.Status = pulumi.StringPtr(mr.Status)
		}
		canvasArgs.ModelRegisterSettings = mrArgs
	}

	if tsf := canvas.TimeSeriesForecastingSettings; tsf != nil {
		tsfArgs := &sagemaker.DomainDefaultUserSettingsCanvasAppSettingsTimeSeriesForecastingSettingsArgs{}
		if tsf.AmazonForecastRoleArn.GetValue() != "" {
			tsfArgs.AmazonForecastRoleArn = pulumi.StringPtr(tsf.AmazonForecastRoleArn.GetValue())
		}
		if tsf.Status != "" {
			tsfArgs.Status = pulumi.StringPtr(tsf.Status)
		}
		canvasArgs.TimeSeriesForecastingSettings = tsfArgs
	}

	if ws := canvas.WorkspaceSettings; ws != nil {
		wsArgs := &sagemaker.DomainDefaultUserSettingsCanvasAppSettingsWorkspaceSettingsArgs{}
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
