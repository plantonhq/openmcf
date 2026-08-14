---
title: "VDM Analytics"
description: "This preset layers the Virtual Deliverability Manager on top of the reputation defaults: engagement dashboards plus Guardian delivery optimization."
type: "preset"
rank: "02"
presetSlug: "02-vdm-analytics"
componentSlug: "ses-account-settings"
componentTitle: "SES Account Settings"
provider: "aws"
icon: "package"
order: 2
---

# VDM Analytics

This preset layers the Virtual Deliverability Manager on top of the
reputation defaults: engagement dashboards plus Guardian delivery
optimization.

## When to Use

- Sending programs large enough that deliverability analytics pay for
  themselves (VDM carries its own AWS pricing)
- Teams actively tuning open/click engagement

## What You Get

- The VDM dashboard tracking engagement metrics
- Guardian adjusting delivery behavior to protect deliverability
- The suppression defaults from preset 01

## Customize

- Drop `optimizedSharedDelivery` to `false` to keep analytics without
  Guardian's behavior changes
- Destroy resets VDM to disabled (billing stops); suppression
  persists per its own contract
