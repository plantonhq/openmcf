# Azure Traffic Manager Profile -- Operational Guide

Judgment calls that matter when you run Traffic Manager in production.

## The TTL is your failover clock

Traffic Manager stops handing out an unhealthy endpoint immediately -- but clients keep using cached answers until the TTL drains. Detection time is `interval x tolerated failures + timeout`; ADD the TTL on top for the real user-visible failover window. A 30s TTL with default probing means roughly two minutes worst-case. Buy faster failover with the 10s fast interval (billed extra, and set the timeout explicitly to 5-9) and a low TTL -- and remember some resolvers ignore very low TTLs.

## Pick the probe that proves the SERVICE, not the socket

TCP probes pass when the port accepts a handshake -- a hung application behind a live listener stays "healthy" forever. Probe HTTP/HTTPS with a real health path wherever the endpoint speaks HTTP, and scope `expected_status_code_ranges` deliberately: expecting only 200 while the app answers 301 to probes is the classic all-endpoints-degraded false alarm. Use a probe `Host` header when the target serves name-based virtual hosts.

## When ALL endpoints are degraded, Traffic Manager answers with all of them

The service's documented fallback: if every endpoint probes unhealthy, queries are answered as if all were healthy (so a monitoring misconfiguration degrades to "no steering" rather than "no answers"). If your endpoints "work" while the portal shows everything degraded, fix the probe -- the answers you are seeing are the fallback, and real failover is silently broken.

## The relative name is a global, permanent choice

`{relative_name}.trafficmanager.net` is shared across every Azure customer and FIXED at creation -- renaming replaces the profile and breaks every CNAME pointing at it. Prefix with your organization, and treat the generated name as infrastructure: point your own domain at it with a CNAME and give users only your domain.

## Disabling beats deleting, at both levels

`enabled: false` on the profile parks the whole ruleset (the name stops steering); the endpoint component's own enabled flag drains one destination. Maintenance windows, migrations, and kill switches are one-field flips -- deletion loses the globally-unique name to a window where anyone can claim it.

## MultiValue is for resolver-side logic, not load balancing

MultiValue answers with up to `max_return` healthy addresses AT ONCE and requires literal-IP external endpoints. Use it when the CLIENT retries across addresses (DNS-savvy SDKs, custom resolvers). For actual traffic distribution, Weighted with per-endpoint weights is the tool.
