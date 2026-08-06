# AliCloudNatGateway

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `alicloud.planton.dev/v1alpha1`

AliCloudNatGatewaySpec defines the configuration for an Alibaba Cloud NAT
Gateway with bundled EIP association and SNAT entries.

A NAT Gateway enables resources in a VPC (typically placed in private
VSwitches with no public IP) to access the internet via SNAT (Source NAT).
This component bundles the NAT Gateway, its EIP association, and SNAT entries
into a single deployable unit (per DD07 composite bundling) because a NAT
Gateway without an EIP and at least one SNAT entry is non-functional.

The bundled flow:
  1. Create the Enhanced NAT Gateway in the specified VPC/VSwitch.
  2. Associate the provided EIP with the NAT Gateway.
  3. Create SNAT entries that map private VSwitch/CIDR traffic to the EIP.

The EIP's IP address is resolved internally by the IaC modules via a data
source lookup on the provided eip_id -- the user does not need to supply the
IP address separately.

Provider resources:
  Terraform: alicloud_nat_gateway + alicloud_eip_association + alicloud_snat_entry
  Pulumi:    vpc.NatGateway + ecs.EipAssociation + vpc.SnatEntry

## Example

```yaml
apiVersion: alicloud.planton.dev/v1alpha1
kind: AliCloudNatGateway
metadata:
  name: alicloudnatgateway-demo
spec:
  region: cn-hangzhou
  vpcId:
    value: vpc-demo123
  vswitchId:
    value: vsw-demo123
  natGatewayName: demo-nat
  description: Demo NAT Gateway for local testing
  eipId:
    value: eip-demo123
  snatEntries:
    - sourceVswitchId:
        value: vsw-app-demo
      snatEntryName: app-zone-a
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.vpcId` | `string \| valueFrom` | yes |  | AliCloudVpc (`status.outputs.vpc_id`) |
| `spec.vswitchId` | `string \| valueFrom` | yes |  | AliCloudVswitch (`status.outputs.vswitch_id`) |
| `spec.natGatewayName` | `string` | yes |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.natType` | `string` |  | `Enhanced` |  |
| `spec.paymentType` | `string` |  | `PayAsYouGo` |  |
| `spec.internetChargeType` | `string` |  | `PayByLcu` |  |
| `spec.specification` | `string` |  |  |  |
| `spec.deletionProtection` | `bool` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |
| `spec.eipId` | `string \| valueFrom` | yes |  | AliCloudEipAddress (`status.outputs.eip_id`) |
| `spec.snatEntries` | `[]AliCloudSnatEntry` |  |  |  |
| `spec.snatEntries[].sourceVswitchId` | `string \| valueFrom` |  |  | AliCloudVswitch (`status.outputs.vswitch_id`) |
| `spec.snatEntries[].sourceCidr` | `string` |  |  |  |
| `spec.snatEntries[].snatEntryName` | `string` |  |  |  |

## Field Details

### spec.region

`string` · required

Alibaba Cloud region where the NAT Gateway will be created.
Must match the region of the VPC and VSwitch.
Examples: "cn-hangzhou", "cn-shanghai", "us-west-1", "ap-southeast-1".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.vpcId

`string | valueFrom` · required

VPC ID that the NAT Gateway belongs to.

- references: AliCloudVpc (`status.outputs.vpc_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AliCloudVpc, name: <that resource's name>, fieldPath: status.outputs.vpc_id}} -- a bare string does not parse

### spec.vswitchId

`string | valueFrom` · required

VSwitch ID where the Enhanced NAT Gateway is placed.
The NAT Gateway consumes an IP from this VSwitch. Must be in the same VPC
as vpc_id. Required for Enhanced NAT Gateways (the default and recommended
type).

