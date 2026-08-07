# GCP Serverless VPC Access Connector

Deploys a Serverless VPC Access connector — the managed bridge that lets serverless workloads ([GcpCloudFunction](/docs/catalog/gcp/gcpcloudfunction), [GcpCloudRun](/docs/catalog/gcp/gcpcloudrun), [GcpCloudRunJob](/docs/catalog/gcp/gcpcloudrunjob), App Engine) send traffic into a VPC network. The connector runs a small fleet of forwarding instances inside the VPC; serverless egress configured to use it reaches private IPs (Cloud SQL private IP, Memorystore, internal load balancers) as if the workload lived in the network.

## What Gets Created

When you deploy a GcpServerlessVpcConnector resource, Planton provisions:

- **Serverless VPC Access connector** — a `google_vpc_access_connector` in the chosen region; the Serverless VPC Access API is enabled automatically
- **A managed instance fleet** — GCP-operated forwarding instances (sized by `machineType` and the `minInstances`–`maxInstances` band) occupying a dedicated `/28`

One connector serves many functions and services in its region — it is shared infrastructure, not per-workload. Note that Cloud Run and Cloud Run jobs also support **Direct VPC egress** (no connector needed); the connector remains the mechanism for Cloud Functions and App Engine, and for organizations whose constraints require it.

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **A VPC network** ([GcpVpcNetwork](/docs/catalog/gcp/gcpvpcnetwork)) — plus an unused `/28` range, or
- **A dedicated /28 subnetwork** ([GcpSubnetwork](/docs/catalog/gcp/gcpsubnetwork)) for subnet placement (required on Shared VPC)

## Quick Start

Create a file `connector.yaml`:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpServerlessVpcConnector
metadata:
  name: svc-egress
spec:
  region: us-central1
  network:
    value: my-vpc
  ipCidrRange: 10.8.0.0/28
```

Deploy:

```shell
planton apply -f connector.yaml
```

Attach it from a Cloud Function:

```yaml
spec:
  serviceConfig:
    vpcConnector:
      valueFrom:
        name: svc-egress
    vpcConnectorEgressSettings: PRIVATE_RANGES_ONLY
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | Region the connector lives in. Serverless workloads can only use a connector in their own region. Immutable. | Required |
| placement | — | Exactly one of `network` (+`ipCidrRange`) or `subnet`. | Enforced pre-deploy |

### Identity

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `projectId` | `StringValueOrRef` | provider default | GCP project. Can reference GcpProject. |
| `connectorName` | `string` | `metadata.name` | Connector name (max **25** characters — a GCP limit stricter than most resources). Immutable. |

### Placement (exactly one)

| Field | Type | Description |
|-------|------|-------------|
| `network` | `StringValueOrRef` | VPC to attach to; the connector carves `ipCidrRange` out of it. Can reference GcpVpcNetwork. Immutable. |
| `ipCidrRange` | `string` | Unused range, **exactly /28** (e.g. `10.8.0.0/28`); must overlap no subnet, peered range, or route. Immutable. |
| `subnet.name` | `StringValueOrRef` | Existing dedicated `/28` subnetwork (short name). Can reference GcpSubnetwork. Required on Shared VPC. Immutable. |
| `subnet.projectId` | `string` | Project owning the subnet (Shared VPC host project). Defaults to the connector's project. |

### Capacity

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `machineType` | `string` | `e2-micro` | `f1-micro` (~100 Mbps/instance), `e2-micro` (~200 Mbps), `e2-standard-4` (~1 Gbps class). Mutable in place. |
| `minInstances` | `int` | GCP default (2) | Floor, 2–9. Must be < `maxInstances`. **Decreasing replaces the connector.** |
| `maxInstances` | `int` | GCP default (10) | Ceiling, 3–10. Must be > `minInstances`. **Decreasing replaces the connector.** |

Two behaviors worth planning around: the fleet **never scales in on its own** — after a traffic burst it stays at the high-water mark; and while *increases* to the instance band apply in place, *decreases* force the connector to be replaced (a brief egress outage for every workload using it).

The legacy throughput-based scaling fields (`min_throughput`/`max_throughput`) are deliberately not modeled: the provider discourages them in favor of the instance band, they conflict with the instance fields, and changing them always replaces the connector.

## Outputs

| Output | Description |
|--------|-------------|
| `name` | Short connector name |
| `selfLink` | Fully qualified resource name (`projects/*/locations/*/connectors/*`) — what serverless workloads attach to |
| `state` | `READY`, `CREATING`, `DELETING`, `ERROR`, `UPDATING` |
| `region` | Plain region name |

## Presets

- [Private egress — basic](presets/01-private-egress-basic.yaml) — the standard shared connector
- [High throughput — production](presets/02-high-throughput.yaml) — `e2-standard-4` with a wide instance band
- [Shared VPC — subnet placement](presets/03-shared-vpc-subnet.yaml) — host-project subnet mode

## See Also

- [GcpCloudFunction](/docs/catalog/gcp/gcpcloudfunction) — attaches by `vpcConnector` reference
- [GcpCloudRun](/docs/catalog/gcp/gcpcloudrun) / [GcpCloudRunJob](/docs/catalog/gcp/gcpcloudrunjob) — attach by reference, or use Direct VPC egress instead
- [GcpVpcNetwork](/docs/catalog/gcp/gcpvpcnetwork), [GcpSubnetwork](/docs/catalog/gcp/gcpsubnetwork) — the placement targets

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
