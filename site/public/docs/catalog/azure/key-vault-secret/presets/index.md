---
title: "Presets"
description: "Ready-to-deploy configuration presets for Key Vault Secret"
type: "preset-list"
componentSlug: "key-vault-secret"
componentTitle: "Key Vault Secret"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-basic-secret"
    rank: "01"
    title: "Basic Secret"
    excerpt: "This preset stores one secret in a Key Vault with the value wired by reference -- the manifest records that the secret exists and where it lives, never what it is. The example wires a storage..."
  - slug: "02-expiring-api-key"
    rank: "02"
    title: "Expiring API Key"
    excerpt: "This preset stores a third-party API key with explicit activation and expiry attributes, making the credential's validity window -- and the rotation it implies -- auditable infrastructure."
---

# Key Vault Secret Presets

Ready-to-deploy configuration presets for Key Vault Secret. Each preset is a complete manifest you can copy, customize, and deploy.
