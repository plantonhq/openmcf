package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/secretsmanager"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// credentialsSecret stores the runner's credentials document in Secrets
// Manager. The container never sees the document in its launch
// configuration: the task definition carries only this secret's ARN, and
// the ECS agent fetches the value at task start through the execution
// role -- so reading the task definition (a common, low-sensitivity
// permission) reveals nothing.
func credentialsSecret(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*secretsmanager.Secret, error) {
	runnerName := locals.AwsPlantonRunner.Metadata.Name

	// Zero recovery window: the document is re-mintable credential
	// material, not data. Secrets Manager's default 30-day soft-delete
	// would block re-creating a same-named runner for a month after a
	// destroy, with no compensating benefit for a value that can simply
	// be re-issued.
	createdSecret, err := secretsmanager.NewSecret(ctx,
		"credentials-secret",
		&secretsmanager.SecretArgs{
			Name:                 pulumi.String(fmt.Sprintf("%s-credentials", runnerName)),
			Description:          pulumi.String(fmt.Sprintf("Credentials document for Planton runner '%s'", runnerName)),
			RecoveryWindowInDays: pulumi.Int(0),
			Tags:                 pulumi.ToStringMap(locals.AwsTags),
		},
		pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create secret")
	}

	_, err = secretsmanager.NewSecretVersion(ctx,
		"credentials-secret-value",
		&secretsmanager.SecretVersionArgs{
			SecretId:     createdSecret.ID(),
			SecretString: pulumi.String(locals.AwsPlantonRunner.Spec.Credentials),
		},
		pulumi.Provider(provider),
		pulumi.Parent(createdSecret))
	if err != nil {
		return nil, errors.Wrap(err, "failed to store secret value")
	}

	return createdSecret, nil
}
