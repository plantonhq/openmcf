---
title: "Presets"
description: "Ready-to-deploy configuration presets for Key Vault Certificate"
type: "preset-list"
componentSlug: "key-vault-certificate"
componentTitle: "Key Vault Certificate"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-self-signed-auto-renew"
    rank: "01"
    title: "Self-Signed Auto-Renewing Certificate"
    excerpt: "This preset is the fully-hands-off shape for internal TLS: the vault generates the key, self-signs a 12-month certificate, and automatically re-enrolls at 80% of lifetime -- no CA account, no renewal..."
  - slug: "02-imported-bundle"
    rank: "02"
    title: "Imported Certificate Bundle"
    excerpt: "This preset imports an existing certificate -- a wildcard or EV certificate purchased from a public CA, delivered as a PFX (PKCS#12) or PEM bundle -- into the vault, which takes over private-key..."
  - slug: "03-ca-pending-csr"
    rank: "03"
    title: "CA-Signed via Pending CSR"
    excerpt: "This preset runs the out-of-band CA flow: the vault generates the private key (which never leaves it -- `exportable: false`), mints a CSR, and holds the operation pending until your CA's signed..."
---

# Key Vault Certificate Presets

Ready-to-deploy configuration presets for Key Vault Certificate. Each preset is a complete manifest you can copy, customize, and deploy.
