package module

import (
	"strconv"

	awslblistenerrulev1alpha1 "github.com/plantonhq/planton/catalog/aws/awslblistenerrule/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsLbListenerRule *awslblistenerrulev1alpha1.AwsLbListenerRule
	AwsTags           map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *awslblistenerrulev1alpha1.AwsLbListenerRuleStackInput) *Locals {
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
