# Azure HA PostgreSQL

The definitive managed PostgreSQL: zone-redundant high availability, VNet-injected private networking with no public endpoint, Microsoft Entra administration, customer-managed-key encryption on by default, an optional read replica, and the three alerts (CPU, storage, failed connections) every serious database carries. HA + private networking + CMK + Entra auth is the checklist auditors ask for — this chart delivers it as the starting point, not the aspiration.

## Who this is for

Any team that needs a relational database it can defend in a security review. The pieces are individually well-documented and collectively tedious: the VNet-injection trio (delegated subnet, private DNS zone, public access off) is CEL-enforced but easy to half-build, CMK needs a vault + key + identity + grant wired in exactly the right order, and Entra administration needs a principal. Deploy this chart and the posture is simply already correct.

## Architecture

```
                     VNet 10.30.0.0/16
   ┌───────────────────────────────────────────────┐
   │  delegated subnet 10.30.1.0/24                │      privatelink.postgres.database.azure.com
   │  (Microsoft.DBforPostgreSQL/flexibleServers)  │◀──── private DNS zone + VNet link
   │                                               │
   │  ┌──────────────┐  sync, auto-failover        │
   │  │ primary      │──────────▶ standby (zone N) │      Key Vault (purge-protected)
   │  │ (zone M)     │                             │        └─ RSA-3072 key (wrap/unwrap)
   │  │              │  async                      │             ▲ versionless reference
   │  │              │──────────▶ read replica     │             │ unwrap grant
   │  └──────┬───────┘  (replica_enabled)          │      admin identity (UAI)
   └─────────┼─────────────────────────────────────┘        ├─ Entra administrator
             │                                               └─ CMK unwrap principal
     diagnostics ──▶ Log Analytics
     alerts: cpu>80% (15m) · storage>85% · connections_failed>10 ──▶ ops email
```

Design decisions worth knowing:

- **HA and the replica are different tools.** The standby is synchronous with automatic failover (zero data loss on zone failure); the replica is asynchronous read scale-out and a warm recovery target. The chart models both and the comments say which to reach for.
- **One identity, two jobs.** The user-assigned identity is the Entra administrator (token-based admin with MFA/PIM/auditing) and, with CMK, the key-unwrap principal. Applications get their own identities and database roles — this one is for administration and infrastructure.
- **CMK's kill switch.** The key is referenced versionless (rotation propagates automatically), and revoking the identity's `Key Vault Crypto Service Encryption User` grant makes the data unreadable to Azure itself. The vault therefore carries purge protection — irreversible on that vault, and the vault name stays reserved for the soft-delete window after teardown. That trade is stated in the parameter, not discovered later.

## Resources

| Kind | Name | Purpose |
| --- | --- | --- |
| AzureResourceGroup | `{env}-ha-postgres` | One container for the estate |
| AzureLogAnalyticsWorkspace | `{env}-postgres-logs` | PostgreSQL logs (slow queries, lock waits) |
| AzureUserAssignedIdentity | `{env}-postgres-admin-identity` | Entra administrator + CMK unwrap principal |
| AzureVirtualNetwork | `{env}-postgres-vnet` | The database network |
| AzureSubnet | `{env}-postgres-subnet` | Delegated flexible-server subnet |
| AzurePrivateDnsZone + link | `{env}-postgres-dns` | Private resolution of the server FQDN |
| AzureKeyVault | `{env}-postgres-vault` | Purge-protected key home (CMK toggle) |
| AzureKeyVaultKey | `{env}-postgres-cmk` | RSA-3072 wrap/unwrap key (CMK toggle) |
| AzureRoleAssignment | `{env}-postgres-cmk-grant` | Crypto Service Encryption User grant (CMK toggle) |
| AzurePostgresqlFlexibleServer | `{env}-postgres` | The primary: HA, injected, Entra, CMK |
| AzurePostgresqlFlexibleServer | `{env}-postgres-replica` | Async read replica (replica toggle) |
| AzureMonitorActionGroup | `{env}-postgres-ops` | Alert routing |
| AzureMonitorMetricAlert | cpu / storage / connections | The three that matter |
| AzureMonitorDiagnosticSetting | `{env}-postgres-diag` | Logs + metrics into the workspace |

## Parameters

| Parameter | Description | Default | Must change |
| --- | --- | --- | --- |
| `region` | Azure region (needs availability zones for HA) | `centralus` | |
| `vnet_cidr` / `db_subnet_cidr` | Network layout | `10.30.0.0/16` / `10.30.1.0/24` | |
| `key_vault_name` | Globally unique vault name (CMK) | `myhapostgreskv` | yes |
| `postgres_admin_login` / `postgres_admin_password` | Password-auth credentials | `pgadmin` / placeholder | password: yes |
| `postgres_version` | Major version | `16` | |
| `sku_name` | Compute SKU (B-series cannot run HA/replicas) | `GP_Standard_D2s_v3` | |
| `storage_mb` | Starting storage rung (auto-grow on) | `65536` | |
| `ha_enabled` | Zone-redundant synchronous standby | `true` | |
| `backup_retention_days` | PITR window (7-35) | `14` | |
| `geo_redundant_backup_enabled` | Paired-region backup copies (fixed at creation) | `false` | |
| `replica_enabled` | Async read replica in the same subnet | `false` | |
| `cmk_enabled` | Customer-managed-key encryption | `true` | |
| `ops_email` | Alert recipient | `ops@example.com` | yes |
| `log_retention_days` | Workspace retention | `30` | |

## After deploying

1. **Connect privately** — the server has no public endpoint. From inside the VNet (or a peered one), connect to the server FQDN; the private DNS zone resolves it to the subnet address. The `fqdn` and `server_id` outputs carry the coordinates.
2. **Bring applications into the network** — peer the application VNet with `{env}-postgres-vnet` **and** link the application VNet to the private DNS zone (resolution does not cross peering by itself).
3. **Use Entra administration** — the admin identity connects with directory tokens (`az account get-access-token --resource-type oss-rdbms`); create per-application database roles from that session rather than sharing the password.
4. **Verify the failover posture** — with HA on, Azure runs planned failovers during maintenance windows transparently; trigger a forced failover in a drill before production depends on it.

## Day 2

- **Password-free** — once every client authenticates with Entra tokens, set `authentication.password_auth_enabled: false` on the server resource and remove the credentials: the CEL contract requires dropping both together.
- **Rotate the CMK** — create a new key version in the vault; the versionless reference picks it up automatically. Add a vault rotation policy to make it scheduled.
- **Promote the replica** — set `replication_role: NONE` on the replica resource to break replication and make it a standalone server (region migrations, blue/green cutovers).
- **Tune deliberately** — `server_parameters` carries the two observability defaults (`log_min_duration_statement`, `log_lock_waits`); add workload tuning (`shared_buffers`, `work_mem`) there so it is code, not portal drift.
