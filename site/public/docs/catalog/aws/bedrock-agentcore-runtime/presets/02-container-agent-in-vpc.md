---
title: "Container Agent in a VPC"
description: "This preset hosts a container-image agent whose sessions attach to your subnets — the pattern for agents that query private databases or internal APIs — with session lifecycle caps and per-session..."
type: "preset"
rank: "02"
presetSlug: "02-container-agent-in-vpc"
componentSlug: "bedrock-agentcore-runtime"
componentTitle: "Bedrock AgentCore Runtime"
provider: "aws"
icon: "package"
order: 2
---

# Container Agent in a VPC

This preset hosts a container-image agent whose sessions attach to your
subnets — the pattern for agents that query private databases or
internal APIs — with session lifecycle caps and per-session scratch
storage.

## When to Use

- Agents that must reach private resources (RDS, internal services)
- Teams that already build and scan their own images

## What You Get

- Session network interfaces in YOUR subnets behind YOUR security
  groups — internet egress only through your VPC's routing
- Idle and hard lifetime caps so abandoned sessions never linger
- Ephemeral `/mnt/scratch` per session, wiped at session end

## Customize

- Add an EFS access-point mount for durable cross-session state
- Point `imageUri` at a version tag, not `latest` — every change creates
  a new runtime version you can roll back through endpoints
