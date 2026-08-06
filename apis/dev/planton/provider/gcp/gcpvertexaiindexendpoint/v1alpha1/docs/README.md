# GcpVertexAiIndexEndpoint: Design Notes

## Service Overview

A Vector Search **index endpoint** is the serving half of Vertex AI Vector Search: the network surface that deployed indexes answer nearest-neighbor queries through. The catalog models the trio as three kinds — the index (data), the index endpoint (serving surface), and the deployed index (the placement joining them) — because each has its own lifecycle, its own provider resource, and is referenced independently.

### What an Index Endpoint Does

- Provides the query surface (public domain, peered-VPC address, or PSC service attachment)
- Hosts one or more deployed indexes, each with its own compute sizing
- Persists across index redeployments — clients keep one stable address

### What an Index Endpoint Does NOT Do

- It does not store vectors (that's `GcpVertexAiIndex`)
- It does not place indexes onto itself (that's `GcpVertexAiDeployedIndex`)
- It does not serve model predictions (that's the online-prediction `GcpVertexAiEndpoint` — a DIFFERENT GCP resource that happens to share a similar name)

## Deployment Landscape

### Terraform: `google_vertex_ai_index_endpoint`

```hcl
resource "google_vertex_ai_index_endpoint" "endpoint" {
  display_name            = "My Index Endpoint"
  region                  = "us-central1"
  public_endpoint_enabled = true
}
```

Key characteristics at the pinned released line:

- `network`, `public_endpoint_enabled`, `private_service_connect_config`, and `region` are ForceNew
- `display_name`, `description`, and `labels` PATCH in place
- The provider documents `network` XOR `private_service_connect_config`; the public arm is additionally exclusive per the API's contract — all three exclusions are CEL rules in the spec
- Timeouts: create/update/delete 20 minutes each — endpoints are fast

### Pulumi: `vertex.AiIndexEndpoint`

Same schema through the bridge. The bridged provider adds a client-side `deletion_policy` (pinned to `DELETE` for parity) plus `encryption_spec` and `psc_automation_configs` that the released 6.x Terraform line does not have (see exclusions).

## Feature Coverage

| Feature | Coverage |
|---------|----------|
| Display name + description + labels | Full |
| Region (location) | Full |
| Public querying arm | Full |
| VPC peering via `network` | Full — references GcpVpcNetwork, self-link normalized |
| Private Service Connect (allowlist) | Full |

### Deliberate Exclusions (with reasons)

| Provider surface | Reason |
|------------------|--------|
| `encryption_spec` (CMEK) | Not in the released 6.x `google` provider line for this resource (7.x-only surface); the catalog floats `~> 6.0`. Revisit at the catalog-wide major bump. |
| `deletion_policy` | 7.x-only client-side flag on the Terraform line; the bridged Pulumi provider has it and pins `DELETE` for identical destroy behavior (PARITY comment in the module). |
| `private_service_connect_config.psc_automation_configs` | The provider marks it online-prediction-only ("used by Online Inference(Prediction) only") — dead surface on an index endpoint. |
| `etag` output | Optimistic-concurrency plumbing, not a composition key. |

## Network Format Normalization

The Vertex AI API expects the RELATIVE network form `projects/{project}/global/networks/{name}` (with `{project}` preferably the project NUMBER) and rejects full compute self-link URLs — a rejection class that only surfaces at live create (the AlloyDB/Service Networking precedent). `GcpVpcNetwork` references resolve to the self-link (the kind's canonical output), so both modules strip the `https://www.googleapis.com/compute/v1/` prefix identically; values already in relative form pass through untouched.

## Design Decisions

1. **A separate kind from `GcpVertexAiEndpoint` (671).** They are different GCP resources with different schemas, different serving planes, and different provider resources. Sharing a kind would mean a fake union type; the docs distinguish them loudly instead.
2. **PSC keeps the API's explicit `enable_private_service_connect` flag.** The provider marks the field Required inside the block; modeling it (with a CEL rule pinning it `true`) keeps the spec bijective with the API instead of inventing block-presence semantics the API doesn't have.
3. **All three connectivity exclusions as CEL.** The API rejects invalid combinations only at create time; three small message-level rules move every combination failure pre-deploy with named messages.
4. **`update_time` exported alongside `create_time`** — the deployed index's `index_sync_time` is compared against the index's update time in operational tooling, so both timestamps ride on every kind of the trio.

## References

- [Vector Search Overview](https://cloud.google.com/vertex-ai/docs/vector-search/overview)
- [Deploy an Index](https://cloud.google.com/vertex-ai/docs/vector-search/deploy-index-public)
- [Private Services Access Setup](https://cloud.google.com/vertex-ai/docs/vector-search/setup/vpc)
- [Terraform Resource](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/vertex_ai_index_endpoint)
- [Pulumi Resource](https://www.pulumi.com/registry/packages/gcp/api-docs/vertex/aiindexendpoint/)
