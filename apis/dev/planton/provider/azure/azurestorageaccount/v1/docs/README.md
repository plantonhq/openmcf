# AzureStorageAccount -- Design Research

## The Resource

An Azure Storage Account (`Microsoft.Storage/storageAccounts`) is the
multi-service storage primitive fronting Blob, Files, Queues, Tables, and
Data Lake Storage Gen2 behind one globally-unique DNS name. The component
maps onto `azurerm_storage_account` (azurerm v4.x,
`internal/services/storage/storage_account_resource.go`) plus two folded
companions -- `azurerm_storage_management_policy` (blob lifecycle) and
`azurerm_storage_account_static_website` -- parity-verified against
pulumi-azure v6 (`storage.Account`, `storage.ManagementPolicy`,
`storage.AccountStaticWebsite`).

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `account_name` | Required, ForceNew, 3-24 lowercase alphanumerics ONLY (no hyphens), globally unique |
| `account_kind` | `account_kind` enum | Unspecified = StorageV2; only Storage→StorageV2 upgrades in place |
| `account_tier` | `account_tier` enum | azurerm requires it; the spec defaults STANDARD (the overwhelming choice), ForceNew |
| `account_replication_type` | `replication_type` enum | RA_GRS→RAGRS etc. in the modules; zonal↔non-zonal family switch is ForceNew |
| `access_tier` | `access_tier` enum | Kind-gated CEL (StorageV2/BlobStorage/FileStorage); unspecified lets Azure apply Hot |
| `provisioned_billing_model_version` | same | "V2" or unset, ForceNew |
| `https_traffic_only_enabled` | same | The v4 rename of `enable_https_traffic_only`; optional-bool default true |
| `min_tls_version` | enum | Unspecified = TLS1_2 |
| `shared_access_key_enabled` / `default_to_oauth_authentication` / `allow_nested_items_to_be_public` / `public_network_access_enabled` / `local_user_enabled` / `sftp_enabled` / `cross_tenant_replication_enabled` | same | Azure-authentic defaults (nested-public is TRUE in v4 -- documented, not silently hardened) |
| `allowed_copy_scope` | enum | AAD / PrivateLink |
| `sas_policy` | message | expiration_period `D.HH:MM:SS` pattern-validated |
| `is_hns_enabled` / `nfsv3_enabled` | same | ForceNew; NFSv3 pairing CELs mirror Create's checks (HNS + tier/kind + LRS/RA_GRS) |
| `large_file_share_enabled` | same | One-way; kind-gated CEL |
| `dns_endpoint_type` | enum | ForceNew; AzureDnsZone × restore_policy conflict CEL |
| `infrastructure_encryption_enabled` | same | ForceNew; kind/tier gate CEL |
| `queue_encryption_key_type` / `table_encryption_key_type` | enums | ForceNew; Account scope rejected on legacy Storage kind (CEL) |
| `identity` | `identity` message | VM-precedent shape (type enum + UAI FK list, ids-match-type CEL) |
| `customer_managed_key` | message | Key FK → `AzureKeyVaultKey.versionless_id` (rotation propagates); required UAI FK (azurerm's v4 contract); requires user-assigned identity (CEL) |
| `network_rules` | message | `default_action` required (azurerm's contract -- no invented default); bypass as repeated closed enum; subnet FKs → `AzureSubnet.subnet_id`; `private_link_access` |
| `blob_properties` | message | versioning/change-feed/soft-deletes/restore/CORS/last-access; restore-prereq CELs (versioning + change feed + soft delete) |
| `share_properties` | message | retention/SMB dials/CORS; kind-tier gate CEL; premium-only multichannel CEL |
| `static_website` | message | Realized via the STANDALONE `azurerm_storage_account_static_website` resource on both engines -- the inline block is deprecated for removal in azurerm v5 |
| `routing` | message | choice enum + endpoint publication flags |
| `custom_domain` | message | name + use_subdomain |
| `azure_files_authentication` | message | AADDS/AADKERB/AD; AD requires active_directory (CEL) |
| `immutability_policy` | message | ForceNew; requires versioning (CEL); one-way LOCKED state documented |
| `azurerm_storage_management_policy` | `lifecycle_rules` | FOLDED: singleton-per-account policy document (ARM name hardcoded "default"), never FK-referenced; per-destination one-basis CELs; tag-filter × snapshot/version exclusion CEL |
| `tags` | `tags` | User tags merged over Planton-derived tags |

## Decomposition Decisions

- **Blob containers are a first-class kind** (`AzureStorageContainer`):
  many-per-account, independent lifecycle, referenced by app-hosting
  deployment patterns and inventory/immutability seams. The old bundled
  `containers` list (and its `container_url_map` output) is dissolved.
- **The lifecycle (management) policy FOLDS in**: ARM models it as ONE
  per-account policy document with no independent lifecycle and nothing
  referencing it -- the Postgres firewall-rules fold verdict.
- **Deferred to a storage-data-services session** (each a real standalone
  Azure resource that passes the split test): file shares, queues,
  tables, encryption scopes (referenced by container
  `default_encryption_scope` -- until then the container models the scope
  as a plain value), ADLS Gen2 filesystems, SFTP local users, object
  replication (spans two accounts), and blob inventory policies
  (references a container by name).

## Recorded Skips (with reasons)

- **Inline `queue_properties` (classic analytics metrics/logging + queue
  CORS)** -- legacy Storage-Analytics surface superseded by Azure Monitor;
  the inline block is deprecated for removal in azurerm v5. Queue-service
  settings arrive with the queue kind in the data-services session.
- **`managed_hsm_key_id`** (CMK from a Managed HSM pool) -- deprecated in
  v4 in favour of `key_vault_key_id`; no Managed HSM kind exists to
  reference. Revisit if a Managed HSM kind is ever forged.
- **`edge_zone`** is modeled but exercised nowhere -- Edge Zones need
  special subscription enrollment; recorded as untested surface.
- **Internet/Microsoft routing endpoint VARIANT outputs** (the ×3
  endpoint matrix) -- niche; the primary/secondary endpoint set covers
  the real consumers. Additive if a routing-variant consumer appears.
- **`immutability_policy` state transitions** (Unlocked→Locked one-way,
  Locked is ForceNew) are provider/ARM-enforced at apply time -- CEL
  validates single manifests, not transitions between applies.

## Design Decisions

- **Explicit `account_name`** (the Key Vault `vault_name` precedent):
  globally-unique DNS-visible names must be deliberate, not derived from
  metadata by silent munging (the old module stripped hyphens and
  truncated -- two accounts could collide).
- **`allow_nested_items_to_be_public` keeps Azure's real default (true)**
  with the security posture documented on the field, rather than
  inventing a hardened default that would surprise anyone comparing
  against Azure's own behavior. The account-level flag only PERMITS
  per-container opt-in.
- **`network_rules.default_action` is required** when the block is
  present -- azurerm's real contract; the old spec's invented
  always-DENY default silently locked accounts down.
- **Static website via the standalone resource on BOTH engines**: the
  spec models it as a block (the user-facing shape is identical), the
  modules realize it via `storage_account_static_website` -- building on
  the deprecated inline block would be tech debt with a published
  removal date.
- **Access keys and connection strings ARE exported** (the ACR
  `admin_password` / AWS IAM secret-key output precedent, prose-documented
  as secret-bearing): Function App and Linux Web App storage bindings
  genuinely consume the key, and their FK defaults now point at
  `primary_access_key`. Entra-based auth stays the documented preference.

## Operational Behavior Worth Knowing

- **Accounts create in ~30-60 seconds**; deletes are fast but Azure
  briefly reserves a just-deleted name -- never destroy-and-recreate the
  same account name in quick succession.
- **ARM control-plane operations bypass `network_rules`** -- they govern
  the data plane only. A DENY-firewalled account still verifies, plans,
  and destroys normally.
- **Changing between zonal (ZRS/GZRS/RA_GZRS) and non-zonal (LRS/GRS/
  RA_GRS) replication families replaces the account**; within-family
  changes are live migrations.
- **CMK requires the identity grant BEFORE account creation**: the
  user-assigned identity needs wrap/unwrap on the key's vault (Key Vault
  Crypto Service Encryption User), and the vault must have purge
  protection enabled.

## Composition

- `resource_group` → `AzureResourceGroup.resource_group_name`
- `network_rules.virtual_network_subnet_ids[]` → `AzureSubnet.subnet_id`
- `customer_managed_key.key_vault_key_id` → `AzureKeyVaultKey.versionless_id`
- `customer_managed_key.user_assigned_identity_id` / `identity.identity_ids[]`
  → `AzureUserAssignedIdentity.identity_id`
- `storage_account_id` output ← referenced by
  `AzureStorageContainer.storage_account_id`
- `storage_account_name` output ← `AzureFunctionApp.storage_account_name`
- `primary_access_key` output ← `AzureFunctionApp.storage_account_access_key`,
  `AzureLinuxWebApp.storage_mounts[].access_key`
- `primary_blob_host` / `primary_web_host` outputs ← CDN / Front Door
  origins and custom-domain CNAMEs
