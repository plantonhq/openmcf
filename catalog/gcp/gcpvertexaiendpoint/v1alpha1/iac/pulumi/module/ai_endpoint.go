package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/vertex"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func aiEndpoint(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpVertexAiEndpoint.Spec

	// Enable the Vertex AI API — the control plane that owns endpoints.
	// DisableOnDestroy stays false: tearing down one endpoint must never
	// disable the API for everything else in the project (other endpoints
	// keep serving).
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
		"gcpvep-aiplatform.googleapis.com", aiplatformApiArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable aiplatform.googleapis.com api")
	}

	// The Vertex AI endpoint — the durable serving surface models deploy
	// onto. `Name` is ALWAYS sent: the API requires a numeric ID (no
	// leading zeros, max 10 digits) and never generates one, and the
	// engine's own auto-naming would produce a non-numeric name the API
	// rejects. locals.EndpointName is the spec's explicit value or the
	// identity-derived ID shared byte-for-byte with the Terraform module.
	args := &vertex.AiEndpointArgs{
		Name:        pulumi.StringPtr(locals.EndpointName),
		DisplayName: pulumi.String(spec.DisplayName),
		Location:    pulumi.String(spec.Location),

		// The provider resolves the Vertex AI API host from `region`
		// (https://{region}-aiplatform.googleapis.com), never from
		// `location` — without it, deploys fail with "Cannot determine
		// region" unless the provider config happens to carry one. The two
		// fields are the same axis for this regional API, so the module
		// pins region to location and the spec keeps a single honest field.
		Region: pulumi.StringPtr(spec.Location),
		Labels: pulumi.ToStringMap(locals.GcpLabels),

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

	// VPC-peered private networking (requires Private Services Access on
	// the network; mutually exclusive with PSC — enforced pre-deploy by
	// the spec's CEL rule).
	if spec.Network != nil && spec.Network.GetValue() != "" {
		args.Network = pulumi.StringPtr(spec.Network.GetValue())
	}

	// CMEK: the endpoint and all sub-resources encrypted under this key.
	if spec.KmsKeyName != nil && spec.KmsKeyName.GetValue() != "" {
		args.EncryptionSpec = &vertex.AiEndpointEncryptionSpecArgs{
			KmsKeyName: pulumi.String(spec.KmsKeyName.GetValue()),
		}
	}

	// Dedicated endpoint DNS: isolated hostname with better performance
	// and reliability than the shared regional DNS.
	if spec.DedicatedEndpointEnabled {
		args.DedicatedEndpointEnabled = pulumi.BoolPtr(true)
	}

	// Private Service Connect: the endpoint exposed via a service
	// attachment. The secure flag adds IAM authorization on top of
	// network reachability.
	if spec.PrivateServiceConnectConfig != nil {
		pscArgs := &vertex.AiEndpointPrivateServiceConnectConfigArgs{
			EnablePrivateServiceConnect: pulumi.Bool(true),
		}
		if spec.PrivateServiceConnectConfig.EnableSecurePrivateServiceConnect {
			pscArgs.EnableSecurePrivateServiceConnect = pulumi.BoolPtr(true)
		}
		if len(spec.PrivateServiceConnectConfig.ProjectAllowlist) > 0 {
			pscArgs.ProjectAllowlists = pulumi.ToStringArray(spec.PrivateServiceConnectConfig.ProjectAllowlist)
		}
		args.PrivateServiceConnectConfig = pscArgs
	}

	// Request/response logging: samples online predictions into BigQuery —
	// the raw material for drift monitoring and audit. sampling_rate 0
	// means "not set" (the API then applies its own default).
	if spec.RequestResponseLoggingConfig != nil {
		loggingArgs := &vertex.AiEndpointPredictRequestResponseLoggingConfigArgs{}
		if spec.RequestResponseLoggingConfig.Enabled {
			loggingArgs.Enabled = pulumi.BoolPtr(true)
		}
		if spec.RequestResponseLoggingConfig.SamplingRate != 0 {
			loggingArgs.SamplingRate = pulumi.Float64Ptr(spec.RequestResponseLoggingConfig.SamplingRate)
		}
		if spec.RequestResponseLoggingConfig.BigqueryDestinationUri != "" {
			loggingArgs.BigqueryDestination = &vertex.AiEndpointPredictRequestResponseLoggingConfigBigqueryDestinationArgs{
				OutputUri: pulumi.StringPtr(spec.RequestResponseLoggingConfig.BigqueryDestinationUri),
			}
		}
		args.PredictRequestResponseLoggingConfig = loggingArgs
	}

	createdEndpoint, err := vertex.NewAiEndpoint(ctx, "vertex-ai-endpoint", args,
		pulumi.Provider(gcpProvider),
		pulumi.DependsOn([]pulumi.Resource{createdAiplatformApi}))
	if err != nil {
		return errors.Wrap(err, "failed to create vertex ai endpoint")
	}

	ctx.Export(OpEndpointId, createdEndpoint.ID())
	ctx.Export(OpDisplayName, createdEndpoint.DisplayName)
	ctx.Export(OpDedicatedEndpointDns, createdEndpoint.DedicatedEndpointDns)
	ctx.Export(OpCreateTime, createdEndpoint.CreateTime)
	ctx.Export(OpEndpointName, createdEndpoint.Name)

	return nil
}
