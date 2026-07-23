---
title: "Presets"
description: "Ready-to-deploy configuration presets for MongoDB"
type: "preset-list"
componentSlug: "mongodb"
componentTitle: "MongoDB"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-single-instance"
    rank: "01"
    title: "Single Instance"
    excerpt: "This preset declares the smallest useful MongoDB cluster: one replica set (`rs0`) with a single member, small storage, no backups. It is a single point of failure by construction — the development..."
  - slug: "02-replica-set"
    rank: "02"
    title: "Replica Set"
    excerpt: "This preset declares the production MongoDB posture: a three-member replica set with automated failover (a new primary is elected in seconds when the current one dies), explicit resources, a..."
---

# MongoDB Presets

Ready-to-deploy configuration presets for MongoDB. Each preset is a complete manifest you can copy, customize, and deploy.
