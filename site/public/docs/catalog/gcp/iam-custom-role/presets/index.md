---
title: "Presets"
description: "Ready-to-deploy configuration presets for IAM Custom Role"
type: "preset-list"
componentSlug: "iam-custom-role"
componentTitle: "IAM Custom Role"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-workload-least-privilege"
    rank: "01"
    title: "Workload Least-Privilege Role"
    excerpt: "This preset defines a custom role with exactly the permissions one workload needs — the standard replacement for granting an over-broad predefined role like `roles/storage.admin` to a service account..."
  - slug: "02-readonly-auditor"
    rank: "02"
    title: "Read-Only Auditor Role"
    excerpt: "This preset defines a read-only role for security audits, compliance dashboards, and inventory tooling — visibility into project configuration, IAM policy, service accounts, and storage layout..."
  - slug: "03-ci-cd-deployer"
    rank: "03"
    title: "CI/CD Deployer Role"
    excerpt: "This preset defines a deployment role for CI/CD pipelines that roll out new Cloud Run revisions — update access to the services being deployed and pull access to the artifact repository, without..."
---

# IAM Custom Role Presets

Ready-to-deploy configuration presets for IAM Custom Role. Each preset is a complete manifest you can copy, customize, and deploy.
