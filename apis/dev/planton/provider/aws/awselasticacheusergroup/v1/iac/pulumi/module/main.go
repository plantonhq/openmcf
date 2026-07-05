package module

import (
	"github.com/pkg/errors"
	awselasticacheusergroupv1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awselasticacheusergroup/v1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources provisions the ElastiCache RBAC user group -- the attachment
// unit between users and caches. Membership is modeled here as references
// (the aws_elasticache_user_group_association glue resource is deliberately
// not used): the group is the single place an application's cache access is
// granted or revoked, and this module never mutates the users it references.
func Resources(ctx *pulumi.Context, stackInput *awselasticacheusergroupv1.AwsElasticacheUserGroupStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Build the AWS provider from the stack input via the shared builder, which
	// resolves the right credential mechanism (static keys, keyless web identity,
	// or ambient chain).
	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, locals.AwsElasticacheUserGroup.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	if err := userGroup(ctx, locals, provider); err != nil {
		return errors.Wrap(err, "failed to create ElastiCache user group")
	}

	return nil
}
