# AzureFrontDoorSecret -- Terraform Module

Creates an `azurerm_cdn_frontdoor_secret` on the referenced Front Door
profile, wrapping the referenced Key Vault certificate.

## Inputs

- `metadata` -- Planton resource metadata (name, org, env, labels)
- `spec` -- see `variables.tf`; mirrors `spec.proto` exactly

## Behavior notes

- **Fully immutable** -- the provider exposes no update; any change
  replaces the secret. Rotation happens in Key Vault: a versionless
  `key_vault_certificate_id` follows the latest certificate version
  automatically.
- **Operational prerequisite** -- Front Door's service principal
  (`Microsoft.AzureFrontDoor-Cdn`) must hold vault read access before
  the create succeeds; Azure otherwise rejects it with an
  access-denied error naming the vault.
- **No Azure tags** -- ARM does not support tags on Front Door secrets.

## Outputs

- `secret_id` -- what a custom domain's `tls.secret_id` references
- `secret_name`
- `subject_alternative_names` -- the DNS names the certificate covers
