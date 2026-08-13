# GCP Vertex AI Index Endpoint

Deploys a Vertex AI Vector Search index endpoint: the serving surface deployed indexes answer nearest-neighbor queries through. The endpoint owns connectivity — public, VPC-peered, or Private Service Connect — while the indexes themselves are separate GcpVertexAiIndex resources placed onto it by GcpVertexAiDeployedIndex. This is a different GCP resource from the online-prediction GcpVertexAiEndpoint (which serves models). The component integrates with Planton's Provider Connections for GCP credential management and supports ValueFromRef wiring to GCP projects and VPC networks.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Vertex AI Index Endpoint** -- a regional Vector Search endpoint with the configured connectivity mode
- **Public Serving Domain** -- when `publicEndpointEnabled` is true, a GCP-managed public domain name (surfaced as the `public_endpoint_domain_name` output) that authenticated clients query over the internet
- **VPC Peering Attachment** -- when `network` is set, the endpoint becomes reachable only inside the peered VPC via Private Services Access; both IaC modules normalize the network self-link to the relative form the Vertex AI API expects
- **PSC Service Attachment Surface** -- when `privateServiceConnectConfig` is present, allowlisted consumer projects create forwarding rules to each deployment's service attachment; no peering needed

## Before You Deploy

### Planton Setup

- **GCP Provider Connection** -- an active connection in the Connect module with credentials for the target GCP project. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or browser OAuth authentication modes.

### GCP Project

- **A GCP project** where the endpoint will be created. Provide the project ID directly or reference a GcpProject Cloud Resource via ValueFromRef.
- **Vertex AI API** enabled in the target project.
- **For VPC-peered serving** -- Private Services Access configured on the network first: a GcpGlobalAddress (purpose VPC_PEERING) composed with a GcpServiceNetworkingConnection.
- **For PSC serving** -- the list of consumer project IDs, decided up front (the allowlist is immutable).

## Deploy

### Console

Open the deployment store, find **GCP Vertex AI Index Endpoint**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Public** preset in the [Presets](#presets) tab for the simplest path to production.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpVertexAiIndexEndpoint
metadata:
  name: catalog-search
  org: acme-corp
  env: prod
spec:
  projectId:
    value: "acme-prod-12345"
  location: us-central1
  displayName: catalog-search-endpoint
  publicEndpointEnabled: true
```

```shell
planton apply -f index-endpoint.yaml
```

A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying the full vector-search composition, deploy the endpoint (and index) before the deployment that joins them:

```yaml
# In a GcpVertexAiDeployedIndex spec:
indexEndpoint:
  valueFrom:
    kind: GcpVertexAiIndexEndpoint
    name: catalog-search
    fieldPath: status.outputs.index_endpoint_id
```

The InfraPipeline resolves the dependency graph, provisions the endpoint first, then the deployment that serves through it.

## Key Configuration

These are the most important decisions when configuring an index endpoint. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Connectivity (the one-way door)** -- exactly one of `publicEndpointEnabled: true`, `network`, or `privateServiceConnectConfig`. All three are mutually exclusive and every choice is immutable — switching modes means recreating the endpoint and re-deploying every index on it.

**Public serving** -- queries hit a GCP-managed domain over the internet, authenticated with Google Cloud credentials. Public reachability, not public access — IAM still gates every query.

**VPC peering** -- `network` references a GcpVpcNetwork (its `network_self_link` output). Requires Private Services Access on the network before the endpoint deploys. Private clients call each deployment's `match_grpc_address`.

**Private Service Connect** -- `projectAllowlist` names the consumer projects allowed to create forwarding rules; include every consumer up front (immutable). Consumers connect through each deployment's `service_attachment`.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **GcpProject** (optional) | `projectId` | `status.outputs.project_id` |
| **GcpVpcNetwork** (optional — VPC-peered mode) | `network` | `status.outputs.network_self_link` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `index_endpoint_id` | Fully qualified endpoint path (`projects/{p}/locations/{l}/indexEndpoints/{id}`) | The exact value a GcpVertexAiDeployedIndex's `indexEndpoint` join consumes |
| `index_endpoint_name` | The GCP-assigned numeric endpoint ID | Display, logging |
| `public_endpoint_domain_name` | Public query domain (populated only when public serving is enabled) | Query clients over the internet |
| `create_time` | RFC3339 creation timestamp | Audit |
| `update_time` | RFC3339 last-update timestamp | Audit |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Public serving** -- The simplest path to production: one flag, a public domain output, IAM-authenticated queries. Start from the **Public** preset.

**VPC-peered private serving** -- Query traffic never leaves private address space; compose with the Private Services Access pair first. Start from the **VPC Peered** preset.

**Private Service Connect** -- Private connectivity across project boundaries without peering — the modern model for multi-project estates. Start from the **PSC** preset.

## Works With

- [**GCP Project**](/cloud-catalog/gcp-project) -- provides the GCP project where the endpoint is created
- [**GCP VPC Network**](/cloud-catalog/gcp-vpc-network) -- the VPC a peered endpoint serves inside
- [**GCP Global Address**](/cloud-catalog/gcp-global-address) + [**GCP Service Networking Connection**](/cloud-catalog/gcp-service-networking-connection) -- the Private Services Access pair a peered endpoint requires
- [**GCP Vertex AI Index**](/cloud-catalog/gcp-vertex-ai-index) -- the vector index deployed onto this endpoint
- [**GCP Vertex AI Deployed Index**](/cloud-catalog/gcp-vertex-ai-deployed-index) -- joins an index to this endpoint via `indexEndpoint`
