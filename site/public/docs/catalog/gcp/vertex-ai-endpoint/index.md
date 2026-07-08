---
title: "Vertex AI Endpoint"
description: "Vertex AI Endpoint deployment documentation"
icon: "package"
order: 100
componentName: "gcpvertexaiendpoint"
---

# GCP Vertex AI Endpoint

Deploys a GCP Vertex AI Endpoint — a stable serving surface for machine learning models with configurable networking (public, VPC-peered, or Private Service Connect), optional CMEK encryption, and optional dedicated DNS for isolated prediction traffic. Model deployment to the endpoint is a separate operational step.

## What Gets Created

When you deploy a GcpVertexAiEndpoint resource, Planton provisions:

- **Vertex AI Endpoint** — a `google_vertex_ai_endpoint` resource in the specified region, labeled with your `labels` merged beneath the platform's attribution labels
- **API Enablement** — the Vertex AI API (`aiplatform.googleapis.com`) is enabled in the target project (never disabled on destroy)
- **Stable Numeric Identity** — when `endpointName` is omitted, both provisioning engines derive the same numeric endpoint ID from the resource identity, so the same manifest always produces the same endpoint reference

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **A VPC network with Private Services Access** if using VPC-peered networking (`network` field)
- **A Cloud KMS key** in the same region as the endpoint if using CMEK encryption (`kmsKeyName` field)
- **IAM permissions** — the Vertex AI service agent must have `roles/cloudkms.cryptoKeyEncrypterDecrypter` on the KMS key if CMEK is enabled

## Quick Start

Create a file `endpoint.yaml`:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpVertexAiEndpoint
metadata:
  name: my-endpoint
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.GcpVertexAiEndpoint.my-endpoint
spec:
  location: us-central1
  displayName: My ML Endpoint
```

Deploy:

```shell
planton apply -f endpoint.yaml
```

This creates a public Vertex AI Endpoint in the provider's default project, accessible via the shared regional DNS with Google-managed encryption.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `location` | `string` | Region for the endpoint (e.g., `us-central1`). Immutable after creation. | Required, min length 1 |
| `displayName` | `string` | Human-readable name for the endpoint. | Required, 1-128 characters |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `projectId` | `StringValueOrRef` | provider default project | GCP project where the endpoint is created. Can reference a GcpProject resource via `valueFrom`. |
| `description` | `string` | `""` | Description of the endpoint. |
| `labels` | `map(string)` | `{}` | User labels for cost attribution and ownership; merged beneath platform labels. Mutable. |
| `network` | `StringValueOrRef` | — | VPC network for private endpoints via VPC peering. Format: `projects/{project}/global/networks/{network}`. Mutually exclusive with `privateServiceConnectConfig`. Immutable. Can reference a GcpVpcNetwork resource via `valueFrom`. |
| `kmsKeyName` | `StringValueOrRef` | — | Cloud KMS key for CMEK encryption. Format: `projects/{p}/locations/{l}/keyRings/{r}/cryptoKeys/{k}`. Immutable. Can reference a GcpKmsKey resource via `valueFrom`. |
| `dedicatedEndpointEnabled` | `bool` | `false` | Enables a dedicated DNS name for better performance and traffic isolation. Mutually exclusive with `privateServiceConnectConfig`. |
| `privateServiceConnectConfig` | `object` | — | Private Service Connect configuration. Mutually exclusive with `network` and `dedicatedEndpointEnabled`. |
| `privateServiceConnectConfig.projectAllowlist` | `string[]` | `[]` | Projects allowed to create forwarding rules targeting this endpoint. |
| `privateServiceConnectConfig.enableSecurePrivateServiceConnect` | `bool` | `false` | Require IAM authorization on PSC connections (secure PSC). |
| `requestResponseLoggingConfig.enabled` | `bool` | `false` | Sample online predictions into BigQuery. |
| `requestResponseLoggingConfig.samplingRate` | `double` | GCP default | Fraction of requests to log, in `(0, 1]`. |
| `requestResponseLoggingConfig.bigqueryDestinationUri` | `string` | — | BigQuery destination: `bq://project`, `bq://project.dataset`, or `bq://project.dataset.table`. |
| `endpointName` | `string` | identity-derived | Numeric-only GCP resource identifier (1-10 digits, no leading zeros). When omitted, a stable ID is derived from the resource identity — identical on both engines. Immutable. |

