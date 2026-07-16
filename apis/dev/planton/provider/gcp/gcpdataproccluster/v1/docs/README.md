# GcpDataprocCluster — Research and Design Documentation

## 1. The Dataproc Landscape

Google Cloud Dataproc is the managed service for Apache Spark, Hadoop, Hive, Pig, and related open-source data processing frameworks. Google offers three execution models under the Dataproc brand:

1. **Standard clusters (GCE-based)** — Dataproc provisions and manages Compute Engine VMs as cluster nodes. The most common deployment model, and this component's `cluster_config` arm.

2. **Dataproc on GKE (virtual clusters)** — Spark workloads run as Kubernetes pods on an existing GKE cluster. This component's `virtual_cluster_config` arm: the GKE cluster and node pools are first-class catalog resources (`GcpGkeCluster`, `GcpGkeNodePool`) composed by reference.

3. **Dataproc Serverless (Batches)** — fully managed, per-job Spark execution without any cluster. A job-submission surface rather than infrastructure; a future kind candidate, out of scope here.

Both cluster models live in one kind because the API itself models them as two mutually exclusive arms of one resource (`google_dataproc_cluster` carries `cluster_config` XOR `virtual_cluster_config`). Splitting them would force the platform to police the exclusivity across two kinds; keeping the API's own shape lets a single CEL rule (`at_most_one_deployment_arm`) express it exactly. Omitting both arms is also valid — GCP creates a default GCE-based cluster with two workers.

## 2. Cluster Architecture (GCE arm)

A standard Dataproc cluster consists of:

- **Master nodes** — HDFS NameNode, YARN ResourceManager, Spark History Server, Hive Metastore. 1 master for standard mode; 3 for HA (a GCP requirement, not a free choice).
- **Primary workers** — HDFS DataNodes and YARN NodeManagers. Persistent, on-demand VMs; the stable base of the cluster.
- **Secondary workers** — optional Spot or preemptible VMs. They join YARN but carry no HDFS DataNode, so preemption never loses data — the property that makes them safe burst capacity.
- **Auxiliary DRIVER node groups** — optional dedicated capacity where Spark drivers run, so heavy drivers (large `collect()`s, many concurrent jobs) don't compete with the master's control-plane duties.

### Image versions and components

The Dataproc image version pins the whole framework stack:

| Image | Spark | Hadoop | Hive | Python |
|---|---|---|---|---|
| 2.2-debian12 | 3.5.x | 3.3.x | 3.1.x | 3.11 |
| 2.1-debian12 | 3.4.x | 3.3.x | 3.1.x | 3.10 |
| 2.0-debian11 | 3.3.x | 3.2.x | 3.1.x | 3.9 |

Optional components (JUPYTER, DOCKER, TRINO, ZEPPELIN, FLINK, HIVE_WEBHCAT, …) install per-cluster; `software_config.properties` overrides any Hadoop/Spark/YARN/Hive property using the `"prefix:property"` key form (e.g. `"spark:spark.executor.memory": "8g"`, or `"dataproc:dataproc.allow.zero.workers": "true"` for a single-node cluster).

## 3. Secondary Workers: Spot Economics and Flexibility

`preemptibility` selects the pricing model: `SPOT` (modern, dynamic pricing — recommended), `PREEMPTIBLE` (legacy fixed-discount), or `NON_PREEMPTIBLE` (on-demand secondaries, unusual). Secondary workers inherit machine configuration from the primary group unless an `instance_flexibility_policy` overrides it — and that policy is the lever that keeps large autoscaled clusters schedulable:

- **`instance_selection_list`** — ranked machine-type preferences. Dataproc provisions from the lowest-rank entry with available capacity, falling back down the list when a type's spot pool dries up.
- **`provisioning_model_mix`** — blends on-demand and spot within the group: `standard_capacity_base` workers are always on-demand, and `standard_capacity_percent_above_base` (0-100) of the capacity above the base is also on-demand. The result is a preemption-proof floor with a cheap burst ceiling.

The released provider line carries flexible provisioning only on the secondary group — masters and primary workers are deliberately not modeled with it (see recorded skips).

## 4. Autoscaling by Policy Reference

