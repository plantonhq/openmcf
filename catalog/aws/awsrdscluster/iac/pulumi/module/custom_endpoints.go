package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/rds"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// customEndpoints provisions the spec's custom cluster endpoints: one
// resource per entry, keyed by the entry's name so endpoints come and
// go without touching the cluster or its instances. Members reference
// the folded instances BY THEIR SPEC NAMES -- resolving them through
// the created instance resources maps each name to its full AWS
// identifier and makes a typo'd member fail the preview instead of
// silently fronting nothing.
func customEndpoints(ctx *pulumi.Context, locals *Locals, provider *aws.Provider,
	createdCluster *rds.Cluster, createdInstances []*rds.ClusterInstance) ([]*rds.ClusterEndpoint, error) {
	spec := locals.AwsRdsCluster.Spec

	instanceByName := make(map[string]*rds.ClusterInstance, len(createdInstances))
	for i, instance := range spec.Instances {
		instanceByName[instance.Name] = createdInstances[i]
	}

	resolveMembers := func(names []string) (pulumi.StringArray, error) {
		members := make(pulumi.StringArray, 0, len(names))
		for _, name := range names {
			createdInstance, ok := instanceByName[name]
			if !ok {
				return nil, errors.Errorf("custom endpoint member %q names no spec.instances entry", name)
			}
			members = append(members, createdInstance.Identifier)
		}
		return members, nil
	}

	createdEndpoints := make([]*rds.ClusterEndpoint, 0, len(spec.CustomEndpoints))
	for _, endpoint := range spec.CustomEndpoints {
		args := &rds.ClusterEndpointArgs{
			ClusterIdentifier:         createdCluster.ClusterIdentifier,
			ClusterEndpointIdentifier: pulumi.Sprintf("%s-%s", locals.ClusterIdentifier, endpoint.Name),
			CustomEndpointType:        pulumi.String(endpoint.Type),
			Tags:                      pulumi.ToStringMap(locals.AwsTags),
		}

		if len(endpoint.StaticMembers) > 0 {
			members, err := resolveMembers(endpoint.StaticMembers)
			if err != nil {
				return nil, err
			}
			args.StaticMembers = members
		}
		if len(endpoint.ExcludedMembers) > 0 {
			members, err := resolveMembers(endpoint.ExcludedMembers)
			if err != nil {
				return nil, err
			}
			args.ExcludedMembers = members
		}

		createdEndpoint, err := rds.NewClusterEndpoint(ctx,
			fmt.Sprintf("custom-endpoint-%s", endpoint.Name), args, pulumi.Provider(provider))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create custom endpoint %s", endpoint.Name)
		}
		createdEndpoints = append(createdEndpoints, createdEndpoint)
	}

	return createdEndpoints, nil
}
