# GcpCloudSql — Design & Research

## What this component is

`GcpCloudSql` models one Cloud SQL instance (`google_sql_database_instance`): a fully managed MySQL, PostgreSQL, or SQL Server database server. One resource is one instance — a primary, or a read replica of another instance. Everything *inside* the instance that has its own lifecycle is its own composable kind: logical databases are `GcpCloudSqlDatabase`, users are `GcpCloudSqlUser`, both referencing the instance by its `instance_name` output.

The component targets the 90/10 line: roughly 90% of real Cloud SQL architectures — private-IP production databases, HA with automatic failover, PITR-protected data, IAM-authenticated access, read-scaling via replicas, CMEK compliance — composable from first-class nodes, without mirroring 100% of the API's long tail.

## Decomposition rationale

The split test (independent lifecycle / FK-referenced / many-per-parent) lands cleanly:

- **Databases** — many per instance, created and dropped constantly as applications come and go, referenced by nothing but their instance. Split.
- **Users** — many per instance, per-application credential lifecycles (rotation!), IAM-type users tie into the workload-identity story. Split.
- **Read replicas** — a replica IS a `google_sql_database_instance` with `master_instance_name` set; the honest model is a replica arm on this kind (each replica is its own `GcpCloudSql` node referencing the primary), NOT a separate kind. Promotion re-shapes an existing instance and is an operational action outside IaC.
- **SSL client certs** (`google_sql_ssl_cert`) — deferred; mutual-TLS client certs are a niche direct-connection pattern, and the Auth Proxy path (which is always TLS + IAM) covers the mainstream. Revisit on concrete pull.

## Connectivity model (the part people get wrong)

Cloud SQL has three connectivity paths, and the spec makes each explicit:

1. **Public IPv4** (`network.ipv4Enabled`) — with an empty `authorizedNetworks` list, a public IP is reachable ONLY through the Cloud SQL Auth Proxy or language connectors, which authenticate via IAM and encrypt via TLS. This is a safe, zero-VPC-plumbing pattern, and it is the component's default when the network block is omitted.
2. **Private VPC IP** (`network.privateNetwork`) — requires *private services access* on the VPC first: a reserved range (`GcpGlobalAddress`, purpose `VPC_PEERING`) plus the peering (`GcpServiceNetworkingConnection`). GCP prechecks this and fails the create otherwise — deliberately, because a failed instance create still burns the reserved instance name. Setting a network enables private IP; the network can be changed but never removed in place.
3. **Private Service Connect** (`network.psc`) — exposes the instance as a PSC service attachment that consumer VPCs connect to via PSC endpoints; private connectivity without VPC peering. The `psc_service_attachment_link` output is the composition key for consumer-side endpoints.

TLS posture on direct connections is `ssl_mode` (up to mutual TLS); `connectorEnforcement: REQUIRED` rejects direct connections entirely. The server CA hierarchy is selectable (per-instance CA, Google-managed CAS, or customer-managed CAS pool with custom SANs).

## Engine-coherence validation

The spec validates cross-field rules GCP would otherwise reject mid-deploy (an expensive failure — instance creates run ~10 minutes and burn the name):

- `databaseVersion` prefix ⇔ `databaseEngine`; SQL Server ⇒ `rootPassword`.
- PITR (`pointInTimeRecoveryEnabled`) is PostgreSQL/SQL Server; MySQL's PITR mechanism is `binaryLogEnabled` — modeled as two fields with cross-engine CEL, because conflating them (one "pitr" flag) would silently do the wrong thing on one engine.
- `REGIONAL` availability requires backups (and binary logs on MySQL).
- The SQL Server-only surface (`timeZone`, `collation`, `threadsPerCore`, audit config, Active Directory domain) is fenced by engine.
- Data cache requires `ENTERPRISE_PLUS`.
- Network coherence: at least one connectivity path; allowlists require the public IP; allocated ranges and private-path require the private network; customer CAS pool pairing.

## Lifecycle semantics worth knowing

- **Name reservation**: a deleted instance's name is reserved ~1 week. E2E and ephemeral environments must use per-run names.
- **In place**: `databaseVersion` (major upgrades!), `tier`, `edition`, availability, disk growth, flags, maintenance, insights.
- **ForceNew**: name, region, CMEK key, disk type, PSC enablement, removing the private network.
- **Never shrink**: disks grow (manually or auto-resize) but never shrink.
- **Two delete guards**: the engine-side `deletionProtection` refuses destroys at plan time; the API-side `deletionProtectionEnabled` blocks deletion from every surface including the console. Production instances set both, plus `retainBackupsOnDelete` as the last-resort recovery path.

## Deliberately unmodeled (with reasons)

- **Restore/clone paths** (`clone`, `restore_backup_context`, `point_in_time_restore_context`, BackupDR restores): operational actions, not steady-state configuration. A declarative spec that embeds "restore from backup X" re-runs the restore on every drift-correction — a foot-gun, not a feature.
- **Read pools** (instance_type/node_count/auto-scale) and several newer fields (`entraid_config`, `data_api_access`, write-only root password, final-backup config, MySQL auto-upgrade, Hyperdisk IOPS dials): not on the released 6.x provider line the modules pin (`~> 6.0`). Recorded here so the next depth pass re-evaluates them against the then-current release.
- **`replication_cluster` (DR pairs)**: Enterprise Plus switchover choreography; revisit on pull.
- **`maintenance_version` pinning** and **`pricing_plan`** (single-valued): drift generators with no configuration value.

## Composition map

- `network.privateNetwork` ← `GcpVpcNetwork.status.outputs.network_id` (the `projects/.../global/networks/...` form the API canonically accepts).
- `network.allocatedIpRange` ← a `GcpGlobalAddress.name` (VPC_PEERING range) when pinning which range the private IP comes from.
- `encryptionKeyName` ← `GcpKmsKey.status.outputs.key_id`.
- `masterInstanceName` ← another `GcpCloudSql.status.outputs.instance_name`.
- Downstream: `GcpCloudSqlDatabase.instance` / `GcpCloudSqlUser.instance` ← `instance_name`; `GcpCloudRun`'s Cloud SQL wiring ← `connection_name`; GCS grants for import/export ← `service_account_email`.

## Operational guidance

- Prefer the Auth Proxy / connectors over IP allowlists; never ship `0.0.0.0/0`.
- Enable Query Insights in production — the overhead is negligible and it is the first tool anyone reaches for during an incident.
- Take a backup before major version upgrades (in place ≠ risk-free).
- For IAM database auth on PostgreSQL, set `databaseFlags: {"cloudsql.iam_authentication": "on"}` on the instance before creating IAM-type users.
- Replicas inherit users from the primary; users are instance-scoped resources on the primary only.
