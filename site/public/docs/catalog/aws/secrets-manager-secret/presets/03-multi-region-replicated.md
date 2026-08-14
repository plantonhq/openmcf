---
title: "Multi-Region Replicated"
description: "This preset creates a secret replicated to additional regions with a resource policy restricting reads to a named application role."
type: "preset"
rank: "03"
presetSlug: "03-multi-region-replicated"
componentSlug: "secrets-manager-secret"
componentTitle: "Secrets Manager Secret"
provider: "aws"
icon: "package"
order: 3
---

# Multi-Region Replicated

This preset creates a secret replicated to additional regions with a resource policy restricting reads to a named application role.

## When to Use

- Multi-region applications that read the same credential everywhere (regional reads, no cross-region latency or dependency)
- Disaster-recovery postures where the standby region must work when the primary is down
- Secrets consumed by regional services (Lambda@Edge origins, regional ECS services)

## Key Configuration Choices

- **Replicas are read-only copies** kept in sync by Secrets Manager — write to the primary, read anywhere by the same name (each region has its own ARN)
- **Per-replica KMS**: replicas here use each region's AWS-managed key; add `kmsKeyId` per replica for customer-managed keys — the key must live in the replica's region
- **Resource policy + `blockPublicPolicy`** — least-privilege reads with the anonymous-access guard on

## Operational Notes

Removing a region from `replicaRegions` deletes that replica. If a same-named secret already exists in a replica region, replication fails loudly unless `forceOverwriteReplicaSecret: true` is set deliberately.
