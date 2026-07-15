# HTTPS Endpoint Health Check

This preset creates the standard internet-facing HTTPS probe: Route 53's global checker fleet requests your application's health endpoint and the check stays healthy while the endpoint answers 2xx/3xx. This is the health check that gates a failover PRIMARY record or drops an unhealthy weighted member out of rotation.

## When to Use

- Gating a failover routing pair's PRIMARY record
- Health-aware weighted routing (unhealthy members drop out of the split)
- Multivalue answer members that should disappear when down

## Key Configuration Choices

- **HTTPS probing** (`checkType: HTTPS`) -- TLS to port 443 by default; the check passes on any 2xx/3xx response
- **Dedicated health path** (`resourcePath: /healthz`) -- probe a purpose-built endpoint, not the homepage; keep it fast and dependency-light
- **Standard cadence** (`requestInterval: 30`, `failureThreshold: 3`) -- failure detected in ~90 seconds; switch to `10`/`2` for ~20–30 second detection at extra cost
- **Internet reachability** -- the checker fleet probes from eight AWS regions; the endpoint (and its security groups) must admit them

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<aws-region>` | Region for provider API calls (checks are global objects) | Your deployment region |
| `<endpoint-domain>` | The publicly resolvable domain to probe | Your application's DNS name |
