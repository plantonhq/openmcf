# AzureKeyVaultKey

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureKeyVaultKeySpec** defines the configuration for creating a
cryptographic key inside an Azure Key Vault -- the customer-managed-key
(CMK) building block of an Azure platform.

A Key Vault key is asymmetric key material the vault guards: the private
part never leaves the vault (or its HSM on the Premium tier); consumers
call the vault to encrypt/decrypt, wrap/unwrap, or sign/verify. This is
how "bring your own key" works across Azure -- Storage accounts, disk
encryption sets, container registries, and database services all encrypt
their data with a data-encryption key that THIS key wraps, which is why
revoking or rotating it revokes access to everything downstream.

Composition seams: the parent vault is a first-class AzureKeyVault
(referenced by key_vault_id); CMK consumers reference this key's
versionless_id output so Azure services follow rotation automatically
(pin key_id only when a compliance regime demands a frozen version). The
deploying credential needs data-plane key permissions on the vault (the
"Key Vault Administrator" or "Key Vault Crypto Officer" RBAC role, or
key permissions in a legacy access policy).

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureKeyVaultKey
metadata:
  name: test-kv-key
  org: test-org
  env: dev
spec:
  name: test-cmk
  key_vault_id:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.KeyVault/vaults/test-kv
  key_type: RSA_HSM
  key_size: 2048
  key_opts:
    - WRAP_KEY
    - UNWRAP_KEY
    - ENCRYPT
    - DECRYPT
  not_before_date: "2027-01-01T00:00:00Z"
  expiration_date: "2029-01-01T00:00:00Z"
  rotation_policy:
    expire_after: P2Y
    notify_before_expiry: P30D
    automatic:
      time_before_expiry: P60D
  tags:
    purpose: hack-testing
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.name` | `string` | yes |  |  |
| `spec.keyVaultId` | `string \| valueFrom` | yes |  | AzureKeyVault (`status.outputs.key_vault_id`) |
| `spec.keyType` | `enum` |  |  |  |
| `spec.keySize` | `int32` |  |  |  |
| `spec.curve` | `enum` |  |  |  |
| `spec.keyOpts` | `[]enum` | yes |  |  |
| `spec.notBeforeDate` | `string` |  |  |  |
| `spec.expirationDate` | `string` |  |  |  |
| `spec.rotationPolicy` | `AzureKeyVaultKeyRotationPolicy` |  |  |  |
| `spec.rotationPolicy.expireAfter` | `string` |  |  |  |
| `spec.rotationPolicy.notifyBeforeExpiry` | `string` |  |  |  |
| `spec.rotationPolicy.automatic` | `AzureKeyVaultKeyRotationPolicyAutomatic` |  |  |  |
| `spec.rotationPolicy.automatic.timeAfterCreation` | `string` |  |  |  |
| `spec.rotationPolicy.automatic.timeBeforeExpiry` | `string` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.name

`string` · required

The key's name within the vault: 1-127 characters of letters, digits,
and hyphens, unique among the vault's keys. Changing the name replaces
the key. A deleted key's name stays reserved for the vault's
soft-delete retention window unless purged.

- rule: {"required":true,"string":{"pattern":"^[0-9a-zA-Z-]{1,127}$"}}

### spec.keyVaultId

`string | valueFrom` · required

The vault the key lives in, by ARM resource ID. Defaults to
referencing an AzureKeyVault's key_vault_id output in composed
environments. Changing the vault replaces the key -- key material
never moves between vaults.

- references: AzureKeyVault (`status.outputs.key_vault_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureKeyVault, name: <that resource's name>, fieldPath: status.outputs.key_vault_id}} -- a bare string does not parse

### spec.keyType

`enum`

The key's algorithm family, fixed at creation. RSA is the
general-purpose choice every Azure CMK integration accepts; EC keys
sign and verify (JWTs, code signing) but do not encrypt. The _HSM
variants keep the private key inside a FIPS 140-2 Level 3 hardware
module and require the vault's PREMIUM SKU.

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `azure_key_vault_key_type_unspecified` -- Not specified -- invalid; pick the key's algorithm family explicitly.
- `RSA` -- RSA, software-protected -- the general-purpose choice every Azure CMK integration accepts.
- `RSA_HSM` -- RSA, HSM-protected (FIPS 140-2 Level 3). Requires the vault's PREMIUM SKU.
- `EC` -- Elliptic curve, software-protected -- signing and verification workloads.
- `EC_HSM` -- Elliptic curve, HSM-protected. Requires the vault's PREMIUM SKU.

### spec.keySize

`int32` · optional (explicit presence)

For RSA / RSA_HSM keys: the modulus size in bits. Key Vault supports
2048 (the baseline every integration accepts), 3072, and 4096. Fixed
at creation. Must be unset for EC keys (set curve instead).

