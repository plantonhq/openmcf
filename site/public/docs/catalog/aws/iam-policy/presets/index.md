---
title: "Presets"
description: "Ready-to-deploy configuration presets for IAM Policy"
type: "preset-list"
componentSlug: "iam-policy"
componentTitle: "IAM Policy"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-s3-read-only"
    rank: "01"
    title: "S3 Read-Only Access"
    excerpt: "This preset creates a managed policy granting read-only access to a single S3 bucket -- the most common shared permission set in a typical AWS estate. Attach it to any role or user through their..."
  - slug: "02-permissions-boundary"
    rank: "02"
    title: "Workload Permissions Boundary"
    excerpt: "This preset creates a permissions-boundary policy: the ceiling on what any principal carrying it can ever do, regardless of what its permission policies grant. Apply it through a role's or user's..."
---

# IAM Policy Presets

Ready-to-deploy configuration presets for IAM Policy. Each preset is a complete manifest you can copy, customize, and deploy.
