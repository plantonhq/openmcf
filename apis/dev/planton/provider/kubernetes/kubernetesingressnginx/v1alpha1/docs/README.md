# KubernetesIngressNginx: Research and Design

## Introduction

ingress-nginx is the Kubernetes-community ingress controller: a Deployment
(or DaemonSet) running an embedded NGINX whose configuration is continuously
regenerated from the cluster's Ingress resources. It is the de facto
standard HTTP(S) entry point for Kubernetes — traffic arrives at a Service
in front of the controller pods, NGINX terminates and routes it, and the
Ingress objects that drive the routing are ordinary namespaced resources
that application teams own. This component installs that controller from
the official Helm chart (`ingress-nginx` at
`https://kubernetes.github.io/ingress-nginx`, default chart 4.15.1 =
controller v1.15.1).

## Upstream Architecture

The controller has three moving parts:

**The watch loop.** The controller watches Ingress resources whose class
matches the IngressClass it owns (`--controller-class`), plus the Services,
Endpoints, and Secrets they reference, and renders them into a single NGINX
configuration. Config changes that only touch upstream endpoints are
applied without a reload; structural changes trigger an NGINX reload.

**The entry Service.** Traffic reaches the controller pods through a
Kubernetes Service — LoadBalancer on managed clouds (the cloud provisions
the actual LB), NodePort or host networking on bare metal. The controller
itself never calls cloud APIs: everything the cloud does for the LB
(scheme, type, target mode, static addresses) is expressed as annotations
on this Service and executed by the cloud's own service controller.

**The admission webhook.** A ValidatingWebhookConfiguration that renders a
candidate NGINX config for every Ingress write and rejects ones that would
break the controller — moving failures to kubectl-apply time instead of
taking down the shared entry point. The webhook's certificate is
bootstrapped by the chart's kube-webhook-certgen hook Jobs.

Leader election coordinates status updates among replicas; the election ID
defaults to `<fullname>-leader` in the chart, which is why pinning the
chart fullname per instance also isolates election per instance.

## The Chart's Value Surface: Typed vs. helm_values

The chart's `values.yaml` is large; the typed spec covers the
configuration that changes deployment behavior in practice:

- **IngressClass identity** (`controller.ingressClassResource.*`,
  `controller.ingressClass`, `controller.watchIngressWithoutClass`) — the
  spec's `ingress_class` block. The modules also keep the legacy
  annotation-vocabulary value (`controller.ingressClass`) in lockstep with
  the class name; the chart defaults it to `nginx` independently, which
  would mis-route annotation-based Ingresses on a non-default-named
  instance
- **Capacity** (`controller.replicaCount`, `controller.autoscaling.*`,
  `controller.resources`) — `replicas`, `autoscaling`, `resources`. When
  autoscaling is enabled the modules do NOT render `replicaCount`: the
  chart's Deployment template omits replicas under autoscaling, and
  rendering both would create a rollout tug-of-war. Chart HPA defaults:
  min 1, max 11, CPU and memory targets 50%. The chart auto-creates a
  PodDisruptionBudget (minAvailable 1) when effective replicas exceed one
  (Deployment only — Kubernetes does not support PDBs for DaemonSets);
  there is deliberately no separate PDB toggle in the spec
- **The controller Service** (`controller.service.*`) — type, annotations,
  externalTrafficPolicy, loadBalancerSourceRanges, loadBalancerClass,
  enableHttp/enableHttps, nodePorts, and the chart's optional second
  internal Service (`controller.service.internal`), which the chart
  refuses to render without at least one annotation — enforced in the spec
  as a validation rule
- **Workload shape** (`controller.kind`, `controller.hostNetwork`,
  `controller.hostPort`) — `controller_kind`, `host_network`, `host_ports`.
  With host networking the modules also set
  `dnsPolicy: ClusterFirstWithHostNet` so in-cluster name resolution keeps
  working. `host_network` and `host_ports` are mutually exclusive in the
  spec — host networking already binds every listener on the node
