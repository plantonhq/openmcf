# GcpCloudRun — Research and Design Documentation

## 1. Introduction

### What Is a Cloud Run Service?

Cloud Run is GCP's fully managed, request-serving container platform: give it a container image and it serves HTTPS traffic, scaling instances from zero to thousands with per-request or per-instance billing. A Cloud Run *service* is the request-serving shape (the sibling primitives — jobs for run-to-completion work, worker pools for pull-based workers — are separate API resources with their own lifecycles).

The v2 API models a service as two layers. The **service layer** owns the stable endpoint: ingress posture, invoker policy, custom audiences, and the traffic split. The **revision template** describes the next revision: containers, volumes, scaling bounds, networking, and hardware. Every template change stamps out a new immutable revision; the traffic block routes requests across revisions by percent and tag. That two-layer shape is what makes blue/green and canary rollouts declarative — the rollout IS the manifest, not a script.

### The Composition Boundary

- **GcpCloudRun** owns the service: its endpoint, its revisions, and its runtime shape.
- **GcpRegionNetworkEndpointGroup** bridges the service into the load-balancer family by referencing its `service_name` output — custom domains, CDN, Cloud Armor, and multi-region serving all compose there.
- **GcpCloudSql**, **GcpGcsBucket**, **GcpServiceAccount**, **GcpVpcNetwork**/**GcpSubnetwork**, and **GcpKmsKey** attach by reference as volumes, identity, egress network, and encryption key.
- Domain mapping (the v1-API convenience that binds a bare domain to a service) is deliberately not folded into this resource: it is a separate API object with its own lifecycle, and the production-grade path is the composed load balancer with `GcpDnsRecord`.

## 2. Deployment Methods Landscape

### Level 0: Cloud Console

The "Deploy container" form. Good for discovery; unrepeatable by construction.

### Level 1: gcloud CLI

`gcloud run deploy` with flags. Scriptable but imperative — every flag omitted silently keeps the live value, so two deploys from two terminals diverge without anyone noticing.

### Level 2: Terraform / OpenTofu

`google_cloud_run_v2_service` is declarative with plan/apply. Its sharp edges:

- **The template is one giant block** — a typo in a nested probe or volume surfaces only at apply time.
- **Cross-field rules live in the API**: connector vs direct-VPC conflict, TCP liveness rejection, traffic percents that must sum to 100, GPU resource minimums — all apply-time failures.
- **Secret references are easy to get wrong**: a literal value and a secret source on the same env var, or a version omitted where the API requires one.

### Level 3: Pulumi

`cloudrunv2.Service` bridges the same schema into real languages — same surface, same sharp edges.

### Level 4: Planton

A validated protobuf spec compiled to BOTH engines with identical behavior. The spec turns the API's apply-time failures into manifest-time errors, wires every cross-resource input as a reference to a first-class kind, and documents every trade-off in the field comments so a wizard or an agent can configure the service from the spec alone.

## 3. The Planton Approach

### Pre-Deploy Coherence Rules

| Rule | What it prevents |
|---|---|
| `allow_unauthenticated` XOR `invoker_iam_disabled` | Granting public access through IAM while also switching the IAM check off — two mechanisms that must not combine |
| Env var: literal `value` XOR `value_from_secret` | The API's env-source exclusivity, rejected before a deploy |
| Volume: exactly one source (proto oneof) | Ambiguous volumes the API would reject |
| VPC access: `connector` XOR `network_interfaces` | The provider's conflicting egress mechanisms |
| Liveness probes: HTTP/gRPC only | Cloud Run rejects TCP liveness probes (TCP is startup-only) |
| Probe `timeout_seconds` ≤ `period_seconds` | The API's probe-timing rule |
| Traffic: REVISION targets name a revision, LATEST targets must not | Dead or ambiguous traffic entries |
| `manual_instance_count` requires MANUAL scaling mode | Dead config on automatic scaling |
| `gpu_zonal_redundancy_disabled` requires an accelerator | A GPU flag on a CPU service |
| Binary authorization: default policy XOR named policy | The provider's policy pairing |
| Revision/service/tag names, CPU/memory formats, egress/medium/stage enums validated by pattern | Typo-driven apply failures |

### Modeled Surface (the 90/10 floor)

Verified against the RELEASED provider line (google 6.50.0 schema dump) and pulumi-gcp v9.29.0 — never the provider's main branch:

| Family | Modeled |
|---|---|
| Identity & metadata | project (ambient default), region, service name (defaults to metadata.name), description, labels, runtime service account by reference |
| Containers | multiple containers with names + `depends_on` ordering, image, command/args, working dir, single serving port with `h2c`, CPU/memory limits, `cpu_idle` (billing model), `startup_cpu_boost` |
| Environment | literal values; Secret Manager references (secret + version) |
| Probes | startup (HTTP/TCP/gRPC) and liveness (HTTP/gRPC) with delays, periods, timeouts, thresholds, custom headers |
| Volumes | Cloud SQL sockets (refs), Secret Manager files (default mode + per-item path/version/mode), empty_dir (MEMORY/DISK + size limit), GCS FUSE (ref + read-only), NFS |
| Scaling | per-revision min/max, service-level scaling (AUTOMATIC/MANUAL + manual count + service minimum), request concurrency, timeout |
| Traffic | percent + tag splits across LATEST/named revisions; explicit revision naming |
| Networking | direct VPC egress (network/subnetwork refs + tags, multiple interfaces), connector (un-defaulted ref pending its kind), egress scope, ingress posture |
| Security & governance | public-invoker grant XOR invoker-check-off, custom audiences, Binary Authorization, CMEK by reference, launch stage |
| Hardware | GPU accelerator per instance, zonal-redundancy opt-out |
| Safety | deletion protection (default true), execution environment, session affinity |

### Deliberately Not Modeled (recorded reasons)

| Excluded | Reason |
|---|---|
| `readiness_probe`, `iap_enabled`, `default_uri_disabled`, `deletion_policy`, `multi_region_settings`, `health_check_disabled`, `gcs.mount_options` | Absent from the released 6.x provider line (main-branch-only). Revisit on the next provider major. |
| `mesh` (Cloud Service Mesh) | Beta-only on the released line. |
| `build_config`, `base_image_uri`, `build_info` | Source-based deploys — build pipelines own image production; the component models image-based revisions. |
| Domain mapping | A separate resource with its own lifecycle; the custom-domain path is the composed LB family + `GcpDnsRecord`. |
| `client`/`client_version`, `annotations` | API-client telemetry and external-tool metadata passthrough. |
| Per-identity invoker IAM | IAM-family territory; the public/authenticated toggle covers the decision every service makes. |

## 4. Implementation Notes

### Both Engines, One Contract

- Both modules enable `run.googleapis.com` before creating the service (`disable_on_destroy = false`) — destroying one service never disables the API project-wide.
- Public access is implemented identically: an IAM member granting `roles/run.invoker` to `allUsers`, created only when `allow_unauthenticated` is true. `invoker_iam_disabled` maps to the API field of the same name — the two are different mechanisms and the spec rejects combining them.
- Enum values pass through as the API's own names (the spec enums use provider-authentic value names), so neither engine carries a translation table that can drift.
- Empty optional strings are omitted (never sent as `""`); presence-carrying messages gate their blocks; `deletion_protection` unset means true in both engines.
- The spec's `timeout_seconds` integer becomes the API's `"Ns"` duration string identically on both sides.
- User labels merge beneath the platform attribution labels.
- The Terraform module runs on plain `google ~> 6.0` — every modeled field is GA on the released line (verified by schema dump).

### Immutability

Only `region`, `service_name`, and the project are immutable — everything else on a Cloud Run service updates in place by minting a new revision. That is the platform's superpower: configuration changes are rollouts, and rollbacks are traffic edits.

### Outputs

`service_name` is the handle serverless NEGs (and gcloud) address the service by; `url`/`urls` are the serving endpoints; `revision` is the latest ready revision (what LATEST traffic serves); `uid` is the server-assigned identifier that survives nothing — it is never reused after deletion, which makes it the honest identity for audit trails.

## 5. Production Best Practices

1. **A dedicated runtime service account per service** — the Compute default SA is over-privileged; grant the minimal set (Secret Manager accessor, SQL client, bucket reader) to a dedicated identity.
2. **Secrets by reference, never by value** — `value_from_secret` keeps material in Secret Manager; pin versions for deterministic rollbacks, use `latest` only where silent rotation is acceptable.
3. **Startup probes on real readiness** — the default TCP check passes as soon as the port opens; probe an endpoint that verifies dependencies so traffic never lands on a half-started instance.
4. **Lock ingress behind the load balancer** — services fronted by the composed HTTPS LB should set `INGRESS_TRAFFIC_INTERNAL_LOAD_BALANCER` so the run.app URL stops being a side door.
5. **Direct VPC egress over connectors** — no infrastructure to size or pay for; keep `PRIVATE_RANGES_ONLY` unless egress inspection requires routing everything.
6. **Canary with tags before percents** — a tagged revision gets a preview URL receiving zero traffic; smoke-test it, then move percents.
7. **`cpu_idle: false` only when needed** — instance-based billing is the right choice for background work and GPU serving, and a silent cost multiplier everywhere else.
8. **Keep `deletion_protection` on** — deleting a service tears down its endpoint and every revision; the default protects against exactly that.
