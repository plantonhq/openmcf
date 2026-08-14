---
title: "Full-Fidelity Audit"
description: "This preset delivers invocation logs to BOTH destinations: CloudWatch for querying (with S3 spillover for oversized payloads) and S3 for retention — the canonical audit posture for production Bedrock..."
type: "preset"
rank: "01"
presetSlug: "01-full-fidelity-audit"
componentSlug: "bedrock-invocation-logging"
componentTitle: "Bedrock Invocation Logging"
provider: "aws"
icon: "package"
order: 1
---

# Full-Fidelity Audit

This preset delivers invocation logs to BOTH destinations: CloudWatch
for querying (with S3 spillover for oversized payloads) and S3 for
retention — the canonical audit posture for production Bedrock
workloads.

## When to Use

- Production regions running Bedrock workloads that need a complete
  audit trail
- Compliance postures where model inputs/outputs must be retained

## What You Get

- Every model call logged: queryable in CloudWatch, archived in S3
- Payloads over the 256 KB CloudWatch event cap preserved via the
  spillover prefix instead of truncated

## Customize

- Point the role at one trusting `bedrock.amazonaws.com` with write
  access to the log group
- The bucket policy must let `bedrock.amazonaws.com` put objects —
  Bedrock writes S3 as its service principal, not through the role
- Drop `videoDataDeliveryEnabled` / `embeddingDataDeliveryEnabled` to
  `false` to trim volume for text-first workloads
