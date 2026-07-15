package module

import (
	"strconv"

	awsrediselasticachev1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsrediselasticache/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals holds pre-computed values derived from the stack input.
type Locals struct {
	AwsRedisElasticache *awsrediselasticachev1.AwsRedisElasticache

	// ReplicationGroupId is metadata.name -- create-only in AWS, and the
	// basis both engines share so a manifest deploys identically on either.
	ReplicationGroupId string

	AwsTags map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *awsrediselasticachev1.AwsRedisElasticacheStackInput) *Locals {
	locals := &Locals{}
	locals.AwsRedisElasticache = stackInput.Target

	metadata := stackInput.Target.Metadata
	locals.ReplicationGroupId = metadata.Name

	// Resource-identity tags match the Terraform module key-for-key.
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsRedisElasticache.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}
