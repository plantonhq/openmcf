# AzureIpGroup -- Design Research

## The Resource

An Azure IP Group (`Microsoft.Network/ipGroups`) is a named, reusable set
of IP addresses and CIDR ranges that Azure Firewall and Firewall Policy
rules reference by ARM id, replacing repeated literal address lists with a
single curated address object. The component maps onto `azurerm_ip_group`
(azurerm v4.x, `internal/services/network/ip_group_resource.go`),
parity-verified against pulumi-azure v6 (`network.IPGroup`).

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `name` | Required, ForceNew |
| `location` | `region` | Required, ForceNew |
| `resource_group_name` | `resource_group` | FK → AzureResourceGroup |
| `cidrs` | `cidrs` | Optional set of addresses/CIDRs; updates in place |
| `tags` | `tags` | User tags merged over Planton-derived tags |
| `id` (computed) | `ip_group_id` output | Join key for rule references |
| `firewall_ids` (computed) | not modeled | Reverse index -- see below |
| `firewall_policy_ids` (computed) | not modeled | Reverse index -- see below |

This is the full v4.80 input surface -- the IP Group is deliberately a
passive address set.

## Recorded Skips

- **`firewall_ids` / `firewall_policy_ids` (computed reverse indexes).**
  Azure maintains, on the group, the list of firewalls and policies that
  currently reference it. These are not inputs and no downstream resource
  consumes them as a composition seam (composition flows the other way:
  rules reference the group). Exporting them would also make the group's
  outputs churn every time an unrelated policy adds a reference. Skipped
  as outputs; lands if a real consumer appears.
- **`azurerm_ip_group_cidr`** (a single CIDR as its own resource) is a
  Terraform-ergonomics construct for teams splitting one group's entries
  across modules -- the entries have no independent Azure lifecycle apart
  from their group document. The `cidrs` list on this kind is the honest
  management surface; a per-entry kind would be decomposition for its own
  sake.

## Decomposition Decisions

- **The group holds addresses, nothing else.** Azure models consumption
  from the rule's side: a firewall policy rule lists
  `source_ip_groups`/`destination_ip_groups`, and intrusion-detection
  traffic bypasses do the same. Planton mirrors this exactly -- the IP
  Group kind is the passive anchor; every consumer references
  `ip_group_id`.
- **First-class, not folded.** The address set is curated on its own
  schedule and referenced by many rules across many policies -- the
  textbook split criterion. Folding addresses inline into every rule is
  exactly the copy-paste problem the ARM resource exists to solve.

## Composition Seams

- `AzureFirewallPolicyRuleCollectionGroup` rules'
  `source_ip_groups`/`destination_ip_groups` → this kind's `ip_group_id`.
- `AzureFirewallPolicy` intrusion-detection `traffic_bypass`
  source/destination IP Groups → this kind (Premium IDPS bypass lists).

## Lifecycle

- `name` and `region` are ForceNew: renaming or moving the group replaces
  it, orphaning every rule reference until re-pointed.
- `cidrs` and `tags` update in place; an address change retargets every
  referencing rule without touching the rules.
- The provider serializes writes against any firewall/policy that
  references the group (named locks) to avoid ARM conflicts -- no
  spec-visible consequence, but explains why updates briefly queue behind
  policy deployments.
- Azure Firewall limits: at most 100 policies may reference one group; at
  most 5,000 entries per group; both enforced by ARM at apply time.

## Parity Notes

- Both engines create a single `IPGroup`/`azurerm_ip_group` with identical
  name, location, resource group, cidrs, and merged tags.
- Tag-shape divergence (lowered enum string vs snake-case literal;
  `resource_id` omission when `metadata.id` is empty) is the family-wide
  `PARITY-EXCEPTION` documented in both modules -- output-neutral, since
  tags never feed stack outputs.
