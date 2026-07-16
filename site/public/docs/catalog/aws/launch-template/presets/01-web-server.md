---
title: "Web Server Fleet"
description: "This preset creates the launch blueprint for an auto-scaled web fleet: IMDSv2 enforced, an encrypted gp3 root volume, detailed monitoring for responsive scaling, and instance identity via an IAM..."
type: "preset"
rank: "01"
presetSlug: "01-web-server"
componentSlug: "launch-template"
componentTitle: "Launch Template"
provider: "aws"
icon: "package"
order: 1
---

# Web Server Fleet

This preset creates the launch blueprint for an auto-scaled web fleet:
IMDSv2 enforced, an encrypted gp3 root volume, detailed monitoring for
responsive scaling, and instance identity via an IAM instance profile
reference. Reference this template's `launch_template_id` output from an
`AwsAutoScalingGroup`.

## When to Use

- Web or API services behind an ALB target group, managed by an
  auto-scaling group
- Any stateless fleet whose instances should be interchangeable and
  hardened by default

## Key Configuration Choices

- **`httpTokens: required`** -- IMDSv2 only; blocks the SSRF-to-credential
  path. Hop limit 2 keeps containerized workloads able to reach metadata
- **Encrypted gp3 root volume** -- gp3's baseline (3000 IOPS / 125 MiB/s)
  outperforms gp2 at lower cost; account-default encryption is made
  explicit
- **`detailedMonitoring: true`** -- 1-minute metrics let target-tracking
  policies react several minutes sooner than the free 5-minute default
- **Instance profile by reference** -- the fleet's AWS identity (SSM
  access, ECR pulls) composes with `AwsIamInstanceProfile` instead of a
  hardcoded ARN

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<fleet-name>` | Name for the launch template (max 125 chars) | Your fleet's name (e.g., `web`, `api`) |
| `<aws-region>` | AWS region code (e.g., `us-east-1`) | Your deployment region |
| `ami-<your-ami-id>` | AMI to boot from (e.g., the current AL2023 image) | EC2 console or SSM public parameters |
| `<instance-profile-resource-name>` | Name of the AwsIamInstanceProfile resource | Your instance-profile manifest's `metadata.name` |
| `<security-group-resource-name>` | Name of the AwsSecurityGroup resource | Your security-group manifest's `metadata.name` |

## Common Additions

- Add `userData` with a bootstrap script (plain text; the modules
  base64-encode it)
- Set `instanceInitiatedShutdownBehavior: terminate` for strictly
  immutable fleet members
- Grow `volumeSizeGb` or add data-volume mappings for disk-heavy services

## Related Presets

- **02-spot-flexible** -- attribute-based Spot capacity instead of a fixed type
- **03-hardened** -- stricter protection flags and a customer-managed KMS key
