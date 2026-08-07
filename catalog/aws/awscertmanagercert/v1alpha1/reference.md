# AwsCertManagerCert

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1alpha1`

AwsCertManagerCertSpec defines an AWS Certificate Manager (ACM)
certificate in any of ACM's three creation modes:

  1. REQUESTED (Amazon-issued) -- ACM issues a public certificate for
     domains you prove ownership of, via DNS or EMAIL validation. This
     is the default mode and what almost every TLS-fronting AWS
     service (ALB/NLB listeners, CloudFront, API Gateway, Cognito,
     OpenSearch) consumes. Set primary_domain_name.
  2. IMPORTED -- bring your own certificate material (issued by an
     external CA) and let ACM distribute it to integrated services.
     Set the imported block.
  3. PRIVATE (ACM-PCA) -- ACM issues a certificate from your AWS
     Private Certificate Authority, with no public validation. Set
     primary_domain_name together with certificate_authority_arn.

Exactly one of primary_domain_name / imported drives the mode; CEL
rules keep the arms honest so a manifest cannot mix them.

ACM certificates have no name in AWS -- identity is the generated
ARN. metadata.name drives the Name identity tag, and consumers
compose with this certificate by referencing its cert_arn output.

Renewal: Amazon-issued certificates renew automatically as long as
the validation DNS records stay in place (one more reason DNS
validation is the recommended method). Imported certificates never
auto-renew -- re-import new material before expiry.

## Example

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsCertManagerCert
metadata:
  name: awscertmanagercert-demo
spec:
  region: us-east-1
  primaryDomainName: example.com
  alternateDomainNames:
    - www.example.com
  validationMethod: DNS
  # No route53HostedZoneId: the certificate rests in PENDING_VALIDATION and
  # the validation records are exported for external DNS -- the leanest
  # deployable shape (no zone prerequisite, nothing to wait on).
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.primaryDomainName` | `string` |  |  |  |
| `spec.alternateDomainNames` | `[]string` |  |  |  |
| `spec.validationMethod` | `string` |  |  |  |
| `spec.validationOptions` | `[]AwsCertManagerCertValidationOption` |  |  |  |
| `spec.validationOptions[].domainName` | `string` | yes |  |  |
| `spec.validationOptions[].validationDomain` | `string` | yes |  |  |
| `spec.keyAlgorithm` | `string` |  |  |  |
| `spec.route53HostedZoneId` | `string \| valueFrom` |  |  | AwsRoute53Zone (`status.outputs.zone_id`) |
| `spec.waitForValidation` | `bool` |  | `true` |  |
| `spec.options` | `AwsCertManagerCertOptions` |  |  |  |
| `spec.options.certificateTransparencyLoggingPreference` | `string` |  |  |  |
| `spec.options.export` | `string` |  |  |  |
| `spec.imported` | `AwsCertManagerCertImported` |  |  |  |
| `spec.imported.certificateBody` | `string` | yes |  |  |
| `spec.imported.privateKey` | `string` (sensitive) | yes |  |  |
| `spec.imported.certificateChain` | `string` |  |  |  |
| `spec.certificateAuthorityArn` | `string` |  |  |  |

## Field Details

### spec.region

`string` · required

The AWS region the certificate is created in. Regional services
(ALB, API Gateway, OpenSearch, ...) require the certificate in
their own region. CloudFront is the notable exception: it only
accepts certificates from "us-east-1", regardless of where the
distribution's origins live.
Example: "us-west-2", "eu-west-1".

- rule: {"string":{"minLen":"1"}}

### spec.primaryDomainName

`string`

The main domain the certificate covers -- an apex ("example.com"),
a subdomain ("app.example.com"), or a wildcard ("*.example.com",
which covers one label level only). Setting this selects the
requested (Amazon-issued) mode, or the private mode when
certificate_authority_arn is also set. Leave empty only for
imported certificates (ACM derives their domains from the
certificate body).

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"pattern":"^(?:\\*\\.[A-Za-z0-9\\-\\.]+|[A-Za-z0-9\\-\\.]+\\.[A-Za-z]{2,})$"}}

### spec.alternateDomainNames

`[]string`

Additional domains (Subject Alternative Names) the certificate
covers, each validated independently. A common pairing is the apex
plus its wildcard ("example.com" + "*.example.com") so one
certificate serves both the bare domain and every subdomain. Do
not repeat primary_domain_name here -- ACM already includes it.

