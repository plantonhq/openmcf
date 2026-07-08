---
title: "Private PSC Endpoint"
description: "Prediction serving exposed through Private Service Connect: per-project access control and IAM-authorized connections, without VPC peering."
type: "preset"
rank: "03"
presetSlug: "03-private-psc"
componentSlug: "vertex-ai-endpoint"
componentTitle: "Vertex AI Endpoint"
provider: "gcp"
icon: "package"
order: 3
---

# Private PSC Endpoint

Prediction serving exposed through Private Service Connect: per-project
access control and IAM-authorized connections, without VPC peering.

## What this preset creates

An endpoint named `Partner Inference` in `us-central1`, reachable only
through PSC forwarding rules created from the two allowlisted consumer
projects. Secure PSC adds IAM authorization on top of network
reachability, and the endpoint is CMEK-encrypted under the referenced
`GcpKmsKey` (`inference-key`).

## When to use

- Serving predictions to specific consumer projects (internal platform
  teams or external partners) without sharing a network
- The strongest isolation posture Vertex AI serving offers
- Multi-tenant architectures where each consumer connects from its own
  VPC

## Constraints

- PSC is mutually exclusive with both VPC peering (`network`) and the
  dedicated DNS (`dedicatedEndpointEnabled`) — the spec enforces both
  pre-deploy.

## Remix ideas

- Drop `enableSecurePrivateServiceConnect` when network-level allowlist
  control is sufficient and minimum latency matters.
- Leave `projectAllowlist` empty to allow any project in the same
  organization to connect.
