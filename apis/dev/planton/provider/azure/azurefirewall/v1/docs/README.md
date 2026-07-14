# AzureFirewall -- Design Research

## The Resource

An Azure Firewall (`Microsoft.Network/azureFirewalls`) is the managed
stateful firewall data plane. The component maps onto `azurerm_firewall`
(azurerm v4.x, `internal/services/firewall/firewall_resource.go`),
parity-verified against pulumi-azure v6 (`network.Firewall` -- zero
bridge lag, including `management_ip_configuration` and
`dns_proxy_enabled`).

## Decomposition Decisions

- **Policy split from firewall**: Azure separates the rule/inspection
  document (`firewallPolicies`) from the data plane; one policy serves
  many firewalls. The firewall models only WHERE enforcement runs.
- **Classic inline rule collections: recorded skip.** azurerm ships three
  legacy resources (`azurerm_firewall_application_rule_collection`,
  `_nat_rule_collection`, `_network_rule_collection`) that mutate rule
  slices ON the firewall -- the pre-policy management model. ARM rejects
  classic collections on a policy-attached firewall, and a greenfield
  catalog modeling both would carry two competing rule surfaces.
  Policy-based management is the modeled path; the classic trio lands
  only if adoption ever demands managing a legacy fleet.
- **Addressing is referenced, never bundled**: subnets and public IPs are
  first-class `AzureSubnet`/`AzurePublicIp` references -- the firewall
  consumes their ARM ids.
- **`virtual_hub_id` stays a bare reference**: no Virtual WAN family
  exists in the catalog yet (recorded to the adoption backlog), so the
  hub id is supplied as a literal or an explicit kind/fieldPath
  reference.

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `name` | Required, ForceNew; FirewallName regex mirrored as CEL |
| `location` | `region` | Required, ForceNew |
| `resource_group_name` | `resource_group` | FK → AzureResourceGroup |
| `sku_name` | `sku_name` | Required in azurerm; enum with unspecified→AZFW_VNet (the catalog's explicit-default pattern) |
| `sku_tier` | `sku_tier` | Required in azurerm; unspecified→Standard |
| `ip_configuration` | `ip_configurations` | subnet FK → AzureSubnet, PIP FK → AzurePublicIp |
| `management_ip_configuration` | `management_ip_configuration` | ForceNew block; subnet + PIP both required |
| `firewall_policy_id` | `firewall_policy_id` | FK → AzureFirewallPolicy |
| `threat_intel_mode` | `threat_intel_mode` | Optional+Computed in azurerm → sent only when specified |
| `dns_servers` / `dns_proxy_enabled` | same | see DNS coupling below |
| `private_ip_ranges` | `private_ip_ranges` | CIDR or literal `IANAPrivateRanges` |
| `virtual_hub` | `virtual_hub` | hub id + public_ip_count (default 1) |
| `zones` | `zones` | ForceNew |
| `tags` | `tags` | merged over Planton-derived |

## Front-Loaded Contracts (all seven as root-message CELs)

The provider has NO CustomizeDiff on this resource; its imperative
apply-time checks plus ARM's structural contracts are front-loaded:

1. **AZFW_HUB ⇔ virtual_hub** (the deployment-model pairing; structural
   ARM reality, visible in the provider's own hub fixtures).
2. **virtual_hub excludes ip_configurations.**
3. **VNet firewall requires ≥1 ip_configuration.**
4. **Exactly ONE subnet-bearing ip_configuration** (the provider's
   `validateFirewallIPConfigurationSettings`, verbatim).
5. **Public IP required without a management path** (the provider's
   documented create contract: "A public ip address is required unless a
   management_ip_configuration block is specified").
6. **Management name ≠ any ip_configuration name** (the provider's
   apply-time collision check).
7. **Policy owns DNS** -- a policy-attached firewall must not set
   `dns_servers`/`dns_proxy_enabled`; ARM rejects firewall-level DNS
   parameters with `AzureFirewallDNSConfigNotAllowedForVhubOrVnetWithPolicy`
   (live-verified). DNS moves to the policy's `dns` block.

Documented but NOT CEL'd (not verifiable from the provider source):
BASIC tier requiring the management configuration (Azure's documented
requirement, ARM-enforced; the spec comment teaches it); tier pairing
with the policy (cross-resource, ARM-enforced).

## Behavioral Notes (from the provider source)

- **DNS is wire-encoded as AdditionalProperties**: `dns_servers` sets
  `Network.DNS.EnableProxy=true` AND the server list -- servers
  implicitly force the proxy on regardless of `dns_proxy_enabled`. The
  spec documents the coupling; the modules pass both through verbatim so
  the coupling stays Azure's.
- **`threat_intel_mode` is Optional+Computed**: sent only when specified;
  Azure defaults it (Alert). Policy-attached firewalls take posture from
  the policy.
- **Subnet name validators**: the provider validates the subnet ID's
  name segment is exactly `AzureFirewallSubnet` (management:
  `AzureFirewallManagementSubnet`) at plan time. Because the spec's
  subnet is a reference (unresolvable at authoring time), the contract
  lives on the FK's documentation and the E2E fixture, and the provider
  still enforces it at plan time after resolution.
- **Virtual-hub scale-down** retains the first N public IPs (provider
  expand logic) -- relevant only to AZFW_HUB updates.
- **Locking**: the provider serializes against the policy, the firewall,
  and every touched VNet/subnet on create/update/delete.
- **Timeouts**: create/update/delete 90 minutes -- the slowest resource
  in the Azure catalog's networking family.

## Composition Seams

- `ip_configurations[].subnet_id` → `AzureSubnet.subnet_id` (the
  AzureFirewallSubnet).
- `ip_configurations[].public_ip_address_id` /
  `management_ip_configuration.public_ip_address_id` →
  `AzurePublicIp.public_ip_id` (Standard/static).
- `firewall_policy_id` → `AzureFirewallPolicy.firewall_policy_id`.
- **`private_ip_address` output** → `AzureRouteTable` routes'
  `next_hop_in_ip_address` (VIRTUAL_APPLIANCE) -- the hub-spoke seam.
- DNAT rules' `destination_address` composes from the referenced
  `AzurePublicIp.ip_address` output.

## Parity Notes

- Both engines send identical payloads: sku pair explicit, threat-intel
  only-when-specified, dns flags verbatim, absent blocks omitted.
- The data-path `private_ip_address` output is read back from the
  subnet-bearing ip configuration on both engines (first non-empty).
- Zero PARITY-EXCEPTIONs beyond the family-wide tag-shape note.
