# AzureContainerAppEnvironmentCertificate - Pulumi Module

Pulumi implementation for the AzureContainerAppEnvironmentCertificate
component.

## Architecture

```
containerapp.EnvironmentCertificate (one BYO certificate on the environment)
```

## Key Design Decisions

- **Exactly one certificate source** -- inline PFX (blob + password) or
  Key Vault secret reference; the unused path's arguments are omitted so
  both engines send identical request bodies. An unset Key Vault
  identity defaults to `"System"`.
- **Certificate facts are exported for rotation monitoring** (subject,
  issuer, issue/expiration dates, thumbprint) because Azure never
  returns the PFX on read -- expiry is the only rotation signal.
- **Only tags update in place**; every other change replaces the
  certificate, briefly rebinding the custom domains that reference it.
- **PARITY-EXCEPTION on tag shape** versus the Terraform module
  (documented in both engines) -- output-neutral.

## Provider

Built via the shared `pulumiazureprovider.Get` builder -- static client
secret, keyless web identity, or ambient chain, resolved from the stack
input. Never construct the provider inline.
