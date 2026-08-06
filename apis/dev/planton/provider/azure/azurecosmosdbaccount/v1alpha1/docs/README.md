# AzureCosmosdbAccount -- Design Research

## The Resource

An Azure Cosmos DB account (`Microsoft.DocumentDB/databaseAccounts`) is
the governance boundary of Azure's globally distributed, multi-model
database: regions, consistency, network posture, encryption, backup,
identity, and billing guardrails are all account-level decisions. The
component maps onto `azurerm_cosmosdb_account` (azurerm v4.x,
`internal/services/cosmos/cosmosdb_account_resource.go`),
parity-verified field-by-field against pulumi-azure v6
(`cosmosdb.Account`).

The account deliberately contains NO databases or containers: those are
standalone ARM child resources with independent lifecycles, throughput,
and billing, modeled as the first-class kinds AzureCosmosdbSqlDatabase,
AzureCosmosdbSqlContainer, AzureCosmosdbMongoDatabase, and
AzureCosmosdbMongoCollection referencing this account.

## Field Mapping (azurerm -> spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `account_name` | Required, ForceNew, `^[-a-z0-9]{3,50}$` (globally unique DNS label) |
| `location` | `region` | Required, ForceNew |
| `resource_group_name` | `resource_group` | FK -> AzureResourceGroup, ForceNew |
| `offer_type` | not modeled | ARM accepts only "Standard" -- nothing to choose; hardcoded in both modules |
| `kind` | enum | GLOBAL_DOCUMENT_DB (default) / MONGO_DB; ForceNew. `Parse` skipped (legacy kind for the retired Parse Server platform) |
| `consistency_policy` | required message | Level enum defaulting SESSION; BoundedStaleness dials 5-86400 / >=10; the multi-region floors (>=300 / >=100000) are spec-level rules mirroring the provider's create-time check |
| `geo_location` | `geo_locations` | Exactly-one-priority-0, unique priorities, unique locations as spec rules (the provider enforces them in expand) |
| `capabilities` | closed enum list | All 20 azurerm values; kind-pairing rules mirror the provider's capability/kind map; MongoDBv3.4-requires-EnableMongo mirrored; in-place add/remove subset documented on the field |
| `free_tier_enabled` | optional bool | ForceNew; one free-tier account per subscription |
| `automatic_failover_enabled` / `multiple_write_locations_enabled` / `public_network_access_enabled` / `is_virtual_network_filter_enabled` | optional bools | Azure defaults preserved (true only for public network access) |
| `virtual_network_rule` | `virtual_network_rules` | subnet FK -> AzureSubnet + `ignore_missing_vnet_service_endpoint` |
| `ip_range_filter` | repeated string | v4's set-of-strings shape; per-item IPv4/CIDR pattern rule |
| `backup` | message | Type enum + Continuous tier XOR Periodic dials as message rules (the provider's contract); PERIODIC->CONTINUOUS in place, reverse ForceNew |
| `mongo_server_version` | enum | ARM's full vocabulary incl. 3.2 |
| `identity` | message | SystemAssigned/UserAssigned/both + UAI FKs; ids-match-type rule |
| `default_identity_type` | `default_identity` message | The provider models a composite string ("UserAssignedIdentity=<id>"); the spec decomposes it into a type enum + a real AzureUserAssignedIdentity reference and the modules compose the wire string -- referenceable composition instead of a magic string |
| `key_vault_key_id` | FK -> AzureKeyVaultKey.versionless_id | ForceNew; versionless so rotation propagates |
| `analytical_storage_enabled` + `analytical_storage` | optional bool + message | One-way: disabling recreates (provider CustomizeDiff) |
| `capacity` | message | `total_throughput_limit` >= -1 |
| `access_key_metadata_writes_enabled` / `local_authentication_enabled` | optional bools (default true) | The spec models the non-deprecated `*_enabled` forms; the TF module uses azurerm's `local_authentication_enabled` directly, while pulumi-azure v6.38 has bridged only the deprecated `localAuthenticationDisabled` form, so the Pulumi module sends the negation (PARITY-EXCEPTION documented in both modules; the wire property and the created account are identical) |
| `minimal_tls_version` | enum (TLS_1_0/TLS_1_1/TLS_1_2) | azurerm v4's running schema accepts all three wire values (`Tls`/`Tls11`/`Tls12`); unset materializes `Tls12` -- Azure's own default since April 2023 -- on both engines. The provider's v5 line restricts the vocabulary to `Tls12` only |
| `network_acl_bypass_for_azure_services` / `network_acl_bypass_ids` | bool + repeated string | Bypass ids stay plain ARM ids -- the list admits many unrelated resource types |
| `burst_capacity_enabled` / `partition_merge_enabled` | optional bools | |
| `cors_rule` | message | One rule (provider MaxItems 1); method vocabulary as a per-item rule |
| `create_mode` + `restore` | enum + message | All ForceNew; create_mode-requires-CONTINUOUS and RESTORE<->restore pairing as spec rules (the provider's contracts) |
| `tags` | `tags` | User tags merged over the platform's identity tags (user wins) |

## Decomposition Decisions

- **Databases/containers are first-class kinds, not folds**: standalone
  ARM child resources, many-per-account, independent lifecycle and
  throughput billing, FK-referenced (containers reference databases).
- **The default identity decomposes**: the provider's
  "UserAssignedIdentity=<resource id>" composite string becomes a type
  enum plus a real identity reference, so composed CMK deployments wire
  the unwrapping identity by reference.
- **Capabilities are declared, never injected**: a MONGO_DB account
  declares ENABLE_MONGO itself. Capabilities are part of the account's
  real, mostly-immutable configuration; a module silently adding one
  would hide exactly the kind of state that recreates accounts.

## Recorded Skips (with reasons)

- **`managed_hsm_key_id`** -- deprecated in favor of `key_vault_key_id`
  (removed in the provider's v5); not modeled.
- **`local_authentication_disabled`** -- the deprecated negative form;
  the spec models `local_authentication_enabled`.
- **`Parse` kind** -- the retired Parse Server platform's legacy kind.
- **Cassandra/Gremlin/Table child resources** -- the capabilities are
  modeled (an advanced org can enable the APIs); dedicated child kinds
  wait for chart/adoption demand.
- **`azurerm_cosmosdb_sql_dedicated_gateway`** -- a real standalone
  resource (provisioned query gateway); evaluated for its own kind on
  demand.
- **SQL data-plane RBAC** (`sql_role_definition`/`sql_role_assignment`)
  and **Mongo RBAC** (`mongo_role_definition`/`mongo_user_definition`)
  -- the Entra keyless data-plane story; committed follow-up kinds.
- **SQL triggers/stored procedures/functions** -- JavaScript code
  artifacts (content, not infrastructure); applications own them.

## Operational Behavior Worth Knowing

- **Capability changes usually recreate the account.** Azure allows only
  a small set to be ADDED in place (DeleteAllItemsByPartitionKey,
  DisableRateLimitingResponses, AllowSelfServeUpgradeToMongo36,
  EnableAggregationPipeline, MongoDBv3.4, mongoEnableDocLevelTTL, the
  EnableMongo* feature switches, EnableFabricNetworkAclBypass) and only
  EnableMongoRetryableWrites / DisableRateLimitingResponses to be
  REMOVED in place. Settle capabilities before production.
- **Analytical storage is one-way**: enabling updates in place,
  disabling recreates.
- **Backup direction matters**: PERIODIC -> CONTINUOUS upgrades in
  place; CONTINUOUS -> PERIODIC recreates.
- **Region changes are in-place**: adding/removing `geo_locations` and
  re-prioritizing failover are updates -- how accounts grow their
  footprint.
- **Restores come from the RESTORABLE account id** (the
  locations/{location}/restorableDatabaseAccounts/{instanceId} form),
  not the plain account id -- list them with
  `az cosmosdb restorable-database-account list`.
- **Provisioning is slow**: accounts take 5-10+ minutes to create and
  delete (the provider's own timeouts run to 3 hours).
- **The keys stop working when local auth is off** -- by design; pair
  `local_authentication_enabled: false` with
  `access_key_metadata_writes_enabled: false` for a fully
  Entra-governed account.

## Composition

- `resource_group` -> `AzureResourceGroup.status.outputs.resource_group_name`
- `virtual_network_rules[].subnet_id` -> `AzureSubnet.status.outputs.subnet_id`
- `key_vault_key_id` -> `AzureKeyVaultKey.status.outputs.versionless_id`
- `identity.identity_ids[]` / `default_identity.user_assigned_identity_id`
  -> `AzureUserAssignedIdentity.status.outputs.identity_id`
- `cosmosdb_account_id` output <- AzureCosmosdbSqlDatabase /
  AzureCosmosdbMongoDatabase parents; AzurePrivateEndpoint's
  `private_connection_resource_id` (subresource "Sql" / "MongoDB")