Dataproc autoscaling is not inline configuration: it's a separate, shareable `GcpDataprocAutoscalingPolicy` resource attached by URI. One policy governs many clusters; a platform team tunes scaling in one place and every attached cluster follows.

`autoscaling_policy_uri` is a `StringValueOrRef` whose default reference resolves a `GcpDataprocAutoscalingPolicy`'s `status.outputs.name` — the full resource name (`projects/{p}/locations/{l}/autoscalingPolicies/{id}`) the API expects. Attaching, swapping, and detaching the policy all update in place, so scaling strategy changes never recreate the cluster. Pair the attachment with:

- `worker_config.min_num_instances` — the autoscaler's floor on primary workers (in-place).
- `graceful_decommission_timeout` — the YARN drain window on scale-down, so running tasks finish before their node disappears.

The policy must live in the cluster's region.

## 5. Security Surfaces

### Kerberos XOR identity mapping

`security_config` accepts exactly one mechanism (`exactly_one_security_mechanism` CEL):

- **`kerberos_config`** — Hadoop Secure Mode. The critical design fact: every secret field (`root_principal_password_uri`, `kdc_db_key_uri`, keystore/truststore passwords, the cross-realm shared password) is a **Cloud Storage URI of a KMS-encrypted file**, decrypted on-cluster by the referenced `kms_key_uri`. The manifest carries paths, never secret material — the API's own contract.
- **`identity_config`** — personal cluster authentication: each user runs workloads as their own mapped service account (`user_service_account_mapping`, min 1 pair) instead of a shared cluster identity.

### VM hardening and placement

- `shielded_instance_config` — secure boot, vTPM, integrity monitoring (production baseline: all three on).
- `confidential_instance_config` — data-in-use memory encryption via AMD SEV; requires N2D machine types.
- `reservation_affinity` — pin capacity to Compute Engine reservations; `SPECIFIC_RESERVATION` requires `key` + `values` (enforced by CEL).
- `node_group_affinity` — sole-tenant node group placement for BYOL/compliance isolation.
- `encryption_kms_key_name` — CMEK on every persistent disk.

## 6. Observability and Shared Services

- **`dataproc_metric_config`** — OSS metric collection into Cloud Monitoring. Sources: `MONITORING_AGENT_DEFAULTS`, `HDFS`, `SPARK`, `YARN`, `SPARK_HISTORY_SERVER`, `HIVESERVER2` (min 1 metric entry; optional per-source `metric_overrides`).
- **`metastore_config`** — attaches the cluster to a persistent Dataproc Metastore service, so table metadata outlives ephemeral clusters. The field accepts a literal service resource name today; references attach when a metastore-service kind lands in the catalog.
- **`endpoint_config`** — the Component Gateway: authenticated HTTPS endpoints for Spark UI, YARN, HDFS, Jupyter, Zeppelin.

## 7. The Virtual Arm (Dataproc on GKE)

The virtual arm runs Dataproc's control plane and Spark workloads as pods on an existing GKE cluster:

