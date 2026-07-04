# GCP Cloud Run

Deploys a Cloud Run service — a fully managed, request-serving container deployment that scales from zero to thousands of instances. One resource carries the whole story: the serving endpoint (ingress, invoker policy, traffic splitting) and the revision template (containers, volumes, scaling, networking) that every deploy stamps into a new immutable revision.

## What Gets Created

When you deploy a GcpCloudRun resource, Planton provisions:

- **Cloud Run service** — a `google_cloud_run_v2_service` in the chosen region; the Cloud Run Admin API is enabled automatically so a fresh project works on the first deploy
- **Containers** — the serving container plus any sidecars, with startup ordering, probes, Secret Manager environment variables, and volume mounts as configured
- **Volumes** — Cloud SQL sockets, Secret Manager files, scratch space, GCS FUSE buckets, and NFS shares as configured
- **Public-invoker grant** — when `allowUnauthenticated` is true, an IAM member granting `roles/run.invoker` to `allUsers`
- **Platform attribution** — organization, environment, and resource labels applied on the service object

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **A container image** in a registry the runtime identity can read ([GcpArtifactRegistryRepo](/docs/catalog/gcp/gcpartifactregistryrepo))
- **A service account** ([GcpServiceAccount](/docs/catalog/gcp/gcpserviceaccount)) for production services — the Compute Engine default SA is used otherwise
- **A VPC and subnetwork** ([GcpVpcNetwork](/docs/catalog/gcp/gcpvpcnetwork), [GcpSubnetwork](/docs/catalog/gcp/gcpsubnetwork)) if using direct VPC egress
- **Secret Manager secrets** already created if referencing them in env vars or volumes (the runtime identity needs `roles/secretmanager.secretAccessor`)

## Quick Start

Create a file `service.yaml`:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpCloudRun
metadata:
  name: my-api
spec:
  region: us-central1
  containers:
    - image: us-docker.pkg.dev/my-project/my-repo/api:1.0.0
      ports:
        containerPort: 8080
  allowUnauthenticated: true
  deletionProtection: false
```

Deploy:

```shell
planton apply -f service.yaml
```

This creates a public, scale-to-zero HTTP service with Cloud Run's defaults (1 CPU / 512Mi, concurrency 80, 300s timeout) and prints its `https://….run.app` URL.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | Region the service is deployed in (e.g. `us-central1`). Immutable. | Required |
| `containers` | `list` | The instance's containers; the first conventionally serves requests. | At least 1; each needs `image` |

### Service Identity & Metadata

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `projectId` | `StringValueOrRef` | provider default | GCP project. Can reference a GcpProject resource. |
| `serviceName` | `string` | `metadata.name` | Service name in GCP (1–63 chars, lowercase). Immutable. |
| `description` | `string` | — | Human-readable description shown in the console. |
| `labels` | `map` | — | User labels (shared with billing); merged beneath platform attribution labels. |
| `serviceAccount` | `StringValueOrRef` | Compute Engine default SA | Runtime identity. Reference a GcpServiceAccount's `email` output. |
| `deletionProtection` | `bool` | `true` | A destroy fails until this is explicitly set false. |

