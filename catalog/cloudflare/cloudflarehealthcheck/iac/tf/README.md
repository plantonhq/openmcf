# CloudflareHealthcheck Terraform Module

Terraform IaC module for a standalone origin health check -- scheduled probes with healthy/unhealthy status, no load balancer required.

## Architecture

```
provider.tf   — Cloudflare provider configuration (~> 5.23)
variables.tf  — Input variables mirroring CloudflareHealthcheckSpec
locals.tf     — http_config / tcp_config shaping (headers wrapper → header)
main.tf       — cloudflare_healthcheck
outputs.tf    — healthcheck_id, zone_id
```

## Behavior

`type` is HTTP, HTTPS, or TCP. `http_config` is sent only for HTTP/HTTPS and `tcp_config` only for TCP (the unused block is never sent -- both are Computed upstream). The spec's `headers` wrapper is unwrapped to the provider's `header` map. Health checks are a paid zone feature; the API enforces the plan gate. Destroy is a real delete.

## Outputs

| Name | Description |
|------|-------------|
| `healthcheck_id` | The created health check's ID |
| `zone_id` | The zone the health check belongs to |

## Provider Version

Uses `cloudflare/cloudflare ~> 5.23`.
