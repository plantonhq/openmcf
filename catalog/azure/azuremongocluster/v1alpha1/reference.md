# AzureMongoCluster

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Deep-shape example for docs and offline validation: a Default-mode
# production-ish cluster with zone-redundant HA, Entra + native auth,
# and two firewall rules. The administrator password is a literal here
# only so the manifest validates standalone -- reference a secret
# store in real use.
apiVersion: azure.planton.dev/v1alpha1
kind: AzureMongoCluster
metadata:
  name: test-mongo-cluster
  id: test-mongo-cluster
  org: test-org
  env: test
spec:
  resourceGroup:
    value: test-rg
  name: acme-orders-db
  region: eastus
  createMode: Default
  administratorUsername: mongoadmin
  administratorPassword:
    value: Test-Passw0rd-Change-Me
  version: "8.0"
  computeTier: M30
  storageSizeInGb: 128
  storageType: PremiumSSD
  shardCount: 1
  highAvailabilityMode: ZoneRedundantPreferred
  authenticationMethods:
    - NativeAuth
    - MicrosoftEntraID
  publicNetworkAccessEnabled: true
  firewallRules:
    - name: office
      startIpAddress: 203.0.113.0
      endIpAddress: 203.0.113.255
    - name: vpn-egress
      startIpAddress: 198.51.100.7
      endIpAddress: 198.51.100.7
  tags:
    costCenter: platform
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.region` | `string` | yes |  |  |
| `spec.createMode` | `string` |  | `Default` |  |
| `spec.administratorUsername` | `string` |  |  |  |
| `spec.administratorPassword` | `string \| valueFrom` (sensitive) |  |  |  |
| `spec.version` | `string` |  |  |  |
| `spec.computeTier` | `string` |  |  |  |
| `spec.storageSizeInGb` | `int32` |  |  |  |
| `spec.storageType` | `string` |  | `PremiumSSD` |  |
| `spec.shardCount` | `int32` |  |  |  |
| `spec.highAvailabilityMode` | `string` |  |  |  |
| `spec.authenticationMethods` | `[]string` |  |  |  |
| `spec.userAssignedIdentityIds` | `[]string \| valueFrom` |  |  | AzureUserAssignedIdentity (`status.outputs.identity_id`) |
| `spec.customerManagedKey` | `AzureMongoClusterCustomerManagedKey` |  |  |  |
| `spec.customerManagedKey.keyVaultKeyId` | `string \| valueFrom` | yes |  | AzureKeyVaultKey (`status.outputs.versionless_id`) |
| `spec.customerManagedKey.userAssignedIdentityId` | `string \| valueFrom` | yes |  | AzureUserAssignedIdentity (`status.outputs.identity_id`) |
| `spec.previewFeatures` | `[]string` |  |  |  |
| `spec.sourceServerId` | `string \| valueFrom` |  |  | AzureMongoCluster (`status.outputs.mongo_cluster_id`) |
| `spec.sourceLocation` | `string` |  |  |  |
| `spec.restore` | `AzureMongoClusterRestore` |  |  |  |
| `spec.restore.pointInTimeUtc` | `string` | yes |  |  |
| `spec.restore.sourceId` | `string \| valueFrom` | yes |  | AzureMongoCluster (`status.outputs.mongo_cluster_id`) |
| `spec.dataApiModeEnabled` | `bool` |  |  |  |
| `spec.publicNetworkAccessEnabled` | `bool` |  | `true` |  |
| `spec.firewallRules` | `[]AzureMongoClusterFirewallRule` |  |  |  |
| `spec.firewallRules[].name` | `string` | yes |  |  |
| `spec.firewallRules[].startIpAddress` | `string` | yes |  |  |
| `spec.firewallRules[].endIpAddress` | `string` | yes |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.resourceGroup

`string | valueFrom` · required

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.name

`string` · required

- rule: Cluster names must be 3-40 characters of lowercase letters, numbers, and hyphens, starting and ending with a letter or number
- rule: {"required":true}

### spec.region

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.createMode

`string` · optional (explicit presence)

- default: `Default`
- rule: {"string":{"in":["Default","GeoReplica","PointInTimeRestore"]}}

### spec.administratorUsername

`string`

### spec.administratorPassword

`string | valueFrom` · sensitive

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.version

`string` · optional (explicit presence)

- rule: {"string":{"in":["5.0","6.0","7.0","8.0"]}}

### spec.computeTier

`string` · optional (explicit presence)

- rule: {"string":{"in":["Free","M10","M20","M25","M30","M40","M50","M60","M80","M200"]}}

### spec.storageSizeInGb

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":32768,"gte":32}}

### spec.storageType

`string` · optional (explicit presence)

- default: `PremiumSSD`
- rule: {"string":{"in":["PremiumSSD","PremiumSSDv2"]}}

### spec.shardCount

`int32` · optional (explicit presence)

- rule: {"int32":{"gte":1}}

### spec.highAvailabilityMode

`string` · optional (explicit presence)

- rule: {"string":{"in":["Disabled","ZoneRedundantPreferred"]}}

### spec.authenticationMethods

`[]string`

- rule: {"repeated":{"unique":true,"items":{"string":{"in":["NativeAuth","MicrosoftEntraID"]}}}}

### spec.userAssignedIdentityIds

`[]string | valueFrom`

- references: AzureUserAssignedIdentity (`status.outputs.identity_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureUserAssignedIdentity, name: <that resource's name>, fieldPath: status.outputs.identity_id}} -- a bare string does not parse