### Containers (`containers[]`)

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `name` | `string` | auto | Container name; required for multi-container services (`dependsOn` refers to it). |
| `image` | `string` | — | Container image URL. Resolved to a digest at revision creation. |
| `command` / `args` | `string[]` | image defaults | Entrypoint and argument overrides. |
| `env[].name` | `string` | — | Variable name. |
| `env[].value` | `string` | — | Literal value (never place credentials here). |
| `env[].valueFromSecret` | object | — | Secret Manager reference: `secret` + `version` (`latest` or a number). Exclusive with `value`. |
| `ports.containerPort` | `int` | `8080` | The single serving port (one container per service may set it). |
| `ports.name` | `string` | `http1` | `h2c` enables end-to-end HTTP/2 (required for gRPC streaming). |
| `resources.cpu` | `string` | `1` | CPU limit: `1`, `2`, `4`, `8`, or fractional (`0.5`, `500m`). |
| `resources.memory` | `string` | `512Mi` | Memory limit with unit suffix (`512Mi`, `2Gi`). Minimums scale with CPU. |
| `resources.cpuIdle` | `bool` | `true` | `true` = request-based billing; `false` = instance-based (CPU always allocated — required for background work). |
| `resources.startupCpuBoost` | `bool` | `false` | Extra CPU during startup; cuts cold-start latency for JIT runtimes. |
| `volumeMounts[]` | list | — | `{name, mountPath}` pairs; Cloud SQL volumes must mount at `/cloudsql`. |
| `workingDir` | `string` | image default | Working directory for the entrypoint. |
| `startupProbe` | object | TCP on the port | Gates instance start: HTTP, TCP, or gRPC with delays/periods/thresholds. |
| `livenessProbe` | object | — | Restarts unhealthy containers: HTTP or gRPC (Cloud Run rejects TCP liveness). |
| `dependsOn` | `string[]` | — | Containers this one waits for — sidecar startup ordering. |

### Volumes (`volumes[]` — exactly one source each)

| Source | Fields | Description |
|--------|--------|-------------|
| `cloudSqlInstance` | `instances[]` (refs → GcpCloudSql `connection_name`) | Managed Unix sockets under the mount path — no sidecar, no VPC needed. |
| `secret` | `secret`, `defaultMode`, `items[]{path, version, mode}` | Secret Manager versions projected as files. |
| `emptyDir` | `medium` (`MEMORY`/`DISK`), `sizeLimit` | Per-instance scratch space; MEMORY counts against the memory limit. |
| `gcs` | `bucket` (ref → GcpGcsBucket), `readOnly` | Cloud Storage FUSE mount; requires GEN2. |
| `nfs` | `server`, `path`, `readOnly` | NFS share (e.g. Filestore); requires GEN2 and VPC access. |

### Scaling & Runtime

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `scaling.minInstanceCount` | `int` | `0` | Warm instances per revision; 0 scales to zero, 1+ eliminates cold starts. |
| `scaling.maxInstanceCount` | `int` | GCP default (100) | Per-revision instance cap — the cost/overload circuit breaker. |
| `serviceScaling.scalingMode` | `string` | `AUTOMATIC` | `MANUAL` pins the total instance count regardless of traffic. |
| `serviceScaling.manualInstanceCount` | `int` | — | Exact instance count in MANUAL mode. |
| `serviceScaling.minInstanceCount` | `int` | — | Service-level minimum, distributed across serving revisions. |
| `maxInstanceRequestConcurrency` | `int` | 80 (1 below 1 CPU) | Concurrent requests per instance (1–1000). |
| `timeoutSeconds` | `int` | `300` | Request timeout (1–3600); also bounds instance startup. |
| `executionEnvironment` | enum | GCP-selected | `GEN2` (full Linux, required for GCS/NFS) or `GEN1` (faster cold starts). |
| `sessionAffinity` | `bool` | `false` | Best-effort same-client-same-instance routing. |
| `revision` | `string` | auto-generated | Explicit next-revision name (must be prefixed with the service name) — enables declarative blue/green. |
| `encryptionKey` | `StringValueOrRef` | Google-managed | CMEK crypto key for deployed images. Reference a GcpKmsKey's `key_id` output. |

