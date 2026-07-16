package module

import (
	"strconv"

	awselasticacheuserv1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awselasticacheuser/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsElasticacheUser *awselasticacheuserv1.AwsElasticacheUser

	// UserId is metadata.name -- the AWS user id is create-time immutable,
	// and metadata.name is the naming basis both engines share so a manifest
	// deploys identically on either.
	UserId string

	AwsTags map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *awselasticacheuserv1.AwsElasticacheUserStackInput) *Locals {
	locals := &Locals{}
	locals.AwsElasticacheUser = stackInput.Target

	metadata := stackInput.Target.Metadata
	locals.UserId = metadata.Name

	// Resource-identity tags match the Terraform module key-for-key.
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsElasticacheUser.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}
