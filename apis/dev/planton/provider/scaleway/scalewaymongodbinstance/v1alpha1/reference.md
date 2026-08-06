# ScalewayMongodbInstance

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `scaleway.planton.dev/v1alpha1`

ScalewayMongodbInstanceSpec defines the specification for a Scaleway Managed
MongoDB instance.

Scaleway MongoDB provides fully managed document databases with automatic
failover (replica set), block storage volumes, TLS certificates, and
Private Network integration.

This is a **composite resource** that bundles:
  1. The MongoDB instance (managed database engine with admin user).
  2. Additional database users with role-based access control.

MongoDB instances are **regional** resources (currently only "fr-par").

**Key differences from ScalewayRdbInstance (PostgreSQL/MySQL):**
  - No explicit database creation: MongoDB databases are created implicitly
    when you first write data. There is no `scaleway_mongodb_database`
    resource.
  - No network ACL: MongoDB has no IP-based access control resource.
    Network security is controlled entirely by the Private Network /
    Public Network endpoint choice.
  - No separate privilege resource: User permissions are expressed as
    role assignments (read, read_write, db_admin) directly on the user
    resource, scoped to a database name or all databases.
  - HA via replica set: node_number = 3 gives a 3-node replica set with
    automatic failover. There is no 2-node HA mode.

**Composition pattern:** The instance accepts a Private Network reference
via `StringValueOrRef`, enabling private connectivity from applications.
Downstream resources can reference `status.outputs.instance_id`.

**Bundling rationale:** Users want a "ready-to-use MongoDB" with application
users and roles pre-configured. Creating an instance without users leaves
only the admin user, requiring manual intervention. Bundling users follows
the same principle as ScalewayRdbInstance.

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.version` | `string` | yes |  |  |
| `spec.nodeType` | `string` | yes |  |  |
| `spec.nodeNumber` | `uint32` | yes | `1` |  |
| `spec.privateNetworkId` | `string \| valueFrom` |  |  | ScalewayPrivateNetwork (`status.outputs.private_network_id`) |
| `spec.enablePublicNetwork` | `bool` |  |  |  |
| `spec.volumeType` | `string` |  | `sbs_5k` |  |
| `spec.volumeSizeInGb` | `uint32` |  |  |  |
| `spec.enableSnapshotSchedule` | `bool` |  |  |  |
| `spec.snapshotScheduleFrequencyHours` | `uint32` |  |  |  |
| `spec.snapshotScheduleRetentionDays` | `uint32` |  |  |  |
| `spec.adminUser` | `string` | yes |  |  |
| `spec.adminPassword` | `string` (sensitive) | yes |  |  |
| `spec.users` | `[]ScalewayMongodbUser` |  |  |  |
| `spec.users[].name` | `string` | yes |  |  |
| `spec.users[].password` | `string` (sensitive) | yes |  |  |
| `spec.users[].roles` | `[]ScalewayMongodbUserRole` |  |  |  |
| `spec.users[].roles[].role` | `string` | yes |  |  |
| `spec.users[].roles[].databaseName` | `string` |  |  |  |
| `spec.users[].roles[].anyDatabase` | `bool` |  |  |  |
| `spec.settings` | `map<string, string>` |  |  |  |

## Field Details

### spec.region

`string` · required

The Scaleway region where the instance will be created.

Currently, Managed MongoDB is only available in "fr-par" (Paris).
Additional regions may become available in the future.

This field is required and cannot be changed after creation.

- rule: {"required":true}

### spec.version

`string` · required

MongoDB engine version.

Format: semantic version string (e.g., "7.0.12").

The Scaleway provider normalizes this to major.minor internally.
Check Scaleway documentation for currently supported versions.

Cannot be changed after creation.

- rule: {"required":true,"string":{"pattern":"^[0-9]+\\.[0-9]+\\.[0-9]+$"}}

### spec.nodeType

`string` · required

Node type determining CPU, RAM, and baseline performance.

Shared vCPU (cost-optimized):
  - "MGDB-PLAY2-NANO"  -- 2 vCPU, 4 GB RAM. Development/testing.
  - "MGDB-PRO2-XXS"    -- 2 vCPU, 8 GB RAM. Light production.
  - "MGDB-PRO2-XS"     -- 4 vCPU, 16 GB RAM.
  - "MGDB-PRO2-S"      -- 8 vCPU, 32 GB RAM.
  - "MGDB-PRO2-M"      -- 16 vCPU, 64 GB RAM.
  - "MGDB-PRO2-L"      -- 32 vCPU, 128 GB RAM.

Dedicated vCPU (production-optimized):
  - "MGDB-POP2-2C-8G"     -- 2 vCPU, 8 GB RAM.
  - "MGDB-POP2-4C-16G"    -- 4 vCPU, 16 GB RAM.
  - "MGDB-POP2-8C-32G"    -- 8 vCPU, 32 GB RAM.
  - "MGDB-POP2-16C-64G"   -- 16 vCPU, 64 GB RAM.
  - "MGDB-POP2-32C-128G"  -- 32 vCPU, 128 GB RAM.
  - "MGDB-POP2-64C-256G"  -- 64 vCPU, 256 GB RAM.

Node type can be changed after creation (vertical scaling).

- rule: {"required":true}

### spec.nodeNumber

`uint32` · required

Number of nodes in the instance.

Valid values:
  - 1: Standalone mode (single node, no redundancy).
  - 3: Replica set (1 primary + 2 secondaries, automatic failover).

There is no 2-node mode. MongoDB requires an odd number of voting
members for replica set elections.

IMPORTANT: Changing from 1 to 3 or 3 to 1 may destroy and recreate
the instance. Plan node count before initial deployment.

- default: `1`
- rule: {"required":true}

### spec.privateNetworkId

`string | valueFrom`

The Private Network to attach the instance to.

When set, the instance receives a private endpoint reachable only
from resources on the same Private Network. This is the recommended
topology for production.

When a Private Network is attached and `enable_public_network` is
false (the default), the instance is private-only -- no public
endpoint is created. This is the most secure configuration.

IPAM (automatic IP assignment) is used by default.

In infra charts, this is typically wired via valueFrom:

  privateNetworkId:
    valueFrom:
      kind: ScalewayPrivateNetwork
      name: app-network
      fieldPath: status.outputs.private_network_id

Optional. If omitted, the instance gets a public endpoint by default.

- references: ScalewayPrivateNetwork (`status.outputs.private_network_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: ScalewayPrivateNetwork, name: <that resource's name>, fieldPath: status.outputs.private_network_id}} -- a bare string does not parse

### spec.enablePublicNetwork

`bool`

Whether to also create a public network endpoint when a Private
Network is attached.

This field only takes effect when `private_network_id` is set.
When `private_network_id` is not set, a public endpoint is always
created (Scaleway's default behavior).

Use cases for enabling both endpoints:
  - Admin access from outside the Private Network
  - Debugging connectivity from a local machine
  - Migration tooling running outside the VPC

WARNING: MongoDB on Scaleway has no IP-based ACL rules (unlike RDB).
The public endpoint is accessible from ANY IP address. Only enable
this if you have other network-level controls (e.g., firewall) or
accept the security risk.

Default: false (private-only when PN is attached).

### spec.volumeType

`string`

Volume type for database storage.

Options:
  - "sbs_5k"  (default) -- Block Storage SSD with 5K IOPS.
  - "sbs_15k"           -- Block Storage SSD with 15K IOPS.
    Higher throughput for I/O-intensive workloads.

Unlike RDB (which supports local SSD), MongoDB uses only block
storage volumes.

Cannot be changed after creation.

- default: `sbs_5k`

### spec.volumeSizeInGb

`uint32`

Volume size in GB.

Must be a multiple of 5. Minimum 5 GB.

If 0 or omitted, uses the default size (5 GB). The size determines
how much data the database can store.

Can only be increased after creation, never decreased.

### spec.enableSnapshotSchedule

`bool`

Whether to enable automatic snapshot scheduling.

When enabled, Scaleway takes periodic snapshots of the database at
the configured frequency and retains them for the configured period.

Default: false (no automatic snapshots).

### spec.snapshotScheduleFrequencyHours

`uint32`

Hours between automatic snapshots.

Only used when `enable_snapshot_schedule` is true.
If 0 or omitted when snapshots are enabled, uses Scaleway's default.

### spec.snapshotScheduleRetentionDays

`uint32`

Days to retain automatic snapshots.

Only used when `enable_snapshot_schedule` is true.
If 0 or omitted when snapshots are enabled, uses Scaleway's default.

### spec.adminUser

`string` · required

Username for the initial admin user created with the instance.

This user is created as part of instance provisioning and has full
administrative privileges. It cannot be removed through this resource.

Must be different from any user in the `users` list.

- rule: {"required":true,"string":{"maxLen":"63"}}

### spec.adminPassword

`string` · required · sensitive

Password for the initial admin user.

Must meet minimum complexity requirements. For production, use a
strong, randomly generated password and manage it through your
organization's secrets workflow.

- rule: {"required":true,"string":{"minLen":"8"}}

### spec.users

`[]ScalewayMongodbUser`

Additional database users to create on the instance.

Each entry creates a `scaleway_mongodb_user` Terraform resource.
Each user's roles define what databases they can access and at what
permission level.

The admin user (specified above) always exists with full access.
Users listed here are typically application-level accounts with
restricted roles (e.g., "read_write" on a specific database).

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

### spec.users[].roles

`[]ScalewayMongodbUserRole`

Role assignments for this user.

Each role grants a specific permission level scoped to a database
name or all databases. Multiple roles can be assigned to one user.

If empty, the user exists but has no database access until roles
are assigned (via this field or manually).

- rule: exactly one of database_name or any_database must be set -- they define the scope of the role

### spec.users[].roles[].role

`string` · required

The role to assign.

Available roles:
  - "read"       -- Read-only access (find, count, aggregate).
  - "read_write" -- Read and write access (insert, update, delete).
  - "db_admin"   -- Database administration (index management,
                    collection management, schema validation).

The "sync" role is intentionally excluded: it is a niche replication
role (aggregates clusterMonitor, backup, and restore) that most
users do not need. It can be added via the Scaleway console if
required.

- rule: {"required":true,"string":{"in":["read","read_write","db_admin"]}}

### spec.users[].roles[].databaseName

`string`

The specific database to scope this role to.

When set, the role applies only to this database. The database does
not need to exist at provisioning time -- MongoDB creates databases
implicitly when data is first written.

Mutually exclusive with `any_database`. Exactly one must be set.

### spec.users[].roles[].anyDatabase

`bool`

Whether to apply this role to ALL databases on the instance.

When true, the role grants access to every current and future
database on the instance. Useful for admin or monitoring accounts.

Mutually exclusive with `database_name`. Exactly one must be set.

### spec.settings

`map<string, string>`

MongoDB-specific configuration settings.

Key-value pairs passed to the MongoDB engine configuration.
Applied on both creation and updates. Keys are engine-specific.

Optional. If empty, Scaleway uses MongoDB defaults optimized for
the node type.

## Validation Rules

- `valid_node_number`: node_number must be 1 (standalone) or 3 (replica set) -- MongoDB requires an odd number of voting members

## Outputs

Reference an output from another manifest as `valueFrom: {kind: ScalewayMongodbInstance, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.instance_id` | `string` | The unique identifier of the created MongoDB instance. Format: regional ID (e.g., "fr-par/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"). This is the primary output referenced by downstream resources: - Snapshot tooling - Monitoring integrations - Management automation |
| `status.outputs.public_dns_record` | `string` | Public endpoint DNS record. The DNS hostname for connecting to the database over the public network. Format: "{id}.mgdb.{region}.scw.cloud" Empty if the instance has no public endpoint (private-only mode). |
| `status.outputs.public_port` | `uint32` | Public endpoint port number. The TCP port for public connections (typically 27017). Zero if the instance has no public endpoint. |
| `status.outputs.private_dns_records` | `[]string` | Private Network endpoint DNS records. DNS hostnames for connecting to the database from resources on the same Private Network. Empty if no Private Network is attached. |
| `status.outputs.private_ips` | `[]string` | Private Network endpoint IP addresses. IPv4 addresses assigned via IPAM for connecting to the database from the Private Network. This is the recommended connection path for application workloads. Empty if no Private Network is attached. |
| `status.outputs.private_port` | `uint32` | Private Network endpoint port number. The TCP port for private connections (typically 27017). Zero if no Private Network is attached. |
| `status.outputs.tls_certificate` | `string` | TLS certificate in PEM format for verifying the database server. Clients should use this CA certificate to establish encrypted connections and verify the database server's identity. MongoDB drivers accept this via the `tlsCAFile` connection option or the `--tlsCAFile` flag for the mongo shell. Always available (Scaleway issues a TLS certificate for every MongoDB instance regardless of endpoint type). |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.privateNetworkId` | ScalewayPrivateNetwork | `status.outputs.private_network_id` |

## See Also

- [Overview](./README.md)
