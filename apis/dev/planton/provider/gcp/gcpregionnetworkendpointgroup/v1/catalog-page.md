# GCP Region Network Endpoint Group

Creates a regional network endpoint group (NEG) — the bridge that puts a serverless workload (Cloud Run, Cloud Functions, App Engine), a Private Service Connect endpoint, or an external origin behind a Google Cloud load balancer. A backend service references the NEG's self-link in `backends[].group`.

## What Gets Created

A single regional NEG whose target is decided by `networkEndpointType`: `SERVERLESS` (a `cloudRun` / `cloudFunction` / `appEngine` block), `PRIVATE_SERVICE_CONNECT`, `INTERNET_IP_PORT` / `INTERNET_FQDN_PORT`, or `GCE_VM_IP_PORTMAP`.

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **An existing GCP project** — referenced via `projectId` (or the provider's default project)
- **IAM permissions** — any role carrying `compute.regionNetworkEndpointGroups.*` on the target project

## Quick Start

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpRegionNetworkEndpointGroup
metadata:
  name: web-neg
spec:
  projectId:
    value: my-gcp-project-123
  region: us-central1
  cloudRun:
    service:
      value: my-cloud-run-service
```

```shell
planton apply -f neg.yaml
```

## Configuration Reference

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `region` | `string` | — (required) | Region the NEG lives in. Immutable. |
| `networkEndpointType` | `string` | `SERVERLESS` | SERVERLESS / PRIVATE_SERVICE_CONNECT / INTERNET_IP_PORT / INTERNET_FQDN_PORT / GCE_VM_IP_PORTMAP. Immutable. |
| `cloudRun` / `cloudFunction` / `appEngine` | object | — | Serverless target — exactly one for a SERVERLESS NEG. |
| `pscTargetService` | `string` | — | Target service URL for PSC/INTERNET NEGs (required for PSC). |
| `network` / `subnetwork` | `StringValueOrRef` | — | VPC/subnet for PSC/INTERNET/PORTMAP NEGs. |
| `projectId` | `StringValueOrRef` | provider default | Project that owns the NEG. Immutable. |
| `networkEndpointGroupName` | `string` | `metadata.name` | Cloud-side name (RFC1035). Immutable. |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `self_link` | Self-link URI — the value a backend service references in `backends[].group` |
| `network_endpoint_group_name` | Name of the NEG in GCP |
| `network_endpoint_type` | The NEG's endpoint type |
| `region` | Region the NEG lives in |

## Related Components

- [GcpBackendService](/docs/catalog/gcp/gcpbackendservice) — references this NEG in `backends[].group`
- [GcpCloudRun](/docs/catalog/gcp/gcpcloudrun) — the serverless workload a SERVERLESS NEG fronts
- [GcpVpcNetwork](/docs/catalog/gcp/gcpvpcnetwork) — the network for PSC/INTERNET NEGs
