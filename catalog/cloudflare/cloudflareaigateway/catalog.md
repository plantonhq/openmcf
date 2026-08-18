# Cloudflare AI Gateway

The control plane in front of AI model traffic: caching, rate limiting, retries, logging, guardrails, DLP, spend limits, and dynamic request routing, all on one gateway whose URL slug your clients call. Free tier -- the gateway is configuration, not compute.

## What Gets Created

When you deploy this resource, the IaC module provisions:

- **AI Gateway** -- one `cloudflare_ai_gateway` (the control plane and its endpoint slug)
- **Dynamic routes** -- one `cloudflare_ai_gateway_dynamic_routing` per `dynamicRoutes` row, attached to the gateway

## Prerequisites

- **A Cloudflare account** (free plans included)
- **A Cloudflare API token** with Account → AI Gateway → Edit
- For BYO provider keys: the account **Secrets Store** and the key secrets in it

## Quick Start

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareAiGateway
metadata:
  name: prod-llm-gateway
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  gatewayId: prod-llm-gateway
  cacheInvalidateOnUpdate: true
  cacheTtl: 300
  collectLogs: true
  rateLimitingInterval: 60
  rateLimitingLimit: 1000
```

```shell
planton apply -f ai-gateway.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `accountId` | string | The Cloudflare account. | Required, 32-hex. |
| `gatewayId` | string | The user-chosen URL slug. | Required; replaces on change (and changes the endpoint URL). |
| `cacheInvalidateOnUpdate` | bool | Invalidate cached responses on config updates. | Required -- explicit choice. |
| `cacheTtl` | int | Cache TTL seconds; 0 disables. | Required, >= 0. |
| `collectLogs` | bool | Collect request/response logs. | Required -- explicit choice. |
| `rateLimitingInterval` | int | Rate-limit window seconds; 0 disables. | Required, >= 0. |
| `rateLimitingLimit` | int | Requests per window; 0 disables. | Required, >= 0. |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `rateLimitingTechnique` | string | Cloudflare default | `fixed` or `sliding`. |
| `retry` | object | unset | `backoff` constant/linear/exponential, `delay` <= 5000ms, `maxAttempts` 1-5. |
| `logManagement` | object | unset | `maxRecords` 10k-10M + `strategy` STOP_INSERTING/DELETE_OLDEST. |
| `authentication` | bool | false | Require the cf-aig-authorization token per request. |
| `logpush` / `logpushPublicKey` | bool / string | unset | Logpush integration, optionally encrypted. |
| `zdr` | bool | false | Zero Data Retention -- never store bodies. |
| `workersAiBillingMode` | string | postpaid | Only `postpaid` today. |
| `storeId` | StringValueOrRef | unset | Secrets Store holding BYO provider keys. |
| `dlp` | object | unset | DLP screening (default action, profiles, policy rows). |
| `guardrails` | object | unset | Per-hazard-code FLAG/BLOCK controls (prompt and response sides both required inside). |
| `otel[]` | list | none | OpenTelemetry export destinations (authorization sensitive). |
| `stripe` | object | unset | Stripe usage reporting (authorization sensitive). |
| `spendLimits` | object | unset | Budget rules; every rule needs its own unique `id`. |
| `dynamicRoutes[]` | list | none | Named routing graphs; a graph edit recreates that route object. |

## Examples

### Gateway with a spend cap and a routing graph

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareAiGateway
metadata:
  name: routed-gateway
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  gatewayId: routed-gateway
  cacheInvalidateOnUpdate: true
  cacheTtl: 0
  collectLogs: true
  rateLimitingInterval: 60
  rateLimitingLimit: 500
  spendLimits:
    enabled: true
    rules:
      - id: monthly-cap
        limit: 500
        limitType: cost
        window: 2592000
  dynamicRoutes:
    - name: cheap-first
      elements:
        - id: start
          type: start
          outputs:
            next:
              elementId: check
        - id: check
          type: conditional
          outputs:
            onTrue:
              elementId: cheap
            onFalse:
              elementId: smart
          properties:
            conditions: '{"metadata.tier": {"$eq": "free"}}'
        - id: cheap
          type: model
          outputs:
            success:
              elementId: done
          properties:
            model: "@cf/meta/llama-3.1-8b-instruct"
            provider: workers-ai
        - id: smart
          type: model
          outputs:
            success:
              elementId: done
          properties:
            model: gpt-4o
            provider: openai
        - id: done
          type: end
          outputs: {}
```

## Destroy Semantics

Destroy is a real delete of the gateway and every dynamic route. Clients calling the endpoint URL fail immediately -- drain traffic first.

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `gateway_id` | string | The gateway's URL slug |
| `dynamic_route_ids` | map | Each managed route's id, keyed by route name |

## Related Components

- [Cloudflare Secrets Store](/docs/catalog/cloudflare/cloudflaresecretsstore) -- the vault behind BYO provider keys
- [Cloudflare Secrets Store Secret](/docs/catalog/cloudflare/cloudflaresecretsstoresecret) -- the keys themselves
- [Cloudflare Worker](/docs/catalog/cloudflare/cloudflareworker) -- Workers calling models through the gateway
