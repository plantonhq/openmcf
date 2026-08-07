# GcpVertexAiDeployedIndex: Design Notes

## Service Overview

A **DeployedIndex** is the join object of Vertex AI Vector Search: it places an index (the data) onto an index endpoint (the serving surface) and gives that placement its own serving compute. The API object literally is `DeployedIndex` — an entry in the endpoint's `deployedIndexes[]` — which is why the kind keeps the provider-authentic name instead of inventing one.

### What a Deployed Index Does

- Makes an index queryable through an endpoint (nothing serves until this exists)
- Owns the serving compute: automatic (Vertex-managed) or dedicated (pinned machine type)
- Pins deployments to predictable IP spaces on peered networks (groups + reserved ranges)
- Optionally requires JWT auth on the private query endpoint

### What a Deployed Index Does NOT Do

- It does not store vectors (that's `GcpVertexAiIndex`)
- It does not own connectivity (that's `GcpVertexAiIndexEndpoint`)
- It is not a first-class GCP resource with its own project/labels — it lives inside the endpoint

## Deployment Landscape

### Terraform: `google_vertex_ai_index_endpoint_deployed_index`

```hcl
resource "google_vertex_ai_index_endpoint_deployed_index" "deployment" {
  deployed_index_id = "products_v1"
  index             = google_vertex_ai_index.index.id
  index_endpoint    = google_vertex_ai_index_endpoint.endpoint.id
  region            = "us-central1"
}
```

Key characteristics at the pinned released line:

- The resource is a projection of the endpoint's `deployIndex`/`undeployIndex`/`mutateDeployedIndex` RPCs; its Terraform ID is `{index_endpoint}/deployedIndex/{deployed_index_id}`
- EVERYTHING is ForceNew except the replica bounds inside the sizing arms (PATCHed via `mutateDeployedIndex`) — including, unusually, `display_name`
- No `labels`, no `project` — the resource class simply doesn't carry them
- `region` resolves the regional API host (`https://{region}-aiplatform.googleapis.com`); without it the deploy cannot be addressed
- Timeouts: create/update 45 minutes, delete 20 — deploys are genuinely slow

### Pulumi: `vertex.AiIndexEndpointDeployedIndex`

Same schema through the bridge, plus a client-side `deletion_policy` (pinned to `DELETE` for parity).

## Feature Coverage

| Feature | Coverage |
|---------|----------|
| Deployment handle (`deployed_index_id`) | Full — CEL regex mirrors the API contract |
| Index / endpoint references | Full — FK refs to the sibling kinds' composition keys |
| Automatic sizing (replica bounds) | Full |
| Dedicated sizing (machine type + bounds) | Full |
| Deployment groups | Full — pairing contract taught in the spec |
| Access logging | Full |
| Reserved IP ranges | Full — FK refs to GcpGlobalAddress names |
| JWT auth (issuers + audiences) | Full — issuers FK-ref GcpServiceAccount emails |

### Deliberate Exclusions (with reasons)

| Provider surface | Reason |
|------------------|--------|
| `deletion_policy` | 7.x-only client-side flag on the Terraform line; the bridged Pulumi provider has it and pins `DELETE` for identical destroy behavior (PARITY comment in the module). |
| `private_endpoints[].psc_automated_endpoints` output | Populated only by PSC automation configs, which the released 6.x endpoint resource cannot create — always empty from this catalog's IaC. |

## Design Decisions

1. **A first-class kind, not a field on the endpoint.** Deployments are many-per-endpoint, join two other kinds, and have their own lifecycle (deploy/undeploy without touching index or endpoint) — the textbook split-test pass. Embedding them in the endpoint would serialize every deployment change through one resource.
2. **`auth_config` flattens the provider's single-child wrapper.** The API nests `deployedIndexAuthConfig.authProvider.{allowed_issuers, audiences}`; the wrapper has exactly one child, so the spec exposes one honest `auth_config` block (the sibling endpoint's flattened `encryption_spec` precedent).
3. **`reserved_ip_ranges` FK-references `GcpGlobalAddress` by NAME.** The API wants compute address *names* under the peered VPC, and Vertex peering ranges are global INTERNAL `VPC_PEERING` addresses — the same composition `GcpServiceNetworkingConnection.reserved_peering_ranges` uses. Regional overrides remain possible via explicit `valueFrom.kind`.
4. **No fake labels.** The resource class has no labels in the API. Inventing platform attribution here (e.g. smuggling it into `display_name`) would corrupt a user-facing field; the honest answer is documented absence.
5. **`location` is required.** The regional API host cannot be derived reliably from provider config, and the sibling endpoint kind's live "Cannot determine region" failure class is the proof. One explicit field, documented as "must match the endpoint's region".
6. **Deployment-group pairing taught, not modeled.** The API's group↔ranges contract ("a non-default group, once used with ranges [a,b,c], can never be used with a different set") is a server-side constraint across resources; the spec documents it loudly rather than pretending a per-resource CEL rule could enforce it.

## Sizing Deep Dive

| | automatic_resources | dedicated_resources |
|---|---|---|
| Machine type | Vertex-managed | User-pinned (shard-size-compatible) |
| Min replicas | default 2 | required, >= 1 |
| Max replicas | default = min, cap 1000 | default = min, cap 1000 |
| SLA | none at min=1 | none at min=1 |
| In-place updates | replica bounds only | replica bounds only |

Machine-type compatibility by shard size (from the provider docs): SMALL shards allow `e2-standard-2` plus everything below; MEDIUM shards allow `e2-standard-16` plus LARGE's types; LARGE shards allow `e2-highmem-16` and `n2d-standard-32`.

## References

- [Deploy an Index](https://cloud.google.com/vertex-ai/docs/vector-search/deploy-index-public)
- [DeployedIndex API Reference](https://cloud.google.com/vertex-ai/docs/reference/rest/v1/projects.locations.indexEndpoints#DeployedIndex)
- [Terraform Resource](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/vertex_ai_index_endpoint_deployed_index)
- [Pulumi Resource](https://www.pulumi.com/registry/packages/gcp/api-docs/vertex/aiindexendpointdeployedindex/)
