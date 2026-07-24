---
title: "Presets"
description: "Ready-to-deploy configuration presets for SeaweedFS"
type: "preset-list"
componentSlug: "seaweedfs"
componentTitle: "SeaweedFS"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-dev-single-node"
    rank: "01"
    title: "Dev single node preset"
    excerpt: "The smallest useful SeaweedFS: one master, one volume server and one filer on small PersistentVolumeClaims, the S3 gateway embedded on the filer with authentication on (the default — the chart..."
  - slug: "02-production-ha"
    rank: "02"
    title: "Production HA preset"
    excerpt: "A highly available store: a 3-master Raft quorum, 3 volume servers on explicit fast storage with the replication code that keeps one extra copy of every object on another server, a dedicated..."
  - slug: "03-artifact-store"
    rank: "03"
    title: "Artifact store preset"
    excerpt: "A bucket-centric store for build and release artifacts: three typed buckets created at install — CI artifacts that expire themselves after 30 days (SeaweedFS TTL, no cleanup job), release artifacts..."
---

# SeaweedFS Presets

Ready-to-deploy configuration presets for SeaweedFS. Each preset is a complete manifest you can copy, customize, and deploy.
