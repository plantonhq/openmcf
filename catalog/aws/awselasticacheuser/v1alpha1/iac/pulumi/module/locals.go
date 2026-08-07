package module

import (
	"strconv"

	awselasticacheuserv1alpha1 "github.com/plantonhq/planton/catalog/aws/awselasticacheuser/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsElasticacheUser *awselasticacheuserv1alpha1.AwsElasticacheUser

	// UserId is metadata.name -- the AWS user id is create-time immutable,
	// and metadata.name is the naming basis both engines share so a manifest
	// deploys identically on either.
	UserId string

	AwsTags map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *awselasticacheuserv1alpha1.AwsElasticacheUserStackInput) *Locals {
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
