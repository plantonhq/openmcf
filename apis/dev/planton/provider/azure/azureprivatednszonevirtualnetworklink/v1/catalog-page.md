# Azure Private DNS Zone Virtual Network Link

Creates a virtual network link on an Azure Private DNS zone -- the attachment that makes the zone's records resolvable from a virtual network, optionally with automatic VM record registration and a resolution-fallback policy. One link per zone-network pair; hub-and-spoke topologies link one zone to many networks.

## What Gets Created

When you deploy an AzurePrivateDnsZoneVirtualNetworkLink resource, Planton provisions:

- **Virtual Network Link** — an `azurerm_private_dns_zone_virtual_network_link` on the referenced zone, binding it to the referenced network

The link is an ARM child of the zone: the zone's name and resource group are derived from the referenced zone's ARM ID, so they can never contradict it.

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **A private DNS zone** to link (an `AzurePrivateDnsZone` in composed environments)
- **A virtual network** to make the zone resolvable from (an `AzureVirtualNetwork` in composed environments)
- **Write rights**: `Microsoft.Network/privateDnsZones/virtualNetworkLinks/write` on the zone and `Microsoft.Network/virtualNetworks/join/action` on the network (Private DNS Zone Contributor + Network Contributor, or Contributor/Owner)

## Quick Start

Create a file `zone-link.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzurePrivateDnsZoneVirtualNetworkLink
metadata:
  name: postgres-zone-hub-link
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzurePrivateDnsZoneVirtualNetworkLink.postgres-zone-hub-link
spec:
  name: hub-vnet
  privateDnsZoneId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/network-rg/providers/Microsoft.Network/privateDnsZones/privatelink.postgres.database.azure.com
  virtualNetworkId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/network-rg/providers/Microsoft.Network/virtualNetworks/hub-vnet
```

Deploy:

```shell
planton apply -f zone-link.yaml
```

Workloads in `hub-vnet` now resolve the zone's records to their private IPs.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `name` | `string` | The link's ARM resource name under the parent zone; name it after the network it attaches. | Required, 1-80 chars, unique per zone |
| `privateDnsZoneId` | `StringValueOrRef` | ARM ID of the parent zone. Defaults to referencing an `AzurePrivateDnsZone`'s `zone_id` output. | Required |
| `virtualNetworkId` | `StringValueOrRef` | ARM ID of the network the zone becomes resolvable from. Defaults to referencing an `AzureVirtualNetwork`'s `virtual_network_id` output. | Required |

### Optional Fields

| Field | Type | Description |
|-------|------|-------------|
| `registrationEnabled` | `bool` | Auto-register A records for VMs in the linked network (created at boot, removed at deallocation). Default `false`. Useful for custom internal zones; keep `false` for privatelink zones (their records are managed by private endpoints). Azure allows only ONE registration-enabled link per network. |
| `resolutionPolicy` | `enum` | `DEFAULT` (answer strictly from the private zone) or `NX_DOMAIN_REDIRECT` (retry unanswerable names against public DNS). Unset lets Azure apply its per-zone-type default. |
| `tags` | `map(string)` | User tags, merged over Planton-derived tags (user wins on collision). |

## Examples

### Composed: Zone + Network by Reference

```yaml
apiVersion: azure.planton.dev/v1
kind: AzurePrivateDnsZoneVirtualNetworkLink
metadata:
  name: postgres-zone-spoke-link
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzurePrivateDnsZoneVirtualNetworkLink.postgres-zone-spoke-link
spec:
  name: spoke-payments
  privateDnsZoneId:
    valueFrom:
      name: postgres-privatelink-zone
  virtualNetworkId:
    valueFrom:
      name: payments-network
```

### Internal Zone with VM Auto-Registration

```yaml
apiVersion: azure.planton.dev/v1
kind: AzurePrivateDnsZoneVirtualNetworkLink
metadata:
  name: corp-zone-hub-link
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzurePrivateDnsZoneVirtualNetworkLink.corp-zone-hub-link
spec:
  name: hub-vnet
  privateDnsZoneId:
    valueFrom:
      name: corp-internal-zone
  virtualNetworkId:
    valueFrom:
      name: hub-network
  registrationEnabled: true
```

Every VM in the hub network is now discoverable by hostname in `corp.internal`.

### Hub-and-Spoke: One Zone, Many Networks

Deploy several link resources referencing the same zone -- one per network. Networks join and leave the resolution audience without touching the zone or each other:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzurePrivateDnsZoneVirtualNetworkLink
metadata:
  name: postgres-zone-spoke2-link
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzurePrivateDnsZoneVirtualNetworkLink.postgres-zone-spoke2-link
spec:
  name: spoke-analytics
  privateDnsZoneId:
    valueFrom:
      name: postgres-privatelink-zone
  virtualNetworkId:
    valueFrom:
      name: analytics-network
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `link_id` | `string` | Full ARM ID of the link (`{zone-id}/virtualNetworkLinks/{name}`) |
| `link_name` | `string` | The link's name as deployed |

## Related Components

- [AzurePrivateDnsZone](/docs/catalog/azure/private-dns-zone) — the zone this link makes resolvable
- [AzureVirtualNetwork](/docs/catalog/azure/virtual-network) — the network gaining resolution
- [AzurePrivateEndpoint](/docs/catalog/azure/private-endpoint) — registers PaaS records in privatelink zones
