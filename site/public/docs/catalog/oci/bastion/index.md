---
title: "Bastion"
description: "Bastion deployment documentation"
icon: "package"
order: 100
componentName: "ocibastion"
---

# Bastion on OCI

Deploys an OCI Bastion -- a managed SSH gateway that provides secure, time-limited access to resources in private subnets without requiring a public IP on the target. The bastion creates a private endpoint in the target subnet and allows authorized clients (by CIDR) to establish managed SSH sessions, port forwarding, and optional SOCKS5 dynamic port forwarding through it. Integrates with Planton's Provider Connections for OCI credential management and ValueFromRef for compartment and subnet wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Bastion** -- a managed STANDARD bastion with a private endpoint in the target subnet, configurable client CIDR allow list, and maximum session TTL enforcement
- **Freeform Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the bastion

## Before You Deploy

### Planton Setup

- **OCI Provider Connection** -- an active connection in the Connect module with credentials for the target OCI tenancy. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### OCI Tenancy

- A compartment to place the bastion in. Provide the compartment OCID directly or reference an OciCompartment Cloud Resource via ValueFromRef.
- A private subnet that the bastion will connect to. The bastion creates its endpoint in this subnet and can reach any resource accessible from it, including resources in peered VCNs if route rules are configured. Provide the subnet OCID directly or reference an OciSubnet Cloud Resource via ValueFromRef.

## Deploy

### Console

Open the deployment store, find **Bastion on OCI**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Standard SSH Gateway** preset in the [Presets](#presets) tab to pre-populate a bastion with a 3-hour session TTL.

### CLI

```yaml
apiVersion: oci.planton.dev/v1
kind: OciBastion
metadata:
  name: platform-bastion
  org: acme-corp
  env: prod
spec:
  compartmentId:
    value: "ocid1.compartment.oc1..example"
  targetSubnetId:
    value: "ocid1.subnet.oc1..example"
  clientCidrBlockAllowList:
    - "10.0.0.0/16"
  maxSessionTtlInSeconds: 10800
```

```shell
planton apply -f bastion.yaml
```

This creates a bastion with a 3-hour maximum session TTL, restricted to clients in the 10.0.0.0/16 range. DNS proxy is not enabled. Sessions are created on-demand via the OCI Console or CLI -- they are ephemeral operational artifacts, not infrastructure managed by this component.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the bastion to a compartment and subnet deployed in the same InfraPipeline:

```yaml
spec:
  compartmentId:
    valueFrom:
      kind: OciCompartment
      name: platform-compartment
      fieldPath: status.outputs.compartmentId
  targetSubnetId:
    valueFrom:
      kind: OciSubnet
      name: private-app-subnet
      fieldPath: status.outputs.subnetId
```

The InfraPipeline resolves the dependency graph, deploys the compartment and subnet first, then provisions the bastion with the resolved values.

## Key Configuration

These are the most important decisions when configuring a bastion. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Target subnet** -- The `targetSubnetId` determines which subnet the bastion creates its private endpoint in. Sessions can reach any resource accessible from this subnet. This field is immutable after creation -- changing it requires recreating the bastion.

**Client CIDR allow list** -- The `clientCidrBlockAllowList` restricts which source IPs can create and connect to sessions. Use your corporate VPN range, office IP ranges, or CI/CD runner subnets. This is the primary access control mechanism and can be updated after creation without recreation.

**Maximum session TTL** -- The `maxSessionTtlInSeconds` sets the upper bound for any session on this bastion. The OCI default of 10800 (3 hours) balances convenience with security. Reduce to 1800 (30 minutes) for high-security environments or increase to 28800 (8 hours) for extended maintenance windows. Updatable after creation.

**DNS proxy** -- Set `isDnsProxyEnabled: true` to allow sessions to target resources by FQDN instead of IP address and to enable SOCKS5 dynamic port forwarding. Essential for environments with dynamic IPs or DNS-based service discovery. This field is immutable after creation.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OciCompartment** | `compartmentId` | `status.outputs.compartmentId` |
| **OciSubnet** | `targetSubnetId` | `status.outputs.subnetId` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `bastion_id` | OCID of the bastion | OCI CLI session creation, monitoring dashboards |
| `private_endpoint_ip_address` | Private IP address of the bastion's endpoint in the target subnet | Network troubleshooting, connectivity validation |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard SSH gateway** -- A bastion with CIDR-restricted client access and 3-hour session TTL. Covers the majority of use cases: SSH to private compute instances, port forwarding to private databases, and secure access to private API endpoints. Start from the **Standard SSH Gateway** preset.

**DNS proxy enabled** -- A bastion with FQDN-based target resolution and SOCKS5 proxy support. Use when targets have dynamic IPs (auto-scaling groups, container instances) or when browsing multiple private web UIs through a single session. Start from the **DNS Proxy Enabled** preset.

## Works With

- [**Compartment on OCI**](/cloud-catalog/oci-compartment) -- provides the compartment that scopes this bastion
- [**Subnet on OCI**](/cloud-catalog/oci-subnet) -- provides the private subnet where the bastion creates its endpoint