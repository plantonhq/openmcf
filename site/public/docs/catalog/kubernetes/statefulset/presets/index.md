---
title: "Presets"
description: "Ready-to-deploy configuration presets for StatefulSet"
type: "preset-list"
componentSlug: "statefulset"
componentTitle: "StatefulSet"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-database"
    rank: "01"
    title: "Single-Instance Database"
    excerpt: "This preset deploys a single-replica database: one pod with a stable name (`my-database-0`), one PersistentVolumeClaim stamped from the `data` template that survives pod restarts and rescheduling,..."
  - slug: "02-ha-quorum-cluster"
    rank: "02"
    title: "Highly Available Quorum Cluster"
    excerpt: "This preset deploys a three-member clustered stateful system — the shape of Kafka, etcd, ZooKeeper, or any consensus-based store. Each member gets a stable name (`my-quorum-cluster-0/-1/-2`), a..."
  - slug: "03-hardened-database"
    rank: "03"
    title: "Hardened Database"
    excerpt: "This preset passes the Kubernetes restricted Pod Security Standard while running persistent storage: non-root with a pinned UID, read-only root filesystem (with a writable EmptyDir for /tmp), all..."
---

# StatefulSet Presets

Ready-to-deploy configuration presets for StatefulSet. Each preset is a complete manifest you can copy, customize, and deploy.
