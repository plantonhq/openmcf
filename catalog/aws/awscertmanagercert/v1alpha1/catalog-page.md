# AWS Certificate

Deploys an AWS Certificate Manager (ACM) certificate in any of ACM's three creation modes: **requested** (Amazon-issued, DNS or EMAIL validated), **imported** (bring-your-own certificate material), or **private** (issued by an ACM Private CA). For DNS-validated certificates whose zone is in Route53, Planton also creates the validation records and waits for issuance — fully hands-off HTTPS.

## What Gets Created

Depending on the creation mode, Planton provisions:

- **ACM Certificate** — the certificate itself: requested from ACM, imported from your PEM material, or issued by your private CA. A certificate swap is create-before-destroy so consumers holding the ARN are never left dangling.
- **Route53 CNAME Records** (managed DNS validation only) — one validation record per certificate domain in the referenced hosted zone, created as UPSERTs so re-runs and shared records never collide.
- **Certificate Validation Waiter** (managed DNS validation only) — blocks the deployment until ACM reports `ISSUED`, so downstream resources attach a ready certificate. Disable with `waitForValidation: false`.

## Prerequisites

- **AWS credentials** configured via a Planton provider config
- For managed DNS validation: a **Route53 public hosted zone** authoritative for every domain on the certificate
- For the private mode: an **ACM Private CA** and its ARN
- For the imported mode: PEM-encoded certificate, private key, and (if needed) chain

## Quick Start

Create a file `cert.yaml`:

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsCertManagerCert
metadata:
  name: my-cert
spec:
  region: us-east-1
  primaryDomainName: example.com
  route53HostedZoneId:
    value: Z0123456789ABCDEFGHIJ
```

Deploy:

```shell
planton apply -f cert.yaml
```

This requests a certificate for `example.com`, creates the DNS validation CNAME in the zone, and waits until ACM issues the certificate.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | AWS region for the certificate. CloudFront only accepts certificates from `us-east-1`. | Required; non-empty |

### Creation Mode (exactly one)

| Field | Type | Description |
|-------|------|-------------|
| `primaryDomainName` | `string` | Requested or private mode: the main domain (apex, subdomain, or `*.wildcard`). |
| `imported` | `object` | Imported mode: `certificateBody` + `privateKey` (sensitive) + optional `certificateChain`, all PEM. |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `alternateDomainNames` | `string[]` | `[]` | Subject Alternative Names; each validated independently. Do not repeat the primary domain. |
| `validationMethod` | `string` | `DNS` | `DNS` (recommended — renewals stay automatic) or `EMAIL` (re-approval every renewal). |
| `validationOptions` | `object[]` | `[]` | Per-domain overrides of where the validation request is sent (`domainName` + `validationDomain`). |
| `keyAlgorithm` | `string` | `RSA_2048` | `RSA_2048`, `RSA_3072`, `RSA_4096`, `EC_prime256v1`, `EC_secp384r1`, `EC_secp521r1`. Create-time immutable. |
| `route53HostedZoneId` | `StringValueOrRef` | — | Route53 zone for automated DNS validation. Unset = external DNS: the records to create are exported as outputs. |
| `waitForValidation` | `bool` | `true` | Wait for `ISSUED` after creating managed validation records. |
| `options.certificateTransparencyLoggingPreference` | `string` | `ENABLED` | Public CT logging (`ENABLED`/`DISABLED`). |
| `options.export` | `string` | `DISABLED` | `ENABLED` makes the certificate exportable off-AWS (additional AWS charge). |
| `certificateAuthorityArn` | `string` | — | ACM-PCA authority ARN — selects the private mode (no public validation). |

## Examples

### Wildcard Certificate with Automated Validation

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsCertManagerCert
metadata:
  name: wildcard-cert
spec:
  region: us-east-1
  primaryDomainName: "*.example.com"
  alternateDomainNames:
    - example.com
  route53HostedZoneId:
    valueFrom:
      kind: AwsRoute53Zone
      name: my-zone
      fieldPath: status.outputs.zone_id
```

### External DNS (zone outside Route53)

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsCertManagerCert
metadata:
  name: external-dns-cert
spec:
  region: us-west-2
  primaryDomainName: app.example.com
```

The deployment finishes with the certificate in `PENDING_VALIDATION`; create the records from the `domain_validation_records` output in your DNS provider and ACM issues automatically.

### Imported Certificate

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsCertManagerCert
metadata:
  name: imported-cert
spec:
  region: eu-west-1
  imported:
    certificateBody: |
      -----BEGIN CERTIFICATE-----
      ...
      -----END CERTIFICATE-----
    privateKey: |
      -----BEGIN PRIVATE KEY-----
      ...
      -----END PRIVATE KEY-----
```

### Private (ACM-PCA) Certificate

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsCertManagerCert
metadata:
  name: internal-cert
spec:
  region: us-west-2
  primaryDomainName: svc.internal.example.com
  certificateAuthorityArn: arn:aws:acm-pca:us-west-2:123456789012:certificate-authority/11111111-2222-3333-4444-555555555555
```

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `cert_arn` | `string` | The certificate ARN — the join key TLS consumers reference (listeners, CloudFront, Cognito, OpenSearch, Client VPN). |
| `status` | `string` | `PENDING_VALIDATION` until ownership is proven, then `ISSUED`. |
| `domain_validation_records` | `list` | The DNS records proving ownership (`domain_name`, `record_name`, `record_type`, `record_value`) — create these when DNS is external; keep them in place for automatic renewal. |
| `not_before` / `not_after` | `string` | The validity window (RFC3339); `not_after` is the re-import deadline for imported certificates. |
| `certificate_type` | `string` | `AMAZON_ISSUED`, `IMPORTED`, or `PRIVATE`. |

## Related Components

- [AwsRoute53Zone](/docs/catalog/aws/awsroute53zone) — provides the hosted zone for automated DNS validation
- [AwsLbListener](/docs/catalog/aws/awslblistener) — terminates TLS on load balancers with this certificate
- [AwsCloudFront](/docs/catalog/aws/awscloudfront) — serves HTTPS on custom domains with a `us-east-1` certificate
- [AwsCognitoUserPool](/docs/catalog/aws/awscognitouserpool) — custom auth domains reference the certificate ARN
- [AwsOpenSearchDomain](/docs/catalog/aws/awsopensearchdomain) — custom endpoints reference the certificate ARN
