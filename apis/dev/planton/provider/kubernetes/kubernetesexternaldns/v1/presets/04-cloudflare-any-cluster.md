# Cloudflare DNS from Any Cluster

This preset installs ExternalDNS publishing to Cloudflare — the canonical
cross-cloud arm. Cloudflare has no workload-identity federation with
Kubernetes clusters, so authentication is always an API token; that also
means this preset works unchanged on ANY host cluster: EKS, GKE, AKS, kind,
or self-managed. Records are proxied (orange cloud) so Cloudflare's CDN/WAF
fronts them.

## When to Use

- Any cluster whose public DNS lives in Cloudflare — regardless of where
  the cluster itself runs
- Self-managed / datacenter / kind clusters with no cloud identity at all
- When you want Cloudflare's proxy (CDN, WAF, DDoS protection) in front of
  the published records

## Key Configuration Choices

- **Token authentication** (`cloudflare.apiToken`) — the module
  materializes the token as a Kubernetes Secret wired into the controller
  as `CF_API_TOKEN`; it never appears in chart values or pod specs. Note
  the controller validates the token at first zone sync, not at startup —
  a bad token surfaces as a crash-looping pod with a Cloudflare 4xx in logs
- **`proxied: true`** — created records go through Cloudflare's proxy;
  override per resource with the
  `external-dns.alpha.kubernetes.io/cloudflare-proxied` annotation
- **`policy: upsert-only`** — never deletes; the safe default for zones
  that also hold hand-managed records. Switch to `sync` for a dedicated zone
- **No `workloadIdentity`** — token-based providers don't use it

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<cloudflare-api-token>` | API token scoped to Zone:Read + DNS:Edit on the managed zones | Cloudflare dashboard → API Tokens |
| `<cloudflare-zone-id>` | Zone ID | Cloudflare dashboard (zone overview) or `CloudflareDnsZone` outputs |
| `<cluster-name>` | Unique owner ID for this instance | Your cluster naming |
| `<example.com>` | Domain suffix this instance manages | Your Cloudflare zone |

## Related Presets

- **01-aws-route53-eks-keyless** / **02-google-cloud-dns-gke** /
  **03-azure-dns-aks** — keyless same-cloud postures
- **05-webhook-provider** — for DNS providers with no in-tree arm
