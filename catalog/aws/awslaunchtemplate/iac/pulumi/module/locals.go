package module

import (
	"strconv"

	awslaunchtemplatev1alpha1 "github.com/plantonhq/planton/catalog/aws/awslaunchtemplate/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsLaunchTemplate *awslaunchtemplatev1alpha1.AwsLaunchTemplate
	AwsTags           map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *awslaunchtemplatev1alpha1.AwsLaunchTemplateStackInput) *Locals {
	locals := &Locals{}
	locals.AwsLaunchTemplate = stackInput.Target

	metadata := stackInput.Target.Metadata
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsLaunchTemplate.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}
