package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/rds"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// roleAssociations attaches the spec's engine feature roles: one
// association resource per spec.iam_roles entry, keyed by the role ARN
// so roles attach and detach without touching the instance. AWS keys
// associations by (instance, role) and REQUIRES the feature name here
// (unlike the cluster-side association, where it is optional) -- the
// spec's CEL mirrors that asymmetry.
func roleAssociations(ctx *pulumi.Context, locals *Locals, provider *aws.Provider,
	createdInstance *rds.Instance) error {
	for _, entry := range locals.AwsRdsInstance.Spec.IamRoles {
		if _, err := rds.NewRoleAssociation(ctx,
			fmt.Sprintf("role-association-%s", entry.Role.GetValue()),
			&rds.RoleAssociationArgs{
				DbInstanceIdentifier: createdInstance.Identifier,
				RoleArn:              pulumi.String(entry.Role.GetValue()),
				FeatureName:          pulumi.String(entry.FeatureName),
			}, pulumi.Provider(provider)); err != nil {
			return errors.Wrapf(err, "failed to associate IAM role %s", entry.Role.GetValue())
		}
	}
	return nil
}
