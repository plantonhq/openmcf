package module

import (
	"fmt"

	"github.com/pkg/errors"
	awssagemakerdomainv1alpha1 "github.com/plantonhq/planton/catalog/aws/awssagemakerdomain/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/sagemaker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// spaces creates the folded aws_sagemaker_space satellites — one per
// spec.spaces entry, keyed by space name. A space's settings tree is
// DELIBERATELY different from the domain's default_space_settings (AWS uses
// distinct types: SpaceSettings vs DefaultSpaceSettings), so the builders
// below are space-specific rather than reusing the domain-baseline builders.
// Returns the spaces' ARNs and URLs keyed by name for the space_arns /
// space_urls outputs.
func spaces(ctx *pulumi.Context, locals *Locals, createdDomain *sagemaker.Domain, createdProfiles []pulumi.Resource, provider *aws.Provider) (pulumi.StringMap, pulumi.StringMap, error) {
	spaceArns := pulumi.StringMap{}
	spaceUrls := pulumi.StringMap{}

	for _, space := range locals.Spec.Spaces {
		args := &sagemaker.SpaceArgs{
			DomainId:  createdDomain.ID(),
			SpaceName: pulumi.String(space.SpaceName),
			Tags:      pulumi.ToStringMap(locals.AwsTags),
		}

		if space.DisplayName != "" {
			args.SpaceDisplayName = pulumi.String(space.DisplayName)
		}

		// ownership_settings and space_sharing_settings travel together
		// (provider RequiredWith, CEL-enforced) and are never sent on
		// update — the provider does not support changing them after
		// create.
		if space.OwnershipSettings != nil {
			args.OwnershipSettings = &sagemaker.SpaceOwnershipSettingsArgs{
				OwnerUserProfileName: pulumi.String(space.OwnershipSettings.OwnerUserProfileName),
			}
		}
		if space.SpaceSharingSettings != nil {
			args.SpaceSharingSettings = &sagemaker.SpaceSpaceSharingSettingsArgs{
				SharingType: pulumi.String(space.SpaceSharingSettings.SharingType),
			}
		}

		if space.SpaceSettings != nil {
			args.SpaceSettings = buildSpaceSettings(space.SpaceSettings)
		}

		// Ownership references the owner profile by NAME (no implicit engine
		// edge) -- spaces wait for every profile so an owner always exists.
		createdSpace, err := sagemaker.NewSpace(ctx,
			fmt.Sprintf("space-%s", space.SpaceName),
			args, pulumi.Provider(provider), pulumi.Parent(createdDomain),
			pulumi.DependsOn(createdProfiles))
		if err != nil {
			return nil, nil, errors.Wrapf(err, "space %s", space.SpaceName)
		}

		spaceArns[space.SpaceName] = createdSpace.Arn
		spaceUrls[space.SpaceName] = createdSpace.Url
	}

	return spaceArns, spaceUrls, nil
}

// buildSpaceSettings renders a space's own settings tree.
func buildSpaceSettings(ss *awssagemakerdomainv1alpha1.AwsSagemakerDomainSpaceSettings) *sagemaker.SpaceSpaceSettingsArgs {
	ssArgs := &sagemaker.SpaceSpaceSettingsArgs{}

	if ss.AppType != nil {
		ssArgs.AppType = pulumi.StringPtr(ss.GetAppType())
	}

	// JupyterLab: the resource spec is required on a space (CEL-enforced);
	// the space idle dial carries only the timeout — no lifecycle_management
	// switch and no min/max guardrails.
	if jl := ss.JupyterLabAppSettings; jl != nil {
		jlArgs := &sagemaker.SpaceSpaceSettingsJupyterLabAppSettingsArgs{
			DefaultResourceSpec: buildSpaceResourceSpec(jl.DefaultResourceSpec),
		}
		if len(jl.CodeRepositories) > 0 {
			var repos sagemaker.SpaceSpaceSettingsJupyterLabAppSettingsCodeRepositoryArray
			for _, repo := range jl.CodeRepositories {
				repos = append(repos, &sagemaker.SpaceSpaceSettingsJupyterLabAppSettingsCodeRepositoryArgs{
					RepositoryUrl: pulumi.String(repo.RepositoryUrl),
				})
			}
			jlArgs.CodeRepositories = repos
		}
		if idle := jl.IdleSettings; idle != nil && idle.IdleTimeoutInMinutes != nil {
			jlArgs.AppLifecycleManagement = &sagemaker.SpaceSpaceSettingsJupyterLabAppSettingsAppLifecycleManagementArgs{
				IdleSettings: &sagemaker.SpaceSpaceSettingsJupyterLabAppSettingsAppLifecycleManagementIdleSettingsArgs{
					IdleTimeoutInMinutes: pulumi.IntPtr(int(idle.GetIdleTimeoutInMinutes())),
				},
			}
		}
		ssArgs.JupyterLabAppSettings = jlArgs
	}

	// Code Editor: the space form (required resource spec, timeout-only idle).
	if ce := ss.CodeEditorAppSettings; ce != nil {
		ceArgs := &sagemaker.SpaceSpaceSettingsCodeEditorAppSettingsArgs{
			DefaultResourceSpec: buildSpaceCodeEditorResourceSpec(ce.DefaultResourceSpec),
		}
		if idle := ce.IdleSettings; idle != nil && idle.IdleTimeoutInMinutes != nil {
			ceArgs.AppLifecycleManagement = &sagemaker.SpaceSpaceSettingsCodeEditorAppSettingsAppLifecycleManagementArgs{
				IdleSettings: &sagemaker.SpaceSpaceSettingsCodeEditorAppSettingsAppLifecycleManagementIdleSettingsArgs{
					IdleTimeoutInMinutes: pulumi.IntPtr(int(idle.GetIdleTimeoutInMinutes())),
				},
			}
		}
		ssArgs.CodeEditorAppSettings = ceArgs
	}

	// Classic Jupyter Server (same shape as the domain baseline;
	// default_resource_spec required on a space, CEL-enforced).
	if js := ss.JupyterServerAppSettings; js != nil {
		jsArgs := &sagemaker.SpaceSpaceSettingsJupyterServerAppSettingsArgs{
			DefaultResourceSpec: buildSpaceJupyterServerResourceSpec(js.DefaultResourceSpec),
		}
		if len(js.LifecycleConfigArns) > 0 {
			jsArgs.LifecycleConfigArns = pulumi.ToStringArray(js.LifecycleConfigArns)
		}
		if len(js.CodeRepositories) > 0 {
			var repos sagemaker.SpaceSpaceSettingsJupyterServerAppSettingsCodeRepositoryArray
			for _, repo := range js.CodeRepositories {
				repos = append(repos, &sagemaker.SpaceSpaceSettingsJupyterServerAppSettingsCodeRepositoryArgs{
					RepositoryUrl: pulumi.String(repo.RepositoryUrl),
				})
			}
			jsArgs.CodeRepositories = repos
		}
		ssArgs.JupyterServerAppSettings = jsArgs
	}

	// KernelGateway (same shape as the domain baseline).
	if kg := ss.KernelGatewayAppSettings; kg != nil {
		kgArgs := &sagemaker.SpaceSpaceSettingsKernelGatewayAppSettingsArgs{
			DefaultResourceSpec: buildSpaceKernelGatewayResourceSpec(kg.DefaultResourceSpec),
		}
		if len(kg.LifecycleConfigArns) > 0 {
			kgArgs.LifecycleConfigArns = pulumi.ToStringArray(kg.LifecycleConfigArns)
		}
		if len(kg.CustomImages) > 0 {
			var images sagemaker.SpaceSpaceSettingsKernelGatewayAppSettingsCustomImageArray
			for _, img := range kg.CustomImages {
				imgArgs := &sagemaker.SpaceSpaceSettingsKernelGatewayAppSettingsCustomImageArgs{
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
		ssArgs.KernelGatewayAppSettings = kgArgs
	}

	// Mounted EFS file systems (by id — the space form has no per-mount
	// path, unlike the domain baseline's config).
	if len(ss.CustomFileSystems) > 0 {
		var fileSystems sagemaker.SpaceSpaceSettingsCustomFileSystemArray
		for _, fs := range ss.CustomFileSystems {
			fileSystems = append(fileSystems, &sagemaker.SpaceSpaceSettingsCustomFileSystemArgs{
				EfsFileSystem: &sagemaker.SpaceSpaceSettingsCustomFileSystemEfsFileSystemArgs{
					FileSystemId: pulumi.String(fs.FileSystemId.GetValue()),
				},
			})
		}
		ssArgs.CustomFileSystems = fileSystems
	}

	// The space's EBS volume (a single concrete size).
	if st := ss.SpaceStorageSettings; st != nil {
		ssArgs.SpaceStorageSettings = &sagemaker.SpaceSpaceSettingsSpaceStorageSettingsArgs{
			EbsStorageSettings: &sagemaker.SpaceSpaceSettingsSpaceStorageSettingsEbsStorageSettingsArgs{
				EbsVolumeSizeInGb: pulumi.Int(int(st.EbsVolumeSizeInGb)),
			},
		}
	}

	return ssArgs
}

func buildSpaceResourceSpec(rs *awssagemakerdomainv1alpha1.AwsSagemakerDomainResourceSpec) *sagemaker.SpaceSpaceSettingsJupyterLabAppSettingsDefaultResourceSpecArgs {
	args := &sagemaker.SpaceSpaceSettingsJupyterLabAppSettingsDefaultResourceSpecArgs{}
	if rs.InstanceType != "" {
		args.InstanceType = pulumi.StringPtr(rs.InstanceType)
	}
	if rs.LifecycleConfigArn != "" {
		args.LifecycleConfigArn = pulumi.StringPtr(rs.LifecycleConfigArn)
	}
	if rs.SagemakerImageArn != "" {
		args.SagemakerImageArn = pulumi.StringPtr(rs.SagemakerImageArn)
	}
	if rs.SagemakerImageVersionAlias != "" {
		args.SagemakerImageVersionAlias = pulumi.StringPtr(rs.SagemakerImageVersionAlias)
	}
	if rs.SagemakerImageVersionArn != "" {
		args.SagemakerImageVersionArn = pulumi.StringPtr(rs.SagemakerImageVersionArn)
	}
	return args
}

func buildSpaceCodeEditorResourceSpec(rs *awssagemakerdomainv1alpha1.AwsSagemakerDomainResourceSpec) *sagemaker.SpaceSpaceSettingsCodeEditorAppSettingsDefaultResourceSpecArgs {
	args := &sagemaker.SpaceSpaceSettingsCodeEditorAppSettingsDefaultResourceSpecArgs{}
	if rs.InstanceType != "" {
		args.InstanceType = pulumi.StringPtr(rs.InstanceType)
	}
	if rs.LifecycleConfigArn != "" {
		args.LifecycleConfigArn = pulumi.StringPtr(rs.LifecycleConfigArn)
	}
	if rs.SagemakerImageArn != "" {
		args.SagemakerImageArn = pulumi.StringPtr(rs.SagemakerImageArn)
	}
	if rs.SagemakerImageVersionAlias != "" {
		args.SagemakerImageVersionAlias = pulumi.StringPtr(rs.SagemakerImageVersionAlias)
	}
	if rs.SagemakerImageVersionArn != "" {
		args.SagemakerImageVersionArn = pulumi.StringPtr(rs.SagemakerImageVersionArn)
	}
	return args
}

func buildSpaceJupyterServerResourceSpec(rs *awssagemakerdomainv1alpha1.AwsSagemakerDomainResourceSpec) *sagemaker.SpaceSpaceSettingsJupyterServerAppSettingsDefaultResourceSpecArgs {
	args := &sagemaker.SpaceSpaceSettingsJupyterServerAppSettingsDefaultResourceSpecArgs{}
	if rs.InstanceType != "" {
		args.InstanceType = pulumi.StringPtr(rs.InstanceType)
	}
	if rs.LifecycleConfigArn != "" {
		args.LifecycleConfigArn = pulumi.StringPtr(rs.LifecycleConfigArn)
	}
	if rs.SagemakerImageArn != "" {
		args.SagemakerImageArn = pulumi.StringPtr(rs.SagemakerImageArn)
	}
	if rs.SagemakerImageVersionAlias != "" {
		args.SagemakerImageVersionAlias = pulumi.StringPtr(rs.SagemakerImageVersionAlias)
	}
	if rs.SagemakerImageVersionArn != "" {
		args.SagemakerImageVersionArn = pulumi.StringPtr(rs.SagemakerImageVersionArn)
	}
	return args
}

func buildSpaceKernelGatewayResourceSpec(rs *awssagemakerdomainv1alpha1.AwsSagemakerDomainResourceSpec) *sagemaker.SpaceSpaceSettingsKernelGatewayAppSettingsDefaultResourceSpecArgs {
	args := &sagemaker.SpaceSpaceSettingsKernelGatewayAppSettingsDefaultResourceSpecArgs{}
	if rs.InstanceType != "" {
		args.InstanceType = pulumi.StringPtr(rs.InstanceType)
	}
	if rs.LifecycleConfigArn != "" {
		args.LifecycleConfigArn = pulumi.StringPtr(rs.LifecycleConfigArn)
	}
	if rs.SagemakerImageArn != "" {
		args.SagemakerImageArn = pulumi.StringPtr(rs.SagemakerImageArn)
	}
	if rs.SagemakerImageVersionAlias != "" {
		args.SagemakerImageVersionAlias = pulumi.StringPtr(rs.SagemakerImageVersionAlias)
	}
	if rs.SagemakerImageVersionArn != "" {
		args.SagemakerImageVersionArn = pulumi.StringPtr(rs.SagemakerImageVersionArn)
	}
	return args
}
