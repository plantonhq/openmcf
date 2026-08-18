package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// directoryBucket creates the S3 directory bucket and exports
// outputs.
//
// Lifecycle facts the render below depends on:
//   - the bucket name is DERIVED (locals.BucketName) - AWS mandates
//     "{base}--{zone_id}--x-s3" and a hand-assembled name could
//     disagree with the location block;
//   - everything except ForceDestroy replaces the bucket;
//   - DataRedundancy is sent only when set; the provider derives the
//     only-valid value from the location type otherwise;
//   - ForceDestroy is config-only at AWS (never read back), so
//     imports do not round-trip it.
func directoryBucket(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	zoneType := spec.ZoneType
	if zoneType == "" {
		zoneType = "AvailabilityZone"
	}

	args := &s3.DirectoryBucketArgs{
		Bucket: pulumi.String(locals.BucketName),
		Location: &s3.DirectoryBucketLocationArgs{
			Name: pulumi.String(spec.ZoneId),
			Type: pulumi.String(zoneType),
		},
		ForceDestroy: pulumi.Bool(spec.ForceDestroy),
		Tags:         pulumi.ToStringMap(locals.AwsTags),
	}
	if spec.DataRedundancy != "" {
		args.DataRedundancy = pulumi.String(spec.DataRedundancy)
	}

	createdBucket, err := s3.NewDirectoryBucket(ctx, "directory_bucket", args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create directory bucket")
	}

	ctx.Export(OpBucketName, createdBucket.Bucket)
	ctx.Export(OpBucketArn, createdBucket.Arn)
	return nil
}
