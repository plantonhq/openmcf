# OciKmsVault

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `oci.planton.dev/v1`

OciKmsVaultSpec defines the specification for an OCI KMS Vault --
an HSM-backed container for encryption keys used by other OCI services
(Compute, Block Volume, Object Storage, Database, etc.).

Vault types:
  - default_vault   = Shared HSM partition, lower cost, lower throughput
  - virtual_private  = Dedicated HSM partition, higher throughput limits
  - external         = External key manager (BYOK/EKMS) via IDCS OAuth

The vault exposes crypto_endpoint and management_endpoint outputs that
downstream OciKmsKey resources consume for key creation and crypto ops.

Excluded from v1:
  - restore_from_file / restore_from_object_store / restore_trigger --
    operational restore, not a deployment concern
  - time_of_deletion -- deletion scheduling, operational concern
  - defined_tags, system_tags -- managed by platform
  - freeform_tags -- auto-populated from metadata labels
  - vault replication -- separate resource with independent lifecycle

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.compartmentId` | `string \| valueFrom` | yes |  | OciCompartment (`status.outputs.compartment_id`) |
| `spec.displayName` | `string` |  |  |  |
| `spec.vaultType` | `enum` |  |  |  |
| `spec.externalKeyManagerMetadata` | `ExternalKeyManagerMetadata` |  |  |  |
| `spec.externalKeyManagerMetadata.externalVaultEndpointUrl` | `string` | yes |  |  |
| `spec.externalKeyManagerMetadata.oauthMetadata` | `OAuthMetadata` | yes |  |  |
| `spec.externalKeyManagerMetadata.oauthMetadata.clientAppId` | `string` | yes |  |  |
| `spec.externalKeyManagerMetadata.oauthMetadata.clientAppSecret` | `string` (sensitive) | yes |  |  |
| `spec.externalKeyManagerMetadata.oauthMetadata.idcsAccountNameUrl` | `string` | yes |  |  |
| `spec.externalKeyManagerMetadata.privateEndpointId` | `string` | yes |  |  |

## Field Details

### spec.compartmentId

`string | valueFrom` · required

OCID of the compartment where the vault will be created.

- references: OciCompartment (`status.outputs.compartment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciCompartment, name: <that resource's name>, fieldPath: status.outputs.compartment_id}} -- a bare string does not parse

### spec.displayName

`string`

Display name for the vault. When omitted, the metadata name is used.

### spec.vaultType

`enum`

Type of vault to create. Cannot be changed after creation.
  - default_vault:  shared HSM partition (lower cost)
  - virtual_private: dedicated HSM partition (higher throughput)
  - external:        external key manager via IDCS OAuth (BYOK/EKMS)

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `unspecified`
- `default_vault` -- Shared HSM partition. Lower cost, suitable for most workloads.
- `virtual_private` -- Dedicated HSM partition. Higher throughput limits, required for high-volume cryptographic operations.
- `external` -- External key manager (BYOK/EKMS). Keys are stored and managed outside OCI using a third-party HSM connected via IDCS OAuth.

### spec.externalKeyManagerMetadata

`ExternalKeyManagerMetadata`

Configuration for an external key manager. Required when vault_type
is external; must not be set otherwise. All sub-fields within this
message are immutable after creation (ForceNew).

### spec.externalKeyManagerMetadata.externalVaultEndpointUrl

`string` · required

URI of the vault on the external key manager system.

- rule: {"string":{"minLen":"1"}}

### spec.externalKeyManagerMetadata.oauthMetadata

`OAuthMetadata` · required

OAuth credentials for authenticating with Oracle IDCS to reach
the external key manager.

- rule: {"required":true}

### spec.externalKeyManagerMetadata.oauthMetadata.clientAppId

`string` · required

Application ID of the client app registered in IDCS.

- rule: {"string":{"minLen":"1"}}

### spec.externalKeyManagerMetadata.oauthMetadata.clientAppSecret

`string` · required · sensitive

Secret of the client app registered in IDCS. This value is
sensitive and will not be returned by the API after creation.

- rule: {"string":{"minLen":"1"}}

### spec.externalKeyManagerMetadata.oauthMetadata.idcsAccountNameUrl

`string` · required

Base URL of the IDCS account (e.g., "https://idcs-xxx.identity.oraclecloud.com").

- rule: {"string":{"minLen":"1"}}

### spec.externalKeyManagerMetadata.privateEndpointId

`string` · required

OCID of the KMS private endpoint used to connect to the external
key manager over a private network path.

- rule: {"string":{"minLen":"1"}}

## Validation Rules

- `external_requires_metadata`: external_key_manager_metadata is required when vault_type is external
- `non_external_forbids_metadata`: external_key_manager_metadata must not be set when vault_type is not external

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OciKmsVault, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.vault_id` | `string` | OCID of the KMS vault. |
| `status.outputs.crypto_endpoint` | `string` | Service endpoint for cryptographic operations (encrypt, decrypt, sign, verify). |
| `status.outputs.management_endpoint` | `string` | Service endpoint for key management operations (create, import, rotate keys). |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.compartmentId` | OciCompartment | `status.outputs.compartment_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| OciKmsKey | `spec.managementEndpoint` | `status.outputs.management_endpoint` |
| OciVaultSecret | `spec.vaultId` | `status.outputs.vault_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
