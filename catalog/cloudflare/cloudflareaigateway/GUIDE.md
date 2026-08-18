# CloudflareAiGateway guide

The judgment this guide protects you from: the slug is the URL every client calls, route graphs replace rather than update, and one leaked provider default can silently collapse your spend budgets.

## The gateway id IS the endpoint URL

`gateway_id` appears verbatim in the endpoint every client calls (`https://gateway.ai.cloudflare.com/v1/<account>/<gateway_id>/...`). It is create-only: renaming replaces the gateway AND breaks every client mid-flight. Choose the slug like a domain name -- boring, stable, forever.

## Route graphs replace, never update

A dynamic route's `elements` list is create-only at the provider: ANY edit to the graph -- one condition, one model swap -- plans a replacement of that route object. This is safe (requests re-resolve the route by name on the next call) but worth knowing when you read the plan: a replace on a route is the designed behavior, not a mistake. Renaming a route is its only in-place update -- and the provider has a decode defect there (computed fields go stale in state until the next refresh), so prefer stable route names too.

## Every spend-limit rule needs its own id -- here is why the spec forces it

The provider schema ships a LEAKED EXAMPLE VALUE as the default rule id: every rule authored without an id gets the same identity and Cloudflare silently collapses them into one budget. This spec makes `id` required and unique per rule, so the failure cannot be authored. Pick short stable slugs ("daily-cap", "monthly-cap") -- they are identities, not descriptions.

## Guardrails: both sides, uppercase wire codes

`guardrails` requires BOTH `prompt` and `response` objects (an empty object means "evaluate nothing on this side"). Each control names a hazard-category code from Cloudflare's Guardrails taxonomy (p1, s1-s13 -- the Cloudflare docs carry what each code covers) and takes FLAG (log) or BLOCK (refuse). On the wire the codes are uppercase; the modules handle the casing -- author lowercase per the spec.

## Two credentials the provider forgot to mark

`otel[].authorization` and `stripe.authorization` are credentials the provider leaves unmarked. This spec treats both as sensitive: provide managed-secret references, and the modules keep them secret in state. Never paste a bearer token or Stripe key as a literal.

## ZDR and logging pull in opposite directions

`zdr: true` (Zero Data Retention) tells Cloudflare to never store request/response bodies; `collect_logs`, `logpush`, and log management all exist to store or ship exactly that data. Cloudflare enforces the legal combinations server-side -- expect API rejections, not plan errors, when you ask for both. Decide the data posture first, then configure logging around it.

## The five required scalars are required on purpose

Cloudflare demands an explicit choice for caching, log collection, and rate limiting -- there are no server defaults. A "minimal" gateway still states all five. `cache_ttl: 0` and `rate_limiting_interval: 0` are the explicit OFF switches.

## Pairs well with

- [CloudflareSecretsStore](../cloudflaresecretsstore/README.md) -- the vault behind `store_id` (BYO provider keys).
- [CloudflareSecretsStoreSecret](../cloudflaresecretsstoresecret/README.md) -- the provider keys themselves, scoped `ai_gateway`.
- [CloudflareWorker](../cloudflareworker/README.md) -- Workers calling models through the gateway endpoint.