- **NGINX behavior** (`controller.config`,
  `controller.allowSnippetAnnotations`) — `nginx_config` passes upstream's
  ConfigMap key/value vocabulary through verbatim (see "Deliberately
  upstream-shaped" below); `allow_snippet_annotations` follows upstream's
  post-CVE-2021-25742 default of false
- **Default TLS certificate** — the chart exposes no first-class key for
  it; upstream's own documented mechanism is the
  `--default-ssl-certificate` controller flag via `controller.extraArgs`,
  which is exactly how both modules render `default_tls_certificate`
  (namespace defaults to the installation namespace)
- **Default backend** (`defaultBackend.*`) — enabled, replicas, a custom
  image (spec carries `repository:tag`; the modules split it for the
  chart), resources
- **Admission webhook** (`controller.admissionWebhooks.*`) — enabled
  (chart default true), failurePolicy (chart default `Fail`),
  timeoutSeconds
- **Metrics** (`controller.metrics.*`) — the metrics endpoint (port 10254)
  and the ServiceMonitor with scrape interval and additional labels. The
  ServiceMonitor requires the Prometheus operator CRDs on the cluster or
  the release fails to install — the spec enforces that `service_monitor`
  implies `enabled`
- **TCP/UDP exposure** (`tcp`, `udp` root values) — maps of external port
  to `"namespace/service:port"`, upstream's exposing-tcp-udp-services
  mechanism; the chart publishes each port on the controller Service
- **Scheduling** (`controller.nodeSelector`, `controller.tolerations`,
  `controller.priorityClassName`) — chart default nodeSelector is
  `kubernetes.io/os: linux`
- **Image registry** (`global.image.registry`) — the air-gapped/mirror
  knob for all chart images (controller, kube-webhook-certgen,
  defaultbackend); empty means `registry.k8s.io`

`helm_values` merges LAST with Helm `-f` semantics on both engines
(Terraform natively via the two-document values list; Pulumi module-side
with the same deep-merge: nested maps merge with the later document
winning, scalars and lists replace).

### Deliberately upstream-shaped: nginx_config

`nginx_config` is a raw string map on purpose. Upstream's ConfigMap
vocabulary is roughly two hundred tuning keys (`proxy-body-size`,
`worker-shutdown-timeout`, `use-forwarded-headers`, `ssl-protocols`, ...)
documented and versioned by the controller itself. Inventing a typed
mirror of that vocabulary would drift from the controller's own
documentation with every controller release; passing keys and values
through verbatim keeps upstream's docs authoritative.

### Deliberate exclusions (reachable via helm_values)

These chart capabilities are intentionally NOT typed fields — all are
expressible through `helm_values`, and all were excluded for low demand:

- **KEDA-based autoscaling** (`controller.keda.*`) — the chart's
  alternative autoscaling arm (mutually exclusive with the HPA). Requires
  a KEDA installation and scaler-specific trigger configuration; the HPA
  arm covers the common case
- **Custom NGINX template** (`controller.customTemplate`) — replacing the
  controller's config template wholesale is an expert maneuver with a
  per-version compatibility burden
- **MaxMind/GeoIP** (`controller.maxmindLicenseKey`) — gates GeoLite2
  database downloads behind a third-party license key; a licensing knob,
  not deployment configuration
- **Namespace scoping** (`controller.scope.*`) — limiting the controller
  to one namespace (or a namespace label selector) is a niche multi-tenant
  posture; the multi-instance pattern with distinct classes covers the
  mainstream separation need
- **PrometheusRule** (`controller.metrics.prometheusRule`) — alert rules
  are operator-specific policy; the ServiceMonitor (scrape discovery) is
  the integration seam this spec owns

## The Multi-Instance Pattern

Running several controllers per cluster — most commonly a public + internal
split — is the upstream-sanctioned pattern, and this component makes it the
default shape rather than an advanced configuration:

- **Release name = `metadata.name`**, never a fixed chart name: each
  manifest is its own release
- **Chart fullname pinned to the release name** (`fullnameOverride`): every
  chart object (Deployment/DaemonSet, Services, RBAC, webhook) carries a
  deterministic, manifest-derived name — `<name>-controller`,
  `<name>-controller-internal` — which verification, imports, and
  downstream composition key off
- **One IngressClass per instance** (`ingress_class.name`, default
  `nginx`), unique per cluster. The class's controller identifier derives
  automatically: the chart default `k8s.io/ingress-nginx` for class
  `nginx`, otherwise `k8s.io/<class-name>` — additional controllers
  isolate without the user inventing a vocabulary
- **Leader election isolates for free**: the chart's election ID defaults
  to `<fullname>-leader`, and the fullname is per-instance
- **At most one `is_default_class`** and at most one
  `watch_ingress_without_class` instance per cluster — both are
  cluster-wide claims

## Environment Injection: Per-Cloud LB Recipes

There is no provider oneof and no workload-identity block in this spec BY
DESIGN: the controller never calls cloud APIs. The host cloud shapes the
load balancer entirely through `service.annotations` — this map is the
component's environment-injection surface, playing the role that
credential/identity blocks play in components that do call cloud APIs.

| Cloud / posture | Annotations |
|---|---|
| AWS NLB (production default on EKS) | `service.beta.kubernetes.io/aws-load-balancer-type: "external"`, `service.beta.kubernetes.io/aws-load-balancer-nlb-target-type: "ip"` |
| AWS internal | `service.beta.kubernetes.io/aws-load-balancer-scheme: "internal"` |
| GCP internal | `networking.gke.io/load-balancer-type: "Internal"` |
| GCP static IP | `networking.gke.io/load-balancer-ip-addresses: "<address-name>"` (or `controller.service.loadBalancerIP` via `helm_values`) |
| Azure internal | `service.beta.kubernetes.io/azure-load-balancer-internal: "true"` |

`load_balancer_class` selects a non-default LB implementation where the
cloud offers several (e.g. `service.k8s.aws/nlb` for the AWS Load Balancer
Controller); `external_traffic_policy: local` preserves client source IPs
and skips the extra node hop — the usual production choice, with the
trade-off that LB health checks then only pass on nodes running a
controller pod.

## Bare-Metal Postures

Without a cloud LB controller there is nothing to provision a
`load_balancer` Service, so the component offers the standard bare-metal
shapes:

- **DaemonSet + hostPorts** (`controller_kind: daemon_set`, `host_ports`):
  one controller per node, listeners on node ports 80/443 via hostPort —
  CNI-friendly, no host networking
- **Host network** (`host_network`): the controller binds the node's
  interfaces directly; the modules set `dnsPolicy:
  ClusterFirstWithHostNet` so cluster DNS keeps working. Mutually
  exclusive with `host_ports`
- **NodePort service** (`service.type: node_port`, optionally with pinned
  `http_node_port`/`https_node_port`): the pattern behind an external LB
  you manage yourself

On such clusters a `load_balancer` service type fails loudly at install:
Helm's readiness wait includes the LB address, and an address that never
arrives times out (300s) with atomic rollback. Both modules document this
as deliberate — the failure names the real problem instead of leaving a
silently Pending entry point.

## Install Semantics

Both engines install a REAL Helm release with wait + atomic +
cleanup-on-fail (300s timeout), and the Terraform module additionally waits
for hook Jobs (`wait_for_jobs`) to cover the admission-webhook certgen
hooks. A controller that never starts — bad image, unschedulable pod,
certgen failure, missing ServiceMonitor CRDs — fails THIS deploy, not the
first Ingress. The modules create only the optional anchor namespace
themselves; the chart owns every controller object. Governance labels are
stamped on the module-created namespace only, never injected into chart
resources.

After the release, both modules read the controller Service back for the
LB address outputs — gated on the `load_balancer` service type (for
`node_port`/`cluster_ip` there is no LB status to read and the address
outputs stay empty by design). The release wait guarantees the address
exists by the time the read runs.

## Outputs as Composition Seams

`ingress_class_name` — what KubernetesIngress resources reference to route
through this controller. `controller_service_name` (`<name>-controller`) —
the traffic entry point; external-dns and manual DNS records point at this
Service's address. `load_balancer_ip` / `load_balancer_hostname` — the
cloud-assigned address (AWS populates a hostname; GCP/Azure populate an
IP; both empty without a cloud LB controller). `internal_service_name` —
populated only when the dual-LB internal Service is enabled.
`release_name` and `namespace` — the handles for verification and
Helm-level operations.

## E2E

The kind-cluster lanes prove the install machinery with `node_port`
services (kind has no cloud LB controller — the LoadBalancer-wait failure
mode is the documented posture, not a lane): a chart-default minimal
install asserting the controller Deployment is Available and the `nginx`
IngressClass exists; a multi-instance coexistence lane deploying a second
controller with its own class while the first is live; and a tuned lane
exercising the typed breadth — custom class, chart-managed HPA, NGINX
ConfigMap tuning, admission-webhook tuning, default backend, metrics
endpoint, TCP exposure, and an escape-hatch value. ServiceMonitor stays off
in E2E (no Prometheus operator CRDs on the kind lane). Real cloud LB
provisioning rides the real-cluster lanes.
