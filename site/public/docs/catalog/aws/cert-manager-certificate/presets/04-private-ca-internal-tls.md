---
title: "Private CA Certificate for Internal TLS"
description: "This preset issues a certificate from your AWS Private Certificate Authority (ACM-PCA) with managed early renewal -- the shape for internal TLS where clients trust your private root instead of a..."
type: "preset"
rank: "04"
presetSlug: "04-private-ca-internal-tls"
componentSlug: "cert-manager-certificate"
componentTitle: "Cert Manager Certificate"
provider: "aws"
icon: "package"
order: 4
---

# Private CA Certificate for Internal TLS

This preset issues a certificate from your AWS Private Certificate Authority (ACM-PCA) with managed early renewal -- the shape for internal TLS where clients trust your private root instead of a public CA.

## When to Use

- Internal service-to-service TLS: service meshes, internal ALBs, east-west APIs on private domains
- Any TLS surface on a domain that never faces the public internet (no public validation is possible or wanted)

## Key Configuration Choices

- **Private issuance** -- setting `certificateAuthorityArn` alongside the domain selects the private mode: the CA issues directly, no DNS/EMAIL validation happens, and validation fields must stay unset (validation rules enforce this)
- **Managed early renewal** (`earlyRenewalDuration: P90D`) -- ACM renews the certificate 90 days before expiry while keeping the same ARN, so consumers are undisturbed; this mechanism is private-CA-only (publicly validated certificates renew on ACM's own schedule, imported ones never renew), and durations under 60 days have no effect
- **Consumers compose by ARN** -- reference the `cert_arn` output from listeners and other TLS-fronting resources exactly as with public certificates

## Values to Replace

The preset ships realistic example values (they document the expected shapes); replace them with your own:

| Field | Description | Where to Find |
| --- | --- | --- |
| `region` | The region the certificate is created in (match the consuming service) | Your deployment plan |
| `primaryDomainName` | The internal domain the certificate covers | Your internal DNS plan |
| `certificateAuthorityArn` | The Private CA that issues the certificate | Your ACM-PCA authority's ARN |

## Related Presets

- **01-single-domain-dns** -- Use for publicly validated certificates on internet-facing domains
- **03-external-dns** -- Use when the domain's DNS lives outside Route53
