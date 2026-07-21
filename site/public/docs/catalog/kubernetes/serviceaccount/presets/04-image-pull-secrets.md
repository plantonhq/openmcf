---
title: "Image Pull Secrets with Automount Hardening"
description: "This preset creates a ServiceAccount that carries private-registry pull credentials and disables the automatic API token mount. Every pod running as this identity pulls images with the attached..."
type: "preset"
rank: "04"
presetSlug: "04-image-pull-secrets"
componentSlug: "serviceaccount"
componentTitle: "ServiceAccount"
provider: "kubernetes"
icon: "package"
order: 4
---

# Image Pull Secrets with Automount Hardening

This preset creates a ServiceAccount that carries private-registry pull credentials and disables the automatic API token mount. Every pod running as this identity pulls images with the attached credential and receives no kube-apiserver token.

## When to Use

- Workloads pulling images from a private registry (GHCR, private Docker Hub, ECR/GCR/ACR with static credentials) without repeating `imagePullSecrets` on every pod spec
- Pods that never call the Kubernetes API — the automount hardening removes an unused credential from every container filesystem (a common security-baseline requirement)

## Key Configuration Choices

- **`imagePullSecrets: [{value: my-registry-cred}]`** — names a `kubernetes.io/dockerconfigjson` secret in the SAME namespace; the kubelet presents it for every image pull by pods running as this identity. Accepts a literal name (as here) or a reference to a KubernetesSecret resource, so the credential and the identity can deploy in one chart
- **`automountServiceAccountToken: false`** — pods running as this identity get no projected API token. An unused mounted token is pure attack surface: any container compromise hands the attacker a cluster credential. The field is tri-state — unset defers to the cluster default (mount), and an individual pod can still override with `automountServiceAccountToken: true` in its own spec if one workload does need API access
- **Attachment at the identity level** — attaching the pull secret here covers every current and future pod running as the ServiceAccount, instead of maintaining `imagePullSecrets` in each pod template

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | Namespace for the ServiceAccount (the pull secret must live in the same namespace) | Your namespace management |
| `my-registry-cred` | Name of the `kubernetes.io/dockerconfigjson` secret holding registry credentials | Your KubernetesSecret resource (docker-registry preset) or existing cluster secret |

## Related Presets

- **01-basic** — identity with no pull secrets or federation
- **02-workload-identity-gke** — bind the identity to a GCP service account
- **03-workload-identity-eks-irsa** — bind the identity to an AWS IAM role
