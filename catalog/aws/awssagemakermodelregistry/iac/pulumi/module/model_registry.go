package module

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/sagemaker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// modelRegistry creates the model package group and its folded resource
// policy and exports outputs.
//
// Lifecycle facts the renders below depend on:
//   - the group is immutable except tags (even the description is
//     ForceNew upstream - a description edit replaces the group);
//   - the policy is an idempotent upsert (PutModelPackageGroupPolicy)
//     that updates in place; removing spec.resource_policy deletes the
//     policy resource;
//   - model package VERSIONS register into the group imperatively
//     (training pipelines, SDK) - never through this module.
func modelRegistry(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	args := &sagemaker.ModelPackageGroupArgs{
		// The component's name IS the group name.
		ModelPackageGroupName: pulumi.String(locals.GroupName),
		Tags:                  pulumi.ToStringMap(locals.AwsTags),
	}
	if spec.Description != "" {
		args.ModelPackageGroupDescription = pulumi.String(spec.Description)
	}

	createdGroup, err := sagemaker.NewModelPackageGroup(ctx, locals.GroupName, args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create model package group")
	}

	// Cross-account sharing policy - present exactly when the spec
	// carries one.
	if spec.ResourcePolicy != nil {
		policyBytes, err := json.Marshal(spec.ResourcePolicy.AsMap())
		if err != nil {
			return errors.Wrap(err, "marshal resource policy")
		}
		if _, err := sagemaker.NewModelPackageGroupPolicy(ctx, locals.GroupName+"-policy",
			&sagemaker.ModelPackageGroupPolicyArgs{
				ModelPackageGroupName: createdGroup.ModelPackageGroupName,
				ResourcePolicy:        pulumi.String(string(policyBytes)),
			}, pulumi.Provider(provider), pulumi.DependsOn([]pulumi.Resource{createdGroup})); err != nil {
			return errors.Wrap(err, "create model package group policy")
		}
	}

	ctx.Export(OpModelPackageGroupName, createdGroup.ModelPackageGroupName)
	ctx.Export(OpModelPackageGroupArn, createdGroup.Arn)

	return nil
}
