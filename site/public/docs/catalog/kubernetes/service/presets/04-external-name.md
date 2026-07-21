---
title: "ExternalName Alias"
description: "This preset creates a pure DNS alias: cluster DNS answers lookups for `prod-db.<namespace>.svc.cluster.local` with a CNAME to `db.prod.example.com`. No proxying, no selectors, no ports — traffic..."
type: "preset"
rank: "04"
presetSlug: "04-external-name"
componentSlug: "service"
componentTitle: "Service"
provider: "kubernetes"
icon: "package"
order: 4
---

# ExternalName Alias

This preset creates a pure DNS alias: cluster DNS answers lookups for `prod-db.<namespace>.svc.cluster.local` with a CNAME to `db.prod.example.com`. No proxying, no selectors, no ports — traffic flows directly from the pod to the external endpoint; Kubernetes is only involved at name-resolution time.

## When to Use

- Giving workloads one stable in-cluster name for a managed database, SaaS API, or service in another cluster — swap the target by editing one field instead of every consumer
- Migration staging: point `prod-db` at the external database today, replace this Service with a selecting one when the database moves in-cluster, and no client changes
- Keeping application config environment-agnostic — every environment connects to `prod-db`, and this Service decides what that means per environment

## Key Configuration Choices

- **`type: external_name`** — the only type that is a DNS alias rather than a proxy. `selector`, `ports`, and `cluster_ip_address` are meaningless here and validation rejects them
- **`external_dns_name`** — the CNAME target; must be a lowercase RFC-1123 hostname. Required for (and only allowed on) this type
- **TLS caveat** — clients connect to the alias but the server certificate names the real host; clients doing hostname verification must either verify against the target name or connect using the target name after resolution
- **Ports caveat** — a CNAME carries no port mapping; clients must use the external endpoint's real port

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | Namespace whose workloads will use the alias | Your namespace management |
| `db.prod.example.com` | The external hostname to alias (working example — replace with your endpoint) | Your database/provider console |

## Related Presets

- **01-cluster-ip-app** — the selecting Service to swap in when the endpoint moves into the cluster
