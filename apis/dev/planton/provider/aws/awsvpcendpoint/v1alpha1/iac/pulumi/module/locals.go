package module

import (
	"strconv"

	awsvpcendpointv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsvpcendpoint/v1alpha1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsVpcEndpoint *awsvpcendpointv1alpha1.AwsVpcEndpoint

	AwsTags map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *awsvpcendpointv1alpha1.AwsVpcEndpointStackInput) *Locals {
	locals := &Locals{}
	locals.AwsVpcEndpoint = stackInput.Target

	// VPC endpoints carry no name parameter in AWS -- identity lives
	// entirely in tags, so the Name tag is what the console displays.
	metadata := stackInput.Target.Metadata
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsVpcEndpoint.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}
