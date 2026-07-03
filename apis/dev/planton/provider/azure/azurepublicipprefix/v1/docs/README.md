# AzurePublicIpPrefix -- Design Research

## The Resource

An Azure Public IP Prefix (`Microsoft.Network/publicIPPrefixes`) reserves a
contiguous range of public IP addresses from which individual public IPs can
be allocated and which NAT gateways can associate for outbound SNAT. The
component maps onto `azurerm_public_ip_prefix` (azurerm v4.x,
`internal/services/network/public_ip_prefix_resource.go`), parity-verified
against pulumi-azure v6 (`network.PublicIpPrefix`).

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `name` | Required, ForceNew |
| `location` | `region` | Required, ForceNew |
| `resource_group_name` | `resource_group` | FK → AzureResourceGroup |
| `prefix_length` | `prefix_length` | Azure default 28 when unset; ForceNew |
| `ip_version` | `ip_version` enum | IPv4 / IPv6; ForceNew |
| `sku` | `sku` enum | Standard / StandardV2; ForceNew |
| `sku_tier` | `sku_tier` enum | Regional / Global; ForceNew |
| `zones` | `zones` | ForceNew |
| `custom_ip_prefix_id` | `custom_ip_prefix_id` | Plain ARM ID (BYOIP); ForceNew |
| `tags` | `tags` | User tags merged over Planton-derived tags; **only updatable field** |
| `ip_prefix` (computed) | `ip_prefix` output | The reserved CIDR, assigned by Azure at creation |
| `id` (computed) | `public_ip_prefix_id` output | Join key for downstream FK references |

## Decomposition Decisions

- **Split OUT from NAT gateway and public IP.** A prefix has its own
  lifecycle (reserved once, consumed by many), is FK-referenced by both
  `AzurePublicIp` and `AzureNatGateway`, and exports a distinct allowlist
  surface (`ip_prefix`). Folding it inline would hide the range from the
  resource graph and force partners to rediscover addresses one at a time.
- **BYOIP stays a plain ARM ID.** Custom IP Prefix onboarding is a rare,
  telco-grade flow with no Planton kind; modeling it as a string keeps the
  common Microsoft-pool path simple without blocking BYOIP users.
- **No inline public IPs.** Addresses carved from the prefix are
  `AzurePublicIp` resources with `public_ip_prefix_id` set -- each keeps
  its own attach lifecycle (load balancer frontend today, gateway tomorrow).

## Design Decisions

- **Closed proto enums for SKU, tier, and IP version** with `defined_only`
  -- ARM's values are stable, and enums give wizards and agents the full
  option set from the spec alone.
- **`global_tier_requires_standard_sku` enforced in the spec** (message-level
  CEL): ARM rejects StandardV2 with the Global tier; catching it at
  validation time beats a mid-apply ARM error.
- **Unset enums and length defer to Azure defaults** (IPv4 / Standard /
  Regional / prefix length 28): both IaC engines send `null` for
  unspecified values so an empty spec and Azure's defaults deploy
  identically.
- **`ip_prefix` is an output, not an input** -- Azure assigns the actual
  range at creation; the only pre-create sizing knob is `prefix_length`.

## Operational Behavior Worth Knowing

- The prefix is **essentially immutable**: everything except tags is fixed
  at creation. Replacing it assigns a **new IP range** -- every partner
  allowlist and firewall rule tied to the old `ip_prefix` must migrate.
- A prefix **cannot be deleted** while public IPs still allocate from it or
  a NAT gateway still associates it.
- You are billed for **every address in the reserved range**, used or not --
  a /28 costs sixteen address-units; a /30 costs four.
- Each address in a NAT-associated prefix contributes **64,512 SNAT ports**;
  a /28 (16 addresses) scales outbound capacity sixteenfold in one
  allowlistable CIDR.
- Zone-redundant prefixes (`zones: ["1","2","3"]`) survive a single
  availability-zone failure; single-zone prefixes pin the range to one zone.

## Composition

- `resource_group` → `AzureResourceGroup.status.outputs.resource_group_name`
- `public_ip_prefix_id` output is consumed by:
  - `AzurePublicIp.public_ip_prefix_id` -- allocate an individual address
    from the range
  - `AzureNatGateway.public_ip_prefix_ids` -- associate the whole range for
    outbound SNAT
- `ip_prefix` output is the partner/firewall allowlist value surfaced to
  operators after deployment
