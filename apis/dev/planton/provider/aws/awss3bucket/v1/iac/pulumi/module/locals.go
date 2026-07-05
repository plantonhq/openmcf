package module

import (
	"strconv"

	awss3bucketv1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awss3bucket/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals holds pre-computed values derived from the stack input.
type Locals struct {
	Target  *awss3bucketv1.AwsS3Bucket
	Spec    *awss3bucketv1.AwsS3BucketSpec
	AwsTags map[string]string
	// BucketName is the bucket's cloud name. S3 bucket names are globally
	// unique and immutable, so metadata.name IS the identity.
	BucketName string
}

func initializeLocals(ctx *pulumi.Context, in *awss3bucketv1.AwsS3BucketStackInput) *Locals {
	locals := &Locals{}
	locals.Target = in.Target
	locals.Spec = in.Target.Spec
	locals.BucketName = in.Target.Metadata.Name

	// Resource-identity tags match the Terraform module key-for-key.
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         locals.Target.Metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: locals.Target.Metadata.Org,
		awstagkeys.Environment:  locals.Target.Metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsS3Bucket.String(),
		awstagkeys.ResourceId:   locals.Target.Metadata.Id,
	}

	return locals
}
