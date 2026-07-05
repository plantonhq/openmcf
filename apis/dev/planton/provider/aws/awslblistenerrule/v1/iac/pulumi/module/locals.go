package module

import (
	"strconv"

	awslblistenerrulev1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awslblistenerrule/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsLbListenerRule *awslblistenerrulev1.AwsLbListenerRule
	AwsTags           map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *awslblistenerrulev1.AwsLbListenerRuleStackInput) *Locals {
	locals := &Locals{}
	locals.AwsLbListenerRule = stackInput.Target

	metadata := stackInput.Target.Metadata
	locals.AwsTags = map[string]string{
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsLbListenerRule.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}
