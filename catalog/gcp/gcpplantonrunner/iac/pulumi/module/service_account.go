package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// serviceAccountIdBudget is GCP's hard cap on a service account's
// account_id (6-30 characters).
const serviceAccountIdBudget = 30

// runtimeServiceAccount resolves the RUNTIME identity: the service account
// the runner runs as — the seam keyless cloud access rides. When the spec
// references a service_account, that account is the runner's identity and
// this module never touches it (modules never mutate resources they merely
// reference). Otherwise a dedicated permissionless account is created so
// the seam always exists: permissions can be granted later without
// replacing the runner. Deliberately NEVER the project's Compute Engine
// default service account — it typically carries broad project access the
// runner should not inherit silently.
//
// Returns the resolved email plus the created account (nil on the
// referenced arm) so the service can order itself after the creation.
func runtimeServiceAccount(ctx *pulumi.Context, locals *Locals, provider *gcp.Provider) (pulumi.StringOutput, pulumi.Resource, error) {
	spec := locals.GcpPlantonRunner.Spec
	runnerName := locals.GcpPlantonRunner.Metadata.Name

	if spec.ServiceAccount.GetValue() != "" {
		return pulumi.String(spec.ServiceAccount.GetValue()).ToStringOutput(), nil, nil
	}

	// FAIL LOUDLY past GCP's account_id budget: a silent truncation could
	// collide two runners' identities. Longer names keep working by
	// composing a GcpServiceAccount resource and referencing it.
	if len(runnerName) > serviceAccountIdBudget {
		return pulumi.StringOutput{}, nil, errors.Errorf(
			"metadata.name %q is %d characters — GCP caps service account ids at %d; reference your own service account via spec.service_account (or use a shorter name)",
			runnerName, len(runnerName), serviceAccountIdBudget)
	}

	accountArgs := &serviceaccount.AccountArgs{
		AccountId:   pulumi.String(runnerName),
		DisplayName: pulumi.String(fmt.Sprintf("Planton runner '%s'", runnerName)),
		Description: pulumi.String(fmt.Sprintf("Runtime identity for Planton runner '%s' -- grant it the roles keyless operations need", runnerName)),
	}
	// An empty project falls back to the provider's default project — the
	// ambient-project contract every GCP kind honors.
	if locals.ProjectId != "" {
		accountArgs.Project = pulumi.String(locals.ProjectId)
	}

	createdAccount, err := serviceaccount.NewAccount(ctx,
		"runtime-service-account",
		accountArgs,
		pulumi.Provider(provider))
	if err != nil {
		return pulumi.StringOutput{}, nil, errors.Wrap(err, "failed to create runtime service account")
	}

	return createdAccount.Email, createdAccount, nil
}
