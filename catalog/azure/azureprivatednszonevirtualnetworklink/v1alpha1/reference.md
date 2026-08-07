# AzurePrivateDnsZoneVirtualNetworkLink

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzurePrivateDnsZoneVirtualNetworkLinkSpec** defines the configuration
for linking an Azure Private DNS zone to a virtual network -- the
attachment that makes the zone's records resolvable from inside that
network. A zone without links answers nobody; each link adds one network
to its audience.

The link is a first-class resource because it is many-per-zone with its
own lifecycle: a hub-and-spoke topology links one zone (say,
"privatelink.postgres.database.azure.com") to the hub and every spoke
network, and networks join and leave the topology without touching the
zone or each other. Compose it with AzurePrivateDnsZone (the zone) and
AzureVirtualNetwork (the network); one link resource per zone-network
pair.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzurePrivateDnsZoneVirtualNetworkLink
metadata:
  name: test-link
  org: test-org
  env: dev
spec:
  name: hub-vnet
  privateDnsZoneId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/privateDnsZones/corp.internal
  virtualNetworkId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/hub
  registrationEnabled: true
  tags:
    cost-center: platform
    owner: network-team
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.name` | `string` | yes |  |  |
| `spec.privateDnsZoneId` | `string \| valueFrom` | yes |  | AzurePrivateDnsZone (`status.outputs.zone_id`) |
| `spec.virtualNetworkId` | `string \| valueFrom` | yes |  | AzureVirtualNetwork (`status.outputs.virtual_network_id`) |
| `spec.registrationEnabled` | `bool` |  | `false` |  |
| `spec.resolutionPolicy` | `enum` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.name

`string` · required

The link's name -- its ARM resource name under the parent zone, unique
per zone; 1-80 characters (alphanumerics, underscores, periods, and
hyphens). Name it after the network it attaches ("hub-vnet",
"spoke-payments"). Changing the name replaces the link (a brief
resolution gap for the affected network, nothing else).

- rule: Link names may contain alphanumerics, underscores, periods, and hyphens, and must start and end with an alphanumeric or underscore
- rule: {"required":true,"string":{"minLen":"1","maxLen":"80"}}

### spec.privateDnsZoneId

`string | valueFrom` · required

The private DNS zone this link is written on. Takes the zone's full
ARM resource ID; defaults to referencing an AzurePrivateDnsZone's
zone_id output in composed environments. The link is a child resource
of the zone -- the zone's name and resource group are derived from
this ID, so they can never contradict it. Changing the zone replaces
the link.

- references: AzurePrivateDnsZone (`status.outputs.zone_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzurePrivateDnsZone, name: <that resource's name>, fieldPath: status.outputs.zone_id}} -- a bare string does not parse

### spec.virtualNetworkId

`string | valueFrom` · required

The virtual network the zone becomes resolvable from. Takes the
network's full ARM resource ID; defaults to referencing an
AzureVirtualNetwork's virtual_network_id output in composed
environments. Changing the network replaces the link.

- references: AzureVirtualNetwork (`status.outputs.virtual_network_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureVirtualNetwork, name: <that resource's name>, fieldPath: status.outputs.virtual_network_id}} -- a bare string does not parse

### spec.registrationEnabled

`bool` · optional (explicit presence)

Whether Azure auto-registers DNS records for virtual machines in the
linked network: each VM gets an A record in the zone at boot, removed
when it goes away. Useful for custom internal zones ("corp.internal")
where machines should be discoverable by hostname. Keep it false (the
default) for Private Link zones ("privatelink.*") -- their records are
managed by private endpoints, not by VM registration -- and note Azure
allows only ONE registration-enabled link per network.

- default: `false`

### spec.resolutionPolicy

`enum`

How resolution failures in the zone behave for this network.
Unspecified lets Azure choose (it applies its own default per zone
type). DEFAULT answers strictly from the private zone.
NX_DOMAIN_REDIRECT retries names the private zone cannot answer
against public DNS -- the fallback pattern for Private Link zones
shared across environments where some records exist only publicly.

Allowed values (use exactly as shown):

- `azure_private_dns_zone_virtual_network_link_resolution_policy_unspecified` -- Not specified: Azure applies its default policy for the zone type.
- `DEFAULT` -- Answer strictly from the private zone; unresolvable names fail.
- `NX_DOMAIN_REDIRECT` -- Retry names the private zone answers NXDOMAIN for against public DNS -- the fallback pattern for shared Private Link zones.

### spec.tags

`map<string, string>`

Free-form tags applied to the link, merged over the Planton-derived
resource tags (organization, environment, resource id); a user tag
with the same key wins. Tags are Azure's governance surface -- Azure
Policy enforces them and Microsoft Cost Management groups by them --
so carry your org's ownership/cost-center conventions here. Updatable
in place.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzurePrivateDnsZoneVirtualNetworkLink, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.link_id` | `string` | The Azure Resource Manager ID of the virtual network link. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Network/privateDnsZones/{zone}/virtualNetworkLinks/{name} |
| `status.outputs.link_name` | `string` | The name of the virtual network link. Echoed from the spec for convenience. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.privateDnsZoneId` | AzurePrivateDnsZone | `status.outputs.zone_id` |
| `spec.virtualNetworkId` | AzureVirtualNetwork | `status.outputs.virtual_network_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
