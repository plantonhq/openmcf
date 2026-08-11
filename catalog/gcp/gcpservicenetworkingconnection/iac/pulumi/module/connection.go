package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/servicenetworking"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// serviceNetworkingConnection provisions the private services access
// connection: a VPC peering between this network and the service producer's
// network, carved from the reserved VPC_PEERING address ranges. This single
// resource is what turns "Cloud SQL with private IP" from an error into a
// working deployment — producers allocate service subnets out of the
// reserved ranges and route them over the peering.
//
// Cardinality is one connection per (network, service) pair — GCP rejects a
// second. Capacity grows by appending range names to reservedPeeringRanges
// on THIS resource (an in-place update that never disturbs subnets the
// producer already provisioned), never by adding another connection.
//
// network and service are immutable (ForceNew in the provider): changing
// either destroys and recreates the connection, severing private
// connectivity for every producer resource on it. Teardown ordering: GCP
// refuses to delete the connection while the producer still holds subnets —
// destroy the private-IP service instances (Cloud SQL, AlloyDB,
// Memorystore, ...) before this resource.
func serviceNetworkingConnection(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpServiceNetworkingConnection.Spec

	// Enable the Service Networking API — the producer-side control plane
	// this connection talks to — and the Compute Engine API (the network and
	// the reserved ranges live in Compute). disable_on_destroy stays false
	// (the provider default): tearing down one connection must never disable
	// the APIs for everything else in the project. The project only scopes
	// this enablement — the connection itself is addressed by the network.
	// Honor the spec contract: an empty project_id falls back to the
	// provider's default project (leaving Project unset lets the gcp
	// provider resolve its own project; an empty string would be sent
	// verbatim and rejected).
	serviceNetworkingApiArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("servicenetworking.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
		DisableOnDestroy:         pulumi.BoolPtr(false),
	}
	computeApiArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("compute.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
		DisableOnDestroy:         pulumi.BoolPtr(false),
	}
	if spec.ProjectId.GetValue() != "" {
		serviceNetworkingApiArgs.Project = pulumi.String(spec.ProjectId.GetValue())
		computeApiArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdServiceNetworkingApi, err := projects.NewService(ctx,
		"snc-servicenetworking.googleapis.com", serviceNetworkingApiArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable servicenetworking.googleapis.com api")
	}
	createdComputeApi, err := projects.NewService(ctx,
		"snc-compute.googleapis.com", computeApiArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable compute.googleapis.com api")
	}

	// The reserved ranges arrive as GcpGlobalAddress NAMES (the resolver
	// flattens references before the module runs) — the API's own addressing
	// for peering ranges, never self-links or CIDRs.
	reservedRanges := make(pulumi.StringArray, 0, len(spec.ReservedPeeringRanges))
	for _, r := range spec.ReservedPeeringRanges {
		reservedRanges = append(reservedRanges, pulumi.String(r.GetValue()))
	}

	args := &servicenetworking.ConnectionArgs{
		Network:               pulumi.String(spec.Network.GetValue()),
		Service:               pulumi.String(locals.Service),
		ReservedPeeringRanges: reservedRanges,
	}

	// Adopts a pre-existing connection for the same pair instead of failing
	// with "Cannot modify allocated ranges" (see the spec comment).
	if spec.UpdateOnCreationFail {
		args.UpdateOnCreationFail = pulumi.BoolPtr(true)
	}
	// Empty defers to the provider default (DELETE).
	if spec.DeletionPolicy != "" {
		args.DeletionPolicy = pulumi.StringPtr(spec.DeletionPolicy)
	}

	createdConnection, err := servicenetworking.NewConnection(ctx, "connection", args,
		pulumi.Provider(gcpProvider),
		pulumi.DependsOn([]pulumi.Resource{createdServiceNetworkingApi, createdComputeApi}))
	if err != nil {
		return errors.Wrap(err, "failed to create service networking connection")
	}

	ctx.Export(OpPeering, createdConnection.Peering)
	ctx.Export(OpNetwork, createdConnection.Network)

	return nil
}
