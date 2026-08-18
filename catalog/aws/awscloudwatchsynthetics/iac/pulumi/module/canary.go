package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/synthetics"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// canary creates the canary and its group joins, and exports the
// canary outputs.
//
// Lifecycle facts the render below depends on:
//   - the canary's name is metadata.name (the canary charset is
//     lowercase letters, digits, hyphens, underscores) and renaming
//     replaces the canary;
//   - StartCanary drives Start/StopCanary around updates - false
//     creates the canary READY but never running (no run costs);
//   - a canary that lands in CREATE_FAILED is deleted and recreated by
//     the provider (AWS offers no other repair);
//   - RunConfig.EnvironmentVariables are WRITE-ONLY at AWS (reads
//     never return them) - never put secrets there;
//   - the association joins the canary by ARN and the group by NAME
//     (create/delete only - both sides replace on change); the
//     association join is untaggable at AWS.
func canary(ctx *pulumi.Context, locals *Locals, provider *aws.Provider, createdGroups map[string]*synthetics.Group) error {
	spec := locals.Spec

	if spec.Canary == nil {
		// Groups-only instance: export empty canary outputs so the
		// contract is stable across arms.
		ctx.Export(OpCanaryName, pulumi.String(""))
		ctx.Export(OpCanaryArn, pulumi.String(""))
		ctx.Export(OpEngineArn, pulumi.String(""))
		ctx.Export(OpSourceLocationArn, pulumi.String(""))
		ctx.Export(OpCanaryStatus, pulumi.String(""))
		return nil
	}

	canarySpec := spec.Canary

	args := &synthetics.CanaryArgs{
		Name:               pulumi.String(locals.Target.Metadata.Name),
		ArtifactS3Location: pulumi.String(locals.ArtifactS3Location),
		ExecutionRoleArn:   pulumi.String(canarySpec.ExecutionRoleArn.GetValue()),
		Handler:            pulumi.String(canarySpec.Handler),
		RuntimeVersion:     pulumi.String(canarySpec.RuntimeVersion),
		S3Bucket:           pulumi.String(canarySpec.Code.S3Bucket.GetValue()),
		S3Key:              pulumi.String(canarySpec.Code.S3Key),
		StartCanary:        pulumi.Bool(canarySpec.StartCanary),
		DeleteLambda:       pulumi.Bool(canarySpec.DeleteLambda),
		Tags:               pulumi.ToStringMap(locals.AwsTags),
	}

	if canarySpec.Code.S3Version != "" {
		args.S3Version = pulumi.String(canarySpec.Code.S3Version)
	}

	scheduleArgs := &synthetics.CanaryScheduleArgs{
		Expression: pulumi.String(canarySpec.Schedule.Expression),
	}
	if canarySpec.Schedule.DurationInSeconds != 0 {
		scheduleArgs.DurationInSeconds = pulumi.Int(int(canarySpec.Schedule.DurationInSeconds))
	}
	if canarySpec.Schedule.MaxRetries != nil {
		scheduleArgs.RetryConfig = &synthetics.CanaryScheduleRetryConfigArgs{
			MaxRetries: pulumi.Int(int(*canarySpec.Schedule.MaxRetries)),
		}
	}
	args.Schedule = scheduleArgs

	if runConfig := canarySpec.RunConfig; runConfig != nil {
		runConfigArgs := &synthetics.CanaryRunConfigArgs{
			ActiveTracing: pulumi.Bool(runConfig.ActiveTracing),
		}
		if len(runConfig.EnvironmentVariables) > 0 {
			runConfigArgs.EnvironmentVariables = pulumi.ToStringMap(runConfig.EnvironmentVariables)
		}
		if runConfig.EphemeralStorage != nil {
			runConfigArgs.EphemeralStorage = pulumi.Int(int(*runConfig.EphemeralStorage))
		}
		if runConfig.MemoryInMb != nil {
			runConfigArgs.MemoryInMb = pulumi.Int(int(*runConfig.MemoryInMb))
		}
		if runConfig.TimeoutInSeconds != nil {
			runConfigArgs.TimeoutInSeconds = pulumi.Int(int(*runConfig.TimeoutInSeconds))
		}
		args.RunConfig = runConfigArgs
	}

	if vpcConfig := canarySpec.VpcConfig; vpcConfig != nil {
		subnetIds := pulumi.StringArray{}
		for _, subnetId := range vpcConfig.SubnetIds {
			subnetIds = append(subnetIds, pulumi.String(subnetId.GetValue()))
		}
		securityGroupIds := pulumi.StringArray{}
		for _, securityGroupId := range vpcConfig.SecurityGroupIds {
			securityGroupIds = append(securityGroupIds, pulumi.String(securityGroupId.GetValue()))
		}
		vpcConfigArgs := &synthetics.CanaryVpcConfigArgs{
			SubnetIds:               subnetIds,
			Ipv6AllowedForDualStack: pulumi.Bool(vpcConfig.Ipv6AllowedForDualStack),
		}
		if len(securityGroupIds) > 0 {
			vpcConfigArgs.SecurityGroupIds = securityGroupIds
		}
		args.VpcConfig = vpcConfigArgs
	}

	if canarySpec.ArtifactEncryptionMode != "" || canarySpec.ArtifactEncryptionKmsKeyArn.GetValue() != "" {
		s3EncryptionArgs := &synthetics.CanaryArtifactConfigS3EncryptionArgs{}
		if canarySpec.ArtifactEncryptionMode != "" {
			s3EncryptionArgs.EncryptionMode = pulumi.String(canarySpec.ArtifactEncryptionMode)
		}
		if canarySpec.ArtifactEncryptionKmsKeyArn.GetValue() != "" {
			s3EncryptionArgs.KmsKeyArn = pulumi.String(canarySpec.ArtifactEncryptionKmsKeyArn.GetValue())
		}
		args.ArtifactConfig = &synthetics.CanaryArtifactConfigArgs{
			S3Encryption: s3EncryptionArgs,
		}
	}

	if canarySpec.FailureRetentionPeriod != nil {
		args.FailureRetentionPeriod = pulumi.Int(int(*canarySpec.FailureRetentionPeriod))
	}
	if canarySpec.SuccessRetentionPeriod != nil {
		args.SuccessRetentionPeriod = pulumi.Int(int(*canarySpec.SuccessRetentionPeriod))
	}

	createdCanary, err := synthetics.NewCanary(ctx, "canary", args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create canary")
	}

	// Group joins: owned groups join after their group resource exists;
	// external names resolve as-is.
	for _, groupName := range spec.GroupNames {
		options := []pulumi.ResourceOption{pulumi.Provider(provider)}
		if createdGroup, isOwned := createdGroups[groupName]; isOwned {
			options = append(options, pulumi.DependsOn([]pulumi.Resource{createdGroup}))
		}
		if _, err := synthetics.NewGroupAssociation(ctx, "association-"+groupName, &synthetics.GroupAssociationArgs{
			CanaryArn: createdCanary.Arn,
			GroupName: pulumi.String(groupName),
		}, options...); err != nil {
			return errors.Wrapf(err, "associate group %s", groupName)
		}
	}

	ctx.Export(OpCanaryName, createdCanary.Name)
	ctx.Export(OpCanaryArn, createdCanary.Arn)
	ctx.Export(OpEngineArn, createdCanary.EngineArn)
	ctx.Export(OpSourceLocationArn, createdCanary.SourceLocationArn)
	ctx.Export(OpCanaryStatus, createdCanary.Status)
	return nil
}
