# KubernetesKarpenterEc2NodeClass: Research and Design

## Introduction

An EC2NodeClass is the AWS provider's answer to "how is a node built":
which AMIs to boot, which subnets and security groups to attach, which
IAM identity nodes assume, how disks are laid out, how the kubelet is
configured, and how the instance metadata service is exposed. NodePools
answer "what nodes may exist" and reference a class for everything
machine-level — one class is typically shared by several pools.

This component renders the CLUSTER-SCOPED `karpenter.k8s.aws/v1`
EC2NodeClass, named after `metadata.name`. The spec holds 100% fidelity
with the upstream CRD at the pinned release (the controller source's
`pkg/apis/crds/karpenter.k8s.aws_ec2nodeclasses.yaml`), and the CRD's
own CEL rules are mirrored into the spec so mistakes surface at validate
time instead of at apply.

## Boundary: Class vs Pool vs Engine

The class is deliberately the SLOW-MOVING half of the fleet declaration:
AMI pipelines, subnet layouts, IAM roles, and disk standards change on
platform cadence, while pool constraints (instance types, capacity type,
taints) change with workload demand. Splitting them means several
NodePools — spot, on-demand, GPU — share one machine template, and an
AMI rollout is one class change instead of N pool changes.

The reference direction: KubernetesKarpenterNodePool's `node_class_ref`
carries a foreign key to this kind's `status.outputs.node_class_name`.
The engine (KubernetesKarpenter) precedes both — the EC2NodeClass CRD is
installed by its CRD release.

## Selector-Term Grammar

AMIs, subnets, security groups, and capacity reservations all use the
same grammar: a list of terms ORed together, with the fields within a
term ANDed. The CRD's exclusivity rules are mirrored as CEL:

- **AMI terms** (1–30): `alias` is `family@version` (families al2,
  al2023, bottlerocket, windows2019/2022/2025; Windows families only
  support `latest`). An alias term is mutually exclusive with every
  other selector field AND must be the list's only term. An `id` term
  (`ami-...`) combines with nothing else. `name` ANDs with `owner`
  (account IDs, `self`, `amazon`, `aws-marketplace`); `ssm_parameter`
  serves custom-AMI pipelines (the parameter holds the image id); `tags`
  select with `'*'` as the any-value wildcard.
- **`ami_family`** is the UserData/bootstrap format. It is inferred from
  an alias (and may then only be set to the alias's own family or
  `Custom`), and REQUIRED when selecting by id/name/tags/ssm-parameter —
  the family cannot be inferred. All three pairing rules are mirrored
  from the CRD.
- **Subnet terms** (1–30): `tags` — `karpenter.sh/discovery: <cluster>`
  is the convention EKS setups tag their private subnets with — or
  explicit `id`. Karpenter spreads across the selected subnets' zones.
- **Security-group terms** (1–30): `tags`, group `name` (the name
  FIELD, not the Name tag), or `id`, each exclusive.
- **Capacity-reservation terms** (0–30; requires the reservedCapacity
  feature gate, chart-default on): `id`, `tags`, or
  `instance_match_criteria` (`open`/`targeted`), with `owner_id` for
  cross-account reservations.

## Node Identity: role XOR instance_profile

Exactly one of `role` / `instance_profile` — mirrored from the CRD's
identity rule. `role` is the recommended arm: Karpenter creates and
manages the instance profile for the named IAM role (the controller
needs `iam:PassRole`). `instance_profile` brings a pre-existing profile
for accounts where Karpenter must not manage profiles. Either way the
nodes' identity is cloud-side composition: on EKS the role must also be
registered in the cluster's access configuration for nodes to join.

## Disks, Kubelet, IMDS

- **`block_device_mappings`** (max 50, at most one `root_volume`): each
  EBS mapping needs `snapshot_id` or `volume_size` — without either
  there is nothing to create (mirrored CEL) — plus type
  (standard/io1/io2/gp2/sc1/st1/gp3), IOPS, gp3 throughput (125–1000
  MiB/s), encryption and KMS key, and the snapshot-only
  `volume_initialization_rate`. Empty means the AMI family's defaults.