- references: AliCloudVswitch (`status.outputs.vswitch_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AliCloudVswitch, name: <that resource's name>, fieldPath: status.outputs.vswitch_id}} -- a bare string does not parse

### spec.natGatewayName

`string` · required

NAT Gateway name. 2-128 characters; must start with a letter or Chinese
character; can contain digits, underscores, periods, and hyphens.

- rule: {"required":true,"string":{"minLen":"2","maxLen":"128"}}

### spec.description

`string`

Human-readable description of the NAT Gateway's purpose.

### spec.natType

`string` · optional (explicit presence)

NAT Gateway type. "Enhanced" is the modern standard with higher
performance and VSwitch placement support. "Normal" is legacy.
Default: "Enhanced"

- default: `Enhanced`
- rule: nat_type must be one of: Enhanced, Normal

### spec.paymentType

`string` · optional (explicit presence)

Billing method for the NAT Gateway.
"PayAsYouGo" for on-demand billing, "Subscription" for reserved pricing.
Default: "PayAsYouGo"

- default: `PayAsYouGo`
- rule: payment_type must be one of: PayAsYouGo, Subscription

### spec.internetChargeType

`string` · optional (explicit presence)

Metering method for the NAT Gateway.
"PayByLcu" bills based on actual capacity units (recommended for most
workloads). "PayBySpec" bills based on a fixed specification tier.
Default: "PayByLcu"

- default: `PayByLcu`
- rule: internet_charge_type must be one of: PayByLcu, PayBySpec

### spec.specification

`string` · optional (explicit presence)

Fixed specification tier for the NAT Gateway. Only applicable when
internet_charge_type is "PayBySpec". Ignored when "PayByLcu".
Valid values: "Small" (1Gbps), "Middle" (5Gbps), "Large" (10Gbps),
"XLarge.1" (20Gbps).

- rule: specification must be one of: Small, Middle, Large, XLarge.1

### spec.deletionProtection

`bool` · optional (explicit presence)

Enable deletion protection to prevent accidental deletion via the console
or API. Must be explicitly disabled before the NAT Gateway can be deleted.

### spec.tags

`map<string, string>`

Tags to apply to the NAT Gateway resource.

### spec.eipId

`string | valueFrom` · required

EIP allocation ID to associate with this NAT Gateway.
The IaC module resolves the EIP's public IP address internally via a data
source lookup and uses it for all SNAT entries. The EIP must already exist
(create one with AliCloudEipAddress).

- references: AliCloudEipAddress (`status.outputs.eip_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AliCloudEipAddress, name: <that resource's name>, fieldPath: status.outputs.eip_id}} -- a bare string does not parse

### spec.snatEntries

`[]AliCloudSnatEntry`

SNAT entries that map private network sources to the NAT Gateway's EIP.
Each entry enables outbound internet access for traffic originating from
the specified VSwitch or CIDR block.

### spec.snatEntries[].sourceVswitchId

`string | valueFrom`

VSwitch whose traffic should be NATed through the EIP.
Use this for the common case of enabling internet access for an entire
VSwitch. Mutually exclusive with source_cidr.

- references: AliCloudVswitch (`status.outputs.vswitch_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AliCloudVswitch, name: <that resource's name>, fieldPath: status.outputs.vswitch_id}} -- a bare string does not parse

### spec.snatEntries[].sourceCidr

`string`

Source CIDR block whose traffic should be NATed. Use this for fine-grained
control when you want to NAT only a subset of a VSwitch's address space.
Mutually exclusive with source_vswitch_id.
Example: "10.0.1.0/24"

### spec.snatEntries[].snatEntryName

`string`

Human-readable name for this SNAT entry. 2-128 characters.
If omitted, a name is generated from the NAT Gateway name and entry index.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AliCloudNatGateway, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.nat_gateway_id` | `string` | The NAT Gateway ID assigned by Alibaba Cloud (e.g., "ngw-xxxxx"). |
| `status.outputs.nat_gateway_name` | `string` | The NAT Gateway name as created. |
| `status.outputs.snat_table_id` | `string` | The SNAT table ID, auto-created with the NAT Gateway. Exposed for advanced users who need to add SNAT entries outside Planton. |
| `status.outputs.forward_table_id` | `string` | The forward (DNAT) table ID, auto-created with the NAT Gateway. Exposed for users who want to add DNAT entries outside Planton. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.vpcId` | AliCloudVpc | `status.outputs.vpc_id` |
| `spec.vswitchId` | AliCloudVswitch | `status.outputs.vswitch_id` |
| `spec.eipId` | AliCloudEipAddress | `status.outputs.eip_id` |
| `spec.snatEntries[].sourceVswitchId` | AliCloudVswitch | `status.outputs.vswitch_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
