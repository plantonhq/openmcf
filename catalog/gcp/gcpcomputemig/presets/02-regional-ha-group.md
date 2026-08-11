# Regional HA Group

A fixed-size, three-zone fleet with application-level auto-healing and
hardened boot — the availability-first posture for serving tiers whose
size is planned rather than demand-driven.

## What it configures

- EVEN distribution across three named zones (two instances per zone at
  the fixed size of 6) — a zone outage removes exactly one third.
- Auto-healing from a dedicated conservative GcpHealthCheck with a
  5-minute cold-start allowance; repairs stay in-zone and never apply a
  new template as a side effect.
- Shielded VM (secure boot + vTPM + integrity monitoring) on every
  instance.
- `PREVENT` teardown posture — destroying the tier is a deliberate,
  two-step act.

## Adjust before deploying

- **distributionPolicy.zones** — your region's zones; quota is needed in
  EVERY named zone.
- **autoHealing.healthCheck** — reference a GcpHealthCheck tuned for
  repairs (long interval, high unhealthy threshold) — never share the
  load balancer's aggressive probe.
- **targetSize** — keep it a multiple of the zone count so EVEN
  distribution is exact.

## After deploying

Watch repairs in the group's console panel; a repair storm means the
health check or `initialDelaySec` is too aggressive for the app's real
cold start.

## When to choose something else

For demand-driven sizing, start from **Autoscaled Web Tier** (the two
sizing modes are mutually exclusive by design). For per-instance
identity and preserved disks, start from **Stateful Group**.
