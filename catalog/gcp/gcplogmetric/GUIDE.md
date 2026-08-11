# GcpLogMetric Guide

Operational judgment for running log-based metrics as code — the things
the spec reference cannot tell you.

## Metrics count from creation, never backwards

A log-based metric ingests entries that arrive AFTER it exists. It will
never answer "how many errors last month" for a metric created today —
that question needs Log Analytics (SQL over a GcpLogBucket) or BigQuery.
Create the metrics you will want to alert on WITH the service, not after
the first incident.

## The label schema is append-only — plan replacement

Changing a label's key or value type REPLACES the metric, and the
history does not follow. Adding a NEW label is safe. Design the label
set small and stable up front (status class, method — not user IDs),
and treat a schema change as a new metric with a migration window where
both exist.

## High-cardinality labels are a cost incident waiting

Every distinct label-value combination is a separate time series that
Monitoring stores and bills. `EXTRACT(httpRequest.status)` is bounded
(~60 values); `EXTRACT(jsonPayload.user_id)` is unbounded and will
melt the bill. If a value is unbounded, it belongs in the log entry,
not on a metric label.

## Distribution metrics need all three parts

`valueType: DISTRIBUTION` alone measures nothing: it needs
`valueExtractor` (where the number comes from) and `bucketOptions` (the
histogram it lands in). The API enforces the pairing server-side at
apply. Exponential buckets fit latencies (each bucket a constant factor
wider); linear fits bounded gauges; explicit is for when you already
know the SLO boundaries you care about.

## Filter cost discipline

The filter runs against every entry in scope as it arrives. Lead with
the cheapest selective term — `resource.type` — before free-text or
regex terms, and scope to the service (`resource.labels.service_name`)
so a noisy neighbor's logs never feed your metric.

## disabled is the pause button, DELETE is not

A misfiring extractor (wrong regex, NaN floods) is silenced with
`disabled: true` — configuration and history survive while you fix it.
Deleting and recreating the metric loses history and breaks every alert
policy and dashboard that references the name.

## Teardown discipline

`DELETE` removes the metric; alert policies watching
`logging.googleapis.com/user/{name}` silently stop evaluating (they do
not fail loudly). `PREVENT` is the honest posture once alerts depend on
the metric — the same reasoning as the SLO kind.