- rule: {"repeated":{"unique":true,"items":{"string":{"pattern":"^(?:\\*\\.[A-Za-z0-9\\-\\.]+|[A-Za-z0-9\\-\\.]+\\.[A-Za-z]{2,})$"}}}}

### spec.validationMethod

`string`

How ACM verifies domain ownership for requested certificates:
"DNS" (recommended -- a CNAME record per domain, kept in place so
renewals are fully automatic) or "EMAIL" (approval mail to the
domain's WHOIS/admin addresses; renewal requires re-approval every
time, so prefer DNS). Empty keeps DNS. Not applicable to imported
or private certificates.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["DNS","EMAIL"]}}

### spec.validationOptions

`[]AwsCertManagerCertValidationOption`

Per-domain overrides for where the validation request is sent.
The classic use is EMAIL-validating a subdomain at its parent:
validating "app.example.com" by mailing the owners of
"example.com". Rarely needed with DNS validation.

### spec.validationOptions[].domainName

`string` · required

The certificate domain this override applies to (the primary
domain or one of the alternate domains).

- rule: {"string":{"minLen":"1"}}

### spec.validationOptions[].validationDomain

`string` · required

The domain the validation request is sent to -- must be the
domain itself or one of its ancestors. Example: validate
"app.example.com" via "example.com".

- rule: {"string":{"minLen":"1"}}

### spec.keyAlgorithm

`string`

The key algorithm for the certificate's key pair, create-time
immutable: "RSA_2048" (the default -- universally compatible),
"RSA_3072", "RSA_4096", "EC_prime256v1" (smaller/faster TLS
handshakes; check client compatibility), "EC_secp384r1", or
"EC_secp521r1". Empty keeps the ACM default (RSA_2048). Not
applicable to imported certificates (the algorithm is baked into
the imported key material).

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["RSA_2048","RSA_3072","RSA_4096","EC_prime256v1","EC_secp384r1","EC_secp521r1"]}}

### spec.route53HostedZoneId

`string | valueFrom`

The Route53 public hosted zone where DNS validation records are
created automatically. When set (DNS validation only), the module
creates the validation CNAMEs in this zone and -- unless
wait_for_validation is false -- waits for the certificate to be
ISSUED before finishing. When unset, the certificate is created in
PENDING_VALIDATION and the required records are exported as the
domain_validation_records output for you to create in your
external DNS; the deployment does not wait. The zone must be
authoritative for every domain on the certificate.

