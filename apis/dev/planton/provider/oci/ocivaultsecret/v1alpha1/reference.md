# OciVaultSecret

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `oci.planton.dev/v1alpha1`

OciVaultSecretSpec defines the specification for an OCI Vault Secret --
a named piece of sensitive data (credential, certificate, API key, etc.)
stored in an OCI KMS Vault and encrypted by a master encryption key.

Secrets support two mutually exclusive content modes:
  - Explicit content via secret_content (user provides base64-encoded data)
  - Auto-generation via enable_auto_generation + secret_generation_context
    (OCI generates the content using a template)

Lifecycle management is controlled through secret_rules (expiry and reuse
policies) and rotation_config (scheduled rotation against a target system).
Content updates create new secret versions automatically; the
current_version_number output tracks the active version.

vault_id, key_id, and secret_name are immutable after creation (ForceNew).

Excluded from v1:
  - replication_config -- cross-region secret replication is a separate
    infrastructure concern (consistent with OciKmsVault excluding vault
    replication); requires vault + key OCIDs from other regions
  - defined_tags, system_tags -- managed by platform
  - freeform_tags -- auto-populated from metadata labels

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.compartmentId` | `string \| valueFrom` | yes |  | OciCompartment (`status.outputs.compartment_id`) |
| `spec.secretName` | `string` | yes |  |  |
| `spec.vaultId` | `string \| valueFrom` | yes |  | OciKmsVault (`status.outputs.vault_id`) |
| `spec.keyId` | `string \| valueFrom` | yes |  | OciKmsKey (`status.outputs.key_id`) |
| `spec.description` | `string` |  |  |  |
| `spec.secretContent` | `SecretContent` |  |  |  |
| `spec.secretContent.content` | `string` |  |  |  |
| `spec.secretContent.name` | `string` |  |  |  |
| `spec.secretContent.stage` | `string` |  |  |  |
| `spec.enableAutoGeneration` | `bool` |  |  |  |
| `spec.secretGenerationContext` | `SecretGenerationContext` |  |  |  |
| `spec.secretGenerationContext.generationType` | `enum` |  |  |  |
| `spec.secretGenerationContext.generationTemplate` | `string` | yes |  |  |
| `spec.secretGenerationContext.passphraseLength` | `int32` |  |  |  |
| `spec.secretGenerationContext.secretTemplate` | `string` |  |  |  |
| `spec.secretRules` | `[]SecretRule` |  |  |  |
| `spec.secretRules[].ruleType` | `enum` |  |  |  |
| `spec.secretRules[].isSecretContentRetrievalBlockedOnExpiry` | `bool` |  |  |  |
| `spec.secretRules[].secretVersionExpiryInterval` | `string` |  |  |  |
| `spec.secretRules[].timeOfAbsoluteExpiry` | `string` |  |  |  |
| `spec.secretRules[].isEnforcedOnDeletedSecretVersions` | `bool` |  |  |  |
| `spec.rotationConfig` | `RotationConfig` |  |  |  |
| `spec.rotationConfig.isScheduledRotationEnabled` | `bool` |  |  |  |
| `spec.rotationConfig.rotationInterval` | `string` |  |  |  |
| `spec.rotationConfig.targetSystemDetails` | `TargetSystemDetails` | yes |  |  |
| `spec.rotationConfig.targetSystemDetails.targetSystemType` | `enum` |  |  |  |
| `spec.rotationConfig.targetSystemDetails.adbId` | `string \| valueFrom` |  |  | OciAutonomousDatabase (`status.outputs.autonomous_database_id`) |
| `spec.rotationConfig.targetSystemDetails.functionId` | `string \| valueFrom` |  |  |  |
| `spec.secretMetadata` | `map<string, string>` |  |  |  |

## Field Details

### spec.compartmentId

`string | valueFrom` · required

OCID of the compartment where the secret will be created.

- references: OciCompartment (`status.outputs.compartment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciCompartment, name: <that resource's name>, fieldPath: status.outputs.compartment_id}} -- a bare string does not parse

### spec.secretName

`string` · required

User-friendly name for the secret. Must be unique within the vault.
Immutable after creation.

- rule: {"string":{"minLen":"1"}}

### spec.vaultId

`string | valueFrom` · required

OCID of the vault that will contain this secret.
Immutable after creation.

- references: OciKmsVault (`status.outputs.vault_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciKmsVault, name: <that resource's name>, fieldPath: status.outputs.vault_id}} -- a bare string does not parse

### spec.keyId

`string | valueFrom` · required

OCID of the master encryption key used to encrypt this secret.
Must be a symmetric key within the specified vault.
Immutable after creation.

- references: OciKmsKey (`status.outputs.key_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_id}} -- a bare string does not parse

### spec.description

`string`

Brief description of the secret.

### spec.secretContent

`SecretContent`

