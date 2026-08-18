# Terraform Module: DigitalOcean Uptime Check

Provisions an uptime probe on an external endpoint plus its alert rules -- the complete `digitalocean_uptime_check` surface with one `digitalocean_uptime_alert` per spec alert row.

## Resources

| Resource | Purpose |
|---|---|
| `digitalocean_uptime_check.check` | The probe: target, protocol, vantage regions |
| `digitalocean_uptime_alert.alerts` | One per spec alert row (`for_each` keyed `"<index>-<alert_name>"`) |

## Inputs

Generated `variables.tf` mirrors the `DigitalOceanUptimeCheckSpec` proto: `check_name`, `target`, optional `type`, required `regions`, optional `enabled`, and the `alerts` rows (each with `alert_name`, `type`, optional `threshold`/`comparison`/`period`, and `notifications` channels). Authentication uses `digitalocean_token` (sensitive).

## Outputs

Exactly the `DigitalOceanUptimeCheckStackOutputs` contract: `check_id`. Alert row ids are not outputs (a check manages many); their import identities are `{check_id},{alert_id}`.

## Behavior notes

- `check_id` on each alert row is wired to the check this module creates -- the upstream mutable-parent corruption class is unrepresentable here.
- `regions` is always sent (spec-required): the provider never reconciles a DigitalOcean-defaulted region set, so omission would leave a perpetual removal diff.
- An unset alert `threshold` is sent as the API's accepted zero (down/down_global rows); validation requires it on latency rows.
- Import: check by `<check_id>`, alert rows by `<check_id>,<alert_id>` (see `iac/import-map.yaml`).
