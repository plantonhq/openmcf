---
title: "Private DNS Zone"
description: "Private DNS Zone deployment documentation"
icon: "package"
order: 100
componentName: "alicloudprivatednszone"
---

# AliCloud Private DNS Zone

Deploys an Alibaba Cloud Private Zone (PVTZ) with bundled VPC attachments and DNS records. Private Zones provide VPC-internal DNS resolution -- records are only visible to resources inside the attached VPCs and are never served to the public internet. The zone, VPC attachments, and records are deployed as a single atomic unit because a Private Zone without at least one VPC attachment has no resolver scope.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Private Zone** -- an `alicloud_pvtz_zone` resource representing the private DNS hosted zone
- **VPC Attachments** -- one `alicloud_pvtz_zone_attachment` per entry in `vpcAttachments`, making the zone resolvable within each attached VPC
- **DNS Records** -- one `alicloud_pvtz_zone_record` per entry in `records`, providing internal name resolution for services, databases, and other private endpoints

## Before You Deploy

### Planton Setup

- **AliCloud Provider Connection** -- an active connection in the Connect module with credentials for the target Alibaba Cloud account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery.

### Alibaba Cloud Account

- **At least one VPC** -- the Private Zone must be attached to a VPC for records to be queryable. Create one with AliCloudVpc.
- **Zone name planning** -- the `zoneName` is immutable after creation. Use a naming convention like "internal.example.com" or "db.corp".
- **Record types** -- Private Zone supports A, CNAME, MX, PTR, SRV, and TXT records. AAAA and NS are not supported.

## Deploy

### Console

Open the deployment store, find **AliCloud Private DNS Zone**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields including zone name, VPC attachments, and DNS records.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: alicloud.planton.dev/v1
kind: AliCloudPrivateDnsZone
metadata:
  name: internal-zone
  org: acme-corp
  env: prod
spec:
  region: cn-hangzhou
  zoneName: internal.acme.com
  remark: Internal service discovery zone
  vpcAttachments:
    - vpcId:
        value: vpc-bp1234567890
  records:
    - rr: db-master
      type: A
      value: "10.0.1.100"
    - rr: cache
      type: A
      value: "10.0.1.200"
```

```shell
planton apply -f alicloud-private-dns.yaml
```

This creates a private zone with two A records resolvable within the attached VPC. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a service infrastructure stack, use ValueFromRef to wire VPC dependencies:

```yaml
spec:
  region: cn-hangzhou
  zoneName: internal.acme.com
  vpcAttachments:
    - vpcId:
        valueFrom:
          kind: AliCloudVpc
          name: platform-vpc
          fieldPath: status.outputs.vpc_id
  records:
    - rr: db-master
      type: A
      value: "10.0.1.100"
```

The InfraPipeline resolves the dependency graph and provisions the VPC before the Private Zone attachment.

## Key Configuration

These are the most important decisions when configuring a Private DNS Zone. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Zone name** -- The `zoneName` field defines the DNS suffix for all records (e.g., "internal.acme.com" makes records resolve as "db.internal.acme.com"). Immutable after creation.

**VPC attachments** -- Each `vpcAttachments` entry binds the zone to a VPC. Cross-region attachments are supported by setting `regionId` on individual entries. At least one attachment is required.

**DNS records** -- Each `records` entry creates an internal DNS record. Supported types are A, CNAME, MX, PTR, SRV, and TXT. Default TTL is 60 seconds.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AliCloudVpc** | `vpcAttachments[].vpcId` | `status.outputs.vpc_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `zone_id` | Private Zone ID assigned by Alibaba Cloud | API references |
| `zone_name` | Zone name as created | Application DNS configuration |
| `is_ptr` | Whether this is a reverse-lookup (PTR) zone | Zone type detection |
| `record_count` | Number of DNS records in the zone | Monitoring |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Internal service discovery** -- A single-VPC Private Zone with A records for internal services. Start from the **Internal Service Discovery** preset.

**Multi-VPC database zone** -- A Private Zone attached to multiple VPCs (including cross-region) with database endpoint records and organizational tags. Start from the **Multi VPC Database Zone** preset.

## Works With

- [**AliCloud VPC**](/cloud-catalog/ali-cloud-vpc) -- the VPCs this Private Zone is attached to
- [**AliCloud RDS Instance**](/cloud-catalog/ali-cloud-rds-instance) -- database endpoints registered as A records
- [**AliCloud Redis Instance**](/cloud-catalog/ali-cloud-redis-instance) -- cache endpoints registered as A records
