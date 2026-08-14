package module

import (
	"github.com/pkg/errors"
	awssecretsmanagersecretv1alpha1 "github.com/plantonhq/planton/catalog/aws/awssecretsmanagersecret/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources orchestrates creation of the Secrets Manager secret and its
// satellites (resource policy, managed version, rotation) and exports
// outputs.
func Resources(ctx *pulumi.Context, stackInput *awssecretsmanagersecretv1alpha1.AwsSecretsManagerSecretStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Build the AWS provider from the stack input via the shared builder, which resolves
	// the right credential mechanism (static keys, keyless web identity, or ambient chain).
	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, locals.Target.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	if err := secret(ctx, locals, provider); err != nil {
		return errors.Wrap(err, "secrets manager secret")
	}

	return nil
}
