package module

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/cloudwatch"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// accountPolicy creates the CloudWatch Logs account policy and exports
// outputs.
//
// Lifecycle facts the render below depends on:
//   - policy_name and policy_type are BOTH identity (changing either
//     replaces the policy; the provider imports the pair as
//     "policy_name:policy_type");
//   - selection_criteria narrows the account-wide scope and also
//     replaces on change; only the document (and AWS's scope argument)
//     update in place;
//   - the provider's `scope` argument is deliberately pinned to its
//     only legal value (ALL) rather than modeled - a recorded
//     exclusion.
func accountPolicy(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	// The document arrives as a Struct (each policy type's own JSON
	// schema); the provider wants a JSON string and diffs it
	// semantically.
	documentJson, err := json.Marshal(spec.PolicyDocument.AsMap())
	if err != nil {
		return errors.Wrap(err, "marshal policy document")
	}

	args := &cloudwatch.LogAccountPolicyArgs{
		PolicyName:     pulumi.String(spec.PolicyName),
		PolicyType:     pulumi.String(spec.PolicyType),
		PolicyDocument: pulumi.String(string(documentJson)),
		// ALL is the only value the provider's Scope enum carries at
		// the pin.
		Scope: pulumi.String("ALL"),
	}

	if spec.SelectionCriteria != "" {
		args.SelectionCriteria = pulumi.String(spec.SelectionCriteria)
	}

	createdPolicy, err := cloudwatch.NewLogAccountPolicy(ctx, "account_policy", args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create account policy")
	}

	ctx.Export(OpPolicyName, createdPolicy.PolicyName)
	ctx.Export(OpPolicyType, createdPolicy.PolicyType)
	return nil
}
