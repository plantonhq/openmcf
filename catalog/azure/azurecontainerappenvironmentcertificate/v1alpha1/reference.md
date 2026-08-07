# AzureContainerAppEnvironmentCertificate

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureContainerAppEnvironmentCertificateSpec** defines a
bring-your-own TLS certificate stored on a Container App Environment.
Certificates live on the ENVIRONMENT and are shared by every app in it:
upload the certificate once, then bind it to any app's custom domain
through AzureContainerAppCustomDomain's certificate reference.

The certificate arrives one of two ways -- set exactly one:

1. **Inline PFX** -- `certificate_blob_base64` + `certificate_password`:
   a base64-encoded PKCS#12 bundle carrying the certificate chain and
   its private key. Rotation is manual (upload the renewed PFX).
2. **Key Vault reference** -- `certificate_key_vault`: the environment
   pulls the certificate from Azure Key Vault and FOLLOWS renewals
   automatically when the reference is versionless. The environment
   authenticates to the vault with a managed identity that needs read
   access to the vault's secrets (Key Vault Secrets User under RBAC).

For certificates Azure should provision and renew end to end (free,
domain-validated), use AzureContainerAppEnvironmentManagedCertificate
instead -- this kind is for certificates you bring: EV/OV chains,
org-mandated CAs, or wildcard certificates.

**ForceNew fields**: everything except `tags` replaces the certificate
(Azure only allows tag updates in place). Replacing a certificate that
a custom domain is bound to briefly re-binds the domain.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureContainerAppEnvironmentCertificate
metadata:
  name: test-env-certificate
spec:
  certificate_name: app.example.com
  container_app_environment_id:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.App/managedEnvironments/test-env
  certificate_blob_base64: dGVzdC1jZXJ0aWZpY2F0ZS1wZngtYmxvYg==
  certificate_password: test-pfx-password
  tags:
    team: platform
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.certificateName` | `string` | yes |  |  |
| `spec.containerAppEnvironmentId` | `string \| valueFrom` | yes |  | AzureContainerAppEnvironment (`status.outputs.environment_id`) |
| `spec.certificateBlobBase64` | `string` (sensitive) |  |  |  |
| `spec.certificatePassword` | `string` (sensitive) |  |  |  |
| `spec.certificateKeyVault` | `AzureContainerAppEnvironmentCertificateKeyVault` |  |  |  |
| `spec.certificateKeyVault.keyVaultSecretId` | `string \| valueFrom` | yes |  | AzureKeyVaultCertificate (`status.outputs.versionless_secret_id`) |
| `spec.certificateKeyVault.identity` | `string \| valueFrom` |  |  | AzureUserAssignedIdentity (`status.outputs.identity_id`) |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.certificateName

`string` · required

The certificate's name on the environment -- how custom-domain
bindings and the portal refer to it. Lowercase letters, digits,
hyphens, and dots (commonly the domain name, e.g. "app.example.com"
or "wildcard-example-com"). Changing it replaces the certificate.

- rule: Certificate name must use lowercase letters, digits, hyphens, and dots, start and end with a letter or digit, and avoid consecutive hyphens -- e.g. app.example.com or wildcard-example-com
- rule: {"required":true}

### spec.containerAppEnvironmentId

`string | valueFrom` · required

The Container App Environment this certificate is stored on, by ARM
ID. References an AzureContainerAppEnvironment's environment_id
output. Changing it replaces the certificate.

- references: AzureContainerAppEnvironment (`status.outputs.environment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureContainerAppEnvironment, name: <that resource's name>, fieldPath: status.outputs.environment_id}} -- a bare string does not parse

### spec.certificateBlobBase64

`string` · sensitive

The certificate as a base64-encoded PKCS#12 (PFX) bundle -- the
certificate chain plus its private key. The blob bundles the private
key, so it is handled as a secret. Set together with
certificate_password; mutually exclusive with certificate_key_vault.

Note: Azure never returns the blob on reads, so drift in this field
is invisible to the platform -- rotating the certificate means
updating this value (which replaces the certificate resource).

