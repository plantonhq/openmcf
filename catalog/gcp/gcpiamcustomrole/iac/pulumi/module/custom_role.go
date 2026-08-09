package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// customRole provisions the project-scoped IAM custom role — a named,
// least-privilege permission bundle.
//
// role_id and project are immutable (ForceNew in the provider): changing
// either destroys and recreates the role, breaking every grant that references
// the old projects/<project>/roles/<role_id> name. title, description, stage,
// and permissions all update in place, and permission edits propagate
// immediately to every existing grant of the role.
//
// GCP soft-deletes custom roles: after destroy, the role_id stays reserved for
// up to 14 days. Re-creating a role with a soft-deleted ID undeletes it and
// patches it to this configuration — the provider handles that flow natively,
// so a destroy/recreate cycle within the window converges rather than failing.
func customRole(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpIamCustomRole.Spec

	// Drop accidental empty strings so the API never sees a blank permission.
	permissions := pulumi.StringArray{}
	for _, permission := range spec.Permissions {
		if permission != "" {
			permissions = append(permissions, pulumi.String(permission))
		}
	}

	args := &projects.IAMCustomRoleArgs{
		RoleId: pulumi.String(spec.RoleId),
		Title:  pulumi.String(spec.Title),
		// Launch stage defaults to GA (via the spec default) — the right label
		// for production roles.
		Stage:       pulumi.String(spec.GetStage()),
		Permissions: permissions,
	}

	// DELETE (provider default) soft-deletes the role on destroy; PREVENT
	// fails the destroy; ABANDON leaves the role active. Sent only when
	// set — mirrors the Terraform module.
	if spec.DeletionPolicy != "" {
		args.DeletionPolicy = pulumi.StringPtr(spec.DeletionPolicy)
	}

	// Omitted description stays unset (matching the Terraform module's null)
	// rather than being sent as an empty string.
	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}

	// Honor the spec contract: an empty project_id falls back to the provider's
	// default project. Leaving Project unset lets the gcp provider resolve its
	// own project (configuration or the GOOGLE_PROJECT / GOOGLE_CLOUD_PROJECT
	// environment chain); an empty string would be sent verbatim and rejected.
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.String(spec.ProjectId.GetValue())
	}

	createdRole, err := projects.NewIAMCustomRole(ctx, "custom-role", args, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to create custom role")
	}

	ctx.Export(OpName, createdRole.Name)
	ctx.Export(OpRoleId, createdRole.RoleId)
	ctx.Export(OpDeleted, createdRole.Deleted)

	return nil
}
