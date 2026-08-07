---
title: "Standard Domain Authorization"
description: "This preset creates a global DNS authorization for one domain — the standard building block for issuing Google-managed certificates before traffic serves."
type: "preset"
rank: "01"
presetSlug: "01-standard-domain"
componentSlug: "cert-manager-dns-authorization"
componentTitle: "Cert Manager DNS Authorization"
provider: "gcp"
icon: "package"
order: 1
---

# Standard Domain Authorization

This preset creates a global DNS authorization for one domain — the
standard building block for issuing Google-managed certificates before
traffic serves.

## When to Use

- Before creating any DNS-authorized `GcpCertManagerCert`
- One per distinct domain (the wildcard is covered implicitly)

## Key Configuration Choices

- **Global location (default)** — pairs with global certificates for
  classic external HTTPS load balancers.
- **Type left to GCP's default** — `FIXED_RECORD` for global
  authorizations.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<gcp-project-id>` | GCP project ID | `GcpProject` outputs |

The sample domain `example.com` is a realistic placeholder for the
pattern-validated `domain` field — replace it with your bare domain (no
`*.` prefix, no trailing dot).

## Related Presets

- **02-shared-per-project** — share one validation record across projects
