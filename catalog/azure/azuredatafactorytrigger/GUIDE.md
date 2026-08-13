# Azure Data Factory Trigger -- Operational Guide

Judgment calls that matter when you run Data Factory triggers in production.

## A started trigger is a running system

`activated: true` (the default) means the trigger starts evaluating the moment it deploys -- a schedule trigger with an omitted `start_time` starts from NOW, and every firing is a real, billed pipeline run. Deploy new triggers with an explicit future `start_time` (or `activated: false`), verify the target pipeline's Debug run, then let it go live. The trigger's state survives your intent: it is started or stopped in Azure, not in your manifest history.

One live-proven bound on "explicit future start_time": ARM refuses to START a schedule trigger whose next execution is more than 18 months out (`InvalidWorkflowTriggerRecurrenceSchedule`). Park a trigger further into the future than that with `activated: false`, not with a distant date -- the distant date deploys but can never be started until it comes within the window. Tumbling window triggers carry no such bound (a far-future window start is accepted).

## Custom event triggers require an Event-Grid-schema topic

When a custom event trigger STARTS, Data Factory creates a webhook event subscription on your Event Grid topic -- and topics using the CloudEvents schema validate webhook subscribers with an HTTP OPTIONS handshake Data Factory's endpoint does not answer. The failure is live-proven and reads as "Webhook endpoint validation failed ... MethodNotAllowed" at Start. Point custom event triggers only at topics whose `input_schema` is "EventGridSchema" (the default); a CloudEvents topic cannot feed Data Factory, whatever the portal lets you save.

## A trigger holds its pipelines hostage at delete time

ARM refuses to delete a pipeline any trigger still references ("The document cannot be deleted since it is referenced by ..."), including triggers in a failed half-started state. Delete triggers before their pipelines; when tearing down a whole factory, deleting the factory itself removes everything inside it in one motion and sidesteps the ordering.

## Updates bounce the trigger -- plan for the gap

Azure forbids editing a started trigger, so every update stops it, applies the change, and starts it again. Schedule firings due DURING that gap are skipped, not queued -- for a once-a-day trigger updated at the wrong moment, that is a missed day. Update triggers outside their firing windows, or reconcile the skipped window manually.

## Schedule for clocks, tumbling window for data

A schedule trigger says "run at 02:00"; a tumbling window trigger says "process 00:00-01:00, then 01:00-02:00, ..." and passes each window's bounds to the pipeline (`@trigger().outputs.windowStartTime`). If the pipeline's work is a function of a time range -- every incremental load -- use tumbling windows: you get backfill by pointing `start_time` at history, per-window retry, `max_concurrency` rate-limiting, and window dependencies. Use schedule triggers for genuinely clock-shaped work (report generation, cache warms).

## Backfills are deliberate, bounded operations

A tumbling window `start_time` in the past runs EVERY window between then and now. Before pointing at history: set `max_concurrency` to what the sink tables tolerate (the 50 default fans out fifty parallel runs), set the pipeline's own `concurrency` guard, and give the run `retry` so transient failures do not strand single windows in a sea of green ones.

## Blob event triggers have a second resource behind them

The blob event type creates an Event Grid subscription on the storage account behind the scenes (hence the Microsoft.EventGrid provider requirement). Two operational consequences: the storage account's event pipeline is now part of your trigger's failure domain, and path filters (`blob_path_begins_with`/`ends_with`) are your only volume guard -- an over-broad filter on a busy account fires the pipeline for every blob. Filter to the narrowest path that works, and use `ignore_empty_blobs` when upstream tools write zero-byte markers.

## Self-dependencies serialize windows

A tumbling window dependency entry with no `trigger_name` makes each window wait for the PREVIOUS window of the same trigger -- the lever for loads that must apply strictly in order. Azure requires a negative `offset` for self-dependencies (e.g. `-24:00:00` for daily windows). Without it, windows run in parallel up to `max_concurrency`, which is faster and fine for independent partitions.
