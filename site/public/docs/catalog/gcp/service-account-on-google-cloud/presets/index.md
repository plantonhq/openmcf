---
title: "Presets"
description: "Ready-to-deploy configuration presets for Service Account on Google Cloud"
type: "preset-list"
componentSlug: "service-account-on-google-cloud"
componentTitle: "Service Account on Google Cloud"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-workload-identity"
    rank: "01"
    title: "Workload Identity Service Account"
    excerpt: "This preset creates a GCP service account designed for GKE Workload Identity. No JSON key is generated -- pods authenticate via KSA-to-GSA binding instead. The account is granted logging, monitoring,..."
  - slug: "02-ci-cd-pipeline"
    rank: "02"
    title: "CI/CD Pipeline Service Account"
    excerpt: "This preset creates a GCP service account with a JSON key for CI/CD pipelines (GitHub Actions, GitLab CI, Jenkins). It has permissions to push container images, deploy to GKE, and deploy Cloud Run..."
  - slug: "03-identity-with-first-class-grants"
    rank: "03"
    title: "Pure Identity with First-Class Grants"
    excerpt: "This preset creates the service account as a pure identity node with no inline role lists. Every grant lives as its own GcpProjectIamMember resource referencing this account's `member` output — the..."
---

# Service Account on Google Cloud Presets

Ready-to-deploy configuration presets for Service Account on Google Cloud. Each preset is a complete manifest you can copy, customize, and deploy.
