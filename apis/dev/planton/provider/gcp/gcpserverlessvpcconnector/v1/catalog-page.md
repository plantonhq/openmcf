# GCP Serverless VPC Access Connector

The managed bridge from serverless into a VPC: Cloud Functions, Cloud Run, and App Engine route egress through the connector to reach private IPs (Cloud SQL private IP, Memorystore, internal load balancers).

**Enum:** 721 · **ID prefix:** `gcpvpcconn` · **Provider:** GCP · **API:** `gcp.planton.dev/v1`

## At a Glance

| | |
|---|---|
| **Creates** | `google_vpc_access_connector` |
| **Consumed by** | GcpCloudFunction, GcpCloudRun, GcpCloudRunJob (by `selfLink` reference) |
| **Placement** | network + /28 range, or a dedicated /28 subnet (Shared VPC) |
| **Engines** | Terraform (~> 6.0) and Pulumi |

## When to Use

- **Cloud Functions reaching private resources** — the connector is the only VPC path for functions
- **Org-mandated connector egress** — Cloud Run supports Direct VPC egress, but some constraints still require connectors
- **Shared VPC serverless egress** — subnet placement against a host-project `/28`

## Quick Example

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpServerlessVpcConnector
metadata:
  name: svc-egress
spec:
  region: us-central1
  network:
    value: my-vpc
  ipCidrRange: 10.8.0.0/28
```

## Key Fields

- `network` + `ipCidrRange` — carve a new dedicated `/28` out of the VPC
- `subnet` — occupy an existing dedicated `/28` (the Shared-VPC mode); exactly one placement
- `machineType` — per-instance throughput class; changes in place
- `minInstances` / `maxInstances` — the scaling band; decreases replace the connector

## Outputs

`name`, `selfLink`, `state`, `region`

## See Also

- [README](README.md) — full configuration reference
- [GcpCloudFunction](/docs/catalog/gcp/gcpcloudfunction) — the primary consumer
- Presets: [basic](presets/01-private-egress-basic.yaml), [high throughput](presets/02-high-throughput.yaml), [Shared VPC](presets/03-shared-vpc-subnet.yaml)
