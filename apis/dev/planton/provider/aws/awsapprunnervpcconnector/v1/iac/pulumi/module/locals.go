package module

import (
	"strconv"

	awsapprunnervpcconnectorv1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsapprunnervpcconnector/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals mirrors Terraform-style locals: the target resource and the identity
// tag set applied to the VPC connector.
type Locals struct {
	AwsAppRunnerVpcConnector *awsapprunnervpcconnectorv1.AwsAppRunnerVpcConnector
	AwsTags                  map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *awsapprunnervpcconnectorv1.AwsAppRunnerVpcConnectorStackInput) *Locals {
	locals := &Locals{}
	locals.AwsAppRunnerVpcConnector = stackInput.Target

	locals.AwsTags = map[string]string{
		awstagkeys.Name:         locals.AwsAppRunnerVpcConnector.Metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: locals.AwsAppRunnerVpcConnector.Metadata.Org,
		awstagkeys.Environment:  locals.AwsAppRunnerVpcConnector.Metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsAppRunnerVpcConnector.String(),
		awstagkeys.ResourceId:   locals.AwsAppRunnerVpcConnector.Metadata.Id,
	}

	return locals
}
