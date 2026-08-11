package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/cognito"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// userGroups creates one aws_cognito_user_group per spec entry. Groups are
// pool-scoped configuration with no independent AWS lifecycle; membership
// (which users are in a group) is data-plane content managed at runtime,
// never from here.
func userGroups(ctx *pulumi.Context, locals *Locals, createdPool *cognito.UserPool, provider *aws.Provider) error {
	for _, group := range locals.Spec.UserGroups {
		args := &cognito.UserGroupArgs{
			UserPoolId: createdPool.ID(),
			Name:       pulumi.String(group.Name),
		}

		if group.Description != "" {
			args.Description = pulumi.StringPtr(group.Description)
		}

		// AWS accepts precedence 0 (the strongest priority) but the provider's
		// own zero-value gating cannot send it, so 0 carries "no precedence"
		// here -- the spec documents 1 as the strongest expressible value.
		if group.Precedence > 0 {
			args.Precedence = pulumi.IntPtr(int(group.Precedence))
		}

		if group.RoleArn.GetValue() != "" {
			args.RoleArn = pulumi.StringPtr(group.RoleArn.GetValue())
		}

		// Keyed by group name (the AWS identity within the pool), parented to
		// the pool so replacements order correctly.
		if _, err := cognito.NewUserGroup(ctx,
			group.Name,
			args,
			pulumi.Provider(provider),
			pulumi.Parent(createdPool)); err != nil {
			return errors.Wrapf(err, "failed to create user group %s", group.Name)
		}
	}

	return nil
}
