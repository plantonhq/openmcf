# DigitalOceanDatabaseCluster

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `digital-ocean.planton.dev/v1alpha1`

DigitalOceanDatabaseClusterSpec defines the essential configuration for creating a managed database cluster on DigitalOcean.
This follows the 80/20 principle: only the most commonly used fields are exposed to keep the API simple.

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.clusterName` | `string` | yes |  |  |
| `spec.engine` | `enum` | yes |  |  |
| `spec.engineVersion` | `string` | yes |  |  |
| `spec.region` | `enum` | yes |  |  |
| `spec.sizeSlug` | `string` | yes |  |  |
| `spec.nodeCount` | `uint32` | yes |  |  |
| `spec.vpc` | `string \| valueFrom` |  |  | DigitalOceanVpc (`status.outputs.vpc_id`) |
| `spec.storageGib` | `uint32` |  |  |  |
| `spec.enablePublicConnectivity` | `bool` |  |  |  |

## Field Details

### spec.clusterName

`string` · required

A human-readable name for the database cluster.
This name will be used as the cluster's identifier in DigitalOcean.

- rule: {"required":true,"string":{"maxLen":"64"}}

### spec.engine

`enum` · required

The database engine for the cluster.
Allowed values include: POSTGRES, MYSQL, REDIS, MONGODB.

- rule: {"required":true}

Allowed values (use exactly as shown):

- `digital_ocean_database_engine_unspecified`
- `pg`
- `mysql`
- `redis`
- `mongodb`

### spec.engineVersion

`string` · required

The engine version for the cluster.
For example, "14" for PostgreSQL 14, "8" for MySQL 8, etc.
Only major (and optionally minor) version numbers are expected.

- rule: {"required":true,"string":{"pattern":"^[0-9]+(\\.[0-9]+)?$"}}

### spec.region

`enum` · required

The DigitalOcean region where the cluster will be created.
Determines the data center location for the cluster.

- rule: {"required":true}

Allowed values (use exactly as shown):

- `digital_ocean_region_unspecified` -- 0: default / unspecified region
- `nyc3` -- new york 3
- `sfo3` -- san francisco 3
- `fra1` -- frankfurt 1
- `sgp1` -- singapore 1
- `lon1` -- london 1
- `tor1` -- toronto 1
- `blr1` -- bangalore 1
- `ams3` -- amsterdam 3

### spec.sizeSlug

`string` · required

The slug identifier for the cluster's node size (e.g., "db-s-2vcpu-4gb").
This defines the CPU/memory resources for each node in the cluster.

- rule: {"required":true}

### spec.nodeCount

`uint32` · required

The number of nodes in the cluster. Allowed values are 1 to 3 for primary nodes.

- rule: {"required":true,"uint32":{"lte":3,"gte":1}}

### spec.vpc

`string | valueFrom`

(Optional) Reference to a DigitalOcean VPC for the database cluster.
If provided, the cluster will be created within the specified private network.
Use a literal VPC UUID or a reference to a DigitalOceanVpc resource.

- references: DigitalOceanVpc (`status.outputs.vpc_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: DigitalOceanVpc, name: <that resource's name>, fieldPath: status.outputs.vpc_id}} -- a bare string does not parse

### spec.storageGib

`uint32`

(Optional) Custom storage size in GiB for the cluster.
If not set, the default storage for the chosen size_slug will be used.

### spec.enablePublicConnectivity

`bool`

(Optional) Whether to enable cluster access to public networking.
When false (default), no public connection is available; the cluster is accessible only via the VPC or DigitalOcean internal network.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: DigitalOceanDatabaseCluster, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.cluster_id` | `string` | The unique identifier (UUID) of the created database cluster. |
| `status.outputs.connection_uri` | `string` | The full connection URI for the database cluster (including credentials and database name). |
| `status.outputs.host` | `string` | The hostname or IP address at which the database cluster is accessible. |
| `status.outputs.port` | `uint32` | The network port that the database cluster is listening on. |
| `status.outputs.database_user` | `string` | The username for the cluster’s default database user. |
| `status.outputs.database_password` | `string` | The password for the cluster’s default database user. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.vpc` | DigitalOceanVpc | `status.outputs.vpc_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