### Networking & Access

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `vpcAccess.networkInterfaces[]` | list | — | Direct VPC egress: `network`/`subnetwork` references + firewall `tags`. Exclusive with `connector`. |
| `vpcAccess.connector` | `StringValueOrRef` | — | Serverless VPC Access connector (the legacy egress mechanism). |
| `vpcAccess.egress` | `string` | `PRIVATE_RANGES_ONLY` | `ALL_TRAFFIC` routes everything through the VPC. |
| `ingress` | enum | `INGRESS_TRAFFIC_ALL` | `INTERNAL_ONLY` or `INTERNAL_LOAD_BALANCER` lock the run.app URL down. |
| `allowUnauthenticated` | `bool` | `false` | Grants `roles/run.invoker` to `allUsers`. Exclusive with `invokerIamDisabled`. |
| `invokerIamDisabled` | `bool` | `false` | Switches the IAM invoker check off entirely (org-policy alternative). |
| `customAudiences` | `string[]` | — | Extra accepted token audiences for authenticated callers. |

### Traffic Splitting (`traffic[]`)

| Field | Type | Description |
|-------|------|-------------|
| `type` | enum | `TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST` or `…_REVISION`. |
| `revision` | `string` | Revision name for REVISION targets. |
| `percent` | `int` | Share of traffic (all entries must sum to 100). |
| `tag` | `string` | Stable `<tag>---<host>` preview URL for smoke-testing before traffic moves. |

Empty `traffic` routes 100% to the latest ready revision.

### GPU & Deploy Gates

| Field | Type | Description |
|-------|------|-------------|
| `nodeSelector.accelerator` | `string` | GPU per instance (e.g. `nvidia-l4`); needs ≥4 CPU / 16Gi. |
| `gpuZonalRedundancyDisabled` | `bool` | Single-zone GPU serving — cheaper capacity, zonal risk. |
| `launchStage` | `string` | `BETA`/`ALPHA` declaration when using preview features. |
| `binaryAuthorization` | object | `useDefault` XOR `policy` + `breakglassJustification` — attestation-gated deploys. |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `url` | Canonical serving URL |
| `service_name` | Service name in GCP — what serverless NEGs reference |
| `revision` | Latest ready revision name |
| `location` | Deployed region |
| `uid` | Server-assigned unique identifier |
| `urls` | Every URL serving the service |

## Common Patterns

- **Public API** — one container, `allowUnauthenticated: true`, scale-to-zero (see the `01-public-api-service` preset)
- **Private backend** — internal ingress, direct VPC egress, Cloud SQL volume, secret env, dedicated identity (see `02-private-vpc-service`)
- **GPU inference** — `nodeSelector.accelerator`, instance-based billing, bounded max instances (see `03-gpu-inference`)
- **Canary rollout** — pin `revision` names and split `traffic` 90/10 with a `canary` tag; promote by editing percents
- **Custom domain** — compose a [GcpRegionNetworkEndpointGroup](/docs/catalog/gcp/gcpregionnetworkendpointgroup) → [GcpBackendService](/docs/catalog/gcp/gcpbackendservice) → [GcpUrlMap](/docs/catalog/gcp/gcpurlmap) → [GcpTargetHttpsProxy](/docs/catalog/gcp/gcptargethttpsproxy) → [GcpGlobalForwardingRule](/docs/catalog/gcp/gcpglobalforwardingrule) with [GcpDnsRecord](/docs/catalog/gcp/gcpdnsrecord)

## Related Components

- [GcpRegionNetworkEndpointGroup](/docs/catalog/gcp/gcpregionnetworkendpointgroup) — the bridge into the global HTTPS load balancer
- [GcpCloudSql](/docs/catalog/gcp/gcpcloudsql) — databases mounted via Cloud SQL volumes
- [GcpServiceAccount](/docs/catalog/gcp/gcpserviceaccount) — the runtime identity
- [GcpVpcNetwork](/docs/catalog/gcp/gcpvpcnetwork) / [GcpSubnetwork](/docs/catalog/gcp/gcpsubnetwork) — direct VPC egress targets
- [GcpGcsBucket](/docs/catalog/gcp/gcpgcsbucket) — GCS FUSE volume origins
- [GcpKmsKey](/docs/catalog/gcp/gcpkmskey) — CMEK image encryption
- [GcpArtifactRegistryRepo](/docs/catalog/gcp/gcpartifactregistryrepo) — image storage
