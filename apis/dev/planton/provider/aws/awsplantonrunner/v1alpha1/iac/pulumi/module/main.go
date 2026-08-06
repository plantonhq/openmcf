package module

import (
	"github.com/pkg/errors"
	awsplantonrunnerv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsplantonrunner/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources provisions the runner appliance. The pieces, in dependency
// order: the credentials secret, the two IAM roles (setup vs runtime --
// kept separate so neither accumulates the other's permissions), the log
// group, the outbound-only security group, and finally the compute
// (cluster, task definition, service) that runs the container.
func Resources(ctx *pulumi.Context, stackInput *awsplantonrunnerv1alpha1.AwsPlantonRunnerStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Build the AWS provider from the stack input via the shared builder, which
	// resolves the right credential mechanism (static keys, keyless web identity,
	// or ambient chain).
	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, locals.AwsPlantonRunner.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	createdSecret, err := credentialsSecret(ctx, locals, provider)
	if err != nil {
		return errors.Wrap(err, "failed to create credentials secret")
	}

	createdExecutionRole, err := executionRole(ctx, locals, provider, createdSecret)
	if err != nil {
		return errors.Wrap(err, "failed to create execution role")
	}

	taskRoleArn, err := runtimeRole(ctx, locals, provider)
	if err != nil {
		return errors.Wrap(err, "failed to resolve runtime role")
	}

	createdSecurityGroup, err := securityGroup(ctx, locals, provider)
	if err != nil {
		return errors.Wrap(err, "failed to create security group")
	}

	if err := runnerCompute(ctx, locals, provider, createdSecret, createdExecutionRole, taskRoleArn, createdSecurityGroup); err != nil {
		return errors.Wrap(err, "failed to create runner compute")
	}

	return nil
}
