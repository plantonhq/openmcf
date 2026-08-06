# CivoDatabase

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `civo.planton.dev/v1alpha1`

CivoDatabaseSpec defines the essential configuration for creating a managed database instance on Civo.
Following the 80/20 principle, it exposes only the most commonly used fields.

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.dbInstanceName` | `string` | yes |  |  |
| `spec.engine` | `enum` | yes |  |  |
| `spec.engineVersion` | `string` | yes |  |  |
| `spec.region` | `enum` | yes |  |  |
| `spec.sizeSlug` | `string` | yes |  |  |
| `spec.replicas` | `uint32` |  |  |  |
| `spec.networkId` | `string \| valueFrom` | yes |  | CivoVpc (`status.outputs.network_id`) |
| `spec.firewallIds` | `[]string \| valueFrom` |  |  | CivoFirewall (`status.outputs.firewall_id`) |
| `spec.storageGib` | `uint32` |  |  |  |
| `spec.tags` | `[]string` |  |  |  |

## Field Details

### spec.dbInstanceName

`string` · required

A human-readable name for the database instance.
This name must be unique within the chosen Civo region.

- rule: {"required":true,"string":{"maxLen":"64"}}

### spec.engine

`enum` · required

The database engine for the instance (MySQL or PostgreSQL).

- rule: {"required":true}

Allowed values (use exactly as shown):

- `civo_database_engine_unspecified`
- `mysql`
- `postgres`

### spec.engineVersion

`string` · required

The engine version for the database.
For example: "8.0" for MySQL 8.0, "14" for PostgreSQL 14.
Only major (and optionally minor) version numbers are expected.

- rule: {"required":true,"string":{"pattern":"^[0-9]+(\\.[0-9]+)?$"}}

### spec.region

`enum` · required

The Civo region where the database will be created.

- rule: {"required":true}

Allowed values (use exactly as shown):

- `civo_region_unspecified` -- 0: default / unspecified region
- `lon1` -- london 1
- `lon2` -- london 2
- `fra1` -- frankfurt 1
- `nyc1` -- new york 1
- `phx1` -- phoenix 1
- `mum1` -- mumbai 1

### spec.sizeSlug

`string` · required

The plan or size identifier for the database instance (e.g., "g3.db.small").
This defines CPU, memory, and base storage for the instance.

- rule: {"required":true}

### spec.replicas

`uint32`

The number of replica nodes to add to the database (0 means no replicas, just the primary).
Typically use 0, 2, or 4 for a total of 1, 3, or 5 nodes in the cluster.

- rule: {"uint32":{"lte":4}}

### spec.networkId

`string | valueFrom` · required

The target private network for the database instance.
Provide a literal network UUID or a reference to a CivoNetwork resource.

- references: CivoVpc (`status.outputs.network_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: CivoVpc, name: <that resource's name>, fieldPath: status.outputs.network_id}} -- a bare string does not parse

### spec.firewallIds

`[]string | valueFrom`

(Optional) Firewall rules to attach to this database instance.
Provide one or more firewall IDs or references to CivoFirewall resources for access control.

- references: CivoFirewall (`status.outputs.firewall_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: CivoFirewall, name: <that resource's name>, fieldPath: status.outputs.firewall_id}} -- a bare string does not parse

### spec.storageGib

`uint32`

(Optional) Custom storage size in GiB for the database, if different from the default provided by size_slug.

### spec.tags

`[]string`

(Optional) Tags to assign to the database instance for organizational purposes.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: CivoDatabase, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.database_id` | `string` | The unique identifier of the created database instance (UUID). |
| `status.outputs.host` | `string` | The host name or IP address of the primary database endpoint. |
| `status.outputs.port` | `uint32` | The network port on which the database is listening. |
| `status.outputs.username` | `string` | The username for the default database user. |
| `status.outputs.password_secret_ref` | `string` | A reference to a secret containing the default user's password. |
| `status.outputs.replica_endpoints` | `[]string` | The host addresses of replica nodes, if any (empty if no replicas were configured). |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.networkId` | CivoVpc | `status.outputs.network_id` |
| `spec.firewallIds` | CivoFirewall | `status.outputs.firewall_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
