package module

import (
	"encoding/base64"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/sagemaker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// notebookInstance creates the notebook instance and its folded
// lifecycle configuration and exports outputs.
//
// Lifecycle facts the renders below depend on:
//   - most instance changes ride the provider's stop-update-start
//     choreography (SageMaker requires a Stopped instance for
//     UpdateNotebookInstance) - budget several minutes per change;
//   - growing volume_size updates in place, SHRINKING replaces the
//     instance (provider-enforced, mirroring AWS's no-shrink rule);
//   - the lifecycle scripts are sent base64-encoded (the module
//     encodes; the spec carries plain shell) and run as root with a
//     5-minute limit;
//   - clearing a script upstream does NOT clear it in AWS (the
//     provider's update omits empty fields) - replacing the text is the
//     reliable path, taught on the spec fields.
func notebookInstance(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	args := &sagemaker.NotebookInstanceArgs{
		// The component's name IS the notebook instance name.
		Name:         pulumi.String(locals.NotebookName),
		InstanceType: pulumi.String(spec.InstanceType),
		RoleArn:      pulumi.String(spec.RoleArn.GetValue()),
		Tags:         pulumi.ToStringMap(locals.AwsTags),
	}

	if spec.VolumeSizeGb != nil {
		args.VolumeSize = pulumi.Int(int(*spec.VolumeSizeGb))
	}
	if spec.SubnetId.GetValue() != "" {
		args.SubnetId = pulumi.String(spec.SubnetId.GetValue())
	}
	if len(spec.SecurityGroupIds) > 0 {
		var securityGroups pulumi.StringArray
		for _, s := range spec.SecurityGroupIds {
			securityGroups = append(securityGroups, pulumi.String(s.GetValue()))
		}
		args.SecurityGroups = securityGroups
	}
	if spec.KmsKeyArn.GetValue() != "" {
		args.KmsKeyId = pulumi.String(spec.KmsKeyArn.GetValue())
	}
	// Changing either replaces the instance (provider-enforced).
	if spec.DirectInternetAccess != "" {
		args.DirectInternetAccess = pulumi.String(spec.DirectInternetAccess)
	}
	if spec.PlatformIdentifier != "" {
		args.PlatformIdentifier = pulumi.String(spec.PlatformIdentifier)
	}
	if spec.RootAccess != "" {
		args.RootAccess = pulumi.String(spec.RootAccess)
	}
	if spec.DefaultCodeRepository != "" {
		args.DefaultCodeRepository = pulumi.String(spec.DefaultCodeRepository)
	}
	if len(spec.AdditionalCodeRepositories) > 0 {
		args.AdditionalCodeRepositories = pulumi.ToStringArray(spec.AdditionalCodeRepositories)
	}
	if spec.ImdsMinimumVersion != "" {
		args.InstanceMetadataServiceConfiguration = &sagemaker.NotebookInstanceInstanceMetadataServiceConfigurationArgs{
			MinimumInstanceMetadataServiceVersion: pulumi.String(spec.ImdsMinimumVersion),
		}
	}

	// The folded lifecycle configuration (bootstrap scripts).
	var lifecycleConfigName pulumi.StringOutput
	hasLifecycle := spec.LifecycleConfig != nil
	if hasLifecycle {
		lifecycleArgs := &sagemaker.NotebookInstanceLifecycleConfigurationArgs{
			Name: pulumi.String(locals.LifecycleConfigName),
			Tags: pulumi.ToStringMap(locals.AwsTags),
		}
		if spec.LifecycleConfig.OnCreate != "" {
			lifecycleArgs.OnCreate = pulumi.String(base64.StdEncoding.EncodeToString([]byte(spec.LifecycleConfig.OnCreate)))
		}
		if spec.LifecycleConfig.OnStart != "" {
			lifecycleArgs.OnStart = pulumi.String(base64.StdEncoding.EncodeToString([]byte(spec.LifecycleConfig.OnStart)))
		}
		createdLifecycleConfig, err := sagemaker.NewNotebookInstanceLifecycleConfiguration(ctx,
			locals.LifecycleConfigName, lifecycleArgs, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrap(err, "create lifecycle configuration")
		}
		lifecycleConfigName = createdLifecycleConfig.Name
		args.LifecycleConfigName = lifecycleConfigName
	}

	createdInstance, err := sagemaker.NewNotebookInstance(ctx, locals.NotebookName, args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create notebook instance")
	}

	ctx.Export(OpNotebookInstanceName, createdInstance.Name)
	ctx.Export(OpNotebookInstanceArn, createdInstance.Arn)
	ctx.Export(OpUrl, createdInstance.Url)
	ctx.Export(OpNetworkInterfaceId, createdInstance.NetworkInterfaceId)
	if hasLifecycle {
		ctx.Export(OpLifecycleConfigName, lifecycleConfigName)
	} else {
		ctx.Export(OpLifecycleConfigName, pulumi.String(""))
	}

	return nil
}
