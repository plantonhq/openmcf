# OciNosqlTable

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `oci.planton.dev/v1`

OciNosqlTableSpec defines the specification for an OCI NoSQL table --
a fully managed, serverless NoSQL database with provisioned or on-demand
throughput capacity.

Table schema is defined via a DDL statement (e.g. CREATE TABLE ...),
which is OCI NoSQL's native schema definition mechanism. The DDL
supports columns of type INTEGER, LONG, FLOAT, DOUBLE, NUMBER,
STRING, BOOLEAN, BINARY, TIMESTAMP, JSON, ENUM, ARRAY, MAP, RECORD,
and UUID. Primary keys and shard keys are declared within the DDL.

Secondary indexes are bundled as sub-resources. Each index is
immutable; changes require recreation.

Excluded from v1:
  - defined_tags, system_tags -- managed by platform via freeform_tags
  - freeform_tags -- auto-populated from metadata labels
  - multi-region replicas -- advanced replication feature
  - child table support (names containing ".") -- separate deployments

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.compartmentId` | `string \| valueFrom` | yes |  | OciCompartment (`status.outputs.compartment_id`) |
| `spec.name` | `string` | yes |  |  |
| `spec.ddlStatement` | `string` | yes |  |  |
| `spec.tableLimits` | `TableLimits` | yes |  |  |
| `spec.tableLimits.capacityMode` | `enum` |  |  |  |
| `spec.tableLimits.maxReadUnits` | `int32` |  |  |  |
| `spec.tableLimits.maxWriteUnits` | `int32` |  |  |  |
| `spec.tableLimits.maxStorageInGbs` | `int32` |  |  |  |
| `spec.isAutoReclaimable` | `bool` |  |  |  |
| `spec.indexes` | `[]Index` |  |  |  |
| `spec.indexes[].name` | `string` | yes |  |  |
| `spec.indexes[].keys` | `[]IndexKey` | yes |  |  |
| `spec.indexes[].keys[].columnName` | `string` | yes |  |  |
| `spec.indexes[].keys[].jsonFieldType` | `string` |  |  |  |
| `spec.indexes[].keys[].jsonPath` | `string` |  |  |  |

## Field Details

### spec.compartmentId

`string | valueFrom` · required

OCID of the compartment where the NoSQL table will be created.

- references: OciCompartment (`status.outputs.compartment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciCompartment, name: <that resource's name>, fieldPath: status.outputs.compartment_id}} -- a bare string does not parse

### spec.name

`string` · required

Table name. Must match the table name used in ddl_statement.
Changing this forces recreation.

- rule: {"string":{"minLen":"1"}}

### spec.ddlStatement

`string` · required

DDL statement defining the table schema. For new tables use
CREATE TABLE; for schema evolution use ALTER TABLE.
Column order should not change; new columns can only be appended.
Example: "CREATE TABLE users (id INTEGER, name STRING, PRIMARY KEY(id))"

- rule: {"string":{"minLen":"1"}}

### spec.tableLimits

`TableLimits` · required

Throughput and storage limits for the table.

- rule: {"required":true}

### spec.tableLimits.capacityMode

`enum`

Capacity mode determines how throughput is managed.
When on_demand, max_read_units and max_write_units are ignored.
When unset, defaults to provisioned.

Allowed values (use exactly as shown):

- `capacity_mode_unspecified`
- `provisioned` -- User specifies max_read_units and max_write_units.
- `on_demand` -- Throughput scales automatically; read/write units are ignored.

### spec.tableLimits.maxReadUnits

`int32`

Maximum sustained read throughput limit.
Required when capacity_mode is provisioned (or unspecified).

### spec.tableLimits.maxWriteUnits

`int32`

Maximum sustained write throughput limit.
Required when capacity_mode is provisioned (or unspecified).

### spec.tableLimits.maxStorageInGbs

`int32`

Maximum amount of storage in GB that the table can use.

- rule: {"int32":{"gte":1}}

### spec.isAutoReclaimable

`bool`

When true, the table can be automatically reclaimed by OCI after
an idle period. Changing this forces recreation.

### spec.indexes

`[]Index`

Secondary indexes on the table. Each index is immutable; any
change to an index requires its recreation.

### spec.indexes[].name

`string` · required

Index name. Used as the resource key in IaC modules.

- rule: {"string":{"minLen":"1"}}

### spec.indexes[].keys

`[]IndexKey` · required

Columns included in the index.

- rule: {"repeated":{"minItems":"1"}}

### spec.indexes[].keys[].columnName

`string` · required

Name of the column to index.

- rule: {"string":{"minLen":"1"}}

### spec.indexes[].keys[].jsonFieldType

`string`

If the column is of type JSON, the scalar type of the JSON
field being indexed (e.g. "STRING", "INTEGER", "NUMBER").

### spec.indexes[].keys[].jsonPath

`string`

If the column is of type JSON, the dot-separated path to the
field being indexed (e.g. "address.zipCode").

## Validation Rules

- `provisioned_requires_throughput`: table_limits.max_read_units and table_limits.max_write_units must be greater than zero when capacity_mode is provisioned or unspecified

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OciNosqlTable, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.table_id` | `string` | OCID of the NoSQL table. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.compartmentId` | OciCompartment | `status.outputs.compartment_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
