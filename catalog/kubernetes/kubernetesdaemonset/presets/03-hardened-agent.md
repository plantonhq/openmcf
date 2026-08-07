# Hardened Agent

This preset deploys a per-node agent that passes the Kubernetes restricted Pod Security Standard: non-root with a pinned UID, read-only root filesystem (with a writable EmptyDir for /tmp), all Linux capabilities dropped, the runtime-default seccomp filter, and no API token mount. It is the template for agents that need per-node presence but NOT host access — network egress proxies, per-node caches, config watchers.

Identity is composed — pods run as a `KubernetesServiceAccount` resource referenced by output. If the agent needs Kubernetes API permissions, grant them with a KubernetesRbac resource targeting that identity and set `automountServiceAccountToken` to `true`.

## When to Use

- Per-node agents in clusters that enforce Pod Security Standards
- Multi-tenant clusters where DaemonSets from different teams share nodes
- Any node agent that does not need HostPath, host namespaces, or elevated capabilities

## Key Configuration Choices

- **No host access at all** — no HostPath mounts, no `hostNetwork`/`hostPid`, no added capabilities; a compromised agent is contained to its own pod. Contrast with `02-node-monitor`, which trades containment for node visibility
- **`runAsNonRoot` + pinned UID/GID 10001** — refuses to start images that silently default to root
- **`readOnlyRootFilesystem` + EmptyDir /tmp** — the container cannot modify its own image; the size-limited EmptyDir covers legitimate scratch writes
- **`drop: ["ALL"]` + `RuntimeDefault` seccomp** — the restricted-profile baseline
- **`automountServiceAccountToken: false`** — an agent that never calls the Kubernetes API should not carry credentials for it
- **ServiceAccount by reference** — the identity deploys with the chart; permissions and any workload-identity binding live on the identity, not the workload

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | Target namespace | Your `KubernetesNamespace` resource |
| `<your-container-registry>/<your-agent-image>` | Agent container image | Your container registry |
| `<your-image-tag>` | Image tag or version | Your registry or CI/CD pipeline output |
| `<your-service-account-resource>` | The `KubernetesServiceAccount` resource pods run as | Your identity resources |

## Related Presets

- **01-log-collector** — HostPath log shipping (needs node filesystem access)
- **02-node-monitor** — Host-namespace metrics agent (needs node visibility)
