---
title: "Presets"
description: "Ready-to-deploy configuration presets for Deployment"
type: "preset-list"
componentSlug: "deployment"
componentTitle: "Deployment"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-web-service"
    rank: "01"
    title: "Web Service Deployment"
    excerpt: "This preset deploys a single-replica web application fronted by a ClusterIP Service, with a readiness probe so traffic only reaches pods that are ready. It is the most common Kubernetes Deployment..."
  - slug: "02-web-service-with-hpa"
    rank: "02"
    title: "Web Service with Autoscaling and Zero-Downtime Rollouts"
    excerpt: "This preset is the production baseline for a stateless web service: CPU-based horizontal autoscaling between 2 and 10 replicas, a rollout strategy that never drops below the desired replica count, a..."
  - slug: "03-worker"
    rank: "03"
    title: "Background Worker"
    excerpt: "This preset deploys a queue consumer or background processor: no ports, no Service — just the container, its environment, and a secret pulled from an existing Kubernetes Secret. When the app..."
  - slug: "04-hardened-production"
    rank: "04"
    title: "Hardened Production Service"
    excerpt: "This preset passes the Kubernetes restricted Pod Security Standard: non-root with a pinned UID, read-only root filesystem (with a writable EmptyDir for /tmp), all Linux capabilities dropped, the..."
---

# Deployment Presets

Ready-to-deploy configuration presets for Deployment. Each preset is a complete manifest you can copy, customize, and deploy.
