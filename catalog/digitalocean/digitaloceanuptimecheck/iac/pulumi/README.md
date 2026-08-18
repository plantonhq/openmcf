# Pulumi Module: DigitalOcean Uptime Check

Provisions an uptime probe on an external endpoint plus its alert rules -- the complete `digitalocean_uptime_check` surface with one `digitalocean_uptime_alert` per spec alert row, at 100% behavioral parity with the Terraform module (same arguments, same outputs).

## Layout

- `main.go` -- entrypoint (`package main`), loads the stack input and calls the module
- `module/main.go` -- orchestration: locals, provider, resources
- `module/uptime_check.go` -- the `UptimeCheck` resource, the per-row `UptimeAlert` children, and output exports
- `module/locals.go` -- target handle (a check has no tag surface, so no label set applies)
- `module/outputs.go` -- output key constants (the `DigitalOceanUptimeCheckStackOutputs` contract)

## Behavior notes

- Alert rows parent the check (`pulumi.Parent`) and take its id as `CheckId` -- the upstream mutable-parent corruption class is unrepresentable here.
- The SDK keeps the provider's unbounded notifications list, but the provider reads only the first element -- exactly one is sent per row.
- The Slack webhook URL is not secret-flagged by the SDK, so the module wraps it in `pulumi.ToSecret`.
- Resource names carry the row index (`alert-<idx>-<name>`) so two rows may share a display name without colliding.
