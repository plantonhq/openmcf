package module

import (
	"strconv"

	awseksaddonv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awseksaddon/v1alpha1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsEksAddon *awseksaddonv1alpha1.AwsEksAddon

	AwsTags map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *awseksaddonv1alpha1.AwsEksAddonStackInput) *Locals {
	locals := &Locals{}
	locals.AwsEksAddon = stackInput.Target

	metadata := stackInput.Target.Metadata
	locals.AwsTags = map[string]string{
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsEksAddon.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}
