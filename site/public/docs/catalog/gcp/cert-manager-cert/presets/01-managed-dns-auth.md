---
title: "Managed Certificate with DNS Authorization"
description: "This preset creates a Google-managed certificate validated through a referenced `GcpCertManagerDnsAuthorization`. DNS authorization lets the certificate reach ACTIVE before any traffic serves — the..."
type: "preset"
rank: "01"
presetSlug: "01-managed-dns-auth"
componentSlug: "cert-manager-cert"
componentTitle: "Cert Manager Cert"
provider: "gcp"
icon: "package"
order: 1
---

# Managed Certificate with DNS Authorization

This preset creates a Google-managed certificate validated through a
referenced `GcpCertManagerDnsAuthorization`. DNS authorization lets the
certificate reach ACTIVE before any traffic serves — the zero-downtime
migration path — and is the only validation mode that supports wildcards.

## When to Use

- Issuing a certificate for a domain before cutting traffic over to GCP
- Any certificate that must renew automatically with no serving dependency
- The standard production pattern for external Application Load Balancers

## Key Configuration Choices

- **Managed arm** — Google provisions and renews; no key material to hold.
- **DNS authorization by reference** — the authorization is its own
  composable resource; its validation record composes into the zone via a
  `GcpDnsRecord`.
- **Global scope (default)** — correct for classic external HTTPS load
  balancers.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<gcp-project-id>` | GCP project ID | `GcpProject` outputs |
| `<dns-authorization-resource-name>` | Name of the GcpCertManagerDnsAuthorization resource | Your DNS authorization manifest |

The sample domain `app.example.com` is a realistic placeholder for the
pattern-validated `domains` entries — replace it with your domain (no
trailing dot).

## Related Presets

- **02-wildcard-cert** — wildcard plus apex under one certificate
- **03-self-managed-pem** — bring-your-own PEM certificate
