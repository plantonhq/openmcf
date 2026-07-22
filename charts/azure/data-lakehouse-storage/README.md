# Azure Data Lakehouse Storage

A governed ADLS Gen2 foundation: the medallion zones (raw, curated, sandbox) as first-class filesystems with POSIX inheritance posture, per-zone encryption scopes under a customer-managed key by default, least-privilege data-plane RBAC scoped to each zone, lifecycle tiering that pays hot rates only where the lake earns them, and the per-operation data-access audit trail a governed lake is defined by. Keyless throughout — the two workload identities are the only data-plane principals.

## Who this is for

A data platform team standing up the storage layer under Spark, Databricks, Synapse, or dbt. The medallion layout is a whiteboard sketch everyone agrees on and a permissions model few wire correctly: zone-scoped grants (not account-wide), ACL inheritance that propagates without per-object management, crypto separation an auditor can verify, and the answer to "who read this data." Deploy this chart and the lake starts governed instead of growing toward governance.

## Architecture

```
   engineering UAI ──── Storage Blob Data Contributor ──▶ raw, curated, sandbox
   analysts UAI ─────── Storage Blob Data Reader ───────▶ curated (only)
                        Storage Blob Data Contributor ──▶ sandbox
                                    │ (grants scope to each zone's
                                    │  container-proxy ARM id)
   ┌────────────────────────────────▼──────────────────────────────┐
   │  HNS storage account (ZRS, OAuth-default, soft delete)        │
   │  ├─ raw/      scope "rawzone"      eng rwx · others ---      │
   │  ├─ curated/  scope "curatedzone"  eng rwx · analysts r-x    │
   │  └─ sandbox/  scope "sandboxzone"  both rwx · TTL-expired    │
   │  lifecycle: raw hot→cool({raw_cool_after_days}d)→archive     │
   └───────┬───────────────────────────────────────────────────────┘
           │ blob-service diagnostics (StorageRead/Write/Delete)
           ▼                          Key Vault (purge-protected)
     Log Analytics                      └─ RSA-3072 zone key ◀─ account identity
                                           (cmk_enabled, versionless)   unwrap grant
```

Design decisions worth knowing:

