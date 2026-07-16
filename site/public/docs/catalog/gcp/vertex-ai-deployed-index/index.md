---
title: "Vertex AI Deployed Index"
description: "Vertex AI Deployed Index deployment documentation"
icon: "package"
order: 100
componentName: "gcpvertexaideployedindex"
---

# GCP Vertex AI Deployed Index

Deploys a Vertex AI Vector Search index onto an index endpoint — the placement that makes nearest-neighbor queries actually servable, with automatic or dedicated serving compute, deployment-group IP partitioning, reserved-range pinning, and optional JWT authentication. The index and the endpoint are separate resources (`GcpVertexAiIndex`, `GcpVertexAiIndexEndpoint`); this component joins them.

## What Gets Created

When you deploy a GcpVertexAiDeployedIndex resource, Planton provisions:

- **Deployed Index** — a `google_vertex_ai_index_endpoint_deployed_index` resource placing the referenced index onto the referenced endpoint with the chosen serving compute

Note: the GCP API gives this resource class no labels and no project field (it lives inside the endpoint resource), so platform label attribution is not possible here.

## Prerequisites

- **A `GcpVertexAiIndex`** in the same region (the `index` reference)
- **A `GcpVertexAiIndexEndpoint`** in the same region (the `indexEndpoint` reference)
- **Reserved `VPC_PEERING` global addresses** if pinning `reservedIpRanges` (peered endpoints only)
- **GCP credentials** configured via environment variables or Planton provider config

## Quick Start

Create a file `deployed-index.yaml`:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpVertexAiDeployedIndex
metadata:
  name: products-v1
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: ml-platform
    pulumi.planton.dev/stack.name: dev.GcpVertexAiDeployedIndex.products-v1
spec:
  location: us-central1
  deployedIndexId: products_v1
  index:
    valueFrom:
      kind: GcpVertexAiIndex
      name: product-embeddings
      fieldPath: status.outputs.index_id
  indexEndpoint:
    valueFrom:
      kind: GcpVertexAiIndexEndpoint
      name: vector-serving
      fieldPath: status.outputs.index_endpoint_id
```

Deploy:

```shell
planton apply -f deployed-index.yaml
```

This deploys the index with automatic serving compute at GCP's default bounds. Deployment is a long-running operation — expect tens of minutes.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `location` | `string` | Region of the index endpoint (the Vertex AI API host is regional; deployments cannot cross regions). Immutable. | Required, min length 1 |
| `deployedIndexId` | `string` | User-chosen deployment handle, unique within the project. Immutable. | Letter start; letters, numbers, underscores; up to 128 chars |
| `index` | `StringValueOrRef` | The index being deployed (full resource path). Can reference a GcpVertexAiIndex via `valueFrom`. Immutable. | Required |
| `indexEndpoint` | `StringValueOrRef` | The endpoint being deployed onto (full resource path). Can reference a GcpVertexAiIndexEndpoint via `valueFrom`. Immutable. | Required |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `displayName` | `string` | `""` | Display name — unusually, IMMUTABLE on this resource class (changing it redeploys). Up to 128 chars. |
| `automaticResources` | `object` | GCP default when both arms omitted | Vertex-managed serving compute. Mutually exclusive with `dedicatedResources`. |
| `automaticResources.minReplicaCount` | `int32` | `2` | Minimum replicas (no SLA at 1). Mutable in place. |
| `automaticResources.maxReplicaCount` | `int32` | = min | Maximum replicas, up to 1000. Mutable in place. |
| `dedicatedResources` | `object` | — | Pinned serving compute. Mutually exclusive with `automaticResources`. |
| `dedicatedResources.machineType` | `string` | API default | Machine type (e.g. `e2-standard-16`); must be compatible with the index's `shardSize`. Immutable. |
| `dedicatedResources.minReplicaCount` | `int32` | — | Required, >= 1. Mutable in place. |
| `dedicatedResources.maxReplicaCount` | `int32` | = min | Up to 1000. Mutable in place. |
| `deploymentGroup` | `string` | `default` | IP-space partitioning group (<= 64 chars, at most 5 besides `default`). The API permanently pairs a non-default group with the exact reserved-range set it first ships with. Immutable. |
| `enableAccessLogging` | `bool` | `false` | Send private-endpoint access logs to Cloud Logging. Immutable. |
| `reservedIpRanges` | `StringValueOrRef[]` | — | Names of reserved `VPC_PEERING` global addresses to deploy into (peered endpoints only). Can reference GcpGlobalAddress resources via `valueFrom`. Immutable. |
| `authConfig` | `object` | — | JWT auth on the private query endpoint. Immutable. |
| `authConfig.allowedIssuers` | `StringValueOrRef[]` | — | Service-account emails whose signed JWTs are accepted. Can reference GcpServiceAccount resources via `valueFrom`. |
| `authConfig.audiences` | `string[]` | — | Accepted JWT audiences. |

## Examples

### Dedicated Compute with Explicit Bounds

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpVertexAiDeployedIndex
metadata:
  name: prod-products
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: ml-platform
    pulumi.planton.dev/stack.name: prod.GcpVertexAiDeployedIndex.prod-products
spec:
  location: us-central1
  deployedIndexId: products_prod
  index:
    value: projects/my-project/locations/us-central1/indexes/1234567890
  indexEndpoint:
    value: projects/my-project/locations/us-central1/indexEndpoints/9876543210
  displayName: Products Production
  dedicatedResources:
    machineType: e2-standard-16
    minReplicaCount: 2
    maxReplicaCount: 10
```

