---
title: "Service Connection Policy"
description: "Service Connection Policy deployment documentation"
icon: "package"
order: 100
componentName: "gcpserviceconnectionpolicy"
---

# GCP Service Connection Policy

Authorizes Google's service connectivity automation to place Private Service Connect endpoints in your VPC for a producer's service class — required before PSC-first managed services (Memorystore for Valkey, Redis Cluster, and others) can create instances in a region.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Service Connection Policy** — per (network, service class, region) authorization with PSC subnet address space and optional connection limit / producer hierarchy allowlist

## Before You Deploy

### Planton Setup

- **GCP Provider Connection** with credentials for the target project
- **Planton Runner** when using Runner-based credential delivery

### GCP Prerequisites

- **GcpVpcNetwork** — consumer VPC (`network_id` resource path, not self link)
- **GcpSubnetwork** — at least one subnet in the policy region for endpoint IPs

## Deploy

### Console

Open the deployment store, find **GCP Service Connection Policy**, and click **Deploy**. Start from the **Memorystore Valkey** preset in the [Presets](#presets) tab.

### CLI

```yaml
apiVersion: gcp.planton.dev/v1
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

## Key Configuration

**Identity & placement** — Policy name, region, network, and service class are create-time. One policy per (network, service class, region).

**PSC address space** — Subnets the automation draws endpoint IPs from (same region as the policy). Optional connection limit caps instance attachments.

**Producer authorization** — Default accepts any producer; custom hierarchy restricts to listed projects, folders, or organizations.

## Outputs and Dependencies

### Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **GcpProject** (optional) | `projectId` | `status.outputs.project_id` |
| **GcpVpcNetwork** | `network` | `status.outputs.network_id` |
| **GcpSubnetwork** | `pscConfig.subnetworks` | `status.outputs.subnetwork_self_link` |

### Provides

| Output | Description |
|--------|-------------|
| `policy_id` | Fully qualified policy resource path |
| `name` | Short policy name |
| `infrastructure` | Connectivity mechanism (PSC) |
| `etag` | Server-computed change token |

## Presets

**Memorystore Valkey** — Authorizes `gcp-memorystore` in one region with a single subnet. Start from the **Memorystore Valkey** preset.

**Shared VPC Guarded** — Adds a connection limit for shared-network guardrails. Start from the **Shared VPC Guarded** preset.

**Producer Allowlist** — Custom hierarchy mode restricting which producer projects may attach. Start from the **Producer Allowlist** preset.

## Works With

- [**GCP VPC Network**](/cloud-catalog/gcp-vpc-network) — authorized consumer network
- [**GCP Subnetwork**](/cloud-catalog/gcp-subnetwork) — endpoint IP address space
- [**GCP Memorystore Instance**](/cloud-catalog/gcp-memorystore-instance) — deploy after the matching policy
