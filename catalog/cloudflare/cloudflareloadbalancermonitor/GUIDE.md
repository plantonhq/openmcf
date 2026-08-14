# CloudflareLoadBalancerMonitor guide

Operational judgment for load-balancer health checks. The README covers what each field is; this covers how the pieces interact.

## The paid add-on gates this too

Monitors, pools, and load balancers all ride the Load Balancing add-on. Without it every call returns `403`. Enable the subscription before deploying any of the three.

## Type decides which fields matter

`http`/`https` use `path`, `method`, `expected_codes`, `expected_body`, `follow_redirects`, `allow_insecure`, and headers. `tcp` (and `udp_icmp` / `smtp`) require a `port > 0` and ignore the HTTP knobs. A tcp monitor with no port fails validation; an http monitor with a port is fine (it overrides the default 80/443).

## Zero means "Cloudflare's default"

`interval`, `timeout`, `retries`, `consecutive_up`, `consecutive_down`, and `port` left at 0 are omitted so the provider applies the server default. After import the API returns those defaults as numbers, so a post-import plan may show a diff on port / consecutive_* even though nothing operationally changed.

## Pools reference monitors, not the other way around

A monitor has no knowledge of the pools that use it. Deleting a monitor that a pool still names will fail or leave the pool unmonitored (origins then look permanently healthy). Delete pools first, or clear the pool's `monitor` field, then the monitor.
