# KubernetesCilium: Research and Design

## Introduction

Cilium is the eBPF-based networking, network-security, and observability
engine for Kubernetes: it is the cluster's CNI. The agent attaches eBPF
programs to every pod to wire networking, enforce NetworkPolicy (standard
Kubernetes policies plus Cilium's own L7-aware policies), optionally replace
kube-proxy entirely with eBPF service load-balancing, and stream flow-level
observability through Hubble. This component installs it from the official
Helm chart (`cilium` at `https://helm.cilium.io`; the pinned default chart
1.19.6 ships Cilium 1.19.6 — chart and app versions move together).

## Upstream Architecture

An installation is three cooperating pieces, all named FIXED by the chart
regardless of the release name:

1. **The agent** — DaemonSet `cilium`, one pod per node. It generates the
   node's CNI configuration, attaches the eBPF programs, and answers every
   per-packet decision (connectivity, policy, load-balancing, encryption).
2. **The operator** — Deployment `cilium-operator`. It handles the
   cluster-scoped work: IPAM allocation (carving per-node pod CIDRs in
   cluster-pool mode, managing ENIs in eni mode), CRD lifecycle, and —
   when enabled — the Gateway API control plane.
3. **Hubble** — the observability layer built into the agent, optionally
   extended with the `hubble-relay` Deployment (cluster-wide flow
   aggregation) and the `hubble-ui` Deployment (service-map console). The
   relay and UI Services are also fixed chart-template names
   (`hubble-relay`, `hubble-ui` — no release-derived prefix), which is what
   lets the stack outputs be pure functions of the spec toggles.

The agent DaemonSet, the operator, and the generated CNI configuration are
cluster singletons: one dataplane per cluster is an upstream constraint,
so the Helm release name is FIXED to `cilium` and never derives from
`metadata.name`. `kube-system` is the installation convention — the agent
is cluster infrastructure and several chart defaults assume it.

## The Load-Bearing Choice: Primary CNI vs Chaining

Every other decision in the spec follows from this one:

- **Primary CNI** (default): Cilium owns pod networking. The cluster must
  be created WITHOUT another CNI (kind: `disableDefaultCNI: true`; AKS:
  `--network-plugin none`; kubeadm: no CNI addon). Nodes sit NotReady until
  Cilium installs — by design: the kubelet flips a node Ready only once a
  CNI configuration exists, so the install itself is what unblocks
  scheduling.
- **CNI chaining** (`cni.chaining_mode`): Cilium runs ON TOP of an existing
  CNI. The incumbent keeps IPAM and interface wiring; Cilium attaches eBPF
  programs for policy enforcement, load-balancing, and observability. The
  upstream chaining vocabulary is closed: `none`, `aws-cni` (EKS with the
  AWS VPC CNI), `flannel`, `generic-veth` (any veth-based CNI), `portmap`.
  `cni.chaining_target` names a CNI network to chain into (implying
  generic-veth); `aws-cni` chaining implies its own target.

Chaining has one hard corollary the spec CEL-enforces: `cni.exclusive` must
be false. Exclusive mode (chart default true) makes Cilium take ownership
of `/etc/cni/net.d` and rename non-Cilium CNI configurations — which would
destroy the very configuration the chain depends on.

## Per-Environment Posture

The `cloud` block is exactly the environment-injection surface: the same
spec deploys to any cluster, and these arms adapt the datapath to the host
cloud. At most one may be enabled (CEL-enforced) — they configure mutually
exclusive datapaths.

| Environment | Posture | Key fields |
|---|---|---|
| kind / local dev | Primary CNI (cluster created with `disableDefaultCNI: true`) | `ipam.mode: kubernetes` (kind allocates PodCIDRs), `operator.replicas: 1` on single-node |
| Self-managed (kubeadm, Cluster API, bare metal) | Primary CNI, optionally kube-proxy-free | `ipam.mode: cluster-pool` with a deliberate CIDR plan; `kube_proxy_replacement` + `k8s_service_host`/`k8s_service_port` when kube-proxy is skipped |
| EKS — keep the AWS VPC CNI | Chaining | `cni.chaining_mode: aws-cni`, `cni.exclusive: false`; `cloud` stays empty (the VPC CNI IS the cloud integration) |
| EKS / EC2 — Cilium as primary CNI | ENI datapath | `cloud.aws_eni` + `ipam.mode: eni` — pods draw VPC-routable IPs from ENIs Cilium manages; the agent uses the node instance role (or IRSA) for EC2 ENI calls |
| AKS | BYOCNI | `cloud.aks_byocni` on clusters created with `--network-plugin none` — the supported primary-CNI path on AKS; the legacy azure-IPAM integration (service-principal credentials) is deliberately not typed |
| GKE | GKE integration | `cloud.gke` — configures the agent for GKE node images and routing; pairs with `ipam.mode: kubernetes` and native routing per upstream's GKE guide |

## IPAM and Routing

`ipam.mode` mirrors the agent's own accepted vocabulary: `cluster-pool`
(chart default — the operator carves per-node pod CIDRs from
`cluster_pool_ipv4_pod_cidrs`, chart default `10.0.0.0/8` with a /24 per
node), `kubernetes` (consume the node's Kubernetes-assigned PodCIDR — the
right mode on kind and kubeadm, whose control planes already allocate),
`eni`, `azure`, `alibabacloud`, `multi-pool`, and `delegated-plugin` (IPAM
stays with the chained CNI). Cluster-pool CIDRs must overlap neither the
node network nor the Service CIDR.

`routing.mode` is `tunnel` (chart default — encapsulate inter-node pod
traffic with `vxlan` (default) or `geneve`; works on any network) or
`native` (route pod CIDRs directly on the underlying fabric — lower
overhead, requires the network to carry pod CIDRs). The native-routing
knobs are CEL-fenced to native mode: `ipv4_native_routing_cidr` (the range
that is not masqueraded because the fabric can route it back) and
`auto_direct_node_routes` (each agent installs direct routes to peer pod
CIDRs — only correct on a shared L2 segment).

## Kube-Proxy Replacement

`kube_proxy_replacement` makes Cilium serve ClusterIP/NodePort/LoadBalancer
traffic in eBPF instead of iptables chains. The cluster must then be
created WITHOUT kube-proxy (kind: `kubeProxyMode: none`; kubeadm:
`--skip-phases=addon/kube-proxy`). One bootstrap subtlety makes
`k8s_service_host` mandatory (CEL-enforced): before Cilium's own service
load-balancing is up there is no kube-proxy to resolve the in-cluster
`kubernetes.default` Service, so the agent needs the API server's real
address. On kind that is the control-plane node's container name; on
managed clouds the API endpoint hostname; 6443 is the common port.

A rendering trap lives here: `kubeProxyReplacement` is a STRING in the
chart's values (`"false"` is the declared default — historically it took
`"strict"`/`"partial"`), and `k8sServicePort` is also a string (default
`""`). Both modules render them as strings to keep the values document
byte-identical with what the chart declares.

