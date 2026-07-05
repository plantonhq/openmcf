package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/elasticache"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// subnetGroup manages the cluster's ElastiCache subnet group when the spec
// brings raw subnets; an existing group name short-circuits it. The group
// itself is pure glue (a named list of subnets), which is why it stays
// inside this module instead of being its own node — the referenced subnets
// are first-class AwsSubnet nodes this module never modifies.
func subnetGroup(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*elasticache.SubnetGroup, error) {
	spec := locals.Spec
	if spec.SubnetGroupName != "" || len(spec.SubnetIds) == 0 {
		return nil, nil
	}

	subnetIds := pulumi.StringArray{}
	for _, subnetId := range spec.SubnetIds {
		if subnetId.GetValue() != "" {
			subnetIds = append(subnetIds, pulumi.String(subnetId.GetValue()))
		}
	}
	if len(subnetIds) == 0 {
		return nil, nil
	}

	createdSubnetGroup, err := elasticache.NewSubnetGroup(ctx, "subnet-group",
		&elasticache.SubnetGroupArgs{
			Name:        pulumi.String(locals.ClusterIdentifier),
			Description: pulumi.Sprintf("ElastiCache subnet group for %s", locals.ClusterIdentifier),
			SubnetIds:   subnetIds,
			Tags:        pulumi.ToStringMap(locals.AwsTags),
		}, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "create subnet group")
	}

	ctx.Export(OpSubnetGroupName, createdSubnetGroup.Name)
	return createdSubnetGroup, nil
}
