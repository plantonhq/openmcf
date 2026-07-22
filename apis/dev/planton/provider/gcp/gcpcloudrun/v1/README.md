# GCP Cloud Run

Deploys a Cloud Run service via `google_cloud_run_v2_service` (Terraform) or Pulumi `cloudrunv2.Service` — a fully managed, request-serving container deployment that scales from zero to thousands of instances, with the full revision-template surface: sidecars, probes, volumes, traffic splitting, direct VPC egress, and GPU inference.

## Overview

A Cloud Run service is two things at once: a stable serving endpoint (ingress posture, invoker policy, traffic splitting) and a revision template (containers, volumes, scaling, networking) describing the next revision to stamp out. Every template change creates a new immutable revision; the `traffic` block decides how requests split between revisions — which makes blue/green and canary rollouts a first-class, declarative operation rather than a scripted one.

Sidecar containers are first-class: a service runs its application container alongside collectors, proxies, or auth helpers, with explicit startup ordering (`depends_on`) and shared volumes. Private egress uses direct VPC networking (`network_interfaces`) — no connector infrastructure to size or pay for. The service honors the ambient-project contract: an empty `project_id` falls back to the provider's default project.

Custom domains compose rather than embed: a serverless network endpoint group (`GcpRegionNetworkEndpointGroup`) bridges the service into the backend service → URL map → target proxy → forwarding rule family, with DNS managed by `GcpDnsRecord`.

## Purpose

- **The serverless request-serving primitive**: public APIs, internal backends, and GPU inference endpoints share one honest resource shape.
- **Pre-deploy coherence**: the API's cross-field rules (public-grant XOR invoker-check-off, connector XOR direct VPC, value XOR secret reference per env var, exactly one volume source, no TCP liveness probes, GPU-flag-needs-accelerator) are enforced by validation before any cloud call.
- **Cost levers as first-class fields**: scale-to-zero, request-based vs instance-based billing (`cpu_idle`), startup CPU boost, MANUAL scaling mode, and per-revision instance caps are all modeled.

## Key Features

- **Containers**: multiple containers per instance with `depends_on` startup ordering, command/args overrides, working dir, single serving port with `h2c` (end-to-end HTTP/2 for gRPC), CPU/memory limits with `cpu_idle` and `startup_cpu_boost`
- **Environment**: literal values and Secret Manager references (secret + version) per variable
- **Probes**: startup (HTTP/TCP/gRPC) and liveness (HTTP/gRPC) with delays, periods, timeouts, thresholds, and custom headers
- **Volumes**: Cloud SQL Unix sockets (by `GcpCloudSql` reference), Secret Manager files with per-item paths/modes, in-memory or disk `empty_dir`, GCS FUSE buckets (by `GcpGcsBucket` reference), NFS shares
- **Scaling**: per-revision min/max instances (scale-to-zero), service-level scaling across revisions including MANUAL mode with a pinned instance count, per-instance request concurrency
- **Traffic**: percent splits and tagged preview URLs across latest/named revisions; deterministic revision naming for declarative blue/green
- **Networking**: direct VPC egress (network/subnetwork references + firewall tags) or a Serverless VPC Access connector; `ALL_TRAFFIC` vs `PRIVATE_RANGES_ONLY` egress; ingress posture down to internal-plus-load-balancer
- **Security**: runtime identity by `GcpServiceAccount` reference, public access via additive IAM grant XOR invoker-check-off, custom token audiences, Binary Authorization, CMEK image encryption by `GcpKmsKey` reference
- **GPU**: one accelerator per instance (`node_selector`), zonal-redundancy opt-out for cheaper capacity
- **Safety**: `deletion_protection` defaults to true — a destroy fails until the manifest opts out

## Stack Outputs

| Output | Description |
|---|---|
| `url` | Canonical serving URL (`https://<service>-<hash>-<region>.run.app`) |
| `service_name` | Name of the service as created in GCP — the handle serverless NEGs reference |
| `revision` | Latest ready revision name |
| `location` | Region the service is deployed in |
| `uid` | Server-assigned unique identifier, never reused |
| `urls` | Every URL serving this service |

## Deliberately not modeled (recorded reasons)

| Excluded Feature | Why |
|---|---|
| `readiness_probe`, `iap_enabled`, `default_uri_disabled`, `deletion_policy`, `multi_region_settings`, `health_check_disabled`, `gcs.mount_options`, `empty_dir` size semantics beyond the released surface | Exist only on the provider's unreleased main line, not on the released 6.x major the GCP modules pin. Revisit on the next provider major. |
| `mesh` (Cloud Service Mesh) | Beta-only on the released line; returns when it reaches the GA surface. |
| `build_config`, `base_image_uri`, `build_info` | Source-based deploys (`gcloud run deploy --source`) — build pipelines own image production; this component models image-based revisions, the shape every CI/CD integration produces. |
| Domain mapping (`google_cloud_run_domain_mapping`) | A separate resource with its own lifecycle, deliberately not folded in as a toggle. The production-grade custom-domain path is the composed load-balancer family (serverless NEG → backend service → URL map → HTTPS proxy → forwarding rule) with `GcpDnsRecord`. |
| `client` / `client_version` | API-client telemetry strings with no user-facing behavior. |
| `annotations` (service and revision) | External-tool metadata passthrough; returns on concrete integration pull. |
| Service IAM beyond the public-invoker grant | Fine-grained invoker grants to specific identities are IAM-family territory; the modeled toggle covers the public/authenticated split every service decides. |

## Related Components

- **GcpRegionNetworkEndpointGroup** — bridges this service into the global HTTPS load balancer (references `service_name`)
- **GcpCloudSql** — databases mounted via the Cloud SQL volume (references `connection_name`)
- **GcpServiceAccount** — the runtime identity (`service_account` reference)
- **GcpVpcNetwork** / **GcpSubnetwork** — the network direct VPC egress places instances into
- **GcpGcsBucket** — buckets mounted via GCS FUSE volumes
- **GcpKmsKey** — CMEK key for container image encryption
- **GcpArtifactRegistryRepo** — where the deployed images live

## Additional Resources

- [Cloud Run Documentation](https://cloud.google.com/run/docs)
- [Cloud Run Admin API v2](https://cloud.google.com/run/docs/reference/rest/v2/projects.locations.services)

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
