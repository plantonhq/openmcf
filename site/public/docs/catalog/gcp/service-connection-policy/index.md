---
title: "Service Connection Policy"
description: "Service Connection Policy deployment documentation"
icon: "package"
order: 100
componentName: "gcpserviceconnectionpolicy"
---

# GCP Service Connection Policy

Creates a service connection policy — the per-network authorization that lets Google's service connectivity automation create Private Service Connect endpoints in your VPC on a managed-service producer's behalf. PSC-first services (Memorystore for Valkey, Redis Cluster) require a policy for their service class in the instance's region before any instance can be created.

## What Gets Created

One policy per (network, service class, region) triple. When a producer instance is created, the automation places PSC forwarding rules in the listed subnets and the instance's endpoint IPs surface on the instance.

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **A VPC network** — referenced via `network` (a `GcpVpcNetwork` resource or a literal resource path)
- **At least one subnet** in the policy's region — referenced in `pscConfig.subnetworks`
- **IAM permissions** — `networkconnectivity.serviceConnectionPolicies.create` (e.g. `roles/networkconnectivity.consumerNetworkAdmin`)

## Quick Start

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpServiceConnectionPolicy
metadata:
  name: memorystore-valkey-policy
spec:
  location: us-central1
  network:
    valueFrom:
      kind: GcpVpcNetwork
      name: prod-vpc
      fieldPath: status.outputs.network_id
  serviceClass: gcp-memorystore
  pscConfig:
    subnetworks:
      - valueFrom:
          kind: GcpSubnetwork
          name: prod-subnet
          fieldPath: status.outputs.subnetwork_self_link
```

```shell
planton apply -f scp.yaml
```

## Configuration Reference

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `location` | `string` | — (required) | Region the policy applies to. Immutable. |
| `network` | `StringValueOrRef` | — (required) | Consumer VPC network. Can reference a GcpVpcNetwork. Immutable. |
| `serviceClass` | `string` | — (required) | The producer's published class (e.g. `gcp-memorystore`). Immutable. |
| `pscConfig.subnetworks` | `StringValueOrRef[]` | — (min 1 when set) | Subnets endpoint IPs come from. Mutable. |
| `pscConfig.limit` | `int32` | GCP default | Max PSC connections under this policy. Mutable. |
| `policyName` | `string` | `metadata.name` | Policy resource name. Immutable. |
| `projectId` | `StringValueOrRef` | provider default | Project owning the network and policy. |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `policy_id` | Fully qualified policy resource path |
| `name` | Short policy name |
| `infrastructure` | Underlying connectivity mechanism (PSC) |
| `etag` | Server-computed etag |

## Related Components

- [GcpVpcNetwork](/docs/catalog/gcp/vpc) — the network the policy authorizes connections into
- [GcpSubnetwork](/docs/catalog/gcp/subnetwork) — supplies the endpoint IP space
- [GcpMemorystoreInstance](/docs/catalog/gcp/memorystore-instance) — the first PSC-first consumer of this policy
