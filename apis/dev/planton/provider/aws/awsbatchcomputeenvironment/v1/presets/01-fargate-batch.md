# Fargate Batch

A serverless Fargate compute environment — zero instances to manage, per-second billing, scale-to-zero when queues are empty. The right default for containerized batch jobs that fit Fargate's sizing (up to 16 vCPU / 120 GiB per job).

## When to Use

- Bursty or unpredictable batch workloads where idle instances would waste money
- Teams that want no AMI, patching, or capacity management at all
- Jobs that need no GPUs, no privileged mode, and no custom kernel settings

## What It Configures

- **FARGATE type** — AWS manages all compute; only `maxVcpus`, subnets, and security groups apply
- **256 vCPU ceiling** — total concurrent capacity across all running jobs
- **Two private subnets** — task ENIs spread across Availability Zones
- **Security group** — required for Fargate task ENIs

## What to Customize

- Replace `<aws-region>` and the subnet/security-group placeholders (or use `valueFrom` references to `AwsSubnet` / `AwsSecurityGroup` resources)
- Raise `maxVcpus` for higher concurrency; switch `type` to `FARGATE_SPOT` for interruptible jobs at Spot pricing
- Add an `AwsBatchJobQueue` mapped onto this environment to start submitting jobs
