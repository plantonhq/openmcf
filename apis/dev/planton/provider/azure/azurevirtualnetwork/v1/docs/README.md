# AzureVirtualNetwork -- Design Research

## The Resource

An Azure Virtual Network (`Microsoft.Network/virtualNetworks`) is the
regional, private, isolated address space Azure networking is built on.
Everything network-attached -- VMs, AKS nodes, private endpoints, VNet-
integrated PaaS -- exists inside one. The component maps 1:1 onto
`azurerm_virtual_network` (azurerm v4.x, `internal/services/network/
virtual_network_resource.go`), parity-verified against pulumi-azure v6
(`network.VirtualNetwork`).

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `name` | Required, ForceNew |
| `location` | `region` | Required, ForceNew |
| `resource_group_name` | `resource_group` | FK → AzureResourceGroup |
| `address_space` | `address_spaces` | ExactlyOneOf with `ip_address_pool` → message-level CEL XOR |
| `ip_address_pool` | `ip_address_pools` | Max 2 (one per IP version); `number_of_ip_addresses` is a string (IPv6 counts exceed int range) and can only grow |
| `dns_servers` | `dns_servers` | Empty = Azure default resolver (168.63.129.16) |
| `bgp_community` | `bgp_community` | `asn:community`; azurerm validates both segments 1-65534; the ASN is always 12076 (Microsoft's) today |
| `ddos_protection_plan` | `ddos_protection_plan` | `id` + `enable` -- ARM keeps attach and activate distinct |
| `encryption.enforcement` | `encryption` enum | Block-of-one-enum flattened to a proto enum; unspecified = block absent = encryption off |
| `flow_timeout_in_minutes` | `flow_timeout_in_minutes` | 4-30; optional int, unset = ARM default (4) |
| `private_endpoint_vnet_policies` | `private_endpoint_vnet_policies` enum | ARM default `Disabled` = unspecified; only `Basic` is sent |
| `edge_zone` | `edge_zone` | Optional, ForceNew |
| `tags` | `tags` | User tags merged over Planton-derived tags |
| `guid` (computed) | output `guid` | ARM's stable network GUID |

## Deliberately Not Modeled (recorded reasons)

- **Inline `subnet` blocks.** azurerm supports subnets both inline on the
  network and as standalone `azurerm_subnet` resources, and warns the two
  paths conflict. The standalone resource is the modern grain and matches
  the catalog's `AzureSubnet` kind -- one composable node per subnet, each
  with delegations/service endpoints/policies. Modeling inline subnets
  would create two competing owners for the same child resource.
- **`azurerm_virtual_network_dns_servers`** (the standalone DNS-servers
  resource): redundant with the inline `dns_servers` field this spec
  carries; azurerm itself warns against using both.

## Design Decisions

- **Encryption modeled as a flattened enum**, not a one-field nested
  message: presence-of-block == enabled maps cleanly to unspecified ==
  disabled, and matches the catalog's enum convention for opt-in modes
  (unspecified is always the provider default). The provider's own note
  that ARM currently accepts only `AllowUnencrypted` is recorded in the
  field comment; `DROP_UNENCRYPTED` is modeled because the API defines it.
- **`ddos_protection_plan.id` is a plain string**, not a StringValueOrRef:
  no DDoS-plan kind exists in the catalog to reference (the plan is a
  shared, org-level, separately billed resource typically managed
  centrally). If a first-class plan kind is ever added, this field is the
  FK seam.
- **`address_spaces` output echoes the ACTUAL ranges** rather than the
  spec, because with IPAM pools the ranges are provisioned at deploy time
  and the output is the only place downstream planning can read them.
- **Name in spec (not derived from metadata)**: the ARM name is the
  network's identity and its subnets join on it; an explicit field keeps
  the contract visible and renameable-with-intent.

## Operational Behavior Worth Knowing

- Address space edits are in-place, but ARM rejects removing/shrinking a
  block that subnets are carved from.
- DNS server changes propagate on DHCP lease renewal -- running VMs need a
  restart to pick them up immediately.
- `flow_timeout_in_minutes` governs connection tracking for intra-network
  flows; long-lived idle connections (database sessions, message-bus
  consumers) are the reason to raise it.
- Deleting a virtual network requires it to be empty; the composed model
  (subnets as separate resources) means destroy ordering is handled by the
  dependency graph.

## Composition

- `resource_group` → `AzureResourceGroup.status.outputs.resource_group_name`
- Consumed by: `AzureSubnet.vnet_id`,
  `AzurePrivateDnsZoneVirtualNetworkLink.virtual_network_id` (both →
  `status.outputs.virtual_network_id`)
- Future networking-wave kinds (VNet peering) will reference the same
  output.
