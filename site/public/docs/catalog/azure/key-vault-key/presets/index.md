---
title: "Presets"
description: "Ready-to-deploy configuration presets for Key Vault Key"
type: "preset-list"
componentSlug: "key-vault-key"
componentTitle: "Key Vault Key"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-rsa-cmk"
    rank: "01"
    title: "RSA Customer-Managed Key"
    excerpt: "This preset creates the baseline CMK: a software-protected RSA-2048 key whose capability list is exactly the envelope-encryption pair (`WRAP_KEY` + `UNWRAP_KEY`). Every Azure CMK integration --..."
  - slug: "02-hsm-auto-rotating"
    rank: "02"
    title: "HSM Auto-Rotating Production CMK"
    excerpt: "This preset creates the production-hardened CMK: HSM-protected key material (FIPS 140-2 Level 3, never leaves the hardware module) with a fully automatic rotation policy -- Azure mints a new version..."
  - slug: "03-ec-signing"
    rank: "03"
    title: "EC Signing Key"
    excerpt: "This preset creates an elliptic-curve signing key: the private key signs inside the vault (`SIGN`), anyone with the public half verifies (`VERIFY`). The `public_key_pem` output exports the public..."
---

# Key Vault Key Presets

Ready-to-deploy configuration presets for Key Vault Key. Each preset is a complete manifest you can copy, customize, and deploy.
