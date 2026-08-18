# AWS AppSync API

Managed APIs without servers: GraphQL APIs that resolve fields straight out of DynamoDB, Lambda, OpenSearch, EventBridge, Aurora, or any HTTPS endpoint — and Events APIs that give browsers and mobile apps real-time pub/sub over WebSockets with nothing to operate.

## What Gets Managed

- The API, as exactly one of AWS's two AppSync models:
  - **GraphQL**: authorization (one primary + additional providers), the SDL schema, individually managed types, APPSYNC_JS/VTL functions and resolvers, the server-side cache, limits, X-Ray, metrics, logging, a WAF web ACL, and the MERGED federation variant with its source APIs.
  - **Events**: per-phase authorization (connect / publish / subscribe) and channel namespaces with inline code handlers or direct data-source integrations.
- Data sources for either model: DynamoDB, Lambda, HTTP(S), OpenSearch, EventBridge, Aurora Data API, Bedrock runtime, or NONE (local logic).
- API keys (for the API_KEY auth mode) with expiry.
- A custom domain fronting the API (ACM certificate in us-east-1).

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with AppSync permissions (plus IAM PassRole for data-source service roles, and WAF permissions when attaching a web ACL).

### AWS Prerequisites

- Backends for real data sources (a table, function, bus, domain, or cluster) and an IAM role AppSync can assume to reach them — NONE and unsigned HTTP data sources need neither.
- For a custom domain: an ACM certificate in us-east-1 covering the domain.

## After You Deploy

- GraphQL: applications point at the `graphql_url` output (subscriptions at `realtime_url`); fetch API key secrets from the AWS console/CLI — AWS shows them only at creation.
- Events: clients publish through `events_http_endpoint` and subscribe through `events_realtime_endpoint`.
- With a custom domain: point your DNS at `appsync_domain_name` (a CNAME, or a Route53 alias via `domain_hosted_zone_id`).

## Day-2 Operations

- Schema, resolvers, functions, auth, limits, and logging all update in place; resolver changes on one API apply serially (AWS's own consistency lock) — slow on big batches, not stuck.
- The cache replaces on encryption-flag changes (a cold cache, not an outage) and bills per hour while it exists.
- EventBridge data sources are replace-to-change at this provider pin (rename the entry); other data sources update in place.
- API keys rotate by adding a new key, rolling clients, then removing the old one; maximum key life is 365 days.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
