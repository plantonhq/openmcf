---
title: "Webhook Provider (Out-of-Tree DNS)"
description: "This preset installs ExternalDNS with the webhook provider — upstream's extension architecture for every DNS provider that is not in-tree (Akamai, OVH, Scaleway, RFC2136, Hetzner, and many more). The..."
type: "preset"
rank: "05"
presetSlug: "05-webhook-provider"
componentSlug: "external-dns"
componentTitle: "External DNS"
provider: "kubernetes"
icon: "package"
order: 5
---

# Webhook Provider (Out-of-Tree DNS)

This preset installs ExternalDNS with the webhook provider — upstream's
extension architecture for every DNS provider that is not in-tree (Akamai,
OVH, Scaleway, RFC2136, Hetzner, and many more). The provider's own webhook
implementation runs as a sidecar container next to the controller, serving
the webhook API on localhost; the controller runs with
`--provider=webhook`. This works on any host cluster.

## When to Use

- Your DNS provider has no in-tree arm (`aws_route53`, `google_cloud_dns`,
  `azure_dns`, `cloudflare`) but publishes an external-dns webhook image
- You are migrating to a provider whose in-tree support upstream has moved
  out of tree

## Key Configuration Choices

- **Webhook sidecar** (`webhook.imageRepository` + `imageTag`) — the
  provider's published webhook image; consult the provider's external-dns
  documentation for the exact image and its configuration surface
- **Provider configuration via `env` / `args`** — passed to the sidecar
  as-is. Values land in chart values verbatim, so put SECRETS in Kubernetes
  Secrets and reference them through `helmValues` env entries with
  `valueFrom` instead of inlining them
- **`policy: upsert-only`** — the safe default until you've verified the
  provider's behavior; switch to `sync` for a dedicated zone
- **No `workloadIdentity`** — the webhook arm authenticates however the
  sidecar's provider requires, not through the controller's cloud identity

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<webhook-provider-image>` | The provider's external-dns webhook image repository | The DNS provider's external-dns webhook documentation |
| `<webhook-provider-image-tag>` | Image tag to pin | Same documentation / release page |
| `<PROVIDER_SETTING>` / `<value>` | Provider-specific environment variables (non-secret) | Same documentation |
| `<cluster-name>` | Unique owner ID for this instance | Your cluster naming |
| `<example.com>` | Domain suffix this instance manages | Your zone's domain |

## Related Presets

- **04-cloudflare-any-cluster** — if your provider is Cloudflare, use the
  in-tree arm instead
