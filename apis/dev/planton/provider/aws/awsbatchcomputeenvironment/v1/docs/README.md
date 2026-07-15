# AWS Batch Compute Environment: Concepts

A compute environment is the capacity plane of AWS Batch. This reference covers how the Batch resource graph fits together, the in-place-update envelope that decides whether changes are cheap or destructive, and the deliberate modeling boundaries.

## The Batch Resource Graph

AWS Batch splits batch processing into four independently-owned resources, and the catalog mirrors that split:

| Resource | Owns | Referenced by |
|----------|------|---------------|
| **Compute environment** (this kind) | Capacity: instance selection, vCPU bounds, VPC placement, Spot strategy | Job queues (`compute_environment_order`) |
| **[Job queue](../../awsbatchjobqueue/v1/docs/README.md)** | Routing: which environments run jobs, in what preference order | SubmitJob callers, EventBridge Batch targets |
| **[Job definition](../../awsbatchjobdefinition/v1/docs/README.md)** | The workload: image, command, sizing, identities, retries | SubmitJob callers, EventBridge Batch targets |
| **[Scheduling policy](../../awsbatchschedulingpolicy/v1/docs/README.md)** | Fairness: capacity division across share identifiers | Job queues (`scheduling_policy`) |

The split is what makes the canonical patterns expressible:

- **Spot-first with On-Demand overflow** — one queue, two environments in preference order.
- **Zero-downtime environment replacement** — associate the new environment on the queue, drain and remove the old one.
- **Shared capacity, fair division** — many teams submit to one queue; the scheduling policy divides the environment's capacity by share.

## The In-Place-Update Envelope

The single most important operational fact about compute environments: **whether a change updates in place or replaces the environment depends on two settings**, verified in the provider's own update logic.

In-place infrastructure updates are possible only when BOTH hold:

1. `service_role` is unset (the environment uses the `AWSServiceRoleForBatch` service-linked role), and
2. `allocation_strategy` is `BEST_FIT_PROGRESSIVE`, `SPOT_CAPACITY_OPTIMIZED`, or `SPOT_PRICE_CAPACITY_OPTIMIZED`.

Inside the envelope, changes to instance types, key pair, AMI configuration, launch-template version, security groups, subnets, and instance tags roll through an infrastructure update governed by `update_policy` (wait for running jobs vs terminate after a timeout). Outside it, the same changes **replace the whole environment** — which also breaks any queue association until the replacement is VALID.

Always ForceNew regardless of the envelope: the environment's name, MANAGED/UNMANAGED type, `eks_configuration`, `placement_group`, `spot_iam_fleet_role`, and adding/removing the launch-template block.

Practical guidance: leave `service_role` unset and pick `BEST_FIT_PROGRESSIVE` (EC2) or `SPOT_PRICE_CAPACITY_OPTIMIZED` (SPOT) unless you have a specific reason not to — that keeps day-2 changes cheap.

## Fargate vs EC2 Honesty

AWS rejects the EC2-only knobs (min/desired vCPUs, instance types, allocation strategy, instance role, key pair, bid percentage, Spot Fleet role, launch template, EC2 configuration, placement group, instance tags) on Fargate environments, and requires security groups for Fargate task ENIs. The spec encodes both as CEL rules so a bad manifest fails at validation, not twenty minutes into a deploy.

## Spot Semantics

- The **Spot Fleet role** (`spot_iam_fleet_role`) is only used by the legacy `BEST_FIT` allocation strategy. The capacity-optimized strategies drive Spot directly and need no role — the spec's conditional requiredness mirrors AWS's actual rule rather than blanket-requiring the role for all SPOT environments.
- `bid_percentage` defaults to 100% of On-Demand when omitted; with capacity-optimized strategies the realized price is usually far below the cap, so setting a low bid mostly increases interruptions without saving much.

## Deliberate Modeling Boundaries

- **UNMANAGED environments are not modeled.** They exist for orgs that run their own ECS container instances; Batch then only tracks capacity. That is a fundamentally different operating model (you own scaling), rarely used, and additive later as a spec arm if concrete demand appears.
- **`compute_resources.image_id` (deprecated) is skipped.** The provider marks it deprecated in favor of `ec2_configuration.image_id_override`, which the spec models.
- **Batch-on-EKS is modeled at the environment level only** (`eks_configuration`: cluster ref + namespace). EKS job workloads (`eksProperties` job definitions) are a long-tail arm of the job-definition kind — an EKS-attached environment remains usable with job definitions registered outside this graph.

## Composition

| This kind references | Via | Output consumed |
|----------------------|-----|-----------------|
| AwsSubnet | `compute_resources.subnet_ids` | `subnet_id` |
| AwsSecurityGroup | `compute_resources.security_group_ids` | `security_group_id` |
| AwsIamInstanceProfile | `compute_resources.instance_role` | `instance_profile_arn` |
| AwsIamRole | `service_role`, `spot_iam_fleet_role` | `role_arn` |
| AwsLaunchTemplate | `compute_resources.launch_template` | `launch_template_id` |
| AwsEksCluster | `eks_configuration.eks_cluster_arn` | `cluster_arn` |

| Consumed by | Via | Output referenced |
|-------------|-----|-------------------|
| AwsBatchJobQueue | `compute_environment_order[].compute_environment` | `compute_environment_arn` |
