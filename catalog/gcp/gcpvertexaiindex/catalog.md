# GCP Vertex AI Index

Deploys a Vertex AI Vector Search index: the data structure that holds embedding vectors and answers nearest-neighbor queries. The index is the first resource of the vector-search trio — a GcpVertexAiIndexEndpoint provides the serving surface and a GcpVertexAiDeployedIndex places the index onto it, after which queries can be served. The component integrates with Planton's Provider Connections for GCP credential management and supports ValueFromRef wiring to GCP projects.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Vertex AI Index** -- a regional Vector Search index with the configured embedding geometry (dimensions, distance measure, normalization, shard size)
- **Nearest-Neighbor Algorithm** -- tree-AH approximate search (tuned or GCP-default) or brute-force exact search, chosen by which algorithm block the spec carries; omitting both applies GCP's default tree-AH tuning
- **Initial Data Load** -- when `contentsDeltaUri` points at a Cloud Storage directory, GCP builds the index from those files on the first deploy; otherwise an empty index is created (the norm for stream-update regimes)

## Before You Deploy

### Planton Setup

- **GCP Provider Connection** -- an active connection in the Connect module with credentials for the target GCP project. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or browser OAuth authentication modes.

### GCP Project

- **A GCP project** where the index will be created. Provide the project ID directly or reference a GcpProject Cloud Resource via ValueFromRef.
- **Vertex AI API** enabled in the target project.
- **Embedding model decided** -- the index's dimensions and distance measure must match the embedding model's output exactly, and both are immutable.
- **Data format** (batch loads) -- vector files in Cloud Storage laid out per GCP's Vector Search format, in the same region.

## Deploy

### Console

Open the deployment store, find **GCP Vertex AI Index**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Streaming Tree-AH** preset in the [Presets](#presets) tab for a production-shaped streaming index.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpVertexAiIndex
metadata:
  name: catalog-embeddings
  org: acme-corp
  env: prod
spec:
  projectId:
    value: "acme-prod-12345"
  location: us-central1
  displayName: product-catalog-embeddings
  indexUpdateMethod: STREAM_UPDATE
  config:
    dimensions: 768
    approximateNeighborsCount: 150
    distanceMeasureType: DOT_PRODUCT_DISTANCE
    featureNormType: UNIT_L2_NORM
    treeAhConfig: {}
```

```shell
planton apply -f vertex-index.yaml
```

A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying the full vector-search composition, deploy the index (and endpoint) before the deployment that joins them:

```yaml
# In a GcpVertexAiDeployedIndex spec:
index:
  valueFrom:
    kind: GcpVertexAiIndex
    name: catalog-embeddings
    fieldPath: status.outputs.index_id
```

The InfraPipeline resolves the dependency graph, provisions the index first, then the deployment that serves it.

## Key Configuration

These are the most important decisions when configuring an index. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Geometry** -- `config.dimensions` (required) is the embedding model's output size; `distanceMeasureType` must match how the model was trained (blank keeps GCP's dot-product default). The ENTIRE config block is immutable — changing any of it recreates the index.

**Algorithm** -- `treeAhConfig` (approximate, the production choice) and `bruteForceConfig` (exact, for small corpora and ground truth) are mutually exclusive; omitting both applies GCP's default tree-AH tuning. `approximateNeighborsCount` is required whenever tree-AH is configured.

**Update regime** -- `indexUpdateMethod` is immutable: `BATCH_UPDATE` (bulk rebuilds from Cloud Storage, GCP's default) or `STREAM_UPDATE` (near-real-time upserts at a higher cost per GB).

**Data source** -- `contentsDeltaUri` points at a Cloud Storage DIRECTORY (not file). It is write-only on read and travels in its own single-field update, so out-of-band data loads show as a one-field diff on the next plan — expected, not drift.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **GcpProject** (optional) | `projectId` | `status.outputs.project_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `index_id` | Fully qualified index path (`projects/{p}/locations/{l}/indexes/{id}`) | The exact value a GcpVertexAiDeployedIndex's `index` join consumes |
| `index_name` | The GCP-assigned numeric index ID | Display, logging |
| `metadata_schema_uri` | Cloud Storage URI of the index's metadata schema | Tooling that inspects index metadata |
| `create_time` | RFC3339 creation timestamp | Audit |
| `update_time` | RFC3339 last-update timestamp | Data-freshness checks |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Streaming tree-AH** -- A stream-update index with tuned tree-AH for production corpora that must reflect writes within seconds. Start from the **Streaming Tree-AH** preset.

**Batch from Cloud Storage** -- A batch-update index built from embedding files in a GCS directory, with explicit shard sizing for predictable serving machine types. Start from the **Batch From GCS** preset.

**Brute-force evaluation** -- An exact-search twin used to measure a production index's recall against ground truth. Start from the **Brute Force Eval** preset.

## Works With

- [**GCP Project**](/cloud-catalog/gcp-project) -- provides the GCP project where the index is created
- [**GCP Vertex AI Index Endpoint**](/cloud-catalog/gcp-vertex-ai-index-endpoint) -- the serving surface this index is deployed onto
- [**GCP Vertex AI Deployed Index**](/cloud-catalog/gcp-vertex-ai-deployed-index) -- joins this index to an endpoint via `index`; the final resource that makes queries servable
- [**GCP GCS Bucket**](/cloud-catalog/gcp-gcs-bucket) -- holds the vector data files a batch index loads from (compose the bucket name into `contentsDeltaUri`)
