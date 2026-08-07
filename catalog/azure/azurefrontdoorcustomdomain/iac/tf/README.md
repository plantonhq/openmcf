# AzureFrontDoorCustomDomain -- Terraform Module

Creates an `azurerm_cdn_frontdoor_custom_domain` on the referenced
Front Door profile. The domain lands in the pending-validation state;
the `validation_token` output is the DNS TXT challenge.

## Inputs

- `metadata` -- Planton resource metadata (name, org, env, labels)
- `spec` -- see `variables.tf`; mirrors `spec.proto` exactly. Enum
  fields arrive as the spec enum's FULL value names and are mapped in
  `locals.tf`.

## Behavior notes

- **`minimum_version` is the constant `TLS12`** -- Azure retired TLS
  1.0/1.1, the provider accepts exactly one value, so the module sends
  it unconditionally (never the deprecated `minimum_tls_version`
  alias).
- **`certificate_type` materializes `ManagedCertificate`** when
  unspecified (tfvars drops unset fields; the documented default).
- **The cipher-suite block is sent only when configured**; an empty
  pinned tls13 list normalizes to null so Azure's TLS 1.3 defaults
  apply.
- **No Azure tags** -- ARM does not support tags on custom domains.

## Outputs

- `custom_domain_id`, `host_name`, `validation_token`, `expiration_date`
