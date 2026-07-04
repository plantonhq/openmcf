---
title: "Presets"
description: "Ready-to-deploy configuration presets for Workload Identity Pool"
type: "preset-list"
componentSlug: "workload-identity-pool"
componentTitle: "Workload Identity Pool"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-ci-federation-pool"
    rank: "01"
    title: "CI Federation Pool"
    excerpt: "This preset creates the standard keyless-auth trust boundary for CI/CD: a FEDERATION_ONLY pool that GitHub Actions, GitLab CI, or any other OIDC issuer attaches to via a..."
  - slug: "02-locked-down-pool"
    rank: "02"
    title: "Locked-Down (Staged) Pool"
    excerpt: "This preset provisions the pool in a disabled state: the resource exists, its name is stable, IAM bindings and providers can be prepared against it — but every token exchange is rejected until..."
---

# Workload Identity Pool Presets

Ready-to-deploy configuration presets for Workload Identity Pool. Each preset is a complete manifest you can copy, customize, and deploy.
