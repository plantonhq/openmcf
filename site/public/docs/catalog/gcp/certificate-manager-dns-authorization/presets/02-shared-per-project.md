---
title: "Shared Per-Project Authorization"
description: "This preset creates a DNS authorization whose validation record is scoped per (domain, project) — the shape that lets multiple teams and projects issue certificates for the same domain without..."
type: "preset"
rank: "02"
presetSlug: "02-shared-per-project"
componentSlug: "certificate-manager-dns-authorization"
componentTitle: "Certificate Manager DNS Authorization"
provider: "gcp"
icon: "package"
order: 2
---

# Shared Per-Project Authorization

This preset creates a DNS authorization whose validation record is scoped
per (domain, project) — the shape that lets multiple teams and projects
issue certificates for the same domain without coordinating record
ownership.

## When to Use

- Several projects issue certificates for the same domain
- Regional certificates (PER_PROJECT_RECORD is the regional default)

## Key Configuration Choices

- **`type: PER_PROJECT_RECORD`** — one validation record per (domain,
  project) instead of one per authorization.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<gcp-project-id>` | GCP project ID | `GcpProject` outputs |

The sample domain `example.com` is a realistic placeholder for the
pattern-validated `domain` field — replace it with your bare domain.

## Related Presets

- **01-standard-domain** — the classic fixed-record authorization
