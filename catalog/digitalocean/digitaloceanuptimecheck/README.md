# DigitalOcean Uptime Check

Built for 100% parity with the Terraform DigitalOcean provider's `digitalocean_uptime_check` and `digitalocean_uptime_alert` resources at the pinned provider version.

## What this component models

An availability/latency probe on an EXTERNAL endpoint, run from DigitalOcean's global vantage regions, plus its alert rules. The alert rows are composed here -- one alert resource per row -- because they cannot exist without the check, and because DigitalOcean's standalone alert resource leaves the parent check id mutable (re-pointing an alert orphans it on its old check, a corruption class this composition makes unrepresentable).

The component covers both resources' full argument surfaces:

- `check_name` -- the display name
- `target` -- a URL for http/https probes, a hostname or IP for ping (DigitalOcean enforces the pairing)
- `type` -- optional; `ping`, `http`, or `https` (unset defers to DigitalOcean's default, https)
- `regions` -- the vantage regions (`us_east`, `us_west`, `eu_west`, `se_asia`); required here although DigitalOcean can default it, because the provider never reconciles a defaulted set (see below)
- `enabled` -- optional; unset defers to DigitalOcean's default (enabled)
- `alerts` -- the rules: `alert_name`, `type` (`latency`, `down`, `down_global`, `ssl_expiry`), optional `threshold` (milliseconds for latency -- required there by validation -- and days-before-expiry for ssl_expiry), optional `comparison` (`greater_than`/`less_than`), optional `period` (`2m`...`1h`), and `notifications` (emails and/or Slack rows; at least one channel required)

## Quick start

```yaml
apiVersion: digital-ocean.planton.dev/v1alpha1
kind: DigitalOceanUptimeCheck
metadata:
  name: homepage-check
spec:
  checkName: homepage
  target: https://www.example.com
  regions:
    - us_east
    - eu_west
  alerts:
    - alertName: homepage-down
      type: down_global
      notifications:
        emails:
          - ops@example.com
```

Deploy with either provisioner; both produce identical resources and outputs.

## Outputs

| Output | Description |
|---|---|
| `check_id` | UUID of the check (the API identity, and the import id) |

## Behavior worth knowing

- **Regions are always declared.** Upstream, an omitted region set is filled by the API and then read back into state, leaving every subsequent plan trying to remove what DigitalOcean chose -- a perpetual diff. Requiring the field kills that class outright.
- **Alert rows destroy with the check.** DigitalOcean deletes a check's alerts with it; the composed rows mirror that lifecycle exactly.
- **Alert row ids are not outputs.** The composed rows import as `{check_id},{alert_id}`; the import map says where to find each alert id.
- **Slack webhook URLs are credentials.** DigitalOcean does not mark them sensitive; this spec does, and both provisioners keep them out of plain-text state rendering.
- **The two comparison spellings never mix.** Uptime alerts speak snake_case (`greater_than`); monitor alerts speak CamelCase (`GreaterThan`). Two different DigitalOcean APIs, deliberately not unified.

## Module layout

- `iac/tf/` -- OpenTofu/Terraform module (provider pinned `~> 2.99`)
- `iac/pulumi/` -- Pulumi module (Go, pulumi-digitalocean SDK)
- Both engines wire the same spec fields and export the same outputs; behavioral parity is the contract.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
