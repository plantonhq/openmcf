---
title: "Launch-Template Instance"
description: "This preset launches an instance from an `AwsLaunchTemplate` -- the org's golden baseline of AMI, hardening (IMDSv2, encrypted volumes), and storage defaults -- and overrides only what makes this..."
type: "preset"
rank: "02"
presetSlug: "02-launch-template"
componentSlug: "ec2-instance"
componentTitle: "EC2 Instance"
provider: "aws"
icon: "package"
order: 2
---

# Launch-Template Instance

This preset launches an instance from an `AwsLaunchTemplate` -- the org's golden baseline of AMI, hardening (IMDSv2, encrypted volumes), and storage defaults -- and overrides only what makes this instance different. The template centralizes the opinionated parts; the instance spec carries just the deviations.

## When to Use

- Organizations that maintain a golden launch template for all EC2 compute
- Standalone instances that should inherit the same baseline as the auto-scaling fleets
- Keeping per-instance manifests minimal: AMI rotations and hardening changes ship through the template, not through every instance spec

## Key Configuration Choices

- **Template by reference** (`launchTemplate.id`) -- Reference an `AwsLaunchTemplate`'s `launch_template_id` output; the template supplies the AMI, metadata options, and block devices
- **`$Default` version** -- Both Planton launch-template modules promote each new template version to default, so this instance picks up template releases on its next stop/start or replacement
- **Inline overrides win** -- `instanceType` and `subnetId` set here override the template for this one instance; everything unset inherits

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<aws-region>` | AWS region (must match the template's region) | AWS region list |
| `<launch-template-id>` | The launch template ID (e.g., `lt-0123456789abcdef0`) | `AwsLaunchTemplate` status outputs |
| `<private-subnet-id>` | Private subnet ID where the instance will launch | `AwsSubnet` status outputs |

## Related Presets

- **01-ssm-managed** -- Fully inline configuration when no template exists
- **03-spot-worker** -- Interruption-tolerant Spot capacity for batch and CI workloads
