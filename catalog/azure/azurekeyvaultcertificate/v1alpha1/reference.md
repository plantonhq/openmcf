# AzureKeyVaultCertificate

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureKeyVaultCertificateSpec** defines the configuration for creating an
X.509 certificate inside an Azure Key Vault -- the TLS building block the
vault manages end to end: private key custody, enrollment or import,
renewal, and expiry notification.

A vault certificate is really three linked objects sharing one name: the
certificate (public part + policy), its private KEY, and a SECRET whose
value is the full bundle (cert + key). That third face is the composition
seam -- Azure services that terminate TLS (Application Gateway listeners,
App Service custom domains) consume the certificate THROUGH its secret
identifier, which is why the secret_id / versionless_secret_id outputs
exist and are what downstream kinds reference.

Two ways to get a certificate in, combinable per azurerm's contract:
- **Generate** (certificate_policy only): the vault creates the key and
  either self-signs (issuer "Self") or forwards a CSR to a configured CA.
  Self-signed + auto-renew is the fully-hands-off shape for internal TLS.
- **Import** (certificate): bring an existing PFX/PEM bundle; a policy may
  accompany it to govern the imported material (exportability, renewal
  actions). Without an explicit policy Azure derives one from the bundle.

The deploying credential needs data-plane certificate permissions on the
vault (the "Key Vault Administrator" or "Key Vault Certificates Officer"
RBAC role, or certificate permissions in a legacy access policy).

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureKeyVaultCertificate
metadata:
  name: test-kv-cert
  org: test-org
  env: dev
spec:
  name: test-tls
  key_vault_id:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.KeyVault/vaults/test-kv
  certificate_policy:
    issuer_name: Self
    key_properties:
      exportable: true
      key_type: EC_HSM
      curve: P_384
      reuse_key: true
    lifetime_actions:
      - action_type: AUTO_RENEW
        trigger:
          lifetime_percentage: 80
      - action_type: EMAIL_CONTACTS
        trigger:
          days_before_expiry: 14
    secret_properties:
      content_type: PEM
    x509_certificate_properties:
      subject: CN=hack.example.com
      subject_alternative_names:
        dns_names:
          - hack.example.com
          - alt.example.com
        emails:
          - admin@example.com
        upns:
          - user@example.com
      key_usage:
        - DIGITAL_SIGNATURE
        - KEY_ENCIPHERMENT
        - KEY_AGREEMENT
      extended_key_usage:
        - 1.3.6.1.5.5.7.3.1
      validity_in_months: 12
  tags:
    purpose: hack-testing
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.name` | `string` | yes |  |  |
| `spec.keyVaultId` | `string \| valueFrom` | yes |  | AzureKeyVault (`status.outputs.key_vault_id`) |
| `spec.certificate` | `AzureKeyVaultCertificateImport` |  |  |  |
| `spec.certificate.contents` | `string` (sensitive) | yes |  |  |
| `spec.certificate.password` | `string` (sensitive) |  |  |  |
| `spec.certificatePolicy` | `AzureKeyVaultCertificatePolicy` |  |  |  |
| `spec.certificatePolicy.issuerName` | `string` | yes |  |  |
| `spec.certificatePolicy.keyProperties` | `AzureKeyVaultCertificateKeyProperties` | yes |  |  |
| `spec.certificatePolicy.keyProperties.exportable` | `bool` |  |  |  |
| `spec.certificatePolicy.keyProperties.keyType` | `enum` |  |  |  |
| `spec.certificatePolicy.keyProperties.keySize` | `int32` |  |  |  |
| `spec.certificatePolicy.keyProperties.curve` | `enum` |  |  |  |
| `spec.certificatePolicy.keyProperties.reuseKey` | `bool` |  |  |  |
| `spec.certificatePolicy.lifetimeActions` | `[]AzureKeyVaultCertificateLifetimeAction` |  |  |  |
| `spec.certificatePolicy.lifetimeActions[].actionType` | `enum` |  |  |  |
| `spec.certificatePolicy.lifetimeActions[].trigger` | `AzureKeyVaultCertificateLifetimeTrigger` | yes |  |  |
| `spec.certificatePolicy.lifetimeActions[].trigger.daysBeforeExpiry` | `int32` |  |  |  |
| `spec.certificatePolicy.lifetimeActions[].trigger.lifetimePercentage` | `int32` |  |  |  |
| `spec.certificatePolicy.secretProperties` | `AzureKeyVaultCertificateSecretProperties` | yes |  |  |
| `spec.certificatePolicy.secretProperties.contentType` | `enum` |  |  |  |
| `spec.certificatePolicy.x509CertificateProperties` | `AzureKeyVaultCertificateX509Properties` |  |  |  |
| `spec.certificatePolicy.x509CertificateProperties.subject` | `string` | yes |  |  |
| `spec.certificatePolicy.x509CertificateProperties.subjectAlternativeNames` | `AzureKeyVaultCertificateSubjectAlternativeNames` |  |  |  |
| `spec.certificatePolicy.x509CertificateProperties.subjectAlternativeNames.dnsNames` | `[]string` |  |  |  |
| `spec.certificatePolicy.x509CertificateProperties.subjectAlternativeNames.emails` | `[]string` |  |  |  |
| `spec.certificatePolicy.x509CertificateProperties.subjectAlternativeNames.upns` | `[]string` |  |  |  |
| `spec.certificatePolicy.x509CertificateProperties.keyUsage` | `[]enum` | yes |  |  |
| `spec.certificatePolicy.x509CertificateProperties.extendedKeyUsage` | `[]string` |  |  |  |
| `spec.certificatePolicy.x509CertificateProperties.validityInMonths` | `int32` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.name

`string` · required

The certificate's name within the vault: 1-127 characters of letters,
digits, and hyphens, unique among the vault's certificates. Changing
the name replaces the certificate. A deleted certificate's name stays
reserved for the vault's soft-delete retention window unless purged.

- rule: {"required":true,"string":{"pattern":"^[0-9a-zA-Z-]{1,127}$"}}

### spec.keyVaultId

`string | valueFrom` · required

The vault the certificate lives in, by ARM resource ID. Defaults to
referencing an AzureKeyVault's key_vault_id output in composed
environments. Changing the vault replaces the certificate.

- references: AzureKeyVault (`status.outputs.key_vault_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureKeyVault, name: <that resource's name>, fieldPath: status.outputs.key_vault_id}} -- a bare string does not parse

