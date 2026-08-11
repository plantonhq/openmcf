# Public HTTPS Check

The canonical availability monitor: probe a public URL over TLS from all
regions every five minutes, with certificate validation so an expired
certificate fails the probe instead of passing silently.

## What it configures

- `monitoredResource` of type `uptime_url` — the host to probe.
- `httpCheck` with `useSsl` + `validateSsl` — HTTPS with the handshake
  JUDGED, which is what makes this the certificate-expiry monitor too.
- All regions (no `selectedRegions`) — maximum coverage, zero config.

## Adjust before deploying

- **host** — the domain to probe (no scheme, no path — path lives on
  `httpCheck.path`).
- Consider a `contentMatchers` entry on a known-good body token so an
  error page served with status 200 still fails the probe.

## The other half

A check only measures. Pair it with a GcpMonitoringAlertPolicy threshold
condition on `uptime_check_passed` filtered by this check's
`uptime_check_id` output — wired via valueFrom so recreation never
strands the filter.
