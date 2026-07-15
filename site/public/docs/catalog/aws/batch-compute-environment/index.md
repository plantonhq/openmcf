---
title: "Batch Compute Environment"
description: "Batch Compute Environment deployment documentation"
icon: "package"
order: 100
componentName: "awsbatchcomputeenvironment"
---

# AWS Batch Compute Environment

Deploys a MANAGED AWS Batch compute environment — the elastic pool of EC2 On-Demand, EC2 Spot, Fargate, or Fargate Spot compute that AWS Batch scales up and down to run submitted jobs. Job queues (`AwsBatchJobQueue`) map onto the environment in preference order, and job definitions (`AwsBatchJobDefinition`) describe what runs on it.

## What Gets Created

When you deploy an AwsBatchComputeEnvironment resource, Planton provisions:

- **Compute Environment** — an `aws_batch_compute_environment` of type `MANAGED` with the specified resource type (EC2/SPOT/FARGATE/FARGATE_SPOT), vCPU bounds, VPC placement, instance selection, optional launch template and AMI configuration, optional EKS attachment, and optional update policy

## Prerequisites

- **AWS credentials** configured via environment variables or Planton provider config
- **At least one VPC subnet** (private subnets recommended) — use `AwsSubnet` to provision
- **A security group** — **required for FARGATE/FARGATE_SPOT**; use `AwsSecurityGroup` to provision
- **For EC2/SPOT types**: an ECS instance profile — use `AwsIamInstanceProfile` wrapping a role with `AmazonEC2ContainerServiceforEC2Role`
- **For SPOT + BEST_FIT strategy only**: a Spot Fleet IAM role with `AmazonEC2SpotFleetTaggingRole` (the capacity-optimized strategies need no role)

## Quick Start

Create a file `batch-ce.yaml`:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsBatchComputeEnvironment
metadata:
  name: etl-fargate
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AwsBatchComputeEnvironment.etl-fargate
spec:
  region: us-west-2
  computeResources:
    type: FARGATE
    maxVcpus: 256
    subnetIds:
      - value: subnet-0a1b2c3d4e5f00001
      - value: subnet-0a1b2c3d4e5f00002
    securityGroupIds:
      - value: sg-0a1b2c3d4e5f00001
```

Deploy:

```shell
planton apply -f batch-ce.yaml
```

This creates a serverless Fargate environment with up to 256 vCPUs of capacity. Add an `AwsBatchJobQueue` mapped onto it to start submitting jobs.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | AWS region (e.g., `us-west-2`). | Required; non-empty |
| `computeResources.type` | string | `EC2`, `SPOT`, `FARGATE`, or `FARGATE_SPOT` | Required |
| `computeResources.maxVcpus` | int32 | Scale-out ceiling — the one knob updatable on every environment | >= 1 |
| `computeResources.subnetIds` | list(StringValueOrRef → AwsSubnet) | VPC placement; spread across AZs | At least 1 |
| `computeResources.securityGroupIds` | list(StringValueOrRef → AwsSecurityGroup) | Attached to instances / Fargate task ENIs | **Required for Fargate types** |
| `computeResources.instanceRole` | StringValueOrRef → AwsIamInstanceProfile | ECS instance profile | **Required for EC2/SPOT** |

### Optional Fields — Top Level

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `state` | string | `ENABLED` | `ENABLED` / `DISABLED` — the drain switch for maintenance or replacement. |
| `serviceRole` | StringValueOrRef → AwsIamRole | Service-linked role | Leave unset (recommended): the service-linked role is also a precondition for in-place infrastructure updates. |
| `eksConfiguration.eksClusterArn` | StringValueOrRef → AwsEksCluster | — | Attach the environment to an EKS cluster (Batch-on-EKS). Create-time only. |
| `eksConfiguration.kubernetesNamespace` | string | — | Namespace Batch launches job pods into. Create-time only. |
| `updatePolicy.terminateJobsOnUpdate` | bool | false | Terminate running jobs when instances are replaced during in-place updates. |
| `updatePolicy.jobExecutionTimeoutMinutes` | int32 | 30 (AWS) | How long to wait for running jobs before replacing instances anyway (1-360). |

### Optional Fields — Compute Resources (EC2/SPOT only)

| Field | Type | Description |
|-------|------|-------------|
| `minVcpus` | int32 | vCPU floor (default 0 — scale to zero when idle). |
| `desiredVcpus` | int32 | Initial vCPU target; Batch adjusts it continuously. |
| `instanceTypes` | list(string) | `["optimal"]` or explicit families/sizes. |
| `allocationStrategy` | string | `BEST_FIT`, `BEST_FIT_PROGRESSIVE`, `BEST_FIT_PROGRESSIVE_ORDERED`, `SPOT_CAPACITY_OPTIMIZED`, `SPOT_PRICE_CAPACITY_OPTIMIZED`, `SPOT_CAPACITY_OPTIMIZED_PRIORITIZED`. Only the middle three support in-place infrastructure updates. |
| `ec2KeyPair` | string | SSH key pair name (prefer SSM Session Manager). |
| `bidPercentage` | int32 | SPOT only: max % of On-Demand price (0-100; omit for 100). |
| `spotIamFleetRole` | StringValueOrRef → AwsIamRole | SPOT + BEST_FIT only: the Spot Fleet role. |
| `launchTemplate.launchTemplateId` | StringValueOrRef → AwsLaunchTemplate | Custom AMI / user data / IMDSv2 posture. |
| `launchTemplate.version` | string | Version number, `$Latest`, or `$Default`. |
| `ec2Configurations` | list(object), max 2 | Image family (`ECS_AL2023`, `ECS_AL2_NVIDIA`, `EKS_AL2023`, …), AMI override, EKS AMI Kubernetes version. |
| `placementGroup` | string | EC2 placement group for tightly-coupled multi-node jobs. Create-time only. |
| `resourceTags` | map(string) | Tags applied to launched EC2 instances / Spot requests. |

The spec rejects EC2-only fields on Fargate environments at validation time (AWS would reject them at deploy time).

## Outputs

| Output | Type | Description |
|--------|------|-------------|
| `compute_environment_arn` | string | What job queues reference in `computeEnvironmentOrder`. |
| `compute_environment_name` | string | The environment's name (from `metadata.name`). |
| `ecs_cluster_arn` | string | The ECS cluster Batch provisions behind the environment. |
| `status` | string | `VALID` / `INVALID` — queues can only associate VALID environments. |

## Presets

| Name | Description |
|------|-------------|
| [01-fargate-batch](presets/01-fargate-batch.yaml) | Serverless Fargate — zero instance management |
| [02-ec2-managed-batch](presets/02-ec2-managed-batch.yaml) | EC2 with optimal instances and an update policy |
| [03-spot-cost-optimized-batch](presets/03-spot-cost-optimized-batch.yaml) | Spot with the price-capacity-optimized strategy |

## Related Components

- [AwsBatchJobQueue](/docs/catalog/aws/batch-job-queue) — maps onto this environment in preference order
- [AwsBatchJobDefinition](/docs/catalog/aws/batch-job-definition) — the container blueprint jobs run from
- [AwsBatchSchedulingPolicy](/docs/catalog/aws/batch-scheduling-policy) — fair-share capacity division for queues
- [AwsLaunchTemplate](/docs/catalog/aws/launch-template) — custom AMI / user data for EC2 environments
- [AwsIamInstanceProfile](/docs/catalog/aws/iam-instance-profile) — the ECS instance profile EC2/SPOT environments require
