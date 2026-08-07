---
title: "Multi-Domain Certificate (Apex + WWW)"
description: "One Google-managed SSL certificate covering multiple hostnames — typically the apex domain and its `www` alias, or several related subdomains on the same load balancer."
type: "preset"
rank: "02"
presetSlug: "02-multi-domain"
componentSlug: "managed-ssl-certificate-on-google-cloud"
componentTitle: "Managed SSL Certificate on Google Cloud"
provider: "gcp"
icon: "package"
order: 2
---

# Multi-Domain Certificate (Apex + WWW)

One Google-managed SSL certificate covering multiple hostnames — typically the apex domain and its `www` alias, or several related subdomains on the same load balancer.

## When to Use

- Both `example.com` and `www.example.com` (or similar pairs) on the same HTTPS load balancer
- Several related subdomains that share one target HTTPS proxy
- Reducing proxy certificate list size by consolidating hostnames onto one managed cert

## Key Configuration Choices

- **Multiple entries in `domains`** — each must be a fully-qualified domain name; wildcards are not supported
- **One certificate, one proxy reference** — attach a single `self_link` to the target HTTPS proxy; all listed domains must point DNS at the same load balancer IP
- **Shared description** — documents the cert's purpose across all hostnames

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<gcp-project-id>` | GCP project ID | GCP Console or `GcpProject` outputs |
| `example.com` / `www.example.com` | Hostnames to secure | Your DNS configuration |

## Remix Notes

- Every domain in the list must have DNS pointing at the load balancer before provisioning completes
- Adding or removing a domain destroys and recreates the certificate (ForceNew) — plan create-before-destroy if a proxy already references it
- For unrelated hostnames on different load balancers, use separate certificates instead

## Related Presets

- **01-single-domain** — Simplest single-hostname certificate
- **03-explicit-name** — When the GCP certificate name must differ from the Planton resource name
