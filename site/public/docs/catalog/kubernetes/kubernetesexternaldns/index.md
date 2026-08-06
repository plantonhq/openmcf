---
title: "KubernetesExternalDns"
description: "KubernetesExternalDns deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesexternaldns"
---

# KubernetesExternalDns

Installs ExternalDNS — the controller that publishes DNS records for your Services, Ingresses, and Gateway API routes — from the official Helm chart, with a typed spec over the chart's meaningful configuration surface. One installation manages one DNS provider; keyless cloud identity is built in, and cross-cloud combinations (EKS cluster publishing to Cloudflare, GKE cluster publishing to Route 53) are first-class.

## What Gets Created

- **Namespace** (optional) — the installation namespace, created and owned when `create_namespace` is set
- **Credential Secrets** (when needed) — declared static credentials (Cloudflare token, AWS keys, GCP key, Azure config) materialize as Kubernetes Secrets wired into the controller; keyless installs materialize nothing
- **Helm Release** — the ExternalDNS controller (plus, on the webhook arm, the provider's webhook sidecar), named after `metadata.name` so multiple instances coexist per cluster

## Prerequisites

- A Kubernetes cluster (EKS, GKE, AKS, kind, or any conformant cluster)
- A DNS zone in the target provider (Route 53, Cloud DNS, Azure DNS, Cloudflare, or a webhook-served provider)
- For keyless access: the cloud-side identity half (IAM role / GCP service account / Azure managed identity) — composable in the same infra chart

## Provider Arms

- **AWS Route 53** — keyless on EKS via IRSA; assume-role for cross-account zones; static keys as fallback
- **Google Cloud DNS** — keyless on GKE via Workload Identity; service-account key as fallback
- **Azure DNS** — public or private zones; Workload Identity, managed identity, or service principal
- **Cloudflare** — token-authenticated on ANY cluster; optional proxied (orange-cloud) records
- **Webhook** — upstream's extension architecture for every out-of-tree provider
- **In-memory** — sandbox provider for evaluating sources/filters/policies without a real zone

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesExternalDns
metadata:
  name: external-dns
spec:
  namespace:
    value: external-dns
  createNamespace: true
  awsRoute53:
    region: us-east-1
  workloadIdentity:
    eks:
      roleArn:
        value: arn:aws:iam::123456789012:role/external-dns
  policy: sync
  txtOwnerId: my-cluster
```

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | Installation namespace |
| `release_name` | Helm release name (equals `metadata.name`) |
| `service_account_name` | Controller ServiceAccount — bind this identity cloud-side for keyless provider access |

## Next Steps

Annotate Services and Ingresses with hostnames (or let their rules speak for themselves) and ExternalDNS publishes the records. Pair with **KubernetesCertManager** for certificates over the same names, and give each instance sharing a zone its own `txt_owner_id`.