- **`gke_cluster_target`** (required) — the target GKE cluster, referenced from `GcpGkeCluster` (resolves to the fully qualified cluster resource name).
- **`node_pool_target[]`** — maps pools to roles: `DEFAULT` (catch-all), `CONTROLLER` (Dataproc control plane), `SPARK_DRIVER`, `SPARK_EXECUTOR`. Pools may pre-exist (composed via `GcpGkeNodePool`) or be created by Dataproc using `node_pool_config` (locations required, autoscaling bounds with max >= min, preemptible XOR spot). When no targets are given, Dataproc creates and manages a default pool.
- **`kubernetes_software_config`** (required) — `component_version` with at least one pair; SPARK is the required component (e.g. `"3.5-dataproc-17"`, matched to the GKE version compatibility matrix).
- **`kubernetes_namespace`** — where workloads land; derived from the cluster name when omitted.
- **`auxiliary_services_config`** — a persistent Hive metastore and a Spark History Server hosted on another Dataproc cluster (referenced by that cluster's `cluster_id` output), so job history outlives the ephemeral workload cluster.

Two behavioral facts distinguish the arm:

1. **Fully immutable** — any change replaces the virtual cluster. The replacement does not touch the underlying GKE cluster or node pools; it re-deploys Dataproc's pods onto them.
2. **No user labels** — the Dataproc API rejects user labels on virtual clusters. The spec enforces this pre-deploy (`labels_unsupported_on_virtual_clusters`) rather than letting the API fail mid-apply.

## 8. Mutability Profile

| Surface | Mutability |
|---|---|
| `labels` | In place |
| `worker_config.num_instances` / `secondary_worker_config.num_instances` | In place (manual scaling) |
| `worker_config.min_num_instances` | In place (autoscaler floor) |
| `autoscaling_policy_uri` attach/swap/detach | In place |
| `lifecycle_config.idle_delete_ttl` / `auto_delete_time` | In place |
| `region`, `cluster_name`, everything else on the GCE arm | Recreate |
| The entire virtual arm | Recreate (GKE substrate untouched) |

The in-place set is exactly the operational surface: scale, scaling strategy, cost-control TTLs, and labels. Everything structural recreates — plan production changes around blue-green cluster swaps, which the ephemeral-cluster model makes routine.

## 9. 90/10 Coverage

### What the spec models

| Feature | Rationale |
|---------|-----------|
| `region`, `cluster_name`, `project_id` (ambient) | Core identity |
| Two deployment arms (GCE XOR GKE) | The API's own shape; one CEL rule enforces exclusivity |
| Master/worker/secondary node groups with disks, accelerators, CPU platform, image URI | The sizing surface every real cluster touches |
| `local_ssd_interface` (scsi/nvme) | Shuffle-heavy Spark tuning |
| Secondary `instance_flexibility_policy` + `provisioning_model_mix` | The spot-capacity survival levers for large autoscaled clusters |
| `cluster_tier` | Premium-tier opt-in |
| GCE networking (network XOR subnetwork, internal-only, tags, metadata, scopes, zone) | Production network placement |
| Shielded/confidential VMs, reservation and sole-tenant affinity | Security and placement hardening |
| `software_config`, init actions | Framework pinning and node bootstrap |
| `autoscaling_policy_uri` as a reference | First-class composition with GcpDataprocAutoscalingPolicy |
| CMEK (`encryption_kms_key_name`) | Compliance baseline |
| `security_config` (Kerberos XOR identity) | Both in-cluster auth models, path-based secrets |
| `endpoint_config`, `lifecycle_config`, `graceful_decommission_timeout` | The operational trio: UIs, cost control, safe scale-down |
| `metastore_config`, `dataproc_metric_config`, `auxiliary_node_groups` | Persistent metadata, observability, driver isolation |
| Virtual arm: GKE target, node-pool role mapping, namespace, component versions, auxiliary services | Complete Dataproc-on-GKE composition against first-class GKE kinds |
| `labels` | User labels beneath platform attribution labels (GCE arm only — API limitation on the virtual arm) |

### Recorded skips (with reasons)

| Skipped | Reason |
|---------|--------|
| **`cluster_type` / `engine`** | Absent from the released provider 6.x line (schema-verified). |
| **`resource_manager_tags`** on `gce_cluster_config` | Absent from the released line. |
| **Disk `boot_disk_provisioned_iops` / `boot_disk_provisioned_throughput`** | Absent from the released line. |
| **Lifecycle `idle_stop_ttl` / `auto_stop_time`** | Absent from the released line. |
| **Master/worker `instance_flexibility_policy`** | The released line carries it only on secondary workers. |
| **`deletion_policy`** | Client-side lever conflicting with Planton-managed destroy (catalog-wide skip). |
| **Dataproc IAM member/binding/policy trios** | Resource-scoped IAM deferred. |
| **`google_dataproc_job` / `batch` / `workflow_template` / `session_template`** | Workloads, not infrastructure; Serverless Batches is a future kind candidate. |
| **Dataproc Metastore service** | Future kind candidate; `metastore_config` accepts literal resource names today. |

### Design principle

Coverage targets ~90% of real Dataproc architectures composed from first-class nodes — never 100% of the API. Every skip above is deliberate and carries its reason; nothing is silently absent.

## 10. Production Best Practices

### Cost

- **Always set `idle_delete_ttl`** on non-permanent clusters — it is the single highest-leverage cost control, and it tunes in place.
- **Prefer SPOT secondaries** with an `instance_flexibility_policy`; keep a small on-demand `standard_capacity_base` for preemption-proof capacity.
- **Attach an autoscaling policy** instead of hand-scaling; use `graceful_decommission_timeout` so scale-down never kills running tasks.

### Reliability

- **3 masters for HA** where the cluster is long-lived and jobs must survive a master failure; 1 master for ephemeral batch clusters.
- **Size primary workers for HDFS**, secondaries for compute: only primaries carry DataNodes.
- **`min_num_instances`** guards the autoscaler against scaling below the HDFS-safe floor.

### Security

- **Internal IP only + custom service account** is the production baseline; add Shielded VM hardening (all three toggles) and CMEK where compliance requires.
- **Kerberos** for Hadoop Secure Mode; **identity mapping** for shared notebook clusters where per-user audit trails matter. Never both — the API accepts exactly one.

### Observability

- **Enable `dataproc_metric_config`** with SPARK and YARN sources on any cluster SREs are on call for.
- **Component Gateway** gives authenticated UI access without SSH tunnels or public IPs.

## 11. Common Pitfalls

1. **Both arms set** — `cluster_config` and `virtual_cluster_config` are mutually exclusive; the spec rejects the combination before the API does.
2. **Labels on virtual clusters** — the API rejects them; the spec catches it pre-deploy. Keep `spec.labels` empty on the virtual arm.
3. **Assuming worker-count changes recreate** — they don't; counts, the floor, the policy attachment, TTLs, and labels are the in-place set. Everything else does recreate.
4. **`SPECIFIC_RESERVATION` without key/values** — the reservation must be named via `key` (`compute.googleapis.com/reservation-name`) and `values`; validated pre-deploy.
5. **Zone-suffixed regions** — `region` takes `us-central1`, never `us-central1-a`; the pattern also admits multi-digit regions like `europe-west12`.
6. **Confidential VMs on non-N2D machines** — `confidential_instance_config` requires an N2D machine type; the API rejects other families at create time.
7. **Autoscaling policy in another region** — a cluster can only attach policies in its own region.
8. **Secondary-worker accelerators/machine types** — secondaries inherit from the primary group; per-group overrides come only through the flexibility policy's machine-type list.

## 12. Implementation Landscape

### Pulumi module

`iac/pulumi/module/` — resource `dataproc.Cluster` (pulumi-gcp v9):

- `main.go` — provider setup and orchestration
- `locals.go` — the user-beneath-platform label merge
- `dataproc_cluster.go` — both arms mapped to `dataproc.Cluster` args; exports
- `outputs.go` — output keys matching `stack_outputs.proto`

### Terraform module

`iac/tf/` — resource `google_dataproc_cluster` on the plain `google` provider (`~> 6.0`); every modeled field is GA on the released line. `software_config.properties` maps to the provider's `override_properties` (the writable surface; the provider's `properties` attribute is the computed resolved set). The staging-bucket output resolves from whichever arm is active.

Both modules export the same three outputs: `cluster_id`, `cluster_name`, `staging_bucket`.

## 13. References

- [Dataproc Documentation](https://cloud.google.com/dataproc/docs)
- [Cluster REST API](https://cloud.google.com/dataproc/docs/reference/rest/v1/projects.regions.clusters)
- [Autoscaling clusters](https://cloud.google.com/dataproc/docs/concepts/configuring-clusters/autoscaling)
- [Dataproc on GKE overview](https://cloud.google.com/dataproc/docs/guides/dpgke/dataproc-gke-overview)
- [Dataproc-on-GKE version compatibility](https://cloud.google.com/dataproc/docs/guides/dpgke/dataproc-gke-version-compatibility)
- [Secondary workers](https://cloud.google.com/dataproc/docs/concepts/compute/secondary-vms)
- [Terraform google_dataproc_cluster](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/dataproc_cluster)
- [Pulumi dataproc.Cluster](https://www.pulumi.com/registry/packages/gcp/api-docs/dataproc/cluster/)