### spec.customerManagedKey

`AzureMongoClusterCustomerManagedKey`

### spec.customerManagedKey.keyVaultKeyId

`string | valueFrom` · required

- references: AzureKeyVaultKey (`status.outputs.versionless_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureKeyVaultKey, name: <that resource's name>, fieldPath: status.outputs.versionless_id}} -- a bare string does not parse

### spec.customerManagedKey.userAssignedIdentityId

`string | valueFrom` · required

- references: AzureUserAssignedIdentity (`status.outputs.identity_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureUserAssignedIdentity, name: <that resource's name>, fieldPath: status.outputs.identity_id}} -- a bare string does not parse

### spec.previewFeatures

`[]string`

- rule: {"repeated":{"unique":true,"items":{"string":{"minLen":"1"}}}}

### spec.sourceServerId

`string | valueFrom`

- references: AzureMongoCluster (`status.outputs.mongo_cluster_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureMongoCluster, name: <that resource's name>, fieldPath: status.outputs.mongo_cluster_id}} -- a bare string does not parse

### spec.sourceLocation

`string`

### spec.restore

`AzureMongoClusterRestore`

### spec.restore.pointInTimeUtc

`string` · required

- rule: point_in_time_utc must be an RFC 3339 timestamp like 2026-08-01T12:00:00Z
- rule: {"required":true}

### spec.restore.sourceId

`string | valueFrom` · required

- references: AzureMongoCluster (`status.outputs.mongo_cluster_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureMongoCluster, name: <that resource's name>, fieldPath: status.outputs.mongo_cluster_id}} -- a bare string does not parse

### spec.dataApiModeEnabled

`bool` · optional (explicit presence)

### spec.publicNetworkAccessEnabled

`bool` · optional (explicit presence)

- default: `true`

### spec.firewallRules

`[]AzureMongoClusterFirewallRule`

### spec.firewallRules[].name

`string` · required

- rule: Firewall rule names must be 1-80 characters, start with a letter or number, and contain only letters, numbers, dots, hyphens, and underscores
- rule: {"required":true}

### spec.firewallRules[].startIpAddress

`string` · required

- rule: {"required":true,"string":{"ipv4":true}}

### spec.firewallRules[].endIpAddress

`string` · required

- rule: {"required":true,"string":{"ipv4":true}}

### spec.tags

`map<string, string>`

## Validation Rules

- `mongo_cluster_default_mode_requirements`: Default-mode clusters (create_mode unset or 'Default') require administrator_username, version, compute_tier, storage_size_in_gb, shard_count, and high_availability_mode
- `mongo_cluster_geo_replica_requirements`: GeoReplica-mode clusters require source_server_id and source_location
- `mongo_cluster_restore_requirements`: PointInTimeRestore-mode clusters require the restore block
- `mongo_cluster_admin_credentials_pair`: administrator_username and administrator_password must be set together
- `mongo_cluster_source_location_pairs_with_source_server`: source_location requires source_server_id
- `mongo_cluster_burstable_tier_limits`: The Free and M25 tiers cannot use ZoneRedundantPreferred high availability and cannot have more than one shard
- `mongo_cluster_data_api_default_mode_only`: data_api_mode_enabled can only be set on Default-mode clusters
- `mongo_cluster_storage_type_requires_size`: storage_type requires storage_size_in_gb
- `mongo_cluster_cmk_requires_identity`: customer_managed_key requires the unwrap identity to be listed in user_assigned_identity_ids
- `mongo_cluster_firewall_rule_names_unique`: firewall_rules names must be unique -- the name is the rule's identity

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureMongoCluster, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.mongo_cluster_id` | `string` |  |
| `status.outputs.mongo_cluster_name` | `string` |  |
| `status.outputs.connection_string` | `string` |  |
| `status.outputs.connection_strings` | `map<string, string>` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |
| `spec.userAssignedIdentityIds` | AzureUserAssignedIdentity | `status.outputs.identity_id` |
| `spec.customerManagedKey.keyVaultKeyId` | AzureKeyVaultKey | `status.outputs.versionless_id` |
| `spec.customerManagedKey.userAssignedIdentityId` | AzureUserAssignedIdentity | `status.outputs.identity_id` |
| `spec.sourceServerId` | AzureMongoCluster | `status.outputs.mongo_cluster_id` |
| `spec.restore.sourceId` | AzureMongoCluster | `status.outputs.mongo_cluster_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureMongoCluster | `spec.sourceServerId` | `status.outputs.mongo_cluster_id` |
| AzureMongoCluster | `spec.restore.sourceId` | `status.outputs.mongo_cluster_id` |
| AzureMongoClusterUser | `spec.mongoClusterId` | `status.outputs.mongo_cluster_id` |

## See Also

- [Overview](../README.md)
