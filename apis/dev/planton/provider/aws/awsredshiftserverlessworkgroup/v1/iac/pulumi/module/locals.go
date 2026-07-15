package module

import (
	"strconv"

	awsredshiftserverlessworkgroupv1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsredshiftserverlessworkgroup/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsRedshiftServerlessWorkgroup *awsredshiftserverlessworkgroupv1.AwsRedshiftServerlessWorkgroup

	// WorkgroupName is metadata.name -- create-only in AWS, and the
	// basis both engines share so a manifest deploys identically on
	// either.
	WorkgroupName string

	AwsTags map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *awsredshiftserverlessworkgroupv1.AwsRedshiftServerlessWorkgroupStackInput) *Locals {
	locals := &Locals{}
	locals.AwsRedshiftServerlessWorkgroup = stackInput.Target

	metadata := stackInput.Target.Metadata
	locals.WorkgroupName = metadata.Name

	// Resource-identity tags match the Terraform module key-for-key.
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsRedshiftServerlessWorkgroup.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}
