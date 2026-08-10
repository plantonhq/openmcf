# GcpMonitoringSlo Guide

Operational judgment for running SLOs as code — the things the spec
reference cannot tell you.

## An SLO without a burn-rate alert is a report card, not a guardrail

The SLO itself pages nobody. Pair every SLO with a
GcpMonitoringAlertPolicy whose threshold condition filters on
`select_slo_burn_rate("{slo_name}", "3600s")` (the `slo_name` output is
exactly the handle that filter needs). The classic pairing is two
alerts: a fast-burn page (14.4× over 1h — the budget dies in a day) and
a slow-burn ticket (1× over 3 days).

## goal caps at 0.9999 by API decree

GCP refuses five nines and beyond — not a provider quirk, the API's own
bound. If the business asks for 99.999%, the honest answer is that GCP
cannot even EXPRESS that target; negotiate the SLO down or measure it
outside Cloud Monitoring.

## The service binding is for life

`service`, `slo_id`, and the created services' `service_id` are all
replace-on-change: renaming a service REPLACES the SLO under it, and the
error-budget history dies with it. Name services for what they ARE
(checkout, search), never for teams or quarters — history outlives both.

## basic_sli only works where GCP understands the telemetry

`basicSli` (availability/latency with no filters) needs a service type
whose telemetry GCP natively maps — App Engine, Cloud Endpoints, Istio.
On a `customService` the API accepts the SLO and measures NOTHING
meaningful. Custom services take `requestBasedSli`, whose filters say
exactly what good means.

## Rolling windows for engineering, calendar for contracts

`rollingPeriodDays: 30` answers "how are we doing lately" and is what
burn-rate alerting assumes. `calendarPeriod: MONTH` matches how customer
contracts and SLAs are written, but the budget reset at month boundary
makes day-1 incidents look artificially cheap. Pick per audience; run
both for the same service when engineering and contract views disagree.

## The good/total derivation is exactly-two on purpose

`goodTotalRatio` takes exactly two of good/bad/total and derives the
third. Prefer good+total (bad-counting misses timeouts that never wrote
a bad event); use bad+total only when the good signal does not exist.

## Teardown discipline

Deleting an SLO deletes its error-budget history (the underlying metrics
survive). `PREVENT` is the right posture once burn-rate alerts reference
the SLO — a deleted SLO silently breaks every alert built on it. The
kind's deletion policy also covers any service it created; a service
consumed by OTHER SLOs refuses deletion server-side until they go.
