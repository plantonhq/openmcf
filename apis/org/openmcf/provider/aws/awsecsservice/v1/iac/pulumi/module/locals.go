package module

import (
	"strconv"

	"github.com/plantonhq/openmcf/apis/org/openmcf/shared/cloudresourcekind"

	awsecsservicev1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/aws/awsecsservice/v1"
	"github.com/plantonhq/openmcf/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals defines local variables used throughout our Pulumi module.
// This includes the AWS ECS Service resource definition (AwsEcsService)
// and a map of AWS tags to apply to resources.
type Locals struct {
	AwsEcsService *awsecsservicev1.AwsEcsService
	AwsTags       map[string]string
}

// initializeLocals pulls values from the stack input (AwsEcsServiceStackInput)
// and populates the Locals struct. Similar to Terraform "locals" concept.
func initializeLocals(ctx *pulumi.Context, stackInput *awsecsservicev1.AwsEcsServiceStackInput) *Locals {
	locals := &Locals{
		AwsEcsService: stackInput.Target,
	}

	locals.AwsTags = map[string]string{
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: locals.AwsEcsService.Metadata.Org,
		awstagkeys.Environment:  locals.AwsEcsService.Metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsEcsService.String(),
		awstagkeys.ResourceId:   locals.AwsEcsService.Metadata.Id,
	}

	return locals
}
