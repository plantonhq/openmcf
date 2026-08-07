# ScalewayRdbInstance

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `scaleway.planton.dev/v1alpha1`

ScalewayRdbInstanceSpec defines the specification for a Scaleway Managed
Database (RDB) instance.

Scaleway RDB provides fully managed PostgreSQL and MySQL database engines
with automated backups, high availability, and private network integration.

This is a **composite resource** that bundles several Scaleway resources
into a single declarative unit:
  1. The RDB instance (managed database engine).
  2. Logical databases within the instance.
  3. Database users with optional per-database privileges.
  4. Network ACL rules restricting which IPs can connect.

RDB instances are **regional** resources (e.g., "fr-par", "nl-ams").

**Composition pattern**: The instance accepts a Private Network reference
via `StringValueOrRef`, enabling private connectivity from applications
and other resources on the same network. Downstream resources can
reference `status.outputs.instance_id` for read replicas or monitoring.

**Bundling rationale**: Users want a "ready-to-use database", not just an
engine. Creating an instance without databases, users, or access rules
leaves it inaccessible and useless. Bundling these resources into one
kind follows the same principle as ScalewayLoadBalancer (LB + backends
+ frontends) and DigitalOceanDatabaseCluster.

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.engine` | `string` | yes |  |  |
| `spec.nodeType` | `string` | yes |  |  |
| `spec.privateNetworkId` | `string \| valueFrom` |  |  | ScalewayPrivateNetwork (`status.outputs.private_network_id`) |
| `spec.isHaCluster` | `bool` |  |  |  |
| `spec.volumeType` | `string` |  | `lssd` |  |
| `spec.volumeSizeInGb` | `uint32` |  |  |  |
| `spec.disableBackup` | `bool` |  |  |  |
| `spec.backupScheduleFrequencyHours` | `uint32` |  |  |  |
| `spec.backupScheduleRetentionDays` | `uint32` |  |  |  |
| `spec.encryptionAtRest` | `bool` |  |  |  |
| `spec.aclRules` | `[]ScalewayRdbAclRule` |  |  |  |
| `spec.aclRules[].ip` | `string` | yes |  |  |
| `spec.aclRules[].description` | `string` |  |  |  |
| `spec.adminUser` | `string` | yes |  |  |
| `spec.adminPassword` | `string` (sensitive) | yes |  |  |
| `spec.databases` | `[]ScalewayRdbDatabase` |  |  |  |
| `spec.databases[].name` | `string` | yes |  |  |
| `spec.users` | `[]ScalewayRdbUser` |  |  |  |
| `spec.users[].name` | `string` | yes |  |  |
| `spec.users[].password` | `string` (sensitive) | yes |  |  |
| `spec.users[].isAdmin` | `bool` |  |  |  |
| `spec.users[].privileges` | `[]ScalewayRdbUserPrivilege` |  |  |  |
| `spec.users[].privileges[].databaseName` | `string` | yes |  |  |
| `spec.users[].privileges[].permission` | `string` | yes |  |  |
| `spec.settings` | `map<string, string>` |  |  |  |
| `spec.initSettings` | `map<string, string>` |  |  |  |

## Field Details

### spec.region

`string` · required

The Scaleway region where the instance will be created.
Examples: "fr-par", "nl-ams", "pl-waw"

Determines which availability zones are used for HA replicas.

This field is required and cannot be changed after creation.

- rule: {"required":true}

### spec.engine

`string` · required

Database engine and major version.

Format: "{Engine}-{MajorVersion}" (case-sensitive).
Examples:
  - "PostgreSQL-16" -- PostgreSQL version 16
  - "PostgreSQL-15" -- PostgreSQL version 15
  - "MySQL-8"       -- MySQL version 8

The engine determines which SQL dialect, extensions, and features
are available. The major version determines the feature set and
upgrade path.

IMPORTANT: Cannot be changed after creation. Changing the engine
requires creating a new instance and migrating data.

- rule: {"required":true,"string":{"pattern":"^(PostgreSQL|MySQL)-[0-9]+$"}}

### spec.nodeType

`string` · required

Node type determining CPU, RAM, and baseline performance.

Development types:
  - "DB-DEV-S"  -- 2 vCPU, 2 GB RAM. Cheapest, for development only.
  - "DB-DEV-M"  -- 4 vCPU, 4 GB RAM.

General-purpose types:
  - "db-gp-xs"  -- 4 vCPU, 16 GB RAM. Entry production.
  - "db-gp-s"   -- 8 vCPU, 32 GB RAM.
  - "db-gp-m"   -- 16 vCPU, 64 GB RAM.

Node type can be changed after creation (vertical scaling).

- rule: {"required":true}

### spec.privateNetworkId

`string | valueFrom`

The Private Network to attach the instance to.

When set, the instance receives a private endpoint reachable only
from resources on the same Private Network. This is the recommended
topology for production: applications connect via private IPs,
and the public endpoint can be locked down with ACL rules.

IPAM (automatic IP assignment) is used by default when a Private
Network is attached.

In infra charts, this is typically wired via valueFrom:

  privateNetworkId:
    valueFrom:
      kind: ScalewayPrivateNetwork
      name: app-network
      fieldPath: status.outputs.private_network_id

Optional. If omitted, the instance is accessible only via its
public endpoint (restricted by ACL rules).

- references: ScalewayPrivateNetwork (`status.outputs.private_network_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: ScalewayPrivateNetwork, name: <that resource's name>, fieldPath: status.outputs.private_network_id}} -- a bare string does not parse

