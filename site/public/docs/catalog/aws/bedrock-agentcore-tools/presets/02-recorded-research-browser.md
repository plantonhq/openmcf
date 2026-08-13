---
title: "Recorded Research Browser"
description: "This preset gives an agent a cloud browser whose every session is recorded to S3 and whose traffic is cryptographically signed — plus a saved profile so sessions start already logged in to your docs..."
type: "preset"
rank: "02"
presetSlug: "02-recorded-research-browser"
componentSlug: "bedrock-agentcore-tools"
componentTitle: "Bedrock AgentCore Tools"
provider: "aws"
icon: "package"
order: 2
---

# Recorded Research Browser

This preset gives an agent a cloud browser whose every session is
recorded to S3 and whose traffic is cryptographically signed — plus a
saved profile so sessions start already logged in to your docs portal.

## When to Use

- Agents that browse authenticated or third-party sites
- Compliance regimes that ask "what did the agent actually do?"

## What You Get

- Session replays in your bucket under `browser-sessions/` — the audit
  trail for every page the agent touched
- Signed traffic so sites can verify an AWS-managed browser
- A reusable logged-in profile (treat it like the credential it holds)

## Customize

- Add `enterprisePolicies` (Chrome policy JSON in S3, type MANAGED) to
  allow-list URLs or kill downloads
- Switch to `mode: VPC` with subnets/security groups for intranet
  research
- Lifecycle the recordings prefix — replays accumulate fast
