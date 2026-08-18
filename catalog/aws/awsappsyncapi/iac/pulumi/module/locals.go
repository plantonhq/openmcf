package module

import (
	"strconv"

	awsappsyncapiv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsappsyncapi/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals holds pre-computed values derived from the stack input.
type Locals struct {
	Target *awsappsyncapiv1alpha1.AwsAppSyncApi
	Spec   *awsappsyncapiv1alpha1.AwsAppSyncApiSpec

	// IsMerged: a MERGED API is declared by the merged block's
	// presence; the provider's api_type argument is derived from it.
	IsMerged bool

	AwsTags map[string]string
}

func initializeLocals(_ *pulumi.Context, in *awsappsyncapiv1alpha1.AwsAppSyncApiStackInput) *Locals {
	locals := &Locals{}
	locals.Target = in.Target
	locals.Spec = in.Target.Spec
	locals.IsMerged = locals.Spec.GetGraphql().GetMerged() != nil

	metadata := in.Target.Metadata

	// Resource-identity tags match the Terraform module key-for-key.
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsAppSyncApi.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}
