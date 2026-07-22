---
title: "KubernetesExternalSecretsOperator"
description: "KubernetesExternalSecretsOperator deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesexternalsecretsoperator"
---

# KubernetesExternalSecretsOperator

Installs the External Secrets Operator — the machinery that syncs secrets from external stores (AWS Secrets Manager, GCP Secret Manager, Azure Key Vault, Vault/OpenBao, ...) into Kubernetes Secrets — from the official Helm chart, with a typed spec over the chart's meaningful configuration surface and keyless cloud identity built in. One installation per cluster serves every store and every synced secret.

## What Gets Created

- **Namespace** (optional) — the installation namespace, created and owned when `create_namespace` is set
- **Helm Release** — ESO controller, webhook, and cert-controller, with CRDs installed by default and kept on uninstall (SecretStore/ExternalSecret objects are never cascade-deleted)

## Prerequisites

- A Kubernetes cluster (EKS, GKE, AKS, kind, or any conformant cluster)
- For keyless ambient store access: the cloud-side identity half (IAM role / GCP service account / Azure managed identity) — composable in the same infra chart

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesExternalSecretsOperator
metadata:
  name: external-secrets-operator
spec:
  namespace:
    value: external-secrets
  createNamespace: true
```

With ambient identity on EKS (one IAM role for every store without its own auth):

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesExternalSecretsOperator
metadata:
  name: external-secrets-operator
spec:
  namespace:
    value: external-secrets
  createNamespace: true
  workloadIdentity:
    eks:
      roleArn:
        value: arn:aws:iam::123456789012:role/external-secrets
```

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | Installation namespace — home for credential Secrets cluster-scoped stores read |
| `release_name` | Helm release name (always `external-secrets`) |
| `controller_service_account` | Controller ServiceAccount — bind this identity cloud-side for ambient keyless store access |

## Next Steps

Create a **KubernetesClusterSecretStore** (or namespace-scoped **KubernetesSecretStore**) to define which backend holds the secrets, then **KubernetesExternalSecret** resources for each secret to sync.
