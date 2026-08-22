# Shared Inbound-User Session Pool

This preset creates DigitalOcean's "inbound user" pool: the `user` field is omitted, so every client authenticates with its OWN database credentials and the pool proxies them. Session mode holds a server connection for the client's whole session.

## When to Use

- Shared pools serving several teams or tools, each with their own users
- Workloads needing session state: LISTEN/NOTIFY, session-scoped prepared statements, advisory locks
- Avoiding a single shared pool identity (every client keeps its own grants and audit trail)

## Key Configuration Choices

- **`user` omitted** -- the safer default for shared pools; the pool's `password` output is legitimately empty (clients bring credentials).
- **`mode: session`** -- one server connection per client session; size the pool closer to the expected concurrent CLIENT count than in transaction mode.
- **`dbName: defaultdb`** -- the cluster's built-in database; point it at a DigitalOceanDatabaseDb name for a workload-specific shared pool.

## What You Get

A pool endpoint any existing cluster user can authenticate through, with per-client identity preserved end to end.
