# Cluster Operator

This preset grants an operator-style ServiceAccount cluster-wide read access to nodes and namespaces, plus the `/metrics` endpoint. It creates a ClusterRole with two rules and a ClusterRoleBinding — the shape used by monitoring agents, autoscalers, and controllers that observe the whole cluster.

## When to Use

- Monitoring/observability agents that scrape `/metrics` and enumerate nodes and namespaces (Prometheus-style collectors, cost analyzers)
- Any controller whose targets are cluster-scoped resources — `nodes` and `namespaces` exist outside every namespace, so no namespace-scoped Role can ever grant them

## Key Configuration Choices

- **`clusterScope: {}`** — the grant produces ClusterRole + ClusterRoleBinding. Cluster scope is REQUIRED here twice over: the resources are cluster-scoped, and non-resource URLs cannot be namespaced. The empty `{}` is the whole scope — cluster scope carries no parameters
- **Two independent rules** — a rule grants either resources or non-resource URLs, never both, so the resource read (`nodes`, `namespaces`) and the endpoint read (`/metrics`) are separate rules. Rules are independent and unordered; any matching rule allows a request
- **`nonResourceUrls: ["/metrics"]`** — grants HTTP access to the API server's metrics path. Paths like `/healthz`, `/livez`, and wildcards like `/api/*` (trailing `*` as a full segment only) follow the same pattern
- **ServiceAccount subject WITH `namespace`** — required: a ServiceAccount always lives in some namespace, and a cluster-scoped grant has no namespace to default from. The schema rejects cluster-scope ServiceAccount subjects without one
- **Cluster-wide means every namespace** — the binding grants these permissions across the whole cluster; keep the rules read-only and narrow

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-service-account>` | Name of the operator's ServiceAccount | Your KubernetesServiceAccount resource or the workload's `spec.serviceAccountName` |
| `<your-service-account-namespace>` | Namespace where that ServiceAccount lives (mandatory in cluster scope) | Your namespace management |

Also rename `cluster-monitor` (`metadata.name`, which becomes the ClusterRole name) to reflect your operator.

## Related Presets

- **01-namespace-app-reader** — namespace-confined read access for a workload
- **02-grant-builtin-view** — built-in `view` role for team access
- **04-aggregated-clusterrole** — a label-composed ClusterRole published without a binding
