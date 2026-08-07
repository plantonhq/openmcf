# Deploying the OpenTelemetry Collector on Kubernetes: The Pipeline Is the Product

## Introduction: One Binary, Unbounded Configurations

The OpenTelemetry Collector is the vendor-neutral workhorse of modern observability: a single binary that receives telemetry (logs, metrics, traces) in dozens of protocols, processes it (batching, memory protection, Kubernetes metadata enrichment, filtering), and exports it to wherever it needs to go. It is deliberately NOT an application with a fixed shape — it is a pipeline engine, and the pipeline you configure is the product you deploy.

This creates an unusual abstraction problem. Most infrastructure components have a bounded configuration surface that a platform can model as typed fields: a database has storage, replicas, and credentials; an ingress controller has classes and default certificates. The collector's configuration surface is **unbounded by design** — the OpenTelemetry component registry contains hundreds of receivers, processors, exporters, connectors and extensions, and custom distributions add more. Any attempt to model "the useful subset" as typed fields would be obsolete before it shipped and would turn the platform into a translation layer that must chase upstream forever.

This document explains how the collector is deployed on Kubernetes, why the operator pattern is the right vehicle for it, and why Planton's `KubernetesOtelCollector` treats the collector's own configuration document as a first-class field rather than hiding it.

## The Deployment Landscape

### Raw manifests: four workload shapes, hand-maintained

The collector legitimately runs in four different Kubernetes shapes, and real clusters usually need more than one:

- A **Deployment** for the gateway pattern — a scalable, load-balanced tier that applications push OTLP to, which fans out to backends.
- A **DaemonSet** for node-local collection — container log files under `/var/log/pods`, host metrics, and kubelet stats only exist on each node, so something must run on every node to read them.
- A **StatefulSet** when stable pod identities matter — the target allocator for Prometheus scraping, or persistent sending queues.
- A **sidecar** injected into application pods for workloads that need a collector endpoint on localhost.

Maintaining these by hand means maintaining four different YAML stacks that differ in workload kind, Service wiring, volume mounts and RBAC — plus a ConfigMap for each carrying the collector config, with checksums or restarts wired up so config changes actually roll the pods.

### The Helm chart: an installer for one shape at a time

The community `opentelemetry-collector` chart parameterizes one collector installation. It works, but each release is one workload with one config, and the Day 2 story is Helm's generic one: no admission validation of the collector config (a typo ships to the cluster and crash-loops), no automatic Service port derivation, no sidecar injection at all.

### The operator: collectors as custom resources

The OpenTelemetry Operator introduces the `OpenTelemetryCollector` CRD and turns each collector into a single declaration. From one CR the operator creates and continuously reconciles the workload (Deployment, DaemonSet, StatefulSet — or registers the config for sidecar injection), the collector Service, a headless Service for gRPC-aware clients, a monitoring Service for the collector's own metrics on port 8888, and the rendered config ConfigMap.

Three operator behaviors matter most in practice:

1. **Admission validation.** The operator's webhook validates the collector configuration's shape at apply time — the class of "typo'd config crash-loops in production" errors becomes an apply-time rejection.
2. **Image injection.** A CR that declares no image gets the operator's default collector image (the `opentelemetry-collector-k8s` distribution at the operator's paired version). The fleet upgrades when the operator does, and a fleet-wide override exists on the operator kind's `default_collector_image`.
3. **Service port derivation.** The operator reads the declared receivers and derives the collector Service's ports from them — declare the standard `otlp` receiver and the Service carries 4317/4318 without any port plumbing. Only receivers the operator cannot infer need explicit ports.

### v1alpha1 vs v1beta1: the config becomes structured

The CRD's `v1alpha1` version carried the collector config as an opaque **string**. The `v1beta1` storage version carries it as a **structured object** — the API server and the operator see the document's actual shape, not a blob. Planton's modules target `v1beta1` and parse `config_yaml` up front (Terraform's `yamldecode` at plan; the Pulumi module's unmarshal before preview): a document that is not valid YAML fails loudly before anything touches the cluster, and the operator's webhook then validates the collector semantics at apply.

## Planton's Approach: The Operator's Surface, as Two Kinds

Planton models the collector the way the operator itself splits the problem. `KubernetesOtelOperator` installs the engine — the controller, its webhooks, and the CRDs — once per cluster (it watches every namespace). `KubernetesOtelCollector` declares each collector: the one `opentelemetry.io/v1beta1` CR the operator reconciles. A running `KubernetesOtelOperator` is the hard prerequisite; nothing reconciles the declaration without it.

### The pipeline document is a first-class field, not an escape hatch

`config_yaml` is required and carries the collector's own configuration document — receivers, processors, exporters, connectors, extensions, and the service pipelines wiring them together — on the collector's own open contract. This is a deliberate inversion of the usual platform posture: for a component whose registry is unbounded by design, the upstream document IS the right grain. What the spec adds around it are the Kubernetes-shaped decisions the document cannot express:

- **`mode`** — deployment (default), daemonset, statefulset, or sidecar.
- **`replicas` / `autoscaler`** — scaling, in the modes where scaling means anything.
- **`env` / `env_from_secrets`** — the credential path (below).
- **`volumes`** — host and Secret/ConfigMap mounts the pipelines need.
- **`resources`, `scheduling`, `service_account`, `image`, `additional_ports`** — the pod-level surface.

### Scaling semantics are mode semantics — validated, not documented

