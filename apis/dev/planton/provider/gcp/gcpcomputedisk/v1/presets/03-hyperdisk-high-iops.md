# Hyperdisk High IOPS

Storage for workloads whose performance needs will change: a
hyperdisk-balanced volume where IOPS and throughput are dials, not
consequences of size.

## What this preset creates

A 1000 GB `hyperdisk-balanced` disk provisioned at 6000 IOPS and
290 MB/s throughput. Unlike the pd family — where performance is a
function of disk size — hyperdisk buys capacity and performance
independently, and both `provisionedIops` and `provisionedThroughput`
update **in place** (at most every 4 hours).

## Prerequisites

- The instance that attaches this disk must run a machine family that
  supports hyperdisk (e.g. C3, N4, M3 series). The pd family works
  everywhere; hyperdisk does not.

## Composing the attachment

Same pattern as any data disk — the instance references the disk's
`self_link` output:

```yaml
attachedDisks:
  - source:
      valueFrom:
        kind: GcpComputeDisk
        name: analytics-data
        fieldPath: status.outputs.self_link
```

## Remix ideas

- Retune `provisionedIops` / `provisionedThroughput` as monitoring
  dictates — the change applies without replacing the disk (respect the
  4-hour spacing).
- `hyperdisk-throughput` for streaming/scan-heavy workloads where MB/s
  matters more than IOPS; `hyperdisk-extreme` for the highest IOPS
  ceilings.
- `hyperdisk-ml` with `accessMode: READ_ONLY_MANY` serves one shared
  dataset to many VMs at once.
- Add `kmsKey` and `createSnapshotBeforeDestroy: true` when the data
  becomes precious.
