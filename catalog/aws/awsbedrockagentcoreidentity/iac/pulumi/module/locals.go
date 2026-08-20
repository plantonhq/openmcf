package module

import (
	"strconv"

	awsbedrockagentcoreidentityv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsbedrockagentcoreidentity/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals holds pre-computed values derived from the stack input.
type Locals struct {
	Target *awsbedrockagentcoreidentityv1alpha1.AwsBedrockAgentCoreIdentity
	Spec   *awsbedrockagentcoreidentityv1alpha1.AwsBedrockAgentCoreIdentitySpec

	AwsTags map[string]string
}

func initializeLocals(_ *pulumi.Context, in *awsbedrockagentcoreidentityv1alpha1.AwsBedrockAgentCoreIdentityStackInput) *Locals {
	locals := &Locals{}
	locals.Target = in.Target
	locals.Spec = in.Target.Spec

	metadata := in.Target.Metadata

	// Resource-identity tags match the Terraform module key-for-key.
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsBedrockAgentCoreIdentity.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}
