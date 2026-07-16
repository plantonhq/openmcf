---
title: "Imported CA-Issued Certificate"
description: "The core self-managed pattern: upload a certificate chain issued by your own CA (or purchased commercially) with its private key, for a global external HTTPS load balancer."
type: "preset"
rank: "01"
presetSlug: "01-imported-cert"
componentSlug: "ssl-certificate-self-managed"
componentTitle: "SSL Certificate (Self-Managed)"
provider: "gcp"
icon: "package"
order: 1
---

# Imported CA-Issued Certificate

The core self-managed pattern: upload a certificate chain issued by your own CA (or purchased commercially) with its private key, for a global external HTTPS load balancer.

## When to Use

- Wildcard certificates (`*.example.com`) — Google-managed certificates cannot issue them
- EV/OV certificates or a corporate/private CA your clients already trust
- Serving TLS before public DNS points at the load balancer (managed certs cannot finish provisioning without DNS)

## Key Configuration Choices

- **`certificate`** — the full PEM chain, leaf first, then intermediates (at least one intermediate; at most 5 certificates — GCP rejects longer chains)
- **`privateKey`** — the matching unencrypted PEM key (RSA-2048+ or ECDSA P-256); it is write-only in GCP and never appears in outputs
- **Global scope** — no `region`, so global external HTTPS proxies can attach it

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<gcp-project-id>` | GCP project ID where the certificate will live | GCP Console or `GcpProject` outputs |
| `certificate` PEM body | Your real certificate chain | Your CA / certificate vendor |
| `privateKey` PEM body | The matching unencrypted private key | Wherever the CSR keypair was generated |

## Remix Notes

- Reference this certificate's `self_link` from a target HTTPS proxy's `sslCertificates` list (use an explicit `valueFrom` kind — the list defaults to Google-managed certificates)
- Nothing renews itself: watch the `expire_time` output and rotate ahead of it
- Every field is immutable — rotation is create-before-destroy (see **03-rotation-versioned-name**)

## Related Presets

- **02-regional-cert** — The same import for regional external / internal ALB proxies
- **03-rotation-versioned-name** — Version the GCP name to make rotations explicit
