---
title: "Presets"
description: "Ready-to-deploy configuration presets for IAM Instance Profile"
type: "preset-list"
componentSlug: "iam-instance-profile"
componentTitle: "IAM Instance Profile"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-ec2-role-delivery"
    rank: "01"
    title: "EC2 Role Delivery"
    excerpt: "This preset wraps an `AwsIamRole` in an instance profile so EC2 instances can receive the role's temporary credentials through the instance metadata service. Reference this profile's..."
  - slug: "02-existing-role"
    rank: "02"
    title: "Wrap an Existing Role"
    excerpt: "This preset creates an instance profile carrying a role that already exists outside Planton -- useful when adopting EC2 workloads incrementally while IAM roles are still managed elsewhere."
---

# IAM Instance Profile Presets

Ready-to-deploy configuration presets for IAM Instance Profile. Each preset is a complete manifest you can copy, customize, and deploy.
