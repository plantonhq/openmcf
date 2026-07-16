---
title: "KMS Key Ring"
description: "KMS Key Ring deployment documentation"
icon: "package"
order: 100
componentName: "gcpkmskeyring"
---

# GCP KMS Key Ring

Creates a Cloud KMS key ring — the permanent, location-anchored container that groups cryptographic keys. IAM granted on the ring flows down to every key inside it, making the ring the blast-radius boundary for customer-managed encryption.

## What Gets Created

- The Cloud KMS API is enabled on the target project (never disabled on destroy)
- A `google_kms_key_ring` resource in the specified project and location

The key ring API has no labels surface, so no attribution labels are stamped — identically on both engines.

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **A GCP project** (the Cloud KMS API is enabled automatically)
- **IAM permissions** — `roles/cloudkms.admin` or `roles/cloudkms.keyRingCreator`

## Quick Start

Create a file `key-ring.yaml`:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpKmsKeyRing
metadata:
  name: prod-encryption
spec:
  keyRingName: prod-encryption
  location: us-central1
```

Deploy:

```shell
planton apply -f key-ring.yaml
```

With `projectId` omitted, the ring lands in the provider connection's default project.

## Configuration Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `projectId` | `StringValueOrRef` | No | GCP project for the ring. Omitted rides the provider default; reference a `GcpProject` to compose. Immutable. |
| `keyRingName` | `string` | Yes | The permanent GCP name (1-63 chars: letters, digits, hyphens, underscores). Immutable — and never reusable within the project+location, because rings cannot be deleted. |
| `location` | `string` | Yes | Region (`us-central1`), multi-region (`us`, `europe`, `asia`), or `global`. Immutable. Keys inherit it, and most CMEK integrations require co-location with the protected data. |

## Permanence

**Key rings cannot be deleted from GCP.** Destroying this resource removes the ring from IaC state only; the (free, inert) ring remains in the project forever and its name is permanently consumed. Every field is ForceNew — a change abandons the old ring and creates a new one.

## Design Guidance

- **One ring per environment or data domain**, not one per key — ring-level IAM is the unit of access review.
- **Location follows data**: regional rings for regional CMEK, `us`/`europe`/`asia` for multi-region BigQuery/Spanner, `global` for signing keys consumed everywhere.

## Stack Outputs

| Output | Description |
|--------|-------------|
| `key_ring_id` | Fully qualified ring path (`projects/{p}/locations/{l}/keyRings/{name}`) — what a `GcpKmsKey`'s `keyRingId` reference consumes |
| `key_ring_name` | The short name of the ring |
| `location` | The ring's location, for consumers that take a bare name + location pair |

## Related Components

- [GcpKmsKey](../kms-key/) — the cryptographic keys grouped by this ring
- [GcpProject](../project/) — provides the GCP project by reference
