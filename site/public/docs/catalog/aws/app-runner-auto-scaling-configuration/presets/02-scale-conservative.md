---
title: "Cost-Conscious"
description: "A scaling posture for internal tools and low-traffic services: one warm instance, maximum request packing, and a tight scale-out cap so a traffic anomaly can never triple the bill."
type: "preset"
rank: "02"
presetSlug: "02-scale-conservative"
componentSlug: "app-runner-auto-scaling-configuration"
componentTitle: "App Runner Auto Scaling Configuration"
provider: "aws"
icon: "package"
order: 2
---

# Cost-Conscious

A scaling posture for internal tools and low-traffic services: one warm instance, maximum request packing, and a tight scale-out cap so a traffic anomaly can never triple the bill.

## When to Use

- Internal dashboards, admin tools, and staging environments
- Services where a queueing delay under burst is acceptable

## What It Configures

- **`minSize: 1`** — the smallest warm floor App Runner allows
- **`maxConcurrency: 200`** — the maximum request packing per instance (the top of AWS's dial), so scale-out happens as late as possible
- **`maxSize: 3`** — a hard cost ceiling; excess requests queue rather than launching a fleet

## What to Customize

- Replace `<aws-region>` with your region
- Raise `maxSize` if the service ever sheds load during legitimate spikes
