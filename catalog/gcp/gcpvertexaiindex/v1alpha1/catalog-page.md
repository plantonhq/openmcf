# GCP Vertex AI Index

Deploys a GCP Vertex AI Vector Search index — the managed data structure that holds embedding vectors and answers nearest-neighbor queries for semantic search, recommendations, and RAG retrieval. Configure the search geometry (dimensions, algorithm, distance measure, sharding), pick a batch or streaming update regime, and optionally point at initial vector data in Cloud Storage. Serving queries is a separate step: deploy the index onto a `GcpVertexAiIndexEndpoint` with a `GcpVertexAiDeployedIndex`.

## What Gets Created

When you deploy a GcpVertexAiIndex resource, Planton provisions:

- **Vector Search Index** — a `google_vertex_ai_index` resource in the specified region, labeled with your `labels` merged beneath the platform's attribution labels
- **API Enablement** — the Vertex AI API (`aiplatform.googleapis.com`) is enabled in the target project (never disabled on destroy)

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **A Cloud Storage directory of vector files** if seeding initial data via `contentsDeltaUri` (format: https://cloud.google.com/vertex-ai/docs/vector-search/setup/format-structure)

## Quick Start

Create a file `index.yaml`:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpVertexAiIndex
metadata:
  name: product-embeddings
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: ml-platform
    pulumi.planton.dev/stack.name: dev.GcpVertexAiIndex.product-embeddings
spec:
  location: us-central1
  displayName: Product Embeddings
  indexUpdateMethod: STREAM_UPDATE
  config:
    dimensions: 768
    approximateNeighborsCount: 150
    treeAhConfig: {}
```

Deploy:

```shell
planton apply -f index.yaml
```

This creates an empty streaming index in the provider's default project, ready for vector upserts via the Vertex AI API. Creation of an empty streaming index takes a few minutes; large batch builds can take hours.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `location` | `string` | Region for the index (e.g., `us-central1`). Deployed indexes must live in the same region. Immutable. | Required, min length 1 |
| `displayName` | `string` | Human-readable name for the index (the numeric resource ID is GCP-assigned). Mutable. | Required, 1-128 characters |
| `config` | `object` | Vector-search geometry. The whole block is immutable. | Required |
| `config.dimensions` | `int32` | The embedding model's output size (e.g., 768). | Required, >= 1 |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `projectId` | `StringValueOrRef` | provider default project | GCP project where the index is created. Can reference a GcpProject resource via `valueFrom`. |
| `description` | `string` | `""` | Description of the index. Mutable. |
| `indexUpdateMethod` | `string` | `BATCH_UPDATE` | `BATCH_UPDATE` (bulk rebuilds from Cloud Storage) or `STREAM_UPDATE` (near-real-time upserts via the API). Immutable. |
| `contentsDeltaUri` | `string` | — | Cloud Storage DIRECTORY holding initial/delta vector data (`gs://bucket/path/`). May be set later to load data. Changes travel in their own single-field update; GCP never reports the value back on read. |
| `isCompleteOverwrite` | `bool` | `false` | When an update carries `contentsDeltaUri`: `true` replaces the whole index contents, `false` applies the files as a delta. |
| `config.approximateNeighborsCount` | `int32` | — | Candidates fetched by approximate search before exact reordering. Required when `treeAhConfig` is set. |
| `config.shardSize` | `string` | GCP picks | `SHARD_SIZE_SMALL` (2 GB), `SHARD_SIZE_MEDIUM` (20 GB), or `SHARD_SIZE_LARGE` (50 GB). Determines which machine types can serve the index. Immutable. |
| `config.distanceMeasureType` | `string` | `DOT_PRODUCT_DISTANCE` | `SQUARED_L2_DISTANCE`, `L1_DISTANCE`, `COSINE_DISTANCE`, or `DOT_PRODUCT_DISTANCE`. Match the embedding model's training. |
| `config.featureNormType` | `string` | `NONE` | `UNIT_L2_NORM` or `NONE`. Unit-normalized vectors with dot-product distance rank like cosine similarity. |
| `config.treeAhConfig` | `object` | — | Tree-AH approximate search (the production choice for large corpora). Mutually exclusive with `bruteForceConfig`. |
| `config.treeAhConfig.leafNodeEmbeddingCount` | `int32` | `1000` | Embeddings per leaf node. |
| `config.treeAhConfig.leafNodesToSearchPercent` | `int32` | `10` | Percent of leaves any query may search, 1-100. Raises recall at the cost of latency. |
| `config.bruteForceConfig` | `object` | — | Exact (exhaustive) search: perfect recall, linear cost. For small corpora and ground-truth evaluation. Mutually exclusive with `treeAhConfig`. |
| `labels` | `map(string)` | `{}` | User labels for cost attribution and ownership; merged beneath platform labels. Mutable. |

## Examples

### Streaming Index for RAG Retrieval

Near-real-time upserts with cosine-equivalent ranking:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpVertexAiIndex
metadata:
  name: rag-chunks
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: ml-platform
    pulumi.planton.dev/stack.name: prod.GcpVertexAiIndex.rag-chunks
spec:
  projectId:
    value: my-gcp-project
  location: us-central1
  displayName: RAG Document Chunks
  indexUpdateMethod: STREAM_UPDATE
  config:
    dimensions: 1536
    approximateNeighborsCount: 150
    distanceMeasureType: DOT_PRODUCT_DISTANCE
    featureNormType: UNIT_L2_NORM
    shardSize: SHARD_SIZE_MEDIUM
    treeAhConfig:
      leafNodeEmbeddingCount: 1000
      leafNodesToSearchPercent: 10
```

### Batch Index Built from Cloud Storage

Bulk-built from a directory of embedding files:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpVertexAiIndex
metadata:
  name: catalog-embeddings
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: ml-platform
    pulumi.planton.dev/stack.name: prod.GcpVertexAiIndex.catalog-embeddings
spec:
  projectId:
    value: my-gcp-project
  location: us-central1
  displayName: Catalog Embeddings
  contentsDeltaUri: gs://ml-embeddings/catalog/
  config:
    dimensions: 768
    approximateNeighborsCount: 100
    treeAhConfig: {}
```

### Brute-Force Ground-Truth Index

Exact search for measuring an approximate index's recall:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpVertexAiIndex
metadata:
  name: eval-ground-truth
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: ml-platform
    pulumi.planton.dev/stack.name: dev.GcpVertexAiIndex.eval-ground-truth
spec:
  projectId:
    value: my-gcp-project
  location: us-central1
  displayName: Eval Ground Truth
  contentsDeltaUri: gs://ml-embeddings/eval-sample/
  config:
    dimensions: 768
    bruteForceConfig: {}
```

### Using Foreign Key References

Reference other Planton-managed resources for composable infrastructure:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpVertexAiIndex
metadata:
  name: composed-index
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: ml-platform
    pulumi.planton.dev/stack.name: prod.GcpVertexAiIndex.composed-index
spec:
  projectId:
    valueFrom:
      kind: GcpProject
      name: ml-project
      fieldPath: status.outputs.project_id
  location: us-central1
  displayName: Composed Vector Index
  indexUpdateMethod: STREAM_UPDATE
  config:
    dimensions: 768
    approximateNeighborsCount: 150
    treeAhConfig: {}
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `index_id` | `string` | Fully qualified index resource path: `projects/{project}/locations/{location}/indexes/{indexId}` — the value a `GcpVertexAiDeployedIndex` passes as its `index` reference |
| `index_name` | `string` | The GCP-assigned numeric index ID (the last path segment of `index_id`) |
| `metadata_schema_uri` | `string` | Cloud Storage URI of the YAML schema describing index-type-specific metadata |
| `create_time` | `string` | RFC3339 timestamp of when the index was created |
| `update_time` | `string` | RFC3339 timestamp of when the index was last updated |

## Related Components

- [GcpProject](/docs/catalog/gcp/gcpproject) — provides the GCP project for the index
- [GcpVertexAiIndexEndpoint](/docs/catalog/gcp/gcpvertexaiindexendpoint) — the serving surface the index deploys onto
- [GcpVertexAiDeployedIndex](/docs/catalog/gcp/gcpvertexaideployedindex) — places this index onto an index endpoint for querying
- [GcpGcsBucket](/docs/catalog/gcp/gcpgcsbucket) — holds the embedding files referenced by `contentsDeltaUri`
