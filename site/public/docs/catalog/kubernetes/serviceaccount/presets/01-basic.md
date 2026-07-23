---
title: "Basic ServiceAccount"
description: "This preset creates a plain ServiceAccount — a dedicated in-cluster identity for a workload, with no cloud federation and no pull secrets. Point workloads at it with `spec.serviceAccountName` and..."
type: "preset"
rank: "01"
presetSlug: "01-basic"
componentSlug: "serviceaccount"
componentTitle: "ServiceAccount"
provider: "kubernetes"
icon: "package"
order: 1
---

# Basic ServiceAccount

This preset creates a plain ServiceAccount — a dedicated in-cluster identity for a workload, with no cloud federation and no pull secrets. Point workloads at it with `spec.serviceAccountName` and attach permissions with a KubernetesRbac grant.

## When to Use

- Giving a workload its own identity instead of the namespace `default` ServiceAccount (the least-privilege baseline — grants to `default` leak to every pod in the namespace)
- Creating the RBAC anchor before or alongside a KubernetesRbac grant that targets it
- Workloads that stay inside the cluster: no cloud API access, public images

## Key Configuration Choices

- **Name and namespace only** — the minimal useful ServiceAccount; everything else (pull secrets, automount, workload identity) has sensible absent-defaults
- **Namespace set explicitly** — it participates in the identity's RBAC name (`system:serviceaccount:<namespace>:<name>`) and any future cloud federation subject; relying on the implicit `default` namespace is almost never intended
- **Token automount left unset** — defers to the cluster default (mount); set `automountServiceAccountToken: false` if the pods never call the kube-apiserver

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | Namespace for the ServiceAccount (must match the workloads that will run as it) | Your namespace management |

Also rename `app-identity` (in `metadata.name` and `spec.name`) to match your workload, e.g. `<app>-identity` or simply the app's name.

## Related Presets

- **02-workload-identity-gke** — bind the identity to a GCP service account
- **03-workload-identity-eks-irsa** — bind the identity to an AWS IAM role
- **04-image-pull-secrets** — attach private registry credentials and harden token automount
