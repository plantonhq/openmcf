---
title: "gRPC Health Service Probe"
description: "A health check calling the standard gRPC health protocol (`grpc.health.v1.Health/Check`) — the native probe for gRPC microservices behind a load balancer. The backend passes only while it reports..."
type: "preset"
rank: "03"
presetSlug: "03-grpc-service"
componentSlug: "health-check-on-google-cloud"
componentTitle: "Health Check on Google Cloud"
provider: "gcp"
icon: "package"
order: 3
---

# gRPC Health Service Probe

A health check calling the standard gRPC health protocol (`grpc.health.v1.Health/Check`) — the native probe for gRPC microservices behind a load balancer. The backend passes only while it reports `SERVING`.

## When to Use

- gRPC backends behind global external or internal load balancers
- Multi-service gRPC servers where one service's health should gate traffic (set `grpcServiceName`)

## Remix Notes

- The backend MUST implement the gRPC health service, and if `grpcServiceName` is set it must recognize exactly that string — an unimplemented health service shows as permanently unhealthy. `enableLogging` (on in this preset) is the fastest way to spot that wiring mistake.
- Remove `grpcServiceName` to probe the server's overall health instead of one service.
- For TLS-only gRPC servers, switch the block to `grpcTls` (note: named ports are not supported there).
- Health-check logs carry per-transition detail; turn logging off once thresholds are settled if log volume matters.
