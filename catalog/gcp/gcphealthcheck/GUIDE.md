# GcpHealthCheck Guide

The judgment this guide protects: a health check is the load balancer's
eyes — probe the wrong thing and failover fires on database hiccups, or
never fires at all. One shared check often guards MANY backend services;
edit it like shared infrastructure.

## Probe a cheap truth, not a deep one

Point HTTP checks at a dependency-free endpoint (/healthz that answers
from memory). A handler that touches the database turns database latency
into load-balancer failovers — the check ejects healthy web servers
because a slow query timed out. Deep dependency checking belongs in
alerting (uptime checks), not in the traffic path.

## Scope is chosen once

Empty `region` = GLOBAL (global external LBs); set = REGIONAL (regional
backend services, MIG auto-healing) — and a check cannot move between
scopes. Regional backend services can only reference checks in their own
region. When in doubt about a future regional ALB, make the regional
check up front.

## Failover math is two dials multiplied

Detection time ≈ `checkIntervalSec` × `unhealthyThreshold` (default
5s × 2 = 10s); recovery ≈ interval × `healthyThreshold`. Tighten both
for latency-sensitive tiers, but remember every prober probes at that
interval from multiple sources — a 1s interval is real backend load.
`timeoutSec` must not exceed the interval (spec-enforced; the API
rejects it too).

## USE_SERVING_PORT is usually right

Probing the backend's own serving port removes the classic drift where
the app moves ports and the check keeps probing the old one. Fixed ports
are for backends that serve health on a dedicated port; named ports only
work on instance groups (and not at all for gRPC-with-TLS — the spec
enforces that asymmetry).

## source_regions hardens global verdicts

Pinning exactly 3 prober regions means one regional outage cannot flip a
global health verdict. The API's fine print is enforced by the spec: only
HTTP/HTTPS/TCP, interval ≥ 30s, no proxy header, no TCP payload — and a
check with source_regions cannot drive MIG auto-healing.

## Teardown discipline

GCP refuses to delete a check a backend service still references, so
destroy order is consumers-first. `PREVENT` suits the shared probe many
backend services point at (its recreation also briefly breaks every
referencing self-link); `ABANDON` keeps it probing while dropping
management. Enable `enableLogging` while tuning thresholds — every state
transition is a log line — and turn it off on large fleets afterward.
