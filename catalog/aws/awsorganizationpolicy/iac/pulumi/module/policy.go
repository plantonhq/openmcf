package module

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/organizations"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// policy creates the policy with its folded attachments and exports
// outputs.
//
// Lifecycle facts the render below depends on:
//   - the policy type is immutable (forces replacement); name,
//     content, and description update in place (content diffs are
//     JSON-equivalence-suppressed by the provider);
//   - both attachment leaves are immutable - changing a target
//     re-attaches (detach + attach);
//   - the type must be enabled on the organization
//     (AwsOrganization.enabled_policy_types) before any attachment
//     succeeds - AWS state, not validation, is the referee;
//   - the provider's skip_destroy escape hatches are deliberately not
//     modeled - destroy means detach and delete;
//   - AWS-managed policies (FullAWSAccess, ...) can never be adopted -
//     the provider refuses to import them.
func policy(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	// The spec carries the document structured; it is serialized to
	// JSON here.
	contentJson, err := json.Marshal(spec.Content.AsMap())
	if err != nil {
		return errors.Wrap(err, "marshal policy content")
	}

	args := &organizations.PolicyArgs{
		Name:    pulumi.StringPtr(spec.PolicyName),
		Content: pulumi.String(string(contentJson)),
		Tags:    pulumi.ToStringMap(locals.AwsTags),
	}

	// Empty keeps the provider default (SERVICE_CONTROL_POLICY).
	if spec.Type != "" {
		args.Type = pulumi.StringPtr(spec.Type)
	}
	if spec.Description != "" {
		args.Description = pulumi.StringPtr(spec.Description)
	}

	createdPolicy, err := organizations.NewPolicy(ctx, "policy", args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create policy")
	}

	// Each attachment binds the policy to one root, OU, or member
	// account, keyed by its resolved target (each imports as
	// "{target_id}:{policy_id}").
	for _, attachment := range spec.Attachments {
		targetId := attachment.TargetId.GetValue()
		_, err := organizations.NewPolicyAttachment(ctx, "attachment-"+targetId,
			&organizations.PolicyAttachmentArgs{
				PolicyId: createdPolicy.ID(),
				TargetId: pulumi.String(targetId),
			}, pulumi.Provider(provider), pulumi.Parent(createdPolicy))
		if err != nil {
			return errors.Wrapf(err, "attachment %s", targetId)
		}
	}

	ctx.Export(OpPolicyId, createdPolicy.ID())
	ctx.Export(OpArn, createdPolicy.Arn)
	return nil
}
