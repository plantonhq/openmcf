# AzureMssqlServer -- Design Research

## The Resource

An Azure SQL Database logical server (`Microsoft.Sql/servers`) is the
administrative container of the Azure SQL PaaS family: it carries the
login endpoint, authentication, networking posture, the
transparent-data-encryption protector, auditing, and Microsoft Defender
settings -- and NO compute. The component maps onto `azurerm_mssql_server`
(azurerm v4.x, `internal/services/mssql/mssql_server_resource.go`) plus
the server-scoped satellite resources folded below, parity-verified
against pulumi-azure v6 (`mssql.Server*`).

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `server_name` | Required, ForceNew, globally unique (DNS name); 1-63 chars, no edge hyphens |
| `location` | `region` | Required, ForceNew |
| `resource_group_name` | `resource_group` | FK → AzureResourceGroup |
| `version` | `version` | "2.0"/"12.0", default "12.0", ForceNew |
| `administrator_login` / `administrator_login_password` | same | AtLeastOneOf with azuread-only (mirrored as message CELs); reserved-name rule mirrored; login ForceNew once set |
| `azuread_administrator` | `azuread_administrator` | login_username/object_id/tenant_id/azuread_authentication_only; tenant falls back to the deploying credential's on both engines |
| `identity` | `identity` | SYSTEM_ASSIGNED / USER_ASSIGNED / SYSTEM_AND_USER_ASSIGNED; ids FK → AzureUserAssignedIdentity |
| `primary_user_assigned_identity_id` | same | RequiredWith identity → CEL |
| `transparent_data_encryption_key_vault_key_id` | same | VERSIONED Key Vault key (azurerm validates VersionTypeVersioned) → FK defaults to AzureKeyVaultKey.key_id, NOT versionless_id |
| `connection_policy` | `connection_policy` enum | Default/Proxy/Redirect; unspecified not sent |
| `minimum_tls_version` | same | Only "1.2" on 4.x (older floors retired); ARM rejects removal once set (comment, not CEL) |
| `public_network_access_enabled` | same | optional bool default true |
| `outbound_network_restriction_enabled` | same | Paired with the outbound rules by CEL |
| `express_vulnerability_assessment_enabled` | same | Defender's agentless scanning |
| `tags` | `tags` | User tags merged over Planton-derived tags |

Server-scoped satellite resources folded into the spec (1:1 singletons or
trivial name-only children -- none has an independent lifecycle):

| azurerm | spec |
|---|---|
| `azurerm_mssql_firewall_rule` | `firewall_rules[]` (IPv4 ranges) |
| `azurerm_mssql_virtual_network_rule` | `virtual_network_rules[]` (subnet FK + ignore_missing_vnet_service_endpoint) |
| `azurerm_mssql_outbound_firewall_rule` | `outbound_firewall_rules[]` (name-only FQDN children) |
| `azurerm_mssql_server_extended_auditing_policy` | `extended_auditing` (storage endpoint + sensitive key, retention 0-3285, log_monitoring default true, predicate, action groups) |
| `azurerm_mssql_server_security_alert_policy` | `security_alert_policy` (state enum, closed detector enum with ARM's Snake_Pascal wire vocabulary, emails, retention, storage pair CEL) |

## Decomposition Decisions

- **Databases DISSOLVED out of the server spec** (the prior shape bundled
  them). In Azure SQL's model the database is the unit of compute and
  billing with a ~40-field surface of its own, independent lifecycle,
  many-per-server, and is FK-referenced (copy/secondary sources, pool
  membership) -- a first-class kind (AzureMssqlDatabase). The server's
  `database_ids` map output dissolved with the bundle.
- **Elastic pools are a first-class kind** (AzureMssqlElasticPool):
  independent billing container, many-per-server, referenced by pooled
  databases.
- **Firewall/VNet/outbound rules and the auditing/Defender singletons
  FOLD**: 1:1 with the server or trivial name-only children, never
  FK-referenced, administered with the server.

## Recorded Skips (with reasons)

- **`administrator_login_password_wo` / `_wo_version`** -- Terraform's
  write-only-attribute ergonomic for the same password; no pulumi analog,
  a second way to say one thing.
- **`azurerm_mssql_server_dns_alias`** -- a stable CNAME spanning
  server migrations; a migration-day tool, not a durable topology node.
  Adoption backlog if migration tooling becomes a story.
- **`azurerm_mssql_server_microsoft_support_auditing_policy`** --
  audits MICROSOFT-support-operator access; niche compliance surface
  with the same shape as extended auditing. Revisit on demand.
- **`azurerm_mssql_server_vulnerability_assessment`** (classic) --
  requires a storage account + the alert policy; superseded by
  `express_vulnerability_assessment_enabled` (agentless), which IS
  modeled. azurerm itself steers to express.
- **Elastic Jobs** (`mssql_job_*`) -- a scheduling subsystem (agent,
  credentials, jobs, steps, schedules), not part of the server's 90/10
  surface. Beyond the stopping line.
- **`azurerm_mssql_failover_group`** -- a genuine composable DR node
  spanning TWO servers; deliberately deferred to its own kind (honest
  E2E needs a second server + cross-region replication). Recorded as a
  follow-up, not silently dropped.

## Design Decisions

- **The TDE CMK FK points at the VERSIONED `key_id`** -- azurerm
  validates a versioned nested-item ID at the server level (ARM pins the
  version), unlike the versionless rotation seams used elsewhere in the
  catalog. The contrast is documented on the field.
- **The at-least-one-auth contract is three CELs** mirroring azurerm's
  AtLeastOneOf/conflict matrix in both directions: no-auth rejected,
  credentials paired, Entra-only forbids credentials.
- **Two update-time contracts stay comments** (not CEL-able -- they
  compare against prior state): the TLS floor cannot be removed once
  set, and the password cannot change while Entra-only auth is on.
- **`minimum_tls_version` keeps its single-value vocabulary ("1.2")** --
  ARM retired the older floors; the field stays so the future "Disabled"
  posture (5.0 surface) lands as vocabulary widening, not a new field.

## Operational Behavior Worth Knowing

- **Logical servers create in ~1-2 minutes** -- there is no compute to
  provision. Destroys are similarly fast.
- **The server name frees on destroy** but is globally unique while
  alive.
- **`administrator_login` is immutable once set**; the password rotates
  in place (except while Entra-only auth is on).
- **Defender's alert policy propagates to every database** unless a
  database-scoped threat-detection policy overrides it.

## Composition

- `resource_group` → `AzureResourceGroup.status.outputs.resource_group_name`
- `identity.identity_ids[]` / `primary_user_assigned_identity_id` →
  `AzureUserAssignedIdentity.status.outputs.identity_id`
- `azuread_administrator.object_id` →
  `AzureUserAssignedIdentity.status.outputs.principal_id` (or a literal
  user/group directory object ID)
- `transparent_data_encryption_key_vault_key_id` →
  `AzureKeyVaultKey.status.outputs.key_id` (versioned)
- `virtual_network_rules[].subnet_id` →
  `AzureSubnet.status.outputs.subnet_id`
- `server_id` output ← referenced by `AzureMssqlDatabase.server_id`,
  `AzureMssqlElasticPool.server_id`, and
  `AzurePrivateEndpoint.private_connection_resource_id`
- `identity_principal_id` output → the `AzureRoleAssignment` seam
