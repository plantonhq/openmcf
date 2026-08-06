---
kinds:
  - KubernetesNamespace
  - KubernetesKubePrometheusStack
  - KubernetesGrafana
  - KubernetesLoki
  - KubernetesTempo
  - KubernetesOtelCollector
  - KubernetesSignoz
  - KubernetesClickHouse
---

# Observability: the Assembled Stack vs the All-in-One

"Give me observability" has two honest answers in this catalog, and the
choice spans five or more kinds. This pattern is the comparison's single
home; each kind's `GUIDE.md` carries only its own judgment.

## The two shapes

- **The assembled stack** — best-of-breed pieces wired together:
  KubernetesKubePrometheusStack (metrics, alerting, the monitoring CRDs),
  KubernetesLoki (logs) fed by a KubernetesOtelCollector in daemonset
  mode, KubernetesTempo (traces) fed OTLP by apps or a collector, and
  KubernetesGrafana as the query hub reading all three. Each piece scales,
  upgrades, and fails independently; each is swappable; the operational
  surface is four components.
- **The all-in-one** — KubernetesSignoz: traces, metrics and logs in one
  UI, one query engine, one alert system — backed by a composed
  KubernetesClickHouse (plus its KubernetesAltinityOperator). One product
  to learn and operate, but the telemetry store is a real database whose
  lifecycle you own, and swapping any single concern means leaving the
  product.

| Choose | When |
|---|---|
| Assembled | The cluster already runs kube-prometheus-stack (most do); teams want Grafana; pieces must scale or be swapped independently; monitoring CRDs (ServiceMonitor et al.) are expected by other components |
| Signoz | One team wants one tool for all three signals; ClickHouse expertise exists (or the ClickHouse pair is composed anyway); minimizing the number of moving products outweighs per-piece flexibility |

Neither is a workaround — both are first-class. What is NOT first-class:
proposing pieces from both shapes for the same signals (two metrics
pipelines, two alert systems) without saying why.

## The assembled wiring, typed end to end

Every join in the assembled stack is a real reference the platform
validates and draws — this composition is fully `valueFrom`-wired:

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesNamespace
metadata:
  name: observability
spec:
  name: observability
---
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesKubePrometheusStack
metadata:
  name: monitoring
spec:
  namespace:
    valueFrom:
      kind: KubernetesNamespace
      name: observability
      fieldPath: spec.name
---
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesLoki
metadata:
  name: logs
spec:
  namespace:
    valueFrom:
      kind: KubernetesNamespace
      name: observability
      fieldPath: spec.name
---
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesTempo
metadata:
  name: traces
spec:
  namespace:
    valueFrom:
      kind: KubernetesNamespace
      name: observability
      fieldPath: spec.name
---
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesGrafana
metadata:
  name: dashboards
spec:
  namespace:
    valueFrom:
      kind: KubernetesNamespace
      name: observability
      fieldPath: spec.name
  datasources:
    - name: prometheus
      url:
        valueFrom:
          kind: KubernetesKubePrometheusStack
          name: monitoring
          fieldPath: status.outputs.prometheus_endpoint
    - name: loki
      type: loki
      url:
        valueFrom:
          kind: KubernetesLoki
          name: logs
          fieldPath: status.outputs.gateway_endpoint
    - name: tempo
      type: tempo
      url:
        valueFrom:
          kind: KubernetesTempo
          name: traces
          fieldPath: status.outputs.http_endpoint
```

Ingestion completes the picture: a KubernetesOtelCollector in daemonset
mode ships cluster logs to Loki's `gateway_endpoint`, applications (or a
gateway collector) send OTLP to Tempo's `otlp_grpc_endpoint`, and metric
scraping is declared through the stack's ServiceMonitor machinery. Loki's
`ruler.alertmanagerUrl` can reference the stack's Alertmanager, so even
log-driven alerts route through the one alerting system.

## On the diagram

The assembled shape renders as a hub: Grafana with three datasource edges
into the stack, Loki and Tempo, the collector's edge into Loki, and every
component's namespace edge into the shared observability namespace — the
telemetry topology is reviewable at a glance. The Signoz shape renders
smaller — Signoz plus its ClickHouse (and the operator in the shared
layer) — with application OTLP converging on one ingestion gateway.

## When the answer is BOTH

Signoz for one product team's application telemetry can coexist with a
cluster-level kube-prometheus-stack (platform components' ServiceMonitors
need the stack's CRDs regardless). That is a scoped decision, not
double-tooling — say which signals go where.

## See also

Each kind's guide beside its reference page: the stack (CRD singleton +
the serviceMonitor seam), Grafana (hub vs bundled; ephemeral state), Loki
(nothing ships logs by itself), Tempo (the replica/storage floor), Signoz
(the composed-ClickHouse chain), and the collector (modes and RBAC).
