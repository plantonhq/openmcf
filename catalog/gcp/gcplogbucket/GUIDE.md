# GcpLogBucket Guide

Operational judgment for running Cloud Logging buckets as code — the
things the spec reference cannot tell you.

## Storage and routing are different kinds on purpose

A bucket stores; a sink routes. GCP models them separately and so does
the catalog: this kind never creates the GcpLoggingSink that feeds it.
The composition is one chart — the sink's `rawUri` destination takes
this kind's `bucket_name` output as
`logging.googleapis.com/{bucket_name}` — and the split is what lets one
bucket receive from many sinks (and one sink fan into a bucket another
team owns).

## Adopt-on-match is the feature, not a surprise

Creating a bucket whose `bucketId` already exists ADOPTS it. That is the
designed path to managing `_Default` (every project has one, 30-day
retention, unlocked) — adopt it and set retention deliberately. Two
consequences to respect: a typo'd `bucketId` that collides with a
stranger's bucket adopts THAT (names are scope-unique, so review plans
on shared projects), and destroying an adopted `_Default`/`_Required`
only un-manages it — they are undeletable by API decree.

## The three one-way doors deserve a change-review each

- `locked: true` freezes retention FOREVER and blocks deletion until
  every entry ages out. Locking a 3650-day bucket is a ten-year
  commitment no one can undo — arm it in a reviewed change, never a
  quick edit.
- `enableAnalytics: true` never turns off. It is also the gate for
  `linkedBigqueryDataset` — plan both together.
- `cmekKmsKey` never reverts to Google-managed encryption. Key ROTATION
  is allowed; disabling CMEK is not. And grant the Logging service
  account on the key BEFORE the apply, or creation fails mid-plan.

## Views are the access-control story

Never grant `roles/logging.viewer` on a compliance bucket when a team
needs one slice: add a `logViews` entry with the filter and grant
`roles/logging.viewAccessor` on THAT view. Views update in place for
filter/description but REPLACE on rename — treat `viewId` as permanent.

## Retention changes are cheap; shortening is deletion

Raising `retentionDays` is free. LOWERING it schedules deletion of
entries older than the new window — GCP applies the shorter retention
to existing data. Treat a retention decrease as a data-deletion change
in review.

## Teardown discipline

`DELETE` on a custom bucket deletes its stored entries — there is no
recycle bin. The linked BigQuery dataset and views go with it (the
kind's deletion policy fans out). `PREVENT` is the honest posture for
any bucket a compliance clause names.
