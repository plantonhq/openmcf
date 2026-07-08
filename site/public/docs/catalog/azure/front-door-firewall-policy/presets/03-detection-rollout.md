---
title: "Detection-Mode Rollout"
description: "This preset creates a Premium WAF policy in DETECTION mode -- the managed rule sets run and log everything they would have blocked, but no traffic is affected. It is the safe first step of every WAF..."
type: "preset"
rank: "03"
presetSlug: "03-detection-rollout"
componentSlug: "front-door-firewall-policy"
componentTitle: "Front Door Firewall Policy"
provider: "azure"
icon: "package"
order: 3
---

# Detection-Mode Rollout

This preset creates a Premium WAF policy in DETECTION mode -- the
managed rule sets run and log everything they would have blocked, but
no traffic is affected. It is the safe first step of every WAF rollout.

## When to Use

- Introducing a WAF in front of an existing production application:
  watch real traffic for a week before any request is blocked
- Trialing a managed rule set version upgrade (e.g. 2.0 to 2.1) beside
  the tuning you already trust

## Key Configuration Choices

- **`mode: DETECTION`** -- the single switch that separates observing
  from enforcing; flipping to PREVENTION later is an in-place update,
  no replacement
- **`action: RULE_SET_BLOCK` kept as the target posture** -- in
  detection mode the action only shapes the logs, so configuring the
  intended production action now means the flip to PREVENTION changes
  exactly one field

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-resource-name>` | The AzureResourceGroup's Planton resource name | Your Azure composition |

## The Rollout Loop

1. Deploy in DETECTION, associated with your domains through an
   AzureFrontDoorSecurityPolicy.
2. Review the WAF logs for matches on legitimate traffic.
3. Add `exclusions` (or per-rule `OVERRIDE_LOG` overrides) for each
   false positive.
4. Flip `mode` to PREVENTION.
