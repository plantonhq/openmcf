package module

import (
	"github.com/pkg/errors"
	awsbatchjobdefinitionv1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsbatchjobdefinition/v1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources registers the AWS Batch job definition and exports its outputs.
func Resources(ctx *pulumi.Context, stackInput *awsbatchjobdefinitionv1.AwsBatchJobDefinitionStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Build the AWS provider from the stack input via the shared builder, which resolves
	// the right credential mechanism (static keys, keyless web identity, or ambient chain).
	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, locals.AwsBatchJobDefinition.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	createdJobDefinition, err := jobDefinition(ctx, locals, provider)
	if err != nil {
		return errors.Wrap(err, "job definition")
	}

	// --- Exports ---
	// The revision-carrying ARN is the primary handle: it changes on every
	// registered revision, which is what rolls referencing consumers.
	ctx.Export(OpJobDefinitionArn, createdJobDefinition.Arn)
	ctx.Export(OpArnWithoutRevision, createdJobDefinition.ArnPrefix)
	ctx.Export(OpJobDefinitionName, createdJobDefinition.Name)
	ctx.Export(OpRevision, createdJobDefinition.Revision)

	return nil
}
