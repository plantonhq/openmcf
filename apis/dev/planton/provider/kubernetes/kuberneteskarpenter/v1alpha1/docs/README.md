# KubernetesKarpenter: Research and Design

## Introduction

Karpenter replaces the node-group model of cluster scaling with direct
provisioning: the controller watches for unschedulable pods, computes the
cheapest machine shapes that fit them from the live instance catalog,
launches those machines in seconds, consolidates under-used nodes away,
and drains ahead of spot interruptions. This component installs the
engine from the official OCI-served Helm charts
(`oci://public.ecr.aws/karpenter`): the `karpenter` controller chart plus
the companion `karpenter-crd` chart, both pinned to `1.14.0` by default —
the two charts version together with the controller, so one version pins
both releases.

A packaging caveat worth knowing: the chart source vendored in the
controller repository (`charts/karpenter`) lags the OCI-served artifact —
at the `1.14.0` served charts the vendored copy still reads `1.13.0`.
The served charts are the authority for values and defaults; claims below
are verified against the served `1.14.0` values.yaml and templates.

## Upstream Architecture

Karpenter is per-cloud upstream: each cloud ships its own controller,
chart, and NodeClass CRD. This kind installs the AWS provider
(`aws/karpenter-provider-aws`) — the canonical, mature implementation —
and carries the cloud-specific settings in a typed `cloud` oneof with
`aws` as the sole arm today, so future providers land as additive
siblings without redesign.

One installation is:

1. **The `karpenter-crd` release** — the CRDs (NodePool, NodeClaim,
   EC2NodeClass, NodeOverlay, CapacityBuffer; the files under the
   controller source's `pkg/apis/crds/`) as their own Helm release.
   This is upstream's supported mechanism for keeping CRDs upgradable:
   Helm installs the copies bundled inside a main chart once and NEVER
   upgrades them, so a dedicated CRD release is what keeps them current
   across chart upgrades.
2. **The `karpenter` controller release** — the controller Deployment
   (chart default 2 replicas: an active leader plus a warm standby,
   spread across zones by the chart's topology constraints), installed
   with its bundled CRD copies skipped unconditionally so the CRD
   release stays the single owner.

Both release names are FIXED. Karpenter owns the cluster-wide
`karpenter.sh` label domain, its CRDs, and node lifecycle — one
installation per cluster is an upstream constraint, so the names never
derive from `metadata.name`. The fixed release name also pins the
service-account name: with `serviceAccount.name` unset the chart derives
it from the fullname template, and because the release name
(`karpenter`) contains the chart name, the fullname — and therefore the
service account — is `karpenter`. That determinism is why the name is a
stack output: it is the subject every IRSA trust policy and EKS Pod
Identity association is written against.

## Engine vs Fleet Declarations

This component installs and configures the ENGINE. WHAT to provision is
declared separately, per fleet:

- **KubernetesKarpenterNodePool** — the fleet shape: instance-type /
  zone / capacity-type constraints, taints and labels, node lifetime,
  consolidation policy, resource ceilings.
- **KubernetesKarpenterEc2NodeClass** — the cloud-level machine
  template: AMIs, subnets, security groups, IAM identity, disks, kubelet
  configuration, metadata-service posture.

A NodePool references its NodeClass through a foreign key
(`node_class_ref.name` defaults to KubernetesKarpenterEc2NodeClass's
`status.outputs.node_class_name`), so an infra chart deploys the chain in
dependency order: engine → NodeClass → NodePool. An installation without
at least one NodePool provisions nothing — by design, the engine and the
fleet declarations have different lifecycles (cluster infrastructure vs
capacity policy that churns with the workloads).

## CRD Lifecycle: the Keep Mechanism

The `karpenter-crd` chart serves CRDs as ordinary templates, which means
Helm OWNS them — a plain uninstall would cascade-delete the CRDs and with
them every NodePool/EC2NodeClass/NodeClaim record in the cluster. The
CRD chart's whole values surface is one knob (`additionalAnnotations`,
stamped onto every CRD it templates), and the spec's
`crds.keep_on_uninstall` (default true) rides it as the standard
`helm.sh/resource-policy: keep` annotation. A CEL rule refuses
`keep_on_uninstall: true` with `install: false` — nothing to keep when
something else manages the CRDs.

## Cluster Identity

`cluster.name` is the one value the chart REFUSES to render without (the
deployment template wraps it in `required`). On EKS, `eks_control_plane`
lets the controller discover the endpoint and CA from the
DescribeCluster API, so `endpoint` can stay empty; non-EKS control
planes need the `https://` endpoint declared. `ca_bundle` is for TLS
bootstrap of provisioned nodes and is almost never set — empty means the
controller's own API-server TLS configuration is used.

## The AWS Arm

- **`irsa_role_arn`** annotates the `karpenter` service account with
  `eks.amazonaws.com/role-arn` so the controller calls
  EC2/EKS/SQS/Pricing without stored keys. Leaving it EMPTY is itself a
  posture: EKS Pod Identity associations are configured entirely on the
  AWS side and need no annotation, so the module renders nothing.
- **`interruption_queue`** names the SQS queue receiving EC2
  interruption events (spot interruptions, scheduled maintenance,
  instance rebalance). Empty disables interruption handling —
  provisioning still works, but Karpenter cannot drain ahead of an
  interruption.
- **`isolated_vpc`** tells Karpenter the cluster has no
  internet-reachable AWS endpoints beyond the provisioned VPC endpoints;
  it then avoids services without one — including the pricing API, so
  price-aware consolidation falls back to static data.
- **`reserved_enis`** (chart default 0) reserves ENIs outside the
  max-pods and kube-reserved math, for VPC CNI custom networking.
- **`enable_zonal_shift`** respects AWS zonal-shift signals when placing
  NodeClaims (requires the ARC permissions on the controller role).
- **`vm_memory_overhead_percent`** (chart default 0.075) is subtracted
  from every instance type's reported memory. Two type-fidelity details
  the modules preserve: the chart's `reservedENIs` default is the STRING
  "0" while `vmMemoryOverheadPercent` is the NUMBER 0.075, so the
  modules render one as a string and parse the other into a number —
  keeping the rendered values byte-compatible with the served chart's
  types on both engines.

## Batching and Scheduling Posture

`batching` controls how long Karpenter gathers pending pods before
computing a provisioning decision: `max_duration` (chart default 10s) is
the ceiling, `idle_duration` (chart default 1s) closes a batch early
when no new pods arrive. Longer windows consider more pods at once and
usually produce fewer, larger nodes.

`scheduling` sets the simulation posture: `preference_policy` (`Respect`
chart default — preferred affinities and ScheduleAnyways topology
constraints shape the node; `Ignore` counts only hard requirements,
which can pick cheaper shapes) and `min_values_policy` (`Strict` chart
default — scheduling fails when a requirement's minValues cannot be met;
`BestEffort` relaxes them rather than leave pods pending).

## Feature Gates

The spec mirrors the chart's `featureGates` block at the pinned release.
`reserved_capacity` is BETA and chart-default TRUE (NodeClasses can
select On-Demand Capacity Reservations and Karpenter launches into them
first); the other five — `node_repair`, `node_overlay`,
`spot_to_spot_consolidation`, `static_capacity`, `capacity_buffer` — are
ALPHA and default off. The modules render the whole six-key map with
defaults applied, because the deployment template composes the
FEATURE_GATES environment variable from all keys unconditionally —
explicit is safer than sparse, especially with one default diverging
from the rest. Gates follow upstream's graduation process; flipping one
on opts into pre-GA behavior (`static_capacity` in particular is what
NodePool `replicas` requires).

## Controller Scheduling: Where Karpenter Itself Runs

The chart pins controller pods away from Karpenter-provisioned nodes (a
required node affinity on `karpenter.sh/nodepool` DoesNotExist) — the
controller never disrupts its own machine, so it must run on capacity it
does NOT manage. The typed fields layer on top: `priority_class_name`
(chart default `system-cluster-critical`), `node_selector` (MERGES onto
the chart's `kubernetes.io/os=linux` — Helm deep-merges maps),
`tolerations` (REPLACE the chart's CriticalAddonsOnly default — Helm
replaces lists wholesale, so the modules render them only when set), and
`host_network` for the chicken-and-egg cluster whose CNI cannot serve
pod IPs before Karpenter runs.

## Typed Surface vs Escape Hatch

The typed spec covers namespace and lifecycle, chart version, CRD
lifecycle, cluster identity, the AWS provider arm, controller sizing,
batching windows, scheduling posture, feature gates, controller-pod
scheduling, and own telemetry.

`helm_values` merges LAST with Helm `-f` semantics on both engines (maps
deep-merge with the later document winning, lists replace). Deliberately
unmodeled as typed fields (all reachable via `helm_values`):

- **Image digests and registries** — private-mirror overrides, each with
  a correct chart default
- **Extra volumes, sidecars, DNS config** — the chart's operational long
  tail
- **Pod disruption budget tuning** (`podDisruptionBudget`) — the chart
  default matches the two-replica posture
- **Per-component affinity/topology overrides** — the chart's own
  anti-affinity and zone spread are the correct posture for almost every
  cluster

## Install Semantics

Both engines install the two releases as REAL Helm releases, atomically,
with cleanup on fail and a 600s timeout (covering image pulls on cold
clusters), CRDs strictly before the controller — the controller
reconciles NodePools from the moment it starts, so the CRDs must exist
first. One engine asymmetry worth recording: Pulumi's Helm release
resolves `oci://` registries through the joined chart reference
(`oci://public.ecr.aws/karpenter/<chart>`), while the Terraform provider
takes the repository and bare chart name separately — same chart bytes,
different wiring. The module (not Helm) owns namespace creation via
`create_namespace`, so a namespace it creates carries the standard
governance labels and is deleted with the resource; the flag stays false
for the recommended `kube-system` home.

## Outputs

`namespace`, `release_name` (fixed `karpenter`), `crd_release_name`
(fixed `karpenter-crd`, empty when `crds.install` is false), and
`service_account_name` (the chart-derived `karpenter` — the subject
every IRSA trust policy and EKS Pod Identity association is written
against).

## E2E

The behavioral facts are properties of the platform, not of any one test
run:

- The engine alone provisions nothing: a healthy install with no
  NodePool is the correct steady state, and the readiness proof is the
  controller Deployment becoming Available.
- The full behavioral proof needs the chain — engine, then an
  EC2NodeClass, then a NodePool referencing it, then an unschedulable
  pod; a NodeClaim materializing and a node registering under
  `karpenter.sh/nodepool` is the provisioning proof.
- Fleet declarations cannot precede the engine: NodePool/EC2NodeClass
  applies fail before the CRD release exists.
- The ServiceMonitor arm fails the release on clusters without the
  Prometheus operator CRDs, by design.
- Uninstalling keeps the CRDs (and every NodePool/EC2NodeClass/NodeClaim
  record) by default via the keep annotation on the CRD release.
