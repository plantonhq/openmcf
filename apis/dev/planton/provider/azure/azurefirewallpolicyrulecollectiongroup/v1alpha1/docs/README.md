# AzureFirewallPolicyRuleCollectionGroup -- Design Research

## The Resource

A rule collection group
(`Microsoft.Network/firewallPolicies/ruleCollectionGroups`) is the
ordered rules document nested under an Azure Firewall Policy. The
component maps onto `azurerm_firewall_policy_rule_collection_group`
(azurerm v4.x,
`internal/services/firewall/firewall_policy_rule_collection_group_resource.go`),
parity-verified against pulumi-azure v6
(`network.FirewallPolicyRuleCollectionGroup` -- zero bridge lag,
including rule-level `http_headers` and `description`).

## Decomposition Decisions

- **The group is its own kind** (not folded into the policy): a policy
  carries many groups, each with an independent lifecycle and its own
  nested ARM identity -- the textbook split. Azure's own portal deploys
  rules this way.
- **Rules FOLD inside the group**: a rule has no ARM identity, nothing
  references an individual rule, and ordering is the group's semantic --
  one ordered document (the Front Door rule-set shape).
- **The DNAT collection's action is NOT modeled**: azurerm's schema
  accepts exactly one value ("Dnat"); a one-value vocabulary is a
  constant, not a knob. Both engines send the provider's schema literal
  unconditionally. Re-enable trigger: ARM adding a second NAT action.

## Field Mapping Highlights (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `firewall_policy_id` | `firewall_policy_id` | Required, ForceNew; FK → AzureFirewallPolicy |
| `name` | `name` | Required, ForceNew; provider regex mirrored as CEL |
| `priority` | `priority` | IntBetween(100, 65000) mirrored |
| `application_rule_collection` | `application_rule_collections` | action Allow/Deny; full rule surface incl. http_headers, terminate_tls, web_categories |
| `network_rule_collection` | `network_rule_collections` | protocols Any/TCP/UDP/ICMP; FQDN destinations |
| `nat_rule_collection` | `nat_rule_collections` | protocols TCP/UDP only; translated XOR; destination_ports MaxItems 1 |

Rule-level source/destination `*_ip_groups` are repeated
`StringValueOrRef` FKs onto `AzureIpGroup.ip_group_id`.

## Front-Loaded Contracts (the provider's imperative validation as CELs)

The provider's ONLY imperative validation lives in its NAT-rule expander;
everything else is schema-level. Mirrored:

- **NAT translated target XOR** -- exactly one of
  `translated_address`/`translated_fqdn` (the provider errors on both or
  neither at apply time; a message CEL front-loads it).
- **NAT protocols TCP/UDP only** -- schema `StringInSlice`; a CEL over
  the shared protocol enum.
- **NAT destination_ports MaxItems 1**, entries 1-64000, no wildcard --
  mirrored as `max_items: 1` (the list shape kept because ARM's own model
  is a list that may lift the cap).
- **Priorities 100-65000** on the group and every collection.
- **Application protocol port 0-64000** (the provider's bound).

Contracts documented but NOT invented as CELs (ARM-side, not in the
provider): at-least-one source and at-least-one destination per rule;
Premium gating of destination_urls/web_categories/terminate_tls; FQDN
network rules requiring DNS proxy. The provider sends these shapes
verbatim and ARM rejects them with actionable errors.

## Recorded Skips

- **DNAT collection `action`** -- one-value constant (above).
- No computed attributes exist on the resource; the only outputs are the
  nested ARM id and name.

## Composition Seams

- `firewall_policy_id` → `AzureFirewallPolicy.firewall_policy_id`
  (parent).
- Rule `source_ip_groups`/`destination_ip_groups` →
  `AzureIpGroup.ip_group_id` (six reference sites across the three rule
  types).
- A DNAT rule's `destination_address` is one of the firewall's public
  IPs -- compose it from the referenced `AzurePublicIp.ip_address` output
  in charts/scenarios.

## Lifecycle

- ForceNew: `name`, `firewall_policy_id`. Priority and collections update
  in place.
- The provider takes a named lock on the parent policy for every
  create/update/delete -- concurrent groups on one policy serialize.
- Timeouts: create/update/delete 30 min.
- Read-path nuance in the provider: filter collections are re-classified
  application-vs-network by the concrete type of their FIRST rule; a
  collection with zero rules is skipped on read -- one more reason the
  spec requires min 1 rule per collection (the provider's schema does
  too).

## Parity Notes

- Both engines send identical payloads: same wire vocabularies
  (Allow/Deny, Dnat, Any/TCP/UDP/ICMP, Http/Https/Mssql), same
  only-when-set handling for descriptions and NAT optional fields.
- The Pulumi bridge models NAT `destination_ports` as a singular string;
  the module sends the list's single entry (the spec caps the list at
  one, so the shapes coincide exactly).
