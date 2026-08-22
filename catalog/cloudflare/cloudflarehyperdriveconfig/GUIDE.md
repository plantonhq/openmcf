# CloudflareHyperdriveConfig guide

Operational judgment for Hyperdrive configs. The README covers what each field is; this covers how the pieces interact.

## Creation is a live connection test

Cloudflare dials the origin database when the config is created — an unreachable host, a firewall rule, or wrong credentials fail the DEPLOY, not the first query. Order infrastructure so the database (and its network path from Cloudflare) exists before this config, and read a create failure as a connectivity report, not a module defect.

## The password is write-only

The API never returns the origin password. An imported config lands with an empty password in state, and the first plan after import shows a password diff — that re-assert is expected, not drift. Always provide the password as a managed-secret reference so the platform resolves it just-in-time.

## Caching trades staleness for origin load

Cached query results are served for up to `max_age` seconds (plus a stale-while-revalidate window). That is ideal for read-heavy dashboards and catalogs, and wrong for anything that must read its own writes — disable caching for those configs rather than tuning the windows down.

## mTLS and VPC Services are alternative trust paths

A VPC Service origin manages TLS on the VPC side, so the spec forbids combining `service_id` with the `mtls` block. Pick one: public/Access-fronted origins may add mTLS with pre-uploaded certificates; VPC-routed origins bring their own.

## Connection limits are plan-bounded

`origin_connection_limit` floors at 5 and ceilings by plan (about 20 on free, up to 100 on paid). Leaving it 0 takes the plan default — set it only when the origin database's own max_connections budget demands it.
