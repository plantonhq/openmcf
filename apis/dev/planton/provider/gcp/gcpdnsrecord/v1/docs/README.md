# GcpDnsRecord: Design Notes

## What This Component Models

`GcpDnsRecord` is a 1:1 model of Cloud DNS's record set — the unit the API
manages: one (name, type) pair inside a managed zone. The zone itself is a
separate component (`GcpDnsZone`), because zones and records have different
lifecycles and owners: a platform team typically owns the zone, while
application teams own the records their services need. Cloud DNS stores one
record set per (name, type) pair, so exactly one resource should own each
pair.

## The Two Answering Modes

A record set answers queries in one of two mutually exclusive ways, and the
API enforces the exclusivity — the spec mirrors it with a pre-deploy
validation rule:

1. **Static values (`values`)** — the classic mode. Multiple values answer
   as a round-robin set. This covers the overwhelming majority of records:
   A/AAAA service endpoints, CNAME aliases, MX mail routing, TXT
   verification and policy records.

2. **Routing policy (`routingPolicy`)** — Cloud DNS steers each query.
   Exactly one style per record:
   - **Weighted round robin (`wrr`)** — traffic splits by weight ratio.
     A zero-weight entry is valid and receives no traffic, which makes
     staged rollouts expressible: add the new target at weight 0, then
     shift weight gradually.
   - **Geolocation (`geo`)** — each entry answers queries originating
     nearest its GCP location. With `enableGeoFencing`, a location with
     unhealthy targets keeps answering rather than spilling to the
     next-closest location — useful when data residency matters more than
     availability.
   - **Primary/backup (`primaryBackup`)** — queries are answered from the
     global primary targets while any are healthy, then fall back to a
     regional geo policy. `trickleRatio` sends a slice of traffic to the
     backups even while primaries are healthy, keeping the failover path
     warm and observable.

## Health-Checked Targets

Routing policy entries can reference targets whose health determines
whether they are answered (A/AAAA records only):

- **Internal load balancer frontends** carry their own implicit health
  signal — Cloud DNS reads the load balancer's health directly, so no
  health check resource is needed. The frontend IP is referenceable to a
  `GcpAddress` resource and the network to a `GcpVpcNetwork` resource,
  keeping the private-DNS failover graph fully composable.
- **External endpoints** (public IPs) require the routing policy's
  `healthCheck` — referenceable to a `GcpHealthCheck` resource.

One provider-documented sharp edge, taught in the spec comments: if the
zone has DNSSEC enabled, a WRR entry may set only one of `values` or
`healthCheckedTargets`, not both.

## Record Type as a Free String

`type` is a validated free string (uppercase alphanumeric), not an enum.
Cloud DNS accepts a growing set of record types — HTTPS and SVCB records
are now standard practice for HTTP/3 discovery, and DS/DNSKEY appear in
DNSSEC delegations. An enum would go stale silently and reject valid
records; the API itself is the authority on which types exist.

## TTL Semantics

`ttlSeconds` defaults to 300 and accepts any non-negative value, including
an explicit 0 (no resolver caching) and the 172800 (2-day) convention for
NS records. The value is a caching hint to resolvers: low TTLs make
failover and canary shifts take effect quickly at the cost of query
volume.

## Scope Boundaries

- **`deletion_policy`** is not modeled — the attribute is not present in
  the released provider line the modules pin (`google ~> 6.0`).
- Record sets carry no labels in GCP; platform attribution lives on the
  zone.
