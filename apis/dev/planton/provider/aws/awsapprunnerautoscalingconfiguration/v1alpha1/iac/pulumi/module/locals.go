package module

import (
	"strconv"

	awsapprunnerautoscalingconfigurationv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsapprunnerautoscalingconfiguration/v1alpha1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals mirrors Terraform-style locals: the target resource and the identity
// tag set applied to the configuration.
type Locals struct {
	AwsAppRunnerAutoScalingConfiguration *awsapprunnerautoscalingconfigurationv1alpha1.AwsAppRunnerAutoScalingConfiguration
	AwsTags                              map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *awsapprunnerautoscalingconfigurationv1alpha1.AwsAppRunnerAutoScalingConfigurationStackInput) *Locals {
	locals := &Locals{}
	locals.AwsAppRunnerAutoScalingConfiguration = stackInput.Target

	locals.AwsTags = map[string]string{
		awstagkeys.Name:         locals.AwsAppRunnerAutoScalingConfiguration.Metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: locals.AwsAppRunnerAutoScalingConfiguration.Metadata.Org,
		awstagkeys.Environment:  locals.AwsAppRunnerAutoScalingConfiguration.Metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsAppRunnerAutoScalingConfiguration.String(),
		awstagkeys.ResourceId:   locals.AwsAppRunnerAutoScalingConfiguration.Metadata.Id,
	}

	return locals
}
