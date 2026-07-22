---
title: "KubernetesClusterIssuer"
description: "KubernetesClusterIssuer deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesclusterissuer"
---

# KubernetesClusterIssuer

Creates one cert-manager ClusterIssuer — a cluster-wide certificate signing authority that Certificates in ANY namespace can request from. Full signing-backend surface: ACME (Let's Encrypt and friends, HTTP-01 + DNS-01 across nine providers), CA, self-signed, and Vault.

## What Gets Created

- **ClusterIssuer** — named after the resource; the name Certificates and `cert-manager.io/cluster-issuer` annotations reference
- **Credential Secrets** — API tokens and static keys declared in the spec are materialized as Secrets in cert-manager's cluster-resource namespace (never hand-created)

## Prerequisites

- cert-manager on the cluster (**KubernetesCertManager**)
- For keyless Route53/Cloud DNS/Azure DNS: workload identity configured on the cert-manager controller

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesClusterIssuer
metadata:
  name: letsencrypt-production
spec:
  certManagerNamespace:
    value: cert-manager
  config:
    acme:
      email: platform@example.com
      solvers:
        - dns01:
            cloudflare:
              apiToken:
                token: <cloudflare-api-token>
```

## Stack Outputs

| Output | Description |
|---|---|
| `cluster_issuer_name` | The issuer handle Certificates reference |
| `secrets_namespace` | Where credential Secrets were materialized |
| `acme_account_key_secret_name` | ACME account key Secret (empty for non-ACME) |

## Next Steps

Create **KubernetesCertificate** resources referencing this issuer's `cluster_issuer_name` output.
