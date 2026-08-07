---
title: "Presets"
description: "Ready-to-deploy configuration presets for Project IAM Member on Google Cloud"
type: "preset-list"
componentSlug: "project-iam-member-on-google-cloud"
componentTitle: "Project IAM Member on Google Cloud"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-service-account-grant"
    rank: "01"
    title: "Service Account Grant (Predefined Role)"
    excerpt: "This preset grants one predefined role to a service account — the most common IAM grant in infrastructure code. The member references a GcpServiceAccount resource, so the access relationship is a..."
  - slug: "02-custom-role-grant"
    rank: "02"
    title: "Custom Role Grant (Fully Composed)"
    excerpt: "This preset is the full least-privilege composition: a GcpIamCustomRole defines exactly the permissions a workload needs, a GcpServiceAccount is the workload's identity, and this grant is the edge..."
  - slug: "03-conditional-grant"
    rank: "03"
    title: "Conditional Grant (Time-Bound Access)"
    excerpt: "This preset grants a role that only applies while a CEL condition evaluates true — here, an expiry timestamp that makes the access self-revoking. Conditions also support resource-name prefixes and..."
---

# Project IAM Member on Google Cloud Presets

Ready-to-deploy configuration presets for Project IAM Member on Google Cloud. Each preset is a complete manifest you can copy, customize, and deploy.
