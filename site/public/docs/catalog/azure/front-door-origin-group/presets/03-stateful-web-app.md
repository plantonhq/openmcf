---
title: "Stateful Web App Backends"
description: "This preset creates a group tuned for session-based web applications: sticky sessions on, gentle HEAD probes against the root path, and a long traffic-restore ramp so recovering origins warm up..."
type: "preset"
rank: "03"
presetSlug: "03-stateful-web-app"
componentSlug: "front-door-origin-group"
componentTitle: "Front Door Origin Group"
provider: "azure"
icon: "package"
order: 3
---

# Stateful Web App Backends

This preset creates a group tuned for session-based web applications:
sticky sessions on, gentle HEAD probes against the root path, and a
long traffic-restore ramp so recovering origins warm up before taking
their full share.

## When to Use

- Web applications holding sessions in origin memory (no shared session
  store)
- Backends with cold-start costs -- large in-process caches, JIT
  warmup, connection pools -- that suffer when full traffic arrives at
  once

## Key Configuration Choices

- **`sessionAffinityEnabled: true`** -- cookie-based stickiness;
  matches Azure's default but is spelled explicitly here because it is
  the load-bearing choice for this shape
- **20-minute restore ramp** -- twice Azure's default; a recovered
  origin receives a growing slice of traffic instead of an avalanche
- **HEAD probes at 60 s** -- light-touch health checking; the affinity
  cookie keeps clients on their origin between probe cycles

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<front-door-profile-resource-name>` | The AzureFrontDoorProfile's Planton resource name | Your Front Door composition |
| `originGroupName` (example value) | 2-90 chars -- rename to your convention | Your naming convention |
