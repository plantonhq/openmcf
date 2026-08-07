# GcpCloudRunJob — Research and Design Documentation

## 1. Introduction

### What Is a Cloud Run Job?

Cloud Run jobs are GCP's run-to-completion container primitive: define a task template, then each **execution** stamps out `task_count` tasks with up to `parallelism` running concurrently. Tasks exit when the container process finishes — there is no HTTP endpoint, no traffic split, and no request-driven autoscaling.

The v2 API nests two templates. The **execution template** holds `task_count`, `parallelism`, and execution-level labels. The **task template** holds containers, volumes, identity, networking, timeout, and retries — the shape Planton flattens into `spec.template` plus top-level `taskCount`/`parallelism` for manifest clarity.

### The Composition Boundary

- **GcpCloudRunJob** owns the job definition (the template future executions run).
- **Executions** are separate API objects — triggered by `gcloud`, Scheduler, or Eventarc; not managed by this resource.
- **GcpCloudRun** is the request-serving sibling (services, revisions, traffic).
- **GcpCloudSql**, **GcpGcsBucket**, **GcpServiceAccount**, **GcpVpcNetwork**/**GcpSubnetwork**, and **GcpKmsKey** attach by reference as volumes, identity, egress network, and encryption key.

## 2. Deployment Methods Landscape

### Level 2: Terraform / OpenTofu

`google_cloud_run_v2_job` mirrors the service resource's container/volume vocabulary but drops serving fields. Sharp edges:

- **Double-nested template** — outer execution template + inner task template; easy to mis-nest containers.
- **Timeout as duration string** — `"600s"` in TF; Planton spec uses integer seconds and modules convert.
- **Parallelism vs task_count** — parallelism must not exceed task_count.

### Level 4: Planton

Validated protobuf spec compiled to BOTH engines. Pre-deploy rules catch env value/secret conflicts, VPC connector vs direct-VPC mutual exclusion, GPU zonal redundancy without accelerator, and binary authorization policy conflicts.

## 3. Deliberate Omissions (vs GcpCloudRun)

| Service field | Job rationale |
|---|---|
| `ingress`, `traffic`, `allowUnauthenticated` | Jobs are not request-serving |
| `probes` | Batch tasks exit on process completion; no traffic gating |
| `scaling`, `sessionAffinity`, `maxInstanceRequestConcurrency` | Execution parallelism replaces autoscaling |
| `cpu_idle`, `startup_cpu_boost` | Job container resources are limits-only |

## 4. Terraform Provider Floor

Designed from `google_cloud_run_v2_job` on the released Terraform Google provider 6.x line (`~> 6.0`). Both Planton engines enable `run.googleapis.com` before creating the job.

## 5. Registry

- **Enum:** 720 (opens the 720–729 GCP serverless overflow block)
- **ID prefix:** `cldrunjob`
