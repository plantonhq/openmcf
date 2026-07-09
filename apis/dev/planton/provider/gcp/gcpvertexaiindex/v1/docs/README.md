# GcpVertexAiIndex: Design Notes

## Service Overview

Vertex AI Vector Search (formerly Matching Engine) is Google Cloud's managed vector database. The **index** is its storage half: a data structure holding embedding vectors, organized for nearest-neighbor retrieval. The serving half is separate — an index endpoint plus a deployment placing the index onto it — which is why the catalog models the trio as three kinds (`GcpVertexAiIndex`, `GcpVertexAiIndexEndpoint`, `GcpVertexAiDeployedIndex`).

### What an Index Does

- Stores embedding vectors and organizes them for fast similarity search
- Supports approximate search (tree-AH) that keeps query latency flat as the corpus grows
- Supports exact search (brute-force) for small corpora and recall evaluation
- Accepts bulk data loads from Cloud Storage (batch) or per-vector upserts (streaming)

### What an Index Does NOT Do

- It does not serve queries — that requires deployment onto an index endpoint
- It does not generate embeddings (that's a model, e.g. on a `GcpVertexAiEndpoint`)
- It does not manage the endpoint or the deployment (separate kinds)

## Deployment Landscape

### Terraform: `google_vertex_ai_index`

```hcl
resource "google_vertex_ai_index" "index" {
  display_name = "My Index"
  region       = "us-central1"

  metadata {
    contents_delta_uri = "gs://bucket/embeddings/"
    config {
      dimensions                  = 768
      approximate_neighbors_count = 150
      algorithm_config {
        tree_ah_config {}
      }
    }
  }
}
```

Key characteristics at the pinned released line:

- `metadata.config` and everything inside it are ForceNew; `region` and `index_update_method` are ForceNew
- `display_name`, `description`, and `labels` PATCH in place
- The provider's `tree_ah_config`/`brute_force_config` `ExactlyOneOf` lists are EMPTY — the constraint is not enforced client-side, so the spec's CEL rule is the only pre-deploy guard
- `contents_delta_uri` is write-only (the GET never returns it) and a change to it travels in its own single-field PATCH
- Timeouts: create/update/delete 180 minutes each — index builds are genuinely slow at scale

### Pulumi: `vertex.AiIndex`

Same schema through the bridge. The bridged provider adds a client-side `deletion_policy` (pinned to `DELETE` for parity) and an `encryption_spec` block that the released 6.x Terraform line does not have (see exclusions).

## Feature Coverage

| Feature | Coverage |
|---------|----------|
| Display name + description + labels | Full |
| Region (location) | Full |
| Update regime (BATCH_UPDATE / STREAM_UPDATE) | Full |
| Initial/delta data from Cloud Storage | Full |
| Complete-overwrite semantics | Full |
| Dimensions, neighbors count, shard size | Full |
| Distance measure + feature norm | Full |
| Tree-AH tuning (leaf size, search percent) | Full |
| Brute-force algorithm | Full |

### Deliberate Exclusions (with reasons)

| Provider surface | Reason |
|------------------|--------|
| `encryption_spec` (CMEK) | Not in the released 6.x `google` provider line for this resource (7.x-only surface); the catalog floats `~> 6.0`. Revisit at the catalog-wide major bump. |
| `deletion_policy` | 7.x-only client-side flag on the Terraform line; the bridged Pulumi provider has it and pins `DELETE` for identical destroy behavior (PARITY comment in the module). |
| `deployed_indexes` output | Volatile operational state: entries appear/disappear as deployments come and go, so exporting it would perma-diff. The deployment relationship is modeled first-class by `GcpVertexAiDeployedIndex`. |
| `index_stats` output | Volatile data-plane state (vector/shard counts change with every upsert); an infrastructure output would always be stale. |
| `etag` output | Optimistic-concurrency plumbing, not a composition key. |

## Design Decisions

1. **`config` is required.** The provider marks `metadata` optional, but the API rejects an index without it (the bridged SDK documents this explicitly). A required `config` with required `dimensions` turns a guaranteed live failure into a pre-deploy message.
2. **`metadata` flattened into the spec.** The provider's `metadata` block is an API packaging artifact (the Index.metadata blob). Nesting a `metadata` field inside `spec` alongside the KRM resource `metadata` would be a permanent source of confusion; `contents_delta_uri`, `is_complete_overwrite`, and `config` live at the spec top level instead. Both modules reassemble the provider block.
3. **Algorithm XOR enforced by CEL.** The provider's `ExactlyOneOf` is empty (enforces nothing) and the API errors only at create time. `at_most_one_algorithm` + `tree_ah_requires_neighbors_count` move both failures pre-deploy.
4. **`contents_delta_uri` is a plain string.** It is a `gs://` DIRECTORY URI; no GCS kind output matches that shape (same precedent as the sibling endpoint's `bigquery_destination_uri`). The format and both provider quirks are taught on the field.
5. **No fake attribution on immutable geometry.** Everything mutable (display_name, description, labels) PATCHes in place; everything else is documented immutable so a changed manifest replaces the index knowingly.

## Update-Regime Deep Dive

### BATCH_UPDATE (default)

Contents are (re)built from Cloud Storage files. A rebuild is a long-running operation — tens of minutes to hours. `is_complete_overwrite` chooses between replace-all and delta semantics when new files are supplied.

### STREAM_UPDATE

Vectors are upserted/removed via `indexes.upsertDatapoints` within seconds, and deployed replicas pick the changes up in near-real-time. Streaming indexes cost more per GB and periodically compact in the background. The regime is chosen at creation and cannot be changed.

## References

- [Vector Search Overview](https://cloud.google.com/vertex-ai/docs/vector-search/overview)
- [Input Data Format](https://cloud.google.com/vertex-ai/docs/vector-search/setup/format-structure)
- [Index Configuration](https://cloud.google.com/vertex-ai/docs/vector-search/configuring-indexes)
- [Terraform Resource](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/vertex_ai_index)
- [Pulumi Resource](https://www.pulumi.com/registry/packages/gcp/api-docs/vertex/aiindex/)
