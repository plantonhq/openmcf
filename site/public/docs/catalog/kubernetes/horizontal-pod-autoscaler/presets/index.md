---
title: "Presets"
description: "Ready-to-deploy configuration presets for Horizontal Pod Autoscaler"
type: "preset-list"
componentSlug: "horizontal-pod-autoscaler"
componentTitle: "Horizontal Pod Autoscaler"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-cpu-autoscale"
    rank: "01"
    title: "CPU Autoscale"
    excerpt: "This preset is the standard workhorse autoscaler: hold the workload's average CPU at 60% of its requests, running between 2 and 10 replicas. CPU is the reliable scaling signal — it rises and falls..."
  - slug: "02-container-isolated"
    rank: "02"
    title: "Container Isolated"
    excerpt: "This preset scales on ONE container's CPU instead of the pod-level average. A plain `resource` metric blends every container in the pod — the app, the service-mesh proxy, the log shipper — into one..."
  - slug: "03-queue-driven"
    rank: "03"
    title: "Queue Driven"
    excerpt: "This preset sizes a worker fleet by its backlog instead of its CPU: roughly one pod per 30 ready messages in one queue. Queue depth is the honest signal for workers — a worker grinding through a..."
  - slug: "04-behavior-tuned"
    rank: "04"
    title: "Behavior Tuned"
    excerpt: "This preset keeps the fast scale-UP defaults and makes scale-DOWN deliberately slow: wait 10 minutes of consistently lower recommendations, then remove at most 10% of the fleet per minute. The..."
---

# Horizontal Pod Autoscaler Presets

Ready-to-deploy configuration presets for Horizontal Pod Autoscaler. Each preset is a complete manifest you can copy, customize, and deploy.
