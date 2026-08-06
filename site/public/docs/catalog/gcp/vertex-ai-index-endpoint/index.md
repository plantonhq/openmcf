---
title: "Vertex AI Index Endpoint"
description: "Vertex AI Index Endpoint deployment documentation"
icon: "package"
order: 100
componentName: "gcpvertexaiindexendpoint"
---

# GCP Vertex AI Index Endpoint

Deploys a GCP Vertex AI Vector Search index endpoint — the serving surface that deployed indexes answer nearest-neighbor queries through, with public, VPC-peered, or Private Service Connect connectivity. This is a different GCP resource from the online-prediction Vertex AI endpoint (which serves models). Placing an index onto the endpoint is a separate `GcpVertexAiDeployedIndex` resource.

## What Gets Created

When you deploy a GcpVertexAiIndexEndpoint resource, Planton provisions:

- **Vector Search Index Endpoint** — a `google_vertex_ai_index_endpoint` resource in the specified region, labeled with your `labels` merged beneath the platform's attribution labels
- **API Enablement** — the Vertex AI API (`aiplatform.googleapis.com`) is enabled in the target project (never disabled on destroy)

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **A VPC network with Private Services Access** if using VPC-peered connectivity (`network` field)

## Quick Start

Create a file `index-endpoint.yaml`:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpVertexAiIndexEndpoint
metadata:
  name: vector-serving
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: ml-platform
    pulumi.planton.dev/stack.name: dev.GcpVertexAiIndexEndpoint.vector-serving
spec:
  location: us-central1
  displayName: Vector Serving
  publicEndpointEnabled: true
```

Deploy:

```shell
planton apply -f index-endpoint.yaml
```

This creates a public index endpoint in the provider's default project. Once a `GcpVertexAiDeployedIndex` places an index onto it, queries go to the `public_endpoint_domain_name` output.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `location` | `string` | Region for the endpoint (e.g., `us-central1`). Deployed indexes must live in the same region. Immutable. | Required, min length 1 |
| `displayName` | `string` | Human-readable name (the numeric resource ID is GCP-assigned). Mutable. | Required, 1-128 characters |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `projectId` | `StringValueOrRef` | provider default project | GCP project where the endpoint is created. Can reference a GcpProject resource via `valueFrom`. |
| `description` | `string` | `""` | Description of the endpoint. Mutable. |
| `publicEndpointEnabled` | `bool` | `false` | Serve queries over a public domain name. Mutually exclusive with `network` and `privateServiceConnectConfig`. Immutable. |
| `network` | `StringValueOrRef` | — | VPC network to peer into. The modules normalize a self-link to the API's relative form (`projects/{project}/global/networks/{name}`). Requires Private Services Access. Mutually exclusive with the other arms. Immutable. Can reference a GcpVpcNetwork resource via `valueFrom`. |
| `privateServiceConnectConfig` | `object` | — | Private Service Connect configuration. Mutually exclusive with the other arms. Immutable. |
| `privateServiceConnectConfig.enablePrivateServiceConnect` | `bool` | — | Must be `true` when the block is present (the API's enablement flag). |
| `privateServiceConnectConfig.projectAllowlist` | `string[]` | `[]` | Projects allowed to create forwarding rules targeting the endpoint's service attachment. |
| `labels` | `map(string)` | `{}` | User labels for cost attribution and ownership; merged beneath platform labels. Mutable. |

## Examples

### Public Endpoint

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpVertexAiIndexEndpoint
metadata:
  name: search-serving
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: ml-platform
    pulumi.planton.dev/stack.name: prod.GcpVertexAiIndexEndpoint.search-serving
spec:
  projectId:
    value: my-gcp-project
  location: us-central1
  displayName: Search Serving
  publicEndpointEnabled: true
```

### VPC-Peered Private Endpoint

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpVertexAiIndexEndpoint
metadata:
  name: private-search
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: ml-platform
    pulumi.planton.dev/stack.name: prod.GcpVertexAiIndexEndpoint.private-search
spec:
  projectId:
    value: my-gcp-project
  location: us-central1
  displayName: Private Search
  network:
    value: projects/123456789/global/networks/prod-vpc
```

### Private Service Connect Endpoint

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpVertexAiIndexEndpoint
metadata:
  name: psc-search
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: ml-platform
    pulumi.planton.dev/stack.name: prod.GcpVertexAiIndexEndpoint.psc-search
spec:
  projectId:
    value: my-gcp-project
  location: us-central1
  displayName: PSC Search
  privateServiceConnectConfig:
    enablePrivateServiceConnect: true
    projectAllowlist:
      - consumer-project-a
      - consumer-project-b
```

### Using Foreign Key References

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpVertexAiIndexEndpoint
metadata:
  name: composed-serving
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: ml-platform
    pulumi.planton.dev/stack.name: prod.GcpVertexAiIndexEndpoint.composed-serving
spec:
  projectId:
    valueFrom:
      kind: GcpProject
      name: ml-project
      fieldPath: status.outputs.project_id
  location: us-central1
  displayName: Composed Vector Serving
  network:
    valueFrom:
      kind: GcpVpcNetwork
      name: ml-vpc
      fieldPath: status.outputs.network_self_link
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `index_endpoint_id` | `string` | Fully qualified endpoint resource path: `projects/{project}/locations/{location}/indexEndpoints/{id}` — the value a `GcpVertexAiDeployedIndex` passes as its `index_endpoint` reference |
| `index_endpoint_name` | `string` | The GCP-assigned numeric endpoint ID (the last path segment of `index_endpoint_id`) |
| `public_endpoint_domain_name` | `string` | Domain name for public querying. Populated only when `publicEndpointEnabled` is `true` |
| `create_time` | `string` | RFC3339 timestamp of when the endpoint was created |
| `update_time` | `string` | RFC3339 timestamp of when the endpoint was last updated |

## Related Components

- [GcpProject](/docs/catalog/gcp/project) — provides the GCP project for the endpoint
- [GcpVertexAiIndex](/docs/catalog/gcp/vertex-ai-index) — the vector index deployed onto this endpoint
- [GcpVertexAiDeployedIndex](/docs/catalog/gcp/vertex-ai-deployed-index) — places an index onto this endpoint for querying
- [GcpVpcNetwork](/docs/catalog/gcp/vpc) — provides the VPC network for peered private endpoints
