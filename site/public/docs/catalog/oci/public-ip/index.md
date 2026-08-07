---
title: "Public IP"
description: "Public IP deployment documentation"
icon: "package"
order: 100
componentName: "ocipublicip"
---

# Public IP on OCI

Deploys a reserved or ephemeral public IPv4 address on Oracle Cloud Infrastructure. Reserved IPs are persistent and region-scoped, surviving instance termination and VNIC detachment. Ephemeral IPs are tied to the assigned entity's lifecycle and released automatically when that entity is terminated. The component supports optional assignment to a private IP at creation time and optional BYOIP (Bring Your Own IP) pool allocation. The component integrates with Planton's Provider Connections for OCI credential management and supports ValueFromRef wiring to compartments.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Public IP** -- a `core.PublicIp` with the specified lifetime (RESERVED or EPHEMERAL), optional private IP assignment, and optional BYOIP pool allocation
- **Freeform Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the public IP

## Before You Deploy

### Planton Setup

- **OCI Provider Connection** -- an active connection in the Connect module with credentials for the target OCI tenancy. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### OCI Tenancy

- A compartment to place the public IP in. For ephemeral IPs, this must match the compartment of the target private IP. Provide the compartment OCID directly or reference an OciCompartment Cloud Resource via ValueFromRef.
- For assigned IPs: a private IP OCID on an existing VNIC. For ephemeral IPs, this must be a primary private IP.
- For BYOIP: a public IP pool OCID from which the address will be allocated instead of Oracle's pool.

## Deploy

### Console

Open the deployment store, find **Public IP on OCI**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Reserved Unassigned** preset in the [Presets](#presets) tab to pre-allocate a stable IP for DNS records or firewall allowlists.

### CLI

```yaml
apiVersion: oci.planton.dev/v1
kind: OciPublicIp
metadata:
  name: bastion-ip
  org: acme-corp
  env: prod
spec:
  compartmentId:
    value: "ocid1.compartment.oc1..example"
  lifetime: RESERVED
  displayName: bastion-ip
```

```shell
planton apply -f public-ip.yaml
```

This creates a reserved public IP that is not assigned to any resource. The IP persists until explicitly deleted and can be attached to a private IP later.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the public IP to a compartment deployed in the same InfraPipeline:

```yaml
spec:
  compartmentId:
    valueFrom:
      kind: OciCompartment
      name: networking
      fieldPath: status.outputs.compartmentId
```

The InfraPipeline resolves the dependency graph, deploys the compartment first, then provisions the public IP with the resolved compartment OCID.

## Key Configuration

These are the most important decisions when configuring a public IP. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Lifetime** -- Set `lifetime` to `RESERVED` for persistent, region-scoped IPs that survive instance termination (production, DNS records, firewall allowlists). Set to `EPHEMERAL` for IPs tied to the assigned entity's lifecycle (dev/test, temporary workloads). Lifetime is immutable after creation.

**Private IP assignment** -- Set `privateIpId` to bind the public IP to a specific VNIC's private IP at creation time. Required for ephemeral IPs. Optional for reserved IPs -- when omitted, the IP is created unassigned and can be attached later. This allows pre-provisioning stable addresses before the target infrastructure exists.

**BYOIP pool** -- Set `publicIpPoolId` to allocate the address from your own IP range instead of Oracle's pool. Immutable after creation. Used when organizational policy requires specific IP ranges for external communication.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OciCompartment** | `compartmentId` | `status.outputs.compartmentId` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `public_ip_id` | OCID of the public IP resource | Network auditing, resource management |
| `ip_address` | The allocated IPv4 address (e.g., `203.0.113.2`) | DNS A-records, firewall allowlists, client configuration |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Reserved unassigned** -- A persistent public IP created without binding to any resource. Suitable for pre-provisioning stable addresses for DNS records, firewall allowlists, or partner integrations. Start from the **Reserved Unassigned** preset.

**Reserved assigned** -- A persistent public IP immediately bound to an existing private IP on a VNIC. Suitable for bastion hosts, VPN endpoints, and any resource that needs a stable public address from creation. Start from the **Reserved Assigned** preset.

## Works With

- [**OCI Compartment**](/cloud-catalog/oci-compartment) -- provides the compartment that scopes this public IP