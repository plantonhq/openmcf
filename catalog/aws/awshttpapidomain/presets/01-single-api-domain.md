# Single API Domain

## When to Use

Use this preset to give one HTTP API a branded production URL: `https://api.example.com/` serving the API's `$default` stage at the domain root.

## Key Configuration Choices

- **Certificate by reference** -- the ACM certificate rotates on its own lifecycle without touching the domain.
- **Root mapping** -- an empty mapping key serves the API at the domain root; add keyed mappings later to compose more APIs under the same hostname.
- **DNS composed downstream** -- publish the domain with an `AwsRoute53DnsRecord` alias pointing at the exported `target_domain_name` / `hosted_zone_id`.

## What to Customize

1. **`<domain-resource-name>`** — Planton resource name (e.g., `api-example-com`)
2. **`<api.example.com>`** — The fully qualified domain (lowercase; must be covered by the certificate)
3. **`<certificate-resource-name>`** — Name of the AwsCertManagerCert resource in the same region
4. **`<http-api-resource-name>`** — Name of the AwsHttpApiGateway resource to serve

## After Deploying

Create the alias record so the domain resolves:

```yaml
recordType: A
alias:
  dnsName:
    valueFrom:
      kind: AwsHttpApiDomain
      name: <domain-resource-name>
      fieldPath: status.outputs.target_domain_name
  hostedZoneId:
    valueFrom:
      kind: AwsHttpApiDomain
      name: <domain-resource-name>
      fieldPath: status.outputs.hosted_zone_id
```

Also consider `disableExecuteApiEndpoint: true` on the API so callers cannot bypass the domain.
