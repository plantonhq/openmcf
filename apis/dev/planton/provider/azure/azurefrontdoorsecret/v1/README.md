# AzureFrontDoorSecret

A secret inside an Azure Front Door profile: the bring-your-own TLS
certificate node. It wraps a Key Vault certificate so custom domains
can terminate TLS with it -- the domain's `tls.secret_id` references
this secret, and this secret references the AzureKeyVaultCertificate
that holds the key material.

The secret is first-class because one certificate -- typically a
wildcard or multi-SAN cert -- serves many custom domains, and rotating
it must be one operation. With a VERSIONLESS certificate reference (the
default), rotation is not even that: Front Door follows the
certificate's latest Key Vault version automatically.

## When to Use

Use AzureFrontDoorSecret when you need:

- **Wildcard custom domains** -- Azure's managed certificates cover
  exact names only, so `*.example.com` requires BYO
- **EV/OV or org-CA certificates** on any Front Door hostname
- **One certificate across many domains** -- every tenant subdomain
  references the same secret

## Key Configuration

- `profile_id` -- the parent profile; ForceNew
- `secret_name` -- 2-260 characters, unique within the profile;
  ForceNew
- `key_vault_certificate_id` -- the wrapped certificate, defaulting to
  the AzureKeyVaultCertificate's `versionless_id` output
  (rotation-follows-latest); reference the versioned `certificate_id`
  to pin one exact version instead

The whole resource is immutable (Azure exposes no update) -- safe,
because rotation happens in Key Vault, not here.

## One-Time Prerequisite (per tenant)

Front Door reads Key Vault with Microsoft's own service principal (the
`Microsoft.AzureFrontDoor-Cdn` enterprise application). Grant it read
access on the vault -- e.g. "Key Vault Secrets User" on an RBAC-mode
vault -- before the first secret deploys.

## Certificate Content Requirement

Azure rejects SELF-SIGNED certificates for Front Door BYO TLS: the
wrapped certificate must be CA-issued with a complete chain (leaf plus
issuer -- at least two certificates). Enroll the Key Vault certificate
through a CA integration, or import a PKCS#12 that carries its full
chain.

## Composition

```yaml
keyVaultCertificateId:
  valueFrom:
    kind: AzureKeyVaultCertificate
    name: my-wildcard-cert
    fieldPath: status.outputs.versionless_id
```

Custom domains terminate TLS through the secret's `secret_id` output.

## Documentation

- [Design research](docs/README.md) -- field mapping, recorded skips
- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
