---
title: "Rotated Secret"
description: "A credential with a rotation pipeline behind it: GCP publishes a reminder to Pub/Sub every 30 days, the subscriber rotates and adds a new version, and the `current` alias re-points consumers..."
type: "preset"
rank: "03"
presetSlug: "03-rotated-secret"
componentSlug: "secret-manager-secret"
componentTitle: "Secret Manager Secret"
provider: "gcp"
icon: "package"
order: 3
---

# Rotated Secret

A credential with a rotation pipeline behind it: GCP publishes a
reminder to Pub/Sub every 30 days, the subscriber rotates and adds a
new version, and the `current` alias re-points consumers atomically.

## What it configures

- `rotation` (30-day period) + `topics` — GCP's rotation feature
  publishes REMINDERS; it rotates nothing itself. The Pub/Sub
  subscriber (a Cloud Function, a pipeline) performs the rotation.
- `versionDestroyTtl: 604800s` — destroyed versions get a 7-day
  disable-first restore window, the undo buffer for a rotation gone
  wrong.
- `versionAliases.current` — consumers address `versions/current` and
  never notice rotations.

## Adjust before deploying

- **topics** — reference a GcpPubSubTopic via valueFrom (its `topic_id`
  output); the Secret Manager service agent needs
  `roles/pubsub.publisher` on it.
- **nextRotationTime** — the first reminder's timestamp; GCP advances
  it by the period thereafter.
- Build the subscriber — rotation without one is a calendar, not a
  rotation.

## When to choose something else

Secrets that never rotate (third-party keys with no API) skip the
rotation block and rely on `versionDestroyTtl` alone for safety.
