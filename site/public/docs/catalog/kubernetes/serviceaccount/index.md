---
title: "ServiceAccount"
description: "ServiceAccount deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesserviceaccount"
---

# Kubernetes ServiceAccount

Deploys a Kubernetes ServiceAccount — the in-cluster identity pods run as — through a single declarative manifest. The component anchors three concerns: RBAC (grants attach to the identity), registry authentication (`imagePullSecrets`), and cloud workload identity (typed GKE/EKS/AKS federation that the IaC module translates into the exact annotations each cloud's webhook expects).

## What Gets Created

When you deploy a KubernetesServiceAccount resource, Planton provisions:

- **ServiceAccount** — a Kubernetes ServiceAccount with the configured pull secrets and token automount setting
- **Workload-identity annotations** — generated from the selected `workloadIdentity` arm: `iam.gke.io/gcp-service-account` (GKE), `eks.amazonaws.com/role-arn` (EKS), or `azure.workload.identity/client-id` plus optional `azure.workload.identity/tenant-id` (AKS)
- **Labels** — standard Planton tracking labels (`managed-by`, `resource`, `resource-kind`) merged with any user-provided labels
- **Annotations** — user-provided annotations merged after the generated ones (collisions with generated workload-identity annotations are rejected)

## Prerequisites

- **Kubernetes credentials** configured via environment variables or Planton provider config
- **A Kubernetes namespace** that already exists, or a `KubernetesNamespace` resource referenced from `spec.namespace`
- **For workload identity**: the cluster feature enabled (GKE Workload Identity, EKS OIDC provider, or the AKS workload-identity add-on) and the cloud-side trust configured on the identity being bound — the annotation alone grants nothing
- **For image pull secrets**: `kubernetes.io/dockerconfigjson` secret(s) in the same namespace, or `KubernetesSecret` resources referenced from `spec.imagePullSecrets`

## Quick Start

Create a file `service-account.yaml`:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesServiceAccount
metadata:
  name: app-identity
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesServiceAccount.app-identity
spec:
  name: app-identity
  namespace:
    value: backend
```

Deploy:

```shell
planton apply -f service-account.yaml
```

This creates a ServiceAccount named `app-identity` in the `backend` namespace. Point workloads at it with `spec.serviceAccountName: app-identity`.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `spec.name` | `string` | Name of the ServiceAccount. Workloads reference it in `spec.serviceAccountName`; RBAC subjects reference it as `system:serviceaccount:<namespace>:<name>`. | 1–253 characters, matches `^[a-z0-9]([a-z0-9.-]{0,251}[a-z0-9])?$` |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `spec.namespace` | `StringValueOrRef` | `default` | Namespace where the ServiceAccount is created. Accepts a literal name (`{ value: my-namespace }`) or a reference to a `KubernetesNamespace` resource. Participates in the RBAC name and every cloud federation subject — set explicitly. |
| `spec.imagePullSecrets` | `list<StringValueOrRef>` | `[]` | Names of `kubernetes.io/dockerconfigjson` secrets in the SAME namespace; the kubelet presents them when pulling images for pods running as this identity. Literal names or `KubernetesSecret` references. |
| `spec.automountServiceAccountToken` | `bool` (tri-state) | unset | Unset defers to the cluster/pod default (mount). `false` hardens pods that never call the kube-apiserver. `true` makes the mount explicit. Pods can override individually. |
| `spec.workloadIdentity` | `oneof` | — | Binds the ServiceAccount to at most one cloud identity: `gke`, `eks`, or `aks` (see below). Omit for identities that never leave the cluster. |
| `spec.workloadIdentity.gke.serviceAccountEmail` | `StringValueOrRef` | — | GCP service account email (e.g. `app@my-project.iam.gserviceaccount.com`). Emitted as the `iam.gke.io/gcp-service-account` annotation. Required within the `gke` arm. |
| `spec.workloadIdentity.eks.roleArn` | `StringValueOrRef` | — | AWS IAM role ARN (e.g. `arn:aws:iam::123456789012:role/app-role`). Emitted as the `eks.amazonaws.com/role-arn` annotation. Required within the `eks` arm. |
| `spec.workloadIdentity.aks.clientId` | `StringValueOrRef` | — | Client ID (GUID) of the Azure user-assigned managed identity or Entra application. Emitted as the `azure.workload.identity/client-id` annotation. Required within the `aks` arm. |
| `spec.workloadIdentity.aks.tenantId` | `string` (UUID) | — | Entra tenant ID for cross-tenant scenarios. Emitted as the `azure.workload.identity/tenant-id` annotation when set. |
| `spec.labels` | `map<string, string>` | `{}` | Additional labels merged with standard Planton labels. |
| `spec.annotations` | `map<string, string>` | `{}` | Additional annotations. Cloud workload-identity annotations belong in `spec.workloadIdentity`, not here. |

### Cross-Cloud Workload Identity

All three arms follow the same mechanism: the cluster's OIDC issuer vouches for the ServiceAccount, and the cloud exchanges that token for cloud credentials — no keys stored in the cluster.

| Arm | Annotation(s) produced | Cloud-side trust required |
|-----|------------------------|---------------------------|
| `gke` | `iam.gke.io/gcp-service-account: <email>` | `roles/iam.workloadIdentityUser` binding for `serviceAccount:<project>.svc.id.goog[<namespace>/<name>]` |
| `eks` | `eks.amazonaws.com/role-arn: <role-arn>` | IAM role trust policy on the cluster's OIDC provider, conditioned on `system:serviceaccount:<namespace>:<name>` |
| `aks` | `azure.workload.identity/client-id: <client-id>` (+ `azure.workload.identity/tenant-id` when set) | Federated identity credential for subject `system:serviceaccount:<namespace>:<name>` against the cluster's OIDC issuer |

> **AKS pods need one more thing:** the `azure.workload.identity/use: "true"` **pod label** — Azure's webhook only injects the token into labeled pods. GKE and EKS key entirely off the ServiceAccount.

## Examples

### GKE Workload Identity

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesServiceAccount
metadata:
  name: dns-manager
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.KubernetesServiceAccount.dns-manager
spec:
  name: dns-manager
  namespace:
    value: dns-system
  workloadIdentity:
    gke:
      serviceAccountEmail:
        value: dns-manager@my-project.iam.gserviceaccount.com
```

### EKS IRSA

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesServiceAccount
metadata:
  name: s3-reader
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.KubernetesServiceAccount.s3-reader
spec:
  name: s3-reader
  namespace:
    value: data-pipeline
  workloadIdentity:
    eks:
      roleArn:
        value: arn:aws:iam::123456789012:role/s3-reader
```

### Hardened Identity with Registry Credentials

Pull secrets attached at the identity level, token automount disabled for pods that never call the kube-apiserver:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesServiceAccount
metadata:
  name: web-frontend
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.KubernetesServiceAccount.web-frontend
spec:
  name: web-frontend
  namespace:
    value: frontend
  imagePullSecrets:
    - value: ghcr-pull-secret
  automountServiceAccountToken: false
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `serviceAccountName` | `string` | Name of the created ServiceAccount — the value workloads set in `spec.serviceAccountName` |
| `namespace` | `string` | Namespace where the ServiceAccount was created |
| `rbacSubject` | `string` | The fully-qualified `system:serviceaccount:<namespace>:<name>` string — the exact value cloud trust configuration matches on and the username the kube-apiserver sees |
| `workloadIdentityHandle` | `string` | The bound cloud identity (GCP email, IAM role ARN, or Azure client ID); empty when workload identity is not configured |

## Related Components

- [KubernetesRbac](/docs/catalog/kubernetes/rbac) — grants permissions to this identity as a ServiceAccount subject
- [KubernetesSecret](/docs/catalog/kubernetes/secret) — provides the `kubernetes.io/dockerconfigjson` secrets referenced from `spec.imagePullSecrets`
- [KubernetesNamespace](/docs/catalog/kubernetes/namespace) — provides the target namespace; reference it from `spec.namespace` to deploy both in one run
