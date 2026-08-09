# GcpMonitoringAlertPolicy Guide

The judgment this guide protects: an alert policy is code that wakes
humans. Every default here trades paging noise against missed outages,
and the honest posture is explicit: sustained durations, deliberate
missing-data behavior, and runbooks in the notification itself.

## Duration is the flap guard

`duration: 0s` pages on a single violating point; `300s` requires five
sustained minutes. Almost every noisy policy is a `0s` policy. The GCP
API additionally requires whole minutes — a `90s` duration fails at
apply, not at validation.

## Aggregation is not optional in practice

Without an `aggregations` entry, GCP compares every raw point of every
series — high-cardinality metrics then page on any single instance's
blip. `ALIGN_MEAN` over a 60s window is the standard floor;
`REDUCE_MEAN` across series turns "any instance" into "the fleet".

## Missing data is a decision

`evaluationMissingData` defaults to no-op: a metric that STOPS reporting
keeps whatever incident state it had. `_ACTIVE` (missing data violates)
is the paranoid setting for heartbeat-like metrics; `_INACTIVE` closes
incidents on silence. For true "silence is failure" semantics, prefer a
dedicated `conditionAbsent` — it exists for exactly that.

## Log-based policies carry a mandatory rate limit

`conditionMatchedLog` without `alertStrategy.notificationRateLimit`
fails at the API — and with a tight limit, a log storm becomes one page
per period instead of hundreds. This pairing is enforced by the API, not
the provider, so it surfaces at apply time; the spec documents it where
the field lives.

## MQL is a dead end

Google deprecated MQL in favor of PromQL. The arm stays modeled because
the API still serves it, but new policies belong on
`conditionPrometheusQueryLanguage` — and ported Prometheus rules carry
their `for:` as `duration` and their labels verbatim.

## Teardown discipline

`PREVENT` is the right `deletionPolicy` for policies that page
production — destroying one is equivalent to deleting the monitoring it
provides. `ABANDON` keeps it evaluating (and paging) unmanaged, which is
usually worse than either alternative: prefer disabling (`enabled:
false`) while deciding.
