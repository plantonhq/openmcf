# Observability Foundation

The day-2 monitoring pack for an Azure subscription. One deploy produces
the telemetry destination every other composition routes into -- a Log
Analytics workspace with its Application Insights face -- plus the alert
pack teams postpone until the first outage: service-health paging,
failed-admin-operation auditing, application error and exception alerts,
and an optional outside-in availability probe. Deploy it once per
environment, before the workloads it will watch.

## The architecture

Telemetry flows in, alerts flow out, and the routing is explicit:

- **The workspace** (`AzureLogAnalyticsWorkspace`) is the destination.
  Diagnostic settings on any resource, Container Insights on any
  cluster, and the Application Insights beside it all store here -- so
  one KQL query joins application symptoms to platform events.
- **Application Insights** (workspace-based, the only mode Azure still
  supports) is where instrumented applications report requests, traces,
  and exceptions; the query alerts watch the tables it feeds.
- **Two action groups** separate urgency from noise: the operations
  group receives routine signals (failed admin operations, error
  spikes); the critical group receives pages (service-health incidents,
  availability failures). They are separate resources so a
  maintenance-window disable of routine noise can never silence a page.
- **Activity-log alerts** watch the subscription's control-plane record
  -- the events that never appear as metrics. Service-health incidents
  and security advisories page; failed administrative operations (the
  early signature of expired credentials or broken automation) notify
  operations.
- **Scheduled query alerts** run KQL against the workspace on a
  schedule: failed-request spikes every 5 minutes, exception surges over
  15 -- both auto-resolving when the condition clears.
- **The availability probe** (optional) requests a public URL from five
  continents and pages when several locations fail at once -- the only
  signal in the pack measured from where users actually stand.

## What is on by default

- **Service health pages** (`critical_email`): incidents and security
  advisories affecting the subscription go to the critical group.
  Planned-maintenance events stay in the portal.
- **Failed admin operations notify** (`ops_email`): Error/Critical
  control-plane failures across the whole subscription. Succeeded
  changes are history, not alerts.
- **Error and exception alerts auto-resolve** (severity 2, warning):
  they fire into the operations group when thresholds trip and close
  themselves when the condition clears.
- **Query validation is skipped at create time** on the query alerts:
  the `AppRequests`/`AppExceptions` tables materialize only once the
  first application reports -- the rules are correct from the moment
  data flows.
- **Uncapped ingestion** (`daily_quota_gb: -1`): a cap is a cost guard
  that is also a data-loss dial; set one deliberately or not at all.
- **The availability probe is off** (`web_test_enabled`): it needs a
  real public URL to watch -- enable it together with `web_test_url`.

## Parameters worth understanding

- **`subscription_id`** is the one parameter that MUST change: the
  activity-log alerts watch this subscription's control-plane record,
  and service-health events exist only at subscription scope. The
  placeholder default deploys, but watches nothing useful.
- **`ops_email` / `critical_email`**: routing is the whole game --
  operations goes to a team, critical goes to an on-call intake. Azure
  sends each address a confirmation mail on first deploy.
- **`error_spike_threshold`**: shared by the failed-request and
  exception alerts. Tune it to the estate's traffic; the goal is "worth
  a look", not "every blip".
- **`web_test_locations` / `web_test_failed_location_threshold`**: five
  locations with a threshold of three tolerates two simultaneous
  probe-side blips while still catching partial outages. Location IDs
  are Azure web-test IDs (e.g. `us-va-ash-azr`), not region names.

## After deployment

Everything in this chart provisions in under five minutes.

- **Confirm the notification channels**: Azure emails a confirmation to
  each receiver address on first deploy -- make sure it did not land in
  spam.
- **Route telemetry in**: point workloads' diagnostic settings
  (`AzureMonitorDiagnosticSetting`) at the workspace's `workspace_id`
  output, and wire applications to the Application Insights
  `connection_string` output. The alert pack starts watching the moment
  data arrives.
- **Test the pipeline honestly**: fire a synthetic failure (a burst of
  failed requests against a test app) and watch it travel query → alert
  → action group → inbox.
- **Natural next steps**: add per-service `AzureMonitorMetricAlert`s
  against the platform metrics that matter to each workload
  (CPU, queue depth, RU throttling); add `AzureMonitorActionGroup`
  webhook receivers when an incident tool arrives; and once workloads
  multiply, split per-team action groups -- alert rules are the volatile
  edge, groups are the stable routing nodes.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
