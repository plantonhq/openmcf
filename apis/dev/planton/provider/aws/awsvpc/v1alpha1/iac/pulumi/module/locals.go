package module

import (
	"strconv"

	awsvpcv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsvpc/v1alpha1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsVpc  *awsvpcv1alpha1.AwsVpc
	AwsTags map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *awsvpcv1alpha1.AwsVpcStackInput) *Locals {
	locals := &Locals{}
	locals.AwsVpc = stackInput.Target

	metadata := stackInput.Target.Metadata
	locals.AwsTags = map[string]string{
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsVpc.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}
