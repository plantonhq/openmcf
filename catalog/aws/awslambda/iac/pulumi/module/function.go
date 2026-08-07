package module

import (
	"github.com/pkg/errors"
	"github.com/plantonhq/planton/internal/valuefrom"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// function provisions the Lambda function itself. Create-time immutable
// in AWS: the function name and the package type (zip vs container
// image). Everything else -- code, memory, timeout, VPC attachment,
// layers, logging -- edits in place (code changes roll the function
// without replacing it).
func function(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*lambda.Function, error) {
	spec := locals.AwsLambda.Spec

	args := &lambda.FunctionArgs{
		// The function name is metadata.name on both engines, so a
		// manifest deploys the same function regardless of engine --
		// and the default log group "/aws/lambda/<name>" is
		// predictable (AWS would otherwise auto-name with a random
		// suffix).
		Name: pulumi.String(locals.FunctionName),
		Role: pulumi.String(spec.RoleArn.GetValue()),
		Tags: pulumi.ToStringMap(locals.AwsTags),
	}

	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}

	// Zip vs container image is a create-time fork (CEL guarantees
	// exactly one source): package_type cannot change on a live
	// function.
	if spec.ImageUri != "" {
		args.PackageType = pulumi.String("Image")
		args.ImageUri = pulumi.String(spec.ImageUri)
	} else {
		args.PackageType = pulumi.String("Zip")
		args.S3Bucket = pulumi.String(spec.S3.Bucket.GetValue())
		args.S3Key = pulumi.String(spec.S3.Key)
		if spec.S3.ObjectVersion != "" {
			args.S3ObjectVersion = pulumi.String(spec.S3.ObjectVersion)
		}
		// Runtime/handler drive zip execution; images carry their own
		// (CEL keeps them empty in image mode).
		args.Runtime = pulumi.String(spec.Runtime)
		args.Handler = pulumi.String(spec.Handler)
	}

	// The hash is what makes code rolls declarative: a new value rolls
	// the function even when the S3 key is rewritten in place; an
	// unchanged value is a no-op.
	if spec.SourceCodeHash != "" {
		args.SourceCodeHash = pulumi.String(spec.SourceCodeHash)
	}

	// BYOK for the deployment package itself -- distinct from
	// KmsKeyArn, which encrypts environment variables.
	if spec.SourceKmsKeyArn.GetValue() != "" {
		args.SourceKmsKeyArn = pulumi.String(spec.SourceKmsKeyArn.GetValue())
	}

	// The provider models architecture as a single-element list. Empty
	// keeps the AWS default (x86_64).
	if spec.Architecture != "" {
		args.Architectures = pulumi.StringArray{pulumi.String(spec.Architecture)}
	}

	// 0 keeps the AWS defaults (128 MB / 3 s / 512 MB scratch).
	if spec.MemorySizeMb != 0 {
		args.MemorySize = pulumi.Int(int(spec.MemorySizeMb))
	}
	if spec.TimeoutSeconds != 0 {
		args.Timeout = pulumi.Int(int(spec.TimeoutSeconds))
	}
	if spec.EphemeralStorageMb != 0 {
		args.EphemeralStorage = &lambda.FunctionEphemeralStorageArgs{
			Size: pulumi.Int(int(spec.EphemeralStorageMb)),
		}
	}

	// nil in the spec means "draw from the unreserved account pool",
	// which the provider expresses as -1; 0 is the explicit kill
	// switch.
	if spec.ReservedConcurrentExecutions != nil {
		args.ReservedConcurrentExecutions = pulumi.Int(int(*spec.ReservedConcurrentExecutions))
	} else {
		args.ReservedConcurrentExecutions = pulumi.Int(-1)
	}

	args.Publish = pulumi.Bool(spec.Publish)

	if len(spec.Environment) > 0 {
		args.Environment = &lambda.FunctionEnvironmentArgs{
			Variables: pulumi.ToStringMap(spec.Environment),
		}
	}

	if spec.KmsKeyArn.GetValue() != "" {
		args.KmsKeyArn = pulumi.String(spec.KmsKeyArn.GetValue())
	}

	if spec.CodeSigningConfigArn != "" {
		args.CodeSigningConfigArn = pulumi.String(spec.CodeSigningConfigArn)
	}

	if len(spec.LayerArns) > 0 {
		args.Layers = pulumi.ToStringArray(valuefrom.ToStringArray(spec.LayerArns))
	}

	// VPC attachment travels as a set (CEL): subnets + security groups
	// together, IPv6 egress only on top of them.
	if len(spec.SubnetIds) > 0 {
		args.VpcConfig = &lambda.FunctionVpcConfigArgs{
			SubnetIds:               pulumi.ToStringArray(valuefrom.ToStringArray(spec.SubnetIds)),
			SecurityGroupIds:        pulumi.ToStringArray(valuefrom.ToStringArray(spec.SecurityGroupIds)),
			Ipv6AllowedForDualStack: pulumi.Bool(spec.Ipv6AllowedForDualStack),
		}
	}

	if spec.DeadLetterTargetArn.GetValue() != "" {
		args.DeadLetterConfig = &lambda.FunctionDeadLetterConfigArgs{
			TargetArn: pulumi.String(spec.DeadLetterTargetArn.GetValue()),
		}
	}

	if spec.TracingMode != "" {
		args.TracingConfig = &lambda.FunctionTracingConfigArgs{
			Mode: pulumi.String(spec.TracingMode),
		}
	}

	// EFS mounts require the VPC attachment to reach the file system's
	// mount targets (CEL enforces the coupling).
	if spec.FileSystemConfig != nil {
		args.FileSystemConfig = &lambda.FunctionFileSystemConfigArgs{
			Arn:            pulumi.String(spec.FileSystemConfig.AccessPointArn.GetValue()),
			LocalMountPath: pulumi.String(spec.FileSystemConfig.LocalMountPath),
		}
	}

	if spec.ImageConfig != nil {
		imageConfigArgs := &lambda.FunctionImageConfigArgs{}
		if len(spec.ImageConfig.EntryPoint) > 0 {
			imageConfigArgs.EntryPoints = pulumi.ToStringArray(spec.ImageConfig.EntryPoint)
		}
		if len(spec.ImageConfig.Command) > 0 {
			imageConfigArgs.Commands = pulumi.ToStringArray(spec.ImageConfig.Command)
		}
		if spec.ImageConfig.WorkingDirectory != "" {
			imageConfigArgs.WorkingDirectory = pulumi.String(spec.ImageConfig.WorkingDirectory)
		}
		args.ImageConfig = imageConfigArgs
	}

	// SnapStart snapshots published versions only (CEL couples it to
	// publish) -- "None" is the provider's off state, so the block is
	// only rendered when enabled.
	if spec.SnapStart {
		args.SnapStart = &lambda.FunctionSnapStartArgs{
			ApplyOn: pulumi.String("PublishedVersions"),
		}
	}

	if spec.LoggingConfig != nil {
		// The provider requires log_format whenever the block is
		// present; "Text" is the AWS default.
		logFormat := spec.LoggingConfig.LogFormat
		if logFormat == "" {
			logFormat = "Text"
		}
		loggingArgs := &lambda.FunctionLoggingConfigArgs{
			LogFormat: pulumi.String(logFormat),
		}
		if spec.LoggingConfig.ApplicationLogLevel != "" {
			loggingArgs.ApplicationLogLevel = pulumi.String(spec.LoggingConfig.ApplicationLogLevel)
		}
		if spec.LoggingConfig.SystemLogLevel != "" {
			loggingArgs.SystemLogLevel = pulumi.String(spec.LoggingConfig.SystemLogLevel)
		}
		if spec.LoggingConfig.LogGroup.GetValue() != "" {
			loggingArgs.LogGroup = pulumi.String(spec.LoggingConfig.LogGroup.GetValue())
		}
		args.LoggingConfig = loggingArgs
	}

	createdFunction, err := lambda.NewFunction(ctx, "function", args, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Lambda function")
	}
	return createdFunction, nil
}
