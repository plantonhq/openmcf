package module

import (
	"strconv"

	awssecuritygroupv1alpha1 "github.com/plantonhq/planton/catalog/aws/awssecuritygroup/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals holds the AWS Security Group resource definition from the stack input
// and a map of AWS tags to apply to resources.
type Locals struct {
	AwsSecurityGroup *awssecuritygroupv1alpha1.AwsSecurityGroup
	AwsTags          map[string]string
}

// initializeLocals is similar to Terraform "locals" usage. It reads
// values from AwsSecurityGroupStackInput to build a Locals instance.
func initializeLocals(ctx *pulumi.Context, stackInput *awssecuritygroupv1alpha1.AwsSecurityGroupStackInput) *Locals {
	locals := &Locals{}

	locals.AwsSecurityGroup = stackInput.Target

	// Resource-identity tags match the Terraform module key-for-key.
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         locals.AwsSecurityGroup.Metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: locals.AwsSecurityGroup.Metadata.Org,
		awstagkeys.Environment:  locals.AwsSecurityGroup.Metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsSecurityGroup.String(),
		awstagkeys.ResourceId:   locals.AwsSecurityGroup.Metadata.Id,
	}

	return locals
}
