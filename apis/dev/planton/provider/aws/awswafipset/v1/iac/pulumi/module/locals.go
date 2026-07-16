package module

import (
	"strconv"

	awswafipsetv1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awswafipset/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals mirrors Terraform-style locals: the target resource and the identity
// tag set (the Name tag is what the WAF console displays alongside the set).
type Locals struct {
	AwsWafIpSet *awswafipsetv1.AwsWafIpSet
	AwsTags     map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *awswafipsetv1.AwsWafIpSetStackInput) *Locals {
	locals := &Locals{}
	locals.AwsWafIpSet = stackInput.Target

	locals.AwsTags = map[string]string{
		awstagkeys.Name:         locals.AwsWafIpSet.Metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: locals.AwsWafIpSet.Metadata.Org,
		awstagkeys.Environment:  locals.AwsWafIpSet.Metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsWafIpSet.String(),
		awstagkeys.ResourceId:   locals.AwsWafIpSet.Metadata.Id,
	}

	return locals
}
