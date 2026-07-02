package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/organizations"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// iam attaches the roles listed in spec.project_iam_roles and spec.org_iam_roles.
// All grants use the ADDITIVE member form: each one manages exactly one
// (role, member) pair and never clobbers other members' bindings on the same
// role — safe to compose with grants made elsewhere.
func iam(
	ctx *pulumi.Context,
	locals *Locals,
	createdServiceAccount *serviceaccount.Account,
	gcpProvider *gcp.Provider,
) error {

	// === Project-level roles ===
	for idx, role := range locals.GcpServiceAccount.Spec.ProjectIamRoles {
		bindingName := fmt.Sprintf("%s-project-%d", locals.GcpServiceAccount.Metadata.Name, idx)

		// The grant targets the account's home project; deriving it from the
		// created account (rather than re-reading spec) keeps the grant correct
		// when the project fell back to the provider default.
		createdProjectIamMember, err := projects.NewIAMMember(
			ctx,
			bindingName,
			&projects.IAMMemberArgs{
				Project: createdServiceAccount.Project,
				Role:    pulumi.String(role),
				Member:  createdServiceAccount.Member,
			},
			pulumi.Provider(gcpProvider),
			pulumi.DependsOn([]pulumi.Resource{createdServiceAccount}),
		)
		if err != nil {
			return errors.Wrapf(err, "failed to add project role %s", role)
		}
		_ = createdProjectIamMember
	}

	// === Org-level roles ===
	if len(locals.GcpServiceAccount.Spec.OrgIamRoles) > 0 {
		if locals.GcpServiceAccount.Spec.OrgId == "" {
			return errors.New("org_iam_roles specified but org_id is empty")
		}

		for idx, role := range locals.GcpServiceAccount.Spec.OrgIamRoles {
			bindingName := fmt.Sprintf("%s-org-%d", locals.GcpServiceAccount.Metadata.Name, idx)

			createdOrgIamMember, err := organizations.NewIAMMember(
				ctx,
				bindingName,
				&organizations.IAMMemberArgs{
					OrgId:  pulumi.String(locals.GcpServiceAccount.Spec.OrgId),
					Role:   pulumi.String(role),
					Member: createdServiceAccount.Member,
				},
				pulumi.Provider(gcpProvider),
				pulumi.DependsOn([]pulumi.Resource{createdServiceAccount}),
			)
			if err != nil {
				return errors.Wrapf(err, "failed to add org role %s", role)
			}
			_ = createdOrgIamMember
		}
	}

	return nil
}
