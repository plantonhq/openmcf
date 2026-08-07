package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/elasticache"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// parameterGroup manages a custom ElastiCache parameter group only when the
// spec carries inline parameters — a named parameter list is configuration
// owned by exactly this cluster, not a composable node. Bringing an existing
// parameter_group_name short-circuits it (CEL-enforced mutual exclusion with
// inline parameters).
func parameterGroup(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*elasticache.ParameterGroup, error) {
	spec := locals.Spec
	if spec.ParameterGroupName != "" || len(spec.Parameters) == 0 || spec.ParameterGroupFamily == "" {
		return nil, nil
	}

	parameters := elasticache.ParameterGroupParameterArray{}
	for _, parameter := range spec.Parameters {
		parameters = append(parameters, &elasticache.ParameterGroupParameterArgs{
			Name:  pulumi.String(parameter.Name),
			Value: pulumi.String(parameter.Value),
		})
	}

	createdParameterGroup, err := elasticache.NewParameterGroup(ctx, "parameter-group",
		&elasticache.ParameterGroupArgs{
			Name:        pulumi.String(locals.ClusterIdentifier),
			Family:      pulumi.String(spec.ParameterGroupFamily),
			Description: pulumi.Sprintf("Custom parameter group for %s", locals.ClusterIdentifier),
			Parameters:  parameters,
			Tags:        pulumi.ToStringMap(locals.AwsTags),
		}, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "create parameter group")
	}

	ctx.Export(OpParameterGroupName, createdParameterGroup.Name)
	return createdParameterGroup, nil
}
