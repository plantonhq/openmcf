package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/vertex"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func deployedIndex(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpVertexAiDeployedIndex.Spec

	// The deployed index — the placement joining a Vector Search index to
	// an index endpoint, with its own serving compute. Deploying is a
	// long-running operation (tens of minutes). Only the replica bounds
	// inside the sizing arm update in place (the provider PATCHes them
	// via mutateDeployedIndex); every other change undeploys and
	// redeploys.
	//
	// No aiplatform API enablement here, deliberately: a deployment
	// cannot exist without its index endpoint, and creating the endpoint
	// (its own kind) already enabled the API. This resource also carries
	// no project field — the project rides inside the index_endpoint
	// resource path.
	args := &vertex.AiIndexEndpointDeployedIndexArgs{
		DeployedIndexId: pulumi.String(spec.DeployedIndexId),
		Index:           pulumi.String(spec.Index.GetValue()),
		IndexEndpoint:   pulumi.String(spec.IndexEndpoint.GetValue()),

		// The provider resolves the regional Vertex AI API host
		// (https://{region}-aiplatform.googleapis.com) from `region`;
		// without it the deploy fails unless the provider config happens
		// to carry one. Must match the endpoint's own region —
		// deployments cannot cross regions.
		Region: pulumi.StringPtr(spec.Location),

		// PARITY: the bridged provider carries a client-side deletion_policy
		// flag the released 6.x Terraform line does not have. Pinned to
		// DELETE so destroy really undeploys the index on both engines.
		DeletionPolicy: pulumi.StringPtr("DELETE"),
	}

	// Unusually for a display name, the API treats it as immutable on a
	// deployed index.
	if spec.DisplayName != "" {
		args.DisplayName = pulumi.StringPtr(spec.DisplayName)
	}

	// Vertex-managed serving compute. 0 means "not set": GCP then applies
	// its defaults (min 2, max = min).
	if spec.AutomaticResources != nil {
		automaticArgs := &vertex.AiIndexEndpointDeployedIndexAutomaticResourcesArgs{}
		if spec.AutomaticResources.GetMinReplicaCount() > 0 {
			automaticArgs.MinReplicaCount = pulumi.IntPtr(int(spec.AutomaticResources.GetMinReplicaCount()))
		}
		if spec.AutomaticResources.MaxReplicaCount > 0 {
			automaticArgs.MaxReplicaCount = pulumi.IntPtr(int(spec.AutomaticResources.MaxReplicaCount))
		}
		args.AutomaticResources = automaticArgs
	}

	// Explicitly pinned serving compute. machine_spec is a required block
	// even when machine_type is left to the API's default;
	// min_replica_count is required by the API (>= 1).
	if spec.DedicatedResources != nil {
		machineSpecArgs := &vertex.AiIndexEndpointDeployedIndexDedicatedResourcesMachineSpecArgs{}
		if spec.DedicatedResources.MachineType != "" {
			machineSpecArgs.MachineType = pulumi.StringPtr(spec.DedicatedResources.MachineType)
		}
		dedicatedArgs := &vertex.AiIndexEndpointDeployedIndexDedicatedResourcesArgs{
			MachineSpec:     machineSpecArgs,
			MinReplicaCount: pulumi.Int(int(spec.DedicatedResources.MinReplicaCount)),
		}
		if spec.DedicatedResources.MaxReplicaCount > 0 {
			dedicatedArgs.MaxReplicaCount = pulumi.IntPtr(int(spec.DedicatedResources.MaxReplicaCount))
		}
		args.DedicatedResources = dedicatedArgs
	}

	// IP-space partitioning: empty lets GCP default ("default"). The API
	// HOLDS the group↔ranges pairing — a non-default group, once used
	// with a set of reserved ranges, can only ever be used with exactly
	// that set (taught in the spec).
	if spec.GetDeploymentGroup() != "" {
		args.DeploymentGroup = pulumi.StringPtr(spec.GetDeploymentGroup())
	}

	if spec.EnableAccessLogging {
		args.EnableAccessLogging = pulumi.BoolPtr(true)
	}

	// Names of reserved VPC_PEERING address ranges under the endpoint's
	// peered network; only meaningful on a peered endpoint.
	if len(spec.ReservedIpRanges) > 0 {
		reservedRanges := make([]string, 0, len(spec.ReservedIpRanges))
		for _, rangeRef := range spec.ReservedIpRanges {
			reservedRanges = append(reservedRanges, rangeRef.GetValue())
		}
		args.ReservedIpRanges = pulumi.ToStringArray(reservedRanges)
	}

	// JWT auth on the private query endpoint. The provider nests the
	// API's single-child deployedIndexAuthConfig.authProvider wrapper;
	// the spec flattens it to one honest auth_config block.
	if spec.AuthConfig != nil {
		authProviderArgs := &vertex.AiIndexEndpointDeployedIndexDeployedIndexAuthConfigAuthProviderArgs{}
		if len(spec.AuthConfig.AllowedIssuers) > 0 {
			allowedIssuers := make([]string, 0, len(spec.AuthConfig.AllowedIssuers))
			for _, issuerRef := range spec.AuthConfig.AllowedIssuers {
				allowedIssuers = append(allowedIssuers, issuerRef.GetValue())
			}
			authProviderArgs.AllowedIssuers = pulumi.ToStringArray(allowedIssuers)
		}
		if len(spec.AuthConfig.Audiences) > 0 {
			authProviderArgs.Audiences = pulumi.ToStringArray(spec.AuthConfig.Audiences)
		}
		args.DeployedIndexAuthConfig = &vertex.AiIndexEndpointDeployedIndexDeployedIndexAuthConfigArgs{
			AuthProvider: authProviderArgs,
		}
	}

	createdDeployedIndex, err := vertex.NewAiIndexEndpointDeployedIndex(ctx, "vertex-ai-deployed-index", args,
		pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to create vertex ai deployed index")
	}

	ctx.Export(OpName, createdDeployedIndex.Name)
	ctx.Export(OpDeployedIndexId, createdDeployedIndex.DeployedIndexId)
	ctx.Export(OpIndexEndpoint, createdDeployedIndex.IndexEndpoint)
	ctx.Export(OpCreateTime, createdDeployedIndex.CreateTime)
	ctx.Export(OpIndexSyncTime, createdDeployedIndex.IndexSyncTime)

	// Private-endpoint addresses exist only on peered/PSC endpoints;
	// export empty strings otherwise so the output shape is stable —
	// identical to the Terraform module's try() fallbacks.
	ctx.Export(OpMatchGrpcAddress, createdDeployedIndex.PrivateEndpoints.ApplyT(
		func(privateEndpoints []vertex.AiIndexEndpointDeployedIndexPrivateEndpoint) string {
			if len(privateEndpoints) == 0 || privateEndpoints[0].MatchGrpcAddress == nil {
				return ""
			}
			return *privateEndpoints[0].MatchGrpcAddress
		}).(pulumi.StringOutput))
	ctx.Export(OpServiceAttachment, createdDeployedIndex.PrivateEndpoints.ApplyT(
		func(privateEndpoints []vertex.AiIndexEndpointDeployedIndexPrivateEndpoint) string {
			if len(privateEndpoints) == 0 || privateEndpoints[0].ServiceAttachment == nil {
				return ""
			}
			return *privateEndpoints[0].ServiceAttachment
		}).(pulumi.StringOutput))

	return nil
}
