# Data Volume

The default posture for stateful VMs: an empty pd-balanced data disk that
exists as its own resource, so the data outlives whatever instance mounts
it.

## What this preset creates

An empty 100 GB `pd-balanced` disk in a single zone. No source is set, so
the disk is born blank — format and mount it from the instance. The disk
name falls back to `metadata.name` (`app-data`).

## Prerequisites

None beyond project credentials — the Compute Engine API is enabled
automatically.

## Composing the attachment

Attach it to a `GcpComputeInstance` in the same zone via
`attachedDisks[].source`:

```yaml
attachedDisks:
  - source:
      valueFrom:
        kind: GcpComputeDisk
        name: app-data
        fieldPath: status.outputs.self_link
```

Replacing the instance (machine-type change, image upgrade) leaves this
disk — and its data — untouched.

## Remix ideas

- Bump `sizeGb` any time — it grows in place and never shrinks.
- Switch to `pd-ssd` for latency-sensitive databases (type is
  create-time, so decide before data lands).
- Add `createSnapshotBeforeDestroy: true` once the volume holds anything
  you would miss.
