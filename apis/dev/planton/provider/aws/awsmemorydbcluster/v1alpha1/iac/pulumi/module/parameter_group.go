package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/memorydb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// parameterGroup creates the module-managed MemoryDB parameter group when
// the folded parameters arm is used (CEL guarantees the family accompanies
// it). The bring-your-own parameter_group_name arm is handled in cluster.go.
func parameterGroup(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*memorydb.ParameterGroup, error) {
	spec := locals.Spec
	if len(spec.Parameters) == 0 {
		return nil, nil
	}

	var params memorydb.ParameterGroupParameterArray
	for _, p := range spec.Parameters {
		params = append(params, &memorydb.ParameterGroupParameterArgs{
			Name:  pulumi.String(p.Name),
			Value: pulumi.String(p.Value),
		})
	}

	// "-params" suffix keeps the group distinct from the cluster while
	// remaining discoverable by the cluster's name, on both engines.
	pg, err := memorydb.NewParameterGroup(ctx, "parameter-group", &memorydb.ParameterGroupArgs{
		Name:        pulumi.Sprintf("%s-params", locals.ClusterName),
		Family:      pulumi.String(spec.ParameterGroupFamily),
		Description: pulumi.Sprintf("Custom parameter group for %s", locals.ClusterName),
		Parameters:  params,
		Tags:        pulumi.ToStringMap(locals.AwsTags),
	}, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "create parameter group")
	}

	// The parameter_group_name output is exported once, in cluster.go, where
	// all three arms (module-managed / bring-your-own / default) converge.
	return pg, nil
}
