---
title: "Compliance Retention"
description: "A 400-day audit bucket with a status-code index and a PREVENT destroy posture — the storage half of a compliance logging pipeline (pair it with a GcpLoggingSink routing audit entries in)."
type: "preset"
rank: "01"
presetSlug: "01-compliance-retention"
componentSlug: "log-bucket"
componentTitle: "Log Bucket"
provider: "gcp"
icon: "package"
order: 1
---

# Compliance Retention

A 400-day audit bucket with a status-code index and a PREVENT destroy
posture — the storage half of a compliance logging pipeline (pair it
with a GcpLoggingSink routing audit entries in).

## What it configures

- `retentionDays: 400` — thirteen months, the common "one year plus
  slack" audit window.
- An integer index on `jsonPayload.request.status` for fast status-code
  queries without full scans.
- `deletionPolicy: PREVENT` — a bucket a compliance clause names must
  not die in a fat-fingered destroy.

## Adjust before deploying

- **retentionDays** — set from the actual policy document; remember
  LOWERING it later deletes entries older than the new window.
- **locked** — the preset deliberately leaves the bucket unlocked.
  Locking freezes retention forever and blocks deletion until every
  entry ages out — arm it only in a reviewed change once the retention
  number is final.

## When to choose something else

For SQL over your logs and BigQuery access rather than long-term
archive, start from the **Analytics And Views** preset.
