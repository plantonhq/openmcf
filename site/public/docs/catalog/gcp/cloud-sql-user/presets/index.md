---
title: "Presets"
description: "Ready-to-deploy configuration presets for Cloud SQL User"
type: "preset-list"
componentSlug: "cloud-sql-user"
componentTitle: "Cloud SQL User"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-application-user"
    rank: "01"
    title: "Application User (Built-in)"
    excerpt: "This preset creates a classic username/password user for one application, with a lockout policy after repeated failed logins. One user per application — never share the instance's admin user across..."
  - slug: "02-iam-service-account-user"
    rank: "02"
    title: "IAM Service Account User (Passwordless)"
    excerpt: "This preset creates a passwordless database user for a workload's service account. Authentication flows through IAM — no credential to store, leak, or rotate — which is the strongest auth posture..."
---

# Cloud SQL User Presets

Ready-to-deploy configuration presets for Cloud SQL User. Each preset is a complete manifest you can copy, customize, and deploy.
