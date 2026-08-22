package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/secretmanager"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// tokenSecret stores the runner token in Secret Manager and grants exactly
// the runner's own service account read access to exactly this one secret.
// The container never sees the token in its launch configuration: the
// service's env carries only a secret reference, and Cloud Run resolves the
// value at instance start through the runtime service account — so reading
// the service definition (a common, low-sensitivity permission) reveals
// nothing.
//
// The token authorizes JOINING and is never the runner's identity: the
// runner presents it once per enrollment, registers itself, and persists
// the identity it receives. The env reference pins Secret Manager's
// "latest" version deliberately, so rotating the token needs no service
// update — running instances keep serving on their minted identity (the
// token is only read at join), and the next instance start joins with the
// new value.
func tokenSecret(ctx *pulumi.Context, locals *Locals, provider *gcp.Provider,
	serviceAccountEmail pulumi.StringOutput,
) (pulumi.Resource, pulumi.Resource, error) {
	secretArgs := &secretmanager.SecretArgs{
		SecretId: pulumi.String(locals.TokenSecretId),
		// Automatic replication: the token is a single small value read at
		// instance start; regional placement decisions belong to the
		// service, not its bootstrap secret.
		Replication: &secretmanager.SecretReplicationArgs{
			Auto: &secretmanager.SecretReplicationAutoArgs{},
		},
		Labels: pulumi.ToStringMap(locals.GcpLabels),
	}
	if locals.ProjectId != "" {
		secretArgs.Project = pulumi.String(locals.ProjectId)
	}

	createdSecret, err := secretmanager.NewSecret(ctx,
		"token-secret",
		secretArgs,
		pulumi.Provider(provider))
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create secret")
	}

	createdSecretVersion, err := secretmanager.NewSecretVersion(ctx,
		"token-secret-value",
		&secretmanager.SecretVersionArgs{
			Secret:     createdSecret.ID(),
			SecretData: pulumi.String(locals.GcpPlantonRunner.Spec.GetToken()),
		},
		pulumi.Provider(provider),
		pulumi.Parent(createdSecret))
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to store secret value")
	}

	// The accessor grant lives ON the module-owned secret (never on the
	// referenced service account) and names exactly one principal — the
	// least-privilege twin of the AWS execution role's one-secret inline
	// policy. Without it the first instance start fails resolving the env
	// reference, however broad the project's other grants are.
	createdAccessorGrant, err := secretmanager.NewSecretIamMember(ctx,
		"token-secret-accessor",
		&secretmanager.SecretIamMemberArgs{
			Project:  createdSecret.Project,
			SecretId: createdSecret.SecretId,
			Role:     pulumi.String("roles/secretmanager.secretAccessor"),
			Member: serviceAccountEmail.ApplyT(func(email string) string {
				return fmt.Sprintf("serviceAccount:%s", email)
			}).(pulumi.StringOutput),
		},
		pulumi.Provider(provider),
		pulumi.Parent(createdSecret))
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to grant secret access")
	}

	return createdSecretVersion, createdAccessorGrant, nil
}
