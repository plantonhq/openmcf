# AzureFrontDoorCustomDomain -- Pulumi Module

Creates a `cdn.FrontdoorCustomDomain` on the referenced Front Door
profile (pulumi-azure classic v6), through the shared Azure provider
builder (static client secret, keyless web identity, or ambient chain).
The domain lands in the pending-validation state; the
`validation_token` output is the DNS TXT challenge.

## Behavior notes

- **`MinimumVersion` is the constant `TLS12`** -- Azure retired TLS
  1.0/1.1, the provider accepts exactly one value, so the module sends
  it unconditionally (never the deprecated `MinimumTlsVersion` alias).
- **`CertificateType` is sent only when chosen** -- ARM defaults to
  ManagedCertificate.
- **The cipher-suite block is sent only when configured**; a pinned
  tls13 list is sent only when non-empty so Azure's TLS 1.3 defaults
  apply otherwise.
- **No Azure tags** -- ARM does not support tags on custom domains.

## Outputs

- `custom_domain_id`, `host_name`, `validation_token`, `expiration_date`

## Build

```shell
make build
```
