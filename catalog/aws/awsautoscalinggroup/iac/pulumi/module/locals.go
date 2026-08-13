package module

import (
	"strconv"

	awsautoscalinggroupv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsautoscalinggroup/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsAutoScalingGroup *awsautoscalinggroupv1alpha1.AwsAutoScalingGroup
	AwsTags             map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *awsautoscalinggroupv1alpha1.AwsAutoScalingGroupStackInput) *Locals {
	locals := &Locals{}
	locals.AwsAutoScalingGroup = stackInput.Target

	metadata := stackInput.Target.Metadata
	// Resource-identity tags match the Terraform module key-for-key. On an
	// auto-scaling group these are emitted through the native tag blocks with
	// propagate_at_launch enabled (see groupTags), so every launched instance
	// carries them.
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsAutoScalingGroup.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}
