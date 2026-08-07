# AzureFrontDoorCustomDomain

A custom domain inside an Azure Front Door profile: your own hostname
(e.g. www.example.com, or a wildcard) served by Front Door instead of
the generated *.azurefd.net endpoint hostname. TLS is always on --
either Azure's free managed certificate (issued, hosted, and
auto-rotated by Azure) or a bring-your-own certificate wrapped in an
AzureFrontDoorSecret.

A domain deploys immediately in a pending-validation state and exports
a `validation_token`; publishing that token as a DNS TXT record at
`_dnsauth.<host_name>` proves ownership and flips the domain to
approved. Routes then serve the domain through their
`custom_domain_ids` -- the route side owns the attachment.

## When to Use

Use AzureFrontDoorCustomDomain when you need:

- **Production hostnames on Front Door** -- every real site serves its
  own domain, not *.azurefd.net
- **Wildcard tenant domains** (`*.example.com`) with a shared BYO
  certificate -- the multi-tenant platform shape
- **Compliance TLS postures** -- pinned cipher suites via the
  `cipher_suite` block (predefined hardened sets or hand-picked)

## Key Configuration

- `profile_id` -- the parent profile; ForceNew
- `domain_name` -- the ARM resource name (hyphens, no dots); ForceNew
- `host_name` -- the actual hostname; ForceNew (it IS the identity).
  Managed certificates cap it at 64 characters and forbid wildcards
- `dns_zone_id` -- optional AzureDnsZone reference when DNS lives in
  Azure DNS
- `tls.certificate_type` -- MANAGED_CERTIFICATE (default) or
  CUSTOMER_CERTIFICATE paired with `tls.secret_id`
- `tls.cipher_suite` -- optional hardening: TLS12_2022 / TLS12_2023
  predefined sets, or CUSTOMIZED with explicit suites

The minimum TLS version is not configurable: Azure retired TLS 1.0/1.1,
so every domain floors at TLS 1.2 (TLS 1.3 always offered on top).

## Composition

```yaml
tls:
  certificateType: CUSTOMER_CERTIFICATE
  secretId:
    valueFrom:
      kind: AzureFrontDoorSecret
      name: my-wildcard-cert-secret
      fieldPath: status.outputs.secret_id
```

Routes serve the domain through its `custom_domain_id` output.

## Documentation

- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
