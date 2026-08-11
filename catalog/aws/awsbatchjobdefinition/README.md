# AwsBatchJobDefinition

An AWS Batch job definition — the versioned container blueprint jobs are submitted from: image, command, sizing, IAM identities, retries, and timeout.

## What It Is

A job definition is the workload half of AWS Batch: the [compute environment](../awsbatchcomputeenvironment/README.md) provides capacity, the [job queue](../awsbatchjobqueue/README.md) routes, and the job definition describes WHAT runs. It is referenced at submission time (SubmitJob) and by EventBridge Batch targets.

Job definitions are **revisioned and immutable**: every change registers a NEW revision of the name rather than mutating the old one. Because the `job_definition_arn` output carries the revision, an EventBridge rule that references it by output picks up each new revision on its next deployment — "change the image tag, the schedule runs the new code" falls out of the composition.

This kind models both arms of AWS type `container` job definitions: single-container **ECS-based jobs** (EC2 and Fargate — the shape nearly every Batch workload uses) via `container`, and **Batch-on-EKS pod jobs** via `eks` (the workload half of a compute environment attached to an EKS cluster with `eksConfiguration`). Exactly one arm is set per definition. Multi-node parallel jobs (`multinode`) and multi-container ECS jobs remain deliberately unmodeled long-tail shapes.

## When to Use It

Every Batch workload needs a job definition. Define one per distinct workload (an ETL step, a transcode task, a simulation runner), then parameterize runs with `parameters` placeholders and per-job SubmitJob overrides instead of multiplying definitions.

## Key Facts

- **Sizing goes through resource requirements.** `vcpus` (fractional for Fargate: 0.25-16) and `memory_mib` are required; `gpus` reserves whole GPUs on EC2 GPU compute.
- **Two IAM identities by design.** `job_role` is the code's runtime identity (S3, DynamoDB, …); `execution_role` is the agent's setup identity (pull image, resolve secrets, write logs). Fargate REQUIRES the execution role.
- **Zero-config logging.** Without `log_configuration`, container logs land in CloudWatch under `/aws/batch/job` automatically.
- **Secrets never live in the definition.** `secrets` maps env-var names to Secrets Manager / SSM ARNs, resolved by the agent at job start.
- **Retry with discrimination.** `retry_strategy.evaluate_on_exit` can RETRY Spot reclaims (`"Host EC2*"` status reasons) while EXITing on real application failures.
- **Platform gating is validated up front.** Fargate jobs reject the EC2-only knobs (GPUs, privileged, ulimits, Linux parameters); EC2 jobs reject the Fargate-only ones (platform version, public IP, ephemeral storage, runtime platform); EKS jobs reject `platform_capabilities` and `propagate_tags` (ECS concepts).
- **EKS jobs are Kubernetes-native.** Sizing is cpu/memory quantities under `eks.containers[].resources`; AWS permissions come from `service_account_name` (IRSA / Pod Identity) instead of `job_role`; secrets mount as Kubernetes secret volumes. Note AWS's Batch-pod default is `hostNetwork: true` — set an explicit `false` for VPC-CNI pod networking.
- **Deleting the resource deregisters every ACTIVE revision** of the name.

## Spec Overview

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `region` | string | **Yes** | Region the definition is registered in. |
| `container` XOR `eks` | message | **Yes (exactly one)** | The workload arm: ECS-based container job or Batch-on-EKS pod job. |
| `container.image` | string | **Yes (container arm)** | Full image reference (ECR or public registry). |
| `container.vcpus` | double | **Yes (container arm)** | vCPUs (fractional on Fargate). |
| `container.memory_mib` | int | **Yes (container arm)** | Memory hard limit in MiB. |
| `container.command` | string[] | No | CMD override; supports `Ref::<key>` placeholders. |
| `container.job_role` / `execution_role` | ref → AwsIamRole | Fargate: execution required | Runtime vs setup identity. |
| `container.environment` / `secrets` | map | No | Plain env vars / Secrets-Manager-resolved env vars. |
| `container.volumes` + `mount_points` | list | No | EFS (refs) and host-path volumes. |
| `container.gpus`, `ulimits`, `linux_parameters`, `privileged` | — | No | EC2-only. |
| `container.runtime_platform`, `fargate_platform_version`, `assign_public_ip`, `ephemeral_storage_gib` | — | No | Fargate-only. |
| `eks.containers` | list (1-10) | **Yes (eks arm)** | Pod containers: image, command/args, env, Kubernetes cpu/memory/GPU quantities, security context, volume mounts. |
| `eks.init_containers` | list (0-10) | No | Setup containers run to completion before the main ones. |
| `eks.host_network`, `dns_policy`, `service_account_name`, `pod_labels`, `image_pull_secret_names`, `share_process_namespace` | — | No | Pod-level networking, DNS, identity (IRSA), metadata, and registry credentials. |
| `eks.volumes` | list | No | Kubernetes volumes: `empty_dir` (node or Memory), node `host_path`, or a `secret` — exactly one backing each. |
| `platform_capabilities` | string[] | No (default EC2; container arm only) | `EC2` and/or `FARGATE`. |
| `parameters` | map | No | Default `Ref::` placeholder values. |
| `retry_strategy` | message | No | 1-10 attempts + up to 5 evaluate-on-exit conditions. |
| `timeout.attempt_duration_seconds` | int | No | Hard wall-clock limit per attempt (≥60). |
| `scheduling_priority` | int | No | 0-9999; consulted only on fair-share queues. |
| `deregister_on_new_revision` | bool | No (default true) | Keep exactly one ACTIVE revision. |

## Outputs

| Field | Description |
|-------|-------------|
| `job_definition_arn` | Revision-carrying ARN — the primary handle; changes every revision. |
| `arn_without_revision` | Revisionless ARN for latest-ACTIVE consumers. |
| `job_definition_name` | The name revisions register under (from `metadata.name`). |
| `revision` | The revision number this deployment registered. |

## Example

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsBatchJobDefinition
metadata:
  name: nightly-etl
  org: my-org
spec:
  region: us-west-2
  container:
    image: 123456789012.dkr.ecr.us-west-2.amazonaws.com/etl:1.4.2
    command: ["python", "run.py", "Ref::dataset"]
    vcpus: 2
    memoryMib: 4096
    jobRole:
      valueFrom:
        kind: AwsIamRole
        name: etl-job-role
        fieldPath: status.outputs.role_arn
  parameters:
    dataset: s3://data-lake/default
  retryStrategy:
    attempts: 3
    evaluateOnExit:
      - action: RETRY
        onStatusReason: Host EC2*     # Spot reclaims re-run
      - action: EXIT
        onExitCode: "1*"              # real failures do not
  timeout:
    attemptDurationSeconds: 7200
```

The spec's field comments carry the full contract, including revision semantics and the modeling boundaries.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
