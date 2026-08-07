# AwsCloudwatchLogGroup — Research Notes

## What AWS models

A CloudWatch Logs log group is the container and control point for log
streams: retention, encryption, class (pricing/feature tier), and deletion
protection live on the group. Around it, AWS models a family of
group-scoped satellite resources, each keyed by the group's name:

- **Metric filters** (`PutMetricFilter`, many per group) — extract CloudWatch
  metrics from matching log events. The transformation carries the metric
  name/namespace/value, an optional default value for quiet periods (which
  AWS forbids combining with dimensions), up to 3 dimensions, and a
  StandardUnit. This is the log-derived-alerting primitive.
- **Subscription filters** (`PutSubscriptionFilter`, max 2 per group) —
  real-time delivery of matching events to a Kinesis stream, Firehose
  delivery stream, Lambda function, or cross-account logs destination.
  Kinesis/Firehose destinations need an IAM role trusting
  `logs.amazonaws.com`; Lambda destinations authorize via a Lambda
  permission instead. `distribution` chooses per-stream ordering vs
  throughput for Kinesis; `emit_system_fields` can enrich events with the
  source account/region for centralized destinations.
- **Data protection policy** (one per group) — a JSON policy with paired
  Audit + Deidentify statements over managed/custom data identifiers; masks
  PII at ingestion (visible only with `logs:Unmask`).
- **Field index policy** (one per group) — a `{"Fields": [...]}` document
  (max 20 fields) that accelerates and cheapens Logs Insights queries
  filtering on the indexed fields.

Log group classes: `STANDARD` (full features), `INFREQUENT_ACCESS` (~50%
cheaper storage; NO metric/subscription filters or Contributor Insights),
`DELIVERY` (AWS-managed retention for vended service logs). Class is
create-time (ForceNew). Retention accepts a fixed value set only; the KMS
key association updates in place and requires the key policy to admit
`logs.<region>.amazonaws.com`.

## Fold vs. split decisions

All four satellites are FOLDED into this spec: they are strictly
group-scoped, share the group's lifecycle, and nothing else references them
— the many-per-parent ones (filters) materialize as one provider resource
per named entry in both engines. Deliberately NOT modeled here:

- **Log streams** — created by agents/services at write time (data plane).
- **Transformer** — a ~20-processor pipeline surface; deferred on demand.
- **Log anomaly detector** — spans multiple groups (standalone resource).
- **Delivery source/destination/delivery** — the vended-log delivery plane.
- **Cross-account destinations + policies** — aggregation plumbing;
  subscription filters compose to them by literal ARN.
- **Query definitions, account policies** — account-scoped tooling.
- **`skip_destroy`** — contradicts honest lifecycle management;
  `deletion_protection_enabled` serves the protection need.
- **`name_prefix`** — naming basis is `metadata.name`.

## Operational notes

- The provider strips the API's `:*` ARN suffix on read; consumers needing
  the wildcard form (Step Functions logging) append it themselves.
- Deletion protection is applied via a separate
  `PutLogGroupDeletionProtection` call; destroy fails while it is set.
- Metric filters publishing dimensions create one custom metric per unique
  dimension-value combination — keep cardinality low.
