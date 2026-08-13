# Data Warehouse Events

This preset runs a Redshift Data API statement whenever matching events arrive — here, calling an ingestion procedure each time new data lands in an S3 bucket. It demonstrates the event-driven ELT pattern: no poller, no scheduler drift, no idle compute — the warehouse reacts to data the moment it exists.

## When to Use

- Loading or transforming data in Redshift as soon as it lands (event-driven ELT)
- Refreshing materialized views when upstream data changes
- Triggering warehouse maintenance from operational events instead of fixed schedules
- Any pattern where SQL must run in response to an event, not on a timer

## Key Configuration Choices

- **Redshift target** — the target `arn` is the CLUSTER; the statement, database, and credentials live in `redshiftTarget`.
- **Secrets Manager auth** — the statement authenticates with credentials from Secrets Manager (`secretsManagerArn`). For temporary-credential auth, use `dbUser` instead.
- **`withEvent: true`** — the matched event is delivered to the statement as execution context, so the procedure can read what triggered it.
- **Invocation role** — EventBridge assumes `roleArn` to call the Redshift Data API; it needs `redshift-data:ExecuteStatement` plus read access to the credentials secret.
- **Retry policy + DLQ** — failed deliveries retry for up to an hour, then land in the dead letter queue instead of disappearing.

## Placeholders to Replace

| Placeholder | Description | Example |
|-------------|-------------|---------|
| `<landing-bucket-name>` | S3 bucket whose Object Created events trigger ingestion | `acme-landing-zone` |
| `<redshift-cluster-arn>` | ARN of the Redshift cluster | `arn:aws:redshift:us-east-1:123456789012:cluster:analytics` |
| `<events-redshift-role-arn>` | Role EventBridge assumes to call the Data API | `arn:aws:iam::123456789012:role/events-redshift-data` |
| `<redshift-credentials-secret-arn>` | Secrets Manager secret with the database credentials | `arn:aws:secretsmanager:us-east-1:123456789012:secret:redshift-etl` |
| `<dlq-queue-arn>` | ARN of the SQS dead letter queue | `arn:aws:sqs:us-east-1:123456789012:warehouse-events-dlq` |

## Common Additions

- Switch to `scheduleExpression` for fixed-cadence warehouse maintenance instead of event-driven runs
- Use `dbUser` for temporary-credential auth when the cluster allows IAM-mapped database users
- Add an `inputTransformer` to pass extracted event fields into the statement's context
- Fan out with additional targets (e.g., an SQS audit queue alongside the warehouse call)

## Related Presets

- **02-event-pattern-sqs** — use when events should queue for application-side processing instead of running SQL
- **03-multi-target-fanout** — use when one event must reach several services at once
