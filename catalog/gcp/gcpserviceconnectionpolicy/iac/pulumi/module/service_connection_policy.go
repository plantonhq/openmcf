package module

import (
	"strconv"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/networkconnectivity"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// serviceConnectionPolicy provisions the per-network authorization that
// lets Google's service connectivity automation create PSC endpoints in
// the listed subnets on a producer's behalf. PSC-first managed services
// (Memorystore for Valkey, Redis Cluster) refuse to create instances on
// a network until a policy for their service class exists in the
// instance's region — this resource is that prerequisite.
//
// Cardinality is one policy per (network, service class, region) triple;
// GCP rejects a second. location, network, serviceClass, and the policy
// name are all immutable (ForceNew in the provider) — only the pscConfig
// contents, description, and labels update in place. Keep the policy
// alive as long as any instance depends on it: deleting it strands
// existing endpoints and blocks new ones.
func serviceConnectionPolicy(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpServiceConnectionPolicy.Spec

	// Enable the Network Connectivity API (the control plane that owns
	// service connection policies) and the Compute Engine API (the
	// network and subnets live in Compute, and the automation's
	// forwarding rules are Compute-side objects). disable_on_destroy
	// stays false: tearing down one policy must never disable the APIs
	// for everything else in the project. Honor the spec contract: an
	// empty project_id falls back to the provider's default project
	// (leaving Project unset lets the gcp provider resolve its own
	// project; an empty string would be sent verbatim and rejected).
	networkConnectivityApiArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("networkconnectivity.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
		DisableOnDestroy:         pulumi.BoolPtr(false),
	}
	computeApiArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("compute.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
		DisableOnDestroy:         pulumi.BoolPtr(false),
	}
	if spec.ProjectId.GetValue() != "" {
		networkConnectivityApiArgs.Project = pulumi.String(spec.ProjectId.GetValue())
		computeApiArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdNetworkConnectivityApi, err := projects.NewService(ctx,
		"scp-networkconnectivity.googleapis.com", networkConnectivityApiArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable networkconnectivity.googleapis.com api")
	}
	createdComputeApi, err := projects.NewService(ctx,
		"scp-compute.googleapis.com", computeApiArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable compute.googleapis.com api")
	}

	args := &networkconnectivity.ServiceConnectionPolicyArgs{
		Name:         pulumi.String(locals.PolicyName),
		Location:     pulumi.String(spec.Location),
		Network:      pulumi.String(toResourcePath(spec.Network.GetValue())),
		ServiceClass: pulumi.String(spec.ServiceClass),
		Labels:       pulumi.ToStringMap(locals.GcpLabels),
	}
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}

	if spec.PscConfig != nil {
		// Subnets may arrive as self-link URLs (the GcpSubnetwork
		// reference's canonical output) — normalized to relative
		// resource paths, identically to the Terraform module.
		subnetworks := make(pulumi.StringArray, 0, len(spec.PscConfig.Subnetworks))
		for _, s := range spec.PscConfig.Subnetworks {
			subnetworks = append(subnetworks, pulumi.String(toResourcePath(s.GetValue())))
		}

		pscConfig := &networkconnectivity.ServiceConnectionPolicyPscConfigArgs{
			Subnetworks: subnetworks,
		}
		// The API types the connection limit as a string-encoded
		// integer; 0 in the spec means "leave GCP's default".
		if spec.PscConfig.Limit > 0 {
			pscConfig.Limit = pulumi.StringPtr(strconv.Itoa(int(spec.PscConfig.Limit)))
		}
		if spec.PscConfig.ProducerInstanceLocation != "" {
			pscConfig.ProducerInstanceLocation = pulumi.StringPtr(spec.PscConfig.ProducerInstanceLocation)
		}
		if len(spec.PscConfig.AllowedGoogleProducersResourceHierarchyLevels) > 0 {
			levels := make(pulumi.StringArray, 0, len(spec.PscConfig.AllowedGoogleProducersResourceHierarchyLevels))
			for _, l := range spec.PscConfig.AllowedGoogleProducersResourceHierarchyLevels {
				levels = append(levels, pulumi.String(l))
			}
			pscConfig.AllowedGoogleProducersResourceHierarchyLevels = levels
		}
		args.PscConfig = pscConfig
	}
	// Empty defers to the provider default (DELETE).
	if spec.DeletionPolicy != "" {
		args.DeletionPolicy = pulumi.StringPtr(spec.DeletionPolicy)
	}

	createdPolicy, err := networkconnectivity.NewServiceConnectionPolicy(ctx, "policy", args,
		pulumi.Provider(gcpProvider),
		pulumi.DependsOn([]pulumi.Resource{createdNetworkConnectivityApi, createdComputeApi}))
	if err != nil {
		return errors.Wrap(err, "failed to create service connection policy")
	}

	ctx.Export(OpPolicyId, createdPolicy.ID())
	ctx.Export(OpName, pulumi.String(locals.PolicyName))
	ctx.Export(OpInfrastructure, createdPolicy.Infrastructure)
	ctx.Export(OpEtag, createdPolicy.Etag)

	return nil
}
