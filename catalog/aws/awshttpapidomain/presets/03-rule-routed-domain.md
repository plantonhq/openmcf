# Preset: Rule-Routed Domain

## When to Use

Use this preset when one hostname must route requests to DIFFERENT APIs by path or header — a service split (`/orders` to the orders API), tenant pinning (beta tenants to a canary API), or a gradual migration — instead of the static one-path-one-API mapping model.

## Key Configuration Choices

- **`ROUTING_RULE_THEN_API_MAPPING`** — rules are evaluated first by ascending priority; anything unmatched falls back to the api_mappings, so the domain keeps a well-defined default.
- **Priorities are sparse (10, 20, ...)** — leave room to insert rules later without renumbering; AWS rejects two rules with the same priority.
- **`stripBasePath: true` on the path rule** — the orders API receives `/list`, not `/orders/list`, so the API's routes stay path-prefix-agnostic.
- **Conditions on one rule are ANDed** — add a `basePaths` condition and a `header` condition to the same rule to require both; separate rules are ORed by priority order.

## What to Customize

1. **`<domain-resource-name>`** — Planton resource name (e.g., `api-example-com`)
2. **`api.example.com`** — The fully qualified domain (lowercase; must be covered by the certificate)
3. **`<certificate-resource-name>`** — Name of the AwsCertManagerCert resource in the same region
4. **`<orders-api-resource-name>` / `<canary-api-resource-name>` / `<main-api-resource-name>`** — The AwsHttpApiGateway resources the rules and fallback route to
5. **The conditions** — base paths are exact first-path-segment matches (case-sensitive); header globs support `prefix-*`, `*-suffix`, and `*infix*` forms

## After Deploying

Publish DNS exactly as in the single-API preset (an `AwsRoute53DnsRecord` alias on the exported `target_domain_name` / `hosted_zone_id`), and consider `disableExecuteApiEndpoint: true` on every routed API so callers cannot bypass the domain's routing policy.
