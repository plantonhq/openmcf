# AzureFrontDoorFirewallPolicy

A Web Application Firewall (WAF) policy for Azure Front Door: the edge
rule set that inspects HTTP traffic at Microsoft's edge network, before
requests ever reach an origin. Custom match and rate-limit rules run
first, then (on PREMIUM) Microsoft's managed rule sets --
Microsoft_DefaultRuleSet for OWASP-class attacks and the Bot Manager
for bot classification -- with scoped exclusions and per-rule tuning.

The policy is a GLOBAL, resource-group-scoped resource -- a different
ARM type than the regional WAF policy Application Gateways use, with
its own rule vocabulary. It also enforces nothing on its own: an
AzureFrontDoorSecurityPolicy associates it with a profile's domains,
and one policy is commonly shared by many profiles.

## When to Use

Use AzureFrontDoorFirewallPolicy when you need:

- **OWASP-class protection at the edge** -- SQL injection, XSS, RCE,
  protocol attacks blocked by Microsoft-maintained rules before
  traffic reaches your origins (PREMIUM)
- **Per-client rate limiting** -- throttle abusive clients globally,
  with the counting done at the edge
- **Bot management** -- classify good/bad/unknown bots and gate the
  unknown ones behind a JavaScript challenge or CAPTCHA (PREMIUM)
- **IP/geo allow- and denylists** enforced once for every domain the
  policy protects

## Key Configuration

- `policy_name` -- 1-128 letters/digits (no hyphens), unique within the
  resource group; ForceNew
- `sku` -- STANDARD (custom rules only) or PREMIUM (managed rules +
  challenges); ForceNew, must MATCH the profile's sku at association
  time, and Azure refuses a PREMIUM-to-STANDARD downgrade outright
- `mode` -- DETECTION (log only, the tuning mode) or PREVENTION (act,
  the production posture); flips in place
- `custom_rules[]` -- your own rules, by ascending priority: match
  conditions (9 request variables, 12 operators, transforms) plus an
  action (allow/block/log/redirect/JS challenge/CAPTCHA)
- `managed_rules[]` -- Microsoft's curated sets with set-wide
  `exclusions` and per-group/per-rule `overrides` (PREMIUM only)
- `log_scrubbing` -- redact sensitive request parts (auth headers, PII,
  client IPs) from the WAF logs

## Composition

```yaml
resourceGroup:
  valueFrom:
    kind: AzureResourceGroup
    name: my-resource-group
    fieldPath: status.outputs.resource_group_name
```

An AzureFrontDoorSecurityPolicy attaches the policy to a profile's
endpoints and custom domains through its `firewall_policy_id` output --
that association is what turns enforcement on.

## Documentation

- [Design research](docs/README.md) -- field mapping, recorded skips
- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