### Peered Endpoint with Reserved Ranges and JWT Auth

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpVertexAiDeployedIndex
metadata:
  name: private-products
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: ml-platform
    pulumi.planton.dev/stack.name: prod.GcpVertexAiDeployedIndex.private-products
spec:
  location: us-central1
  deployedIndexId: products_private
  index:
    valueFrom:
      kind: GcpVertexAiIndex
      name: product-embeddings
      fieldPath: status.outputs.index_id
  indexEndpoint:
    valueFrom:
      kind: GcpVertexAiIndexEndpoint
      name: private-search
      fieldPath: status.outputs.index_endpoint_id
  deploymentGroup: prod
  enableAccessLogging: true
  reservedIpRanges:
    - valueFrom:
        kind: GcpGlobalAddress
        name: vertex-ai-range-a
        fieldPath: status.outputs.name
  authConfig:
    allowedIssuers:
      - valueFrom:
          kind: GcpServiceAccount
          name: query-client
          fieldPath: status.outputs.email
    audiences:
      - vector-search-clients
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `name` | `string` | Name of the DeployedIndex resource as the provider reports it |
| `deployed_index_id` | `string` | The user-chosen deployment handle |
| `create_time` | `string` | RFC3339 timestamp of when the deployment was created |
| `index_sync_time` | `string` | Timestamp up to which the deployment reflects the source index's updates — at least the index's `update_time` means in sync |
| `match_grpc_address` | `string` | Private gRPC address for match queries (peered endpoints only) |
| `service_attachment` | `string` | PSC service attachment consumers target with forwarding rules (PSC endpoints only) |
| `index_endpoint` | `string` | Full resource path of the endpoint this deployment lives on — query clients pair it with `deployed_index_id` |

## Related Components

- [GcpVertexAiIndex](/docs/catalog/gcp/vertex-ai-index) — the vector index being deployed
- [GcpVertexAiIndexEndpoint](/docs/catalog/gcp/vertex-ai-index-endpoint) — the serving surface being deployed onto
- [GcpGlobalAddress](/docs/catalog/gcp/global-address) — reserved `VPC_PEERING` ranges for `reservedIpRanges`
- [GcpServiceAccount](/docs/catalog/gcp/service-account) — JWT issuers for `authConfig.allowedIssuers`
