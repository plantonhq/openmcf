# AliCloud CEN Instance

Deploys an Alibaba Cloud Cloud Enterprise Network (CEN) instance with bundled child-instance attachments. CEN is a global networking service that provides high-quality, low-latency private connectivity between VPCs in different regions, or between VPCs and on-premises data centers via Virtual Border Routers (VBR). Unlike most Alibaba Cloud resources, CEN is region-agnostic -- the instance itself is not bound to a single region.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **CEN Instance** -- a `alicloud_cen_instance` resource that serves as the global backbone for connecting child instances
- **CEN Attachments** -- one `alicloud_cen_instance_attachment` per entry in `attachments`, connecting VPCs, VBRs, or CCNs to the CEN instance for private inter-network communication

## Before You Deploy

### Planton Setup

- **AliCloud Provider Connection** -- an active connection in the Connect module with credentials for the target Alibaba Cloud account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery.

### Alibaba Cloud Account

- **At least two VPCs** (or VBRs/CCNs) to connect -- CEN is only useful with multiple child instances.
- **CIDR planning** -- by default, CEN rejects overlapping CIDR blocks between attached VPCs. Set `protectionLevel` to "REDUCED" if you need overlapping CIDRs with route map control.
- **Cross-account considerations** -- attaching VPCs from different Alibaba Cloud accounts requires cross-account authorization (not managed by this component).

## Deploy

### Console

Open the deployment store, find **AliCloud CEN Instance**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields including CEN name, protection level, and child-instance attachments.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: alicloud.planton.dev/v1
kind: AliCloudCenInstance
metadata:
  name: platform-cen
  org: acme-corp
  env: prod
spec:
  region: cn-hangzhou
  cenInstanceName: platform-backbone
  description: Connects production and staging VPCs
  attachments:
    - childInstanceId:
        value: vpc-prod-001
      childInstanceRegionId: cn-hangzhou
    - childInstanceId:
        value: vpc-staging-001
      childInstanceRegionId: cn-shanghai
```

```shell
planton apply -f alicloud-cen.yaml
```

This creates a CEN instance connecting two VPCs across cn-hangzhou and cn-shanghai. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-VPC topology, use ValueFromRef to wire VPC dependencies:

```yaml
spec:
  region: cn-hangzhou
  cenInstanceName: platform-backbone
  attachments:
    - childInstanceId:
        valueFrom:
          kind: AliCloudVpc
          name: prod-vpc
          fieldPath: status.outputs.vpc_id
      childInstanceRegionId: cn-hangzhou
    - childInstanceId:
        valueFrom:
          kind: AliCloudVpc
          name: staging-vpc
          fieldPath: status.outputs.vpc_id
      childInstanceRegionId: cn-shanghai
```

The InfraPipeline resolves the dependency graph and provisions VPCs before the CEN attachments.

## Key Configuration

These are the most important decisions when configuring a CEN instance. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Protection level** -- The `protectionLevel` field controls CIDR overlap behavior. Leave empty for strict mode (rejects overlaps). Set to "REDUCED" to allow overlapping CIDRs between attached VPCs -- useful when migrating workloads or using route maps for traffic steering.

**Child instance type** -- Each attachment's `childInstanceType` defaults to "VPC". Use "VBR" for Virtual Border Router connections (Express Connect to on-premises) or "CCN" for Cloud Connect Networks (Smart Access Gateway).

**Region** -- The CEN's `region` field is for API routing only. Attached child instances can reside in any region -- each attachment declares its own `childInstanceRegionId`.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AliCloudVpc** | `attachments[].childInstanceId` | `status.outputs.vpc_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `cen_id` | CEN instance ID (e.g., cen-xxxxx) | Bandwidth packages, route maps |
| `cen_instance_name` | CEN instance name as configured | Display and tagging |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Multi-VPC same region** -- Connects multiple VPCs within one region for private communication without public internet exposure. Start from the **Multi VPC Same Region** preset.

**Cross-region backbone** -- Connects VPCs across multiple regions (e.g., China and international) with REDUCED protection level for flexible routing. Start from the **Cross Region Backbone** preset.

## Works With

- [**AliCloud VPC**](/cloud-catalog/ali-cloud-vpc) -- the VPCs connected by this CEN instance
- [**AliCloud VPN Gateway**](/cloud-catalog/ali-cloud-vpn-gateway) -- alternative connectivity for site-to-site VPN
