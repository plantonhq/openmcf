# AwsRestApiGateway — Component Guide

Authored operational judgment for the REST API Gateway component: the
design decisions behind the spec's shape, and what to know before
running REST APIs in production.

## Design decisions

- **Typed routes XOR OpenAPI.** AWS accepts both; mixing them silently
  overwrites. The spec enforces exactly one so a stray `routes` block
  cannot clobber an imported document.
- **The resource tree is derived, not declared.** Paths like
  `/users/{id}/orders` become nested API Gateway resources (max five
  segments). Authors write the HTTP surface; the modules own the
  parent/child wiring.
- **Explicit deployment, hashed from the definition.** REST APIs do
  not auto-deploy. The modules create one deployment whose trigger is
  a hash of the full API definition, so every spec change redeploys —
  the declarative behavior a Planton component owes its users.
- **One stage.** Planton resources are already environment-scoped.
  Canary traffic shifting needs two live deployments and is a deploy
  workflow, not a resource field.
- **Satellites are named and referenced.** Authorizers, models, and
  validators live on the spec and routes point at them by name, so
  the same TOKEN authorizer can guard many methods without duplication.
- **Account-level CloudWatch logging is a region singleton.** Stage
  access logs need no account role; method-level execution logging
  (`method_settings.logging_level`) does. That role is not modeled
  here.

## Running REST APIs in production

- **Prefer typed routes for day-to-day APIs.** OpenAPI import is the
  right tool when the contract already lives in a document; typed
  routes are the right tool when Planton is the source of truth.
- **MOCK integrations are first-class.** Use them for health checks
  and contract stubs — they cost nothing and prove the tree without a
  backend.
- **Keep TOKEN authorizers on a short TTL.** `result_ttl_seconds`
  caches the Lambda's decision; a stolen token lives that long.
- **Resource policies are the PRIVATE-endpoint contract.** An
  `endpoint_configuration.type` of PRIVATE without a policy that
  admits the VPC endpoint is an API nobody can call.
- **Client certificates are for backend trust, not caller auth.**
  Generate one when the HTTP backend must verify that traffic came
  through this API; distribute the PEM to the backend trust store.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
