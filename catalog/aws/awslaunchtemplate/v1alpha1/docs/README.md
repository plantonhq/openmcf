# AWS Launch Template: The Blueprint EC2 Fleets Launch From

## What a Launch Template Is

An EC2 launch template captures everything RunInstances needs -- AMI,
instance type, storage, networking, IAM identity, metadata posture,
purchase options -- as a named, versioned object. It is not a detail of the
fleet that uses it: a template has its own ID and ARN, its own lifecycle,
and a one-to-many relationship with its consumers. One template can back
several auto-scaling groups, an EKS managed node group, and a Batch compute
environment at once; a mixed-instances group can draw from several
templates in a single fleet.

That independence is what makes the template the composition anchor of EC2
compute. Fleet managers decide *how many, where, and when*; the template
decides *what*. `AwsLaunchTemplate` models it as its own component so the
"what" is defined once, hardened once, and referenced everywhere.

## Versioning: Immutable Versions, Promoted Defaults

Launch template versions are immutable in AWS -- an update never edits an
existing version, it creates the next one. Both IaC modules set
`update_default_version`, so every applied change also promotes the new
version to the template's default. The consequences are worth
internalizing:

- Consumers referencing **`$Default`** (the common ASG and node-group
  setup) pick up every applied change on their next launch -- or
  immediately, when the group runs an instance refresh triggered by the
  template change.
- Consumers **pinned to a numeric version** are immune to template changes
  until their pin moves; the pin is an explicit, reviewable rollout step.
- **Rollback** is applying the previous spec: that creates yet another
  version whose content matches the old one and promotes it. History is
  linear and additive, exactly like the underlying AWS model.

Only the template *name* is create-only (it replaces the template and every
downstream reference). `latest_version` and `default_version` are exported
as stack outputs rather than accepted as inputs: in a declarative model the
spec describes one desired configuration, so the newest version is always
the intended default.

## Everything Is a Launch-Time Default

A template field is a default, not a constraint: any consumer may override
it (an ASG mixed-instances override replaces the instance type; an EKS node
group injects its own AMI). Leaving a field unset omits it from the
template so the consumer or the AWS account default decides. That is how an
org-wide golden template works -- set the opinionated parts, leave the
workload parts open:

- Opinionated: `metadataOptions.httpTokens: required` (IMDSv2),
  encrypted `gp3` root volume with a customer-managed KMS key,
  `detailedMonitoring: true`.
- Open: `imageId` and `instanceType` unset, supplied by each consumer.

Note the one asymmetry: an auto-scaling group requires the template it
references to carry an AMI, so ASG-consumed templates set `imageId` while
EKS-node-group-consumed templates usually leave it out.

## Attribute-Based Selection over Type Lists

`instanceRequirements` replaces "a list of instance types we hope stays
current" with a description of what the workload needs: memory and vCPU
ranges (required), CPU manufacturers (`amazon-web-services` selects
Graviton -- pair with an arm64 AMI), generations, accelerator filters,
local-storage needs, and price protection. AWS resolves the matching set at
launch, which:

- keeps the fleet current as new instance families ship, with no template
  edit;
- widens Spot to dozens of pools, so a single pool interruption cannot
  starve the group;
- expresses intent ("memory-optimized, current generation, no bare metal")
  instead of an inventory.

`instanceType` and `instanceRequirements` are mutually exclusive, mirroring
the AWS API. The same message shape appears in
`AwsAutoScalingGroup`'s mixed-instances overrides; each kind carries its own
copy because every component's proto surface is self-contained.

## The Security Surface

Three template decisions carry most of the fleet's security posture:

- **IMDSv2** (`httpTokens: required`): blocks the classic SSRF-to-
  credential-theft path. Hop limit 1 confines metadata to the host; use 2
  only when containers legitimately need instance credentials.
- **EBS encryption**: `encrypted: true` per mapping, with `kmsKeyId`
  referencing an `AwsKmsKey` for revocable, auditable customer-managed
  keys. Accounts with encryption-by-default are already covered; setting it
  in the template makes the posture explicit and portable.
- **Instance identity**: `instanceProfile` references an
  `AwsIamInstanceProfile`; keyless access via SSM Session Manager
  (no `keyName`) is the modern default.

## Deliberately Not Modeled

Bounded by the 90/10 rule; each skip is additive later if real
architectures pull for it:

- **`kernel_id` / `ram_disk_id`** -- paravirtual-era options with no
  current-generation use.
- **`security_group_names`** -- EC2-Classic addressing; retired by AWS.
- **License Manager specifications** (`license_specification`) -- niche
  license-tracking surface with no Planton License Manager kind to compose
  with.
- **Secondary interfaces** (`secondary_interfaces`) -- multi-card
  secondary-subnet plumbing for specialized network appliances.
- **Elastic GPU / Elastic Inference** -- removed by AWS (end-of-life);
  absent from the current provider schema too.
- **`network_performance_options.bandwidth_weighting`** -- niche bandwidth
  rebalancing on a handful of instance types.
- **Capacity reservation targeting** -- deferred until a capacity
  reservation kind exists to reference; a literal-only field would be dead
  weight in the graph.
- **User tag pass-through** (`tag_specifications` beyond identity tags) --
  custom user tags are a single platform-wide decision, not per-kind
  surface. The modules DO emit Planton identity tags on the template,
  instances, and volumes, because template tags do not propagate on their
  own and untagged fleet instances escape cost/cleanup queries.

## Dual-Engine Implementation

`AwsLaunchTemplate` ships both a Terraform/OpenTofu module and a Pulumi
(Go) module at behavioral parity. Both truncate the name to AWS's
125-character limit, base64-encode plain-text user data, promote new
versions to default, emit identity tags in the same three places, express
the nullable tri-states (`ebs_optimized`, public-IP association,
delete-on-termination) identically, and export the same outputs
(`launch_template_id`, `launch_template_arn`, `latest_version`,
`default_version`). Whichever engine a team standardizes on, the template
behaves identically.
