# AzureContainerAppEnvironmentCertificate -- Design Research

## The Resource

A bring-your-own certificate on a Container App Environment
(`Microsoft.App/managedEnvironments/certificates`). The component maps
onto `azurerm_container_app_environment_certificate` (azurerm v4.x,
`internal/services/containerapps/container_app_environment_certificate_resource.go`,
Microsoft.App API 2025-07-01), parity-verified against pulumi-azure v6
(`containerapp.EnvironmentCertificate`) -- zero bridge lag.

## Field Mapping (azurerm -> spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `certificate_name` | The provider's certificate-name validator (lowercase alnum/hyphen/dot, no leading/trailing hyphen, no `--`) as a CEL. ForceNew |
| `container_app_environment_id` | `container_app_environment_id` | FK to `AzureContainerAppEnvironment.environment_id`. ForceNew |
| `certificate_blob_base64` | `certificate_blob_base64` | Sensitive (bundles the private key; azurerm leaves it unflagged but never reads it back -- modeled sensitive regardless). ForceNew |
| `certificate_password` | `certificate_password` | Sensitive. Provider `RequiredWith` the blob, front-loaded as a pairing CEL. ForceNew |
| `certificate_key_vault.key_vault_secret_id` | `certificate_key_vault.key_vault_secret_id` | Versionless secret URL; FK-defaults to `AzureKeyVaultCertificate.versionless_secret_id`. ForceNew |
| `certificate_key_vault.identity` | `certificate_key_vault.identity` | "System" (provider default, sent explicitly by both engines when unset) or a UAI ARM id (FK default). ForceNew |
| `tags` | `tags` | The only in-place update |

Computed attributes exported as outputs: `subject_name`, `issuer`,
`issue_date`, `expiration_date`, `thumbprint`.

## Decomposition Decision

Separate from `AzureContainerAppEnvironmentManagedCertificate` -- two
different ARM types (`certificates` vs `managedCertificates`) with
disjoint inputs (uploaded material vs a domain to validate) and opposite
lifecycles (synchronous PUT vs an operation that blocks on public-DNS
domain validation). Unifying them would be an invented abstraction; the
polymorphic-dispatch pattern applies only to byte-identical ARM types.

## Front-Loaded Contracts

- **Inline PFX XOR Key Vault** (message CEL) -- the provider's
  ExactlyOneOf.
- **Blob and password travel together** (message CEL) -- the provider's
  RequiredWith in both directions.
- The identity-must-ride-the-environment and vault-read-permission
  requirements are cross-resource contracts Azure checks at deploy time;
  documented on the spec field.

## Recorded Skips (with reasons)

Nothing skipped: the azurerm surface is exactly the fields above.

## Operational Behavior Worth Knowing

- Azure never returns the PFX blob on reads, so blob drift is invisible
  to refresh -- rotation is a spec update (which replaces the resource).
- Replacing a certificate a custom domain is bound to briefly re-binds
  the domain.
- Container Apps accepts self-signed uploads (useful for pre-production
  testing); browsers will still reject them, of course.
