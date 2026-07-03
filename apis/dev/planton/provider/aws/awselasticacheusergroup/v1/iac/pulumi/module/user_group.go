package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/elasticache"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// userGroup provisions the aws_elasticache_user_group. AWS refuses the
// create unless the membership includes a user whose user NAME is "default"
// (CEL cannot prove name-data across resources, so that constraint surfaces
// at deploy time -- the spec comment and presets steer users to include one).
func userGroup(ctx *pulumi.Context, locals *Locals, provider pulumi.ProviderResource) error {
	spec := locals.AwsElasticacheUserGroup.Spec

	// Membership refs arrive pre-resolved to plain user ids (the platform
	// flattens valueFrom references before the module runs), so GetValue()
	// is the whole extraction -- literals and resolved refs look identical.
	userIds := make([]string, 0, len(spec.UserIds))
	for _, ref := range spec.UserIds {
		userIds = append(userIds, ref.GetValue())
	}

	created, err := elasticache.NewUserGroup(ctx, locals.UserGroupId, &elasticache.UserGroupArgs{
		// The AWS user group id is create-time immutable and doubles as
		// the Pulumi resource name -- metadata.name on both engines.
		UserGroupId: pulumi.String(locals.UserGroupId),
		Engine:      pulumi.String(spec.Engine),
		UserIds:     pulumi.ToStringArray(userIds),
		Tags:        pulumi.ToStringMap(locals.AwsTags),
	}, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "failed to create user group")
	}

	ctx.Export(OpUserGroupId, created.UserGroupId)
	ctx.Export(OpArn, created.Arn)

	return nil
}
