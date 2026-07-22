---
title: "Presets"
description: "Ready-to-deploy configuration presets for KubernetesSecretStore"
type: "preset-list"
componentSlug: "kubernetessecretstore"
componentTitle: "KubernetesSecretStore"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-team-gcp-secret-manager"
    rank: "01"
    title: "Team-Scoped GCP Secret Manager (Keyless)"
    excerpt: "This preset creates a namespaced store connected to GCP Secret Manager, authenticating through a ServiceAccount in the team's own namespace whose GKE Workload Identity binding authorizes the reads...."
  - slug: "02-vault-approle"
    rank: "02"
    title: "Team-Scoped Vault KV with AppRole"
    excerpt: "This preset creates a namespaced store connected to a HashiCorp Vault (or OpenBao — same API) KV v2 engine, authenticating with AppRole machine identity: a role-id plus a secret-id. The secret-id is..."
  - slug: "03-fake-sandbox"
    rank: "03"
    title: "Fake Backend Sandbox (Test-Only)"
    excerpt: "This preset creates a namespaced store backed by ESO's built-in fake backend: the store serves the literal key/value entries declared in the spec — no external account, no network, fully..."
---

# KubernetesSecretStore Presets

Ready-to-deploy configuration presets for KubernetesSecretStore. Each preset is a complete manifest you can copy, customize, and deploy.
