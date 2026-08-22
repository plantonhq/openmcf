# Same-Region Read Replica

This preset creates a read replica beside its primary -- same region, matching size -- to take read traffic off the primary. Point read-heavy consumers (reports, search indexers, read APIs) at the replica's own endpoint outputs.

## When to Use

- Read scaling: the primary saturates on reads before writes
- Isolating heavy analytical reads from transactional latency
- A warm copy in-region (DigitalOcean's console can promote a replica manually)

## Key Configuration Choices

- **Region matches the primary** -- region is required by design (never inherited silently); write the primary's region for a local replica.
- **Size matches the primary** -- DigitalOcean requires replica size >= primary size, so the primary's slug is the natural floor. Size can grow in place later; it can never shrink.
- **No tags yet** -- replica tags are create-only upstream: a retag REPLACES the replica. Settle tagging before production.

## What You Get

A single-node read-only endpoint (`host`/`port`/`uri` outputs, credentials included as secrets) that follows the primary's data -- billed like a second cluster node of the same slug.
