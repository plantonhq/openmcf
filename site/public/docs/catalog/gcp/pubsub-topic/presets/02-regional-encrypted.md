---
title: "Regional Encrypted Topic"
description: "The data-residency posture: customer-managed encryption plus hard region pinning for regulated event streams."
type: "preset"
rank: "02"
presetSlug: "02-regional-encrypted"
componentSlug: "pubsub-topic"
componentTitle: "Pub/Sub Topic"
provider: "gcp"
icon: "package"
order: 2
---

# Regional Encrypted Topic

The data-residency posture: customer-managed encryption plus hard region
pinning for regulated event streams.

## What this preset creates

A topic whose messages are encrypted with a customer-managed KMS key
(referenced from a `GcpKmsKey` resource) and persisted only in
`us-central1`. With `enforceInTransit: true`, publish calls from
non-allowed regions are rejected outright instead of being rerouted —
the strictest residency guarantee Pub/Sub offers.

## Prerequisites

- A `GcpKmsKey` named `audit-events-key` (or swap in a literal key path).
- Grant the Pub/Sub service agent
  (`service-{project_number}@gcp-sa-pubsub.iam.gserviceaccount.com`)
  `roles/cloudkms.cryptoKeyEncrypterDecrypter` on the key — publishes
  fail without it.

## Remix ideas

- Widen `allowedPersistenceRegions` to an approved region set while
  keeping the in-transit guarantee.
- Drop `enforceInTransit` to allow cross-region publishers whose
  messages are still stored in-region.
