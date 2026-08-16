# CloudflareSnippetRules Pulumi Module

Pulumi (Go) IaC module for a zone's snippet routing table -- the ordered list of expressions deciding which requests invoke which snippet. One manifest owns the whole table.

## Architecture

```
main.go                    — Entrypoint loading the stack input
module/main.go             — Resources(): provider setup, resource, outputs
module/locals.go           — Locals initialization
module/snippet_rules.go    — cloudflare.SnippetRules
module/outputs.go          — Stack output keys
```

## Behavior

Mirrors the Terraform module's contract exactly: full-replacement PUT, destroy wipes the whole table, `enabled` defaults true, `zone_id` stack output.

## Outputs

| Name | Description |
|------|-------------|
| `zone_id` | The zone whose snippet routing table is managed (the singleton's identity) |

## Provider Version

Uses `pulumi-cloudflare` SDK v6 (Cloudflare Terraform provider v5 line).
