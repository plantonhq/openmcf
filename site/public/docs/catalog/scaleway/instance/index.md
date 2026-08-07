---
title: "Instance"
description: "Instance deployment documentation"
icon: "package"
order: 100
componentName: "scalewayinstance"
---

# Scaleway Instance

Deploys a compute instance on Scaleway as a composite resource that bundles the server, an optional dedicated Flexible IP (public IPv4), optional additional local volumes, and an optional Private Network attachment into a single declarative unit. Configurable cloud-init bootstrapping, deletion protection, and security group assignment provide production-ready instance management. Supports ValueFromRef for Private Network and security group dependency wiring in InfraCharts.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Instance Server** -- a zonal compute instance with the configured instance type, OS image, root volume, and optional cloud-init script
- **Flexible IP** -- created only when `publicIp` is set; a dedicated public IPv4 address with independent lifecycle that survives instance replacement
- **Additional Local Volumes** -- created only when `additionalVolumes` entries are provided; local SSD or scratch volumes attached to the instance
- **Private Network NIC** -- created only when `privateNetworkId` is set; an inline network interface on the specified Private Network for private connectivity
- **Scaleway Tags** -- resource metadata tags (resource name, kind, organization, environment) applied automatically for tracking and governance

## Before You Deploy

### Scaleway Account

- **A Scaleway account** with an active project and API access key pair.
- **A Private Network** in the target region for internal connectivity. Provide the Private Network UUID directly or reference a ScalewayPrivateNetwork Cloud Resource via ValueFromRef. Optional for isolated development instances.
- **A security group** with appropriate inbound/outbound rules. Provide the security group UUID directly or reference a ScalewayInstanceSecurityGroup Cloud Resource via ValueFromRef. If omitted, Scaleway assigns its default allow-all security group.

## Deploy

### Console

Open the deployment store, find **Scaleway Instance**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Development Instance** preset in the [Presets](#presets) tab to get a small instance with a public IP running quickly.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: scaleway.planton.dev/v1
kind: ScalewayInstance
metadata:
  name: web-01
  org: acme-corp
  env: prod
spec:
  zone: fr-par-1
  type: DEV1-S
  image: ubuntu_jammy
  publicIp: {}
```

```shell
planton apply -f scaleway-instance.yaml
```

This creates a DEV1-S instance with Ubuntu 22.04 and a public IP in the Paris zone. No Private Network, security group, or additional volumes are configured. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the instance to a Private Network and security group deployed in the same InfraPipeline:

```yaml
spec:
  securityGroupId:
    valueFrom:
      kind: ScalewayInstanceSecurityGroup
      name: web-sg
      fieldPath: status.outputs.security_group_id
  privateNetworkId:
    valueFrom:
      kind: ScalewayPrivateNetwork
      name: app-network
      fieldPath: status.outputs.private_network_id
```

The InfraPipeline resolves the dependency graph, deploys the VPC, Private Network, and security group first, then provisions the instance with the resolved values.

## Key Configuration

These are the most important decisions when configuring an instance. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Instance type** -- The `type` field sets CPU, RAM, and local storage. Use `DEV1-S` or `DEV1-M` for development, `PRO2-S` or `PRO2-M` for production workloads with guaranteed resources. Type can be changed after creation (the instance is stopped, migrated, and restarted automatically).

**Public IP** -- Include `publicIp` to create a dedicated Flexible IP for internet-reachable instances. Omit it for production workloads behind a Load Balancer or Public Gateway -- keeping instances off the public internet reduces attack surface. The Flexible IP survives instance replacement, preserving DNS records.

**Root volume** -- Configure `rootVolume` to override the image defaults. Use `l_ssd` for high-performance local storage or `sbs_volume` for network-attached storage with snapshots and resizing. Changing `volumeType` after creation recreates the instance.

**Deletion protection** -- Set `protected` to true for production instances to prevent accidental deletion via the API. The instance cannot be terminated without first disabling protection.

**Cloud-init** -- Provide a `cloudInit` script to install packages, configure services, or join configuration management systems on first boot. Maximum size is approximately 127 KB.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **ScalewayInstanceSecurityGroup** (optional) | `securityGroupId` | `status.outputs.security_group_id` |
| **ScalewayPrivateNetwork** (optional) | `privateNetworkId` | `status.outputs.private_network_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `server_id` | Zoned identifier of the created instance | Monitoring dashboards, management API calls |
| `public_ip_address` | Public IPv4 address (empty if no public IP configured) | ScalewayDnsRecord A records, SSH access, external whitelisting |
| `public_ip_id` | Flexible IP resource identifier (empty if no public IP configured) | IP lifecycle management, reassignment to replacement instances |
| `private_ip_address` | Private IP on the attached Private Network (empty if no network configured) | Load Balancer backend servers, internal service discovery, database allowlists |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Development instance** -- A DEV1-S instance with Ubuntu 22.04 and a public IP for quick development and prototyping. No Private Network or security group for simplicity. Start from the **Development Instance** preset.

**Production private instance** -- A PRO2-S instance on a Private Network with no public IP, an explicit security group, and deletion protection enabled. Only reachable through the Private Network via a Public Gateway, Load Balancer, or VPN. Start from the **Production Private** preset.

## Works With

- [**Scaleway Instance Security Group**](/cloud-catalog/scaleway-instance-security-group) -- provides firewall rules controlling inbound and outbound traffic
- [**Scaleway Private Network**](/cloud-catalog/scaleway-private-network) -- provides private connectivity to other resources in the same network