package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/sagemaker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// image creates the SageMaker image and its folded versions and exports
// outputs.
//
// Lifecycle facts the renders below depend on:
//   - the provider sleeps ~1 minute before CreateImage (IAM
//     propagation) - every create is at least a minute;
//   - version numbers are AWS-assigned, monotonic, never reused: a
//     changed base_image REPLACES the version under a NEW number.
//     Entries are keyed by position - append-only, taught on the spec;
//   - image_version carries NO tags upstream (by provider design);
//   - AWS serializes version creation per image (the provider holds a
//     mutex) - versions attach one at a time.
func image(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	args := &sagemaker.ImageArgs{
		// The component's name IS the image name.
		ImageName: pulumi.String(locals.ImageName),
		RoleArn:   pulumi.String(spec.RoleArn.GetValue()),
		Tags:      pulumi.ToStringMap(locals.AwsTags),
	}
	if spec.DisplayName != "" {
		args.DisplayName = pulumi.String(spec.DisplayName)
	}
	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}

	createdImage, err := sagemaker.NewImage(ctx, locals.ImageName, args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create image")
	}

	// Versions keyed by their stable POSITION (the append-only contract
	// taught on the spec) - identical keys to the Terraform module's
	// for_each.
	versionNumbers := pulumi.StringMap{}
	for i, v := range spec.Versions {
		key := fmt.Sprintf("%d", i)
		versionArgs := &sagemaker.ImageVersionArgs{
			ImageName: createdImage.ImageName,
			// The version's identity - changing it replaces the version
			// under a new AWS-assigned number.
			BaseImage: pulumi.String(v.BaseImage),
		}
		if len(v.Aliases) > 0 {
			versionArgs.Aliases = pulumi.ToStringArray(v.Aliases)
		}
		if v.Horovod {
			versionArgs.Horovod = pulumi.Bool(true)
		}
		if v.JobType != "" {
			versionArgs.JobType = pulumi.String(v.JobType)
		}
		if v.MlFramework != "" {
			versionArgs.MlFramework = pulumi.String(v.MlFramework)
		}
		if v.Processor != "" {
			versionArgs.Processor = pulumi.String(v.Processor)
		}
		if v.ProgrammingLang != "" {
			versionArgs.ProgrammingLang = pulumi.String(v.ProgrammingLang)
		}
		if v.ReleaseNotes != "" {
			versionArgs.ReleaseNotes = pulumi.String(v.ReleaseNotes)
		}
		if v.VendorGuidance != "" {
			versionArgs.VendorGuidance = pulumi.String(v.VendorGuidance)
		}
		createdVersion, err := sagemaker.NewImageVersion(ctx, "version-"+key, versionArgs,
			pulumi.Provider(provider), pulumi.DependsOn([]pulumi.Resource{createdImage}))
		if err != nil {
			return errors.Wrapf(err, "create image version %s", key)
		}
		versionNumbers[key] = createdVersion.Version.ApplyT(func(version int) string {
			return fmt.Sprintf("%d", version)
		}).(pulumi.StringOutput)
	}

	ctx.Export(OpImageName, createdImage.ImageName)
	ctx.Export(OpImageArn, createdImage.Arn)
	ctx.Export(OpVersionNumbers, versionNumbers)

	return nil
}
