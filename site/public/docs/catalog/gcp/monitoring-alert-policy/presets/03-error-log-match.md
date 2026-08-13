---
title: "Error Log Match"
description: "Pages on matching log entries — the application-truth tier of alerting, where a logged panic pages before any metric moves."
type: "preset"
rank: "03"
presetSlug: "03-error-log-match"
componentSlug: "monitoring-alert-policy"
componentTitle: "Monitoring Alert Policy"
provider: "gcp"
icon: "package"
order: 3
---

# Error Log Match

Pages on matching log entries — the application-truth tier of alerting,
where a logged panic pages before any metric moves.

## What it configures

- A `conditionMatchedLog` on ERROR-severity entries containing "panic",
  extracting the instance id into the incident as a label.
- The MANDATORY `notificationRateLimit` (one notification per 5 minutes)
  — the GCP API rejects log-based policies without it, and the limit is
  what keeps a log storm from becoming a page storm.

## Adjust before deploying

- **filter** — match your application's actual failure signature
  (`jsonPayload.level="fatal"`, a resource.type scope). Broad filters
  page on noise.
- **labelExtractors** — pull the identifiers your runbook needs straight
  into the notification.

## When to choose something else

If the signal already exists as a metric (error-rate SLIs), a threshold
condition evaluates continuously and supports ratios/forecasting; keep
log-match for signals that only exist in logs.
