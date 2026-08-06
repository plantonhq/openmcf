package module

import (
	"strconv"

	awsapprunnerservicev1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsapprunnerservice/v1alpha1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals mirrors Terraform-style locals: the target resource and the identity
// tag set applied to the service.
type Locals struct {
	AwsAppRunnerService *awsapprunnerservicev1alpha1.AwsAppRunnerService
	AwsTags             map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *awsapprunnerservicev1alpha1.AwsAppRunnerServiceStackInput) *Locals {
	locals := &Locals{}
	locals.AwsAppRunnerService = stackInput.Target

	locals.AwsTags = map[string]string{
		awstagkeys.Name:         locals.AwsAppRunnerService.Metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: locals.AwsAppRunnerService.Metadata.Org,
		awstagkeys.Environment:  locals.AwsAppRunnerService.Metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsAppRunnerService.String(),
		awstagkeys.ResourceId:   locals.AwsAppRunnerService.Metadata.Id,
	}

	return locals
}
