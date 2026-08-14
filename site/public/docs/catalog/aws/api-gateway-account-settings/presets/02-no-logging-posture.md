---
title: "No-Logging Posture"
description: "This preset explicitly manages the region with NO API Gateway logging role — applying it clears any role previously set by anyone."
type: "preset"
rank: "02"
presetSlug: "02-no-logging-posture"
componentSlug: "api-gateway-account-settings"
componentTitle: "API Gateway Account Settings"
provider: "aws"
icon: "package"
order: 2
---

# No-Logging Posture

This preset explicitly manages the region with NO API Gateway logging
role — applying it clears any role previously set by anyone.

## When to Use

- Governance: pin the region to "no API Gateway logging" as declared
  state, instead of leaving the setting unmanaged
- Cleaning up after ad-hoc console experiments set a role

## What You Get

- The account-level role cleared and kept clear on every apply

## Customize

- Swap in preset 01 when the region should start logging
