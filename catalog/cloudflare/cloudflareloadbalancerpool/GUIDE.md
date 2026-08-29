# CloudflareLoadBalancerPool guide

Operational judgment for load-balancer origin pools. The README covers what each field is; this covers how the pieces interact.

## The paid add-on gates this too

Pools, monitors, and load balancers all ride the Load Balancing add-on. Without it every call returns `403`. Enable the subscription before deploying any of the three.

## A monitor is optional and that is a production foot-gun

A pool with no monitor treats every origin as healthy. That is acceptable for a throwaway or a pool you health-check some other way; it is not acceptable for production traffic. Attach a CloudflareLoadBalancerMonitor when the pool will actually serve.

## Origins under a monitor must be globally routable

Cloudflare's health checks run from its data centers. RFC 5737 documentation ranges and other non-routable addresses are rejected once a monitor is attached — the create fails with `400` code `1002` "Origin address is invalid or is not globally routable and has health monitoring enabled" (measured live). TEST-NET addresses are fine for an unmonitored pool and wrong for a monitored one.

## Names are short and unique per account

Pool names max out at 32 characters (`[A-Za-z0-9_-]+`) and must be unique in the account. A load balancer references pools by id, not by name, so a rename is cosmetic to the LB but still a unique-name constraint against every other pool.

## Steering lives on the load balancer

Latitude/longitude on the pool matter only when a load balancer uses `proximity` steering. `check_regions` is capped by plan tier — omit it to probe from every region rather than guessing a list that will 400.
