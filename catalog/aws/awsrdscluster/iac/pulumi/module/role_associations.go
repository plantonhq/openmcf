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
// so roles attach and detach without touching the cluster or each
// other. AWS keys associations by (cluster, role) -- one association
// per role -- and the optional feature_name links the role to a
// specific engine capability (s3Import, Lambda, SageMaker, ...). The
// cluster resource's inline IamRoles argument is deliberately unused:
// it cannot carry feature names and, per the provider's own warning,
// mixing it with association resources overwrites the associations.
func roleAssociations(ctx *pulumi.Context, locals *Locals, provider *aws.Provider,
	createdCluster *rds.Cluster) error {
	for _, entry := range locals.AwsRdsCluster.Spec.IamRoles {
		args := &rds.ClusterRoleAssociationArgs{
			DbClusterIdentifier: createdCluster.ClusterIdentifier,
			RoleArn:             pulumi.String(entry.Role.GetValue()),
		}
		if entry.FeatureName != "" {
			args.FeatureName = pulumi.String(entry.FeatureName)
		}

		if _, err := rds.NewClusterRoleAssociation(ctx,
			fmt.Sprintf("role-association-%s", entry.Role.GetValue()),
			args, pulumi.Provider(provider)); err != nil {
			return errors.Wrapf(err, "failed to associate IAM role %s", entry.Role.GetValue())
		}
	}
	return nil
}
