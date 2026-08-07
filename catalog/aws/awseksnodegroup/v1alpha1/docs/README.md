# AWS EKS Node Group: The Managed Worker Fleet

## What This Component Is

A managed EKS node group is an EC2 fleet AWS operates on your behalf:
it provisions the instances, registers them as cluster workers, replaces
what fails health checks, and rolls them node by node on version
changes. `AwsEksNodeGroup` models one pool -- clusters run several
(general, Spot, GPU, dedicated tiers), each with its own lifecycle,
which is exactly why the pool is a first-class node rather than a field
on the cluster.

Everything attaches by reference: the cluster by its `name` output, the
node role as a referenced `AwsIamRole`, the subnets as `AwsSubnet`
nodes, and optionally an `AwsLaunchTemplate` for the launch mechanics.

## Two Launch Styles, Honestly Separated

AWS supports two mutually-exclusive ways to define what a node group
launches, and the spec keeps them honest:

- **Inline**: `instanceTypes`, `diskSizeGb`, optional `remoteAccess`.
  AWS builds the launch mechanics. The simple path for ordinary pools.
- **Launch template**: a reference to an `AwsLaunchTemplate` that owns
  the AMI, instance type, storage, IMDSv2 posture, and instance tags.
  AWS then *rejects* the inline knobs ("...or the node group deployment
  will fail" in AWS's own words), so the spec rejects them first --
  three CEL rules enforce the exclusions at validation time instead of
  twenty minutes into a deploy.

`amiType` deliberately stays valid alongside a template: a template that
carries user-data or tags but no AMI still needs the family declared.
Only a template that pins a custom AMI makes `amiType` wrong, and that
combination fails loudly at AWS.

The template-driven pool is the fleet-rollout pattern from the EC2
world applied to Kubernetes: `version: $Default` means promoting a new
template version rolls the pool; a numeric pin freezes it until the pin
moves.

## The Node Role Carries Its Own Policies

The node role needs `AmazonEKSWorkerNodePolicy`,
`AmazonEC2ContainerRegistryReadOnly`, and `AmazonEKS_CNI_Policy`. Those
attachments belong on the `AwsIamRole` itself (`managedPolicyArns`) --
this module never modifies a role it merely references. A bare role
fails node registration with a clear AWS-side error, which is correct:
the role's configuration is the role's responsibility.

## Rollouts, Repair, and Version Discipline

- **`updateConfig`** bounds how version updates roll the fleet: exactly
  one of `maxUnavailable` / `maxUnavailablePercentage` (mirroring AWS's
  ExactlyOneOf), and `updateStrategy: MINIMAL` surges -- replacements
  launch before old nodes drain, so capacity never dips.
- **`nodeRepairConfig`** is AWS's managed auto-repair: unhealthy nodes
  are replaced or rebooted within declared parallelism bounds, pausing
  when too many nodes are unhealthy at once (a systemic-problem signal).
  Count and percentage forms are mutually exclusive, enforced in CEL.
- **`version` / `releaseVersion`** pin the Kubernetes minor and the
  exact AMI release. The operational pattern: pin during a control-plane
  upgrade, then bump to roll on your schedule. Scaling allows
  min=desired=0, so a dormant pool is expressible.

## Deliberately Not Modeled

Bounded by the 90/10 rule; each skip is additive later if real
architectures pull for it:

- **`node_group_name_prefix`** -- Planton names deterministically from
  `metadata.name`; prefix-generated names would break the by-name
  composition model.
- **Windows-specific remote access nuances** -- the `remoteAccess` block
  models AWS's actual API (key + source security groups); RDP vs SSH
  port selection is AWS-side behavior, not configuration.

## Immutability and Naming

The pool name comes from `metadata.name` (AWS limit 63 characters,
truncated deterministically). Create-only at AWS: cluster, role,
subnets, `instanceTypes`, `amiType`, `capacityType`, `diskSizeGb`, and
the whole `remoteAccess` block -- changing any of them replaces the
pool, which is routine for node groups (workloads drain to the
replacement). Labels, taints, scaling, and update/repair configuration
update in place.

## Dual-Engine Implementation

`AwsEksNodeGroup` ships both a Terraform/OpenTofu module and a Pulumi
(Go) module at behavioral parity: the same launch-template version
defaulting (`$Default` when unset), the same enum mapping for capacity
types, identical null-guarding of unset optionals so AWS defaults keep
applying, identity tags on the group in both engines, and the same four
outputs -- including the real `asg_name` and `remote_access_sg_id`
values read from the node group's `resources` attribute. No
parity exceptions are required at pulumi-aws v7.35.0 / hashicorp/aws v6.
