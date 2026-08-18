package module

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/cloudwatch"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// crossAccountDestination creates the legacy cross-account subscription
// destination and its access policy.
//
// Lifecycle facts the render below depends on:
//   - the destination's access policy PERSISTS at AWS when only the
//     policy resource is destroyed (a provider no-op delete);
//     destroying the destination itself is real;
//   - the first create retries for up to two minutes while the
//     logs.amazonaws.com trust on the role propagates;
//   - the provider creates the destination first and tags it in a
//     separate call (tags on the Put call break AWS's test-message
//     delivery check).
func crossAccountDestination(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec.CrossAccountDestination

	if spec == nil {
		ctx.Export(OpCrossAccountDestinationName, pulumi.String(""))
		ctx.Export(OpCrossAccountDestinationArn, pulumi.String(""))
		return nil
	}

	createdDestination, err := cloudwatch.NewLogDestination(ctx, "cross_account_destination", &cloudwatch.LogDestinationArgs{
		Name:      pulumi.String(spec.Name),
		RoleArn:   pulumi.String(spec.RoleArn.GetValue()),
		TargetArn: pulumi.String(spec.TargetArn.GetValue()),
		Tags:      pulumi.ToStringMap(locals.AwsTags),
	}, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create log destination")
	}

	policyJson, err := json.Marshal(spec.AccessPolicy.AsMap())
	if err != nil {
		return errors.Wrap(err, "marshal access policy")
	}

	policyArgs := &cloudwatch.LogDestinationPolicyArgs{
		DestinationName: createdDestination.Name,
		AccessPolicy:    pulumi.String(string(policyJson)),
	}
	if spec.ForceUpdate != nil {
		policyArgs.ForceUpdate = pulumi.Bool(*spec.ForceUpdate)
	}

	if _, err := cloudwatch.NewLogDestinationPolicy(ctx, "cross_account_destination_policy", policyArgs, pulumi.Provider(provider)); err != nil {
		return errors.Wrap(err, "create log destination policy")
	}

	ctx.Export(OpCrossAccountDestinationName, createdDestination.Name)
	ctx.Export(OpCrossAccountDestinationArn, createdDestination.Arn)
	return nil
}
