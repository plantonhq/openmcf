# OciKmsKey

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `oci.planton.dev/v1`

OciKmsKeySpec defines the specification for an OCI KMS Key --
a cryptographic key stored inside an OCI KMS Vault, used to encrypt
data at rest across OCI services (Block Volume, Object Storage, Database,
File Storage, etc.).

Key shape is immutable after creation -- algorithm, length, and curve
cannot be changed. Protection mode (HSM vs SOFTWARE vs EXTERNAL) is also
immutable and determines whether the key material is stored in a
hardware security module, software, or an external key manager.

The management_endpoint ties this key to a specific vault. It must match
the vault's management endpoint output.

Excluded from v1:
  - oci_kms_key_version -- creating a key version is manual rotation,
    an operational concern; auto-rotation handles this declaratively
  - desired_state -- operational toggle (ENABLED/DISABLED); keys always
    create as ENABLED
  - restore_from_file / restore_from_object_store / restore_trigger --
    operational restore, not a deployment concern
  - time_of_deletion -- deletion scheduling, operational concern
  - defined_tags, system_tags -- managed by platform
  - freeform_tags -- auto-populated from metadata labels

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.compartmentId` | `string \| valueFrom` | yes |  | OciCompartment (`status.outputs.compartment_id`) |
| `spec.displayName` | `string` |  |  |  |
| `spec.managementEndpoint` | `string \| valueFrom` | yes |  | OciKmsVault (`status.outputs.management_endpoint`) |
| `spec.keyShape` | `KeyShape` | yes |  |  |
| `spec.keyShape.algorithm` | `enum` |  |  |  |
| `spec.keyShape.length` | `int32` |  |  |  |
| `spec.keyShape.curveId` | `enum` |  |  |  |
| `spec.protectionMode` | `enum` |  |  |  |
| `spec.isAutoRotationEnabled` | `bool` |  |  |  |
| `spec.autoKeyRotationDetails` | `AutoKeyRotationDetails` |  |  |  |
| `spec.autoKeyRotationDetails.rotationIntervalInDays` | `int32` |  |  |  |
| `spec.autoKeyRotationDetails.timeOfScheduleStart` | `string` |  |  |  |
| `spec.externalKeyReference` | `ExternalKeyReference` |  |  |  |
| `spec.externalKeyReference.externalKeyId` | `string` | yes |  |  |

## Field Details

### spec.compartmentId

`string | valueFrom` · required

OCID of the compartment where the key will be created.

- references: OciCompartment (`status.outputs.compartment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciCompartment, name: <that resource's name>, fieldPath: status.outputs.compartment_id}} -- a bare string does not parse

### spec.displayName

`string`

Display name for the key. When omitted, the metadata name is used.

### spec.managementEndpoint

`string | valueFrom` · required

Management endpoint of the vault that will contain this key.
This is the vault-specific API endpoint for key management operations.

- references: OciKmsVault (`status.outputs.management_endpoint`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciKmsVault, name: <that resource's name>, fieldPath: status.outputs.management_endpoint}} -- a bare string does not parse

### spec.keyShape

`KeyShape` · required

Cryptographic properties of the key. Immutable after creation.

- rule: {"required":true}
- rule: AES key length must be 16 (128-bit), 24 (192-bit), or 32 (256-bit) bytes
- rule: RSA key length must be 256 (2048-bit), 384 (3072-bit), or 512 (4096-bit) bytes
- rule: ECDSA key requires a matching curve_id and length: nist_p256/32, nist_p384/48, or nist_p521/66
- rule: curve_id must not be set for AES or RSA algorithms

### spec.keyShape.algorithm

`enum`

Encryption algorithm. Immutable after creation.

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `unspecified`
- `aes`
- `rsa`
- `ecdsa`

### spec.keyShape.length

`int32`

Key length in bytes. Valid values depend on the algorithm:
  AES:   16 (128-bit), 24 (192-bit), 32 (256-bit)
  RSA:   256 (2048-bit), 384 (3072-bit), 512 (4096-bit)
  ECDSA: 32 (P-256), 48 (P-384), 66 (P-521)

- rule: {"int32":{"gt":0}}

### spec.keyShape.curveId

`enum`

Elliptic curve identifier. Required for ECDSA; must not be set
for AES or RSA.

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `curve_unspecified`
- `nist_p256`
- `nist_p384`
- `nist_p521`

### spec.protectionMode

`enum`

Protection mode determines where key material is stored.
When omitted, OCI defaults to HSM (hardware security module).
Immutable after creation.

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `unspecified`
- `hsm`
- `software`
- `external`

### spec.isAutoRotationEnabled

`bool`

Whether automatic key rotation is enabled.

### spec.autoKeyRotationDetails

`AutoKeyRotationDetails`

Schedule configuration for automatic key rotation.
Only valid when is_auto_rotation_enabled is true.

### spec.autoKeyRotationDetails.rotationIntervalInDays

`int32`

Rotation interval in days. When omitted, OCI uses its default
rotation interval.

### spec.autoKeyRotationDetails.timeOfScheduleStart

`string`

RFC 3339 timestamp for when the first automatic rotation should
occur. When omitted, OCI schedules based on creation time.

### spec.externalKeyReference

`ExternalKeyReference`

Reference to an external key on a third-party key manager.
Required when protection_mode is external; must not be set otherwise.

### spec.externalKeyReference.externalKeyId

`string` · required

Identifier of the key on the external key manager.

- rule: {"string":{"minLen":"1"}}

## Validation Rules

- `external_requires_key_reference`: external_key_reference is required when protection_mode is external
- `non_external_forbids_key_reference`: external_key_reference must not be set when protection_mode is not external
- `auto_rotation_details_requires_enabled`: auto_key_rotation_details can only be set when is_auto_rotation_enabled is true

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OciKmsKey, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.key_id` | `string` | OCID of the KMS key. |
| `status.outputs.current_key_version` | `string` | OCID of the currently active key version. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.compartmentId` | OciCompartment | `status.outputs.compartment_id` |
| `spec.managementEndpoint` | OciKmsVault | `status.outputs.management_endpoint` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| OciComputeInstance | `spec.sourceDetails.kmsKeyId` | `status.outputs.key_id` |
| OciContainerEngineCluster | `spec.kmsKeyId` | `status.outputs.key_id` |
| OciContainerEngineCluster | `spec.imagePolicyConfig.keyDetails[].kmsKeyId` | `status.outputs.key_id` |
| OciFunctionsApplication | `spec.imagePolicyConfig.keyDetails[].kmsKeyId` | `status.outputs.key_id` |
| OciQueue | `spec.customEncryptionKeyId` | `status.outputs.key_id` |
| OciStreamPool | `spec.kmsKeyId` | `status.outputs.key_id` |
| OciVaultSecret | `spec.keyId` | `status.outputs.key_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
