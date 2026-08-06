package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/redshift"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// clusterParameterGroup manages a parameter group only when the spec
// carries inline parameters -- a named parameter list is configuration
// owned by exactly this cluster, not a composable node. The family
// defaults to redshift-1.0 (accepted on every cluster);
// parameter_group_family selects the redshift-2.0 generation when the
// group should track it.
func clusterParameterGroup(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*redshift.ParameterGroup, error) {
	spec := locals.AwsRedshiftCluster.Spec
	if len(spec.Parameters) == 0 {
		return nil, nil
	}

	parameters := redshift.ParameterGroupParameterArray{}
	for _, parameter := range spec.Parameters {
		parameters = append(parameters, &redshift.ParameterGroupParameterArgs{
			Name:  pulumi.String(parameter.Name),
			Value: pulumi.String(parameter.Value),
		})
	}

	family := "redshift-1.0"
	if spec.ParameterGroupFamily != "" {
		family = spec.ParameterGroupFamily
	}

	createdParameterGroup, err := redshift.NewParameterGroup(ctx, "parameter-group",
		&redshift.ParameterGroupArgs{
			Name:       pulumi.String(locals.ClusterIdentifier),
			Family:     pulumi.String(family),
			Parameters: parameters,
			Tags:       pulumi.ToStringMap(locals.AwsTags),
		}, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Redshift parameter group")
	}
	return createdParameterGroup, nil
}
