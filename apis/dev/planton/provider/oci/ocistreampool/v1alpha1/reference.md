# OciStreamPool

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `oci.planton.dev/v1alpha1`

OciStreamPoolSpec defines the specification for an OCI Stream Pool --
the organizational container for OCI Streaming, a Kafka-compatible
managed event-streaming service.

A stream pool groups streams under shared Kafka compatibility settings,
optional KMS encryption, and optional private networking. Streams
within the pool inherit the pool's configuration and are bundled as
sub-resources.

Key behaviors:
  - kafka_settings and kms_key_id are updatable after creation
  - private_endpoint_settings is entirely ForceNew (changing any field
    forces pool recreation)
  - Streams are ForceNew for name, partitions, and retention_in_hours

Excluded from v1:
  - defined_tags, system_tags -- managed by platform
  - freeform_tags -- auto-populated from metadata labels
  - security_attributes -- Oracle ZPR, very low adoption

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.compartmentId` | `string \| valueFrom` | yes |  | OciCompartment (`status.outputs.compartment_id`) |
| `spec.kafkaSettings` | `KafkaSettings` |  |  |  |
| `spec.kafkaSettings.autoCreateTopicsEnable` | `bool` |  |  |  |
| `spec.kafkaSettings.logRetentionHours` | `int32` |  |  |  |
| `spec.kafkaSettings.numPartitions` | `int32` |  |  |  |
| `spec.kmsKeyId` | `string \| valueFrom` |  |  | OciKmsKey (`status.outputs.key_id`) |
| `spec.privateEndpointSettings` | `PrivateEndpointSettings` |  |  |  |
| `spec.privateEndpointSettings.subnetId` | `string \| valueFrom` | yes |  | OciSubnet (`status.outputs.subnet_id`) |
| `spec.privateEndpointSettings.nsgIds` | `[]string \| valueFrom` |  |  | OciSecurityGroup (`status.outputs.network_security_group_id`) |
| `spec.privateEndpointSettings.privateEndpointIp` | `string` |  |  |  |
| `spec.streams` | `[]Stream` |  |  |  |
| `spec.streams[].name` | `string` | yes |  |  |
| `spec.streams[].partitions` | `int32` |  |  |  |
| `spec.streams[].retentionInHours` | `int32` |  |  |  |

## Field Details

### spec.compartmentId

`string | valueFrom` · required

OCID of the compartment where the stream pool will be created.

- references: OciCompartment (`status.outputs.compartment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciCompartment, name: <that resource's name>, fieldPath: status.outputs.compartment_id}} -- a bare string does not parse

### spec.kafkaSettings

`KafkaSettings`

Kafka compatibility settings for the stream pool. All fields are
optional and updatable; OCI applies sensible defaults when omitted.

- rule: log_retention_hours must be between 24 and 168 when specified

### spec.kafkaSettings.autoCreateTopicsEnable

`bool` · optional (explicit presence)

Enable automatic topic creation when a Kafka producer publishes
to a non-existent topic.

### spec.kafkaSettings.logRetentionHours

`int32` · optional (explicit presence)

Default hours to retain log data before deletion (24-168).
OCI defaults to 24 when omitted.

### spec.kafkaSettings.numPartitions

`int32` · optional (explicit presence)

Default number of partitions for auto-created topics.

### spec.kmsKeyId

`string | valueFrom`

OCID of the KMS master encryption key for encrypting streams in
this pool. When omitted, Oracle-managed encryption is used.
Updatable -- the pool can be re-encrypted with a different key.

- references: OciKmsKey (`status.outputs.key_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OciKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_id}} -- a bare string does not parse

### spec.privateEndpointSettings

`PrivateEndpointSettings`

Private endpoint configuration for the stream pool. When provided,
the pool is accessible only from within the specified subnet.
Entire block is ForceNew -- any change forces pool recreation.

### spec.privateEndpointSettings.subnetId

`string | valueFrom` · required

OCID of the subnet for the private endpoint.

- references: OciSubnet (`status.outputs.subnet_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.privateEndpointSettings.nsgIds

`[]string | valueFrom`

OCIDs of network security groups for the private endpoint.

- references: OciSecurityGroup (`status.outputs.network_security_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OciSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.network_security_group_id}} -- a bare string does not parse

### spec.privateEndpointSettings.privateEndpointIp

`string`

Specific IP address within the subnet CIDR for the private
endpoint. When omitted, OCI auto-assigns an available IP.

### spec.streams

`[]Stream`

Streams within this pool. Each stream is identified by its name,
which is used as the resource key in IaC modules.

- rule: retention_in_hours must be between 24 and 168 when specified

### spec.streams[].name

`string` · required

Stream name. Used as the resource key in IaC modules. ForceNew.

- rule: {"string":{"minLen":"1"}}

### spec.streams[].partitions

`int32`

Number of partitions in the stream. ForceNew.

- rule: {"int32":{"gte":1}}

### spec.streams[].retentionInHours

`int32` · optional (explicit presence)

Retention period in hours (24-168). ForceNew.
When omitted, OCI defaults to 24 hours.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OciStreamPool, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.stream_pool_id` | `string` | OCID of the stream pool. |
| `status.outputs.endpoint_fqdn` | `string` | FQDN for accessing streams in the pool. For private pools, this resolves only within the associated subnet. |
| `status.outputs.kafka_bootstrap_servers` | `string` | Kafka-compatible bootstrap server string for producing and consuming messages via Kafka clients. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.compartmentId` | OciCompartment | `status.outputs.compartment_id` |
| `spec.kmsKeyId` | OciKmsKey | `status.outputs.key_id` |
| `spec.privateEndpointSettings.subnetId` | OciSubnet | `status.outputs.subnet_id` |
| `spec.privateEndpointSettings.nsgIds` | OciSecurityGroup | `status.outputs.network_security_group_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