## Examples

### Public Endpoint with Description

A public endpoint for development or non-sensitive workloads:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpVertexAiEndpoint
metadata:
  name: dev-recommendations
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: ml-platform
    pulumi.planton.dev/stack.name: dev.GcpVertexAiEndpoint.dev-recommendations
spec:
  projectId:
    value: my-gcp-project
  location: us-central1
  displayName: Dev Recommendation Engine
  description: Development endpoint for A/B testing recommendation models
```

### VPC-Peered Private Endpoint with CMEK

Production endpoint with network isolation and customer-managed encryption:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpVertexAiEndpoint
metadata:
  name: prod-scoring
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: ml-platform
    pulumi.planton.dev/stack.name: prod.GcpVertexAiEndpoint.prod-scoring
spec:
  projectId:
    value: my-gcp-project
  location: us-central1
  displayName: Production Scoring Endpoint
  network:
    value: projects/123456789/global/networks/prod-vpc
  kmsKeyName:
    value: projects/my-gcp-project/locations/us-central1/keyRings/ml-ring/cryptoKeys/endpoint-key
  dedicatedEndpointEnabled: true
```

### Private Service Connect Endpoint

Strongest network isolation using PSC with an explicit project allowlist:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpVertexAiEndpoint
metadata:
  name: psc-inference
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: ml-platform
    pulumi.planton.dev/stack.name: prod.GcpVertexAiEndpoint.psc-inference
spec:
  projectId:
    value: my-gcp-project
  location: us-central1
  displayName: PSC Inference Endpoint
  privateServiceConnectConfig:
    projectAllowlist:
      - consumer-project-a
      - consumer-project-b
  kmsKeyName:
    value: projects/my-gcp-project/locations/us-central1/keyRings/ml-ring/cryptoKeys/endpoint-key
```

### Using Foreign Key References

Reference other Planton-managed resources for composable infrastructure:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpVertexAiEndpoint
metadata:
  name: composed-endpoint
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: ml-platform
    pulumi.planton.dev/stack.name: prod.GcpVertexAiEndpoint.composed-endpoint
spec:
  projectId:
    valueFrom:
      kind: GcpProject
      name: ml-project
      fieldPath: status.outputs.project_id
  location: us-central1
  displayName: Composed ML Endpoint
  network:
    valueFrom:
      kind: GcpVpcNetwork
      name: ml-vpc
      fieldPath: status.outputs.network_self_link
  kmsKeyName:
    valueFrom:
      kind: GcpKmsKey
      name: ml-encryption-key
      fieldPath: status.outputs.key_id
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `endpoint_id` | `string` | Fully qualified endpoint resource path: `projects/{project}/locations/{location}/endpoints/{name}` |
| `display_name` | `string` | Display name of the endpoint |
| `dedicated_endpoint_dns` | `string` | DNS of the dedicated endpoint. Populated only when `dedicatedEndpointEnabled` is `true`. Format: `https://{endpointId}.{region}-{projectNumber}.prediction.vertexai.goog` |
| `create_time` | `string` | RFC3339 timestamp of when the endpoint was created |
| `endpoint_name` | `string` | The numeric endpoint ID (explicit or identity-derived) — the value model-deployment tooling passes as the endpoint reference |

## Related Components

- [GcpProject](/docs/catalog/gcp/project) — provides the GCP project for the endpoint
- [GcpVpcNetwork](/docs/catalog/gcp/vpc) — provides the VPC network for VPC-peered private endpoints
- [GcpKmsKey](/docs/catalog/gcp/kms-key) — provides the encryption key for CMEK
- [GcpVertexAiNotebook](/docs/catalog/gcp/vertex-ai-notebook) — commonly co-deployed for ML development workflows
- [GcpServiceAccount](/docs/catalog/gcp/service-account) — provides the service identity for model serving
