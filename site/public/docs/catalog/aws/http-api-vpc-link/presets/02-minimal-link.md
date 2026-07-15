---
title: "Preset: Minimal Link"
description: "Use this preset for development or experimentation: a single-subnet, no-security-group link that gets a private integration working with the least ceremony. Reachability is then governed solely by..."
type: "preset"
rank: "02"
presetSlug: "02-minimal-link"
componentSlug: "http-api-vpc-link"
componentTitle: "HTTP API VPC Link"
provider: "aws"
icon: "package"
order: 2
---

# Preset: Minimal Link

## When to Use

Use this preset for development or experimentation: a single-subnet, no-security-group link that gets a private integration working with the least ceremony. Reachability is then governed solely by the target's security groups.

## Key Configuration Choices

- **Single subnet** -- the minimum AWS accepts. The link can only reach targets in this AZ, so it suits single-AZ dev environments.
- **No security groups** -- AWS applies no filtering on the link side; the backend's security group decides what the link can reach.

## What to Customize

1. **`<vpc-link-name>`** — Link name (e.g., `dev-link`)
2. **`<subnet-id>`** — A subnet ID in the target VPC (or switch to a `valueFrom` reference to an AwsSubnet resource)

## Graduating to Production

Move to the private-alb-link preset: two AZs and an explicit egress-scoped security group. Note that changing subnets or security groups **replaces** the link (AWS has no update API for them) -- referencing integrations follow the new link ID through the resource graph.
