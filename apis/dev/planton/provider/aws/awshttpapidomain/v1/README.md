# AwsHttpApiDomain

Deploy and manage an API Gateway v2 custom domain name using Planton -- the production front door that binds an owned domain (e.g. `api.example.com`) to an ACM certificate and maps one or more HTTP APIs onto it.

## Overview

A custom domain replaces the default `https://{api-id}.execute-api.{region}.amazonaws.com` endpoint with your own DNS name and TLS certificate. API mappings then publish APIs under the domain, optionally namespaced by a path key -- one API at the root, another under `/orders`, a third under `/billing`.

The domain is deliberately its own resource rather than a field on the API:

- **A domain outlives any one API** -- APIs come and go; the domain, its certificate binding, and its DNS records persist.
- **Many APIs, one domain** -- mappings compose multiple APIs under distinct path keys.
- **Independent certificate lifecycle** -- the ACM certificate rotates without touching any API.

DNS is composed, not embedded: the domain exports `target_domain_name` and `hosted_zone_id`; create an alias record with `AwsRoute53DnsRecord` pointing your domain at those outputs.

## When to Use

- Give a production API a stable, branded URL with your own TLS certificate.
- Consolidate several microservice APIs under one hostname with path-key namespacing.
- Enforce **mutual TLS** for B2B / machine-to-machine APIs (clients must present certificates chaining to your CA truststore).

## Prerequisites

- An ACM certificate **in the same region** covering the domain name (exact or wildcard match) -- compose with `AwsCertManagerCert`.
- An HTTP API to map -- compose with `AwsHttpApiGateway` (its `api_id` and `stage_name` outputs are the mapping's inputs).
- (Optional) A Route 53 zone to publish the alias record -- compose with `AwsRoute53Zone` + `AwsRoute53DnsRecord`.

## Minimal Example

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsHttpApiDomain
metadata:
  name: api-example-com
  org: my-org
  env: prod
  id: api-example-com-prod
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

## Multi-API Example

```yaml
spec:
  region: us-east-1
  domainName: api.example.com
  certificateArn:
    valueFrom:
      kind: AwsCertManagerCert
      name: api-example-cert
      fieldPath: status.outputs.cert_arn
  apiMappings:
    # https://api.example.com/ -> the public API
    - apiId:
        valueFrom:
          kind: AwsHttpApiGateway
          name: public-api
          fieldPath: status.outputs.api_id
      stage: $default
    # https://api.example.com/orders/... -> the orders service
    - apiId:
        valueFrom:
          kind: AwsHttpApiGateway
          name: orders-api
          fieldPath: status.outputs.api_id
      stage: $default
      apiMappingKey: orders
```

## Spec Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `region` | string | Yes | AWS region. Must match the certificate's and the mapped APIs' region. |
| `domainName` | string | Yes | The fully qualified custom domain (lowercase; wildcard `*.example.com` allowed). Immutable -- changing it replaces the domain. |
| `certificateArn` | StringValueOrRef | Yes | ACM certificate covering the domain, in the same region. |
| `ipAddressType` | string | No | `ipv4` or `dualstack` (AWS default). |
| `mutualTls.truststoreUri` | string | Conditional | S3 URI of the PEM CA bundle (e.g. `s3://bucket/truststore.pem`). Required when the block is present. |
| `mutualTls.truststoreVersion` | string | No | S3 object version pin -- makes CA rotation an explicit change. |
| `apiMappings[].apiId` | StringValueOrRef | Yes | The API to map (references `AwsHttpApiGateway.status.outputs.api_id`). |
| `apiMappings[].stage` | string | Yes | The API stage to serve (typically `$default`). |
| `apiMappings[].apiMappingKey` | string | No | Path key (no slashes). Empty serves the API at the domain root. Keys must be unique across mappings. |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `domain_name` | The custom domain name (the domain's join key). |
| `domain_name_arn` | ARN of the domain name resource. |
| `target_domain_name` | The API Gateway regional domain to alias/CNAME to (e.g. `d-abc123.execute-api.us-east-1.amazonaws.com`). |
| `hosted_zone_id` | Route 53 hosted zone ID of the regional endpoint -- the alias target zone. |

## Publishing DNS

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsRoute53DnsRecord
metadata:
  name: api-alias
spec:
  # ... zone reference ...
  recordType: A
  alias:
    dnsName:
      valueFrom:
        kind: AwsHttpApiDomain
        name: api-example-com
        fieldPath: status.outputs.target_domain_name
    hostedZoneId:
      valueFrom:
        kind: AwsHttpApiDomain
        name: api-example-com
        fieldPath: status.outputs.hosted_zone_id
```

## Security Notes

- When mTLS is enabled, also set `disableExecuteApiEndpoint: true` on every mapped API -- otherwise callers can bypass mTLS via the default execute-api endpoint.
- API Gateway v2 domains accept only the REGIONAL endpoint type and TLS 1.2 security policy; the modules set both (they are not spec fields because AWS accepts no other value).

## Deliberately Omitted

- **Routing mode / routing rules** (`aws_apigatewayv2_routing_rule`): header/path rule-based routing is a separate resource family; the default API-mapping-only mode covers the mapping surface this component models. Revisit on concrete pull.
- **Edge-optimized endpoint type**: a REST API (v1) feature; v2 domains are REGIONAL only.
- **Per-kind tags**: identity tags derive from metadata.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