Explicit secret content (base64-encoded). Mutually exclusive with
auto-generation (enable_auto_generation + secret_generation_context).
Updating this field creates a new secret version.

### spec.secretContent.content

`string`

Base64-encoded secret data.

### spec.secretContent.name

`string`

Optional version name. Must be unique across versions of this secret.

### spec.secretContent.stage

`string`

Rotation state of this content version.
Valid values: "CURRENT" (default when omitted), "PENDING".

- rule: {"string":{"in":["","CURRENT","PENDING"]}}

### spec.enableAutoGeneration

`bool`

Enable automatic secret content generation. When true,
secret_generation_context must be provided and secret_content
must not be set.

### spec.secretGenerationContext

`SecretGenerationContext`

Configuration for automatic secret content generation.
Required when enable_auto_generation is true; must not be set otherwise.

- rule: passphrase_length must be greater than 0 when generation_type is passphrase

### spec.secretGenerationContext.generationType

`enum`

Type of secret content to generate.

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `unspecified`
- `bytes`
- `passphrase`
- `ssh_key`

### spec.secretGenerationContext.generationTemplate

`string` · required

Name of the generation template. Template names are provider-defined
and vary by generation_type.

- rule: {"string":{"minLen":"1"}}

### spec.secretGenerationContext.passphraseLength

`int32`

Length of the passphrase to generate. Required when generation_type
is passphrase.

### spec.secretGenerationContext.secretTemplate

`string`

Optional template structure for storing the generated secret.
Supports predefined placeholders for the generated values.

### spec.secretRules

`[]SecretRule`

Lifecycle rules controlling secret expiry and content reuse.

### spec.secretRules[].ruleType

`enum`

Discriminator selecting the rule type.

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `unspecified`
- `secret_expiry_rule`
- `secret_reuse_rule`

### spec.secretRules[].isSecretContentRetrievalBlockedOnExpiry

`bool`

Block retrieval of secret content after the version expires.
Applies when rule_type is secret_expiry_rule.

### spec.secretRules[].secretVersionExpiryInterval

`string`

Duration after which each secret version expires, in ISO 8601
format (e.g., "P30D" for 30 days). Range: 1-90 days.
Applies when rule_type is secret_expiry_rule.

### spec.secretRules[].timeOfAbsoluteExpiry

`string`

Absolute expiry timestamp in RFC 3339 format. Range: 1-365 days
from creation. Applies when rule_type is secret_expiry_rule.

### spec.secretRules[].isEnforcedOnDeletedSecretVersions

`bool`

Enforce the reuse rule even on deleted secret versions.
Applies when rule_type is secret_reuse_rule.

### spec.rotationConfig

`RotationConfig`

Configuration for scheduled secret rotation against a target system.

### spec.rotationConfig.isScheduledRotationEnabled

`bool`

Enable scheduled automatic rotation.

### spec.rotationConfig.rotationInterval

`string`

Rotation interval in ISO 8601 duration format (e.g., "P30D").
Range: 1-360 days. Required when is_scheduled_rotation_enabled
is true.

### spec.rotationConfig.targetSystemDetails

`TargetSystemDetails` · required

Target system that will be updated during rotation.

- rule: {"required":true}

### spec.rotationConfig.targetSystemDetails.targetSystemType

`enum`

Type of target system.

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `unspecified`
- `adb`
- `function`

### spec.rotationConfig.targetSystemDetails.adbId

`string | valueFrom`

OCID of the Autonomous Database whose credentials are rotated.
Required when target_system_type is adb.

- references: OciAutonomousDatabase (`status.outputs.autonomous_database_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OciAutonomousDatabase, name: <that resource's name>, fieldPath: status.outputs.autonomous_database_id}} -- a bare string does not parse

### spec.rotationConfig.targetSystemDetails.functionId

`string | valueFrom`

OCID of the OCI Functions function invoked during rotation.
Required when target_system_type is function.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.secretMetadata

`map<string, string>`

Additional metadata key-value pairs for administrative context
(e.g., rotation notes, target system hints). Named secret_metadata
to distinguish from Planton's CloudResourceMetadata.

## Validation Rules

- `content_excludes_auto_generation`: secret_content must not be set when enable_auto_generation is true
- `auto_generation_requires_context`: secret_generation_context is required when enable_auto_generation is true
- `no_context_without_auto_generation`: secret_generation_context must not be set when enable_auto_generation is false

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OciVaultSecret, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.secret_id` | `string` | OCID of the Vault Secret. |
| `status.outputs.current_version_number` | `string` | Version number of the currently active secret version. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.compartmentId` | OciCompartment | `status.outputs.compartment_id` |
| `spec.vaultId` | OciKmsVault | `status.outputs.vault_id` |
| `spec.keyId` | OciKmsKey | `status.outputs.key_id` |
| `spec.rotationConfig.targetSystemDetails.adbId` | OciAutonomousDatabase | `status.outputs.autonomous_database_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
