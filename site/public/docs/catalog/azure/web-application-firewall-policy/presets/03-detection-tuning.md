---
title: "Detection-Mode Tuning with Exclusions and Overrides"
description: "This preset is the false-positive workflow: the policy runs in Detection mode (logging matches without blocking), a scoped exclusion skips the session cookie that trips SQL-injection rules, two..."
type: "preset"
rank: "03"
presetSlug: "03-detection-tuning"
componentSlug: "web-application-firewall-policy"
componentTitle: "Web Application Firewall Policy"
provider: "azure"
icon: "package"
order: 3
---

# Detection-Mode Tuning with Exclusions and Overrides

This preset is the false-positive workflow: the policy runs in Detection
mode (logging matches without blocking), a scoped exclusion skips the
session cookie that trips SQL-injection rules, two protocol-enforcement
rules are tuned, and log scrubbing redacts credentials from the very logs
the tuning reads.

## When to Use

- Rolling a WAF onto an existing application: run Detection for a week,
  review the matches, tune, then flip to Prevention
- Diagnosing production false positives without dropping protection
  elsewhere

## Key Configuration Choices

- **`DETECTION` mode** -- every managed and custom rule logs instead of
  blocking; flip `mode` to `PREVENTION` (or remove it -- Prevention is the
  default) once tuning converges
- **Exclusion over disablement** -- skipping ONE cookie for ONE rule group
  keeps SQL-injection protection intact for everything else; disabling
  rule 942440 outright would blind the whole policy to that pattern
- **Override semantics** -- a rule listed WITHOUT `enabled: true` is
  disabled (that is the tuning gesture); `enabled: true` plus an `action`
  retunes what the rule does
- **Log scrubbing** -- `Authorization` headers and all cookies are redacted
  from WAF logs before they land in Log Analytics, so tuning never leaks
  credentials into log storage

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-name>` | The resource group to create the policy in | The resource group's `status.outputs.resource_group_name` |
| `<policy-name>` | The policy's name, unique within the resource group | Your naming convention |
| `<session-cookie-name>` | The cookie that trips SQLI rules (e.g. a signed session token) | Your WAF logs' `details_data` field |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |
