---
title: "Durable Workflow"
description: "This preset creates a durable Lambda function: AWS checkpoints the function's progress so a long-running workflow — here, a week-long order-fulfillment saga — survives interruption and resumes where..."
type: "preset"
rank: "03"
presetSlug: "03-durable-workflow"
componentSlug: "lambda"
componentTitle: "Lambda"
provider: "aws"
icon: "package"
order: 3
---

# Durable Workflow

This preset creates a durable Lambda function: AWS checkpoints the function's progress so a long-running workflow — here, a week-long order-fulfillment saga — survives interruption and resumes where it stopped. The classic 15-minute invocation cap stops applying; the workflow's end-to-end budget is set by `executionTimeoutSeconds` (up to 366 days).

## When to Use

- Multi-step business workflows that wait on external systems (payments, shipping, human approval)
- Long-running orchestration you would otherwise move to Step Functions but want to keep in code
- Retry-heavy integrations where progress must never restart from zero
- Any invocation whose real-world duration exceeds the classic 15-minute cap

## Key Configuration Choices

- **`durableConfig`** — the durability contract: a 7-day execution budget with 30 days of retained execution state and history. Adding or removing this block REPLACES the function (an AWS constraint); the values inside update in place.
- **Published head pointer** — `publish: true` + `publishTo: LATEST_PUBLISHED` keep an addressable `$LATEST.PUBLISHED` qualifier on the newest published version, so in-flight durable executions keep their code while new ones pick up releases.
- **Failure destination** — exhausted asynchronous invocations deliver their full invocation record to a queue for inspection and re-drive.
- **Generous timeout** — `timeoutSeconds: 900` sizes each RESUMED segment, not the whole workflow; the workflow budget lives in `durableConfig`.

## Placeholders to Replace

| Placeholder | Description | Example |
|-------------|-------------|---------|
| `<aws-region>` | Region for the function | `us-west-2` |
| `<failure-queue-arn>` | SQS queue receiving failed invocation records | `arn:aws:sqs:us-west-2:123456789012:workflow-failures` |

The `valueFrom` references assume an `order-workflow-exec` role and a `build-artifacts` bucket managed as Planton resources; replace with your own names or literal values.

## Common Additions

- Add `aliases` with weighted routing to canary new workflow versions
- Add `scalingConfigs` (with a Managed Instances capacity provider) to pin execution-environment capacity for bursty workflow starts
- Add `environment` variables for downstream service endpoints
- Raise `retentionPeriodDays` (up to 90) when post-mortems need older execution history

## Related Presets

- **01-zip-basic** — use for standard request/response functions without durability
- **02-container-basic** — use when the workflow's dependency tree needs a container image
