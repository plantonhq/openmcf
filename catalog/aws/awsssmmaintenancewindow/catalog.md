# AWS SSM Maintenance Window

A recurring change window for your fleet: when it opens, which
machines are in scope, and which tasks run — patch installs,
automation runbooks, Lambda functions, Step Functions executions —
with rate controls and a hard cutoff.

## What Gets Managed

- The window itself: schedule (cron/rate + timezone + offset), how
  long it stays open, and when new task starts stop (the cutoff).
- Registered targets: instances by tag or ID, or resource groups.
- Tasks of all four types with priorities, per-task rate controls,
  cutoff behavior, and type-specific invocation parameters (Run
  Command output to S3/CloudWatch, SNS notifications, runbook
  parameters, Lambda payloads, Step Functions input).

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with SSM permissions.

### AWS Account

- Managed nodes for the targets to select (the window deploys fine
  before they exist).
- For LAMBDA / STEP_FUNCTIONS tasks: the function
  ([AWS Lambda](/cloud-catalog/aws-lambda)) or state machine
  ([AWS Step Function](/cloud-catalog/aws-step-function)) to invoke.

## Deploy

### Console

Create the resource from the AWS catalog, set the schedule and
duration/cutoff, register a target, add a task, and deploy.

### CLI

```bash
planton apply -f maintenance-window.yaml
```

## After Deploy

- Window executions and per-task results appear under Maintenance
  Windows in the Systems Manager console.
- Registered target/task IDs surface as `target_ids`/`task_ids`
  outputs.
- Windows, targets, and tasks are free — what the tasks run is what
  costs.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
