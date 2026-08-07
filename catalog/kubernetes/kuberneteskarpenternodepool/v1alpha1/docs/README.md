# KubernetesKarpenterNodePool: Research and Design

## Introduction

A Karpenter NodePool is the fleet declaration: it answers "what nodes
may exist" — the scheduling constraints new machines must satisfy, the
taints and labels they carry, how long they live, how aggressively
under-used ones are consolidated away, and the resource ceiling the
whole pool may reach. The engine (KubernetesKarpenter) watches these
declarations and provisions against them; clusters routinely run several
pools — a default on-demand pool, a spot pool, a GPU pool — and
Karpenter picks among them by weight.

This component renders the CLUSTER-SCOPED `karpenter.sh/v1` NodePool,
named after `metadata.name`. The spec holds 100% fidelity with the
upstream CRD at the pinned release (the controller source's
`pkg/apis/crds/karpenter.sh_nodepools.yaml`), and the CRD's own
validation rules are mirrored into the spec as CEL so mistakes surface
at validate time instead of at apply.

## Boundary: Pool vs Class vs Engine

Three kinds share the provisioning story, with distinct lifecycles:

- **KubernetesKarpenter** — the ENGINE: controller, CRDs, cloud
  integration. Cluster infrastructure with a fixed identity.
- **KubernetesKarpenterEc2NodeClass** — the machine TEMPLATE: AMIs,
  subnets, security groups, IAM, disks, kubelet. One class is typically
  shared by several pools.
- **KubernetesKarpenterNodePool** (this kind) — the fleet SHAPE:
  constraints, taints, lifetime, disruption policy, ceilings.

The pool references its class through `node_class_ref` — a foreign key
whose defaults address the AWS provider (`group: karpenter.k8s.aws`,
`kind: EC2NodeClass`, and a name defaulting to
KubernetesKarpenterEc2NodeClass's `status.outputs.node_class_name`).
Wired with `valueFrom` in an infra chart, the pool deploys after its
machine template; a literal `value:` serves NodeClasses not managed by
the platform. `group` and `kind` are immutable on the live API after
creation — mirrored from the CRD.

## The NodeClaim Template

`template` is what each launched node looks like:

- **`requirements`** (1–100, at least one required — an unconstrained
  pool launches anything) are node-selector terms with Karpenter's
  extended operators: `In`/`NotIn` (values required for `In`),
  `Exists`/`DoesNotExist` (no values), and `Gt`/`Lt`/`Gte`/`Lte`
  (exactly one non-negative integer). Keys are the well-known scheduling
  keys (`karpenter.sh/capacity-type`,
  `node.kubernetes.io/instance-type`, `topology.kubernetes.io/zone`,
  `kubernetes.io/arch`) and the `karpenter.k8s.aws/instance-*` discovery
  keys. The CRD's domain restrictions are mirrored:
  `kubernetes.io/hostname` and `karpenter.sh/nodepool` can never be
  constrained, and the `karpenter.sh` domain only allows
  `karpenter.sh/capacity-type`.
- **`min_values`** (ALPHA, 1–50, only on `In`) forces the launched fleet
  to span at least that many distinct values — the
  instance-type-diversity knob for spot pools, and a requirement's
  values list must be at least that long (mirrored CEL).
- **`labels`** (max 100, same restricted domains) and **`annotations`**
  stamp every node.
- **`taints`** implement the dedicated-pool pattern — pods must tolerate
  them to schedule onto the pool. **`startup_taints`** cover
  initialization ordering: applied at startup, expected to be removed by
  a tolerating DaemonSet, and pods do NOT need to tolerate them for the
  pool to provision.
- **`expire_after`** (CRD default `720h`, or `Never`) is Karpenter's
  eventually-consistent node-recycling mechanism;
  **`termination_grace_period`** is the cluster admin's guarantee that
  nodes CAN be cycled — it overrides pod terminationGracePeriod math and
  bypasses blocking PodDisruptionBudgets once reached.

## Disruption Policy

`consolidation_policy` picks which nodes consolidation may touch:
`WhenEmptyOrUnderutilized` (CRD default — replace or remove whenever
cheaper capacity fits the pods), `WhenEmpty` (only nodes with no
workload pods), or `Balanced`. `consolidate_after` (CRD default `0s`;
`Never` disables consolidation) is the dwell time before acting.

`budgets` bound concurrent disruption: each budget is an absolute count
(`"2"`) or percentage (`"10%"`), optionally active only on a cron
schedule with a duration (the two must be set together — the schedule
starts the window, the duration ends it — mirrored CEL), and optionally
scoped to reasons (`Underutilized`, `Empty`, `Drifted`). `"0"` blocks
disruption entirely while active — the business-hours freeze pattern.
With no budgets the CRD default is one always-active budget of 10%;
multiple active budgets apply the most restrictive value.

## Ceilings, Weight, and Static Capacity

`limits` caps the whole pool as Kubernetes quantities keyed by resource
name (`cpu`, `memory`, `nodes`, ...) — provisioning stops at the limit.
`weight` (1–100; absent means 0) ranks pools when several could satisfy
a scale-up.

`replicas` switches the pool into ALPHA static-capacity mode: a FIXED
node count instead of demand-driven provisioning, requiring the
staticCapacity feature gate on the Karpenter controller
(KubernetesKarpenter `feature_gates.static_capacity`). The CRD's
static-mode rules are mirrored as CEL: only the `nodes` limit is
supported and `weight` is forbidden when `replicas` is set. The modules
render `replicas` presence-aware — 0 is a valid static size and must
survive, while absence means demand-driven mode.

## Render Fidelity

Both engines create the NodePool as a typed custom resource (the Pulumi
module through a crd2pulumi-generated SDK, catching field-name and
structure errors at compile time). Unset optionals are omitted entirely
so the API server applies the CRD's own defaults — `weight` in
particular is presence-sensitive (an absent weight means 0, but 0 is not
an accepted literal value), and `limits` values stay strings because the
CRD field is int-or-string and the string form round-trips every
quantity.

## Outputs

`node_pool_name` — the cluster-scoped object's name (equals
`metadata.name`). It is also the value of the `karpenter.sh/nodepool`
label on every node the pool launches, making it the join key for
pool-scoped scheduling and monitoring.

## E2E

The behavioral facts are properties of the platform, not of any one test
run:

- A NodePool cannot be applied before the Karpenter CRDs exist — the
  engine (or at least its CRD release) strictly precedes fleet
  declarations in any fixture ordering.
- The provisioning proof is behavioral: an unschedulable pod matching
  the pool's requirements produces a NodeClaim and a registered node
  carrying `karpenter.sh/nodepool: <name>`.
- Static mode (`replicas`) requires the staticCapacity gate on the
  controller; without it the pool is accepted by the API server but
  never provisions a fixed fleet.
- The spec's mirrored CEL rules mean an invalid manifest (an `In`
  requirement without values, a budget schedule without a duration,
  weight on a static pool) fails platform validation before any
  apply is attempted.
