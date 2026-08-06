package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/docdb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// clusterParameterGroup manages a cluster parameter group only when the
// spec carries inline parameters -- a named parameter list is
// configuration owned by exactly this cluster, not a composable node.
// The family is derived from the pinned engine_version (locals.go); a
// version change that crosses families forces a new group, which is
// exactly AWS's own constraint (parameter families are
// engine-major-scoped).
func clusterParameterGroup(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*docdb.ClusterParameterGroup, error) {
	spec := locals.AwsDocumentDb.Spec
	if len(spec.Parameters) == 0 {
		return nil, nil
	}

	parameters := docdb.ClusterParameterGroupParameterArray{}
	for _, parameter := range spec.Parameters {
		parameterArgs := &docdb.ClusterParameterGroupParameterArgs{
			Name:  pulumi.String(parameter.Name),
			Value: pulumi.String(parameter.Value),
		}
		// "immediate" applies dynamic parameters now; static parameters
		// must be "pending-reboot". Empty defers to the provider default
		// (immediate).
		if parameter.ApplyMethod != "" {
			parameterArgs.ApplyMethod = pulumi.String(parameter.ApplyMethod)
		}
		parameters = append(parameters, parameterArgs)
	}

	createdParameterGroup, err := docdb.NewClusterParameterGroup(ctx, "cluster-parameter-group",
		&docdb.ClusterParameterGroupArgs{
			Name:       pulumi.String(locals.ClusterIdentifier),
			Family:     pulumi.String(engineFamily(spec.EngineVersion)),
			Parameters: parameters,
			Tags:       pulumi.ToStringMap(locals.AwsTags),
		}, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create cluster parameter group")
	}
	return createdParameterGroup, nil
}
