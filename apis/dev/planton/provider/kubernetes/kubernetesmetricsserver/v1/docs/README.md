# KubernetesMetricsServer: Research and Design

## Introduction

metrics-server is the reference implementation of the Kubernetes resource
metrics API: a single lightweight Deployment that scrapes every kubelet's
resource endpoint and serves live CPU/memory usage back through the
aggregated API server. This component installs it from the official Helm
chart (`metrics-server` at
`https://kubernetes-sigs.github.io/metrics-server/`; the pinned default
chart 3.13.1 ships metrics-server 0.8.1).

## The Resource-Metrics Pipeline

The pipeline has four hops, and this component is the middle of it:

1. **Kubelet** — each node's kubelet exposes instantaneous CPU/memory usage
   for the node and its pods on its own resource-metrics endpoint.
2. **metrics-server** — scrapes every kubelet on a fixed interval
   (`metric_resolution`, default 15s), aggregates the values in memory.
3. **The metrics API** — metrics-server registers the cluster-wide
   `v1beta1.metrics.k8s.io` APIService; the API server proxies
   `metrics.k8s.io` requests to the metrics-server Service.
4. **Consumers** — `kubectl top nodes/pods` reads the API directly;
   HorizontalPodAutoscalers with Resource-type utilization targets poll it
   through the controller manager.

Without metrics-server, the consumers do not error out visibly: HPAs deploy
cleanly but never receive metric values and never scale, and `kubectl top`
reports the metrics API as unavailable. That silent-failure mode is why the
modules insist on verified readiness at install time (see below).

metrics-server is NOT a monitoring system: it keeps no history, serves only
the latest scraped values, and upstream explicitly cautions against using it
as a metrics source for monitoring pipelines. Metric history, dashboards,
and alerting are a monitoring stack's job (kube-prometheus-stack); custom
and external metrics for autoscaling are Prometheus-adapter / KEDA
territory.

## One Installation Per Cluster

