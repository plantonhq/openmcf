# AzureNatGateway

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureNatGatewaySpec** defines the configuration for creating an Azure
NAT Gateway: the managed source-network-address-translation (SNAT)
service that gives every workload in its attached subnets stable,
scalable outbound internet connectivity -- the production answer to
Azure retiring implicit default outbound access.

The gateway is deliberately just the gateway. What it composes with is
referenced, never created here:
- The public addresses it SNATs through are first-class AzurePublicIp
  and AzurePublicIpPrefix resources (public_ip_ids / public_ip_prefix_ids)
  -- visible in the resource graph, allowlistable, and reusable.
- The subnets it serves declare the attachment themselves
  (AzureSubnet's nat_gateway_id), matching Azure's model: one gateway
  serves many subnets without listing them.

A gateway with no addresses deploys but cannot translate anything --
associate at least one IP or prefix for it to actually carry traffic.
Each address provides 64,512 SNAT ports; a /28 prefix (16 addresses)
scales that by 16 in one contiguous, allowlistable range.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureNatGateway
metadata:
  name: test-nat
spec:
  region: eastus
  resourceGroup:
    value: test-rg
  name: test-nat
  zones:
    - "1"
  publicIpIds:
    - value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/publicIPAddresses/test-nat-ip
  publicIpPrefixIds:
    - value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/publicIPPrefixes/test-egress-prefix
  tags:
    purpose: hack-test
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.skuName` | `enum` |  |  |  |
| `spec.idleTimeoutInMinutes` | `int32` |  | `4` |  |
| `spec.zones` | `[]string` |  |  |  |
| `spec.publicIpIds` | `[]string \| valueFrom` |  |  | AzurePublicIp (`status.outputs.public_ip_id`) |
| `spec.publicIpPrefixIds` | `[]string \| valueFrom` |  |  | AzurePublicIpPrefix (`status.outputs.public_ip_prefix_id`) |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.region

`string` · required

The Azure region where the NAT gateway will be created, e.g. "eastus",
"westeurope". A gateway only serves subnets in its own region.
Changing the region replaces the gateway.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.resourceGroup

`string | valueFrom` · required

The Azure resource group the NAT gateway will be created in. Can be a
literal resource-group name or a reference to an AzureResourceGroup's
name output.

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.name

`string` · required

The name of the NAT gateway, unique within the resource group. 1-80
characters (alphanumerics, underscores, periods, and hyphens; must
start with a letter or number and end with a letter, number, or
underscore). Changing the name replaces the gateway, briefly
interrupting egress for every attached subnet.

- rule: NAT gateway names start with a letter or number, end with a letter, number, or underscore, and may contain alphanumerics, underscores, periods, and hyphens
- rule: {"required":true,"string":{"minLen":"1","maxLen":"80"}}

### spec.skuName

`enum`

The SKU. Unspecified applies Azure's default (STANDARD) -- a zonal
gateway, optionally pinned to one availability zone. STANDARD_V2 is
Azure's next-generation SKU: zone-redundant by default (zones must be
left empty) and requiring StandardV2 public IPs/prefixes. Fixed at
creation.

Allowed values (use exactly as shown):

- `azure_nat_gateway_sku_name_unspecified` -- Not specified: Azure's default (Standard).
- `STANDARD` -- The zonal production SKU: lives in one availability zone (or non-zonal), pairs with Standard public IPs.
- `STANDARD_V2` -- Azure's next-generation SKU: zone-redundant automatically (zones must be empty), pairs with StandardV2 public IPs/prefixes.

### spec.idleTimeoutInMinutes

`int32` · optional (explicit presence)

How long the gateway keeps an idle outbound TCP connection's SNAT port
reserved, in minutes (4-120). Unset uses Azure's default (4). Raise it
only for long-lived idle connections that must not be re-established;
higher values hold ports longer and hasten SNAT exhaustion. Updatable
in place.

- default: `4`
- rule: {"int32":{"lte":120,"gte":4}}

### spec.zones

`[]string`

The availability zone to pin a STANDARD gateway to (e.g. ["1"]). A
STANDARD gateway is zonal: it lives in one zone (or non-zonal when
empty), and zone-resilient architectures deploy one gateway per zone
with per-zone subnets. Must be empty for STANDARD_V2, which is
zone-redundant automatically. Fixed at creation.

- rule: {"repeated":{"items":{"string":{"in":["1","2","3"]}}}}

### spec.publicIpIds

`[]string | valueFrom`

The public IP addresses the gateway SNATs through, by ARM ID. Each
address adds 64,512 SNAT ports. References first-class AzurePublicIp
resources so the egress addresses are visible in the resource graph
and reusable (a StandardV2 gateway needs StandardV2 addresses).
Updatable in place.

- references: AzurePublicIp (`status.outputs.public_ip_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzurePublicIp, name: <that resource's name>, fieldPath: status.outputs.public_ip_id}} -- a bare string does not parse

### spec.publicIpPrefixIds

`[]string | valueFrom`

The public IP prefixes (contiguous reserved ranges) the gateway SNATs
through, by ARM ID -- the scalable, allowlistable alternative to
individual addresses. Updatable in place.

- references: AzurePublicIpPrefix (`status.outputs.public_ip_prefix_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzurePublicIpPrefix, name: <that resource's name>, fieldPath: status.outputs.public_ip_prefix_id}} -- a bare string does not parse

### spec.tags

`map<string, string>`

Free-form tags applied to the NAT gateway, merged over the
Planton-derived resource tags (organization, environment, resource
id); a user tag with the same key wins. Tags are Azure's governance
surface -- Azure Policy enforces them and Microsoft Cost Management
groups by them. Updatable in place.

## Validation Rules

- `standard_v2_is_zone_redundant`: zones must be empty for a STANDARD_V2 NAT gateway -- it is zone-redundant automatically

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureNatGateway, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.nat_gateway_id` | `string` | The Azure Resource Manager ID of the NAT gateway. This is the primary output: AzureSubnet's nat_gateway_id references it to attach the gateway to a subnet. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Network/natGateways/{name} |
| `status.outputs.nat_gateway_name` | `string` | The name of the NAT gateway. |
| `status.outputs.resource_guid` | `string` | The immutable GUID ARM assigns the gateway -- useful when correlating with Azure billing, monitoring, or support data that keys on the GUID rather than the ARM ID. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |
| `spec.publicIpIds` | AzurePublicIp | `status.outputs.public_ip_id` |
| `spec.publicIpPrefixIds` | AzurePublicIpPrefix | `status.outputs.public_ip_prefix_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureSubnet | `spec.natGatewayId` | `status.outputs.nat_gateway_id` |

## See Also

- [Overview](../README.md)
