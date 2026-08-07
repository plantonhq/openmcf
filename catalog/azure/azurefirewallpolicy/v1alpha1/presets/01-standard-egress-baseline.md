# Standard Egress Baseline

This preset creates the STANDARD-tier policy most hub-spoke deployments
start from: Microsoft threat intelligence in DENY (known-malicious
destinations are blocked, not just logged) and the DNS proxy enabled so
FQDN-based network rules resolve deterministically.

It is the "parent" half of the enterprise pattern: the security team owns
this baseline, application teams attach rule collection groups (or child
policies referencing it via `basePolicyId`) on their own schedules.

## When to Use

- The first firewall policy in a new hub-spoke network
- A security baseline that application-scoped child policies will extend
- Any deployment where FQDN rules matter (the DNS proxy is what makes
  them reliable)

## Key Configuration Choices

- **`threatIntelligenceMode: DENY`** -- the production posture; drop to
  ALERT while tuning if a partner endpoint trips false positives (and
  exempt it via `threatIntelligenceAllowlist`)
- **`dns.proxyEnabled: true`** -- point spoke DNS at the firewall's
  private IP so the firewall and clients resolve names identically
- **`sku: STANDARD`** -- upgrade to PREMIUM only for TLS inspection/IDPS;
  the tier is fixed at creation and must match the firewall's

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-name>` | The resource group to create the policy in | The resource group's `status.outputs.resource_group_name` |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |

## Downstream Wiring

Attach the policy to a firewall and nest rules under it:

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
