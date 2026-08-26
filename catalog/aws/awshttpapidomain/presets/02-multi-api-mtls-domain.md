# Multi-API mTLS Domain

## When to Use

Use this preset for B2B / machine-to-machine API surfaces: several APIs composed under one hostname, with mutual TLS requiring partner clients to present certificates chaining to your CA truststore.

## Key Configuration Choices

- **Mutual TLS** -- API Gateway rejects any client that does not present a certificate chaining to a CA in the S3-hosted truststore, before the request reaches any API.
- **Pinned truststore version** -- CA rotation becomes an explicit spec change instead of a silent side effect of overwriting the S3 object.
- **Path-key composition** -- distinct APIs (primary + webhooks) share the hostname under unique keys; each keeps its own routes, authorizers, and lifecycle.

## What to Customize

1. **`<partner-api.example.com>`** — The mTLS-protected domain (must be covered by the certificate)
2. **`<security-bucket>` / `<s3-object-version>`** — The truststore bundle location and pinned version
3. **`<primary-api-resource-name>` / `<webhooks-api-resource-name>`** — The AwsHttpApiGateway resources to compose
4. **`<certificate-resource-name>`** — Name of the AwsCertManagerCert resource in the same region

## Production Checklist

- [ ] `disableExecuteApiEndpoint: true` on EVERY mapped API — otherwise callers bypass mTLS via the default execute-api endpoint
- [ ] Truststore bucket versioning enabled (the version pin depends on it)
- [ ] Truststore bundle contains only the intended issuing CAs
- [ ] Route 53 alias record published from the domain's outputs