### spec.certificate

`AzureKeyVaultCertificateImport`

An existing certificate bundle to IMPORT. Changing the contents
imports a new version of the certificate. Omit to have the vault
generate the certificate from certificate_policy instead.

### spec.certificate.contents

`string` · required · sensitive

The base64-encoded certificate bundle: a PFX (PKCS#12) or PEM
containing BOTH the certificate chain and its private key. This is
secret material -- it carries the private key.

- rule: {"required":true}

### spec.certificate.password

`string` · optional (explicit presence) · sensitive

The password protecting the bundle. Leave unset for unprotected PEMs
and passwordless PFX bundles.

### spec.certificatePolicy

`AzureKeyVaultCertificatePolicy`

How the vault generates and manages the certificate: key shape,
issuer, X.509 content, renewal actions. Required when generating;
optional alongside an import (Azure derives a policy from the bundle
when omitted). Changing any part except lifetime_actions creates a
new certificate version.

### spec.certificatePolicy.issuerName

`string` · required

Who issues the certificate:
- "Self": the vault self-signs -- zero external dependencies, the
  shape for internal/service-mesh TLS where consumers trust the cert
  explicitly.
- "Unknown": the vault mints a CSR and holds it for an out-of-band CA
  to sign (the pending-signature flow).
- The name of a CA issuer configured on the vault (DigiCert and
  GlobalSign have first-party integrations) for automated public
  issuance.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.certificatePolicy.keyProperties

`AzureKeyVaultCertificateKeyProperties` · required

The shape of the certificate's private key.

- rule: {"required":true}

### spec.certificatePolicy.keyProperties.exportable

`bool`

Whether the private key may leave the vault (inside the secret face
consumers read). TLS terminators outside Key Vault integrations need
true; keep false when only vault-integrated services (which fetch
through Key Vault references) consume the certificate.

### spec.certificatePolicy.keyProperties.keyType

`enum`

The key's algorithm family. RSA is the universal TLS choice; EC for
modern, smaller handshakes; the _HSM variants require the vault's
PREMIUM SKU; OCT (symmetric) exists in the API surface but has no
real X.509 use.

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `azure_key_vault_certificate_key_type_unspecified` -- Not specified -- invalid; pick the key's algorithm family explicitly.
- `RSA` -- RSA, software-protected -- the universal TLS choice.
- `RSA_HSM` -- RSA, HSM-protected. Requires the vault's PREMIUM SKU.
- `EC` -- Elliptic curve, software-protected.
- `EC_HSM` -- Elliptic curve, HSM-protected. Requires the vault's PREMIUM SKU.
- `OCT` -- Symmetric key -- present in Azure's API surface; no real X.509 use.

### spec.certificatePolicy.keyProperties.keySize

`int32` · optional (explicit presence)

For RSA keys: the modulus size in bits (2048, 3072, or 4096). For EC
keys Azure accepts the curve's own size (256, 384, 521) or derives it
from the curve when unset.

- rule: {"int32":{"in":[256,384,521,2048,3072,4096]}}

### spec.certificatePolicy.keyProperties.curve

`enum`

For EC keys: the elliptic curve. Unspecified lets Azure derive it
(P-256 baseline).

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_key_vault_certificate_key_curve_unspecified` -- Not specified: Azure derives the curve (P-256 baseline).
- `P_256` -- NIST P-256 (secp256r1) -- the interoperable default.
- `P_256K` -- SECG secp256k1.
- `P_384` -- NIST P-384 (secp384r1).
- `P_521` -- NIST P-521 (secp521r1).

### spec.certificatePolicy.keyProperties.reuseKey

`bool`

Whether renewals re-use the existing private key (true) or mint a
fresh one (false). Fresh keys are the safer default posture; re-use
exists for pinned-key deployments.

### spec.certificatePolicy.lifetimeActions

`[]AzureKeyVaultCertificateLifetimeAction`

What the vault does as the certificate approaches expiry. One action
per trigger; the classic shape is a single AUTO_RENEW at 80% lifetime
(or EMAIL_CONTACTS for CA flows that need a human). Updatable without
creating a new certificate version.

### spec.certificatePolicy.lifetimeActions[].actionType

`enum`

What the vault does when the trigger fires. AUTO_RENEW re-enrolls
automatically (self-signed and integrated-CA issuers); EMAIL_CONTACTS
notifies the vault's certificate contacts instead -- the only action
available for "Unknown"-issuer certificates, whose renewal is
inherently manual.

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `azure_key_vault_certificate_lifetime_action_type_unspecified` -- Not specified -- invalid; pick the action explicitly.
- `AUTO_RENEW` -- Automatically re-enroll the certificate.
- `EMAIL_CONTACTS` -- Notify the vault's certificate contacts.

### spec.certificatePolicy.lifetimeActions[].trigger

`AzureKeyVaultCertificateLifetimeTrigger` · required

When the action fires: a fixed number of days before expiry, or a
percentage of the certificate's lifetime elapsed. Azure accepts
exactly one trigger per action.

- rule: {"required":true}
- rule: set exactly one trigger -- days_before_expiry or lifetime_percentage, not both

### spec.certificatePolicy.lifetimeActions[].trigger.daysBeforeExpiry

`int32` · optional (explicit presence)

Fire this many days before the certificate expires.

- rule: {"int32":{"gte":1}}

### spec.certificatePolicy.lifetimeActions[].trigger.lifetimePercentage

`int32` · optional (explicit presence)

Fire when this percentage of the certificate's lifetime has elapsed
(e.g. 80 -- the conventional renewal point).

- rule: {"int32":{"lte":99,"gte":1}}

### spec.certificatePolicy.secretProperties

`AzureKeyVaultCertificateSecretProperties` · required

How the certificate's secret face is encoded -- what consumers reading
the secret get.

- rule: {"required":true}

### spec.certificatePolicy.secretProperties.contentType

`enum`

The media type consumers reading the certificate's secret get:
PKCS12 (a .pfx bundle -- what Application Gateway and Windows-world
consumers expect) or PEM (concatenated text -- the Linux/OpenSSL
world).

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `azure_key_vault_certificate_content_type_unspecified` -- Not specified -- invalid; pick the encoding consumers expect.
- `PKCS12` -- PKCS#12 (.pfx) -- Application Gateway and Windows-world consumers.
- `PEM` -- PEM text -- the Linux/OpenSSL world.

### spec.certificatePolicy.x509CertificateProperties

`AzureKeyVaultCertificateX509Properties`

The X.509 content of the certificate: subject, SANs, usages,
validity. Required when the vault generates the certificate; omit for
imports (derived from the bundle).

### spec.certificatePolicy.x509CertificateProperties.subject

`string` · required

The certificate's subject distinguished name, e.g.
"CN=internal.example.com". For TLS, modern clients validate SANs --
list every hostname in subject_alternative_names and keep the subject
CN as the primary name.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.certificatePolicy.x509CertificateProperties.subjectAlternativeNames

`AzureKeyVaultCertificateSubjectAlternativeNames`

The additional identities the certificate is valid for.

- rule: subject_alternative_names needs at least one entry across dns_names, emails, or upns -- omit the block entirely when the subject alone suffices

### spec.certificatePolicy.x509CertificateProperties.subjectAlternativeNames.dnsNames

`[]string`

DNS names (the usual TLS case), e.g. "internal.example.com".

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.certificatePolicy.x509CertificateProperties.subjectAlternativeNames.emails

`[]string`

Email addresses (S/MIME-style identities).

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.certificatePolicy.x509CertificateProperties.subjectAlternativeNames.upns

`[]string`

User principal names (smart-card / client-auth identities).

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.certificatePolicy.x509CertificateProperties.keyUsage

`[]enum` · required

The X.509 key-usage extensions stamped on the certificate. TLS
server certificates conventionally carry DIGITAL_SIGNATURE +
KEY_ENCIPHERMENT.

- rule: {"repeated":{"minItems":"1","items":{"enum":{"definedOnly":true,"notIn":[0]}}}}

Allowed values (use exactly as shown):

- `azure_key_vault_certificate_key_usage_unspecified` -- Not specified -- invalid; list usages explicitly.
- `CRL_SIGN` -- Sign certificate revocation lists.
- `DATA_ENCIPHERMENT` -- Encrypt user data directly.
- `DECIPHER_ONLY` -- Decipher-only key agreement.
- `DIGITAL_SIGNATURE` -- Verify digital signatures -- half of the conventional TLS-server pair.
- `ENCIPHER_ONLY` -- Encipher-only key agreement.
- `KEY_AGREEMENT` -- Key-agreement protocols (ECDH).
- `KEY_CERT_SIGN` -- Sign other certificates (CA use).
- `KEY_ENCIPHERMENT` -- Encrypt keys in transport -- the other half of the conventional TLS-server pair.
- `NON_REPUDIATION` -- Non-repudiation / content commitment.

### spec.certificatePolicy.x509CertificateProperties.extendedKeyUsage

`[]string`

Extended key usage OIDs, e.g. "1.3.6.1.5.5.7.3.1" (TLS server
authentication) and "1.3.6.1.5.5.7.3.2" (client authentication).
Empty leaves EKU to Azure's defaults.

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.certificatePolicy.x509CertificateProperties.validityInMonths

`int32`

How long each issued certificate is valid, in months (e.g. 12).
Renewals issue for the same validity.

- rule: {"int32":{"lte":1200,"gte":1}}

### spec.tags

`map<string, string>`

Free-form tags applied to the certificate, merged over the
Planton-derived resource tags (organization, environment, resource
id); a user tag with the same key wins. Updatable in place.

## Validation Rules

- `key_vault_certificate_import_or_policy`: provide certificate (to import an existing bundle) and/or certificate_policy (to have the vault generate one) -- with neither there is nothing to create
- `key_vault_certificate_generated_needs_x509`: a vault-generated certificate needs certificate_policy.x509_certificate_properties -- Azure cannot issue without a subject; x509 properties are only derivable when importing a bundle

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureKeyVaultCertificate, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.certificate_id` | `string` | The certificate's versioned data-plane ID: https://{vault}.vault.azure.net/certificates/{name}/{version}. Pins consumers to THIS version -- renewals do not follow. |
| `status.outputs.versionless_id` | `string` | The certificate's versionless data-plane ID: https://{vault}.vault.azure.net/certificates/{name}. Follows renewals automatically. |
| `status.outputs.secret_id` | `string` | The versioned ID of the certificate's SECRET face: https://{vault}.vault.azure.net/secrets/{name}/{version}. The value TLS terminators consume -- Application Gateway listeners and App Service certificate bindings take exactly this (or its versionless sibling below, to follow renewals). |
| `status.outputs.versionless_secret_id` | `string` | The versionless ID of the certificate's secret face: https://{vault}.vault.azure.net/secrets/{name}. The reference that keeps TLS terminators on the CURRENT certificate across renewals -- prefer it over secret_id wherever the consumer supports versionless references. |
| `status.outputs.certificate_name` | `string` | The certificate's name within the vault. |
| `status.outputs.version` | `string` | The current version identifier (the trailing segment of certificate_id). |
| `status.outputs.thumbprint` | `string` | The SHA-1 thumbprint of the current certificate, hex-encoded -- the fingerprint integrations match on. |
| `status.outputs.resource_manager_id` | `string` | The certificate's versioned ARM resource ID -- the control-plane identity used by ARM-level integrations. |
| `status.outputs.resource_manager_versionless_id` | `string` | The certificate's versionless ARM resource ID. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.keyVaultId` | AzureKeyVault | `status.outputs.key_vault_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureApplicationGateway | `spec.sslCertificates[].keyVaultSecretId` | `status.outputs.versionless_secret_id` |
| AzureApplicationGateway | `spec.trustedRootCertificates[].keyVaultSecretId` | `status.outputs.versionless_secret_id` |
| AzureContainerAppEnvironmentCertificate | `spec.certificateKeyVault.keyVaultSecretId` | `status.outputs.versionless_secret_id` |
| AzureFirewallPolicy | `spec.tlsCertificate.keyVaultSecretId` | `status.outputs.versionless_secret_id` |
| AzureFrontDoorSecret | `spec.keyVaultCertificateId` | `status.outputs.versionless_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
