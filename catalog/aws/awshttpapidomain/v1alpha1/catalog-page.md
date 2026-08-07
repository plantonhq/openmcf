# AWS HTTP API Domain

Deploys an API Gateway v2 custom domain name with its API mappings -- binding an owned domain (e.g. `api.example.com`) to an ACM certificate and publishing one or more HTTP APIs under it, optionally namespaced by path keys.

## What Gets Created

When you deploy an AwsHttpApiDomain resource, Planton provisions:

- **Custom domain name** -- an API Gateway v2 domain with your ACM certificate, REGIONAL endpoint, TLS 1.2 policy, and optional mutual TLS
- **API mappings** -- one mapping per entry, binding an API's stage under an optional path key (empty key = the domain root)

DNS is composed downstream: the exported `target_domain_name` / `hosted_zone_id` feed a Route 53 alias record (`AwsRoute53DnsRecord`).

## Prerequisites

- **An ACM certificate** in the same region covering the domain name (exact or wildcard match) -- compose with `AwsCertManagerCert`
- **An HTTP API** to map -- compose with `AwsHttpApiGateway`
- (Optional) **A Route 53 zone** to publish the alias record

## Quick Start

Create a file `domain.yaml`:

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsHttpApiDomain
metadata:
  name: api-example-com
spec:
  region: us-east-1
  domainName: api.example.com
  certificateArn:
    valueFrom:
      kind: AwsCertManagerCert
      name: api-example-cert
      fieldPath: status.outputs.cert_arn
  apiMappings:
    - apiId:
        valueFrom:
          kind: AwsHttpApiGateway
          name: orders-api
          fieldPath: status.outputs.api_id
      stage: $default
```

Deploy:

```shell
planton apply -f domain.yaml
```

Then create a Route 53 alias record pointing `api.example.com` at the exported `target_domain_name` / `hosted_zone_id`.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | AWS region. Must match the certificate's and mapped APIs' region. | Non-empty |
| `domainName` | `string` | The fully qualified custom domain. Immutable -- changing it replaces the domain. Wildcards (`*.example.com`) allowed. | Lowercase; 1-512 chars |
| `certificateArn` | `StringValueOrRef` | ACM certificate covering the domain, issued/imported in the same region. Can reference `AwsCertManagerCert` via `valueFrom`. | Required |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `ipAddressType` | `string` | AWS default (dualstack) | `"ipv4"` or `"dualstack"` |
| `mutualTls.truststoreUri` | `string` | — | S3 URI of the PEM CA bundle clients must chain to (e.g. `s3://bucket/truststore.pem`). Bucket must be in the same region. |
| `mutualTls.truststoreVersion` | `string` | — | S3 object version pin -- makes CA rotation an explicit, auditable change |
| `apiMappings[].apiId` | `StringValueOrRef` | — | The API to map. Can reference `AwsHttpApiGateway` via `valueFrom`. |
| `apiMappings[].stage` | `string` | — | The API stage to serve (typically `"$default"`) |
| `apiMappings[].apiMappingKey` | `string` | `""` (root) | Path key under which the API is served (no slashes). Keys must be unique across mappings; only one API can serve the root. |

## Examples

### Multi-API Domain with Path Keys

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsHttpApiDomain
metadata:
  name: api-example-com
spec:
  region: us-east-1
  domainName: api.example.com
  certificateArn:
    value: arn:aws:acm:us-east-1:123456789012:certificate/abc-123
  apiMappings:
    - apiId:
        value: a1b2c3d4    # https://api.example.com/
      stage: $default
    - apiId:
        value: e5f6g7h8    # https://api.example.com/orders/...
      stage: $default
      apiMappingKey: orders
```

### Mutual TLS for Machine-to-Machine APIs

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsHttpApiDomain
metadata:
  name: partner-api-domain
spec:
  region: us-east-1
  domainName: partner-api.example.com
  certificateArn:
    value: arn:aws:acm:us-east-1:123456789012:certificate/abc-123
  mutualTls:
    truststoreUri: s3://example-security/truststore.pem
    truststoreVersion: 3JpbWFyeSB0cnVzdHN0b3Jl
  apiMappings:
    - apiId:
        value: a1b2c3d4
      stage: $default
```

When mTLS is enabled, also set `disableExecuteApiEndpoint: true` on every mapped API -- otherwise callers can bypass mTLS via the default execute-api endpoint.

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `domain_name` | `string` | The custom domain name (the domain's join key) |
| `domain_name_arn` | `string` | ARN of the domain name resource |
| `target_domain_name` | `string` | The API Gateway regional domain to alias/CNAME to |
| `hosted_zone_id` | `string` | Route 53 hosted zone ID of the regional endpoint (the alias target zone) |

## Related Components

- [AwsHttpApiGateway](/docs/catalog/aws/awshttpapigateway) — The HTTP APIs mapped onto this domain
- [AwsCertManagerCert](/docs/catalog/aws/awscertmanagercert) — The TLS certificate the domain terminates with
- [AwsRoute53DnsRecord](/docs/catalog/aws/awsroute53dnsrecord) — The alias record publishing the domain in DNS
- [AwsRoute53Zone](/docs/catalog/aws/awsroute53zone) — The hosted zone carrying the alias record
