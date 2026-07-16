---
title: "Self-Managed (Uploaded) Certificate"
description: "This preset uploads a PEM certificate chain and its private key as a Certificate Manager certificate. Renewal before expiry is YOUR responsibility — Google serves the material but never rotates it."
type: "preset"
rank: "03"
presetSlug: "03-self-managed-pem"
componentSlug: "certificate-manager-certificate"
componentTitle: "Certificate Manager Certificate"
provider: "gcp"
icon: "package"
order: 3
---

# Self-Managed (Uploaded) Certificate

This preset uploads a PEM certificate chain and its private key as a
Certificate Manager certificate. Renewal before expiry is YOUR
responsibility — Google serves the material but never rotates it.

## When to Use

- Certificates issued by a corporate or third-party CA that Google cannot
  manage
- Extended-validation (EV/OV) certificates
- Migrating existing certificates without re-issuing

## Key Configuration Choices

- **Self-managed arm** — exactly one of `managed` / `selfManaged` per
  certificate.
- **In-place rotation** — updating both PEM fields rotates the
  certificate without changing any consumer.
- **Pre-deploy framing validation** — swapped certificate/key material is
  rejected before anything deploys.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<gcp-project-id>` | GCP project ID | `GcpProject` outputs |
| `<leaf-certificate-pem>` | PEM body of the leaf certificate (+ intermediates below it) | Your CA |
| `<private-key-pem>` | PEM body of the leaf's private key | Your key store |

## Related Presets

- **01-managed-dns-auth** — let Google provision and renew instead
