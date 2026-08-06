package module

import (
	"strconv"

	awsdynamodbv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsdynamodb/v1alpha1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsDynamodb *awsdynamodbv1alpha1.AwsDynamodb

	// TableName is metadata.name -- create-only in AWS, and the basis
	// both engines share so a manifest deploys identically on either.
	TableName string

	AwsTags map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *awsdynamodbv1alpha1.AwsDynamodbStackInput) *Locals {
	locals := &Locals{}
	locals.AwsDynamodb = stackInput.Target

	metadata := stackInput.Target.Metadata
	locals.TableName = metadata.Name

	// Resource-identity tags match the Terraform module key-for-key.
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsDynamodb.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}