## Gateway API

`gateway_api` enables Cilium's Gateway API implementation: north-south
routing with Envoy embedded in the agent, registering the fixed `cilium`
GatewayClass. Two dependencies are real:

- **kube-proxy replacement is REQUIRED.** The operator checks it as a
  precondition and disables Gateway API support with only a log warning
  when it is off — the gateway would silently never program. The spec makes
  the dependency loud with a CEL rule instead.
- **The Gateway API CRDs must already be on the cluster**
  (KubernetesGatewayApiCrds) — Cilium implements the API but does not ship
  its CRDs.

## Hubble

Hubble is enabled in the agent by chart default (the spec only renders an
explicit false); the deployable extras are opt-in. `relay` deploys the
cluster-wide flow aggregation service; `ui` deploys the service-map console
and reads flows exclusively through the relay (CEL: ui requires relay).
`metrics` is upstream's LIST of metric families (`dns`, `drop`, `tcp`,
`flow`, `icmp`, `http`, with option syntax like `dns:query;ignoreAAAA`) —
the chart key is `hubble.metrics.enabled`, a list despite the name, and an
empty list means Hubble metrics stay disabled. `metrics_service_monitor`
requires a non-empty metrics list (nothing to scrape otherwise) and the
Prometheus operator CRDs on the cluster — the release fails to install
without them.

## Encryption

`encryption` turns on transparent pod-to-pod encryption: `wireguard`
(kernel WireGuard, automatic key management — the low-friction choice and
the spec's default TYPE when the block is enabled) or `ipsec` (requires a
pre-created key Secret, supports key-rotation policies). `node_encryption`
extends encryption to pure node-to-node traffic and is WireGuard-only (the
upstream constraint, CEL-enforced). Upstream also accepts a third type,
`ztunnel`, for Istio ambient interop — deliberately not typed; reach it via
`helm_values` alongside an ambient mesh.

An asymmetry in the CEL rules is deliberate: an explicit `type` with
`enabled: false` is tolerated, because `type` carries a default annotation
and the platform's defaulting middleware can materialize it in a block a
user wrote as `{enabled: false}`. The modules render the encryption block
only when enabled, so a defaulted type with enabled false renders nothing.

## Policy Enforcement, Bandwidth, Sizing

- **`policy_enforcement_mode`** (chart default `default`): enforce only
  where policies select pods; `always` is default-deny; `never` disables
  enforcement (observability only). Closed upstream enum.
- **`bandwidth_manager`**: enforce pods' `kubernetes.io/egress-bandwidth`
  annotations in eBPF instead of token-bucket qdiscs; `bbr` additionally
  enables BBR congestion control for pods (requires a 5.18+ kernel per
  upstream guidance) and is CEL-gated on the manager being enabled.
