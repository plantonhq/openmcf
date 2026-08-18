package module

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/cloudwatch"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// resourcePolicy creates the CloudWatch Logs resource policy and
// exports outputs.
//
// Lifecycle facts the render below depends on:
//   - the spec guarantees exactly one of policy_name (account scope)
//     and resource_arn (resource scope); both are identity - changing
//     either replaces the policy;
//   - updates pass AWS's revision ID from state (optimistic
//     concurrency), so concurrent out-of-band edits fail loudly
//     instead of being overwritten;
//   - resource-scoped deletes REQUIRE the tracked revision ID.
func resourcePolicy(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	documentJson, err := json.Marshal(spec.PolicyDocument.AsMap())
	if err != nil {
		return errors.Wrap(err, "marshal policy document")
	}

	args := &cloudwatch.LogResourcePolicyArgs{
		PolicyDocument: pulumi.String(string(documentJson)),
	}

	if spec.PolicyName != "" {
		args.PolicyName = pulumi.String(spec.PolicyName)
	}
	if spec.ResourceArn.GetValue() != "" {
		args.ResourceArn = pulumi.String(spec.ResourceArn.GetValue())
	}

	createdPolicy, err := cloudwatch.NewLogResourcePolicy(ctx, "resource_policy", args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create resource policy")
	}

	ctx.Export(OpPolicyId, createdPolicy.ID())
	ctx.Export(OpPolicyScope, createdPolicy.PolicyScope)
	ctx.Export(OpRevisionId, createdPolicy.RevisionId)
	return nil
}
