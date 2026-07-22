# KubernetesKeda: Research and Design

## Introduction

KEDA (Kubernetes Event-Driven Autoscaling) extends the Kubernetes
autoscaling model from "how hot is the pod" to "how much work is waiting":
70+ scalers read queue depths, stream lags, database rows, cron schedules,
and cloud metric sources, and drive workloads up — and down to ZERO, which
plain HPA cannot do. This component installs the engine from the official
Helm chart (`keda` at `https://kedacore.github.io/charts`; the pinned
default chart 2.20.1 ships KEDA 2.20.1 — chart and app versions move
together).

## Upstream Architecture

An installation is three Deployments, all named FIXED by the chart's
default values regardless of the release name:

1. **`keda-operator`** — the controller. It watches
   ScaledObject/ScaledJob resources, evaluates their triggers on a loop,
   and drives each workload's HPA (creating and owning it), including the
   0↔1 transition plain HPA cannot make.
2. **`keda-operator-metrics-apiserver`** — the external metrics API
   server. It registers the cluster-wide
   `v1beta1.external.metrics.k8s.io` APIService; the HPA controller reads
   trigger values through it on every reconcile loop.
3. **`keda-admission-webhooks`** — validating webhooks that catch broken
   scale targets and conflicting HPAs at apply time instead of letting
   them misbehave at runtime.

The external-metrics APIService is a cluster singleton — Kubernetes allows
exactly one external metrics provider — so one installation per cluster is
an upstream constraint and the Helm release name is FIXED to `keda`.

## Engine vs Declarations

This component installs and configures the ENGINE. The scaling
declarations — ScaledObject (scale a Deployment/StatefulSet), ScaledJob
(spawn Jobs per pending work item), TriggerAuthentication /
ClusterTriggerAuthentication (per-trigger credentials) — are KEDA custom
resources deployed per workload, next to the workloads they scale
(KubernetesManifest carries them today). The split matters for lifecycle:
the engine is cluster infrastructure with a fixed identity; the
declarations churn with the applications.

## CRD Lifecycle: the Keep Mechanism

The chart templates its CRDs (`crds.install`, default true), which means
Helm OWNS them — and a plain uninstall would cascade-delete the CRDs and
with them every ScaledObject/ScaledJob/TriggerAuthentication in the
cluster. The chart has no native keep knob, so the spec's
`crds.keep_on_uninstall` (default true) rides the standard Helm
`helm.sh/resource-policy: keep` annotation onto the CRDs through the
chart's `crds.additionalAnnotations` passthrough. The keep annotation only
makes sense when this release owns the CRDs, so it renders only when
`install && keep` — and a CEL rule refuses `keep_on_uninstall: true` with
`install: false` (nothing to keep).

## Component Sizing: Standbys, Not Capacity

`operator.replicas` and `metrics_server.replicas` (chart default 1) do not
scale horizontally: per upstream HA guidance, one operator instance leads
and one metrics-server instance serves at a time — extra replicas are warm
standbys that cut failover time. The metrics server matters most: the HPA
controller reads `external.metrics.k8s.io` through it on every reconcile
loop, so its downtime stalls every KEDA-driven HPA. Chart-default container
resources are requests 100m/100Mi, limits 1/1000Mi per component.

`watch_namespace` (empty by chart default) fences the operator to one
namespace instead of the normal cluster-wide watch — the single-team
posture.

## Admission Webhooks

`webhooks.enabled` (chart default true) validates ScaledObjects at apply
time. `failure_policy` is the sharpest knob: `Ignore` (chart default) lets
applies proceed unvalidated when the webhook is unreachable; `Fail`
rejects them until it is back — stricter, but a webhook outage then blocks
ScaledObject changes. Disabling the webhooks entirely moves every mistake
to runtime scaling failures.

## Pod Identity: the Environment-Injection Surface

`pod_identity` is how KEDA's scalers authenticate to cloud metric sources
without stored credentials: each arm annotates/labels KEDA's own service
accounts with the platform's keyless-identity contract. The arms configure
independent chart blocks and adapt one spec to the host cloud:

| Environment | Arm | Mechanism | Chart block |
|---|---|---|---|
| EKS / AWS | `aws_irsa` (role ARN, CEL-validated shape) | IAM Roles for Service Accounts via the cluster's OIDC provider — scalers read SQS/CloudWatch/Kinesis without keys | `podIdentity.aws.irsa` |
| AKS / Azure | `azure_workload_identity` (client + tenant ID, both required) | Federated Entra credentials — Service Bus/Event Hubs/Monitor without stored secrets | `podIdentity.azureWorkload` |
| GKE / GCP | `gcp_workload_identity` (service-account email, format-validated) | Workload Identity binding to a GCP IAM service account — Pub/Sub/Stackdriver without keys | `podIdentity.gcp` (`gcpIAMServiceAccount`) |

The cloud-side half of each contract (trust policy, federated credential,
WI binding) is written against the chart's FIXED operator service-account
name, `keda-operator` — which is why it is a stack output. Cross-cloud
cases (an EKS cluster scaling on Azure Service Bus) need no arm here: that
is a TriggerAuthentication with explicit credentials, next to the
workload. Ambient identity covers the ambient cloud; per-trigger
authentication covers everything else.

