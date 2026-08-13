package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/neptune"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// clusterParameterGroup manages a cluster parameter group only when the
// spec carries inline parameters -- a named parameter list is
// configuration owned by exactly this cluster, not a composable node.
// The family is derived from the pinned engine_version (locals.go); a
// version change that crosses families forces a new group, which is
// exactly AWS's own constraint (parameter families are
// engine-major-scoped).
func clusterParameterGroup(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*neptune.ClusterParameterGroup, error) {
	spec := locals.AwsNeptuneCluster.Spec
	if len(spec.Parameters) == 0 {
		return nil, nil
	}

	parameters := neptune.ClusterParameterGroupParameterArray{}
	for _, parameter := range spec.Parameters {
		parameterArgs := &neptune.ClusterParameterGroupParameterArgs{
			Name:  pulumi.String(parameter.Name),
			Value: pulumi.String(parameter.Value),
		}
		// "immediate" applies dynamic parameters now; static parameters
		// must be "pending-reboot". Empty defers to the provider default,
		// which is pending-reboot at the pinned provider -- set
		// "immediate" explicitly when a dynamic parameter should land
		// right away.
		if parameter.ApplyMethod != "" {
			parameterArgs.ApplyMethod = pulumi.String(parameter.ApplyMethod)
		}
		parameters = append(parameters, parameterArgs)
	}

	createdParameterGroup, err := neptune.NewClusterParameterGroup(ctx, "cluster-parameter-group",
		&neptune.ClusterParameterGroupArgs{
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

// instanceParameterGroup manages the instance-level twin of the cluster
// parameter group: it exists only when the spec carries inline
// instance_parameters, and every folded instance that does not bring its
// own group adopts it. Same family derivation, same ownership reasoning.
func instanceParameterGroup(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*neptune.ParameterGroup, error) {
	spec := locals.AwsNeptuneCluster.Spec
	if len(spec.InstanceParameters) == 0 {
		return nil, nil
	}

	parameters := neptune.ParameterGroupParameterArray{}
	for _, parameter := range spec.InstanceParameters {
		parameterArgs := &neptune.ParameterGroupParameterArgs{
			Name:  pulumi.String(parameter.Name),
			Value: pulumi.String(parameter.Value),
		}
		// Same apply_method semantics as the cluster group: empty defers
		// to the provider default (pending-reboot at the pinned provider).
		if parameter.ApplyMethod != "" {
			parameterArgs.ApplyMethod = pulumi.String(parameter.ApplyMethod)
		}
		parameters = append(parameters, parameterArgs)
	}

	createdParameterGroup, err := neptune.NewParameterGroup(ctx, "instance-parameter-group",
		&neptune.ParameterGroupArgs{
			Name:       pulumi.Sprintf("%s-instance", locals.ClusterIdentifier),
			Family:     pulumi.String(engineFamily(spec.EngineVersion)),
			Parameters: parameters,
			Tags:       pulumi.ToStringMap(locals.AwsTags),
		}, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create instance parameter group")
	}
	return createdParameterGroup, nil
}
