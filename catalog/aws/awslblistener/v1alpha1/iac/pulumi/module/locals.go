package module

import (
	"strconv"

	awslblistenerv1alpha1 "github.com/plantonhq/planton/catalog/aws/awslblistener/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsLbListener *awslblistenerv1alpha1.AwsLbListener
	AwsTags       map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *awslblistenerv1alpha1.AwsLbListenerStackInput) *Locals {
	locals := &Locals{}
	locals.AwsLbListener = stackInput.Target

	metadata := stackInput.Target.Metadata
	locals.AwsTags = map[string]string{
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsLbListener.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}
