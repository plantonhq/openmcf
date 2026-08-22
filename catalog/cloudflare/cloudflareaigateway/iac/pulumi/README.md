# CloudflareAiGateway Pulumi Module

Pulumi (Go) IaC module for one AI Gateway and its dynamic routes.

## Architecture

```
main.go                     — Entrypoint loading the stack input
module/main.go              — Resources(): provider setup, gateway then routes
module/locals.go            — Locals initialization
module/ai_gateway.go        — cloudflare.AiGateway + nested-tree builders
module/dynamic_routes.go    — cloudflare.AiGatewayDynamicRouting per route row
module/outputs.go           — Stack output keys
```

## Behavior

Mirrors the Terraform module's contract exactly: the spec's grouped retry{} / log_management{} fan out to the provider's flat arguments; the spec's on_true/on_false route edges map to the wire's true/false; spend-limit rule ids are always sent explicitly (the provider default is a leaked example value that collapses rules); otel and stripe authorization values are kept secret in Pulumi state (`pulumi.ToSecret`) though the provider leaves them unmarked. Routes attach through the created gateway's id, so they create after it and die with it.

## Outputs

| Name | Description |
|------|-------------|
| `gateway_id` | The gateway's URL slug |
| `dynamic_route_ids` | Each managed route's id, keyed by route name |

## Provider Version

Uses `pulumi-cloudflare` SDK v6 (Cloudflare Terraform provider v5 line).
