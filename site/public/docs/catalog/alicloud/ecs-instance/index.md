---
title: "ECS Instance"
description: "ECS Instance deployment documentation"
icon: "package"
order: 100
componentName: "alicloudecsinstance"
---

# AliCloud ECS Instance

Deploys an Elastic Compute Service (ECS) instance on Alibaba Cloud with configurable instance types, system and data disks, optional public IP allocation, spot pricing, cloud-init user data, and RAM role attachment. The instance integrates with Planton's Provider Connections for AliCloud credential management and supports ValueFromRef wiring to VSwitches and security groups for network placement.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **ECS Instance** -- an `alicloud_instance` with the selected instance type, OS image, system disk, and network configuration, placed in the specified VSwitch
- **System Disk** -- the boot disk for the instance with configurable category (cloud_essd by default), size, ESSD performance level, and optional KMS encryption
- **Data Disks** -- created only when `dataDisks` entries are provided; up to 16 additional disks attached inline with the instance, each with independent category, size, encryption, and lifecycle settings
- **Public IP** -- allocated only when `internetMaxBandwidthOut` is greater than 0; provides internet-facing connectivity with configurable bandwidth and billing method
- **AliCloud Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) merged with user-provided `tags`

## Before You Deploy

### Planton Setup

- **AliCloud Provider Connection** -- an active connection in the Connect module with credentials for the target Alibaba Cloud account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### Alibaba Cloud Account

- **A VSwitch** in the target region and availability zone. The ECS instance inherits its VPC and AZ from the VSwitch. Provide the VSwitch ID directly or reference an AliCloudVswitch Cloud Resource via ValueFromRef.
- **At least one security group** -- the instance must belong to 1-5 security groups. Provide security group IDs directly or reference AliCloudSecurityGroup Cloud Resources via ValueFromRef.
- **An SSH key pair** (recommended) or password for instance authentication. SSH key pairs provide passwordless login and are mutually exclusive with the password field.
- **An OS image ID** -- determines the operating system and architecture (e.g., `ubuntu_22_04_x64_20G_alibase_20230515.vhd`).

## Deploy

### Console

Open the deployment store, find **AliCloud ECS Instance**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Basic Development** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: alicloud.planton.dev/v1
kind: AliCloudEcsInstance
metadata:
  name: web-server
  org: acme-corp
  env: prod
spec:
  region: cn-hangzhou
  vswitchId:
    value: "vsw-abc123"
  securityGroupIds:
    - value: "sg-abc123"
  instanceType: ecs.g7.large
  imageId: ubuntu_22_04_x64_20G_alibase_20230515.vhd
  keyName: my-ssh-key
```

```shell
planton apply -f ecs-instance.yaml
```

This creates an ECS instance with a 40 GB cloud ESSD system disk, no public IP, no data disks, PostPaid billing, and SSH key authentication. Spot pricing, cloud-init user data, RAM role, and deletion protection are not configured.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the ECS instance to a VSwitch and security group deployed in the same InfraPipeline:

```yaml
spec:
  vswitchId:
    valueFrom:
      kind: AliCloudVswitch
      name: app-vswitch
      fieldPath: status.outputs.vswitch_id
  securityGroupIds:
    - valueFrom:
        kind: AliCloudSecurityGroup
        name: app-sg
        fieldPath: status.outputs.security_group_id
```

The InfraPipeline resolves the dependency graph, deploys the VSwitch and security group first, then provisions the ECS instance with the resolved values.

## Key Configuration

These are the most important decisions when configuring an ECS instance. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Instance type** -- Determines CPU and memory allocation. Must start with `ecs.` (e.g., `ecs.g7.large` for 2 vCPU/8 GiB general purpose, `ecs.c7.xlarge` for 4 vCPU/8 GiB compute optimized). The instance type also determines available disk types and network bandwidth caps.

**Public IP and bandwidth** -- Set `internetMaxBandwidthOut` to a value greater than 0 (up to 100 Mbps) to allocate a public IP address. When `internetMaxBandwidthOut` is 0 (default), the instance has no public IP. Choose `internetChargeType: PayByTraffic` for variable workloads or `PayByBandwidth` for predictable usage.

**Spot instances** -- Set `spotStrategy` to `SpotAsPriceGo` (market price) or `SpotWithPriceLimit` (with `spotPriceLimit` as a maximum hourly price cap). Spot instances can reduce costs up to 90% but may be reclaimed when capacity is needed. Suitable for stateless, fault-tolerant workloads.

**Disk composition** -- The system disk defaults to 40 GB cloud ESSD. Configure `systemDisk` to change category, size, or ESSD performance level (PL0-PL3). Add `dataDisks` for application storage, databases, or logs. Each data disk supports independent category, size, encryption, and snapshot-based creation. Set `deleteWithInstance: false` on data disks to retain them when the instance is released.

**Authentication** -- Use `keyName` for SSH key pair authentication (recommended) or `password` for password-based login. These are mutually exclusive. Attach a `roleName` to grant the instance an IAM profile for calling Alibaba Cloud APIs without embedded credentials.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AliCloudVswitch** | `vswitchId` | `status.outputs.vswitch_id` |
| **AliCloudSecurityGroup** | `securityGroupIds` | `status.outputs.security_group_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `instance_id` | ECS instance ID (e.g., `i-bp1xxxxx`) | Monitoring dashboards, audit references |
| `private_ip` | Primary private IP address within the VSwitch | Application connection strings, security group rules, DNS records |
| `public_ip` | Public IP address (empty when `internetMaxBandwidthOut` is 0) | DNS records, external access configuration |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Basic development** -- A small instance type with minimal system disk, SSH key authentication, no public IP, and PostPaid billing for development and testing. Start from the **Basic Development** preset.

**Production web server** -- A general-purpose instance with public IP, larger system disk, deletion protection enabled, and cloud-init user data for automated provisioning. Start from the **Production Web Server** preset.

**Spot batch worker** -- A spot instance with `SpotAsPriceGo` strategy for cost-optimized batch processing or CI/CD workloads. Suitable for stateless, fault-tolerant processing. Start from the **Spot Batch Worker** preset.

## Works With

- [**AliCloud VSwitch**](/cloud-catalog/ali-cloud-vswitch) -- provides the VSwitch for VPC and availability zone placement
- [**AliCloud Security Group**](/cloud-catalog/ali-cloud-security-group) -- provides network access control for the instance