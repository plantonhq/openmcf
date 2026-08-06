# OciLogGroup

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `oci.planton.dev/v1alpha1`

OciLogGroupSpec defines the specification for an OCI Log Group --
the organizational container for OCI Logging service logs. A log
group holds zero or more individual logs, which can be either
service logs (auto-collected from OCI services) or custom logs
(ingested via the Logging Ingestion API).

The component bundles log group + logs because logs cannot exist
outside a log group, and the common pattern is to create a group
with its constituent logs in a single deployment.

Key behaviors:
  - Log group attributes (compartment_id, display_name, description)
    are all updatable
  - log_type and the entire configuration block on individual logs
    are ForceNew (changing them forces log recreation)
  - is_enabled and retention_duration on logs are updatable
  - display_name must be unique within the log group

Excluded from v1:
  - defined_tags, system_tags -- managed by platform
  - freeform_tags -- auto-populated from metadata labels
  - configuration.source.source_type -- hardcoded to "OCISERVICE"
    (the only valid value)

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.compartmentId` | `string \| valueFrom` | yes |  | OciCompartment (`status.outputs.compartment_id`) |
| `spec.description` | `string` |  |  |  |
| `spec.logs` | `[]Log` |  |  |  |
| `spec.logs[].displayName` | `string` | yes |  |  |
| `spec.logs[].logType` | `enum` |  |  |  |
| `spec.logs[].isEnabled` | `bool` |  |  |  |
| `spec.logs[].retentionDuration` | `int32` |  |  |  |
| `spec.logs[].configuration` | `ServiceLogConfiguration` |  |  |  |
| `spec.logs[].configuration.service` | `string` | yes |  |  |
| `spec.logs[].configuration.resource` | `string \| valueFrom` | yes |  |  |
| `spec.logs[].configuration.category` | `string` | yes |  |  |
| `spec.logs[].configuration.parameters` | `map<string, string>` |  |  |  |
| `spec.logs[].configuration.compartmentId` | `string \| valueFrom` |  |  | OciCompartment (`status.outputs.compartment_id`) |

## Field Details

### spec.compartmentId

`string | valueFrom` · required

OCID of the compartment where the log group will be created.

- references: OciCompartment (`status.outputs.compartment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciCompartment, name: <that resource's name>, fieldPath: status.outputs.compartment_id}} -- a bare string does not parse

### spec.description

`string`

Description for the log group.

### spec.logs

`[]Log`

Logs within this group. Each log is identified by its display_name,
which is used as the resource key in IaC modules and must be unique
within the log group.

- rule: configuration is required when log_type is 'service'

### spec.logs[].displayName

`string` · required

Display name for the log. Must be unique within the log group.
Used as the for_each key in IaC modules.

- rule: {"string":{"minLen":"1"}}

### spec.logs[].logType

`enum`

Type of log. Must be explicitly set.
custom: entries pushed via Logging Ingestion API.
service: auto-collected from an OCI service.

- rule: log_type must be explicitly set (custom or service)

Allowed values (use exactly as shown):

- `unspecified`
- `custom`
- `service`

### spec.logs[].isEnabled

`bool` · optional (explicit presence)

Whether the log is enabled. When nil, OCI defaults to true.

### spec.logs[].retentionDuration

`int32` · optional (explicit presence)

Retention period in days, in 30-day increments (30, 60, 90, 120,
150, 180). When nil, OCI defaults to 30 days.

- rule: retention_duration must be a 30-day increment between 30 and 180

### spec.logs[].configuration

`ServiceLogConfiguration`

Source configuration for service logs. Required when log_type is
service; ignored for custom logs.

The provider nests this as configuration > source, but we flatten
by removing the source wrapper (it is the only child alongside an
optional compartment override). source_type is hardcoded to
"OCISERVICE" in IaC modules (the only valid value).

### spec.logs[].configuration.service

`string` · required

OCI service generating the log.
Examples: "objectstorage", "flowlogs", "apigateway",
"loadbalancer", "functionsInvoke".

- rule: {"string":{"minLen":"1"}}

### spec.logs[].configuration.resource

`string | valueFrom` · required

OCID of the resource emitting logs. Polymorphic -- could be a
VCN (for flow logs), a bucket (for object storage), an API
gateway, etc. StringValueOrRef without default_kind because
the resource type varies; valueFrom enables composability
with any Planton component.

- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.logs[].configuration.category

`string` · required

Log category within the service.
Examples: "write", "read", "all" (object storage);
"access" (API gateway); "invoke" (functions).

- rule: {"string":{"minLen":"1"}}

### spec.logs[].configuration.parameters

`map<string, string>`

Additional parameters for the log source. Pass-through to OCI.

### spec.logs[].configuration.compartmentId

`string | valueFrom`

Optional compartment override for source resource lookup.
When omitted, the log group's compartment is used.

- references: OciCompartment (`status.outputs.compartment_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OciCompartment, name: <that resource's name>, fieldPath: status.outputs.compartment_id}} -- a bare string does not parse

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OciLogGroup, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.log_group_id` | `string` | OCID of the log group. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.compartmentId` | OciCompartment | `status.outputs.compartment_id` |
| `spec.logs[].configuration.compartmentId` | OciCompartment | `status.outputs.compartment_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
