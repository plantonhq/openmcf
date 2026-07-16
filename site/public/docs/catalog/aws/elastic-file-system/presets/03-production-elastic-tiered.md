---
title: "Production Elastic EFS with Lifecycle Tiering and DR Replication"
description: "Regional, encrypted, elastic throughput, full lifecycle tiering (IA + Archive + warm-on-access), daily backups, and a cross-region replica for disaster recovery."
type: "preset"
rank: "03"
presetSlug: "03-production-elastic-tiered"
componentSlug: "elastic-file-system"
componentTitle: "Elastic File System"
provider: "aws"
icon: "package"
order: 3
---

# Production Elastic EFS with Lifecycle Tiering and DR Replication

Regional, encrypted, elastic throughput, full lifecycle tiering (IA + Archive + warm-on-access), daily backups, and a cross-region replica for disaster recovery.

## When to Use

- Production workloads with unpredictable or spiky I/O (elastic throughput scales automatically)
- Mixed hot/cold data where lifecycle tiering meaningfully cuts storage cost
- Workloads with a DR requirement — the replica stays read-only and in sync in another region
- Shared storage for ECS tasks and Lambda functions (add `AwsEfsAccessPoint` resources per application)

## What It Configures

- **Elastic throughput** — Scales up/down with the workload; no burst-credit management
- **Lifecycle tiering** — Files idle 30 days move to IA (~92% cheaper); IA files idle 90 more days move to Archive (~96% cheaper); accessed files warm back to Standard
- **Backup enabled** — Daily backups via AWS Backup
- **Two mount targets** — One per AZ so clients mount locally and avoid cross-AZ data charges
- **Cross-region replication** — A read-only replica in us-west-2, kept in sync automatically

## What to Customize

- Replace placeholders: `<subnet-id-az-a>`, `<subnet-id-az-b>`, `<security-group-id>`
- Change `replication.destinationRegion`, or use `destinationAvailabilityZoneName` for a cheaper One Zone replica
- Tune the lifecycle windows (`AFTER_7_DAYS` … `AFTER_365_DAYS`) to your access patterns
- Add a `policy` enforcing encryption in transit (deny `aws:SecureTransport: false`)
- Create `AwsEfsAccessPoint` resources referencing this file system for per-application POSIX isolation
