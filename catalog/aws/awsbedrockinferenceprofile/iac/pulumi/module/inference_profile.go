package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/bedrock"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// inferenceProfile creates the Bedrock application inference profile and
// exports outputs. Every argument is create-time-immutable -- changing one
// replaces the profile and its ARN.
func inferenceProfile(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	args := &bedrock.InferenceProfileArgs{
		Name: pulumi.String(locals.ProfileName),
		// The model source this profile routes to. AWS never echoes it
		// back (GetInferenceProfile reports the resolved models list
		// instead), so the provider pins the configured value in state.
		ModelSource: &bedrock.InferenceProfileModelSourceArgs{
			CopyFrom: pulumi.String(spec.SourceArn),
		},
		Tags: pulumi.ToStringMap(locals.AwsTags),
	}

	// Sent only when set (1-200 characters; ForceNew at the provider).
	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}

	createdProfile, err := bedrock.NewInferenceProfile(ctx, locals.ProfileName, args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create inference profile")
	}

	ctx.Export(OpInferenceProfileArn, createdProfile.Arn)
	ctx.Export(OpInferenceProfileId, createdProfile.ID())
	ctx.Export(OpStatus, createdProfile.Status)
	ctx.Export(OpType, createdProfile.Type)

	return nil
}