- **`kubelet`** is the upstream-supported subset merged into bootstrap
  user data: cluster DNS, CPU CFS quota, hard/soft eviction thresholds
  over the six eviction signals (every soft signal must pair with a
  grace period — mirrored CEL), image GC thresholds (high must exceed
  low), kube/system reserved (keys cpu, memory, ephemeral-storage, pid),
  max pods, and pods per core.
- **`metadata_options`** defaults to the EKS security best practice, per
  the CRD's own defaults: IMDSv2 `required`, hop limit 1 (pods cannot
  reach the node's instance credentials), HTTP endpoint enabled, IPv6
  endpoint disabled. The `http_tokens` field carries a mode enum, not a
  token value.
- **`user_data`** is merged INTO Karpenter's generated bootstrap
  configuration — format follows the AMI family (MIME multi-part for
  AL2/AL2023, TOML for Bottlerocket, PowerShell for Windows); Karpenter
  adds its own required bootstrap fields either way.

## Advanced Networking and Placement

`network_interfaces` (max 150) declares explicit multi-card / EFA
layouts: a primary `interface` at device index 0 on card 0 is required,
`(card, device)` pairs must be unique, and each card carries at most one
`efa-only` interface — all mirrored from the CRD.
`placement_group_selector` picks a placement group by name XOR id.
`connection_tracking` tunes ENI idle timeouts (at least one of the three
timeouts when the block is present). `ip_prefix_count` enables prefix
delegation for pod density; `instance_store_policy: RAID0` stripes
ephemeral NVMe into one volume for pod ephemeral storage;
`associate_public_ip_address` overrides the subnet's default only when
set.

## Tags and Controller-Owned Keys

`tags` land on the EC2 resources Karpenter creates (instances, launch
templates, volumes). The controller-owned keys are rejected, mirrored
from the CRD: `eks:eks-cluster-name`, `kubernetes.io/cluster/*`,
`karpenter.sh/nodepool`, `karpenter.sh/nodeclaim`, and
`karpenter.k8s.aws/ec2nodeclass`.

## Render Fidelity

Both engines create the EC2NodeClass as a typed custom resource (the
Pulumi module through a crd2pulumi-generated SDK, catching field-name
and structure errors at compile time). Unset optionals are omitted
entirely so CRD-side defaults — the IMDS block in particular — are never
overridden by rendered zero values. Several spec fields carry an
explicit `json_name` because the CRD's keys use acronym casing
(`associatePublicIPAddress`, `httpProtocolIPv6`, `clusterDNS`,
`cpuCFSQuota`, `imageGCHighThresholdPercent`, `kmsKeyID`, `snapshotID`,
`ownerID`) that protojson would otherwise derive differently — and the
API server rejects miscased keys as undeclared fields.

## Outputs

`node_class_name` — the cluster-scoped object's name (equals
`metadata.name`), the value NodePools reference through their
`node_class_ref`.

## E2E

The behavioral facts are properties of the platform, not of any one test
run:

- An EC2NodeClass cannot be applied before the Karpenter CRDs exist —
  the engine's CRD release strictly precedes it in fixture ordering.
- A class alone launches nothing; the behavioral proof requires a
  NodePool referencing it and an unschedulable pod. The class's own
  readiness signal is its status conditions (subnet/security-group/AMI
  resolution), which the controller populates without any pool.
- The mirrored CEL rules mean an invalid manifest (an alias term
  combined with tags, both role and instance_profile, an EBS mapping
  with neither snapshot nor size, a soft-eviction signal without a
  grace period) fails platform validation before any apply is
  attempted.
- Discovery-tag selection is account convention: fixtures relying on
  `karpenter.sh/discovery` tags only resolve in accounts whose
  subnets/security groups carry them.
