# AzureFirewallPolicy

## Overview

`AzureFirewallPolicy` provisions an Azure Firewall Policy: the reusable
rule-and-inspection document that Azure Firewall instances enforce. The
policy carries WHAT the firewall does -- threat intelligence posture, DNS
proxying, TLS inspection, intrusion detection (IDPS), SNAT behavior --
while the firewall instance (`AzureFirewall`) carries WHERE enforcement
runs. One policy attaches to many firewalls across regions, so security
policy is authored once and enforced everywhere.

Rules live in `AzureFirewallPolicyRuleCollectionGroup` children -- many
per policy, each deployed independently.

## Why a First-Class Resource?

- **One policy, many firewalls** -- the policy is authored on the security
  team's schedule and enforced by every attached firewall; folding it into
  the firewall would force rule changes through data-plane deployments
- **Inheritance** -- a policy can extend a base policy (`base_policy_id`),
  the enterprise pattern of a global baseline plus per-application child
  policies
- **Independent children** -- rule collection groups nest under the policy
  with their own lifecycles (one per team or application)

## Key Features

- **Tiers**: BASIC (SMB-scale), STANDARD (production default), PREMIUM
  (TLS inspection + IDPS + URL filtering + web categories). Premium-only
  blocks are validation-gated to the PREMIUM tier at authoring time
- **Threat intelligence** -- Microsoft's feed in ALERT (default) or DENY,
  with an allowlist for false positives
- **DNS proxy** -- required for FQDN network rules to resolve
  deterministically
- **TLS inspection (Premium)** -- the CA certificate referenced from an
  `AzureKeyVaultCertificate`'s versionless secret face (renewals followed
  automatically), read through a user-assigned identity
- **IDPS (Premium)** -- signature-based intrusion detection/prevention
  with per-signature overrides and trusted-flow bypasses (IP Groups
  supported)
- **Policy analytics** -- traffic analysis into Log Analytics with
  per-region workspace routing

## Spec Fields (top level)

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `region` | string | Yes | -- | Azure region (policy is regional; attachable anywhere) |
| `resource_group` | StringValueOrRef | Yes | -- | Resource group (defaults to AzureResourceGroup reference) |
| `name` | string | Yes | -- | Policy name, 2-80 chars; ForceNew |
| `sku` | enum | No | `STANDARD` | BASIC / STANDARD / PREMIUM; ForceNew |
| `base_policy_id` | StringValueOrRef | No | -- | Parent policy for inheritance (self-referencing FK) |
| `threat_intelligence_mode` | enum | No | `ALERT` | ALERT / DENY / OFF |
| `threat_intelligence_allowlist` | message | No | -- | Exempted IPs/FQDNs (at least one entry) |
| `dns` | message | No | -- | Custom servers + proxy toggle |
| `intrusion_detection` | message | No | -- | IDPS (PREMIUM only, validation-gated) |
| `identity` | message | No | -- | Managed identity (TLS inspection reads Key Vault through it) |
| `tls_certificate` | message | No | -- | TLS inspection CA (PREMIUM only, validation-gated) |
| `insights` | message | No | -- | Policy analytics into Log Analytics |
| `explicit_proxy` | message | No | -- | Forward-proxy mode (ports capped at 35536 by the provider) |
| `sql_redirect_allowed` | bool | No | false | Allow SQL redirect ports in FQDN rules |
| `private_ip_ranges` | list | No | IANA private | SNAT-exempt ranges |
| `auto_learn_private_ranges_enabled` | bool | No | false | Learn SNAT ranges from routes |
| `tags` | map | No | -- | User tags (merged over Planton-derived tags) |

## Outputs

| Output | Description |
|--------|-------------|
| `firewall_policy_id` | Full ARM ID -- referenced by rule collection groups, firewalls, and child policies |
| `firewall_policy_name` | The policy's name as deployed |
| `identity_principal_id` | System-assigned identity principal (empty without one) |

## Quick Example

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureFirewallPolicy
metadata:
  name: egress-baseline
  org: mycompany
  env: production
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: network-rg
  name: egress-baseline
  threatIntelligenceMode: DENY
  dns:
    proxyEnabled: true
```

Attach it from a firewall and nest rules under it:

```yaml
# AzureFirewall
spec:
  firewallPolicyId:
    valueFrom:
      name: egress-baseline

# AzureFirewallPolicyRuleCollectionGroup
spec:
  firewallPolicyId:
    valueFrom:
      name: egress-baseline
  priority: 500
```

## Lifecycle Notes

- `name`, `region`, `resource_group`, and `sku` are **fixed at creation**
- The **tier must match** every attached firewall's tier, and a base
  policy's tier must match its children (ARM validates at apply time)
- Premium features on an existing Standard policy require recreating the
  policy at PREMIUM (sku is ForceNew) -- choose the tier deliberately
- `auto_learn_private_ranges_enabled` only ever writes "Enabled" on the
  wire; disabling is by omission (Azure's own semantics)
- TLS inspection needs: PREMIUM sku + `identity` (user-assigned) +
  Key Vault secret read access ("Key Vault Secrets User") for that
  identity on the referenced certificate's secret