The `v1beta1.metrics.k8s.io` APIService is a cluster-wide singleton — the
resource-metrics API has exactly one group/version, and it routes to exactly
one Service. A second installation would fight the first over that
registration. The Helm release name is therefore FIXED to `metrics-server`
(and the chart's fullname is pinned to it), never derived from
`metadata.name`: a manifest-derived name would only enable a second, broken
install. The fixed fullname also gives every chart object a deterministic
name — the Service the APIService routes to is always `metrics-server`,
which is what the outputs and install verification key off.

## Per-Environment Posture

Managed clouds differ in whether they ship a metrics API and in whether
their kubelets serve CA-signed certificates:

| Environment | Needs this component? | `kubelet_insecure_tls` |
|---|---|---|
| GKE | NO — metrics-server ships built-in | — |
| EKS | Yes | false (kubelet certs verify against the cluster CA) |
| AKS | Yes (resource metrics) | false |
| kind | Yes | **true** (self-signed kubelet certs) |
| k3s | Yes | **true** (self-signed kubelet certs) |
| kubeadm | Yes | **true** unless kubelet certificate rotation/signing is configured |
| On-prem / self-managed | Yes | usually **true** — false only where kubelet serving certs are signed by the cluster CA |

`kube-system` is the upstream-conventional namespace (the APIService is
cluster infrastructure), and `create_namespace` stays false there —
kube-system always exists. A dedicated `metrics-server` namespace with
`create_namespace: true` also works and keeps the install's objects
separable.

## The TLS Trust Model — Two Independent Sides

metrics-server sits in the middle of two TLS relationships, and the spec
keeps them separate because they are decided separately:

### Scrape side: metrics-server verifying kubelets

`kubelet_insecure_tls` (`--kubelet-insecure-tls`) controls whether
metrics-server verifies each kubelet's serving certificate. On clusters
where kubelets serve self-signed certificates (kind, k3s, kubeadm without
kubelet cert rotation, many on-prem setups), verification can never succeed
and the flag is REQUIRED. On EKS/AKS the kubelet certificates verify against
the cluster CA — leave it false.

A wrong posture here fails loudly, by design: both modules install
atomically and wait for the Deployment, and the chart's readiness probe
(`/readyz`) only passes once the first kubelet scrape succeeds. A
metrics-server that cannot scrape kubelets fails THIS deploy with a
readiness timeout instead of surfacing weeks later as HPAs that mysteriously
never scale.

`kubelet_preferred_address_types` is the scrape side's addressing knob: the
order of node address types metrics-server tries when connecting to a
kubelet (upstream default `InternalIP,ExternalIP,Hostname`; vocabulary
`InternalIP`, `ExternalIP`, `Hostname`, `InternalDNS`, `ExternalDNS`).
Reorder it on clusters where the default resolves to unreachable addresses.
`host_network` covers the class of clusters where the API server cannot
reach pod IPs over the overlay network (the upstream example: Weave CNI on
EKS).

### Serving side: the API server verifying metrics-server

When the API server proxies a `metrics.k8s.io` request, it calls the
metrics-server Service over TLS. Two spec blocks govern that relationship:

- **`tls`** decides how metrics-server's serving certificate is provisioned:
  - `self_signed` (default) — metrics-server generates its own self-signed
    certificate at startup. Pairs with
    `api_service.insecure_skip_tls_verify: true` (the chart default), since
    there is no stable CA to pin.
  - `helm` — Helm generates a self-signed certificate at install time and
    wires the APIService `caBundle` to it automatically: a verified chain
    without cert-manager, at the cost of Helm owning renewal.
  - `cert_manager` — cert-manager issues and renews the certificate, and its
    CA injector maintains the APIService `caBundle` (the chart stamps the
    `cert-manager.io/inject-ca-from` annotation). Requires
    KubernetesCertManager on the cluster. `cert_manager_issuer` selects an
    existing Issuer (namespaced, must live in the installation namespace) or
    ClusterIssuer; left empty, the chart creates its own self-signed Issuer
    and root Certificate chain in the installation namespace.
  - `existing_secret` — reuse a `kubernetes.io/tls` Secret you manage
    (`existing_secret_name`).
- **`api_service`** decides how the API server treats that certificate:
  `insecure_skip_tls_verify` defaults true (matching the default self-signed
  certificate); set it false only when the serving certificate chains to a
  CA the API server can see — via the cert-manager CA injector, the `helm`
  type's automatic wiring, or an explicit `ca_bundle`. `create: false` opts
  out of the APIService registration entirely, for the rare cluster where
  something else manages that object (the `api_service_name` output goes
  empty to match).

The spec's CEL rules keep the arms honest: `cert_manager_issuer` is only
accepted with type `cert_manager`, `existing_secret_name` only with type
`existing_secret`, and type `existing_secret` requires the secret name.

## Sizing and Resolution

The chart's default requests are cpu=100m, memory=200Mi with no limits —
upstream sizes that for clusters up to roughly 100 nodes. Memory scales with
cluster size; the working rule of thumb is ~1 MiB per node plus ~2 KiB per
pod, and upstream's scaling guidance adds roughly 1m core and 2 MiB of
memory per node beyond the 100-node envelope. Set `resources` accordingly on
large clusters; avoid tight memory limits, since an OOM-killed
metrics-server takes the metrics API down with it.

`metric_resolution` (default 15s) is the freshness/cost dial: HPA
responsiveness is built around the 15s default, values under 15s burn
kubelet CPU for little gain, and values over 60s make autoscaling sluggish.

For availability, `replicas: 2` plus `pod_disruption_budget: true` (a
PodDisruptionBudget with minAvailable 1) keeps the metrics API serving
through node drains — the APIService fails over between healthy replicas.
The PDB is meaningful only with replicas > 1; with a single replica it would
block every voluntary eviction and wedge node drains. The pods default to
the `system-cluster-critical` PriorityClass — the metrics pipeline should
outlive workload evictions.

## Own Telemetry

`prometheus.enabled` exposes metrics-server's OWN `/metrics` endpoint
(scrape performance, request latencies) — telemetry ABOUT metrics-server,
unrelated to the resource metrics it serves. `prometheus.service_monitor`
adds a ServiceMonitor for scrape discovery; it requires the Prometheus
operator CRDs (e.g. kube-prometheus-stack) on the cluster, and the release
FAILS to install without them — which is why a CEL rule also refuses a
ServiceMonitor without the metrics endpoint it would scrape.

## Typed Surface vs Escape Hatch

The typed spec covers the chart's meaningful configuration surface:
namespace and lifecycle, chart version, replicas + PDB, the two TLS sides,
kubelet addressing and resolution, host networking, resources, scheduling
(node selector, tolerations, priority class), own telemetry, and the image
override (the air-gapped/mirror knob — prefer bumping `chart_version` over
pinning `image.tag`, which decouples the binary from the chart's tested
pair).

One rendering detail is deliberate: the chart concatenates `defaultArgs` +
`args` into the container command line, and metrics-server's flag parsing
lets the last duplicate silently win. Rather than appending duplicate flags,
the modules OWN the chart's `defaultArgs` list — they re-render it with the
typed substitutions (`kubelet_preferred_address_types`,
`metric_resolution`) applied, so the pod spec stays canonical.

`helm_values` merges LAST with Helm `-f` semantics on both engines
(Terraform natively via the two-document values list; Pulumi module-side
with the same deep-merge): maps deep-merge with the later document winning,
lists replace. It is the escape hatch for the chart surface beyond the typed
fields — never the substitute for them.

Deliberately unmodeled as typed fields (all reachable via `helm_values`):

- **addon-resizer nanny** (`addonResizer.*`) — the chart's autosizing
  sidecar for very large clusters; niche enough that a typed arm would
  outweigh its use
- **extraVolumes / extraVolumeMounts** — no mainstream metrics-server use
  case mounts extra volumes
- **updateStrategy** — the chart's rolling-update default is correct for a
  singleton-API backend
- **dnsConfig, schedulerName, topologySpreadConstraints** — expert
  scheduling knobs beyond node selector + tolerations, which cover the real
  cases (control-plane placement)
- **tls.helm lookup knobs** (`tls.helm.lookup`, certificate duration) — the
  `helm` TLS arm works with its chart defaults; tuning certificate reuse
  semantics is an expert move that belongs in `helm_values`

## Install Semantics

Both engines install a REAL Helm release, atomically, with cleanup on fail
and a 300s timeout, waiting for the Deployment to become Available. Because
the readiness probe requires a successful kubelet scrape, a green deploy
means metrics are actually flowing — the install verifies the pipeline, not
just the pod. The module (not Helm) owns namespace creation via
`create_namespace`, so a namespace it creates carries the standard
governance labels and is deleted with the resource.

## Outputs

`namespace`, `release_name` (fixed `metrics-server`), `service_name` (the
Service the APIService routes to — the chart fullname, pinned to the release
name), and `api_service_name` (`v1beta1.metrics.k8s.io`, empty when
`api_service.create` is false — the outputs contract mirrors what actually
exists).

## E2E

Chart-default and tuned installs run on the kind cluster (a
self-signed-kubelet environment, so `kubelet_insecure_tls: true` is exactly
the posture under test), both engines. The verifier asserts all the way to
metric flow: Deployment Available, APIService Available, and `kubectl top
nodes` returning real values. The ServiceMonitor arm is proven offline (the
kind lane has no Prometheus operator CRDs — the release would fail to
install, by design). This component is also the prerequisite fixture for
scenarios that need a working metrics API, such as HorizontalPodAutoscaler
behavioral scaling.
