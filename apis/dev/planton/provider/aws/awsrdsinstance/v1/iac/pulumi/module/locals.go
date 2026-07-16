package module

import (
	"strconv"

	awsrdsinstancev1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsrdsinstance/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsRdsInstance *awsrdsinstancev1.AwsRdsInstance

	// InstanceIdentifier is metadata.name -- the basis both engines share
	// so a manifest deploys identically on either.
	InstanceIdentifier string

	AwsTags map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *awsrdsinstancev1.AwsRdsInstanceStackInput) *Locals {
	locals := &Locals{}
	locals.AwsRdsInstance = stackInput.Target

	metadata := stackInput.Target.Metadata
	locals.InstanceIdentifier = metadata.Name

	// Resource-identity tags match the Terraform module key-for-key.
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsRdsInstance.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}
