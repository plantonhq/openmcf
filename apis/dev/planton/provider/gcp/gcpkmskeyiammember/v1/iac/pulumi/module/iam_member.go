package module

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/kms"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// memberFormatPattern mirrors GCP IAM's accepted member shapes: either
// <type>:<value> (user:, serviceAccount:, group:, domain:, principal://,
// principalSet://, ...) or one of the bare public literals.
var memberFormatPattern = regexp.MustCompile(`^(.+:.+|allUsers|allAuthenticatedUsers)$`)

// cryptoKeyIdPattern accepts the fully-qualified crypto key resource path
// (projects/<project>/locations/<location>/keyRings/<ring>/cryptoKeys/<key> —
// what a GcpKmsKey reference resolves to) plus the provider's two shorthand
// forms: <project>/<location>/<ring>/<key> and <location>/<ring>/<key>
// (the latter rides the provider's default project).
var cryptoKeyIdPattern = regexp.MustCompile(`^(projects/[^/]+/locations/[^/]+/keyRings/[^/]+/cryptoKeys/[^/]+|[^/]+/[^/]+/[^/]+/[^/]+|[^/]+/[^/]+/[^/]+)$`)

// iamMember applies the single ADDITIVE grant ON the crypto key: one role, to
// one member, on one key. Additive means the provider merges this
// (role, member) pair into the key's IAM policy without touching any other
// member's bindings on the same role, and destroy subtracts only this exact
// pair — grants made by other charts, teams, or tools are never clobbered.
//
// Every argument is immutable (ForceNew): IAM grants have no update — any
// change replaces the grant atomically, which is also how the API behaves.
func iamMember(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpKmsKeyIamMember.Spec

	// Both identifiers are validated here at deploy time (not in the proto)
	// because the values usually arrive through references resolved only at
	// deploy time.
	cryptoKeyId := spec.CryptoKeyId.GetValue()
	if !cryptoKeyIdPattern.MatchString(cryptoKeyId) {
		return errors.Errorf("crypto_key_id %q is not a crypto key identifier (projects/<project>/locations/<location>/keyRings/<ring>/cryptoKeys/<key>)", cryptoKeyId)
	}

	member := spec.Member.GetValue()
	if strings.HasPrefix(member, "deleted:") {
		return errors.Errorf("member %q refers to a deleted principal; grants to deleted principals are not supported", member)
	}
	if !memberFormatPattern.MatchString(member) {
		return errors.Errorf("member %q is not in IAM member format (e.g. serviceAccount:<email>, user:<email>, group:<email>) and is not one of allUsers / allAuthenticatedUsers", member)
	}

	// There is no project or location argument: both are embedded in the
	// key's resource path.
	args := &kms.CryptoKeyIAMMemberArgs{
		CryptoKeyId: pulumi.String(cryptoKeyId),
		Role:        pulumi.String(spec.Role.GetValue()),
		Member:      pulumi.String(member),
	}

	// An IAM Condition is part of the grant's identity: the same role granted
	// with and without a condition are two independent bindings in the policy.
	if spec.Condition != nil {
		conditionArgs := &kms.CryptoKeyIAMMemberConditionArgs{
			Title:      pulumi.String(spec.Condition.Title),
			Expression: pulumi.String(spec.Condition.Expression),
		}
		if spec.Condition.Description != "" {
			conditionArgs.Description = pulumi.String(spec.Condition.Description)
		}
		args.Condition = conditionArgs
	}

	createdMember, err := kms.NewCryptoKeyIAMMember(ctx, "iam-member", args, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to create kms crypto key IAM member")
	}

	// crypto_key_id echoes the configured identifier on both engines (the
	// provider normalizes only on import), so the output is byte-identical
	// to the Terraform module's.
	ctx.Export(OpCryptoKeyId, createdMember.CryptoKeyId)
	ctx.Export(OpRole, createdMember.Role)
	ctx.Export(OpMember, createdMember.Member)
	ctx.Export(OpEtag, createdMember.Etag)

	return nil
}
