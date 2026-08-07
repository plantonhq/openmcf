# AWS HTTP API Domain

Deploys an API Gateway v2 custom domain name — the production front door for HTTP APIs ([AwsHttpApiGateway](/cloud-catalog/aws-http-api-gateway)). It binds an owned domain (e.g. `api.example.com`) to an ACM certificate ([AwsCertManagerCert](/cloud-catalog/aws-cert-manager-cert)) and maps one or more APIs onto the domain through API mappings. The domain is deliberately its own resource: it outlives any one API, many APIs compose onto it under distinct path keys ("orders", "billing"), and the certificate rotates on its own lifecycle. DNS is composed, not embedded — the domain exports `target_domain_name` and `hosted_zone_id`, and an [AwsRoute53DnsRecord](/cloud-catalog/aws-route53-dns-record) alias points your name at them.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **API Gateway v2 Custom Domain** -- the REGIONAL endpoint with the TLS_1_2 security policy (v2 domains support only these, so the modules set both)
- **Certificate Binding** -- TLS termination with the referenced ACM certificate (same region, covering the domain name)
- **Mutual TLS Configuration** -- the S3-hosted truststore requiring client certificates, when configured
- **API Mappings** -- one binding per row, each serving an API's stage at the root or under a path key
- **AWS Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied automatically for tracking and governance

## Before You Deploy

### Planton Setup

- **AWS Provider Connection** -- an active connection in the Connect module with credentials for the target AWS account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Certificate first** -- deploy the [AwsCertManagerCert](/cloud-catalog/aws-cert-manager-cert) (in the same region, covering the domain) before the domain and reference its `cert_arn` output.
- **APIs to map** -- the [AwsHttpApiGateway](/cloud-catalog/aws-http-api-gateway) resources whose `api_id` outputs the mappings reference.

### AWS Account

- **A validated certificate** -- an unvalidated ACM certificate cannot bind; DNS-validated certificates in the same environment are the composed path.
- **For mTLS** -- a versioned S3 bucket (same region) holding the PEM truststore bundle, and `disable_execute_api_endpoint` set on the mapped APIs so callers cannot bypass the client-certificate check.

## Deploy

### Console

Open the deployment store, find **AWS HTTP API Domain**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Single API Domain** preset in the [Presets](#presets) tab to pre-populate the common shape.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsHttpApiDomain
metadata:
  name: api-example-com
  org: acme-corp
  env: prod
spec:
  region: us-west-2
  domainName: api.example.com
  certificateArn:
    valueFrom:
      kind: AwsCertManagerCert
      name: api-example-com-cert
      fieldPath: status.outputs.cert_arn
  apiMappings:
    - apiId:
        valueFrom:
          kind: AwsHttpApiGateway
          name: orders-api
          fieldPath: status.outputs.api_id
      stage: $default
```

```shell
planton apply -f http-api-domain.yaml
```

This binds `api.example.com` to its certificate and serves the orders API at the domain root. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, the certificate and APIs deploy first, then the domain, then the alias record that routes traffic:

```yaml
# In an AwsRoute53DnsRecord manifest (an A alias):
spec:
  # api.example.com ALIAS -> the domain's API Gateway endpoint
  aliasTarget:
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

## Key Configuration

These are the most important decisions when configuring a custom domain. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**The domain name is the one-way door** -- create-time immutable; changing it replaces the domain, and the DNS alias must re-point at the new target. A wildcard (`*.example.com`) matches all first-level subdomains but requires a certificate covering the wildcard.

**The certificate contract** -- issued (or imported) in the SAME region as the domain, covering the domain name exactly or by wildcard. Unlike CloudFront-backed REST domains, the certificate lives in the domain's own region — not us-east-1. Swapping certificates rotates in place.

**API mappings compose the domain** -- an empty path key serves an API at the root; distinct keys serve others under `https://<domain>/<key>/`. Keys are unique (only one root), single-segment (no slashes), and edit in place as APIs come and go.

**Mutual TLS is half a lock by itself** -- the S3 truststore requires client certificates at the domain edge, but each mapped API's default execute-api endpoint bypasses it. Disable that endpoint on the mapped APIs whenever mTLS is on.

## Outputs and Dependencies

### What This Component Consumes

References an [AwsCertManagerCert](/cloud-catalog/aws-cert-manager-cert) (`cert_arn`, required) and [AwsHttpApiGateway](/cloud-catalog/aws-http-api-gateway) resources (`api_id`) through its mappings.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `domain_name` | The custom domain name | Join key for downstream automation |
| `domain_name_arn` | Amazon Resource Name of the domain | IAM policies and tag-based governance |
| `target_domain_name` | The API Gateway-managed regional domain to target from DNS | AwsRoute53DnsRecord alias target |
| `hosted_zone_id` | The Route 53 zone of the API Gateway endpoint | AwsRoute53DnsRecord alias target zone |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Single API front door** -- one API at the domain root; the everyday production shape. Start from the **Single API Domain** preset.

**Composed multi-API domain with mTLS** -- several teams' APIs under path keys, client certificates required at the edge. Start from the **Multi API mTLS Domain** preset.

## Works With

- [**AWS HTTP API Gateway**](/cloud-catalog/aws-http-api-gateway) -- the APIs mapped onto this domain (references `api_id`)
- [**AWS Cert Manager Cert**](/cloud-catalog/aws-cert-manager-cert) -- TLS termination (references `cert_arn`)
- [**AWS Route53 DNS Record**](/cloud-catalog/aws-route53-dns-record) -- the alias that routes the domain to API Gateway (consumes `target_domain_name` + `hosted_zone_id`)
- [**AWS HTTP API VPC Link**](/cloud-catalog/aws-http-api-vpc-link) -- lets the mapped APIs reach private backends behind this front door
