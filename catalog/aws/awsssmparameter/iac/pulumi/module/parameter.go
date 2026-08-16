package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ssm"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// parameter creates the SSM parameter and exports outputs.
//
// Lifecycle facts the render below depends on:
//   - the parameter's name is spec.parameter_name (an explicit field -
//     hierarchical names carry slashes metadata.name cannot), and
//     changing it forces replacement;
//   - the spec's plain `value` arm renders as InsecureValue (readable
//     in previews - the arm's whole point) and the secret
//     `secure_value` arm renders as the provider's sensitive Value
//     argument; the spec guarantees exactly one arm is set;
//   - Overwrite renders only when TRUE: the provider's unset behavior
//     (fail on a pre-existing foreign name at create, overwrite own
//     updates) is the safe default, and an explicit false would break
//     the provider's own update path;
//   - an Advanced -> Standard tier downgrade forces replacement (AWS
//     forbids it in place), and Intelligent-Tiering is resolved
//     server-side to Standard or Advanced per write.
func parameter(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	args := &ssm.ParameterArgs{
		Name: pulumi.String(spec.ParameterName),
		Type: pulumi.String(spec.Type),
		Tags: pulumi.ToStringMap(locals.AwsTags),
	}

	if spec.SecureValue != "" {
		args.Value = pulumi.String(spec.SecureValue)
	}
	if spec.Value != "" {
		args.InsecureValue = pulumi.String(spec.Value)
	}
	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}
	if spec.AllowedPattern != "" {
		args.AllowedPattern = pulumi.String(spec.AllowedPattern)
	}
	if spec.Tier != "" {
		args.Tier = pulumi.String(spec.Tier)
	}
	if spec.KeyId.GetValue() != "" {
		args.KeyId = pulumi.String(spec.KeyId.GetValue())
	}
	if spec.DataType != "" {
		args.DataType = pulumi.String(spec.DataType)
	}
	if spec.Overwrite {
		args.Overwrite = pulumi.Bool(true)
	}

	createdParameter, err := ssm.NewParameter(ctx, "parameter", args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create parameter")
	}

	ctx.Export(OpParameterName, createdParameter.Name)
	ctx.Export(OpParameterArn, createdParameter.Arn)
	// The version is an int at the provider; exported as a string to
	// match the Terraform module's tostring() output key-for-key.
	ctx.Export(OpVersion, pulumi.Sprintf("%d", createdParameter.Version))
	ctx.Export(OpTier, createdParameter.Tier)
	return nil
}
