---
title: "Egress Baseline Rules"
description: "This preset creates the platform team's baseline rules group: DNS and NTP as network rules, and the package registries every build needs as application rules. It is the first group most policies..."
type: "preset"
rank: "01"
presetSlug: "01-egress-baseline-rules"
componentSlug: "firewall-policy-rule-collection-group"
componentTitle: "Firewall Policy Rule Collection Group"
provider: "azure"
icon: "package"
order: 1
---

# Egress Baseline Rules

This preset creates the platform team's baseline rules group: DNS and NTP
as network rules, and the package registries every build needs as
application rules. It is the first group most policies carry, and the one
application-scoped groups slot in around.

## When to Use

- The first rules on a fresh firewall policy
- Hub-spoke networks where spokes default-route through the firewall and
  need core plumbing (DNS, NTP, package pulls) explicitly allowed

## Key Configuration Choices

- **Group priority 500** -- leaves room for security-team overrides
  (lower numbers) and app-team groups (higher) without renumbering
- **Network rules for DNS/NTP, application rules for the web** -- L3/L4
  filtering where the destination is an address, L7 FQDN filtering where
  it is a name; the FQDN network rule (NTP) requires the policy's DNS
  proxy
- **`168.63.129.16`** is Azure's virtual public DNS address -- replace
  with your resolvers if you run your own

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<firewall-policy-name>` | The parent AzureFirewallPolicy | That policy's Planton resource name |
| `<spoke-cidr>` | The spoke address space these rules apply to | Your network plan |

## Downstream Wiring

Application teams add their own groups on the same policy without
touching this one:

```yaml
spec:
  firewallPolicyId:
    valueFrom:
      name: <firewall-policy-name>
  name: payments-app
  priority: 600
```
