---
title: "Reputation Defaults"
description: "This preset sets the recommended account-wide suppression posture: hard bounces and spam complaints automatically suppress recipient addresses across every send from the account."
type: "preset"
rank: "01"
presetSlug: "01-reputation-defaults"
componentSlug: "ses-account-settings"
componentTitle: "SES Account Settings"
provider: "aws"
icon: "package"
order: 1
---

# Reputation Defaults

This preset sets the recommended account-wide suppression posture:
hard bounces and spam complaints automatically suppress recipient
addresses across every send from the account.

## When to Use

- Every production SES region (this is the posture AWS recommends)
- Before scaling sending volume — reputation damage is easier to
  prevent than repair

## What You Get

- Bounced and complaining addresses skipped on all future sends
- Protection for your sending reputation account-wide

## Customize

- Add the `vdm` arm for deliverability analytics (preset 02)
- Suppression persists after destroy — apply an empty `reasons` list
  first if you ever mean to stop suppressing
