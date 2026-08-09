package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/vertex"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func vectorIndex(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpVertexAiIndex.Spec

	// Enable the Vertex AI API — the control plane that owns indexes.
	// DisableOnDestroy stays false: tearing down one index must never
	// disable the API for everything else in the project (other Vertex
	// resources keep working).
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
		"gcpvaidx-aiplatform.googleapis.com", aiplatformApiArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable aiplatform.googleapis.com api")
	}

	// Vector-search geometry. config is required by the spec (the API
	// rejects an index without metadata), so this block always exists.
	configArgs := &vertex.AiIndexMetadataConfigArgs{
		Dimensions: pulumi.Int(int(spec.Config.Dimensions)),
	}

	// Required by the API when tree-AH is used (CEL-enforced pre-deploy);
	// 0 means "not set" for brute-force or GCP-default algorithm.
	if spec.Config.ApproximateNeighborsCount > 0 {
		configArgs.ApproximateNeighborsCount = pulumi.IntPtr(int(spec.Config.ApproximateNeighborsCount))
	}

	// Computed when omitted: GCP picks a shard size from the data.
	if spec.Config.ShardSize != "" {
		configArgs.ShardSize = pulumi.StringPtr(spec.Config.ShardSize)
	}
	if spec.Config.GetDistanceMeasureType() != "" {
		configArgs.DistanceMeasureType = pulumi.StringPtr(spec.Config.GetDistanceMeasureType())
	}
	if spec.Config.GetFeatureNormType() != "" {
		configArgs.FeatureNormType = pulumi.StringPtr(spec.Config.GetFeatureNormType())
	}

	// At most one algorithm arm exists (CEL-enforced). Omitting both lets
	// GCP default (tree-AH with default tuning).
	if spec.Config.TreeAhConfig != nil || spec.Config.BruteForceConfig != nil {
		algorithmArgs := &vertex.AiIndexMetadataConfigAlgorithmConfigArgs{}
		if spec.Config.TreeAhConfig != nil {
			treeAhArgs := &vertex.AiIndexMetadataConfigAlgorithmConfigTreeAhConfigArgs{}
			// 0 means "not set": the provider then applies GCP's defaults
			// (1000 embeddings per leaf, 10% searched).
			if spec.Config.TreeAhConfig.GetLeafNodeEmbeddingCount() > 0 {
				treeAhArgs.LeafNodeEmbeddingCount = pulumi.IntPtr(int(spec.Config.TreeAhConfig.GetLeafNodeEmbeddingCount()))
			}
			if spec.Config.TreeAhConfig.GetLeafNodesToSearchPercent() > 0 {
				treeAhArgs.LeafNodesToSearchPercent = pulumi.IntPtr(int(spec.Config.TreeAhConfig.GetLeafNodesToSearchPercent()))
			}
			algorithmArgs.TreeAhConfig = treeAhArgs
		}
		if spec.Config.BruteForceConfig != nil {
			algorithmArgs.BruteForceConfig = &vertex.AiIndexMetadataConfigAlgorithmConfigBruteForceConfigArgs{}
		}
		configArgs.AlgorithmConfig = algorithmArgs
	}

	// The provider models the API's Index.metadata blob as this nested
	// block: data location + geometry together. contents_delta_uri is
	// write-only upstream (GCP never reports it back) and a change to it
	// travels in its own single-field PATCH — both quirks are documented
	// in the spec so a drift-looking diff after an out-of-band data load
	// is expected, not alarming.
	metadataArgs := &vertex.AiIndexMetadataArgs{
		Config: configArgs,
	}
	if spec.ContentsDeltaUri != "" {
		metadataArgs.ContentsDeltaUri = pulumi.StringPtr(spec.ContentsDeltaUri)
		// Only meaningful alongside contents_delta_uri.
		metadataArgs.IsCompleteOverwrite = pulumi.BoolPtr(spec.IsCompleteOverwrite)
	}

	// The Vector Search index — the data structure holding embedding
	// vectors. GCP assigns the numeric resource ID; display_name is the
	// human handle. Region, index_update_method, and the whole config
	// geometry are immutable (ForceNew). Creation is a long-running
	// operation: minutes for an empty streaming index, up to hours for a
	// large batch build.
	args := &vertex.AiIndexArgs{
		DisplayName: pulumi.String(spec.DisplayName),
		Region:      pulumi.StringPtr(spec.Location),
		Labels:      pulumi.ToStringMap(locals.GcpLabels),
		Metadata:    metadataArgs,
	}
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.StringPtr(spec.ProjectId.GetValue())
	}
	if spec.Description != "" {
		args.Description = pulumi.StringPtr(spec.Description)
	}

	// Client-side destroy behavior (DELETE deletes the corpus; PREVENT
	// refuses; ABANDON drops from state but keeps the index standing).
	// Empty follows the provider default (DELETE) — mirrored
	// zero-vs-omit with Terraform.
	if spec.DeletionPolicy != "" {
		args.DeletionPolicy = pulumi.StringPtr(spec.DeletionPolicy)
	}

	// CMEK: the key must be in the index's region and the Vertex AI
	// service agent needs cryptoKeyEncrypterDecrypter on it. Omitted
	// means Google-managed encryption. Immutable.
	if spec.KmsKeyName.GetValue() != "" {
		args.EncryptionSpec = &vertex.AiIndexEncryptionSpecArgs{
			KmsKeyName: pulumi.String(spec.KmsKeyName.GetValue()),
		}
	}
	// Empty means "let GCP default" (BATCH_UPDATE). Sent only when set so
	// the preview diff stays honest about who chose the value.
	if spec.GetIndexUpdateMethod() != "" {
		args.IndexUpdateMethod = pulumi.StringPtr(spec.GetIndexUpdateMethod())
	}

	createdIndex, err := vertex.NewAiIndex(ctx, "vertex-ai-index", args,
		pulumi.Provider(gcpProvider),
		pulumi.DependsOn([]pulumi.Resource{createdAiplatformApi}))
	if err != nil {
		return errors.Wrap(err, "failed to create vertex ai index")
	}

	ctx.Export(OpIndexId, createdIndex.ID())
	ctx.Export(OpIndexName, createdIndex.Name)
	ctx.Export(OpMetadataSchemaUri, createdIndex.MetadataSchemaUri)
	ctx.Export(OpCreateTime, createdIndex.CreateTime)
	ctx.Export(OpUpdateTime, createdIndex.UpdateTime)

	return nil
}
