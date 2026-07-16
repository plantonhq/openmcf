---
title: "Presets"
description: "Ready-to-deploy configuration presets for Health Check"
type: "preset-list"
componentSlug: "health-check"
componentTitle: "Health Check"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-http-serving-port"
    rank: "01"
    title: "HTTP Probe on the Serving Port"
    excerpt: "The workhorse health check for global external HTTP(S) load balancers: an HTTP GET against a dedicated health endpoint, probing whatever port each backend actually serves on (`USE_SERVING_PORT`) so..."
  - slug: "02-regional-tcp"
    rank: "02"
    title: "Regional TCP Probe"
    excerpt: "A regional health check proving TCP connectability on a fixed port — the shape internal load balancers and regional managed instance groups consume. Regional backend services can only reference..."
  - slug: "03-grpc-service"
    rank: "03"
    title: "gRPC Health Service Probe"
    excerpt: "A health check calling the standard gRPC health protocol (`grpc.health.v1.Health/Check`) — the native probe for gRPC microservices behind a load balancer. The backend passes only while it reports..."
---

# Health Check Presets

Ready-to-deploy configuration presets for Health Check. Each preset is a complete manifest you can copy, customize, and deploy.
