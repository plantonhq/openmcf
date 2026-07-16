package module

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/organizations"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// memberFormatPattern mirrors GCP IAM's accepted member shapes: either
// <type>:<value> (user:, serviceAccount:, group:, domain:, principal://,
// principalSet://, ...) or one of the bare public/legacy literals.
var memberFormatPattern = regexp.MustCompile(`^(.+:.+|projectOwners|projectReaders|projectWriters|allUsers|allAuthenticatedUsers)$`)

// iamMember applies the single ADDITIVE project-level IAM grant: one role, to
// one member, on one project. Additive means the provider merges this
// (role, member) pair into the project's IAM policy without touching any other
// member's bindings on the same role, and destroy subtracts only this exact
// pair — grants made by other charts, teams, or tools are never clobbered.
//
// Every argument is immutable (ForceNew): IAM grants have no update — any
// change replaces the grant atomically, which is also how the API behaves.
func iamMember(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpProjectIamMember.Spec

	// Member format is validated here at deploy time (not in the proto)
	// because the value usually arrives through a reference resolved only
	// at deploy time.
	member := spec.Member.GetValue()
	if strings.HasPrefix(member, "deleted:") {
		return errors.Errorf("member %q refers to a deleted principal; grants to deleted principals are not supported", member)
	}
	if !memberFormatPattern.MatchString(member) {
		return errors.Errorf("member %q is not in IAM member format (e.g. serviceAccount:<email>, user:<email>, group:<email>) and is not one of allUsers / allAuthenticatedUsers", member)
	}

	args := &projects.IAMMemberArgs{
		Role:   pulumi.String(spec.Role.GetValue()),
		Member: pulumi.String(member),
	}

	// Honor the spec contract: an empty project_id falls back to the provider's
	// default project. Unlike most GCP resources, project IAM members REQUIRE an
	// explicit project argument, so the fallback is made concrete by reading the
	// provider's resolved project from the client config (the Pulumi counterpart
	// of the Terraform module's google_client_config data source).
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.String(spec.ProjectId.GetValue())
	} else {
		clientConfig, err := organizations.GetClientConfig(ctx, pulumi.Provider(gcpProvider))
		if err != nil {
			return errors.Wrap(err, "failed to read provider client config for the default project")
		}
		if clientConfig.Project == "" {
			return errors.New("project_id is empty and the provider has no default project configured")
		}
		args.Project = pulumi.String(clientConfig.Project)
	}

	// An IAM Condition is part of the grant's identity: the same role granted
	// with and without a condition are two independent bindings in the policy.
	if spec.Condition != nil {
		conditionArgs := &projects.IAMMemberConditionArgs{
			Title:      pulumi.String(spec.Condition.Title),
			Expression: pulumi.String(spec.Condition.Expression),
		}
		if spec.Condition.Description != "" {
			conditionArgs.Description = pulumi.String(spec.Condition.Description)
		}
		args.Condition = conditionArgs
	}

	createdMember, err := projects.NewIAMMember(ctx, "iam-member", args, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to create project IAM member")
	}

	ctx.Export(OpProjectId, createdMember.Project)
	ctx.Export(OpRole, createdMember.Role)
	ctx.Export(OpMember, createdMember.Member)
	ctx.Export(OpEtag, createdMember.Etag)

	return nil
}
