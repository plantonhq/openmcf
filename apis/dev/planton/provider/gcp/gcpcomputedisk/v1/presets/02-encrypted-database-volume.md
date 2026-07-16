# Encrypted Database Volume

The posture for regulated data: a high-IOPS SSD volume under
customer-managed encryption, with a destroy-time snapshot so even a
mistaken teardown cannot lose the data.

## What this preset creates

A 500 GB `pd-ssd` disk encrypted with a customer-managed key (CMEK) from
a `GcpKmsKey` resource, and `createSnapshotBeforeDestroy: true` — the
final snapshot reuses the disk's key.

## Prerequisites

- A `GcpKmsKey` named `db-disk-key` (replace the reference with your
  own).
- The Compute Engine service agent
  (`service-<project-number>@compute-system.iam.gserviceaccount.com`)
  must hold `roles/cloudkms.cryptoKeyEncrypterDecrypter` on that key
  before the disk is created — this is the prerequisite that trips
  everyone once. Note it is the service agent, not your deploy
  credential.

## Composing the attachment

A CMEK disk's attachment must present the same key. On the
`GcpComputeInstance` side:

```yaml
attachedDisks:
  - source:
      valueFrom:
        kind: GcpComputeDisk
        name: db-data
        fieldPath: status.outputs.self_link
    kmsKey:
      valueFrom:
        kind: GcpKmsKey
        name: db-disk-key
        fieldPath: status.outputs.key_id
```

## Remix ideas

- Set `snapshotBeforeDestroyPrefix` to control the final snapshot's name
  (default `<disk-name>-YYYYMMDD-HHmmss`).
- The key choice is immutable — moving to a different key means a new
  disk restored from a snapshot.
- For tunable performance instead of fixed pd-ssd scaling, see the
  hyperdisk preset.
