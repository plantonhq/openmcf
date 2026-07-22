---
title: "KubernetesExternalSecret"
description: "KubernetesExternalSecret deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesexternalsecret"
---

# KubernetesExternalSecret

Declares one secret sync: the External Secrets Operator reads entries from a store's backend (AWS Secrets Manager, GCP Secret Manager, Azure Key Vault, Vault/OpenBao, ...) and materializes them as a Kubernetes Secret, refreshed on an interval. Explicit per-key mapping (`data`), bulk pulls with key rewrites (`dataFrom`), and a target template for typed Secrets (TLS, dockerconfigjson) and value reshaping.

## What Gets Created

- **ExternalSecret** — named after the resource, in the spec's namespace
- **Materialized Kubernetes Secret** — created and refreshed by the operator (not by the IaC modules); named `target.name`, defaulting to the resource name — the handle workloads reference

## Prerequisites

- External Secrets Operator on the cluster (**KubernetesExternalSecretsOperator**)
- A store to sync from (**KubernetesSecretStore** or **KubernetesClusterSecretStore**)

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesExternalSecret
metadata:
  name: app-database-credentials
spec:
  namespace:
    value: my-app
  storeRef:
    name:
      value: aws-secrets-manager
    kind: ClusterSecretStore
  data:
    - secretKey: password
      remoteRef:
        key: prod/app/database
        property: password
```

## Stack Outputs

| Output | Description |
|---|---|
| `external_secret_name` | The ExternalSecret resource name |
| `namespace` | Where the ExternalSecret and its Secret live |
| `secret_name` | The materialized Secret workloads reference |

## Next Steps

Wire workload env `valueFrom` / volume `secretName` references to this resource's `secret_name` output.
