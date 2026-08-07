# Time-Series Table

The classic Bigtable shape: a wide measurements family with age-based
retention, a small metadata family capped by versions, and key-prefix
pre-splits so initial load distributes across tablets.

## When to use

IoT readings, metrics, clickstreams — any append-heavy workload keyed by
entity + timestamp where old data expires rather than being deleted by
the application.

## What to customize

- `columnFamilies[].gcPolicy.maxAge` — the retention window; this is the
  lever that controls storage cost.
- `splitKeys` — match your row-key prefixes; set at creation (changing
  them later replaces the table and its data).

## Composes with

`GcpBigtableInstance` upstream (reference its `instance_name` output).
