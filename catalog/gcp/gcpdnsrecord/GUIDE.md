# GcpDnsRecord Guide

The judgment this guide protects: DNS answers are cached promises — TTL
decides how fast a change (or a mistake) propagates, and routing
policies move real traffic. Edit records with the same care as load
balancer configs, because for your users they are the same thing.

## TTL is your rollback speed

`ttlSeconds` bounds how long resolvers cache an answer — which is
exactly how long a bad change keeps hurting after you fix it. Before a
planned cutover, LOWER the TTL a full old-TTL ahead of time, make the
change, verify, then raise it again. Steady-state 300s is a sensible
default; multi-day TTLs belong only on records that genuinely never
move (NS conventions).

## One record set, one owner

A record set is (name, type) — the values round-robin inside it. Two
manifests writing the same name/type will fight via full overwrites, so
give every record set exactly one owning manifest. The zone reference
is a `StringValueOrRef` to GcpDnsZone — chart wiring keeps record and
zone in one dependency graph.

## Routing policies replace values, not augment them

Static `values` and `routingPolicy` are mutually exclusive (the API's
own contract, enforced pre-deploy). Weighted round robin's zero-weight
entries stage a target without traffic — the canary pattern: add at
weight 0, flip weights gradually. Geo fencing decides failure behavior:
fenced locations keep answering when unhealthy instead of spilling to
the next-closest — choose per record whether local-broken or
cross-region traffic is the lesser evil. Primary/backup's
`trickleRatio` keeps the backup path warm so failover is not its first
production traffic.

## Health-checked targets have two flavors

Internal load balancer targets carry an implicit health signal; public
endpoints need an explicit GcpHealthCheck reference. Either way the
policy only fails over on the signal it has — a "healthy" check that
does not exercise the real dependency chain gives false confidence.

## Teardown discipline

Deleting a record set stops answers as caches expire — user-visible at
TTL speed. `PREVENT` suits names other systems depend on to find you:
MX, domain-verification TXT, service-discovery entries. `ABANDON` keeps
the record answering while dropping management; remember an abandoned
record still overwrites-fights any future manifest that claims the same
name/type.