- **Zones are filesystems, and filesystems are the grant boundary.** Each `AzureRoleAssignment` scopes to a zone's container-proxy ARM id (the `filesystem_id` output) — "analysts read curated and nothing else" is architecture, not convention. Azure evaluates these RBAC grants first; the POSIX ACLs beneath are the fine-grained layer for principals without a zone-wide role, and the DEFAULT entries make new data inherit the posture automatically.
- **Identities in ACLs are wired by reference.** Every named ACL entry and every grant references the workload identity's `principal_id` output — the object id ACLs actually match on (a client id in an ACL silently never matches, a classic lake-debugging story the templates spell out).
- **The data-protection story is stated honestly.** HNS accounts cannot carry blob versioning or point-in-time restore (Azure's platform contract). What the lake gets: blob and container soft delete (the recycle bin), lifecycle economics as policy, and immutable-landing discipline taught in the comments — not a false sense of a feature that is not there.
- **The audit trail targets the blob service, not the account.** Per-operation read/write/delete logs live on the blob service sub-resource; the account-level id exposes only account metrics. The diagnostic setting references the account's `blob_service_id` output — deploy the chart and "who read this data" is a Log Analytics query.
- **Per-zone encryption scopes make crypto separation a fact.** Every zone's data encrypts under its own named scope; with `cmk_enabled` (the default) all three unwrap a purge-protected, versionless-referenced vault key through the account's identity. Rotation propagates automatically.

## Resources

| Kind | Name | Purpose |
| --- | --- | --- |
| AzureResourceGroup | `{env}-lakehouse` | One container for the estate |
| AzureLogAnalyticsWorkspace | `{env}-lakehouse-logs` | The data-access audit trail |
| AzureUserAssignedIdentity | `{env}-lake-engineering` | Pipelines and processing jobs |
| AzureUserAssignedIdentity | `{env}-lake-analysts` | BI tools and notebooks |
| AzureStorageAccount | `{env}-lakehouse-store` | HNS, ZRS, soft delete, lifecycle rules |
| AzureMonitorDiagnosticSetting | `{env}-lakehouse-audit` | Blob-service read/write/delete logs → workspace |
| AzureKeyVault | `{env}-lakehouse-vault` | Purge-protected key home (CMK toggle) |
| AzureKeyVaultKey | `{env}-lakehouse-cmk` | RSA-3072 zone key (CMK toggle) |
| AzureRoleAssignment | `{env}-lakehouse-cmk-grant` | Account identity → Crypto Service Encryption User (CMK toggle) |
| AzureStorageEncryptionScope | rawzone / curatedzone / sandboxzone | Per-zone crypto separation |
| AzureStorageDataLakeGen2Filesystem | raw / curated / sandbox | The medallion zones with POSIX ACLs |
| AzureRoleAssignment | 5 zone grants | Contributor/Reader at container-proxy scopes |

## Parameters

| Parameter | Description | Default | Must change |
| --- | --- | --- | --- |
| `region` | Azure region | `centralus` | |
| `storage_account_name` | Globally unique account name (no hyphens) | `mylakehousestore` | yes |
| `key_vault_name` | Globally unique vault name (CMK) | `mylakehousekv` | yes |
| `soft_delete_days` | Blob + container recycle-bin window | `14` | |
| `raw_cool_after_days` | Raw zone hot→cool (days since modification) | `30` | |
| `raw_archive_after_days` | Raw zone cool→archive | `90` | |
| `sandbox_ttl_days` | Sandbox auto-delete window | `30` | |
| `cmk_enabled` | Customer-managed key for the zone scopes | `true` | |
| `log_retention_days` | Audit-trail retention | `30` | |

## After deploying

1. **Bind the identities** — attach `{env}-lake-engineering` to compute that runs pipelines (Databricks access connector, Synapse workspace identity, AKS workload identity) and `{env}-lake-analysts` to BI/notebook compute. The identities' `client_id` outputs are what the tools configure.
2. **Address the zones** — `abfss://raw@{account}.dfs.core.windows.net/`, and the same for `curated` and `sandbox`. The account's `primary_dfs_endpoint` output carries the host.
3. **Land data immutably** — write raw landings under date-partitioned paths (`raw/source=x/date=.../`) and never rewrite them; the lifecycle rules assume raw is history, and HNS has no versioning to save a rewrite gone wrong.
4. **Query the audit trail** — the `StorageBlobLogs` table in the workspace answers who read, wrote, or deleted what, with caller identity and path.

## Day 2

- **Harden to Entra-only** — once every consumer authenticates with tokens, set `sharedAccessKeyEnabled: false` on the account. It stays on at deploy time because ADLS filesystem provisioning is a data-plane operation that authenticates with the account key; flip it after the zones exist (re-running the chart's filesystem resources afterwards would need it back on).
- **More teams** — add an identity + zone-scoped grants per team; the pattern in `access.yaml` (Contributor for producers, Reader for consumers, per-zone) is the template. Entra security groups work in the ACLs too — pass the group's object id as a literal.
- **Finer-grained directories** — the chart owns each zone's ROOT ACL; per-project directories inside a zone get their own ACLs via SDKs or Storage Explorer, inheriting the DEFAULT entries automatically.
- **Rotate the CMK** — create a new key version; the versionless reference picks it up. Add a vault rotation policy to make it scheduled.
- **Network hardening** — for a private-only lake, add `networkRules` (default DENY + allowed subnets) or a private endpoint per the Cosmos chart's pattern; analytics engines then need VNet paths (private endpoints or service endpoints) to reach the dfs endpoint.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
