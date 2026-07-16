---
title: "Fargate Container Job"
description: "A Fargate job definition with the full production posture: split identities, parameterized command, Spot-safe retry discrimination, and a hard timeout."
type: "preset"
rank: "01"
presetSlug: "01-fargate-container-job"
componentSlug: "batch-job-definition"
componentTitle: "Batch Job Definition"
provider: "aws"
icon: "package"
order: 1
---

# Fargate Container Job

A Fargate job definition with the full production posture: split identities, parameterized command, Spot-safe retry discrimination, and a hard timeout.

## When to Use

- Containerized batch steps on Fargate / Fargate Spot compute environments
- Parameterized pipelines where one definition serves many datasets

## What It Configures

- **Fargate sizing** — 1 vCPU / 2048 MiB from the Fargate size table
- **Split identities** — the execution role pulls the ECR image and writes logs; the job role is the code's own AWS identity
- **`Ref::dataset` placeholder** — defaulted by `parameters`, overridable per SubmitJob call
- **Retry discrimination** — `Host EC2*` status reasons (Spot reclaims) RETRY; `1*` exit codes (application failures) EXIT immediately
- **Two-hour timeout** — runaway attempts are terminated and retried per the strategy

## What to Customize

- Replace the region/image/role placeholders (roles can be `valueFrom` references to `AwsIamRole` resources)
- Adjust `vcpus`/`memoryMib` to a valid Fargate pairing; add `runtimePlatform.cpuArchitecture: ARM64` for Graviton
- Add `secrets` for credentials (Secrets Manager ARNs — never plain environment values)
