# AWS HTTP API Domain: Architecture and Design

## What a Custom Domain Is

API Gateway's default endpoints (`https://{api-id}.execute-api.{region}.amazonaws.com`) are stable but unbranded, tied to the API's lifecycle, and impossible to put behind your own TLS posture. An API Gateway v2 custom domain name binds an owned DNS name to an ACM certificate and serves as the durable front door: APIs are published onto it through API mappings, and DNS points at the AWS-managed regional target it exposes.

## Design Decisions

### 1. First-Class Resource with Folded Mappings

The domain is its own resource because it outlives any one API: certificates rotate, APIs are replaced, the hostname persists. API mappings, by contrast, are pure domain-scoped glue -- a mapping cannot exist without its domain, nothing else references a mapping, and AWS keys them by (domain, mapping key). They are therefore FOLDED into the domain spec as `api_mappings`, each entry materialized as its own `aws_apigatewayv2_api_mapping` resource in both engines, addressed by its path key. This mirrors the settled folded-satellite pattern used across the catalog (per-name materialization).

### 2. DNS Composed, Never Embedded

The domain does not create DNS records. It exports `target_domain_name` (the AWS-managed regional domain) and `hosted_zone_id` (the alias target zone), and an `AwsRoute53DnsRecord` publishes the alias. This keeps zone ownership with the DNS family and lets domains front APIs in accounts or zones the API team does not control.

### 3. Certificate by Reference

The certificate is an `AwsCertManagerCert` reference (or a literal ARN). AWS requires it to be in the same region as the domain and to cover the domain name (exact or wildcard). Certificate issuance, validation, and rotation live entirely with the certificate resource.

### 4. Endpoint Type and Security Policy Are Not Spec Fields

API Gateway v2 domains accept exactly one endpoint type (`REGIONAL`) and one security policy (`TLS_1_2`) -- the provider validates each against a single-value enum. Both are hardcoded in the modules; exposing them would be two decorative knobs a user could never meaningfully turn. (Edge-optimized domains are a REST API v1 feature.)

### 5. Mutual TLS as an Optional Fold

mTLS configuration (an S3-hosted PEM truststore + optional object-version pin) is domain-scoped configuration with the domain's lifecycle -- folded. The version pin matters operationally: without it, overwriting the S3 object silently changes which client CAs are trusted; with it, CA rotation is an explicit spec change. When mTLS is on, the mapped APIs should set `disable_execute_api_endpoint: true`, because the default endpoint bypasses the domain (and therefore mTLS) entirely.

### 6. Routing Mode Deferred

`routing_mode` (API-mapping-only vs routing-rule-first) exists to support the newer `aws_apigatewayv2_routing_rule` resource family -- header/path rule-based routing across APIs. That is a separate routing surface with its own lifecycle and priority model; until it is built, `routing_mode` has no meaningful non-default value, so it is deliberately not modeled.

## Mapping Semantics

- An **empty mapping key** serves the API at the domain root: `https://api.example.com/`.
- A **non-empty key** namespaces the API: key `orders` serves it at `https://api.example.com/orders/...`, with the key stripped before the request reaches the API's routes.
- Keys must be **unique** per domain (spec-enforced) and cannot contain slashes (nested keys are not supported by API Gateway v2).
- The `stage` is the API's stage name -- for `AwsHttpApiGateway`-managed APIs this is the exported `stage_name` output (typically `$default`).

## Lifecycle Notes

- **Creation** waits for the domain to reach `AVAILABLE` (typically under a minute for REGIONAL domains; ownership verification for wildcard/mTLS setups can extend it).
- **domain_name is immutable** -- renaming replaces the resource and its mappings; the DNS alias must follow.
- **Certificate swaps are in-place** -- pointing `certificate_arn` at a new certificate updates the domain without replacement.

## Dependencies

### Upstream

| Resource | Field | Relationship |
|----------|-------|-------------|
| AwsCertManagerCert | `certificate_arn` | Required -- TLS termination |
| AwsHttpApiGateway | `api_mappings[].api_id` | Optional -- the APIs published on the domain |

### Downstream

| Resource | Output Used | Use Case |
|----------|-------------|----------|
| AwsRoute53DnsRecord | `target_domain_name`, `hosted_zone_id` | The alias record publishing the domain |
