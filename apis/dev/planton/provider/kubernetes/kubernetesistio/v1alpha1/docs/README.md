# KubernetesIstio: Research and Design

## Introduction

Istio is the most widely deployed Kubernetes service mesh: a control plane
(istiod) that validates mesh configuration, issues workload certificates, and
programs a data plane of Envoy-based proxies. This component installs that
control plane from the official Helm charts. Gateways and mesh traffic policy
are separate first-class kinds — this component deliberately manages ONLY the
control plane.

## Design Authority

Designed from the pinned upstream charts' `values.yaml` files
(`base`, `istiod`, `cni`, `ztunnel` at
`https://istio-release.storage.googleapis.com/charts`) and validated against
their defaults. Istio versions its charts in lockstep with the product, so one
version pin (default `1.30.3`) drives every release this component installs.
Chart identity is fixed identically in both engines' modules — cross-engine
chart drift would deploy two different products from one manifest.

## The Deployment Landscape

### Sidecar vs Ambient

The data plane mode is the single largest architectural decision, so the spec
makes it first-class (`dataplane_mode`) rather than burying it in chart values:

- **Sidecar** (the default, the classic architecture): a proxy container is
  injected into every workload pod. Full L7 capability per workload, the
  richest policy surface, and per-pod resource cost that multiplies across the
  mesh (`proxy.resources` defaults: 100m/128Mi request, 2/1Gi limit — size
  deliberately on large meshes). Namespaces opt in via the
  `istio-injection=enabled` label; `sidecar_injector` and `proxy` tune the
  injection behavior and proxy defaults.
- **Ambient** (sidecar-less): a per-node **ztunnel** DaemonSet carries
  mTLS and L4 policy for every enrolled pod, and optional **waypoint** proxies
  add L7 where a namespace or service actually needs it. Two extra releases
  install (`istio-cni` and `ztunnel`), rendered through the charts' own
  `profile: ambient` overlay — the same overlay the upstream install path
  uses. Enrollment is per namespace or pod via the
  `istio.io/dataplane-mode=ambient` label. Lower steady-state cost, no pod
  restarts to join the mesh, but a younger L7 story than sidecars.

Switching modes later is a workload-by-workload migration (relabeling,
draining, re-verifying policy), not a flag flip — the spec documents the
choice as install-time for that reason. Validation rules enforce coherence:
`ztunnel` settings require ambient mode, `sidecar_injector` settings are
rejected in ambient mode.

### The CNI node agent

`istio-cni` is dual-natured, and the spec models that directly. In ambient
mode it ALWAYS installs (it is how traffic reaches ztunnel — not optional).
In sidecar mode it is opt-in (`cni.enabled`): it replaces the injected
privileged `istio-init` container with a node-level agent — required on
platforms that forbid NET_ADMIN init containers (e.g. OpenShift), recommended
for tighter pod security everywhere. When enabled in sidecar mode, the istiod
release is also told (`cni.enabled` chart value) so the injector emits pods
that rely on the agent.

### Revisions and the upgrade discipline

Istio supports **sequential, single-minor upgrades only**: an existing mesh
must step through each minor, never skip. The `version` field pins all four
charts at once; upgrading is a redeploy at the next minor.

`revision` names the control plane: the istiod release and Service become
`istiod-<revision>`, and injection scopes to workloads labeled
`istio.io/rev: <revision>`. Empty runs the unnamed DEFAULT revision — the
standard posture. The revision is a naming discipline for organizations that
always run revisioned control planes, NOT a canary mechanism here:
**side-by-side canary control planes are deliberately not modeled**, because
one KubernetesIstio per cluster is the constraint that keeps the design sound
(the CRDs and the validation-webhook plumbing are cluster singletons, and the
`base` release's `defaultRevision` points the default validating webhook at
exactly one control plane).

### CRD ownership (the KubernetesIstioBaseCrds seam)

The Istio CRDs (DestinationRule, AuthorizationPolicy, Telemetry, ...) are
applied by the module itself via **server-side apply, outside the Helm
release** — the `base` release installs with `base.excludedCRDs` covering the
entire pinned bundle, so it templates NO CRDs. The reason is adoption: Helm
refuses to manage CRDs that already exist without its ownership metadata. If
the chart owned the CRDs, a cluster running the CRDs-only
KubernetesIstioBaseCrds kind could never upgrade to the full mesh. Server-side
applied CRDs are co-ownable by both kinds, so that migration is a plain
redeploy — deploy KubernetesIstio, the CRDs reconcile in place, the control
plane arrives.

The CRD bundle version is pinned to `spec.version` (the upstream
`crd-all.gen.yaml` at the release tag), so the installed CRD schema always
matches the control plane and the typed Istio kinds' generated SDK.

Destroying this component removes everything **including the CRDs** (standard
engine semantics for module-owned resources) — mesh configuration objects
cascade with them.

## Chart Anatomy (what the typed fields map to)