### spec.isHaCluster

`bool`

Whether to deploy a multi-node HA cluster with automatic failover.

When true, Scaleway provisions a standby replica that automatically
takes over if the primary fails. Provides zero-downtime failover
but doubles the cost.

Recommended for production workloads. For dev/test, leave false.

Default: false.

### spec.volumeType

`string`

Volume type for database storage.

Options:
  - "lssd" (default) -- Local SSD. Lowest latency, highest IOPS.
    Limited by the node type's local disk capacity.
  - "bssd" -- Block SSD (5K IOPS). Scalable beyond local disk.
    Good balance of performance and flexibility.
  - "sbs_15k" -- Block SSD (15K IOPS). Highest throughput for
    I/O-intensive workloads.

IMPORTANT: Cannot be changed after creation. Changing volume type
requires creating a new instance and migrating data.

- default: `lssd`

### spec.volumeSizeInGb

`uint32`

Volume size in GB.

If 0 or omitted, uses the default size for the node type. The size
determines how much data the database can store.

Can only be increased after creation, never decreased. Plan ahead
for data growth.

### spec.disableBackup

`bool`

Whether to disable automated backups.

Default: false (backups enabled). Set to true for dev/test
instances where backup cost is unnecessary.

When backups are enabled, Scaleway takes automated snapshots at
the configured frequency and retains them for the configured period.

### spec.backupScheduleFrequencyHours

`uint32`

Hours between automated backups.

Range: 1 to 24. If 0 or omitted, uses Scaleway's default (24 hours).
Lower values provide finer recovery point objectives (RPO) but
increase storage cost.

### spec.backupScheduleRetentionDays

`uint32`

Days to retain automated backups.

Range: 1 to 365. If 0 or omitted, uses Scaleway's default (7 days).
Longer retention provides more recovery options but increases
storage cost.

### spec.encryptionAtRest

`bool`

Whether to enable encryption at rest for the database storage.

When true, all data written to disk is encrypted. Provides
defense-in-depth for compliance-sensitive workloads. Minimal
performance impact on modern hardware.

### spec.aclRules

`[]ScalewayRdbAclRule`

Network access control rules.

Each rule allows a CIDR range to connect to the database's public
endpoint. Rules replace the instance's entire ACL (Scaleway's ACL
is a single resource per instance that replaces all rules atomically).

If empty, no ACL resource is created and Scaleway's default applies
(all IPs allowed on the public endpoint). For production instances,
always specify ACL rules to restrict access.

When using Private Network connectivity, ACL rules control the
public endpoint only -- private network access is always allowed.

### spec.aclRules[].ip

`string` · required

CIDR range to allow.

Examples:
  - "10.0.0.0/24" -- Allow a /24 subnet
  - "1.2.3.4/32"  -- Allow a single IP
  - "0.0.0.0/0"   -- Allow all IPs (NOT recommended for production)

- rule: {"required":true}

### spec.aclRules[].description

`string`

Human-readable description for this rule.

Helps operators understand the purpose of each ACL entry.
Examples: "Office IP", "VPN egress", "CI/CD pipeline"

### spec.adminUser

`string` · required

Username for the initial admin user created with the instance.

This user has full administrative privileges on the database engine
(superuser-like). It is created as part of the instance provisioning
and cannot be removed through this resource.

Must be different from any user in the `users` list.

- rule: {"required":true,"string":{"maxLen":"63"}}

### spec.adminPassword

`string` · required · sensitive

Password for the initial admin user.

Must meet minimum complexity requirements. For production, use a
strong, randomly generated password and manage it through your
organization's secrets workflow.

- rule: {"required":true,"string":{"minLen":"8"}}

### spec.databases

`[]ScalewayRdbDatabase`

Logical databases to create on the instance.

Each entry creates a `scaleway_rdb_database` resource. Databases
are where application data is stored. Users connect to a specific
database and are granted privileges on it.

Optional. If empty, only the engine's default database exists
(e.g., "postgres" for PostgreSQL, "mysql" for MySQL). The admin
user can create additional databases manually.

Reserved names (rejected by Scaleway): postgres, mysql, sys,
information_schema, performance_schema, rdb, template0, template1.

### spec.databases[].name

`string` · required

Database name.

Must be 1-63 characters, consisting of alphanumeric characters,
underscores, or dashes. Must start with a letter.

