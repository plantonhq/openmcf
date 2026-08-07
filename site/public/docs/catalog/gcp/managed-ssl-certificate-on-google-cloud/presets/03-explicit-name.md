---
title: "Explicit GCP Certificate Name"
description: "Use when the Planton resource name (`metadata.name`) should differ from the certificate name that appears in GCP — common during rotation workflows where you create `prod-lb-tls-2026` while the..."
type: "preset"
rank: "03"
presetSlug: "03-explicit-name"
componentSlug: "managed-ssl-certificate-on-google-cloud"
componentTitle: "Managed SSL Certificate on Google Cloud"
provider: "gcp"
icon: "package"
order: 3
---

# Explicit GCP Certificate Name

Use when the Planton resource name (`metadata.name`) should differ from the certificate name that appears in GCP — common during rotation workflows where you create `prod-lb-tls-2026` while the Planton resource stays `prod-lb-tls`.

## When to Use

- Certificate rotation: create a new GCP cert with a new `certificateName`, repoint the HTTPS proxy, then destroy the old cert
- Naming conventions where GCP resource names carry a version or date suffix
- Stable Planton resource identity with evolving cloud-side names

## Key Configuration Choices

- **`certificateName` set explicitly** — overrides the default (`metadata.name`) for the GCP API name
- **Immutable name** — changing `certificateName` destroys and recreates the certificate; treat name changes as rotation events
- **Single domain** in this preset — swap or extend `domains` for your hostname

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<gcp-project-id>` | GCP project ID | GCP Console or `GcpProject` outputs |
| `prod-lb-tls-2026` | GCP-side certificate name | Your rotation / naming convention |
| `api.example.com` | Hostname to secure | Your DNS configuration |

## Remix Notes

- Rotation pattern: deploy this preset with a new `certificateName`, update the target HTTPS proxy to reference the new `self_link`, then delete the old certificate resource
- Never change `certificateName` in place on a cert attached to a live proxy — create-before-destroy instead

## Related Presets

- **01-single-domain** — Default naming (GCP name = `metadata.name`)
- **02-multi-domain** — Multiple hostnames on one certificate
