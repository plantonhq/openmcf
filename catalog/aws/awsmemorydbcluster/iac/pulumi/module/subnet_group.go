package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/memorydb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// subnetGroup creates the module-managed MemoryDB subnet group when the
// folded subnet_ids arm is used. The bring-your-own subnet_group_name arm
// (mutually exclusive by CEL) is handled in cluster.go — this module never
// creates a group it does not own.
func subnetGroup(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*memorydb.SubnetGroup, error) {
	spec := locals.Spec
	if len(spec.SubnetIds) == 0 {
		return nil, nil
	}

	// Subnet refs arrive pre-resolved to plain subnet IDs (the platform
	// flattens valueFrom references before the module runs).
	var subnetIds pulumi.StringArray
	for _, s := range spec.SubnetIds {
		if s.GetValue() != "" {
			subnetIds = append(subnetIds, pulumi.String(s.GetValue()))
		}
	}
	if len(subnetIds) == 0 {
		return nil, nil
	}

	// The group carries the cluster's own name: everything the module owns
	// is discoverable by one identifier, on both engines.
	sg, err := memorydb.NewSubnetGroup(ctx, "subnet-group", &memorydb.SubnetGroupArgs{
		Name:        pulumi.String(locals.ClusterName),
		Description: pulumi.Sprintf("MemoryDB subnet group for %s", locals.ClusterName),
		SubnetIds:   subnetIds,
		Tags:        pulumi.ToStringMap(locals.AwsTags),
	}, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "create subnet group")
	}

	// The subnet_group_name output is exported once, in cluster.go, where
	// all three arms (module-managed / bring-your-own / default) converge.
	return sg, nil
}
