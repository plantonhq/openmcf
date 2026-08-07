package module

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// memberFormatPattern mirrors GCP IAM's accepted member shapes: either
// <type>:<value> (user:, serviceAccount:, group:, domain:, principal://,
// principalSet://, ...) or one of the bare public literals.
var memberFormatPattern = regexp.MustCompile(`^(.+:.+|allUsers|allAuthenticatedUsers)$`)

// serviceAccountIdPattern is the fully-qualified resource-name shape the
// service account IAM API addresses: projects/<project>/serviceAccounts/<account>,
// where <account> is the account's email or unique numeric ID.
var serviceAccountIdPattern = regexp.MustCompile(`^projects/[^/]+/serviceAccounts/[^/]+$`)

// iamMember applies the single ADDITIVE grant ON the service account: one
// role, to one member, on one service account resource. Additive means the
// provider merges this (role, member) pair into the account's IAM policy
// without touching any other member's bindings on the same role, and destroy
// subtracts only this exact pair — grants made by other charts, teams, or
// tools are never clobbered.
//
// Every argument is immutable (ForceNew): IAM grants have no update — any
// change replaces the grant atomically, which is also how the API behaves.
func iamMember(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpServiceAccountIamMember.Spec

	// Both identifiers are validated here at deploy time (not in the proto)
	// because the values usually arrive through references resolved only at
	// deploy time.
	serviceAccountId := spec.ServiceAccountId.GetValue()
	if !serviceAccountIdPattern.MatchString(serviceAccountId) {
		return errors.Errorf("service_account_id %q is not a fully-qualified service account resource name (projects/<project>/serviceAccounts/<email>)", serviceAccountId)
	}

	member := spec.Member.GetValue()
	if strings.HasPrefix(member, "deleted:") {
		return errors.Errorf("member %q refers to a deleted principal; grants to deleted principals are not supported", member)
	}
	if !memberFormatPattern.MatchString(member) {
		return errors.Errorf("member %q is not in IAM member format (e.g. serviceAccount:<email>, principalSet://..., user:<email>, group:<email>) and is not one of allUsers / allAuthenticatedUsers", member)
	}

	// There is no project argument: the target account's project is embedded
	// in its fully-qualified resource name.
	args := &serviceaccount.IAMMemberArgs{
		ServiceAccountId: pulumi.String(serviceAccountId),
		Role:             pulumi.String(spec.Role.GetValue()),
		Member:           pulumi.String(member),
	}

	// An IAM Condition is part of the grant's identity: the same role granted
	// with and without a condition are two independent bindings in the policy.
	if spec.Condition != nil {
		conditionArgs := &serviceaccount.IAMMemberConditionArgs{
			Title:      pulumi.String(spec.Condition.Title),
			Expression: pulumi.String(spec.Condition.Expression),
		}
		if spec.Condition.Description != "" {
			conditionArgs.Description = pulumi.String(spec.Condition.Description)
		}
		args.Condition = conditionArgs
	}

	createdMember, err := serviceaccount.NewIAMMember(ctx, "iam-member", args, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to create service account IAM member")
	}

	ctx.Export(OpServiceAccountId, createdMember.ServiceAccountId)
	ctx.Export(OpRole, createdMember.Role)
	ctx.Export(OpMember, createdMember.Member)
	ctx.Export(OpEtag, createdMember.Etag)

	return nil
}
