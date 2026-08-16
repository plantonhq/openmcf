# AWS EventBridge Pipe

Point-to-point event plumbing without glue code: a pipe reads from one source, optionally filters and enriches each event, and delivers to one target — the managed replacement for the "Lambda that moves messages from A to B".

## What Gets Managed

- The pipe: its source (SQS, Kinesis, DynamoDB streams, MSK, self-managed Kafka, ActiveMQ, RabbitMQ) with per-family tuning (batching, start positions, credentials, dead-letter queues), the event filter, an optional enrichment step (Lambda, Step Functions express, API destination), the target (ECS, Batch, Lambda, Step Functions, Kinesis, SQS, Redshift, SageMaker, CloudWatch Logs, event buses, HTTP APIs) with per-family invocation shaping, input templates, encryption, and execution logging.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with EventBridge Pipes permissions and `iam:PassRole` on the execution role.

### AWS Prerequisites

- The source and target resources, and an execution role trusting `pipes.amazonaws.com` that can read the source and write the target.
- For Kafka/MQ sources: broker credentials stored in Secrets Manager (the spec takes the secret's ARN, never the credential).

## After You Deploy

- The pipe starts consuming immediately (RUNNING is the default). Billing is per event processed — an idle pipe on a quiet source costs nothing.
- Creates and state changes ride a provisioning state machine; healthy creates land in a minute or two.

## Common Changes

- Pause/resume: flip `desired_state` — stream positions survive a STOPPED spell.
- Retarget in place: `target` swaps without replacing the pipe; the source never does (changing it resets consumer positions by replacement).
- Tighten the firehose: add `filter_criteria` patterns — only matching events bill.
