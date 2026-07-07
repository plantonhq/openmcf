---
title: "HA Production"
description: "A high-availability Dataproc cluster for production Spark workloads: 3 masters, SSD storage with NVMe local SSDs, Shielded VMs, CMEK encryption, private networking, and OSS metrics into Cloud..."
type: "preset"
rank: "02"
presetSlug: "02-ha-production"
componentSlug: "dataproc-cluster"
componentTitle: "Dataproc Cluster"
provider: "gcp"
icon: "package"
order: 2
---

# HA Production

A high-availability Dataproc cluster for production Spark workloads: 3
masters, SSD storage with NVMe local SSDs, Shielded VMs, CMEK
encryption, private networking, and OSS metrics into Cloud Monitoring.

## When to use

- Production ETL pipelines that must not fail
- Long-running Spark Structured Streaming applications
- Mission-critical data processing with SLA requirements
- Environments requiring CMEK encryption and private networking

## What to customize

- `projectId` / `region` — your project and region.
- The `valueFrom` references — point at your own `GcpSubnetwork`,
  `GcpServiceAccount`, and `GcpKmsKey` resources (or swap in literal
  `value:` strings).
- `workerConfig.numInstances` / `minNumInstances` — both update in
  place, so capacity changes never recreate the cluster.
- Add `autoscalingPolicyUri` referencing a `GcpDataprocAutoscalingPolicy`
  to hand scaling decisions to the YARN-based autoscaler.

## Key configuration

- **3 masters** — high-availability HDFS/YARN with automatic failover
- **5 workers (floor of 3)** — `minNumInstances` is the autoscaler's
  hard floor; updatable in place
- **NVMe local SSDs** — fast shuffle and spill for wide transformations
- **Shielded VMs** — secure boot, vTPM, and integrity monitoring on
  every node
- **Internal IP only** — no public internet exposure; pair with Cloud
  NAT for package downloads
- **CMEK encryption** — customer-managed keys on every persistent disk
- **1-hour graceful decommission** — running tasks complete before a
  scale-down removes their node
- **Spark/YARN metric collection** — the observability surface
  dashboards and alerting build on

## Related presets

- **01-dev-jupyter** — lightweight cluster for development
- **03-cost-optimized-batch** — Spot secondaries and autoscaling for batch jobs
- **04-spark-on-gke** — run Spark as pods on an existing GKE cluster
