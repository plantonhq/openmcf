# Telemetry on Kubernetes

Defines an Istio Telemetry resource: a namespaced configuration of *how the mesh observes* the workloads it selects. It controls distributed tracing (sampling rate, reporting provider, custom span tags), metrics (provider selection and per-metric tag dimensions), and access logging (provider, on/off, and a CEL filter) -- without touching application code. This is how you dial trace sampling up or down, prune high-cardinality metric tags, and log only the requests that matter.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **A Telemetry resource** -- a namespaced Istio policy that shapes traces, metrics, and access logs for the workloads it selects.
- **Kubernetes Labels** -- resource metadata labels (resource name, kind, organization, environment) applied automatically for tracking.

## What It Controls

- **Scope** -- which workloads the configuration applies to: namespace-wide (or mesh-wide from the Istio root namespace), a label selector for a set of pods, or target references attached to a Gateway, Service, or ServiceEntry. At most one of a selector and target references is set.
- **Tracing** -- the sampling percentage, the tracing provider, custom span tags (from a literal, an environment variable, or a request header), and toggles for span reporting and Istio tags. Per traffic direction (client, server, or both).
- **Metrics** -- the metrics provider, the TCP report interval, and ordered overrides that enable/disable a metric (standard or custom) and add or remove tag (dimension) values via CEL.
- **Access logging** -- the logging provider, whether logging is on, and a CEL filter that selects which requests/connections are logged.

## Important Behavior

Every signal is optional -- an empty Telemetry resource changes nothing and the workloads inherit the mesh defaults. Multiple rules are merged by istiod, with later (more specific) entries overriding earlier ones, so list least-specific matches first. Custom span tags **fully replace** any tags inherited from a parent configuration -- they are not merged. For metric tag overrides, `UPSERT` requires a CEL value expression while `REMOVE` takes none. The workload selector is a plain runtime label match (no dependency edge) -- order this configuration after the workloads it observes with `metadata.relationships`. A target reference's **name** is a foreign key defaulting to a Gateway reference: wiring it with `valueFrom` orders this configuration after the gateway it observes, while a literal `value:` covers Services, ServiceEntries, and resources created outside Planton. Waypoint proxies require target references; label selectors are ignored by waypoints.

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** -- an active connection in the Connect module with kubeconfig credentials for the target cluster. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline kubeconfig authentication.

### Kubernetes Cluster

- **Istio installed** -- the Istio CRDs (Istio Base CRDs on Kubernetes) must be present and the Istio control plane (istiod) running. Telemetry is only generated where istiod is active.
- **Telemetry providers configured** -- the tracing/metrics/logging providers you reference by name must be declared in the mesh's `MeshConfig` extension providers. Leave providers empty to use the mesh defaults.
- **Target namespace exists** -- the resource is created in a specific namespace; reference an existing one or create it first.

## Deploy

### Console

Open the deployment store, find **Telemetry on Kubernetes**, and click **Deploy**. The creation wizard walks you through the namespace, the scope, and the three observability signals (tracing, metrics, access logging), with guidance at each step.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesTelemetry
metadata:
  name: prod-apps-observability
  org: acme-corp
  env: prod
spec:
  namespace:
    value: prod-apps
  tracing:
    - randomSamplingPercentage: 10
      customTags:
        cluster_name:
          environment:
            name: ISTIO_META_CLUSTER_ID
  metrics:
    - overrides:
        - match:
            metric: REQUEST_COUNT
          tagOverrides:
            request_path:
              operation: REMOVE
  accessLogging:
    - filter:
        expression: response.code >= 500
```

```shell
planton apply -f telemetry.yaml
```

This samples 10% of traces and tags every span with the cluster ID, drops the high-cardinality `request_path` dimension from request counts, and logs only server 5xx responses -- a sane, cost-aware observability baseline for a whole namespace.

## Key Configuration

- **Namespace** -- the namespace the resource is created in. It is fixed once created and is the default scope.
- **Scope** -- All Workloads (namespace-wide), a Workload Selector (match pod labels), or Target References (attach to a Gateway/Service/ServiceEntry).
- **Tracing** -- **sampling percentage** (0-100, in 0.01% increments), **providers**, **custom span tags** (Literal / Environment / Request Header), and advanced toggles (disable span reporting, enable Istio tags, use Request ID for sampling).
- **Metrics** -- **providers**, **reporting interval** (e.g. `5s`), and **overrides** (a standard or custom metric, a traffic direction, a disable toggle, and **tag overrides** that UPSERT or REMOVE dimensions).
- **Access Logging** -- **traffic direction**, **providers**, a **disable** toggle, and a **CEL filter** selecting which requests are logged.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Description |
|------------|-------------|
| `namespace` | The namespace the resource is created in. Reference an existing Namespace on Kubernetes or supply the name directly. |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources and operators can reference:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `telemetry_name` | Name of the created Telemetry resource (equals `metadata.name`) | Ordering resources that depend on the config being in place |
| `namespace` | The namespace the resource was created in | Confirming where the config applies |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

- **Cost-aware tracing** -- sample a small percentage of requests so traces stay useful without overwhelming the trace backend or its bill. Start by setting a tracing rule's sampling percentage.
- **Metric cardinality control** -- REMOVE a high-cardinality tag (like a request path or user ID) from a standard metric to keep your time-series database lean, or UPSERT a CEL-derived dimension you want to slice by.
- **Error-only access logs** -- attach a CEL filter such as `response.code >= 500` so only failing requests are logged, keeping log volume and cost down while preserving the signal you need to debug.

## Works With

Telemetry is part of the Istio observability family. It requires the Istio Base CRDs on Kubernetes and a running Istio control plane. It pairs naturally with the telemetry providers declared in your mesh's `MeshConfig` (OpenTelemetry, Zipkin, Prometheus), and its scope can attach to a **Gateway** or **Service Entry** via target references or match application workloads by label. To order this configuration after the workloads, Gateway, or Service it observes within an infra chart, express the dependency through `metadata.relationships`.
