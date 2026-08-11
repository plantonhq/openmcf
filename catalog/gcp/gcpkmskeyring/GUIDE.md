# GcpKmsKeyRing Guide

The judgment this guide protects: a key ring is the one genuinely
permanent object in GCP — it cannot be deleted, ever, by anyone. Every
choice here is written in ink.

## Rings are forever; name and place them like it

A key ring has no delete API: destroy only removes it from management,
and the name is occupied in that project+location permanently. Choose
ring names by durable domain (`payments-prod`, `data-platform`), never
by team names or project codenames that reorganize away. The location
is equally immutable and must match where the encrypting services run —
a `us-central1` ring cannot serve `europe-west1` CMEK requirements.

## The ring is an IAM boundary — use few, deliberately

IAM granted on the ring flows to every key in it, which makes the ring
the natural unit of custody: one ring per environment per domain, keys
inside it per purpose. Many rings with one key each forfeits the
boundary's usefulness; one giant ring makes every grant a fleet-wide
grant.

## Destroy semantics

Because the underlying resource cannot be deleted, this kind's destroy
is always effectively an abandon: the ring (and its keys, which are
their own kind) live on. The real destroy decisions live on GcpKmsKey,
where versions — and therefore data — are at stake.
