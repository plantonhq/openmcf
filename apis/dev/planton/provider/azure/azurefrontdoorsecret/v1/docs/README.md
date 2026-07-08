# AzureFrontDoorSecret -- Design Research

## Scope

The secret is Front Door's bring-your-own certificate node: ARM
metadata wrapping a Key Vault certificate so custom domains can serve
it. It is a first-class kind because one certificate serves many
domains (wildcard/multi-SAN), it is FK-referenced by the custom
domain's `tls.secret_id`, and the rotation posture (follow-latest vs
pinned) is a decision made HERE, once, not per domain.

Source of truth: `azurerm_cdn_frontdoor_secret`
(terraform-provider-azurerm v4.80,
`internal/services/cdn/cdn_frontdoor_secret_resource.go`),
parity-verified against pulumi-azure v6 (`cdn.FrontdoorSecret`).

## Field mapping

| Spec field | azurerm attribute | Notes |
| --- | --- | --- |
| `profile_id` | `cdn_frontdoor_profile_id` | FK to AzureFrontDoorProfile; ForceNew |
| `secret_name` | `name` | ForceNew; 2-260, letter/digit edges, hyphens (the provider's regex) |
| `key_vault_certificate_id` | `secret.customer_certificate.key_vault_certificate_id` | FK to AzureKeyVaultCertificate; ForceNew. The provider accepts versioned AND versionless ids: versionless sets `UseLatestVersion` server-side (rotation follows), versioned pins that exact version |

## Shape decisions

- **The nested `secret { customer_certificate { ... } }` wrapper is
  flattened to one field**: the provider exposes exactly one secret
  parameter type (customer_certificate) with exactly one input, so the
  wrapper is TF schema scaffolding, not structure worth mirroring. The
  provider's internal mapping knows other parameter types (URL signing
  keys, first-party certificates) but does not surface them.
- **The FK defaults to `versionless_id`** (the rotation-follows-latest
  posture); the versioned `certificate_id` output remains referenceable
  for pin-one-version rollouts. Both are the vault DATA-PLANE
  identifiers -- exactly what the provider validates
  (`ValidateNestedItemID`, version optional).
- **Fully immutable** -- the provider has no Update (create/read/delete
  only); every field is ForceNew. Documented on the spec rather than
  worked around: rotation is a Key Vault operation, and domains survive
  a secret replacement because they reference it by ARM id.

## Validation contracts

- Name regex (the provider's own).
- Nothing else is statically checkable: the certificate id's shape is
  vault-domain-specific and the provider validates it at plan time.

## Apply-time contracts left to Azure (documented, not CELs)

- **The Front Door service principal grant**: Azure rejects the create
  with an access-denied error when the `Microsoft.AzureFrontDoor-Cdn`
  enterprise application lacks vault read access. A tenant-level
  operational prerequisite -- not expressible in a spec rule, recorded
  on the spec comment, the README, and the presets. The E2E harness
  bootstraps it for the test subscription.
- **CA-issued certificates only**: Azure validates the wrapped
  certificate's CONTENT at create time and rejects self-signed
  certificates with `BadRequest: "The certificate chain includes an
  invalid number of certificates. The number of certificates should be
  at least 2."` (live-verified). The chain requirement -- leaf plus
  issuer -- means the certificate must come from a CA (a Key Vault CA
  integration, or an imported PKCS#12 carrying its full chain). Not
  statically checkable: the spec holds a reference, and the content
  lives in the vault.

## Outputs

| Output | Source | Consumers |
| --- | --- | --- |
| `secret_id` | resource id | `AzureFrontDoorCustomDomain.tls.secret_id` |
| `secret_name` | `name` | human/portal orientation |
| `subject_alternative_names` | read back from the wrapped certificate | operators confirming hostname coverage before attaching |

## Composition

RG -> vault -> certificate -> profile -> secret -> custom domain ->
route. The secret is the seam where Key Vault's certificate lifecycle
meets Front Door's domain lifecycle.
