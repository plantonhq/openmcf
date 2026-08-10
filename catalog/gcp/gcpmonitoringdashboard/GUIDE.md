# GcpMonitoringDashboard Guide

Operational judgment for running Cloud Monitoring dashboards as code —
the things the spec reference cannot tell you.

## Build in the console, commit the export

Nobody writes dashboard JSON from scratch, and nobody should. The
workflow that works: build or edit the dashboard visually in the GCP
console, open dashboard settings → **JSON editor**, and paste the export
into `dashboardJson` verbatim. The provider suppresses server-assigned
keys (etag, name) on read-back, so a console export round-trips with a
clean plan. Treat the console as the editor and the manifest as the
source of truth — once a dashboard is managed here, console edits are
drift that the next apply reverts.

## The diff-suppression trade-off

Because the provider cannot know which JSON keys are server defaults, it
suppresses any key that exists in GCP's response but not in your
document. The consequence: a REMOVE-only change (deleting one widget and
nothing else) can be suppressed as no-diff. Pair every removal with any
trivial change (a title tweak) if a plan shows nothing — or apply twice.
This is the provider's own documented behavior, not a module choice.

## One dashboard per manifest, small dashboards over giant ones

The Dashboard API caps widgets per layout, and giant dashboards are
unreadable at incident time anyway. The presets model the honest split:
one golden-signals page per service, one infrastructure page per fleet.
Compose MANY small dashboards (they are free) instead of one
everything-page.

## Log-based metrics chart like any other metric

A GcpLogMetric named `checkout/errors` charts on a widget with
`metric.type="logging.googleapis.com/user/checkout/errors"` — the
dashboard + log-metric pair is how log-only signals earn a place on the
health page.

## Teardown discipline

Dashboards hold no data — they are views. `DELETE` is always safe (the
metrics keep flowing; only the page disappears). Reserve `PREVENT` for
the incident-response dashboard your runbooks deep-link, where a broken
link at 3am costs real minutes.
