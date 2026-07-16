---
title: "Pre-Provisioned (Disabled) Endpoint"
description: "This preset creates the endpoint dark: fully provisioned, hostname generated, but not accepting traffic. Flipping `enabled` is a fast in-place update -- the launch switch."
type: "preset"
rank: "02"
presetSlug: "02-maintenance-disabled"
componentSlug: "front-door-endpoint"
componentTitle: "Front Door Endpoint"
provider: "azure"
icon: "package"
order: 2
---

# Pre-Provisioned (Disabled) Endpoint

This preset creates the endpoint dark: fully provisioned, hostname
generated, but not accepting traffic. Flipping `enabled` is a fast
in-place update -- the launch switch.

## When to Use

- Staged launches: provision the endpoint and its routes, prepare DNS
  against the generated hostname, then flip `enabled`
- Maintenance windows: disabling stops traffic at the edge without
  deleting any configuration

## Key Configuration Choices

- **`enabled: false` stops traffic at the edge** -- clients get errors
  from Front Door; the backend sees nothing
- **The hostname exists immediately** -- the `host_name` output is
  available while disabled, so CNAME records can be created ahead of
  the cutover

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<front-door-profile-resource-name>` | The AzureFrontDoorProfile's Planton resource name | Your Front Door composition |
| `endpointName` (example value) | 2-46 chars; becomes the public hostname prefix -- rename to your convention | Your naming convention |
