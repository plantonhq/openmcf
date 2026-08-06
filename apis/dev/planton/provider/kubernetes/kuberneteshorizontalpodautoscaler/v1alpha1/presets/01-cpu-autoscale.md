# CPU Autoscale

This preset is the standard workhorse autoscaler: hold the workload's average CPU at 60% of its requests, running between 2 and 10 replicas. CPU is the reliable scaling signal — it rises and falls with load — and utilization targets are measured against the pods' declared CPU requests, so the workload must set requests for the signal to mean anything.

## When to Use

- Request-serving workloads whose load tracks CPU — APIs, web frontends, most stateless services
- Operator-managed or non-Planton Deployments that need autoscaling (for a Planton Deployment's OWN replicas with simple CPU/memory scaling, prefer the workload's built-in `availability.horizontal_pod_autoscaling` block)
- As the starting point before graduating to per-container, custom, or external metrics

## Key Configuration Choices

- **`resource` metric with `utilization: 60`** — the average across all the target's pods is held at 60% of requests: scale-out above, scale-in below (with the default 300-second scale-down stabilization). Requires metrics-server in the cluster
- **`min_replicas: 2`** — the floor is the availability minimum: two replicas mean a pod loss is not an outage. The autoscaler never scales below it
- **`max_replicas: 10`** — the ceiling is a budget decision: the most this workload may cost. Required by the schema
- **`scale_target` defaults to an `apps/v1` Deployment** — only the name is needed; the name also accepts a reference to a `KubernetesDeployment`'s exported name for chart composition
- **One controller per target** — never point this AND the workload's built-in autoscaling block (or KEDA) at the same workload; two controllers fighting over one replica count flaps the fleet. The workload's own `replicas` becomes advisory once this HPA governs it

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | The scale target's own namespace — an HPA cannot scale across namespaces | Your namespace management |
| `<your-workload-name>` | The target Deployment's name | The workload's manifest |

The bounds (2–10) and target (60%) are working defaults — tune them to the workload's availability floor and cost ceiling.

## Related Presets

- **02-container-isolated** — the same signal, isolated to the app container when sidecars skew the average
- **03-queue-driven** — scale on external queue depth instead of CPU
- **04-behavior-tuned** — add conservative scale-down for spiky traffic