### spec.certificatePassword

`string` · sensitive

The password protecting the PFX blob. Leave empty for a
passwordless PFX -- both engines still send the (empty) password
argument Azure expects alongside the blob. Only meaningful with
certificate_blob_base64.

### spec.certificateKeyVault

`AzureContainerAppEnvironmentCertificateKeyVault`

Pull the certificate from Azure Key Vault instead of uploading a
blob. The environment reads the certificate's SECRET face from the
vault and keeps it current -- pair with a versionless reference so
renewals propagate automatically. Mutually exclusive with
certificate_blob_base64.

### spec.certificateKeyVault.keyVaultSecretId

`string | valueFrom` · required

The Key Vault SECRET URL of the certificate -- a certificate's secret
face is what carries its full PFX material. Reference an
AzureKeyVaultCertificate's versionless_secret_id output
(https://{vault}.vault.azure.net/secrets/{name}) so the environment
follows renewals; a versioned URL pins that version forever.

- references: AzureKeyVaultCertificate (`status.outputs.versionless_secret_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureKeyVaultCertificate, name: <that resource's name>, fieldPath: status.outputs.versionless_secret_id}} -- a bare string does not parse

### spec.certificateKeyVault.identity

`string | valueFrom`

The managed identity the environment uses to read the vault secret:
the literal "System" for the environment's system-assigned identity,
or a user-assigned identity's ARM ID -- reference an
AzureUserAssignedIdentity's identity_id output. Unset deploys
"System" on both engines (Azure's own default). Whichever identity is
named here must already be ON the environment's identity block and
hold read access to the vault's secrets (Key Vault Secrets User under
RBAC) -- Azure checks that at deploy time, across two resources, so
it cannot be validated here.

- references: AzureUserAssignedIdentity (`status.outputs.identity_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureUserAssignedIdentity, name: <that resource's name>, fieldPath: status.outputs.identity_id}} -- a bare string does not parse

### spec.tags

`map<string, string>`

Free-form Azure resource tags applied to the certificate, merged over
the platform's metadata-derived tags (user tags win on key collision).
The only field Azure updates in place.

## Validation Rules

- `azure_container_app_environment_certificate_source`: Provide the certificate exactly one way: either certificate_blob_base64 (with certificate_password) for an inline PFX upload, or certificate_key_vault to pull it from Azure Key Vault
- `azure_container_app_environment_certificate_password_requires_blob`: certificate_password only applies to an inline PFX -- set it together with certificate_blob_base64, or remove it when sourcing the certificate from Key Vault

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureContainerAppEnvironmentCertificate, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.certificate_id` | `string` | The Azure Resource Manager ID of the certificate. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.App/managedEnvironments/{env}/certificates/{name} Consumed by AzureContainerAppCustomDomain to bind the certificate to an app's custom domain. |
| `status.outputs.subject_name` | `string` | The certificate's subject name (e.g. "CN=app.example.com") as read back from the uploaded material. |
| `status.outputs.issuer` | `string` | The certificate's issuer (e.g. "CN=R11, O=Let's Encrypt, C=US"). |
| `status.outputs.issue_date` | `string` | When the certificate was issued (RFC 3339 timestamp). |
| `status.outputs.expiration_date` | `string` | When the certificate expires (RFC 3339 timestamp) -- the value to alarm on for inline-PFX certificates, whose rotation is manual. |
| `status.outputs.thumbprint` | `string` | The certificate's SHA-1 thumbprint -- the fingerprint Azure and operational tooling use to identify the installed certificate. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.containerAppEnvironmentId` | AzureContainerAppEnvironment | `status.outputs.environment_id` |
| `spec.certificateKeyVault.keyVaultSecretId` | AzureKeyVaultCertificate | `status.outputs.versionless_secret_id` |
| `spec.certificateKeyVault.identity` | AzureUserAssignedIdentity | `status.outputs.identity_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureContainerAppCustomDomain | `spec.containerAppEnvironmentCertificateId` | `status.outputs.certificate_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