- references: AwsRoute53Zone (`status.outputs.zone_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsRoute53Zone, name: <that resource's name>, fieldPath: status.outputs.zone_id}} -- a bare string does not parse

### spec.waitForValidation

`bool` · optional (explicit presence)

Whether the deployment waits for the certificate to reach ISSUED
after creating the validation records (DNS validation with
route53_hosted_zone_id only -- issuance typically lands within a
few minutes). Set false to create the certificate and records
without blocking; downstream resources that require an ISSUED
certificate (CloudFront, listeners) will then fail until issuance
completes, so keep the default unless you are staging DNS ahead of
time. Ignored when the module does not manage the validation
records.

- default: `true`

### spec.options

`AwsCertManagerCertOptions`

Certificate options for Amazon-issued (requested and private)
certificates.

### spec.options.certificateTransparencyLoggingPreference

`string`

Whether the certificate is recorded in public Certificate
Transparency logs: "ENABLED" (the ACM default -- browsers
increasingly require CT-logged certificates, so keep it unless
you must hide internal hostnames) or "DISABLED". Empty keeps
ENABLED.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["ENABLED","DISABLED"]}}

### spec.options.export

`string`

Whether the certificate's private key may be exported ("ENABLED"
or "DISABLED", the ACM default). Exportable public certificates
let you run the same certificate on non-AWS infrastructure, and
incur an additional AWS charge per certificate. Empty keeps
DISABLED.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["ENABLED","DISABLED"]}}

### spec.imported

`AwsCertManagerCertImported`

Bring-your-own certificate material issued by an external CA.
Setting this block selects the imported mode: ACM stores and
distributes the certificate but never renews it -- re-import new
material before expiry (updates re-import in place, keeping the
same ARN so consumers are undisturbed).

### spec.imported.certificateBody

`string` · required

The PEM-encoded certificate. Public material, not a secret.

- rule: {"required":true}

### spec.imported.privateKey

`string` · required · sensitive

The PEM-encoded, unencrypted private key matching the certificate.

- rule: {"required":true}

### spec.imported.certificateChain

`string`

The PEM-encoded intermediate/root chain, if the issuing CA is not
already trusted by AWS. Optional for certificates issued by
well-known public CAs. Public material, not a secret.

### spec.certificateAuthorityArn

`string`

The ARN of an AWS Private Certificate Authority (ACM-PCA).
Setting this together with primary_domain_name selects the
private mode: the CA issues the certificate directly, no public
validation happens, and validation_method must stay unset.
Private certificates are for internal TLS (service meshes,
internal ALBs) where clients trust your private root.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"pattern":"^arn:aws[a-zA-Z-]*:acm-pca:[a-z0-9-]+:[0-9]{12}:certificate-authority/.+$"}}

## Validation Rules

- `exactly_one_creation_mode`: set exactly one of primary_domain_name (requested or private certificate) or imported (bring-your-own certificate material)
- `imported_excludes_issuance_fields`: imported certificates derive their domains, algorithm, and options from the certificate material -- remove alternate_domain_names, validation_method, validation_options, key_algorithm, options, certificate_authority_arn, and route53_hosted_zone_id
- `private_ca_excludes_validation`: private (ACM-PCA) certificates are issued without public validation -- remove validation_method, validation_options, and route53_hosted_zone_id when certificate_authority_arn is set
- `route53_zone_requires_dns_validation`: route53_hosted_zone_id automates DNS validation records -- it cannot be combined with validation_method EMAIL

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsCertManagerCert, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.cert_arn` | `string` | The certificate ARN -- the join key every TLS consumer references (load-balancer listeners, CloudFront, Cognito custom domains, OpenSearch custom endpoints, Client VPN, ...). |
| `status.outputs.status` | `string` | The certificate status at the end of the deployment: "PENDING_VALIDATION" (created, waiting for domain validation -- the resting state when validation records are managed externally), "ISSUED" (ready to serve traffic), or a failure state ("VALIDATION_TIMED_OUT", "FAILED", "REVOKED", "EXPIRED"). |
| `status.outputs.domain_validation_records` | `[]AwsCertManagerCertDomainValidationRecord` | The DNS records that prove domain ownership, one per domain -- create these in your DNS when route53_hosted_zone_id is not set, and keep them in place afterwards so ACM can renew automatically. Empty for imported and private certificates (no public validation). |
| `status.outputs.domain_validation_records[].domain_name` | `string` | The certificate domain this record validates. |
| `status.outputs.domain_validation_records[].record_name` | `string` | The DNS record name to create (a "_<hash>.<domain>." CNAME). |
| `status.outputs.domain_validation_records[].record_type` | `string` | The DNS record type (always "CNAME" today). |
| `status.outputs.domain_validation_records[].record_value` | `string` | The DNS record value ACM expects the name to resolve to. |
| `status.outputs.not_before` | `string` | Start of the certificate's validity window (RFC3339). Empty until the certificate is issued. |
| `status.outputs.not_after` | `string` | End of the certificate's validity window (RFC3339) -- when an imported certificate must be re-imported by. Empty until the certificate is issued. |
| `status.outputs.certificate_type` | `string` | How the certificate came to be: "AMAZON_ISSUED" (requested), "IMPORTED", or "PRIVATE" (ACM-PCA). |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.route53HostedZoneId` | AwsRoute53Zone | `status.outputs.zone_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AwsClientVpn | `spec.authenticationOptions[].rootCertificateChainArn` | `status.outputs.cert_arn` |
| AwsClientVpn | `spec.serverCertificateArn` | `status.outputs.cert_arn` |
| AwsCloudFront | `spec.viewerCertificate.acmCertificateArn` | `status.outputs.cert_arn` |
| AwsCognitoUserPool | `spec.domain.certificateArn` | `status.outputs.cert_arn` |
| AwsHttpApiDomain | `spec.certificateArn` | `status.outputs.cert_arn` |
| AwsLbListener | `spec.certificateArn` | `status.outputs.cert_arn` |
| AwsLbListener | `spec.additionalCertificateArns` | `status.outputs.cert_arn` |
| AwsOpenSearchDomain | `spec.domainEndpointOptions.customEndpointCertificateArn` | `status.outputs.cert_arn` |

## See Also

- [Overview](../README.md)
