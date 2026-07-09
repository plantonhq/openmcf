package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/vertex"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func indexEndpoint(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpVertexAiIndexEndpoint.Spec

	// Enable the Vertex AI API — the control plane that owns index
	// endpoints. DisableOnDestroy stays false: tearing down one endpoint
	// must never disable the API for everything else in the project
	// (other Vertex resources keep working).
	aiplatformApiArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("aiplatform.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
		DisableOnDestroy:         pulumi.BoolPtr(false),
	}
	// Honor the spec contract: an empty project_id falls back to the
	// provider's default project.
	if spec.ProjectId.GetValue() != "" {
		aiplatformApiArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdAiplatformApi, err := projects.NewService(ctx,
		"gcpvaiep-aiplatform.googleapis.com", aiplatformApiArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable aiplatform.googleapis.com api")
	}

	// The Vector Search index endpoint — the serving surface deployed
	// indexes answer queries through. This is a DIFFERENT GCP resource
	// from the online-prediction vertex.AiEndpoint (which serves models).
	// GCP assigns the numeric resource ID; display_name is the human
	// handle. Every connectivity choice (public / peered network / PSC)
	// is immutable (ForceNew); display_name, description, and labels
	// PATCH in place.
	args := &vertex.AiIndexEndpointArgs{
		DisplayName: pulumi.String(spec.DisplayName),
		Region:      pulumi.StringPtr(spec.Location),
		Labels:      pulumi.ToStringMap(locals.GcpLabels),

		// PARITY: the bridged provider carries a client-side deletion_policy
		// flag the released 6.x Terraform line does not have. Pinned to
		// DELETE so destroy really deletes the endpoint on both engines.
		DeletionPolicy: pulumi.StringPtr("DELETE"),
	}
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.StringPtr(spec.ProjectId.GetValue())
	}
	if spec.Description != "" {
		args.Description = pulumi.StringPtr(spec.Description)
	}

	// Public querying arm: deployed indexes become reachable through
	// public_endpoint_domain_name.
	if spec.PublicEndpointEnabled {
		args.PublicEndpointEnabled = pulumi.BoolPtr(true)
	}

	// VPC-peered private querying (requires Private Services Access on
	// the network; mutually exclusive with the other arms — enforced
	// pre-deploy by the spec's CEL rules). locals.Network carries the
	// API's relative form regardless of whether the spec supplied a
	// self-link.
	if locals.Network != "" {
		args.Network = pulumi.StringPtr(locals.Network)
	}

	// Private Service Connect: consumers reach deployed indexes through a
	// service attachment (surfaced on the GcpVertexAiDeployedIndex
	// outputs once an index is deployed).
	if spec.PrivateServiceConnectConfig != nil {
		pscArgs := &vertex.AiIndexEndpointPrivateServiceConnectConfigArgs{
			EnablePrivateServiceConnect: pulumi.Bool(true),
		}
		if len(spec.PrivateServiceConnectConfig.ProjectAllowlist) > 0 {
			pscArgs.ProjectAllowlists = pulumi.ToStringArray(spec.PrivateServiceConnectConfig.ProjectAllowlist)
		}
		args.PrivateServiceConnectConfig = pscArgs
	}

	createdIndexEndpoint, err := vertex.NewAiIndexEndpoint(ctx, "vertex-ai-index-endpoint", args,
		pulumi.Provider(gcpProvider),
		pulumi.DependsOn([]pulumi.Resource{createdAiplatformApi}))
	if err != nil {
		return errors.Wrap(err, "failed to create vertex ai index endpoint")
	}

	ctx.Export(OpIndexEndpointId, createdIndexEndpoint.ID())
	ctx.Export(OpIndexEndpointName, createdIndexEndpoint.Name)
	ctx.Export(OpPublicEndpointDomainName, createdIndexEndpoint.PublicEndpointDomainName)
	ctx.Export(OpCreateTime, createdIndexEndpoint.CreateTime)
	ctx.Export(OpUpdateTime, createdIndexEndpoint.UpdateTime)

	return nil
}
