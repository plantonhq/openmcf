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

## Outputs

| Name | Description |
|------|-------------|
| `gateway_id` | The gateway's URL slug |
| `dynamic_route_ids` | Each managed route's id, keyed by route name |

## Provider Version

Uses `cloudflare/cloudflare ~> 5.23`.
