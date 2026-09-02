# CloudflareAiGateway Terraform Module

Terraform IaC module for one AI Gateway and its dynamic routes.

## Architecture

```
provider.tf   — Cloudflare provider configuration (~> 5.23)
variables.tf  — Input variables mirroring CloudflareAiGatewaySpec (generated)
locals.tf     — retry/log_management fan-out, guardrail empty-drop,
                route-graph reshaping (on_true/on_false -> true/false,
                provider -> ai_gateway_dynamic_routing_provider)
main.tf       — cloudflare_ai_gateway + cloudflare_ai_gateway_dynamic_routing (for_each by route name)
outputs.tf    — gateway_id, dynamic_route_ids
```

## Behavior

The gateway's `id` argument is the user-chosen URL slug (create-only). The spec's grouped retry{} / log_management{} fan out to the provider's flat arguments; unset guardrail controls are stripped rather than sent as empty strings; spend-limit rule ids are always sent explicitly (the provider default is a leaked example value that collapses rules). Each dynamic_routes row is its own provider object keyed by route name -- the elements list forces replacement on any change, so a graph edit replaces only its own route. Destroy is a real delete of the gateway and every route.

The module always sends `authentication`, `logpush`, `zdr`, and the log-management pair even when the manifest omits them: Cloudflare echoes false / 100000 / DELETE_OLDEST for the unset forms on every read (live-measured 2026-08-27), so an omitted send would drift forever against the echo. The sent defaults ARE Cloudflare's own -- semantics are unchanged.

**Known upstream wall (v5.23.0-v5.24.0, live-measured 2026-08-27)**: the API is write-only for `guardrails`, `spend_limits`, `otel`, and route graphs -- no read returns them, the provider refreshes them anyway, and plans never settle: every gateway re-plans an in-place no-op update forever (the provider re-marks the computed `otel`/`spend_limits` unknown on every plan, with or without config values), and a managed dynamic route is destroyed and recreated on every apply (its un-restorable `elements` list forces replacement). `ignore_changes` cannot suppress a provider-planned unknown (measured). The kind's e2e profile carries the full evidence and unblock condition; the GUIDE carries the customer-facing framing.

## Outputs

| Name | Description |
|------|-------------|
| `gateway_id` | The gateway's URL slug |
| `dynamic_route_ids` | Each managed route's id, keyed by route name |

## Provider Version

Uses `cloudflare/cloudflare ~> 5.23`.
