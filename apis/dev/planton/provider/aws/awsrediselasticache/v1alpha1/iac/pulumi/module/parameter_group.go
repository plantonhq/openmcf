package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/elasticache"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// parameterGroup manages a custom ElastiCache parameter group only when
// the spec carries inline parameters -- a named parameter list is
// configuration owned by exactly this cluster, not a composable node.
// Bringing an existing group name and inline parameters are mutually
// exclusive (CEL-enforced).
func parameterGroup(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*elasticache.ParameterGroup, error) {
	spec := locals.AwsRedisElasticache.Spec
	if len(spec.Parameters) == 0 || spec.ParameterGroupName != "" {
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
			Name:        pulumi.String(locals.ReplicationGroupId),
			Family:      pulumi.String(spec.ParameterGroupFamily),
			Description: pulumi.Sprintf("Custom parameter group for %s", locals.ReplicationGroupId),
			Parameters:  parameters,
			Tags:        pulumi.ToStringMap(locals.AwsTags),
		}, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create ElastiCache parameter group")
	}
	return createdParameterGroup, nil
}
