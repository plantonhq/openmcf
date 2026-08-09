---
title: "IP Address"
description: "IP Address deployment documentation"
icon: "package"
order: 100
componentName: "civoipaddress"
---

# IP Address on Civo

Allocates a static reserved IPv4 address on Civo Cloud that persists independently of compute instances and can be reassigned between resources in the same region. Reserved IPs are region-specific and suitable for production workloads, high-availability failover, and services that require a stable public endpoint. Integrates with Planton's Provider Connections for Civo credential management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Civo Reserved IP** -- a static public IPv4 address allocated in the specified region, available for attachment to compute instances or load balancers

## Before You Deploy

### Planton Setup

- **Civo Provider Connection** -- an active connection in the Connect module with a Civo API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### Civo Account

- **A Civo account** with access to the target region. No additional prerequisites are required -- reserved IPs are standalone resources with no VPC or firewall dependencies.

## Deploy

### Console

Open the deployment store, find **IP Address on Civo**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Standard** preset in the [Presets](#presets) tab for a general-purpose reserved IP.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: civo.planton.dev/v1
kind: CivoIpAddress
metadata:
  name: prod-ip
  org: acme-corp
  env: prod
spec:
  region: lon1
  description: Production web server IP
```

```shell
planton apply -f civo-ip-address.yaml
```

This allocates a static reserved IPv4 address in Civo's London region. The IP is not attached to any resource initially -- attach it to a compute instance or load balancer after provisioning. A Stack Job tracks the allocation in real time.

## Key Configuration

These are the most important decisions when configuring a Civo reserved IP. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Region** -- The `region` field determines where the IP is allocated. Reserved IPs can only be attached to resources in the same region. Use region codes such as `lon1` (London), `nyc1` (New York), or `fra1` (Frankfurt). Choose the region matching your compute workloads.

**Description** -- The `description` field provides a human-readable label for identifying the IP in the Civo dashboard and API responses. Use descriptive names that indicate the purpose (e.g., "Production load balancer IP", "Bastion host IP") to simplify tracking across multiple reserved IPs.

**Attachment** -- Reserved IPs are allocated in an unattached state. After provisioning, attach the IP to a CivoComputeInstance or load balancer via the Civo dashboard or API. The IP persists independently of the attached resource, so it survives instance replacements during maintenance or failover.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `reserved_ip_id` | Unique identifier of the reserved IP in Civo | Civo API operations, instance attachment |
| `ip_address` | Static IPv4 address allocated for this reservation | DNS A records, client configuration, firewall rules |
| `attached_resource_id` | ID of the resource currently attached to this IP (empty if unattached) | Monitoring, resource tracking |
| `created_at_rfc3339` | Reservation creation timestamp in RFC 3339 format | Audit logs, lifecycle tracking |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard reserved IP** -- a static IPv4 address in the target region for production workloads that need a stable public endpoint. Use for load balancers, bastion hosts, or any service where external clients address resources by IP. Start from the **Standard** preset.

## Works With

This component operates independently and does not reference other components.