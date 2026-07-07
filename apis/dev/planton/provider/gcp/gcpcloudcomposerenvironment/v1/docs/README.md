# GCP Cloud Composer Environment — Deep Dive

## The problem this resource solves

Running production Apache Airflow means operating a Kubernetes cluster, a metadata database, a web server, and the glue between them. Cloud Composer collapses that into one managed resource: you declare sizing, networking, software, and security; Google assembles and operates the GKE cluster, Cloud SQL database, web server, and DAG bucket. Declaring the environment as infrastructure makes an organization's orchestration posture — who can reach the UI, what encrypts the data, how long metadata lives — reviewable and reproducible instead of console-only.

## Composer generations

| Aspect | Composer 2.x | Composer 3 |
|--------|--------------|------------|
| **Networking** | VPC peering or Private Service Connect | PSC via network attachment (default) |
| **Private environment** | `private_environment_config` | `enable_private_environment` flag |
| **DAG processor workload** | No | Yes (replica count capped at 3) |
| **Web server plugins mode** | No | `ENABLED` / `DISABLED` |
| **Task log storage mode** | Yes (2.0.32+) | No |
| **Airflow metadata retention** | No | Yes (30-730 days) |

The `image_version` field selects the generation (`composer-2.x.y-airflow-…` vs `composer-3-airflow-…`); the generation decides which surfaces apply. Composer 1 is a deprecated generation and is not modeled at all — see the recorded skips.

## Creation timing

Environment creation takes **25-45 minutes**: Composer assembles a GKE cluster, Cloud SQL database, and web server behind the single resource. Changes to immutable fields recreate the environment and take just as long — plan CI/CD timeouts and migration windows accordingly.

## Mutability profile

| Surface | Mutability |
|---|---|
| `region`, `environment_name` | Immutable |
| Node networking (`network`, `subnetwork`, network attachment, `ip_allocation_policy`) | Immutable |
| `private_environment_config`, `kms_key_name`, `storage_bucket` | Immutable |
| Workloads sizing, `environment_size`, `resilience_mode` | Mutable |
| Software config (packages, overrides, env vars), maintenance window, access control, retention, labels | Mutable |

## Networking models

**Composer 2.x — VPC peering**: nodes land in your subnetwork; the GKE master, Cloud SQL, and Composer internals each need a dedicated non-overlapping CIDR range (`172.16.0.0/28` is the conventional master range). Peering consumes VPC peering quota.

**Composer 2.x — Private Service Connect**: set `connection_type: PRIVATE_SERVICE_CONNECT` with a `cloud_composer_connection_subnetwork`; avoids peering quota on large networks.

**Composer 3**: `node_config.composer_network_attachment` points at a PSC network attachment; `composer_internal_ipv4_cidr_block` (a /20) covers internals. `enable_private_environment` removes public endpoints entirely.

**VPC-native ranges**: `ip_allocation_policy` pins where GKE pods and services get their IPs — per range, either the name of a secondary range your network team pre-carved on the subnetwork or a CIDR for GKE to carve itself, never both (enforced by validation). `enable_ip_masq_agent` SNATs pod traffic to node IPs when the pod CIDR is not routable in the wider network.

## Workload sizing

Each Airflow component is sized independently and updates in place:

| Component | Levers | Notes |
|-----------|--------|-------|
| Scheduler | cpu, memory, storage, count | 2 replicas for production |
| Web server | cpu, memory, storage | Always exactly 1 replica |
| Workers | cpu, memory, storage, min/max count | Autoscale with task queue depth; `max_count >= min_count` enforced |
| Triggerer | cpu, memory, count | All three required by the API when the block is present; powers deferrable operators; count 0 disables |
| DAG processor | cpu, memory, storage, count | Composer 3 only; parses DAGs independently of the scheduler |

`environment_size` (SMALL through EXTRA_LARGE) sizes the managed infrastructure underneath (GKE and database capacity); `workloads_config` sizes the Airflow components on top. They compose — size the environment for the fleet, the workloads for the DAGs.

## Security surfaces

