# AWS EventBridge API Destination

Authenticated webhooks as infrastructure: EventBridge signs in to your SaaS or internal API (api-key, basic, or OAuth client credentials), respects a rate limit you set, and delivers events to an HTTPS endpoint — no credential-holding Lambda in between.

## What Gets Managed

- The connection: the auth mode and its credentials (stored by AWS in a Secrets Manager secret it owns), static parameters added to every invocation, optional customer-managed encryption, and optional private endpoints through VPC Lattice.
- The API destination: the HTTPS endpoint (with `*` path wildcards), the HTTP method, and the invocations-per-second cap. Deploy both together, or share one connection across many destination-only instances.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with EventBridge and Secrets Manager permissions.

### AWS Prerequisites

- The downstream API's credentials in hand — they enter the spec as secret-typed fields and land in an AWS-owned Secrets Manager secret.

## After You Deploy

- The connection walks an auth state machine (CREATING/AUTHORIZING → AUTHORIZED) — budget minutes on first deploys.
- Nothing is invoked until an EventBridge rule, pipe, or schedule targets the destination's ARN. Invocations beyond the rate limit queue for up to 24 hours.

## Common Changes

- Rotate credentials: update the secret fields — AWS re-authorizes the connection in place.
- Throttle-tune: `invocation_rate_limit_per_second` updates in place; size it to what the downstream API tolerates.
- Share the trust anchor: keep one connection-only instance and point new destinations at its `connection_arn` output.
