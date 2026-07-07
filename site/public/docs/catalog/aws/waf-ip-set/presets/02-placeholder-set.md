---
title: "Placeholder IP Set"
description: "An empty REGIONAL IPv4 set — valid in AWS and useful when the web ACL and its rules must exist before the address list is finalized."
type: "preset"
rank: "02"
presetSlug: "02-placeholder-set"
componentSlug: "waf-ip-set"
componentTitle: "WAF IP Set"
provider: "aws"
icon: "package"
order: 2
---

# Placeholder IP Set

An empty REGIONAL IPv4 set — valid in AWS and useful when the web ACL and its rules must exist before the address list is finalized.

## When to Use

- Infrastructure-as-code pipelines that wire the web ACL rule tree first and back-fill addresses later
- Partner integrations where contractually approved ranges arrive after the security baseline ships
- Environments where ops owns the set and app teams own the web ACL

## What It Configures

- **Empty `addresses`** — the set matches nothing until entries are added; referencing allow rules pass no traffic, block rules never trigger on IP match
- **Description as runbook** — documents who owns populating the list

## What to Customize

- Replace `<aws-region>` with your target region
- Add CIDR entries when ranges are known — updates are in-place, no web ACL redeploy required

## Operational Note

An allow rule referencing an empty set is effectively a no-op for IP matching. Pair with a restrictive default action only after the set is populated, or use count mode on the rule while validating.
