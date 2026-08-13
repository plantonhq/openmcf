package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/neptune"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// customEndpoints provisions the folded custom endpoints: an endpoint's
// identity IS this cluster (it fronts a subset of the cluster's own
// instances), so endpoints live here rather than as their own kind. Each
// entry is keyed by its name -- adding, renaming, or removing one endpoint
// never touches its siblings or the cluster. Member lists name
// spec.instances entries and are translated to the derived AWS instance
// identifiers the same way the instance resources build them.
func customEndpoints(ctx *pulumi.Context, locals *Locals, provider *aws.Provider,
	createdCluster *neptune.Cluster,
	createdInstancesByName map[string]*neptune.ClusterInstance) (pulumi.StringMap, error) {
	spec := locals.AwsNeptuneCluster.Spec

	endpointAddresses := pulumi.StringMap{}
	for _, endpointSpec := range spec.CustomEndpoints {
		args := &neptune.ClusterEndpointArgs{
			ClusterIdentifier:         createdCluster.ID(),
			ClusterEndpointIdentifier: pulumi.String(endpointSpec.Name),
			EndpointType:              pulumi.String(endpointSpec.EndpointType),
			Tags:                      pulumi.ToStringMap(locals.AwsTags),
		}

		// Pin to exactly these instances (by their spec.instances entry
		// names, mapped to the derived identifiers)...
		if len(endpointSpec.StaticMembers) > 0 {
			staticMembers := pulumi.StringArray{}
			for _, member := range endpointSpec.StaticMembers {
				staticMembers = append(staticMembers, createdInstancesByName[member].Identifier)
			}
			args.StaticMembers = staticMembers
		}

		// ...or front every instance of the endpoint's type EXCEPT these
		// (mutually exclusive with static_members -- CEL-enforced).
		if len(endpointSpec.ExcludedMembers) > 0 {
			excludedMembers := pulumi.StringArray{}
			for _, member := range endpointSpec.ExcludedMembers {
				excludedMembers = append(excludedMembers, createdInstancesByName[member].Identifier)
			}
			args.ExcludedMembers = excludedMembers
		}

		createdEndpoint, err := neptune.NewClusterEndpoint(ctx,
			fmt.Sprintf("custom-endpoint-%s", endpointSpec.Name), args, pulumi.Provider(provider))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create custom endpoint %s", endpointSpec.Name)
		}
		endpointAddresses[endpointSpec.Name] = createdEndpoint.Endpoint
	}

	return endpointAddresses, nil
}
