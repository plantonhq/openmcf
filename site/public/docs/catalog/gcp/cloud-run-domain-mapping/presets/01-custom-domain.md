---
title: "Custom Domain"
description: "The default shape: a verified domain mapped onto a Cloud Run service with a managed TLS certificate — no load balancer, no certificate handling, just DNS."
type: "preset"
rank: "01"
presetSlug: "01-custom-domain"
componentSlug: "cloud-run-domain-mapping"
componentTitle: "Cloud Run Domain Mapping"
provider: "gcp"
icon: "package"
order: 1
---

# Custom Domain

The default shape: a verified domain mapped onto a Cloud Run service
with a managed TLS certificate — no load balancer, no certificate
handling, just DNS.

## What it configures

- `domain` — the verified custom domain; it becomes the mapping's name
  in GCP.
- `route` — a ValueFromRef to the GcpCloudRun service's `service_name`
  output, so the mapping follows the deployed service.
- `certificateMode: AUTOMATIC` — Cloud Run provisions and renews the
  certificate once the domain's DNS records are published.

## Adjust before deploying

- **domain** — replace with your real domain, and verify it FIRST
  (Search Console / `gcloud domains verify`); GCP rejects the mapping
  for an unverified domain. Subdomains of a verified domain need no
  separate verification.
- **region** — must match the target service's region exactly; mappings
  are regional.
- **route** — point the reference at your service's resource name, or
  set a literal service name for a service deployed outside Planton.

## After deploying

Read the `resource_records` output and publish those records in the
domain's DNS zone (wire them into GcpDnsRecord, or create them at your
external DNS host). The domain serves — and the certificate issues —
only after DNS propagates.

## When to choose something else

Migrating a domain that must keep serving from its old host until DNS
cutover? Start from the **Migration Without Certificate** preset — it
creates the mapping with `certificateMode: NONE` so certificate issuance
never blocks on DNS you haven't switched yet.
