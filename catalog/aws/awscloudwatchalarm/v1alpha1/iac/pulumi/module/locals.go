package module

import (
	"strconv"

	"github.com/plantonhq/planton/shared/cloudresourcekind"

	awscloudwatchalarmv1alpha1 "github.com/plantonhq/planton/catalog/aws/awscloudwatchalarm/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsCloudwatchAlarm *awscloudwatchalarmv1alpha1.AwsCloudwatchAlarm
	AwsTags            map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *awscloudwatchalarmv1alpha1.AwsCloudwatchAlarmStackInput) *Locals {
	locals := &Locals{}
	locals.AwsCloudwatchAlarm = stackInput.Target

	locals.AwsTags = map[string]string{
		awstagkeys.Name:         locals.AwsCloudwatchAlarm.Metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: locals.AwsCloudwatchAlarm.Metadata.Org,
		awstagkeys.Environment:  locals.AwsCloudwatchAlarm.Metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsCloudwatchAlarm.String(),
		awstagkeys.ResourceId:   locals.AwsCloudwatchAlarm.Metadata.Id,
	}

	return locals
}