- **CMEK** (`kms_key_name`): one key encrypts everything Composer manages — GKE node disks, Cloud SQL, the DAG bucket. The key must be in the environment's region and the Composer service agent needs `roles/cloudkms.cryptoKeyEncrypterDecrypter`. Immutable.
- **Web UI allowlist** (`web_server_network_access_control`): CIDR ranges allowed to reach the Airflow UI.
- **Control-plane allowlist** (`master_authorized_networks_config`): CIDR ranges (up to 50) allowed to reach the GKE master that runs the workloads — a distinct perimeter from the web UI.
- **Private endpoints**: `private_environment_config.enable_private_endpoint` (2.x) or `enable_private_environment` (3) removes public access entirely.

## Operational data retention

Long-lived environments accumulate task logs and metadata rows without bound unless told otherwise. `data_retention_config` provides the levers:

- `task_logs_storage_mode` — `CLOUD_LOGGING_ONLY` keeps task logs out of the environment bucket; `CLOUD_LOGGING_AND_CLOUD_STORAGE` mirrors them into it (Composer 2.0.32+, not Composer 3).
- `airflow_metadata_retention_mode` + `airflow_metadata_retention_days` — Composer 3 only; retention days must be 30-730 (validated).

## Bring-your-own DAG bucket

`storage_bucket` points the environment at an existing bucket (a `GcpGcsBucket` reference resolves to its bucket name) instead of the one Composer auto-creates — useful when bucket governance (lifecycle, IAM, CMEK, retention) is managed as its own resource. Immutable after creation.

## 90/10 coverage

The spec covers ~90% of real Composer architectures — everything from a dev sandbox to a CMEK-encrypted, private, retention-bounded enterprise fleet:

| Feature | Rationale |
|---------|-----------|
| `project_id` (optional), `region`, `environment_name` | Core identity; project rides the provider default when omitted |
| `environment_size` (incl. EXTRA_LARGE), `resilience_mode` | Capacity and availability |
| `node_config` incl. `ip_allocation_policy`, `enable_ip_masq_agent` | Node placement, identity, and VPC-native ranges |
| `software_config` incl. `cloud_data_lineage_integration` | Version pinning, packages, overrides, lineage |
| `private_environment_config` (2.x) + Composer 3 flags | Both private networking generations |
| `workloads_config` (all five components) | Per-component sizing |
| `kms_key_name` | CMEK for compliance |
| `maintenance_window`, `recovery_config` | Operational control and DR |
| `web_server_network_access_control`, `master_authorized_networks_config` | Both access perimeters |
| `data_retention_config` | Bounded operational data |
| `storage_bucket` | Governed DAG bucket |
| `labels` | User labels beneath platform attribution labels |

## Recorded skips (with reasons)

| Skipped | Reason |
|---------|--------|
| **`node_count`** | Composer-1-only field; Composer 1 is a deprecated generation. |
| **node_config `zone`, `machine_type`, `disk_size_gb`, `oauth_scopes`** | Composer-1-only node shape; Composer 2/3 size components via `workloads_config` instead. |
| **`ip_allocation_policy.use_ip_aliases`** | Composer-1-only toggle; alias IPs are always on in Composer 2/3. |
| **software_config `python_version`, `scheduler_count`** | Composer-1-only; the image version pins Python, and scheduler replicas live in `workloads_config.scheduler.count`. |
| **`database_config`** | Composer-1-only Cloud SQL machine sizing; Composer 2/3 size the database via `environment_size`. |
| **`web_server_config`** | Composer-1-only web server machine sizing; superseded by `workloads_config.web_server`. |
| **private_environment_config `web_server_ipv4_cidr_block`** | Composer-1-only private range; the 2.x private model has no separate web server range. |
| **`deletion_policy`** | Client-side Terraform lever conflicting with Planton-managed destroy (catalog-wide decision). |

## API and IAM prerequisites

- `composer.googleapis.com` enabled (both modules enable it with `disable_on_destroy=false`).
- A custom node service account must hold `roles/composer.worker`.
- For CMEK, the Composer service agent needs encrypt/decrypt on the key.

## Composes with

- **GcpCloudComposerUserWorkloadsSecret** — credentials for DAGs, delivered into this environment by reference to its `environment_name` output.
- **GcpCloudComposerUserWorkloadsConfigMap** — non-secret DAG configuration, same composition pattern.
- **GcpGcsBucket** — the governed DAG bucket or snapshot destination.
- **GcpVpcNetwork / GcpSubnetwork / GcpServiceAccount / GcpKmsKey** — upstream networking, identity, and encryption.
