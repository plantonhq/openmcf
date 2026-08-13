package module

import (
	"strconv"

	awsegressonlyinternetgatewayv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsegressonlyinternetgateway/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsEgressOnlyInternetGateway *awsegressonlyinternetgatewayv1alpha1.AwsEgressOnlyInternetGateway
	AwsTags                      map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *awsegressonlyinternetgatewayv1alpha1.AwsEgressOnlyInternetGatewayStackInput) *Locals {
	locals := &Locals{}
	locals.AwsEgressOnlyInternetGateway = stackInput.Target

	metadata := stackInput.Target.Metadata
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsEgressOnlyInternetGateway.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}
