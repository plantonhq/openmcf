package module

import (
	"github.com/pkg/errors"
	awss3vectorbucketv1alpha1 "github.com/plantonhq/planton/catalog/aws/awss3vectorbucket/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// vectorBucket creates the vector bucket, its policy, and its indexes,
// and exports outputs.
//
// Lifecycle facts the render below depends on:
//   - EVERY index property is create-only (RequiresReplace) - an
//     index is replaced, not edited, so dimension must match the
//     embedding model before the first vector lands;
//   - bucket encryption is likewise fixed for life;
//   - indexes are named "index-{name}" (stable across list reorders);
//   - the index DataType argument is module-pinned to float32 - the
//     provider's enum holds exactly that one value;
//   - the policy is JSON-normalized by AWS; ForceDestroy is
//     config-only and never round-trips on import.
func vectorBucket(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	bucketArgs := &s3.VectorsVectorBucketArgs{
		VectorBucketName: pulumi.String(locals.Target.Metadata.Name),
		ForceDestroy:     pulumi.Bool(spec.ForceDestroy),
		Tags:             pulumi.ToStringMap(locals.AwsTags),
	}
	if encryption := spec.Encryption; encryption != nil {
		bucketArgs.EncryptionConfigurations = s3.VectorsVectorBucketEncryptionConfigurationArray{
			buildBucketEncryption(encryption),
		}
	}

	createdBucket, err := s3.NewVectorsVectorBucket(ctx, "vector_bucket", bucketArgs, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create vector bucket")
	}

	if spec.Policy != "" {
		if _, err := s3.NewVectorsVectorBucketPolicy(ctx, "bucket_policy", &s3.VectorsVectorBucketPolicyArgs{
			VectorBucketArn: createdBucket.VectorBucketArn,
			Policy:          pulumi.String(spec.Policy),
		}, pulumi.Provider(provider)); err != nil {
			return errors.Wrap(err, "bucket policy")
		}
	}

	indexArns := pulumi.StringMap{}
	for _, index := range spec.Indexes {
		indexArgs := &s3.VectorsIndexArgs{
			IndexName:        pulumi.String(index.Name),
			VectorBucketName: createdBucket.VectorBucketName,
			// The provider's enum holds exactly this one value.
			DataType:       pulumi.String("float32"),
			Dimension:      pulumi.Int(int(index.Dimension)),
			DistanceMetric: pulumi.String(index.DistanceMetric),
			Tags:           pulumi.ToStringMap(locals.AwsTags),
		}
		if encryption := index.Encryption; encryption != nil {
			indexArgs.EncryptionConfigurations = s3.VectorsIndexEncryptionConfigurationArray{
				buildIndexEncryption(encryption),
			}
		}
		if len(index.NonFilterableMetadataKeys) > 0 {
			indexArgs.MetadataConfiguration = &s3.VectorsIndexMetadataConfigurationArgs{
				NonFilterableMetadataKeys: pulumi.ToStringArray(index.NonFilterableMetadataKeys),
			}
		}

		createdIndex, err := s3.NewVectorsIndex(ctx, "index-"+index.Name, indexArgs, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrapf(err, "index %s", index.Name)
		}
		indexArns[index.Name] = createdIndex.IndexArn
	}

	ctx.Export(OpVectorBucketArn, createdBucket.VectorBucketArn)
	ctx.Export(OpIndexArns, indexArns)
	return nil
}

func buildBucketEncryption(encryption *awss3vectorbucketv1alpha1.AwsS3VectorsEncryption) *s3.VectorsVectorBucketEncryptionConfigurationArgs {
	args := &s3.VectorsVectorBucketEncryptionConfigurationArgs{
		SseType: pulumi.String(encryption.SseType),
	}
	if encryption.KmsKeyArn.GetValue() != "" {
		args.KmsKeyArn = pulumi.String(encryption.KmsKeyArn.GetValue())
	}
	return args
}

func buildIndexEncryption(encryption *awss3vectorbucketv1alpha1.AwsS3VectorsEncryption) *s3.VectorsIndexEncryptionConfigurationArgs {
	args := &s3.VectorsIndexEncryptionConfigurationArgs{
		SseType: pulumi.String(encryption.SseType),
	}
	if encryption.KmsKeyArn.GetValue() != "" {
		args.KmsKeyArn = pulumi.String(encryption.KmsKeyArn.GetValue())
	}
	return args
}
