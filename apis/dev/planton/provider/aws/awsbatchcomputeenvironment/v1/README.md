# AwsBatchComputeEnvironment

An AWS Batch MANAGED compute environment — the elastic pool of EC2 On-Demand, EC2 Spot, Fargate, or Fargate Spot compute that AWS Batch scales up and down to run submitted jobs.

## What It Is

A compute environment is the capacity half of AWS Batch. It owns the infrastructure decisions — instance families, vCPU floor and ceiling, VPC placement, Spot strategy — and AWS Batch handles the scaling: zero instances when queues are empty, scale-out when jobs arrive, scale-in as they drain.

It is one node of the Batch resource graph, not the whole story:

- **[AwsBatchJobQueue](../awsbatchjobqueue/v1/README.md)** maps onto one or more compute environments in preference order and is where jobs are submitted.
- **[AwsBatchJobDefinition](../awsbatchjobdefinition/v1/README.md)** is the container blueprint jobs run from.
- **[AwsBatchSchedulingPolicy](../awsbatchschedulingpolicy/v1/README.md)** optionally divides a queue's capacity fairly across teams.

Keeping these first-class lets one queue span an On-Demand primary and a Spot overflow environment, and lets an environment be replaced behind a queue with zero queue downtime.

## When to Use It

| Use Case | Compute Type |
|----------|--------------|
| **Bursty container jobs, zero ops** | `FARGATE` — no instances to manage, per-second billing |
| **Cost-optimized batch processing** | `FARGATE_SPOT` or `SPOT` — up to ~90% cheaper, interruptible |
| **GPU / large-memory / custom-AMI workloads** | `EC2` with `instance_types` + optional launch template |
| **HPC and tightly-coupled jobs** | `EC2` with a `placement_group` |
| **Kubernetes-native batch** | Any type with `eks_configuration` (Batch-on-EKS) |

## Key Facts

- **Only MANAGED environments are modeled.** Batch owns the instance lifecycle. UNMANAGED (bring-your-own ECS container instances) is a different operating model and is deliberately not modeled.
- **The in-place-update envelope matters.** AWS can only update infrastructure in place when the environment uses the **service-linked role** (leave `service_role` unset) AND `allocation_strategy` is `BEST_FIT_PROGRESSIVE`, `SPOT_CAPACITY_OPTIMIZED`, or `SPOT_PRICE_CAPACITY_OPTIMIZED`. Outside that envelope, most compute changes replace the environment.
- **Fargate keeps it minimal.** For `FARGATE`/`FARGATE_SPOT` only `max_vcpus`, `subnet_ids`, and `security_group_ids` apply — the spec rejects the EC2-only knobs at validate time instead of letting AWS reject them mid-deploy.
- **Security groups are required for Fargate**, and for EC2/SPOT whenever the launch template does not supply its own.
- **Spot Fleet role is a BEST_FIT-only concern.** The modern capacity-optimized Spot strategies do not use Spot Fleet; `spot_iam_fleet_role` is only required for `SPOT` + `BEST_FIT` (or no strategy, which defaults to BEST_FIT).
- **Deletion is disable-then-delete.** The environment is disabled, drained, and then deleted; queues referencing it must be deleted or repointed first.

## Spec Overview

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `region` | string | **Yes** | AWS region. |
| `state` | string | No (default `ENABLED`) | `ENABLED` / `DISABLED` — the drain switch. |
| `service_role` | ref → AwsIamRole | No | Leave unset for the service-linked role (recommended). |
| `compute_resources.type` | string | **Yes** | `EC2`, `SPOT`, `FARGATE`, `FARGATE_SPOT`. |
| `compute_resources.max_vcpus` | int | **Yes** | Scale-out ceiling; updatable on every environment. |
| `compute_resources.min_vcpus` | int | No (default 0) | vCPU floor (EC2/SPOT); 0 scales to zero when idle. |
| `compute_resources.subnet_ids` | ref[] → AwsSubnet | **Yes** | VPC placement; spread across AZs. |
| `compute_resources.security_group_ids` | ref[] → AwsSecurityGroup | Fargate: yes | Attached to instances / task ENIs. |
| `compute_resources.instance_types` | string[] | EC2/SPOT | `["optimal"]` or explicit families/sizes. |
| `compute_resources.allocation_strategy` | string | No | Full provider set incl. the ordered and prioritized variants. |
| `compute_resources.instance_role` | ref → AwsIamInstanceProfile | EC2/SPOT: yes | The ECS instance profile. |
| `compute_resources.launch_template` | ref → AwsLaunchTemplate | No | Custom AMI / user data / IMDSv2 posture. |
| `compute_resources.ec2_configurations` | list (max 2) | No | Image family selection (`ECS_AL2023`, `ECS_AL2_NVIDIA`, `EKS_AL2023`, …). |
| `compute_resources.placement_group` | string | No | Low-latency placement for multi-node jobs. |
| `eks_configuration` | message | No | EKS cluster ref + namespace (Batch-on-EKS; create-time only). |
| `update_policy` | message | No | How running jobs are treated during infrastructure updates. |

## Outputs

| Field | Description |
|-------|-------------|
| `compute_environment_arn` | What job queues reference in `compute_environment_order`. |
| `compute_environment_name` | The environment's name (from `metadata.name`). |
| `ecs_cluster_arn` | The ECS cluster Batch provisions behind the environment. |
| `status` | `VALID` / `INVALID` — queues can only associate VALID environments. |

## Example

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsBatchComputeEnvironment
metadata:
  name: etl-fargate
  org: my-org
spec:
  region: us-west-2
  computeResources:
    type: FARGATE
    maxVcpus: 256
    subnetIds:
      - valueFrom:
          kind: AwsSubnet
          name: private-a
          fieldPath: status.outputs.subnet_id
      - valueFrom:
          kind: AwsSubnet
          name: private-b
          fieldPath: status.outputs.subnet_id
    securityGroupIds:
      - valueFrom:
          kind: AwsSecurityGroup
          name: batch-jobs
          fieldPath: status.outputs.security_group_id
```

Map a queue onto it:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsBatchJobQueue
metadata:
  name: etl-queue
spec:
  region: us-west-2
  priority: 10
  computeEnvironmentOrder:
    - order: 1
      computeEnvironment:
        valueFrom:
          kind: AwsBatchComputeEnvironment
          name: etl-fargate
          fieldPath: status.outputs.compute_environment_arn
```

See [docs/README.md](docs/README.md) for the family architecture and update-semantics deep dive.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
