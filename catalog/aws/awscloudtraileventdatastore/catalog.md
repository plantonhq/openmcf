# AWS CloudTrail Event Data Store

CloudTrail Lake: a queryable, immutable database of AWS activity
events you interrogate with SQL — security investigations and audits
without shipping logs to S3 or standing up a SIEM.

> **Availability:** AWS has closed CloudTrail Lake to new customers — `CreateEventDataStore` is rejected on any account that has never created an event data store ("CloudTrail Lake is no longer accepting new customers"). Deploy this component only on an account already using Lake; on other accounts, use AWS CloudTrail's trail delivery instead.

## What Gets Managed

- The event data store: pricing mode, queryable retention window
  (7–2555 days), multi-region and organization ingestion.
- Ingestion scope via advanced event selectors (every selector needs
  an `eventCategory` condition).
- Termination protection and the ingestion pause switch.
- Optional SSE-KMS encryption (fixed at creation).

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with CloudTrail permissions.

### AWS Account

- Nothing else — Lake owns its own storage. For SSE-KMS, a key whose
  policy allows CloudTrail
  ([AWS KMS Key](/cloud-catalog/aws-kms-key)); AWS recommends
  multi-region keys for multi-region stores.

## Deploy

### Console

Create the resource from the AWS catalog, scope the selectors (Lake
bills per GB ingested), pick the retention window, and deploy.

### CLI

```bash
planton apply -f event-data-store.yaml
```

## After Deploy

- Query events in the CloudTrail console's Lake editor or
  `aws cloudtrail start-query` with Lake SQL.
- Pause ingestion any time with `suspend: true` — stored events stay
  queryable and retention keeps counting.
- To destroy: set `termination_protection_enabled: false`, apply, then
  destroy. AWS holds the store in `PENDING_DELETION` for 7 days (the
  name stays reserved).

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
