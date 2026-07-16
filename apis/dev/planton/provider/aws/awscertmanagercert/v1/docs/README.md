# AWS Certificate Manager Certificate

AWS Certificate Manager (ACM) provisions, stores, and renews SSL/TLS certificates for AWS services. This document describes how Planton models ACM certificates: the three creation modes, how DNS validation is automated, and how the certificate composes with other resources.

## Certificate identity

ACM certificates have no name in AWS — identity is the generated ARN. `metadata.name` drives the `Name` identity tag, and consumers join to the certificate through the `cert_arn` stack output (load-balancer listeners, CloudFront, Cognito custom domains, OpenSearch custom endpoints, Client VPN).

## Creation modes

Provide **exactly one** of:

| Mode | Fields | When to use |
|------|--------|-------------|
| **Requested** (default) | `primary_domain_name` (+ `alternate_domain_names`, `validation_method`) | Amazon-issued public certificates — free, auto-renewing, what almost every TLS-fronting AWS service consumes |
| **Imported** | `imported.certificate_body`, `imported.private_key`, optional `imported.certificate_chain` | Certificates issued by an external CA that ACM should distribute to integrated services |
| **Private (ACM-PCA)** | `primary_domain_name` + `certificate_authority_arn` | Internal TLS issued by your AWS Private Certificate Authority — no public validation |

Validation rules keep the arms honest: a manifest cannot mix imported material with issuance-time fields, and private-CA issuance rejects public-validation fields.

## Domain strategy

- A wildcard (`*.example.com`) covers one label level and does **not** cover the apex. The standard production pairing is the apex plus its wildcard — `primary_domain_name: example.com` with `alternate_domain_names: ["*.example.com"]` — so one certificate serves the bare domain and every first-level subdomain.
- Each domain is validated independently; a domain and its wildcard share the same validation record.

## DNS validation: managed or external

DNS validation is the recommended method — the validation CNAMEs stay in place, so ACM renews the certificate automatically forever.

- **Managed (Route53)** — set `route53_hosted_zone_id` (a literal zone ID or a reference to an `AwsRoute53Zone`). The module creates the validation CNAMEs in the zone and, unless `wait_for_validation` is `false`, blocks until the certificate reaches `ISSUED`.
- **External DNS** — leave `route53_hosted_zone_id` unset. The certificate is created in `PENDING_VALIDATION`, and the exact records to create are exported as the `domain_validation_records` stack output. Once the records exist, ACM issues — and keeps renewing — automatically.

EMAIL validation exists for domains where DNS control is unavailable, but every renewal requires manual re-approval; prefer DNS.

## The CloudFront region constraint

ACM certificates are regional, and every regional consumer (ALB, API Gateway, OpenSearch) needs the certificate in its own region. CloudFront is the exception: it only accepts certificates from **us-east-1**, regardless of where the distribution's origins live. Architectures that pair regional load balancers with CloudFront therefore carry one certificate per consumer region plus one in us-east-1.

## Renewal

- **Amazon-issued** certificates renew automatically as long as the validation DNS records stay in place — never delete them.
- **Imported** certificates never auto-renew. Re-importing new material before expiry updates in place and keeps the same ARN, so consumers are undisturbed. The `not_after` output is the re-import deadline.

## Options

- `key_algorithm` — `RSA_2048` (default), `RSA_3072`, `RSA_4096`, `EC_prime256v1`, `EC_secp384r1`, `EC_secp521r1`. Create-time immutable.
- `options.certificate_transparency_logging_preference` — CT logging (`ENABLED` by default; browsers increasingly require it).
- `options.export` — exportable certificates (`ENABLED` allows exporting the private key for use off-AWS, at an additional AWS charge).
- `validation_options` — per-domain overrides of where the validation request is sent (e.g. EMAIL-validating a subdomain at its parent).

## Deletion behavior

ACM refuses to delete a certificate that is still associated with any AWS resource (`ResourceInUseException`). Some associations are indirect — API Gateway custom domains and Cognito custom domains hold references through AWS-managed infrastructure that can take more than 15 minutes to release after the consumer is deleted. Destroy consumers first and expect the certificate deletion to retry through that propagation window.

## Stack outputs

| Output | Description |
|--------|-------------|
| `cert_arn` | The join key every TLS consumer references |
| `status` | `PENDING_VALIDATION` until ownership is proven, then `ISSUED` |
| `domain_validation_records` | The DNS records to create when validation is managed externally |
| `not_before` / `not_after` | The validity window (re-import deadline for imported certificates) |
| `certificate_type` | `AMAZON_ISSUED`, `IMPORTED`, or `PRIVATE` |
