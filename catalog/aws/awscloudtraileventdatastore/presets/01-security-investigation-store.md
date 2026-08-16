# Security Investigation Store

This preset creates the investigation workhorse: a multi-region (the
AWS default) Lake store ingesting every management event, queryable
with SQL for the full 7-year default window, termination-protected.

## When to Use

- "Who changed this?" investigations without shipping logs to a SIEM
- Long-horizon audit queries across the whole account

## What You Get

- Every management API call, queryable with Lake SQL in the CloudTrail
  console or `aws cloudtrail start-query`
- The 2555-day (7-year) default retention on extendable pricing
- Termination protection on (the AWS default) — deletion needs a
  deliberate two-step

## Customize

- Narrow `advancedEventSelectors` to cut the per-GB ingestion bill —
  ingestion is what Lake charges for
- Set `retentionPeriodDays` down (7 minimum) for short-horizon stores
- Add `kmsKeyId` for SSE-KMS (fixed at creation; the key policy must
  allow CloudTrail)