## Internal TLS

KEDA's components talk to each other over TLS, and the external-metrics
APIService needs a caBundle. `certificates.type` picks the provisioning:

- **`operator`** (chart default): the KEDA operator self-generates the
  certificates and patches the APIService caBundle — zero-dependency,
  fine for almost every cluster. The modules render nothing for this arm:
  the chart's own default needs no values.
- **`cert_manager`**: cert-manager issues and renews them (requires
  KubernetesCertManager). With no issuer reference the chart generates its
  own self-signed CA + Issuer chain; `cert_manager_issuer` instead points
  at an existing Issuer (namespaced, in the installation namespace) or
  ClusterIssuer — rendered as the chart's `issuer` block with
  `generate: false` and group `cert-manager.io`. A CEL rule refuses an
  issuer reference with any other certificates type.

## Timeouts, Scheduling, Telemetry

- **`http_timeout_ms`** (chart default 3000): the default timeout for
  scalers that reach external services over raw HTTP. Scalers built on
  cloud SDKs carry their own clients — the value does not necessarily
  apply to them (the upstream caveat).
- **`priority_class_name` / `node_selector` / `tolerations`** apply to all
  KEDA components. The engine belongs in the cluster-infrastructure
  priority tier: pods that scale on KEDA stop scaling without it.
- **`prometheus`** exposes KEDA's OWN `/metrics` (scaler loop latencies,
  trigger errors, HPA interactions) — telemetry ABOUT KEDA, unrelated to
  the external metrics it serves. The chart's layout is per-component, so
  one spec flag fans out identically to the `prometheus.operator` and
  `prometheus.metricServer` blocks; `service_monitor` requires the
  Prometheus operator CRDs, and the release FAILS to install without them.

## Typed Surface vs Escape Hatch

The typed spec covers the chart's meaningful configuration surface:
namespace and lifecycle, chart version, CRD lifecycle, watch scope,
per-component sizing, webhooks, pod identity, internal TLS, scaler HTTP
timeout, scheduling, and own telemetry.

`helm_values` merges LAST with Helm `-f` semantics on both engines
(Terraform natively via the two-document values list; Pulumi module-side
with the same deep-merge): maps deep-merge with the later document
winning, lists replace. Deliberately unmodeled as typed fields (all
reachable via `helm_values` where they are chart values at all):

- **The HTTP add-on** — HTTP request-based scaling is a SEPARATE upstream
  chart (`http-add-on`, with its own interceptor/scaler architecture),
  not a value of this one; modeling it here would misrepresent its
  lifecycle
- **Per-component scheduling overrides** (`operator.affinity`,
  `metricsServer.affinity`, `webhooks.affinity`,
  `topologySpreadConstraints.*`) — the typed scheduling fields apply to
  all components, which covers the real cases; per-component placement is
  an expert move
- **OpenTelemetry** (`opentelemetry.collector.uri`) — pushing KEDA
  telemetry to an OTel collector is an alternative telemetry pipeline
  beyond the Prometheus posture the platform standardizes on
- **Profiling** (`profiling.*`) — pprof endpoints for debugging KEDA
  itself; a development knob, not an operating posture
- **Network policies, PodDisruptionBudgets, upgrade strategies, image
  overrides** — the chart's long tail; each is a niche override with a
  correct chart default

## Install Semantics

Both engines install a REAL Helm release, atomically, with cleanup on fail
and a 300s timeout, waiting for the components to become Available. The
two classic install failures — a ServiceMonitor rendered without the
Prometheus operator CRDs, and broken internal TLS wiring — fail THIS
deploy with a readiness timeout instead of surfacing later as
ScaledObjects that mysteriously never scale. The module (not Helm) owns
namespace creation via `create_namespace`, so a namespace it creates
carries the standard governance labels and is deleted with the resource.

## Outputs

`namespace` (where the engine runs — ScaledObjects live next to their
workloads, not here), `release_name` (fixed `keda`), and
`operator_service_account_name` (the chart's fixed `keda-operator` — the
subject every cloud-side keyless binding is written against).

## E2E

The behavioral facts are properties of the platform, not of any one test
run:

- The aggregation-layer proof is the `v1beta1.external.metrics.k8s.io`
  APIService reporting Available — the exact dependency every
  ScaledObject-driven HPA has.
- A cron-trigger ScaledObject is the deterministic behavioral proof: the
  operator reconciles it into an HPA and drives the target Deployment
  above its baseline replica count with no external metric source in the
  loop.
- Scaling declarations belong AFTER the engine: a ScaledObject cannot be
  applied before the KEDA CRDs exist, so fixtures that precede the
  installation must not carry one.
- The ServiceMonitor arm fails the release on clusters without the
  Prometheus operator CRDs, by design.
- Uninstalling the release keeps the CRDs (and every scaling declaration)
  by default; the APIService itself is chart-owned and goes with the
  release — a dangling external-metrics registration would break every
  future HPA that consumes external metrics.
