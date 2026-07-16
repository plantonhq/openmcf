package module

import (
	"strconv"

	awscloudwatchcompositealarmv1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awscloudwatchcompositealarm/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals mirrors Terraform-style locals: the target resource and the identity
// tag set applied to the composite alarm.
type Locals struct {
	AwsCloudwatchCompositeAlarm *awscloudwatchcompositealarmv1.AwsCloudwatchCompositeAlarm
	AwsTags                     map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *awscloudwatchcompositealarmv1.AwsCloudwatchCompositeAlarmStackInput) *Locals {
	locals := &Locals{}
	locals.AwsCloudwatchCompositeAlarm = stackInput.Target

	locals.AwsTags = map[string]string{
		awstagkeys.Name:         locals.AwsCloudwatchCompositeAlarm.Metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: locals.AwsCloudwatchCompositeAlarm.Metadata.Org,
		awstagkeys.Environment:  locals.AwsCloudwatchCompositeAlarm.Metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsCloudwatchCompositeAlarm.String(),
		awstagkeys.ResourceId:   locals.AwsCloudwatchCompositeAlarm.Metadata.Id,
	}

	return locals
}