- rule: {"int32":{"in":[2048,3072,4096]}}

### spec.curve

`enum`

For EC / EC_HSM keys: the elliptic curve. Unspecified lets Azure
default to P-256. Fixed at creation. Must be unset for RSA keys (set
key_size instead).

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_key_vault_key_curve_unspecified` -- Not specified: Azure defaults to P-256.
- `P_256` -- NIST P-256 (secp256r1) -- the interoperable default.
- `P_256K` -- SECG secp256k1 -- the curve used by blockchain ecosystems.
- `P_384` -- NIST P-384 (secp384r1).
- `P_521` -- NIST P-521 (secp521r1).

### spec.keyOpts

`[]enum` · required

The cryptographic operations this key may perform -- Azure rejects any
operation not listed here, so the list is the key's capability
boundary. CMK/envelope encryption needs WRAP_KEY + UNWRAP_KEY; direct
encryption needs ENCRYPT + DECRYPT; signing needs SIGN + VERIFY.
Grant only what the key's consumers actually perform.

- rule: {"repeated":{"minItems":"1","items":{"enum":{"definedOnly":true,"notIn":[0]}}}}

Allowed values (use exactly as shown):

- `azure_key_vault_key_operation_unspecified` -- Not specified -- invalid; list the operations explicitly.
- `DECRYPT` -- Decrypt ciphertext with the private key.
- `ENCRYPT` -- Encrypt plaintext with the public key.
- `SIGN` -- Sign digests with the private key.
- `UNWRAP_KEY` -- Unwrap (decrypt) a wrapped data-encryption key -- the operation every CMK consumer performs at startup.
- `VERIFY` -- Verify signatures with the public key.
- `WRAP_KEY` -- Wrap (encrypt) a data-encryption key.

### spec.notBeforeDate

`string` · optional (explicit presence)

The UTC instant before which the key must not be used, RFC 3339
format (e.g. "2027-01-01T00:00:00Z"). Leave unset for a key usable
immediately.

- rule: {"string":{"pattern":"^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}(\\.\\d+)?Z$"}}

### spec.expirationDate

`string` · optional (explicit presence)

The UTC instant the key expires, RFC 3339 format. Expired keys refuse
cryptographic operations. Prefer rotation_policy over hard expiry for
keys that encrypt live data -- an expired CMK takes its dependents
down with it. Once set, expiry cannot be fully unset on the underlying
key even across delete/recreate (Azure restores purged names' state).

- rule: {"string":{"pattern":"^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}(\\.\\d+)?Z$"}}

### spec.rotationPolicy

`AzureKeyVaultKeyRotationPolicy`

Automatic rotation policy: when Azure mints a new key version on its
own, and when it notifies (Event Grid) before expiry. Rotation is why
consumers should reference versionless_id -- each rotation creates a
new version and versionless references follow it with zero
intervention.

- rule: a rotation policy needs expire_after and/or an automatic block -- an empty policy configures nothing
- rule: expire_after and notify_before_expiry go together -- set both or neither

### spec.rotationPolicy.expireAfter

`string` · optional (explicit presence)

How long each key version lives before it expires, as an ISO 8601
duration between P28D and P100Y (e.g. "P90D" for 90 days, "P2Y" for
two years). Sets the expiry attribute on every newly created version.
Must be set together with notify_before_expiry.

- rule: {"string":{"pattern":"^P((\\d+Y)(\\d+M)?(\\d+D)?|(\\d+M)(\\d+D)?|(\\d+D))$"}}

### spec.rotationPolicy.notifyBeforeExpiry

`string` · optional (explicit presence)

How far before a version expires Azure raises the near-expiry Event
Grid notification, as an ISO 8601 duration of at least P7D. Must be
set together with expire_after.

- rule: {"string":{"pattern":"^P((\\d+Y)(\\d+M)?(\\d+D)?|(\\d+M)(\\d+D)?|(\\d+D))$"}}

### spec.rotationPolicy.automatic

`AzureKeyVaultKeyRotationPolicyAutomatic`

When Azure automatically rotates to a new version. Without this block
the policy only sets per-version expiry and notification -- rotation
itself stays manual.

- rule: the automatic block needs a trigger -- time_after_creation or time_before_expiry

### spec.rotationPolicy.automatic.timeAfterCreation

`string` · optional (explicit presence)

Rotate this long after each version's creation, as an ISO 8601
duration (e.g. "P83D" to rotate ~every quarter with a week of slack).

- rule: {"string":{"pattern":"^P((\\d+Y)(\\d+M)?(\\d+D)?|(\\d+M)(\\d+D)?|(\\d+D))$"}}

### spec.rotationPolicy.automatic.timeBeforeExpiry

`string` · optional (explicit presence)

Rotate this long before each version's expiry, as an ISO 8601
duration (requires expire_after on the policy so there is an expiry
to count back from).

- rule: {"string":{"pattern":"^P((\\d+Y)(\\d+M)?(\\d+D)?|(\\d+M)(\\d+D)?|(\\d+D))$"}}

### spec.tags

`map<string, string>`

Free-form tags applied to the key, merged over the Planton-derived
resource tags (organization, environment, resource id); a user tag
with the same key wins. Updatable in place.

## Validation Rules

- `key_vault_key_rsa_requires_size`: RSA and RSA_HSM keys need key_size (2048, 3072, or 4096) -- pick 2048 unless a compliance regime demands larger
- `key_vault_key_ec_forbids_size`: key_size applies only to RSA keys -- an EC key's strength is set by its curve
- `key_vault_key_rsa_forbids_curve`: curve applies only to EC keys -- an RSA key's strength is set by key_size

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureKeyVaultKey, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.key_id` | `string` | The key's versioned data-plane ID: https://{vault}.vault.azure.net/keys/{name}/{version}. Pins consumers to THIS version -- rotation does not follow. Reference it only when a compliance regime demands a frozen key version. |
| `status.outputs.versionless_id` | `string` | The key's versionless data-plane ID: https://{vault}.vault.azure.net/keys/{name}. The reference every customer-managed-key consumer should use -- Azure services resolve it to the current version, so rotation propagates with zero intervention. |
| `status.outputs.key_name` | `string` | The key's name within the vault. |
| `status.outputs.version` | `string` | The current version identifier (the trailing segment of key_id). |
| `status.outputs.resource_id` | `string` | The key's versioned ARM resource ID (.../vaults/{vault}/keys/{name}/versions/{version}) -- the control-plane identity used by ARM-level integrations and RBAC scopes. |
| `status.outputs.resource_versionless_id` | `string` | The key's versionless ARM resource ID (.../vaults/{vault}/keys/{name}). |
| `status.outputs.public_key_pem` | `string` | The public half of the key in PEM form -- consumable by anything that verifies signatures or encrypts to this key outside Azure. |
| `status.outputs.public_key_openssh` | `string` | The public half of the key in OpenSSH form (RSA and P-256/P-384/P-521 EC keys) -- usable directly in authorized_keys. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.keyVaultId` | AzureKeyVault | `status.outputs.key_vault_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureAiFoundry | `spec.encryption.keyId` | `status.outputs.key_id` |
| AzureAksCluster | `spec.keyManagementService.keyVaultKeyId` | `status.outputs.key_id` |
| AzureCognitiveAccount | `spec.customerManagedKey.keyVaultKeyId` | `status.outputs.versionless_id` |
| AzureContainerRegistry | `spec.encryption.keyVaultKeyId` | `status.outputs.versionless_id` |
| AzureCosmosdbAccount | `spec.keyVaultKeyId` | `status.outputs.versionless_id` |
| AzureDiskEncryptionSet | `spec.keyVaultKeyId` | `status.outputs.versionless_id` |
| AzureEventHubNamespaceCustomerManagedKey | `spec.keyVaultKeyIds` | `status.outputs.versionless_id` |
| AzureMachineLearningWorkspace | `spec.encryption.keyId` | `status.outputs.versionless_id` |
| AzureManagedRedis | `spec.customerManagedKey.keyVaultKeyId` | `status.outputs.key_id` |
| AzureMssqlDatabase | `spec.transparentDataEncryptionKeyVaultKeyId` | `status.outputs.key_id` |
| AzureMssqlServer | `spec.transparentDataEncryptionKeyVaultKeyId` | `status.outputs.key_id` |
| AzureMysqlFlexibleServer | `spec.customerManagedKey.keyVaultKeyId` | `status.outputs.versionless_id` |
| AzureMysqlFlexibleServer | `spec.customerManagedKey.geoBackupKeyVaultKeyId` | `status.outputs.versionless_id` |
| AzurePostgresqlFlexibleServer | `spec.customerManagedKey.keyVaultKeyId` | `status.outputs.versionless_id` |
| AzurePostgresqlFlexibleServer | `spec.customerManagedKey.geoBackupKeyVaultKeyId` | `status.outputs.versionless_id` |
| AzureServiceBusNamespace | `spec.customerManagedKey.keyVaultKeyId` | `status.outputs.versionless_id` |
| AzureStorageAccount | `spec.customerManagedKey.keyVaultKeyId` | `status.outputs.versionless_id` |
| AzureStorageEncryptionScope | `spec.keyVaultKeyId` | `status.outputs.versionless_id` |

## See Also

- [Overview](../README.md)
