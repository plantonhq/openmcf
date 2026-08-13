# Overview

The GCP DNS Zone API resource provisions a Google Cloud DNS **managed zone** — the authoritative container for a domain. DNS records are **not** embedded in this resource; use the separate **GcpDnsRecord** kind for individual record sets. This split follows the production "split state" pattern where platform teams own zones and application automation (external-dns, cert-manager) owns dynamic records.

## Why We Created This API Resource

Managing DNS zones in GCP involves visibility modes, VPC bindings, DNSSEC, forwarding, and peering — each with distinct constraints. This resource exposes the full `google_dns_managed_zone` surface while keeping record lifecycle in GcpDnsRecord.

- **Composable zones:** No bundled records or project-level IAM bindings that fight with other stacks
- **Public and private:** Standard private zones, forwarding zones, and peering zones
- **Security options:** DNSSEC, query logging, and force-destroy controls
- **Foreign-key wiring:** VPC and GKE cluster refs resolve from sibling resources

## Key Features

### Zone Types

- **Public zones** — internet-facing authoritative DNS; configure returned nameservers at your registrar
- **Private standard zones** — visible to listed VPC networks and/or GKE clusters
- **Forwarding zones** — private zones that forward queries to upstream resolvers
- **Peering zones** — private zones that peer with another VPC's Cloud DNS

### Spec Highlights

- **`dnsName`** — optional FQDN; defaults to `metadata.name` + `.` when omitted
- **`privateVisibilityConfig`** — network refs → GcpVpcNetwork `network_self_link`; GKE refs → GcpGkeCluster `cluster_id`
- **`dnssecConfig`** — enable signing, custom key specs, NSEC/NSEC3
- **`forwardingConfig` / `peeringConfig`** — private-only; mutually exclusive; each forwarding target carries an IPv4 or IPv6 resolver address (never both)
- **`deletionPolicy`** — the second destroy lever beside `forceDestroy`: `DELETE` (default), `PREVENT` (destroy refuses), or `ABANDON` (the zone keeps serving but leaves management)

### Deliberately Removed (Safety)

- **`records[]`** — use GcpDnsRecord instead
- **`iam_service_accounts`** — no project-level `roles/dns.admin` binding; grant least-privilege IAM separately

## Benefits

- **Safe composition:** Zone-only scope avoids authoritative IAM and record-set conflicts
- **Full managed-zone floor:** Matches released `google_dns_managed_zone` capabilities
- **Consistent engines:** Terraform and Pulumi both enable `dns.googleapis.com` before create

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
