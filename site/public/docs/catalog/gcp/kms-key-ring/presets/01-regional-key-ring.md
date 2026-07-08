---
title: "Regional Key Ring"
description: "This preset creates a KMS key ring in a specific GCP region. It is the most common configuration — co-locating encryption keys with the workloads they protect for lowest latency and data residency..."
type: "preset"
rank: "01"
presetSlug: "01-regional-key-ring"
componentSlug: "kms-key-ring"
componentTitle: "KMS Key Ring"
provider: "gcp"
icon: "package"
order: 1
---

# Regional Key Ring

This preset creates a KMS key ring in a specific GCP region. It is the most common configuration — co-locating encryption keys with the workloads they protect for lowest latency and data residency compliance.

## When to Use

- Production workloads that require encryption keys in the same region as data
- GDPR, HIPAA, or other regulatory requirements mandating data residency
- BigQuery datasets, Cloud SQL instances, or Spanner databases that need CMEK in a specific region
- Standard key management setup for any regional GCP deployment

## Key Configuration Choices

- **Regional location** (`us-central1`) — change to match your workload region. Keys are stored exclusively in this region.
- **No multi-region replication** — keys exist in one region only. Use the multi-region preset if you need continental availability.

## Values to Adjust

- `keyRingName` — the permanent GCP name for this ring (1-63 chars,
  letters/digits/hyphens/underscores). The sample `prod-encryption` works
  as written; pick a name you will never need to recycle.
- `location` — your workload region. Multi-digit regions (e.g.
  `europe-west12`) are valid.
- `projectId` — omitted here, so the ring lands in the provider
  connection's default project. Add it (literal or a `GcpProject`
  reference) to target another project.

## Important

Key rings **cannot be deleted** from GCP. The name you choose is permanent within the project and location.

## Related Presets

- **02-global-key-ring** — Key ring accessible from all regions (no data residency)
- **03-multi-region-key-ring** — Key ring replicated across a continent (high availability + data residency)
