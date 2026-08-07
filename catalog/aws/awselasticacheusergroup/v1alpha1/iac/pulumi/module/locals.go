package module

import (
	"strconv"

	awselasticacheusergroupv1alpha1 "github.com/plantonhq/planton/catalog/aws/awselasticacheusergroup/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsElasticacheUserGroup *awselasticacheusergroupv1alpha1.AwsElasticacheUserGroup

	// UserGroupId is metadata.name -- the AWS user group id is create-time
	// immutable, and metadata.name is the naming basis both engines share so
	// a manifest deploys identically on either.
	UserGroupId string

	AwsTags map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *awselasticacheusergroupv1alpha1.AwsElasticacheUserGroupStackInput) *Locals {
	locals := &Locals{}
	locals.AwsElasticacheUserGroup = stackInput.Target

	metadata := stackInput.Target.Metadata
	locals.UserGroupId = metadata.Name

	// Resource-identity tags match the Terraform module key-for-key.
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsElasticacheUserGroup.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}
