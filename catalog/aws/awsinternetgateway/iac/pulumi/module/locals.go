package module

import (
	"strconv"

	awsinternetgatewayv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsinternetgateway/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsInternetGateway *awsinternetgatewayv1alpha1.AwsInternetGateway
	AwsTags            map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *awsinternetgatewayv1alpha1.AwsInternetGatewayStackInput) *Locals {
	locals := &Locals{}
	locals.AwsInternetGateway = stackInput.Target

	metadata := stackInput.Target.Metadata
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsInternetGateway.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}
