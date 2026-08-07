package module

import (
	"strconv"

	awsnatgatewayv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsnatgateway/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsNatGateway *awsnatgatewayv1alpha1.AwsNatGateway
	AwsTags       map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *awsnatgatewayv1alpha1.AwsNatGatewayStackInput) *Locals {
	locals := &Locals{}
	locals.AwsNatGateway = stackInput.Target

	metadata := stackInput.Target.Metadata
	locals.AwsTags = map[string]string{
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsNatGateway.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}
