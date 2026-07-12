---
title: "Presets"
description: "Ready-to-deploy configuration presets for Container App Environment Certificate"
type: "preset-list"
componentSlug: "container-app-environment-certificate"
componentTitle: "Container App Environment Certificate"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-key-vault-certificate"
    rank: "01"
    title: "Key Vault Certificate"
    excerpt: "The production posture: the environment pulls the certificate from Azure Key Vault with its system-assigned identity and follows renewals automatically -- no rotation treadmill, no key material in..."
  - slug: "02-inline-pfx-certificate"
    rank: "02"
    title: "Inline PFX Certificate"
    excerpt: "Upload a certificate directly as a base64-encoded PKCS#12 bundle -- for certificates that live outside Key Vault (a CA-issued file, an org PKI export)."
---

# Container App Environment Certificate Presets

Ready-to-deploy configuration presets for Container App Environment Certificate. Each preset is a complete manifest you can copy, customize, and deploy.