- **`operator.replicas`** (chart default 2): HA via leader election, spread
  by pod anti-affinity. On single-node clusters run 1 — the second replica
  can never schedule and the rollout never settles, which matters because
  the install waits for it.
- **`agent_resources` / `operator.resources`**: the chart sets no
  requests/limits by default — production installs should size
  deliberately. The top-level chart `resources` key is the agent container.

## Cilium's Own Telemetry

`prometheus.enabled` exposes the agent's and the operator's `/metrics`
endpoints; one spec toggle drives BOTH chart blocks (top-level `prometheus`
for the agent, `operator.prometheus` for the operator) — exposing only one
of the two would be a confusing half-telemetry posture. `service_monitor`
adds ServiceMonitors for both and requires the Prometheus operator CRDs —
the release fails to install without them. Hubble metrics are a separate
surface (`hubble.metrics`).

## Typed Surface vs Escape Hatch

The typed spec covers the chart's meaningful configuration surface: cluster
identity, IPAM, routing, kube-proxy replacement, CNI chaining, the three
cloud datapath arms, Hubble, encryption, policy enforcement, Gateway API,
bandwidth manager, operator sizing, agent resources, and own telemetry.

`helm_values` merges LAST with Helm `-f` semantics on both engines
(Terraform natively via the two-document values list; Pulumi module-side
with the same deep-merge): maps deep-merge with the later document winning,
lists replace. It is the escape hatch for the chart's long tail — never the
substitute for the typed fields.

Deliberately unmodeled as typed fields (all reachable via `helm_values`):

- **BGP control plane** (`bgpControlPlane.*`) — advertising pod CIDRs and
  LoadBalancer IPs to physical routers is a network-engineering
  integration whose real configuration lives in CiliumBGP* custom
  resources, not chart values; a typed toggle would model the least
  interesting part of it
- **L2 announcements** (`l2announcements.*`) — the bare-metal
  LoadBalancer ARP/NDP story; niche, and coupled to CiliumL2AnnouncementPolicy
  resources beyond the chart
- **Cluster Mesh** (`clustermesh.*`) — multi-cluster connectivity is an
  estate-level project (shared CA, per-cluster IDs, apiserver exposure)
  that deserves its own deliberate design, not a spec arm; `cluster_name`
  is typed precisely so an installation stays mesh-ready
- **Cilium's ingress controller** (`ingressController.*`) — the Gateway
  API arm is the typed north-south story; carrying both would duplicate
  one role with two APIs
- **Per-image overrides** (`image.*`, `operator.image.*`, ...) — the
  air-gapped/mirror knob; prefer bumping `chart_version` over pinning
  image tags, which decouples the binary from the chart's tested pair
- **Legacy azure-IPAM credentials** (`azure.*` with service-principal
  secrets) — superseded by BYOCNI (`cloud.aks_byocni`) for new clusters,
  and putting cloud credentials in chart values is the wrong shape

## Install Semantics

Both engines install a REAL Helm release, atomically, with cleanup on fail
— and a 600-second timeout instead of the usual 300, because the install
path is heavier than an ordinary workload chart: the agent DaemonSet must
roll out on EVERY node plus the operator, and on a fresh cluster nodes
transition NotReady→Ready only as Cilium wires each one — the rollout
itself unblocks scheduling. A dataplane that never converges fails THIS
deploy instead of surfacing later as pods stuck in ContainerCreating. The
module (not Helm) owns namespace creation via `create_namespace`, which
stays false for kube-system installs.

## Outputs

`namespace`, `release_name` (fixed `cilium`), `cluster_name` (the resolved
identity — the name this cluster carries in Hubble flows and any future
Cluster Mesh), `hubble_relay_service_name` and `hubble_ui_service_name`
(the fixed chart-template Service names, empty when the component is not
deployed — the outputs contract mirrors what actually exists), and
`gateway_class_name` (`cilium` when `gateway_api` is on, empty otherwise).

## E2E

The behavioral facts are properties of the platform, not of any one test
run:

- The NotReady→Ready node transition is the proof that Cilium became the
  cluster's CNI — on a cluster created without a default CNI, nodes cannot
  go Ready until the agent generates a CNI configuration on each of them.
- NetworkPolicy enforcement requires an enforcing CNI: on a cluster whose
  CNI ignores NetworkPolicy (kind's default kindnet, for one), policies
  apply cleanly and enforce nothing. Cilium is what makes deny-policy
  behavior provable.
- On single-node clusters, `operator.replicas: 1` is mandatory for the
  install wait to converge — the chart-default 2 replicas carry pod
  anti-affinity.
- WireGuard encryption depends on the node kernel's module set — a host
  property, not a component property.
- ServiceMonitor toggles (Cilium's own and Hubble's) fail the release on
  clusters without the Prometheus operator CRDs, by design.
- Cloud datapath arms and CNI chaining require the actual cloud CNI or
  fabric to exist — they are proven on real clusters, not local ones.
