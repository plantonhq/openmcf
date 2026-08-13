# Preset: Rule-Routed Domain

## When to Use

Use this preset when one hostname must route requests to DIFFERENT REST APIs by path or header — a service split (`/orders` to the orders API), tenant pinning (beta tenants to a canary API), or a gradual migration — instead of the static one-path-one-API mapping model. API Gateway routing rules invoke ONLY REST-protocol APIs (HTTP/WebSocket targets are rejected); for HTTP APIs, stay on the mapping model.

## Key Configuration Choices

- **`ROUTING_RULE_ONLY`** — rules are evaluated by ascending priority and the first match wins; anything unmatched receives API Gateway's 404. A rule-routed domain carries NO `apiMappings`: AWS rejects HTTP-API mappings on rule-mode domains (live-verified), so make the rule set cover every path family the domain serves.
- **Rule targets are literal REST API ids** — routing rules support only REST-protocol APIs (live-verified), and the catalog does not yet model REST APIs, so `apiId` takes the id value directly rather than a reference.
- **Priorities are sparse (10, 20, ...)** — leave room to insert rules later without renumbering; AWS rejects two rules with the same priority.
- **`stripBasePath: true` on the path rule** — the orders API receives `/list`, not `/orders/list`, so the API's routes stay path-prefix-agnostic.
- **Conditions on one rule are ANDed** — add a `basePaths` condition and a `header` condition to the same rule to require both; separate rules are ORed by priority order.

## What to Customize

1. **`<domain-resource-name>`** — Planton resource name (e.g., `api-example-com`)
2. **`api.example.com`** — The fully qualified domain (lowercase; must be covered by the certificate)
3. **`<certificate-resource-name>`** — Name of the AwsCertManagerCert resource in the same region
4. **The `apiId` values (`a1b2c3d4e5` / `f6a7b8c9d0`) and stages** — Your REST APIs' ids and stage names
5. **The conditions** — base paths are exact first-path-segment matches (case-sensitive); header globs support `prefix-*`, `*-suffix`, and `*infix*` forms

## After Deploying

Publish DNS exactly as in the single-API preset (an `AwsRoute53DnsRecord` alias on the exported `target_domain_name` / `hosted_zone_id`), and consider `disableExecuteApiEndpoint: true` on every routed API so callers cannot bypass the domain's routing policy.
