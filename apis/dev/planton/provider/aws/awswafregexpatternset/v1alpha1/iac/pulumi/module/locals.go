package module

import (
	"strconv"

	awswafregexpatternsetv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awswafregexpatternset/v1alpha1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals mirrors Terraform-style locals: the target resource and the identity
// tag set (the Name tag is what the WAF console displays alongside the set).
type Locals struct {
	AwsWafRegexPatternSet *awswafregexpatternsetv1alpha1.AwsWafRegexPatternSet
	AwsTags               map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *awswafregexpatternsetv1alpha1.AwsWafRegexPatternSetStackInput) *Locals {
	locals := &Locals{}
	locals.AwsWafRegexPatternSet = stackInput.Target

	locals.AwsTags = map[string]string{
		awstagkeys.Name:         locals.AwsWafRegexPatternSet.Metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: locals.AwsWafRegexPatternSet.Metadata.Org,
		awstagkeys.Environment:  locals.AwsWafRegexPatternSet.Metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsWafRegexPatternSet.String(),
		awstagkeys.ResourceId:   locals.AwsWafRegexPatternSet.Metadata.Id,
	}

	return locals
}
