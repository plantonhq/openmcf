---
title: "Presets"
description: "Ready-to-deploy configuration presets for Service Account IAM Member on Google Cloud"
type: "preset-list"
componentSlug: "service-account-iam-member-on-google-cloud"
componentTitle: "Service Account IAM Member on Google Cloud"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-github-workload-identity-impersonation"
    rank: "01"
    title: "GitHub Workload Identity Impersonation"
    excerpt: "This preset grants `roles/iam.workloadIdentityUser` on one service account to the federated principal set of one GitHub repository — the terminal hop of keyless CI/CD. After this grant, workflows in..."
  - slug: "02-token-creator-grant"
    rank: "02"
    title: "Token Creator Grant (Cross-Account Impersonation)"
    excerpt: "This preset grants `roles/iam.serviceAccountTokenCreator` on one service account to another — letting the caller mint short-lived access and ID tokens AS the target. This is the building block of..."
  - slug: "03-deployer-act-as"
    rank: "03"
    title: "Deployer actAs Grant"
    excerpt: "This preset grants `roles/iam.serviceAccountUser` on a runtime service account to a deployer identity. This is the actAs permission that Cloud Run, GCE, Cloud Functions, and Dataflow all check at..."
---

# Service Account IAM Member on Google Cloud Presets

Ready-to-deploy configuration presets for Service Account IAM Member on Google Cloud. Each preset is a complete manifest you can copy, customize, and deploy.
