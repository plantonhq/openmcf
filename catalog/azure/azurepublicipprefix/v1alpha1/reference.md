# AzurePublicIpPrefix

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzurePublicIpPrefixSpec** defines the configuration for creating an
Azure Public IP Prefix: a reserved, CONTIGUOUS range of public IP
addresses. Individual public IPs allocated from a prefix are guaranteed
to come from the same known range -- which is what lets partners and
firewalls allowlist a single CIDR instead of chasing individual
addresses, and what gives a NAT gateway a predictable, scalable SNAT
range.

The prefix is a first-class composable resource: AzurePublicIp allocates
individual addresses from it (public_ip_prefix_id), and AzureNatGateway
associates whole prefixes for outbound SNAT (public_ip_prefix_ids).

The prefix is essentially immutable: every field except tags is fixed at
creation, and a prefix cannot be deleted while any of its addresses are
in use.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzurePublicIpPrefix
metadata:
  name: test-public-ip-prefix
spec:
  region: eastus
  resourceGroup:
    value: test-rg
  name: test-egress-prefix
  prefixLength: 30
  zones:
    - "1"
    - "2"
    - "3"
  tags:
    purpose: hack-test
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.prefixLength` | `int32` |  |  |  |
| `spec.ipVersion` | `enum` |  |  |  |
| `spec.sku` | `enum` |  |  |  |
| `spec.skuTier` | `enum` |  |  |  |
| `spec.zones` | `[]string` |  |  |  |
| `spec.customIpPrefixId` | `string` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.region

`string` · required

The Azure region where the prefix will be created, e.g. "eastus",
"westeurope". Addresses can only be allocated from the prefix by
resources in the same region. Changing the region replaces the prefix
(and changes the actual IP range).

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.resourceGroup

`string | valueFrom` · required

The Azure resource group the prefix will be created in. Can be a
literal resource-group name or a reference to an AzureResourceGroup's
name output.

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.name

`string` · required

The name of the prefix, unique within the resource group. 1-80
characters (alphanumerics, underscores, periods, and hyphens; must
start with a letter or number and end with a letter, number, or
underscore). Changing the name replaces the prefix -- and with it the
actual IP range -- so name it after the egress boundary it represents
("prod-egress", "partner-allowlist").

- rule: Public IP prefix names start with a letter or number, end with a letter, number, or underscore, and may contain alphanumerics, underscores, periods, and hyphens
- rule: {"required":true,"string":{"minLen":"1","maxLen":"80"}}

### spec.prefixLength

`int32` · optional (explicit presence)

The CIDR length of the range to reserve. Azure's default is 28 (16
addresses). IPv4 prefixes span /21 (2,048 addresses) to /31 (2); IPv6
prefixes are much longer (e.g. /124-/127). Smaller numbers reserve
bigger ranges and bill for every reserved address whether used or not.
Fixed at creation.

- rule: {"int32":{"lte":127,"gte":0}}

### spec.ipVersion

`enum`

The IP version of the range. Unspecified applies Azure's default
(IPv4). Fixed at creation.

Allowed values (use exactly as shown):

- `azure_public_ip_prefix_ip_version_unspecified` -- Not specified: Azure's default (IPv4).
- `IPV4` -- An IPv4 range.
- `IPV6` -- An IPv6 range.

### spec.sku

`enum`

The SKU of the prefix. Unspecified applies Azure's default (STANDARD)
-- the production tier every current architecture uses. STANDARD_V2 is
Azure's next-generation SKU, required to associate the prefix with a
StandardV2 NAT gateway. Fixed at creation; a GLOBAL-tier prefix must
keep the STANDARD SKU.

Allowed values (use exactly as shown):

- `azure_public_ip_prefix_sku_unspecified` -- Not specified: Azure's default (Standard).
- `STANDARD` -- The production SKU every current architecture uses.
- `STANDARD_V2` -- Azure's next-generation SKU, required for StandardV2 NAT gateway association. Not valid with the GLOBAL tier.

### spec.skuTier

`enum`

The SKU tier. Unspecified applies Azure's default (REGIONAL) -- correct
for virtually everything, including NAT gateway ranges. GLOBAL exists
solely for cross-region load balancer frontends and requires the
STANDARD SKU. Fixed at creation.

Allowed values (use exactly as shown):

- `azure_public_ip_prefix_sku_tier_unspecified` -- Not specified: Azure's default (Regional).
- `REGIONAL` -- A regional range -- correct for virtually everything, including NAT gateway SNAT ranges.
- `GLOBAL` -- A globally-anycast range for cross-region load balancer frontends. Requires the STANDARD SKU.

### spec.zones

`[]string`

Availability zones the prefix's addresses are anchored to. Multiple
zones ("1","2","3") make the range zone-redundant -- the production
default; a single zone pins it; empty leaves the choice to Azure.
Zone support varies by region. Fixed at creation.

- rule: {"repeated":{"items":{"string":{"in":["1","2","3"]}}}}

### spec.customIpPrefixId

`string`

The ARM ID of a Custom IP Prefix (bring-your-own IP range onboarded to
Azure) to carve this prefix out of. Plain ARM ID: BYOIP onboarding is a
rare, telco-grade flow that is not modeled as a Planton kind. Omit to
allocate from Microsoft's address pool -- the overwhelmingly common
case. Fixed at creation.
Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Network/customIpPrefixes/{name}

### spec.tags

`map<string, string>`

Free-form tags applied to the prefix, merged over the Planton-derived
resource tags (organization, environment, resource id); a user tag with
the same key wins. The ONLY thing on a prefix that updates in place.

## Validation Rules

- `global_tier_requires_standard_sku`: A GLOBAL-tier prefix must keep the STANDARD SKU (ARM rejects StandardV2 with the Global tier)

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzurePublicIpPrefix, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.public_ip_prefix_id` | `string` | The Azure Resource Manager ID of the prefix. Referenced by AzurePublicIp (public_ip_prefix_id) to allocate addresses from the range, and by AzureNatGateway (public_ip_prefix_ids) to associate the whole range for outbound SNAT. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Network/publicIPPrefixes/{name} |
| `status.outputs.ip_prefix` | `string` | The actual reserved CIDR range, e.g. "20.42.0.16/28" -- the value partners and firewalls allowlist. Known only after creation (Azure assigns the range). |
| `status.outputs.public_ip_prefix_name` | `string` | The name of the prefix resource. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureAksCluster | `spec.defaultNodePool.nodePublicIpPrefixId` | `status.outputs.public_ip_prefix_id` |
| AzureAksCluster | `spec.networkProfile.loadBalancerProfile.outboundIpPrefixIds` | `status.outputs.public_ip_prefix_id` |
| AzureAksNodePool | `spec.nodePublicIpPrefixId` | `status.outputs.public_ip_prefix_id` |
| AzureLoadBalancer | `spec.frontendIpConfigurations[].publicIpPrefixId` | `status.outputs.public_ip_prefix_id` |
| AzureNatGateway | `spec.publicIpPrefixIds` | `status.outputs.public_ip_prefix_id` |
| AzurePublicIp | `spec.publicIpPrefixId` | `status.outputs.public_ip_prefix_id` |
| AzureVirtualMachineScaleSet | `spec.networkInterfaces[].ipConfigurations[].publicIpAddress.publicIpPrefixId` | `status.outputs.public_ip_prefix_id` |

## See Also

- [Overview](../README.md)
