# AliCloudCenInstance

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `alicloud.planton.dev/v1`

AliCloudCenInstanceSpec defines the configuration for an Alibaba Cloud
Cloud Enterprise Network (CEN) instance with bundled child-instance
attachments.

CEN is a global networking service that provides high-quality, low-latency
private connectivity between VPCs in different regions or between VPCs and
on-premises data centers (via VBR). Unlike most Alibaba Cloud resources,
CEN is region-agnostic -- the instance itself is not bound to a single
region. The region field below is used only for Alibaba Cloud API routing.

This component bundles the CEN instance and its child-instance attachments
into a single deployable unit (per DD07 composite bundling) because a CEN
instance without attachments connects nothing.

The bundled flow:
  1. Create the CEN instance.
  2. For each attachment entry, attach the specified VPC, VBR, or CCN
     to the CEN instance.

Provider resources:
  Terraform: alicloud_cen_instance + alicloud_cen_instance_attachment
  Pulumi:    cen.Instance + cen.InstanceAttachment

## Example

```yaml
apiVersion: alicloud.planton.dev/v1
kind: AliCloudCenInstance
metadata:
  name: test-cen
spec:
  region: cn-hangzhou
  cenInstanceName: test-cen
  description: Test CEN instance for development
  attachments:
    - childInstanceId:
        value: vpc-test-hangzhou
      childInstanceRegionId: cn-hangzhou
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.cenInstanceName` | `string` | yes |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.protectionLevel` | `string` |  |  |  |
| `spec.resourceGroupId` | `string` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |
| `spec.attachments` | `[]AliCloudCenAttachment` |  |  |  |
| `spec.attachments[].childInstanceId` | `string \| valueFrom` | yes |  | AliCloudVpc (`status.outputs.vpc_id`) |
| `spec.attachments[].childInstanceType` | `string` |  | `VPC` |  |
| `spec.attachments[].childInstanceRegionId` | `string` | yes |  |  |

## Field Details

### spec.region

`string` · required

Alibaba Cloud region used for API routing. CEN is a global service so
this does not restrict where attached networks reside -- each attachment
declares its own region via child_instance_region_id.
Examples: "cn-hangzhou", "cn-shanghai", "us-west-1", "ap-southeast-1".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.cenInstanceName

`string` · required

CEN instance name. 2-128 characters; must start with a letter or Chinese
character; can contain digits, underscores, periods, and hyphens.

- rule: {"required":true,"string":{"minLen":"2","maxLen":"128"}}

### spec.description

`string`

Human-readable description of the CEN instance's purpose.

### spec.protectionLevel

`string` · optional (explicit presence)

CIDR block overlap protection level for attached networks. When set to
"REDUCED", CEN allows overlapping CIDR blocks between attached VPCs
(routing is controlled by route maps). Leave empty for the default
strict mode that rejects CIDR overlaps.

- rule: protection_level must be empty or 'REDUCED'

### spec.resourceGroupId

`string`

Resource group ID for organizational access control. Leave empty to use
the account's default resource group.

### spec.tags

`map<string, string>`

Tags to apply to the CEN instance resource.

### spec.attachments

`[]AliCloudCenAttachment`

Child-instance attachments bundled with this CEN. Each entry attaches a
VPC, VBR (Virtual Border Router), or CCN (Cloud Connect Network) to the
CEN instance, enabling private connectivity between the attached networks.

### spec.attachments[].childInstanceId

`string | valueFrom` · required

ID of the child instance to attach. For VPCs this is the vpc_id; for VBRs
the vbr_id; for CCNs the ccn_id. When the child is a VPC managed by
Planton, use valueFrom to reference its output.

- references: AliCloudVpc (`status.outputs.vpc_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AliCloudVpc, name: <that resource's name>, fieldPath: status.outputs.vpc_id}} -- a bare string does not parse

### spec.attachments[].childInstanceType

`string` · optional (explicit presence)

Type of the child instance being attached.
Default: "VPC"

- default: `VPC`
- rule: child_instance_type must be one of: VPC, VBR, CCN

### spec.attachments[].childInstanceRegionId

`string` · required

Region of the child instance. Required because CEN is global and the
attached child instance can reside in any region.
Examples: "cn-hangzhou", "us-west-1", "eu-central-1".

- rule: {"required":true,"string":{"minLen":"1"}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AliCloudCenInstance, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.cen_id` | `string` | The CEN instance ID assigned by Alibaba Cloud (e.g., "cen-xxxxx"). |
| `status.outputs.cen_instance_name` | `string` | The CEN instance name as configured in the spec. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.attachments[].childInstanceId` | AliCloudVpc | `status.outputs.vpc_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
