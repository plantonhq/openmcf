---
title: "Modern TLS 1.2 Baseline"
description: "The recommended production posture: a TLS 1.2 floor with the MODERN cipher profile. Attach it to every internet-facing HTTPS proxy unless a stricter regime demands RESTRICTED."
type: "preset"
rank: "01"
presetSlug: "01-modern-tls12"
componentSlug: "ssl-policy-on-google-cloud"
componentTitle: "SSL Policy on Google Cloud"
provider: "gcp"
icon: "package"
order: 1
---

# Modern TLS 1.2 Baseline

The recommended production posture: a TLS 1.2 floor with the MODERN cipher profile. Attach it to every internet-facing HTTPS proxy unless a stricter regime demands RESTRICTED.

## When to Use

- Any production frontend that should not accept TLS 1.0/1.1 handshakes
- PCI DSS and similar compliance regimes that mandate TLS 1.2+
- A shared org-wide TLS floor referenced by many proxies

## Key Configuration Choices

- **`profile: MODERN`** — drops broken and legacy ciphers while keeping broad client compatibility (older enterprise clients on TLS 1.2 still connect)
- **`minTlsVersion: TLS_1_2`** — rejects TLS 1.0/1.1 outright; TLS 1.3 remains negotiable (GCP always enables it when the client supports it)
- **Global scope** — no `region`, so global external HTTPS proxies can attach it

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<gcp-project-id>` | GCP project ID where the policy will live | GCP Console or `GcpProject` outputs |

## Remix Notes

- Reference this policy's `self_link` from a target HTTPS proxy's `sslPolicy` field — without one, GCP applies its permissive default (TLS 1.0, COMPATIBLE)
- Tightening the profile later updates in place and applies to every referencing proxy on the next handshake
- Set `region` to create the regional twin for regional external / internal ALB proxies

## Related Presets

- **02-restricted-strict** — Narrow to ciphers with modern security guarantees
- **03-custom-cipher-list** — Hand-pick the exact cipher suites
