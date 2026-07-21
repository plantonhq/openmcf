# Node Monitor

This preset deploys a node-metrics agent (the node-exporter shape) that observes the node itself: it joins the node's network and PID namespaces, reads `/proc` and `/sys` through read-only HostPath mounts, and serves a metrics endpoint on every node's IP via a host port. It gets the visibility it needs through two targeted Linux capabilities instead of full `privileged` mode.

## When to Use

- Node-level metrics collection (CPU, memory, disk, network per node)
- Process-visibility agents that must see host processes, not just their own container

## Key Configuration Choices

- **`hostNetwork` + `hostPid`** — the agent observes the node's real network interfaces and process table rather than the pod's namespaced view; both imply elevated trust in the image
- **`dnsPolicy: ClusterFirstWithHostNet`** — host-network pods otherwise inherit the node's resolver and lose cluster DNS; this keeps `*.svc.cluster.local` names resolving for agents that push to in-cluster sinks
- **`hostPort: 9100`** — the scrape endpoint is reachable at `<node-ip>:9100` on every node; DaemonSets have no Service, so node-IP exposure is the addressing model. Note a host port is exclusive per node — it conflicts with surge updates (`updateStrategy.maxSurge`), so this preset keeps the default one-node-at-a-time rollout
- **Privileged-lite security context** — `drop: ["ALL"]` then add back only `SYS_PTRACE` (inspect host processes via `hostPid`) and `SYS_RESOURCE` (read resource accounting); adjust to exactly what your agent documents needing, and resist `privileged: true`, which is root on the node
- **Read-only `/proc` and `/sys` mounts under `/host/...`** — mounted off-root so the agent distinguishes host stats from its own; read-only because a monitor must not write kernel interfaces
- **Control-plane toleration** — monitoring that skips the control plane is monitoring with a blind spot

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | Target namespace | Your namespace management or `KubernetesNamespace` resource |
| `<your-container-registry>/<your-node-monitor-image>` | Node monitoring agent image (e.g., `prom/node-exporter`) | Your container registry |
| `<your-image-tag>` | Image tag or version | Your registry or CI/CD pipeline output |

## Related Presets

- **01-log-collector** — HostPath-based log shipping without host namespaces
- **03-hardened-agent** — The opposite trust posture: restricted profile, no host access
