package module

import (
	"strconv"

	awsapprunnerobservabilityconfigurationv1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsapprunnerobservabilityconfiguration/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals mirrors Terraform-style locals: the target resource and the identity
// tag set applied to the configuration.
type Locals struct {
	AwsAppRunnerObservabilityConfiguration *awsapprunnerobservabilityconfigurationv1.AwsAppRunnerObservabilityConfiguration
	AwsTags                                map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *awsapprunnerobservabilityconfigurationv1.AwsAppRunnerObservabilityConfigurationStackInput) *Locals {
	locals := &Locals{}
	locals.AwsAppRunnerObservabilityConfiguration = stackInput.Target

	locals.AwsTags = map[string]string{
		awstagkeys.Name:         locals.AwsAppRunnerObservabilityConfiguration.Metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: locals.AwsAppRunnerObservabilityConfiguration.Metadata.Org,
		awstagkeys.Environment:  locals.AwsAppRunnerObservabilityConfiguration.Metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsAppRunnerObservabilityConfiguration.String(),
		awstagkeys.ResourceId:   locals.AwsAppRunnerObservabilityConfiguration.Metadata.Id,
	}

	return locals
}
