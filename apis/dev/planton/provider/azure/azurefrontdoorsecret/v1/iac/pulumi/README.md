# AzureFrontDoorSecret -- Pulumi Module

Creates a `cdn.FrontdoorSecret` on the referenced Front Door profile
(pulumi-azure classic v6), through the shared Azure provider builder
(static client secret, keyless web identity, or ambient chain),
wrapping the referenced Key Vault certificate.

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

## Build

```shell
make build
```
