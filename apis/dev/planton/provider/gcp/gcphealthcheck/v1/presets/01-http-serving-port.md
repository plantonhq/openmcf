# HTTP Probe on the Serving Port

The workhorse health check for global external HTTP(S) load balancers: an HTTP GET against a dedicated health endpoint, probing whatever port each backend actually serves on (`USE_SERVING_PORT`) so the check keeps working when serving ports change.

## When to Use

- Backend services fronting serverless NEGs (Cloud Run, Cloud Functions) — this shape is effectively mandatory there
- Instance groups serving plaintext HTTP
- Any backend with a cheap `/healthz`-style endpoint

## Remix Notes

- Point `requestPath` at an endpoint that does NOT touch databases or downstream services — a health endpoint that inherits its dependencies' latency turns their blips into load-balancer failovers.
- Add `response: ok` to require the body to start with a known string — a guard against a wrong service answering 200 on the probed port.
- Tighten detection by lowering `checkIntervalSec` or `unhealthyThreshold` (detection ≈ interval × threshold), at the cost of more probe traffic.
- Instance-group backends need a firewall rule admitting Google's prober ranges (`35.191.0.0/16`, `130.211.0.0/22`) or every backend shows unhealthy.
