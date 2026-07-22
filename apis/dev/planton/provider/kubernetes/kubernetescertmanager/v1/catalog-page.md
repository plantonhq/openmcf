# KubernetesCertManager

Installs cert-manager — the cluster's certificate machinery — from the official Helm chart, with a typed spec over the chart's meaningful configuration surface and keyless cloud-DNS identity built in. One installation per cluster serves every issuer and certificate.

## What Gets Created

- **Namespace** (optional) — the installation namespace, created and owned when `create_namespace` is set
- **Helm Release** — cert-manager controller, webhook, and cainjector, with CRDs installed by default and kept on uninstall (certificate data is never cascade-deleted)

## Prerequisites

- A Kubernetes cluster (GKE, EKS, AKS, kind, or any conformant cluster)
- For keyless DNS-01: the cloud-side identity half (IAM role / GCP service account / Azure managed identity) — composable in the same infra chart

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesCertManager
metadata:
  name: cert-manager
spec:
  namespace:
    value: cert-manager
  createNamespace: true
```

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | Installation namespace |
| `release_name` | Helm release name (always `cert-manager`) |
| `service_account_name` | Controller ServiceAccount — bind this identity cloud-side for keyless DNS-01 |
| `cluster_resource_namespace` | Where ClusterIssuer credential Secrets live |

## Next Steps

Create a **KubernetesClusterIssuer** (or namespace-scoped **KubernetesIssuer**) to define who signs certificates, then **KubernetesCertificate** resources for the certificates themselves.
