# AWS EventBridge Scheduler

Serverless cron: fire a Lambda, queue a message, run an ECS task, or start a pipeline on a cron, rate, or one-time expression — with timezones, retries, and a dead-letter queue, and none of the instance babysitting of a cron box.

## What Gets Managed

- The schedule: its expression and timezone, start/end window, enabled/disabled state, flexible time window, what happens after a one-time run completes, encryption, and the target — ARN + execution role, static input payload, retry policy, dead-letter queue, and service-specific shaping for ECS, EventBridge, Kinesis, SageMaker, and SQS.
- Optionally its schedule group (created here or joined by name; the group carries the tags — schedules are untaggable at AWS).

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with EventBridge Scheduler permissions and `iam:PassRole` on the execution role.

### AWS Prerequisites

- The target resource and an execution role trusting `scheduler.amazonaws.com` that can invoke it.

## After You Deploy

- The schedule starts firing on its expression immediately (state ENABLED is AWS's default). The free tier covers fourteen million invocations a month — most fleets pay nothing.
- First deploys with a freshly created role ride a two-minute IAM-propagation retry — expected, not a hang.

## Common Changes

- Pause/resume: flip `state` to DISABLED/ENABLED — everything survives.
- Smooth thundering herds: switch the flexible window to FLEXIBLE with a spread (in place).
- Re-group: the group is fixed for life — moving a schedule replaces it (positions and history are not a concern; schedules are stateless).