| Spec surface | Chart | Values |
|---|---|---|
| `revision` | base, istiod, cni, ztunnel | `defaultRevision` (base), `revision` (others) |
| `dataplane_mode: ambient` | istiod, cni | `profile: ambient` (the charts' own overlay) |
| `istiod.replicas` / `istiod.autoscale` | istiod | `replicaCount` + `autoscaleEnabled=false`, or `autoscaleEnabled`/`autoscaleMin`/`autoscaleMax`/`cpu.targetAverageUtilization` |
| `istiod.resources` / `log_level` / `pod_disruption_budget` / `priority_class_name` / scheduling | istiod | `resources`, `global.logging.level`, `global.defaultPodDisruptionBudget.enabled`, `global.priorityClassName`, `nodeSelector`, `tolerations` |
| `mesh_config.*` | istiod | `meshConfig.trustDomain`, `meshConfig.outboundTrafficPolicy.mode`, `meshConfig.accessLogFile`, `meshConfig.enablePrometheusMerge`, `global.multiCluster.clusterName`, `global.network`, `global.meshID` |
| `proxy.*` | istiod | `global.proxy.resources` / `logLevel` / `autoInject` / `clusterDomain` |
| `sidecar_injector.*` | istiod | `sidecarInjectorWebhook.enableNamespacesByDefault` / `rewriteAppHTTPProbe` |
| `cni.*` | cni (+ istiod's `cni.enabled` in sidecar mode) | `excludeNamespaces`, `cniBinDir`, `cniConfDir`, `chained` |
| `ztunnel.*` | ztunnel | `resources`, `logLevel` |
| `gateway_defaults.service_type` | istiod | `gatewayClasses.istio.service.spec.type` — the per-GatewayClass deployment overlay for Gateway API auto-provisioning |
| `images.*` | all four | `global.hub` / `global.variant` / `global.imagePullSecrets` (base/istiod/cni); top-level `hub` / `variant` / `imagePullSecrets` (ztunnel reads them there) |
| non-default `namespace` | base, istiod, ztunnel | `global.istioNamespace` / `istioNamespace` |

`helm_values` merges LAST per release with Helm `-f` semantics on both engines
(Terraform natively via the two-document values list; Pulumi module-side with
the same deep-merge). It covers chart surface beyond the typed fields — pilot
environment variables, the rest of MeshConfig, per-component overrides — and
is never the substitute for them.

## Install Semantics

Both engines install REAL Helm releases in upstream's own order — CRDs, then
`istio-base`, then `istiod` (or `istiod-<revision>`), then `istio-cni` and
`ztunnel` when applicable — atomic, cleanup-on-fail, waiting for readiness.
Waiting on istiod is the whole promise: a control plane whose webhooks and
discovery service are not serving rejects every mesh-config apply and every
injection, so a premature "success" would just move the failure downstream.

## Deliberately Unmodeled Surfaces

- **Multicluster east-west topology.** The identity fields
  (`mesh_config.cluster_name`, `network`, `mesh_id`) are modeled because they
  must be set at install time for a cluster that will ever join a
  multi-network mesh. The rest — east-west gateways, remote secrets, endpoint
  discovery federation — is an operational topology, not a property of one
  control-plane install.
- **External / remote istiod.** Every install here runs its own control
  plane; pointing a data plane at a control plane in another cluster is a
  different product shape.
- **VirtualService-based routing.** The Kubernetes Gateway API is the
  catalog's routing surface: KubernetesGateway (with the `istio`
  GatewayClass) and the route kinds express north-south routing;
  VirtualService remains reachable for mesh-internal routing via
  KubernetesManifest where genuinely needed.
- **Sidecar, WasmPlugin, WorkloadEntry, ProxyConfig CRs.** The CRDs install
  with the bundle, so the objects are applyable via KubernetesManifest on
  demand — but they are expert, niche surfaces that do not warrant typed
  kinds.

## Outputs as Composition Seams

`gateway_class_name` (`istio`) — what a KubernetesGateway selects to make
istiod provision its deployment (named `<gateway>-istio`).
`trust_domain` — the prefix of SPIFFE principals
(`spiffe://<trust_domain>/ns/<ns>/sa/<sa>`) in KubernetesAuthorizationPolicy
rules. `revision` and `istiod_service_name` — the control-plane identity for
`istio.io/rev` selection and discovery-address wiring. `dataplane_mode` —
tells composed resources whether workloads enroll via injection labels or the
ambient dataplane-mode label.

## E2E

Dual-engine install proofs run for BOTH data plane modes: sidecar
(base + istiod) and ambient (base + istiod + cni + ztunnel). Beyond install
verification, the mesh is proven live: a Gateway API Gateway provisioned by
istiod from the `istio` GatewayClass reaches Accepted/Programmed and routes a
real request end-to-end, and an AuthorizationPolicy is enforced live — a
meshed client's request is 403-denied while the policy exists and succeeds
after its removal. Blind state-import round-trips cover the whole module,
CRDs included.