Reserved names are rejected by Scaleway: postgres, mysql, sys,
information_schema, performance_schema, rdb, template0, template1.

- rule: {"required":true,"string":{"maxLen":"63"}}

### spec.users

`[]ScalewayRdbUser`

Additional database users to create on the instance.

Each entry creates a `scaleway_rdb_user` resource. Each user's
privileges (if specified) create `scaleway_rdb_privilege` resources
linking the user to specific databases with specific permission
levels.

The admin user (specified above) always exists and has full access.
Users listed here are typically application-level accounts with
restricted permissions (e.g., a web app user with "readwrite" on
the application database).

Optional. If empty, only the admin user exists.

### spec.users[].name

`string` · required

Username.

Must be unique within the instance and different from the admin_user
specified in the spec. 1-63 characters.

- rule: {"required":true,"string":{"maxLen":"63"}}

### spec.users[].password

`string` · required · sensitive

Password for this user.

Must meet minimum complexity requirements. For production, use a
strong, randomly generated password.

- rule: {"required":true,"string":{"minLen":"8"}}

### spec.users[].isAdmin

`bool`

Whether this user has admin privileges on the database engine.

Admin users have superuser-like access to all databases. Use
sparingly -- most application users should be non-admin with
specific per-database privileges.

Default: false.

### spec.users[].privileges

`[]ScalewayRdbUserPrivilege`

Per-database permission grants for this user.

Each entry creates a `scaleway_rdb_privilege` resource linking
this user to a specific database with a specific permission level.

If empty, the user exists but has no database access until
privileges are granted (via this field or manually).

The database_name can reference a database in the `databases` list
or an engine-default database (e.g., "postgres"). The Scaleway API
validates that the database exists.

### spec.users[].privileges[].databaseName

`string` · required

Name of the database to grant access to.

Must be the name of a database that exists on the instance -- either
one defined in the `databases` list or an engine-default database.
The Scaleway API validates existence at apply time.

- rule: {"required":true}

### spec.users[].privileges[].permission

`string` · required

Permission level for this user on this database.

Options:
  - "readonly"  -- SELECT only. For read-only application accounts
                   or reporting services.
  - "readwrite" -- SELECT, INSERT, UPDATE, DELETE. For typical
                   application accounts.
  - "all"       -- Full access including DDL (CREATE, ALTER, DROP).
                   For migration tools or admin accounts.
  - "none"      -- Explicitly revoke access. Useful for overriding
                   inherited permissions.

- rule: {"required":true,"string":{"in":["readonly","readwrite","all","none"]}}

### spec.settings

`map<string, string>`

Engine-specific runtime settings.

Key-value pairs passed to the database engine's configuration.
Applied on both creation and updates. Keys are engine-specific.

PostgreSQL examples:
  - "work_mem" = "64MB"
  - "max_connections" = "200"
  - "effective_cache_size" = "4GB"

MySQL examples:
  - "max_connections" = "200"
  - "innodb_buffer_pool_size" = "4G"

Optional. If empty, Scaleway uses engine defaults optimized for
the node type.

### spec.initSettings

`map<string, string>`

Engine-specific init settings applied only during instance creation.

These settings cannot be changed after the instance is created.
Use for configuration that must be set at init time.

MySQL example:
  - "lower_case_table_names" = "1"

Optional. Most users don't need init settings.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: ScalewayRdbInstance, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.instance_id` | `string` | The unique identifier of the created RDB instance. Format: regional ID (e.g., "fr-par/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"). This is the primary output referenced by downstream resources: - Read replicas (scaleway_rdb_read_replica) - Monitoring and backup tools - External management integrations |
| `status.outputs.endpoint_ip` | `string` | Public endpoint IP address. The IPv4 address for connecting to the database over the public internet. Subject to ACL rules configured in the spec. Empty if the instance has no public endpoint (rare -- Scaleway provides one by default). |
| `status.outputs.endpoint_port` | `uint32` | Public endpoint port number. The TCP port for public connections. Typically 5432 for PostgreSQL or 3306 for MySQL, but Scaleway may assign different ports. |
| `status.outputs.private_endpoint_ip` | `string` | Private Network endpoint IP address. The IPv4 address for connecting to the database from resources on the same Private Network. This is the recommended connection path for application workloads. Empty if no Private Network is attached. |
| `status.outputs.private_endpoint_port` | `uint32` | Private Network endpoint port number. The TCP port for private connections. Zero if no Private Network is attached. |
| `status.outputs.certificate` | `string` | TLS certificate in PEM format for verifying the database server. Clients should use this CA certificate to establish encrypted connections and verify the database server's identity. Both PostgreSQL (`sslrootcert`) and MySQL (`ssl-ca`) support this. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.privateNetworkId` | ScalewayPrivateNetwork | `status.outputs.private_network_id` |

## See Also

- [Overview](../README.md)
