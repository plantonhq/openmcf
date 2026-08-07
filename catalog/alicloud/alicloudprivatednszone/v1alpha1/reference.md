# AliCloudPrivateDnsZone

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `alicloud.planton.dev/v1alpha1`

AliCloudPrivateDnsZoneSpec defines the configuration for an Alibaba Cloud
Private Zone (PrivateZone / PVTZ) that provides VPC-internal DNS resolution.

A Private Zone is a private DNS hosted zone that resolves domain names
within one or more VPCs. Unlike public Alidns domains, Private Zone records
are only visible to resources inside the attached VPCs -- they are never
served to the public internet.

This component bundles three provider resources:
  1. alicloud_pvtz_zone       -- the private hosted zone
  2. alicloud_pvtz_zone_attachment -- VPC binding(s)
  3. alicloud_pvtz_zone_record     -- DNS records within the zone

At least one VPC attachment is required; without it the zone has no resolver
scope and records cannot be queried.

Provider resources:
  Terraform: alicloud_pvtz_zone + alicloud_pvtz_zone_attachment + alicloud_pvtz_zone_record
  Pulumi:    pvtz.Zone + pvtz.ZoneAttachment + pvtz.ZoneRecord

## Example

```yaml
apiVersion: alicloud.planton.dev/v1alpha1
kind: AliCloudPrivateDnsZone
metadata:
  name: test-private-zone
  org: test-org
  env: development
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: test-org
    pulumi.planton.dev/project: test-project
    pulumi.planton.dev/stack.name: dev.AliCloudPrivateDnsZone.test-private-zone
spec:
  region: cn-hangzhou
  zoneName: test.internal
  remark: Test private zone for development
  vpcAttachments:
    - vpcId:
        value: vpc-test123
  records:
    - rr: api
      type: A
      value: "10.0.1.50"
      ttl: 120
    - rr: db
      type: A
      value: "10.0.2.100"
    - rr: cache
      type: CNAME
      value: redis.test.internal
  tags:
    team: platform
    purpose: testing
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.zoneName` | `string` | yes |  |  |
| `spec.remark` | `string` |  |  |  |
| `spec.resourceGroupId` | `string` |  |  |  |
| `spec.vpcAttachments` | `[]AliCloudPrivateDnsZoneVpcAttachment` | yes |  |  |
| `spec.vpcAttachments[].vpcId` | `string \| valueFrom` | yes |  | AliCloudVpc (`status.outputs.vpc_id`) |
| `spec.vpcAttachments[].regionId` | `string` |  |  |  |
| `spec.records` | `[]AliCloudPrivateDnsZoneRecord` |  |  |  |
| `spec.records[].rr` | `string` | yes |  |  |
| `spec.records[].type` | `string` | yes |  |  |
| `spec.records[].value` | `string` | yes |  |  |
| `spec.records[].ttl` | `int32` |  |  |  |
| `spec.records[].priority` | `int32` |  |  |  |
| `spec.records[].remark` | `string` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.region

`string` · required

Alibaba Cloud region for provider initialization.
Private Zone is a global service, but the provider requires a region.
Examples: "cn-hangzhou", "cn-shanghai", "us-west-1", "ap-southeast-1".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.zoneName

`string` · required

The private zone name (e.g., "internal.example.com", "db.corp").
This is the DNS suffix that records in the zone are created under.
Cannot be changed after creation (ForceNew).

- rule: {"required":true,"string":{"minLen":"1","maxLen":"253"}}

### spec.remark

`string`

Remark or description for the zone.
Visible in the Private Zone console; useful for noting the zone's purpose.

### spec.resourceGroupId

`string`

Alibaba Cloud resource group ID for access control and cost attribution.
If omitted, the zone is placed in the account's default resource group.
Cannot be changed after creation (ForceNew).

### spec.vpcAttachments

`[]AliCloudPrivateDnsZoneVpcAttachment` · required

VPCs to attach this private zone to. The zone will resolve DNS queries
only within attached VPCs. At least one VPC attachment is required.

Cross-region attachments are supported: set region_id on the attachment
to attach a VPC in a different region than spec.region.

- rule: {"repeated":{"minItems":"1"}}

### spec.vpcAttachments[].vpcId

`string | valueFrom` · required

VPC ID to attach. References an AliCloudVpc component's output.

- references: AliCloudVpc (`status.outputs.vpc_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AliCloudVpc, name: <that resource's name>, fieldPath: status.outputs.vpc_id}} -- a bare string does not parse

### spec.vpcAttachments[].regionId

`string`

Region of the VPC. Defaults to the spec.region when omitted.
Set this explicitly when attaching a VPC in a different region than
the provider's region (cross-region private zone attachment).

### spec.records

`[]AliCloudPrivateDnsZoneRecord`

DNS records within the private zone.
Each record maps a resource record name (rr) to a value (IP, hostname, etc.)
that is resolvable only within the attached VPCs.

### spec.records[].rr

`string` · required

Resource record name -- the hostname part under the zone.
Examples: "db" resolves as "db.internal.example.com" if zone is
"internal.example.com"; "@" for the zone apex.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.records[].type

`string` · required

DNS record type. Must be a type supported by Private Zone.

- rule: type must be one of: A, CNAME, MX, PTR, SRV, TXT
- rule: {"required":true}

### spec.records[].value

`string` · required

Record value. Interpretation depends on the record type:
  A     -> IPv4 address (e.g., "10.0.1.100")
  CNAME -> target hostname (e.g., "db-master.internal.example.com")
  MX    -> mail server hostname (e.g., "mail.internal.example.com")
  PTR   -> pointer hostname
  SRV   -> priority weight port target
  TXT   -> text content

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.records[].ttl

`int32`

Time-to-live in seconds. Controls how long VPC resolvers cache this record.
Default: 60.

### spec.records[].priority

`int32`

Priority for MX records. Ignored for other record types.
Range: 1 (highest priority) to 99 (lowest priority). Default: 1.

### spec.records[].remark

`string`

Remark or description for the record.

### spec.tags

`map<string, string>`

Tags to apply to the private zone.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AliCloudPrivateDnsZone, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.zone_id` | `string` | The Private Zone ID assigned by Alibaba Cloud. |
| `status.outputs.zone_name` | `string` | The zone name as created (e.g., "internal.example.com"). |
| `status.outputs.is_ptr` | `bool` | Whether the zone is a reverse-lookup (PTR) zone. Computed by the provider based on the zone name format. |
| `status.outputs.record_count` | `int32` | The number of DNS records in the zone. Computed by the provider after record creation. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.vpcAttachments[].vpcId` | AliCloudVpc | `status.outputs.vpc_id` |

## See Also

- [Overview](../README.md)
