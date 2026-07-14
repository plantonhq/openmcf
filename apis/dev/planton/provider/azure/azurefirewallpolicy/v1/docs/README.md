# AzureFirewallPolicy -- Design Research

## The Resource

An Azure Firewall Policy (`Microsoft.Network/firewallPolicies`) is the
reusable rule-and-inspection document Azure Firewall instances enforce.
The component maps onto `azurerm_firewall_policy` (azurerm v4.x,
`internal/services/firewall/firewall_policy_resource.go`),
parity-verified against pulumi-azure v6 (`network.FirewallPolicy` --
zero bridge lag on this surface).

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `name` | Required, ForceNew; provider regex (start alnum, end alnum/underscore, 2+ chars) mirrored as a CEL |
| `location` | `region` | Required, ForceNew |
| `resource_group_name` | `resource_group` | FK → AzureResourceGroup |
| `sku` | `sku` | Optional, default Standard, ForceNew; closed enum |
| `base_policy_id` | `base_policy_id` | Self-referencing FK → AzureFirewallPolicy |
| `threat_intelligence_mode` | `threat_intelligence_mode` | Default Alert; closed enum |
| `threat_intelligence_allowlist` | `threat_intelligence_allowlist` | AtLeastOneOf(ip_addresses, fqdns) mirrored as a message CEL |
| `dns` | `dns` | servers + proxy_enabled (default false) |
| `intrusion_detection` | `intrusion_detection` | Premium-only (see gating) |
| `identity` | `identity` | SystemAssigned / UserAssigned / both |
| `tls_certificate` | `tls_certificate` | Premium-only; FK → AzureKeyVaultCertificate.versionless_secret_id |
| `insights` | `insights` | enabled presence-tracked (explicit false legal) |
| `explicit_proxy` | `explicit_proxy` | port bound 0-35536 (the provider's published validator; mirrored so deploys cannot fail on it) |
| `sql_redirect_allowed` | `sql_redirect_allowed` | plain bool |
| `private_ip_ranges` | `private_ip_ranges` | MinItems 1 → sent only when non-empty |
| `auto_learn_private_ranges_enabled` | `auto_learn_private_ranges_enabled` | one-way flag (see below) |
| `tags` | `tags` | user tags merged over Planton-derived |
| `child_policies`, `firewalls`, `rule_collection_groups` (computed) | not modeled | reverse indexes -- see recorded skips |

## Premium Gating (front-loaded contract)

The provider carries NO CustomizeDiff and no SKU-aware validation --
`intrusion_detection` and `tls_certificate` on a non-Premium policy are
rejected only by ARM at apply time. Azure's own contract (both features
are Azure Firewall Premium features; the provider's own acceptance
fixtures configure them exclusively on `sku = "Premium"`) is front-loaded
as two root-message CELs so the error is actionable at authoring time.

## Behavioral Notes (from the provider source)

- **`auto_learn_private_ranges_enabled` is one-way on the wire**: the
  provider only ever sends `AutoLearnPrivateRanges = "Enabled"` (a bool
  `false` is treated as absent by its `GetOk`); disabling is by omission.
  Both engines mirror this: the flag is sent only when true.
- **The sku and threat-intel mode are always sent explicitly**
  (Standard/Alert when unspecified) -- both have provider/Azure defaults,
  made deterministic so the engines produce identical payloads.
- **`explicit_proxy` ports are capped at 35536** by the provider's
  published validator (not 65535). The spec mirrors the enforced bound --
  a value the deploy path rejects is not worth accepting earlier -- and
  documents the quirk.
- **Identity `TypeNone` guard**: the provider leaves identity nil rather
  than sending "None"; the spec models the block as optional so absence
  is the none-state.
- The create/update path is a name-scoped mutex + CreateOrUpdate; no
  apply-time validation beyond the schema exists in the provider.

## Recorded Skips

- **`child_policies` / `firewalls` / `rule_collection_groups` (computed
  reverse indexes)**: not modeled as outputs. Composition flows the other
  way (children reference the policy); exporting them would churn the
  policy's outputs every time an unrelated consumer attaches.
- **ARM's IDPS profile surface** (`FirewallPolicyIntrusionDetectionProfileType`
  Basic/Standard/Advanced/Extended and the IDPS signature-query types):
  present in the ARM SDK but with NO azurerm schema field -- the
  neither-engine-can-ship-it class. Lands when azurerm models it.
- **`insights.firewall_location` naming**: kept as azurerm's argument
  name (the ARM wire name is `region`); the spec documents the meaning.

## Composition Seams

- `AzureFirewallPolicyRuleCollectionGroup.firewall_policy_id` → this
  kind's `firewall_policy_id` (rules nest under the policy).
- `AzureFirewall.firewall_policy_id` → this kind (enforcement).
- `base_policy_id` → this kind (inheritance; self-referencing).
- `tls_certificate.key_vault_secret_id` → `AzureKeyVaultCertificate`'s
  `versionless_secret_id` (rotation-follows-latest; pin with the
  versioned `secret_id` explicitly).
- `identity.user_assigned_identity_ids` → `AzureUserAssignedIdentity`.
- `insights.*workspace_id` → `AzureLogAnalyticsWorkspace.workspace_id`.
- IDPS `traffic_bypass` source/destination IP Groups → `AzureIpGroup`.

## Lifecycle

- ForceNew: `name`, `region`, `resource_group`, `sku`. Everything else
  updates in place.
- Tier pairing is ARM-enforced: policy tier == firewall tier; base-policy
  tier == child tier.
- Timeouts: create/update/delete 30 min (a fast control-plane document).

## Parity Notes

- pulumi-azure v6.38 covers every azurerm v4.80 field on this resource
  (verified field-by-field); zero PARITY-EXCEPTIONs beyond the
  family-wide tag-shape note in both modules.
- Both engines send identical payloads for the same stack input: sku and
  threat-intel mode explicit, auto-learn only-when-true, absent blocks
  omitted.
