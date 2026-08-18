# AwsEventBridgeScheduler — Terraform/OpenTofu module

Manages one EventBridge Scheduler schedule (`aws_scheduler_schedule`) with an optionally owned group (`aws_scheduler_schedule_group`).

Module facts worth knowing before editing:

- **Name and group are fixed for life** — both replace the schedule on change; the group resolution is owned group → joined `group_name` → AWS's `default`.
- **The schedule is untaggable at AWS** — identity tags land on the owned group only (the deliberate tag-convention absence).
- **IAM propagation gates first deploys** — the provider retries "must allow AWS EventBridge Scheduler to assume the role" for up to two minutes when the execution role is freshly created.
- **`action_after_completion = DELETE` deletes out from under state** — a completed one-time schedule disappears at AWS and the next deploy recreates it; fire-and-forget only.
- **The retry-policy diff suppression is deliberately narrow at the provider** — updating `maximum_retry_attempts` between 0 and 1 works correctly.
- **The import ID is `{group_name}/{name}`** — both halves are exported.

Outputs mirror the Pulumi module key-for-key: `schedule_arn`, `group_name`, `group_arn`.
