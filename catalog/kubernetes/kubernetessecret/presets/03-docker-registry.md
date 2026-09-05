# Docker Registry Secret

This preset creates a Docker registry authentication secret (`kubernetes.io/dockerconfigjson`) for pulling images from private container registries. A workload names it in `pod.imagePullSecrets` (by reference, so it deploys first), or a `KubernetesServiceAccount` attaches it for every pod that runs as that identity.

## When to Use

- One registry login shared by several workloads in a namespace (a login that belongs to a single workload is declared on that workload's `pod.imageRegistries` instead — the module then owns the Secret)
- Kubernetes nodes whose own identity does not reach the registry (a same-cloud registry — EKS to ECR, GKE to Artifact Registry, AKS to ACR — needs no pull secret at all)
- Multi-registry environments where different namespaces need different credentials

## Key Configuration Choices

- **Docker config JSON type** -- Kubernetes-native `kubernetes.io/dockerconfigjson` type; automatically used by kubelet for image pulls
- **Password by reference** -- the password field accepts only a `$secret/<slug>` reference to an organization secret; the runner resolves it inside your infrastructure at deploy time, so the token never sits in a manifest or a record

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | Target namespace for the secret | Your namespace management |
| `<your-registry-server>` | Registry server URL (e.g., `https://index.docker.io/v1/`, `ghcr.io`, `123456789.dkr.ecr.us-east-1.amazonaws.com`) | Your container registry settings |
| `<your-registry-username>` | Registry username or access key | Your registry's credential management |
| `$secret/replace-with-your-registry-token-secret` | The organization secret holding the registry password or access token (prefer a scoped read-only token) | `planton secret` in your organization |

## Related Presets

- **01-opaque** -- Generic key-value secret for credentials and API keys
- **02-tls** -- TLS certificate and key pair