`replicas` and `autoscaler` (the operator manages an HPA from `min_replicas`/`max_replicas`/CPU/memory targets) apply only to deployment and statefulset modes. A daemonset runs one collector per node — its "replica count" is the node count. A sidecar runs inside the target pods — its count is the annotated-pod count. Declaring scaling fields in those modes is a validation error at apply time, mirroring the CRD's own admission rules rather than letting the operator reject the CR later.

Sidecar mode carries one more mirrored rule: `scheduling.tolerations` and `scheduling.priority_class_name` are rejected there, because in sidecar mode the collector runs inside the TARGET pods, whose scheduling this CR does not control — there is no standalone pod for those fields to apply to. And `autoscaler` excludes a non-default `replicas` everywhere: the autoscaler manages the count.

### Credentials ride environment variables, never the document

The rendered config lands in a ConfigMap — plaintext, readable by anything that can read ConfigMaps in the namespace. So the spec's secrets contract is absolute: never inline credentials in `config_yaml`. Name existing Secrets in `env_from_secrets` and each key of each Secret becomes an environment variable in the collector container; reference them in the config as `${env:VAR_NAME}`. The collector expands the reference at start — the token exists in the Secret and in the running process, never in the rendered document. Plain (non-secret) variables ride `env` the same way.

### The volumes model: one entry, both halves

The CRD splits pod storage into two parallel lists — `volumes` (the pod volume sources) and `volumeMounts` (the container mounts) — that must agree by name. The spec models them as ONE list: each `volumes` entry carries a name, exactly one source (`host_path`, `empty_dir`, `config_map`, `secret`, or `pvc`), and its own `mount_path`/`read_only`. The modules split each entry into the CR's two halves, so the two lists can never disagree.

The canonical use is daemonset log collection: a `hostPath` volume for `/var/log/pods` mounted read-only (the filelog receiver reads every container's log files), plus a second writable `hostPath` for the receiver's checkpoint state so a restarted collector resumes where it left off instead of re-reading every file.

### `additional_ports`: only what the operator cannot infer

Because the operator derives Service ports from the declared receivers, most configs need no port declarations at all — OTLP, jaeger, zipkin and prometheus receivers all produce their ports automatically. `additional_ports` exists for the remainder: a receiver the operator cannot infer (a syslog listener on a custom UDP port, for example) declares its Service port explicitly, with name, port and protocol.

### Permissions: cluster-reading receivers need real RBAC

Receivers that read cluster state — `k8s_events`, `kubeletstats`, `k8s_cluster`, and `filelog` pipelines using the `k8sattributes` processor for Kubernetes enrichment — need permissions the operator's default ServiceAccount does not carry. The composition is deliberate: declare a `KubernetesServiceAccount` and a `KubernetesRbac` granting exactly what the pipeline reads, and set `service_account` to the composed account. RBAC stays a first-class, auditable declaration instead of a side effect.

## Example: A Traces Gateway

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesOtelCollector
metadata:
  name: traces-gateway
spec:
  namespace:
    value: observability
  configYaml: |
    receivers:
      otlp:
        protocols:
          grpc:
            endpoint: 0.0.0.0:4317
          http:
            endpoint: 0.0.0.0:4318
    processors:
      memory_limiter:
        check_interval: 1s
        limit_mib: 400
        spike_limit_mib: 100
      batch: {}
    exporters:
      otlp:
        endpoint: tempo.observability.svc.cluster.local:4317
        tls:
          insecure: true
    service:
      pipelines:
        traces:
          receivers: [otlp]
          processors: [memory_limiter, batch]
          exporters: [otlp]
  replicas: 2
  resources:
    requests:
      cpu: 200m
      memory: 512Mi
    limits:
      memory: 512Mi
```

Note what is absent: no image (the operator injects its default), no Service or port declarations (derived from the `otlp` receiver), no ConfigMap plumbing (the operator renders and rolls it). And note the sizing discipline: the container memory limit and the `memory_limiter` processor's `limit_mib` agree, so under pressure the collector sheds load and reports it instead of being OOM-killed silently.

## Production Practices

- **Always declare the standard `otlp` receiver.** The exported `otlp_grpc_endpoint` and `otlp_http_endpoint` outputs — what other components compose against — assume it on 4317/4318. Every preset declares it.
- **Size memory with the `memory_limiter`.** Set the processor's `limit_mib` below the container memory limit; the collector then refuses data under pressure (visible, retryable) instead of OOMing (invisible, lossy).
- **Keep `metadata.name` at 42 characters or fewer.** The operator derives child names by suffixing (`-collector-monitoring` is the longest stable suffix at 21 characters) and Kubernetes caps names at 63. Both modules fail loudly past the budget.
- **Scrape the monitoring Service.** The operator creates `<name>-collector-monitoring` on port 8888 carrying the collector's own metrics — queue depths, dropped spans, exporter failures.
- **Compose the backends by reference.** A `KubernetesLoki`'s exported OTLP push endpoint (its gateway's `/otlp` route) is where a log pipeline's `otlphttp` exporter points; a Tempo install's OTLP gRPC endpoint is where a trace pipeline's `otlp` exporter points.

## Conclusion

The collector's configuration surface is unbounded on purpose, and the honest abstraction embraces that: `KubernetesOtelCollector` carries the pipeline document on the collector's own contract and models around it exactly the decisions Kubernetes adds — mode, scaling, credentials, volumes, scheduling — with the CRD's own admission rules mirrored as apply-time validation. The operator (installed once, by `KubernetesOtelOperator`) turns each declaration into a reconciled workload with derived Services and a validated config. The pipeline is the product; everything else is delivery.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
