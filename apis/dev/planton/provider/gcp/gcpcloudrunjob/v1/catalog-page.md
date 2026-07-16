# GCP Cloud Run Job

Run-to-completion batch workloads on Cloud Run — parallel tasks, no serving endpoint. The job sibling to [GcpCloudRun](/docs/catalog/gcp/gcpcloudrun): same container/volume/VPC vocabulary, batch semantics instead of request-serving.

**Enum:** 720 · **ID prefix:** `cldrunjob` · **Provider:** GCP · **API:** `gcp.planton.dev/v1`

## At a Glance

| | |
|---|---|
| **Creates** | `google_cloud_run_v2_job` |
| **Trigger** | `gcloud run jobs execute`, Cloud Scheduler, Eventarc |
| **Not included** | Traffic, ingress, probes (service-only concerns) |
| **Engines** | Terraform (~> 6.0) and Pulumi |

## When to Use

- **Scheduled ETL / cleanup** — `taskCount: 1`, trigger on a cron
- **Parallel map-reduce** — `taskCount` + `parallelism` for shardable work
- **GPU batch inference** — `nodeSelector` + resource minimums
- **Private batch workers** — direct VPC egress to reach RFC1918 APIs

## Quick Example

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpCloudRunJob
metadata:
  name: nightly-sync
spec:
  region: us-central1
  template:
    containers:
      - image: us-docker.pkg.dev/my-project/repo/sync:1.0.0
  taskCount: 1
  deletionProtection: false
```

## Key Fields

- `template.containers[]` — worker image, env (literal or Secret Manager), resources
- `template.volumes[]` — Cloud SQL, secrets, scratch, GCS FUSE, NFS
- `taskCount` / `parallelism` — how many tasks and how many at once
- `template.timeoutSeconds` / `maxRetries` — per-attempt limits
- `template.vpcAccess` — connector or direct VPC for private egress

## Outputs

`jobName`, `location`, `uid`, `latestCreatedExecution`

## See Also

- [README](README.md) — full configuration reference
- [GcpCloudRun](/docs/catalog/gcp/gcpcloudrun) — HTTP services
- Presets: [basic ETL](presets/01-batch-etl-basic.yaml), [VPC cleanup](presets/02-parallel-vpc-cleanup.yaml), [GPU batch](presets/03-gpu-batch-inference.yaml)
