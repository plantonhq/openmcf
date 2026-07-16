---
title: "Presets"
description: "Ready-to-deploy configuration presets for Workload Identity Pool Provider"
type: "preset-list"
componentSlug: "workload-identity-pool-provider"
componentTitle: "Workload Identity Pool Provider"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-github-actions-oidc"
    rank: "01"
    title: "GitHub Actions OIDC Provider"
    excerpt: "This preset attaches GitHub Actions as a trusted issuer: workflows exchange their GitHub-minted OIDC tokens for short-lived Google credentials, replacing service-account keys in CI secrets entirely...."
  - slug: "02-aws-account"
    rank: "02"
    title: "AWS Account Provider"
    excerpt: "This preset trusts workloads running in one AWS account: EC2 instances, Lambda functions, and ECS tasks federate into Google Cloud using their native AWS credentials — no tokens to mint, no keys to..."
  - slug: "03-gitlab-ci-oidc"
    rank: "03"
    title: "GitLab CI OIDC Provider"
    excerpt: "This preset attaches GitLab CI as a trusted issuer: pipelines exchange their GitLab-minted ID tokens for short-lived Google credentials. Works for gitlab.com and (with the issuer URI changed)..."
---

# Workload Identity Pool Provider Presets

Ready-to-deploy configuration presets for Workload Identity Pool Provider. Each preset is a complete manifest you can copy, customize, and deploy.
