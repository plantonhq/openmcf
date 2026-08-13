package module

import (
	"strconv"

	awslbtargetgroupv1alpha1 "github.com/plantonhq/planton/catalog/aws/awslbtargetgroup/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsLbTargetGroup *awslbtargetgroupv1alpha1.AwsLbTargetGroup
	AwsTags          map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *awslbtargetgroupv1alpha1.AwsLbTargetGroupStackInput) *Locals {
	locals := &Locals{}
	locals.AwsLbTargetGroup = stackInput.Target

	metadata := stackInput.Target.Metadata
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsLbTargetGroup.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}
