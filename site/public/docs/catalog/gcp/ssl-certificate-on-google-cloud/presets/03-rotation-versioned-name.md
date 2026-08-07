---
title: "Rotation with a Versioned Name"
description: "Encode the issue year (or serial) into the GCP certificate name so rotations are explicit create-before-destroy steps instead of in-place mutations that cannot work."
type: "preset"
rank: "03"
presetSlug: "03-rotation-versioned-name"
componentSlug: "ssl-certificate-on-google-cloud"
componentTitle: "SSL Certificate on Google Cloud"
provider: "gcp"
icon: "package"
order: 3
---

# Rotation with a Versioned Name

Encode the issue year (or serial) into the GCP certificate name so rotations are explicit create-before-destroy steps instead of in-place mutations that cannot work.

## When to Use

- Any production certificate you will rotate (which is every self-managed certificate — nothing renews itself)
- Teams that want the active certificate's vintage visible in the GCP console
- Avoiding name collisions during the overlap window when old and new certificates coexist

## Key Configuration Choices

- **`certificateName` with a version suffix** — every field of a compute SSL certificate is immutable, so the replacement MUST be a new resource; a versioned name lets both exist during the swap
- **`description` documents the rotation contract** — the operator reading the console knows how this certificate is replaced

## The Rotation Sequence

1. Create a new `GcpSslCertificate` resource with the next name (e.g. `prod-app-cert-2027`) and the new PEM material
2. Update the target HTTPS proxy's `sslCertificates` to reference the new certificate (attach-before-detach — GCP swaps the list in place with zero downtime)
3. Destroy the old certificate resource — only after no proxy references it (GCP blocks deletion of an attached certificate)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<gcp-project-id>` | GCP project ID where the certificate will live | GCP Console or `GcpProject` outputs |
| `prod-app-cert-2026` | Versioned cloud-side name | Your naming convention |
| PEM bodies | Your real chain and unencrypted key | Your CA |

## Related Presets

- **01-imported-cert** — The basic import without name versioning
- **02-regional-cert** — The regional variant
