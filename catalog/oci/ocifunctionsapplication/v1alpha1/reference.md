# OciFunctionsApplication

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `oci.planton.dev/v1alpha1`

OciFunctionsApplicationSpec defines the specification for an OCI Functions
Application -- the organizational container for serverless functions in
Oracle Cloud Infrastructure.

An application provides the shared execution environment (subnets, NSGs,
processor architecture, config) for all functions deployed within it.
Individual functions are deployed as code artifacts (via `fn deploy` or
CI/CD), not managed by IaC.

Key behaviors:
  - display_name, subnet_ids, and shape are immutable after creation (ForceNew)
  - config map entries are passed as environment variables to all functions
  - image_policy_config enables container image signature verification
  - trace_config integrates with OCI Application Performance Monitoring

Excluded from v1:
  - defined_tags, system_tags -- managed by platform
  - freeform_tags -- auto-populated from metadata labels
  - security_attributes -- Oracle ZPR, very low adoption

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.compartmentId` | `string \| valueFrom` | yes |  | OciCompartment (`status.outputs.compartment_id`) |
| `spec.subnetIds` | `[]string \| valueFrom` | yes |  | OciSubnet (`status.outputs.subnet_id`) |
| `spec.displayName` | `string` |  |  |  |
| `spec.shape` | `enum` |  |  |  |
| `spec.config` | `map<string, string>` |  |  |  |
| `spec.networkSecurityGroupIds` | `[]string \| valueFrom` |  |  | OciSecurityGroup (`status.outputs.network_security_group_id`) |
| `spec.syslogUrl` | `string` |  |  |  |
| `spec.imagePolicyConfig` | `ImagePolicyConfig` |  |  |  |
| `spec.imagePolicyConfig.isPolicyEnabled` | `bool` |  |  |  |
| `spec.imagePolicyConfig.keyDetails` | `[]ImagePolicyKeyDetail` |  |  |  |
| `spec.imagePolicyConfig.keyDetails[].kmsKeyId` | `string \| valueFrom` | yes |  | OciKmsKey (`status.outputs.key_id`) |
| `spec.traceConfig` | `TraceConfig` |  |  |  |
| `spec.traceConfig.isEnabled` | `bool` |  |  |  |
| `spec.traceConfig.domainId` | `string` |  |  |  |

## Field Details

### spec.compartmentId

`string | valueFrom` · required

OCID of the compartment where the application will be created.

- references: OciCompartment (`status.outputs.compartment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciCompartment, name: <that resource's name>, fieldPath: status.outputs.compartment_id}} -- a bare string does not parse

### spec.subnetIds

`[]string | valueFrom` · required

OCIDs of the subnets in which to run functions. At least one subnet
is required. Functions execute in these subnets and can reach resources
accessible from them. Immutable after creation.

- references: OciSubnet (`status.outputs.subnet_id`)
- rule: {"repeated":{"minItems":"1"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.displayName

`string`

Display name for the application. Must be unique within the compartment.
When omitted, the metadata name is used. Immutable after creation.

### spec.shape

`enum`

Processor architecture for running functions. When omitted, OCI
defaults to GENERIC_X86. Immutable after creation.
  - generic_x86:     Intel/AMD x86-64
  - generic_arm:     Arm-based (Ampere A1)
  - generic_x86_arm: multi-architecture (runs on either)

Allowed values (use exactly as shown):

- `unspecified`
- `generic_x86`
- `generic_arm`
- `generic_x86_arm`

### spec.config

`map<string, string>`

Application configuration passed as environment variables to all
functions in this application. Keys must be ASCII strings consisting
of letters, digits, and underscores; must not begin with a digit.
Maximum total size for all keys and values is 4 KB.

### spec.networkSecurityGroupIds

`[]string | valueFrom`

OCIDs of network security groups applied to the application.

- references: OciSecurityGroup (`status.outputs.network_security_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OciSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.network_security_group_id}} -- a bare string does not parse

### spec.syslogUrl

`string`

Syslog URL to send all function logs. Supports tcp, udp, and tcp+tls
schemes (e.g., "tcp://logserver.example.com:514"). Must be reachable
from the configured subnets. Ignored when OCI Logging service is enabled.

### spec.imagePolicyConfig

`ImagePolicyConfig`

Image signature verification policy. When enabled, only container
images signed by the specified KMS keys can be used for functions.

### spec.imagePolicyConfig.isPolicyEnabled

`bool`

Whether image signature verification is enabled.

### spec.imagePolicyConfig.keyDetails

`[]ImagePolicyKeyDetail`

KMS keys used to verify image signatures. Required when
is_policy_enabled is true (enforced by CEL on the parent message).

### spec.imagePolicyConfig.keyDetails[].kmsKeyId

`string | valueFrom` · required

OCID of the KMS key used to verify image signatures.

- references: OciKmsKey (`status.outputs.key_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_id}} -- a bare string does not parse

### spec.traceConfig

`TraceConfig`

Distributed tracing configuration for integration with OCI
Application Performance Monitoring (APM).

### spec.traceConfig.isEnabled

`bool` · optional (explicit presence)

Whether tracing is enabled for this application.

### spec.traceConfig.domainId

`string`

OCID of the APM domain (collector) where trace events are sent.
Only meaningful when is_enabled is true.

## Validation Rules

- `image_policy_enabled_requires_keys`: key_details must be non-empty when image_policy_config.is_policy_enabled is true

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OciFunctionsApplication, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.application_id` | `string` | OCID of the functions application. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.compartmentId` | OciCompartment | `status.outputs.compartment_id` |
| `spec.subnetIds` | OciSubnet | `status.outputs.subnet_id` |
| `spec.networkSecurityGroupIds` | OciSecurityGroup | `status.outputs.network_security_group_id` |
| `spec.imagePolicyConfig.keyDetails[].kmsKeyId` | OciKmsKey | `status.outputs.key_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
