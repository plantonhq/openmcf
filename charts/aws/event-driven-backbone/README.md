# AWS Event-Driven Backbone

The event backbone decoupled services are built on: producers publish
domain events to a custom EventBridge bus and stop caring who consumes
them; a routing rule delivers every matched event to a durable SQS work
queue with a locked-down dead-letter queue behind it; optional SNS
fan-out gives additional consumer teams their own copies; an optional
Lambda consumer turns the work queue into running code. The one alarm
that matters — dead letters present — is wired to email.

Publishing an event is one `PutEvents` call; nothing else about the
producer changes as consumers come and go. That inversion is the whole
point of the chart.

## Architecture

```
 producers ──▶ AwsEventBridgeBus (custom bus, the domain's namespace)
                 │
        AwsEventBridgeRule (source prefix match, e.g. com.mycompany*)
                 │
    ┌────────────┴──────────────────────────┐
    ▼ (always)                              ▼ (sns_fanout_enabled)
 AwsSqsQueue work-queue                  AwsSnsTopic fanout
    │ policy: this rule only                │ policy: this rule only
    │ redrive: 5 strikes                    ▼
    │                                    AwsSnsSubscription (raw delivery)
    ▼ (lambda_consumer_enabled)             │
 AwsLambdaEventSourceMapping             AwsSqsQueue fanout-queue
    │ batch 10, partial-batch failures      │ policy: this topic only
 AwsLambda consumer ── AwsIamRole           │ redrive: 5 strikes
                                            ▼
              AwsSqsQueue dlq  ◀── every failure path lands here
                 │ 14-day retention, byQueue redrive allowlist
 AwsCloudwatchAlarm dlq-depth ──▶ AwsSnsTopic alerts ──▶ email (alarms_enabled)
```

Delivery permissions are resource policies composed at render time from
the chart's own names — every grant names exactly this chart's rule or
topic, never "any EventBridge rule in the account".

## Included Cloud Resources

| Resource | Kind | Purpose |
|----------|------|---------|
| Event bus | `AwsEventBridgeBus` | The domain's own namespace — producers publish here |
| Routing rule | `AwsEventBridgeRule` | Prefix-matches `source` and delivers to the planes |
| Work queue | `AwsSqsQueue` | Durable buffer the primary consumer group owns |
| Dead-letter queue | `AwsSqsQueue` | 14-day evidence locker for every failure path |
| Fan-out topic | `AwsSnsTopic` | Copies each event to every subscriber queue (conditional) |
| Fan-out queue + subscription | `AwsSqsQueue` + `AwsSnsSubscription` | One consumer team's copy, raw delivery (conditional) |
| Consumer role | `AwsIamRole` | Logs + receive/delete on the work queue only (conditional) |
| Consumer function | `AwsLambda` | Your handler, batch-invoked from the queue (conditional) |
| Event source mapping | `AwsLambdaEventSourceMapping` | The queue→function poller, partial-batch aware (conditional) |
| Alerts topic + subscription | `AwsSnsTopic` + `AwsSnsSubscription` | Alarm delivery to email (conditional) |
| DLQ depth alarm | `AwsCloudwatchAlarm` | Pages the moment any dead letter appears (conditional) |

## Parameters

| Name | Description | Default | Required |
|------|-------------|---------|----------|
| `aws_region` | Region for every resource | `us-east-1` | yes |
| `aws_account_id` | Account id used in render-time policy ARNs | `123456789012` | yes |
| `backbone_name` | Name prefix for every resource | `events` | yes |
| `event_source_prefix` | `source` prefix the rule matches | `com.mycompany` | yes |
| `sns_fanout_enabled` | Add the SNS fan-out plane | `false` | no |
| `lambda_consumer_enabled` | Deploy the Lambda consumer on the work queue | `false` | no |
| `consumer_code_bucket` | S3 bucket with the consumer zip | `my-lambda-artifacts` | consumer on |
| `consumer_code_key` | Object key of the consumer zip | `events/consumer.zip` | consumer on |
| `consumer_runtime` | Consumer runtime | `python3.12` | consumer on |
| `consumer_handler` | Consumer entry point | `consumer.handler` | consumer on |
| `consumer_timeout_seconds` | Per-invocation ceiling (keep ≤ 50 with the shipped queue) | `30` | no |
| `alarms_enabled` | DLQ-depth alarm wired to email | `true` | no |
| `alert_email` | Alert destination (confirm AWS's first email) | `ops@example.com` | when alarms on |

## First events

1. Deploy the chart (set `aws_region`, `aws_account_id`,
   `backbone_name`, and your `event_source_prefix`). If
   `alarms_enabled` is on, click the confirmation link AWS emails to
   `alert_email`.

2. Publish a test event to the bus:

   ```bash
   aws events put-events --entries '[{
     "EventBusName": "events-bus",
     "Source": "com.mycompany.orders",
     "DetailType": "OrderPlaced",
     "Detail": "{\"orderId\": \"o-123\", \"total\": 42}"
   }]'
   ```

3. Watch it arrive on the work queue:

   ```bash
   aws sqs receive-message --queue-url \
     "$(aws sqs get-queue-url --queue-name events-work-queue --output text)"
   ```

The event arrives in EventBridge's envelope (`source`, `detail-type`,
`detail`, ...). With fan-out enabled, the same event also lands on the
fan-out queue — as the bare envelope too, because the subscription uses
raw delivery.

## The consumer contract

With `lambda_consumer_enabled`, batches of up to 10 messages invoke your
handler. The mapping enables partial-batch failure reporting: return the
ids that failed and only those retry — a handler that returns nothing
reports whole-batch success. The minimal shape:

```python
import json

def handler(event, context):
    failures = []
    for record in event["Records"]:
        try:
            envelope = json.loads(record["body"])  # the EventBridge envelope
            process(envelope["detail"])
        except Exception:
            failures.append({"itemIdentifier": record["messageId"]})
    return {"batchItemFailures": failures}
```

After 5 failed receives, a message moves to the DLQ and the alarm pages.

## Day-2 guidance

- **Redriving dead letters.** After fixing the consumer, move messages
  back with SQS's DLQ redrive (console, or
  `aws sqs start-message-move-task --source-arn <dlq-arn>`). The DLQ's
  `redriveAllowPolicy` already permits exactly this chart's queues as
  destinations.
- **Adding a fan-out consumer.** Copy the fan-out queue + subscription
  pair in `templates/fanout.yaml` under a new name (or create the same
  two resources directly), and add the new queue's ARN to the DLQ's
  `redriveAllowPolicy.sourceQueueArns`. Each team gets its own queue —
  never share one queue between two consumer groups expecting copies.
- **Sharper routing.** Add sibling `AwsEventBridgeRule` resources on the
  deployed bus with narrower patterns (`source` equals, `detail-type`
  matches) and their own targets. Remember every new SQS target needs
  the matching `AllowRuleDelivery` statement in that queue's policy —
  delivery without it is silently dropped.
- **Consumers outside Lambda.** A container service polls the work queue
  with the same three IAM actions the consumer role grants
  (`sqs:ReceiveMessage`, `sqs:DeleteMessage`, `sqs:GetQueueAttributes`)
  on the queue's ARN. The backbone does not care what runs the code.
- **Event archive and replay.** EventBridge can archive matched events
  on the bus and replay them into a rule after an incident — a natural
  extension on the deployed bus once the backbone carries traffic worth
  replaying.
- **Cross-account producers.** Grant a producer account `PutEvents` on
  the bus by setting the bus's `resourcePolicy` — the backbone then
  spans team boundaries with the same delivery guarantees.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
