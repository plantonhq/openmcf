package module

import (
	"strconv"

	awsconfigrulev1alpha1 "github.com/plantonhq/planton/catalog/aws/awsconfigrule/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals holds pre-computed values derived from the stack input.
type Locals struct {
	Target *awsconfigrulev1alpha1.AwsConfigRule
	Spec   *awsconfigrulev1alpha1.AwsConfigRuleSpec

	AwsTags map[string]string
}

func initializeLocals(_ *pulumi.Context, in *awsconfigrulev1alpha1.AwsConfigRuleStackInput) *Locals {
	locals := &Locals{}
	locals.Target = in.Target
	locals.Spec = in.Target.Spec

	metadata := in.Target.Metadata

	// Resource-identity tags match the Terraform module key-for-key.
	// Only the account-scoped rule resource is taggable - the
	// organization rule resources carry no tags in the provider.
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsConfigRule.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}
