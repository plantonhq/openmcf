---
title: "AWS Route 53 on EKS (Keyless via IRSA)"
description: "This preset installs ExternalDNS on an EKS cluster publishing to a Route 53 hosted zone, authenticating keylessly through IRSA — no static AWS keys anywhere. It scopes the instance to one zone,..."
type: "preset"
rank: "01"
presetSlug: "01-aws-route53-eks-keyless"
componentSlug: "external-dns"
componentTitle: "External DNS"
provider: "kubernetes"
icon: "package"
order: 1
---

# AWS Route 53 on EKS (Keyless via IRSA)

This preset installs ExternalDNS on an EKS cluster publishing to a Route 53
hosted zone, authenticating keylessly through IRSA — no static AWS keys
anywhere. It scopes the instance to one zone, enables full `sync`
reconciliation (the zone is dedicated to this cluster), and tags ownership
with a per-cluster TXT owner ID. This is the standard production posture on
AWS.

## When to Use

- EKS clusters whose Services/Ingresses publish records into a Route 53 zone
- Zones dedicated to (or safely shareable with) this cluster, where you want
  stale records deleted when workloads disappear
- Production deployments — keyless IRSA is the recommended authentication

## Key Configuration Choices

- **Keyless authentication** (`workloadIdentity.eks.roleArn`) — the
  controller ServiceAccount assumes an IAM role via IRSA; no credential
  Secrets are created
- **Zone scoping** (`zoneIdFilters` + `domainFilters`) — the guardrails
  against touching unrelated zones the role can see
- **`policy: sync`** — creates, updates, AND deletes records this instance
  owns; safe because the TXT registry limits deletes to records tagged with
  this instance's `txtOwnerId`
- **`txtOwnerId`** — distinct per instance sharing a zone; what stops one
  instance from deleting another's records
- **Sources `service` + `ingress`** — the chart defaults, stated explicitly;
  add Gateway API route sources if routes carry your hostnames

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<route53-hosted-zone-id>` | Hosted zone ID (e.g. `Z104533312EOZ72FQZ4TT`) | Route 53 console or `AwsRoute53Zone` outputs |
| `<external-dns-irsa-role-arn>` | IAM role with Route 53 permissions, trusting the cluster OIDC provider for `system:serviceaccount:external-dns:my-external-dns-route53` | IAM console or `AwsIamRole` outputs |
| `<cluster-name>` | Unique owner ID for this instance | Your cluster naming |
| `<example.com>` | Domain suffix this instance manages | Your zone's domain |

## Related Presets

- **02-google-cloud-dns-gke** — the same posture on GKE + Cloud DNS
- **04-cloudflare-any-cluster** — publish to Cloudflare from any cluster (including EKS) with a token
