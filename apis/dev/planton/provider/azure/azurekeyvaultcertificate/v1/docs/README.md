# AzureKeyVaultCertificate -- Design Research

## The Resource

A Key Vault certificate is a data-plane object
(`https://{vault}.vault.azure.net/certificates/{name}`) that Azure models
as three linked faces sharing one name: the certificate (public part +
policy), a private KEY, and a SECRET whose value is the full bundle. The
component maps onto `azurerm_key_vault_certificate` (azurerm v4.x,
`internal/services/keyvault/key_vault_certificate_resource.go`),
parity-verified against pulumi-azure v6 (`keyvault.Certificate`).

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `name` | Required, ForceNew, 1-127 alphanumeric/hyphens |
| `key_vault_id` | `key_vault_id` | FK → AzureKeyVault `key_vault_id` output; ForceNew |
| `certificate.contents` / `.password` | `certificate.contents` / `.password` | BOTH `(sensitive)` -- the bundle carries the private key |
| `certificate_policy.issuer_parameters.name` | `certificate_policy.issuer_name` | "Self" / "Unknown" / CA issuer name; free string (issuer names are vault-configured) |
| `certificate_policy.key_properties` | `key_properties` | exportable, key_type enum (incl. `oct` -- present in Azure's surface, no real X.509 use), key_size (RSA and EC sizes), curve enum, reuse_key |
| `certificate_policy.lifetime_action` | `lifetime_actions` | action_type enum + trigger (days_before_expiry XOR lifetime_percentage) |
| `certificate_policy.secret_properties.content_type` | `secret_properties.content_type` | Closed enum PKCS12/PEM mapped to the `application/x-pkcs12` / `application/x-pem-file` media types (azurerm models the raw string; the enum removes the misspelling trap while covering the only two values Azure issues) |
| `certificate_policy.x509_certificate_properties` | `x509_certificate_properties` | subject, SANs (dns/emails/upns, at-least-one when present), key_usage closed enum set (min 1, azurerm Required), extended_key_usage OIDs, validity_in_months |
| `tags` | `tags` | User tags merged over Planton-derived tags |
| `secret_id` / `versionless_secret_id` / `versionless_id` / `version` / `thumbprint` / `resource_manager_id` / `resource_manager_versionless_id` (computed) | outputs | The composition surface |

## Design Decisions

- **Import and generate are an AT-LEAST-ONE, not an XOR** -- exactly
  azurerm's contract (`AtLeastOneOf`): an imported bundle may carry an
  explicit policy governing renewal and exportability; with no policy,
  Azure derives one from the bundle. The second CEL requires
  `x509_certificate_properties` when generating (Azure cannot issue
  without a subject) while leaving imports free to omit it.
- **The lifetime-action trigger is exactly-one** (days_before_expiry XOR
  lifetime_percentage): Azure's certificate-policy API accepts one trigger
  per action and errors on both -- the CEL mirrors the service contract
  azurerm leaves to deploy time.
- **`content_type` promoted from free string to a closed enum**: Azure
  issues exactly two secret-face encodings; the media-type strings are
  wire detail the vocabularies own.
- **`issuer_name` stays a free string** (required, min 1): beyond the
  `Self`/`Unknown` sentinels it names a CA issuer configured on the vault
  -- an open, user-defined vocabulary no enum can close.
- **Wire vocabularies are exhaustive maps in both engines** -- note the
  three distinct casings Azure is sensitive about: UpperCamel action types
  (`AutoRenew`), lowerCamel key usages (`digitalSignature`), and
  hyphenated key types (`RSA-HSM`, lowercase `oct`).
- **Outputs export all three faces** (certificate, secret, ARM proxy) in
  both versioned and versionless forms; the secret face is the seam TLS
  terminators consume (App Gateway takes `versionless_secret_id` to follow
  renewals).

## Recorded Skips (with reasons)

- **`azurerm_key_vault_certificate_contacts`** -- vault-wide
  expiry-notification contacts (one list per vault, not per certificate);
  joins the adoption backlog with the vault's other data-plane
  side-resources.
- **`azurerm_key_vault_certificate_issuer`** -- the CA-integration
  configuration (provider account credentials for DigiCert/GlobalSign).
  Niche, carries provider secrets, and certificates reference issuers by
  NAME (the `issuer_name` field composes with it once demanded); adoption
  backlog.

## Operational Behavior Worth Knowing

- **Data-plane creation**: the deploying credential needs certificate
  permissions on the vault (Key Vault Administrator / Certificates
  Officer); subscription Owner alone 403s.
- **Issuer `Self` completes synchronously**; a CA issuer keeps the
  operation pending until the CA responds -- expect longer creates on
  integrated-CA policies. Issuer `Unknown` parks the CSR for out-of-band
  signing, and its renewals are inherently manual (`EMAIL_CONTACTS` is
  the only sensible action).
- **Renewal creates a new version** with a fresh key (unless
  `reuse_key: true`); versionless consumers follow automatically.
- **Changing any policy part except lifetime_actions creates a new
  version**; changing `certificate.contents` imports a new version.
- **Soft delete + purge**: a deleted certificate's name stays reserved in
  the vault for the retention window; the providers purge soft-deleted
  certificates on destroy by default.

## Composition

- `key_vault_id` → `AzureKeyVault.status.outputs.key_vault_id`
- `versionless_secret_id` output is consumed by:
  - `AzureApplicationGateway.ssl_certificates[].key_vault_secret_id`
    (the gateway polls the secret face and follows renewals; its identity
    needs secret-GET on the vault)
- `thumbprint` output serves integrations that pin fingerprints
- The deployer's data-plane grant composes as `AzureRoleAssignment`
  ("Key Vault Certificates Officer") scoped at the vault's `key_vault_id`
